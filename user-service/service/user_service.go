package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iqbaludinm/hr-microservice/user-service/config"
	"github.com/iqbaludinm/hr-microservice/user-service/exception"
	"github.com/iqbaludinm/hr-microservice/user-service/helper"
	"github.com/iqbaludinm/hr-microservice/user-service/model/domain"
	"github.com/thanhpk/randstr"

	"github.com/iqbaludinm/hr-microservice/user-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/user-service/model/web"
	"github.com/iqbaludinm/hr-microservice/user-service/repository"
	"github.com/iqbaludinm/hr-microservice/user-service/service/producers"
	"go.uber.org/zap"
)

type UserService interface {
	// With Transaction
	// CreateUser(ctx context.Context, request web.CreateUserRequest) (web.UserResponse, error)
	// UpdateUser(ctx context.Context, id string, request web.UpdateUserRequest) (web.UserResponse, error)
	// UpdatePassword(ctx context.Context, request web.UpdateUserPasswordRequest) (web.UserResponse, error)
	// Delete(ctx context.Context, id string) (web.UserResponse, error)

	// Without Transaction
	FindAllUser(ctx context.Context, filter web.UserQueryFilter) (result []web.UserResponse, totalData int, err error)
	ForgetPasswordEmail(ctx *fiber.Ctx, email string) (domain.ResetPasswordToken, error)
	ResetPassword(ctx *fiber.Ctx, email, token string, request web.ResetPassword) error
	// FindById(ctx context.Context, id string, filter web.UserQueryFilter) (web.UserResponse, error)
	// FindByEmail(ctx context.Context, email string) (web.UserResponse, error)
	// FindByPhoneNumber(ctx context.Context, phone string) (web.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	kafkaProducerService producers.KafkaProducerService
	logger *zap.SugaredLogger
}

func NewUserService(userRepository repository.UserRepository, kafkaProducerService producers.KafkaProducerService, logger *zap.SugaredLogger) UserService {
	return &userService{
		userRepository: userRepository,
		kafkaProducerService: kafkaProducerService,
	}
}


// Password
func (service *userService) ForgetPasswordEmail(ctx *fiber.Ctx, email string) (domain.ResetPasswordToken, error) {
	user, err := service.userRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "email", email)

	if err != nil || user.Email == "" {
		return domain.ResetPasswordToken{}, fiber.NewError(fiber.StatusBadRequest, "Email tidak ditemukan.")
	}
	
	tokens, _ := service.userRepository.CheckTokenWithQueryTx(ctx.Context(), "email", email)

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

		err = service.userRepository.AddTokenTx(ctx.Context(), data)
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

		err = service.userRepository.UpdateTokenTx(ctx.Context(), tokens)
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

func (service *userService) ResetPassword(ctx *fiber.Ctx, email, token string, request web.ResetPassword) error {
	// var decodedByte, _ = base64.StdEncoding.DecodeString(token)
	// var resetToken = string(decodedByte)

	if request.Password != request.PasswordConfirm {
		panic(exception.ErrBadRequest("Password didn't match."))
	}

	checkToken, err := service.userRepository.CheckTokenWithQueryTx(ctx.Context(), "tokens", token)

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

	user, err = service.userRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "email", email)
	if err != nil {
		return exception.ErrNotFound(err.Error())
	}

	if len(request.Password) < 6 {
		return exception.ErrBadRequest("Password length should more then equal 6 character.")
	}

	user.UpdatedAt = time.Now()
	user.SetPassword(request.Password)

	service.userRepository.UpdatePasswordTx(ctx.Context(), user)

	// Update Token
	nol := 0
	checkToken.Attempt = &nol
	checkToken.Tokens = ""
	err = service.userRepository.UpdateTokenTx(ctx.Context(), checkToken)
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


// func (s *userService) CreateUser(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	// validate email exist
// 	if _, err := s.userRepository.FindByEmail(c, request.Email); err != nil {
// 		if strings.Contains(err.Error(), "no rows") {
// 			return web.UserResponse{}, exception.ErrNotFound(fmt.Sprintf("Supervisor %s not found", request.Email))
// 		} else {
// 			return web.UserResponse{}, err
// 		}
// 	}

// 	// convert to domain or model user
// 	user := domain.User{
// 		ID: uuid.New().String(),
// 		Name: request.Name,
// 		Email: request.Email,
// 		Password: request.Password,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	// call the repo for inserting to db
// 	if err := s.userRepository.CreateUser(c, user); err != nil {
// 		s.logger.Infow(err.Error(), "Create User Error")
// 		return web.UserResponse{}, err
// 	}

// 	// get or returning the repo have created to db
// 	newUser, err := s.userRepository.FindById(c, user.ID, domain.UserQueryFilter{
// 		ShowDeleted: false,
// 	})
// 	if err != nil {
// 		return web.UserResponse{}, exception.ErrInternalServer(fmt.Sprintf("Successfully created user, but failed to get the user have created. Error: %s", err.Error()))
// 	}

// 	// produce kafka create-user message
// 	kafkaUserMessage := kafkamodel.NewKafkaUserMessage(newUser)
// 	go s.kafkaProducerService.Produce(kafkaUserMessage, "POST.JOB", config.KafkaTopic)


// 	return newUser.
// }

// func (u *userService) UpdateUser(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) UpdatePassword(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) Delete(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
func (u *userService) FindAllUser(c context.Context, filter web.UserQueryFilter) (result []web.UserResponse, totalData int, err error) {
	repositoryResponse, err := u.userRepository.FindAllUser(c, domain.ToDomainUserQueryFilter(filter))
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, 0, exception.ErrNotFound("User not found")
		} else {
			return nil, 0, err
		}
	}

	// get total-data
	totalData, err = u.userRepository.CountAllUser(c, domain.ToDomainUserQueryFilter(filter))
	if err != nil {
		return nil, 0, err
	}

	// convert to web.JobResponse
	for _, user := range repositoryResponse {
		log.Println(user.Name)
		result = append(result, user.ToUserResponse())
	}
	log.Println(result)

	return result, totalData, err
}
// func (u *userService) FindById(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) FindByEmail(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) FindByPhoneNumber(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }