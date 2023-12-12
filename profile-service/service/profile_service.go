package service

import (
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iqbaludinm/hr-microservice/profile-service/config"
	"github.com/iqbaludinm/hr-microservice/profile-service/exception"
	"github.com/iqbaludinm/hr-microservice/profile-service/helper"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/web"
	"github.com/iqbaludinm/hr-microservice/profile-service/repository"
	"github.com/iqbaludinm/hr-microservice/profile-service/service/producers"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
)

type ProfileServiceImpl struct {
	ProfileRepository repository.ProfileRepository
	Validate          *validator.Validate
}

type ProfileService interface {
	UpdateMyProfile(ctx *fiber.Ctx, id string, request web.UpdateProfileRequest) (domain.User, error)
	ForgetPasswordEmail(ctx *fiber.Ctx, email string) (domain.ResetPasswordToken, error)
	ResetPassword(ctx *fiber.Ctx, email, token string, request web.ResetPassword) error
}

type profileService struct {
	profileRepository    repository.ProfileRepository
	kafkaProducerService producers.KafkaProducerService
	logger               *zap.SugaredLogger
}

func NewProfileService(profileRepository repository.ProfileRepository, kafkaProducerService producers.KafkaProducerService, logger *zap.SugaredLogger) ProfileService {
	return &profileService{
		profileRepository:    profileRepository,
		kafkaProducerService: kafkaProducerService,
	}
}

func (service *profileService) UpdateMyProfile(ctx *fiber.Ctx, id string, request web.UpdateProfileRequest) (domain.User, error) {

	_, err := service.profileRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "id", id)
	if err != nil {
		return domain.User{}, exception.ErrNotFound("User not found.")
	}

	updatedUser := domain.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Phone:     request.Phone,
		UpdatedAt: time.Now(),
	}

	data, err := service.profileRepository.UpdateMyProfileTx(ctx.Context(), id, updatedUser)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "users_email_key") {
				return domain.User{}, exception.ErrBadRequest("Email already exist.")
			} else if strings.Contains(err.Error(), "users_phone_key") {
				return domain.User{}, exception.ErrBadRequest("Phone already exist.")
			}
		}
		return domain.User{}, err
	}

	updatedUser.ID = id
	// produce to kafka
	kafkaRegisterMessage := kafkamodel.NewKafkaUserMessage(updatedUser)
	go service.kafkaProducerService.Produce(kafkaRegisterMessage, "PUT.USER", config.KafkaTopic)

	return data, nil
}

// Password
func (service *profileService) ForgetPasswordEmail(ctx *fiber.Ctx, email string) (domain.ResetPasswordToken, error) {
	user, err := service.profileRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "email", email)

	if err != nil || user.Email == "" {
		return domain.ResetPasswordToken{}, fiber.NewError(fiber.StatusBadRequest, "Email tidak ditemukan.")
	}
	
	tokens, _ := service.profileRepository.CheckTokenWithQueryTx(ctx.Context(), "email", email)

	var data domain.ResetPasswordToken
	urlReset := helper.UrlReset

	if tokens.Id == "" {
		attempt := 1

		token := strings.ToLower(randstr.String(30))
		data.Id = uuid.New().String()
		data.Tokens = token
		data.Email = email
		data.Attempt = &attempt
		data.LastAttempt = time.Now()
		data.URL = urlReset + "?email=" + email + "&token=" + token

		err = service.profileRepository.AddTokenTx(ctx.Context(), data)
		if err != nil {
			return domain.ResetPasswordToken{}, err
		}

		// non-active
		// err = helper.EmailSender2(email, token)
		// if err != nil {
		// 	return nil, fiber.NewError(fiber.StatusBadGateway, "request error")
		// }
	} else {
		var attempt int
		// check in wib
		y1, m1, d1 := time.Now().Add(7 * time.Hour).Date()
		y2, m2, d2 := tokens.LastAttempt.Add(7 * time.Hour).Date()
		log.Println(time.Now(), tokens.LastAttempt, email)
		if !(y1 == y2 && m1 == m2 && d1 == d2) {
			attempt = 1
		} else {
			attempt = *tokens.Attempt + 1
		}

		if attempt > 3 {
			return domain.ResetPasswordToken{}, fiber.NewError(fiber.StatusTooManyRequests, "Terlalu banyak upaya. Coba lagi dalam 1 hari")
		}

		token := strings.ToLower(randstr.String(30))
		tokens.Tokens = token
		tokens.Email = email
		tokens.Attempt = &attempt
		tokens.LastAttempt = time.Now()
		tokens.URL = urlReset + "?email=" + email + "&token=" + token
		data = tokens

		err = service.profileRepository.UpdateTokenTx(ctx.Context(), tokens)
		if err != nil {
			return domain.ResetPasswordToken{}, err
		}

		// err = helper.EmailSender2(email, token)
		if err != nil {
			return domain.ResetPasswordToken{}, fiber.NewError(fiber.StatusBadGateway, "request error")
		}
	}

	return data, nil
}

func (service *profileService) ResetPassword(ctx *fiber.Ctx, email, token string, request web.ResetPassword) error {
	// var decodedByte, _ = base64.StdEncoding.DecodeString(token)
	// var resetToken = string(decodedByte)

	if request.Password != request.PasswordConfirm {
		panic(exception.ErrBadRequest("Password didn't match."))
	}

	checkToken, err := service.profileRepository.CheckTokenWithQueryTx(ctx.Context(), "tokens", token)

	if err != nil {
		return exception.ErrBadRequest("Token invalid.")
	}
	
	if checkToken.Tokens != token || checkToken.Email != email {
		return exception.ErrBadRequest("Token invalid.")
	}

	if time.Since(checkToken.LastAttempt) > 1*time.Hour {
		return exception.ErrBadRequest("Token sudah kadaluarsa.")
	}

	var user domain.User

	user, err = service.profileRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "email", email)
	if err != nil {
		return exception.ErrNotFound(err.Error())
	}

	if len(request.Password) < 6 {
		return exception.ErrBadRequest("Password length should more then equal 6 character.")
	}

	user.UpdatedAt = time.Now()
	user.SetPassword(request.Password)

	service.profileRepository.UpdatePasswordTx(ctx.Context(), user)

	// Update Token
	nol := 0
	checkToken.Attempt = &nol
	checkToken.Tokens = ""
	err = service.profileRepository.UpdateTokenTx(ctx.Context(), checkToken)
	if err != nil {
		if strings.Contains(err.Error(), "PasswordConfirm") && strings.Contains(err.Error(), "required") {
			return exception.ErrBadRequest("Password confirm required.")
		} else if strings.Contains(err.Error(), "Password") && strings.Contains(err.Error(), "required") {
			return exception.ErrBadRequest("Password required.")
		}
		return err
	}

	// produce to kafka
	kafkaRegisterMessage := kafkamodel.NewKafkaUserMessage(user)
	go service.kafkaProducerService.Produce(kafkaRegisterMessage, "POST.USER", config.KafkaTopic)
	go service.kafkaProducerService.Produce(kafkaRegisterMessage, "PUT.USER_PASS", config.KafkaTopic)

	return nil
}
