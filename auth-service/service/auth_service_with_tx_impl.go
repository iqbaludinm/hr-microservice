package service

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iqbaludinm/hr-microservice/auth-service/config"
	"github.com/iqbaludinm/hr-microservice/auth-service/exception"
	"github.com/iqbaludinm/hr-microservice/auth-service/helper"
	"github.com/iqbaludinm/hr-microservice/auth-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/auth-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/auth-service/model/web"
	"github.com/iqbaludinm/hr-microservice/auth-service/repository"
	"github.com/iqbaludinm/hr-microservice/auth-service/service/producers"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	Validate       *validator.Validate
}

type AuthService interface {
	Register(ctx *fiber.Ctx, request web.RegisterRequest) (web.RegisterResponse, error)
	Login(ctx *fiber.Ctx, request web.LoginRequest) (fiber.Cookie, web.LoginResponse, error)
	Logout() fiber.Cookie
	ForgetPasswordEmail(ctx *fiber.Ctx, email string) (domain.ResetPasswordToken, error)
	ResetPassword(ctx *fiber.Ctx, email, token string, request web.ResetPassword) error
	// Refresh(ctx *fiber.Ctx, nik string) (fiber.Cookie, web.LoginResponse, error)
	FindUserNotDeleteByQueryTx(ctx *fiber.Ctx, query, value string) (domain.User, error)
	// FindUserWithNameNotDeleteByQueryTx(ctx *fiber.Ctx, query, value string) (domain.UserWithName, error)

	// CheckToken(ctx *fiber.Ctx, token, email string) error
	// RefreshToResponse(ctx *fiber.Ctx, nik string) (web.TokenResponse, web.LoginResponse, error)
}

type authService struct {
	authRepository       repository.AuthRepository
	kafkaProducerService producers.KafkaProducerService
	logger               *zap.SugaredLogger
}

func NewAuthService(authRepository repository.AuthRepository, kafkaProducerService producers.KafkaProducerService, logger *zap.SugaredLogger) AuthService {
	return &authService{
		authRepository:       authRepository,
		kafkaProducerService: kafkaProducerService,
	}
}

func (service *authService) Register(ctx *fiber.Ctx, request web.RegisterRequest) (web.RegisterResponse, error) {
	createdUser := domain.User{
		ID:        uuid.New().String(),
		Name:      request.Name,
		Email:     request.Email,
		Phone:     request.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser.SetPassword(request.Password)

	id, err := service.authRepository.RegisterTx(ctx.Context(), createdUser)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "users_email_key") {
				return web.RegisterResponse{}, exception.ErrBadRequest("Email already exist.")
			} else if strings.Contains(err.Error(), "users_phone_key") {
				return web.RegisterResponse{}, exception.ErrBadRequest("Phone already exist.")
			}
		}
		return web.RegisterResponse{}, err
	}

	createdUser.ID = id
	// produce to kafka
	kafkaRegisterMessage := kafkamodel.NewKafkaUserMessage(createdUser)
	go service.kafkaProducerService.Produce(kafkaRegisterMessage, "POST.USER", config.KafkaTopic)

	return domain.ToRegisterResponse(createdUser), nil
}

func (service *authService) Login(ctx *fiber.Ctx, request web.LoginRequest) (fiber.Cookie, web.LoginResponse, error) {
	user, err := service.authRepository.LoginTx(ctx.Context(), request.Email)
	if err != nil || user.ID == "" {
		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrNotFound("User tidak ditemukan.")
	}

	err = user.ComparePassword(user.Password, request.Password)
	if err != nil {
		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrBadRequest("Password salah.")
	}

	token, err := helper.GenerateJwt(user.ID, user.Name, user.Email, user.Phone)
	if err != nil {
		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrBadRequest(err.Error())
	}

	session, _ := strconv.Atoi(config.SessionLogin)
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * time.Duration(session)),
		HTTPOnly: true,
	}

	return cookie, domain.ToLoginResponse(user), nil
}

func (service *authService) Logout() fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	return cookie
}

// Password
func (service *authService) ForgetPasswordEmail(ctx *fiber.Ctx, email string) (domain.ResetPasswordToken, error) {
	user, err := service.authRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "email", email)

	if err != nil || user.Email == "" {
		return domain.ResetPasswordToken{}, fiber.NewError(fiber.StatusBadRequest, "Email tidak ditemukan.")
	}
	
	tokens, _ := service.authRepository.CheckTokenWithQueryTx(ctx.Context(), "email", email)

	var data domain.ResetPasswordToken
	urlReset := config.UrlReset

	if tokens.Id == "" {
		attempt := 1

		token := strings.ToLower(randstr.String(30))
		data.Id = uuid.New().String()
		data.Tokens = token
		data.Email = email
		data.Attempt = &attempt
		data.LastAttempt = time.Now()
		data.URL = urlReset + "?email=" + email + "&token=" + token

		err = service.authRepository.AddTokenTx(ctx.Context(), data)
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

		err = service.authRepository.UpdateTokenTx(ctx.Context(), tokens)
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

func (service *authService) ResetPassword(ctx *fiber.Ctx, email, token string, request web.ResetPassword) error {
	// var decodedByte, _ = base64.StdEncoding.DecodeString(token)
	// var resetToken = string(decodedByte)

	if request.Password != request.PasswordConfirm {
		panic(exception.ErrBadRequest("Password didn't match."))
	}

	checkToken, err := service.authRepository.CheckTokenWithQueryTx(ctx.Context(), "tokens", token)

	if err != nil {
		log.Println("di token invalid 1")
		return exception.ErrBadRequest("Token invalid.")
	}
	
	if checkToken.Tokens != token || checkToken.Email != email {
		log.Println("di token invalid 2")
		return exception.ErrBadRequest("Token invalid.")
	}

	if time.Since(checkToken.LastAttempt) > 1*time.Hour {
		return exception.ErrBadRequest("Token sudah kadaluarsa.")
	}

	var user domain.User

	user, err = service.authRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "email", email)
	if err != nil {
		return exception.ErrNotFound(err.Error())
	}

	if len(request.Password) < 6 {
		return exception.ErrBadRequest("Password length should more then equal 6 character.")
	}

	user.UpdatedAt = time.Now()
	user.SetPassword(request.Password)

	service.authRepository.UpdatePasswordTx(ctx.Context(), user)

	// Update Token
	nol := 0
	checkToken.Attempt = &nol
	checkToken.Tokens = ""
	err = service.authRepository.UpdateTokenTx(ctx.Context(), checkToken)
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

	return nil
}

// func (service *AuthServiceImpl) Refresh(ctx *fiber.Ctx, nik string) (fiber.Cookie, web.LoginResponse, error) {

// 	user, err := service.AuthRepository.LoginTx(ctx.Context(), nik)
// 	if err != nil || user.NIK == "" {
// 		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrorNotFound("User not found.")
// 	}

// 	if *user.Status == 0 {
// 		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrorUnauthorize("User not active.")
// 	}

// 	token, err := helper.GenerateJwt(user.NIK, fmt.Sprintf("%d", *user.RoleId), []helper.PermissionRole{})
// 	if err != nil {
// 		return fiber.Cookie{}, web.LoginResponse{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
// 	}

// 	cookie := fiber.Cookie{
// 		Name:     "token",
// 		Value:    token,
// 		Expires:  time.Now().Add(time.Hour * 24),
// 		HTTPOnly: true,
// 	}

// 	return cookie, domain.ToLoginResponse(user), nil
// }

func (service *authService) FindUserNotDeleteByQueryTx(ctx *fiber.Ctx, query, value string) (domain.User, error) {

	user, err := service.authRepository.FindUserNotDeleteByQueryTx(ctx.Context(), query, value)
	if err != nil {
		log.Println(err)
	}

	return user, nil
}

// Token
// func (service *AuthServiceImpl) CheckToken(ctx *fiber.Ctx, token, email string) error {

// 	var decodedByte, _ = base64.StdEncoding.DecodeString(token)
// 	var resetToken = string(decodedByte)

// 	checkToken, err := service.AuthRepository.CheckTokenWithQueryTx(ctx.Context(), "tokens", resetToken)

// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "Token invalid.")
// 	}

// 	if checkToken.Tokens != resetToken || checkToken.Email != email {
// 		return fiber.NewError(fiber.StatusBadRequest, "Token invalid.")
// 	}
// 	return nil
// }

// func (service *AuthServiceImpl) RefreshToResponse(ctx *fiber.Ctx, nik string) (web.TokenResponse, web.LoginResponse, error) {

// 	user, err := service.AuthRepository.LoginTx(ctx.Context(), nik)
// 	if err != nil || user.NIK == "" {
// 		return web.TokenResponse{}, web.LoginResponse{}, exception.ErrorNotFound("User not found.")
// 	}

// 	if *user.Status == 0 {
// 		return web.TokenResponse{}, web.LoginResponse{}, exception.ErrorUnauthorize("User not active.")
// 	}

// 	token, err := helper.GenerateJwt(user.NIK, fmt.Sprintf("%d", *user.RoleId), []helper.PermissionRole{})
// 	if err != nil {
// 		return web.TokenResponse{}, web.LoginResponse{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
// 	}

// 	return web.TokenResponse{Token: token}, domain.ToLoginResponse(user), nil
// }
