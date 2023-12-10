package controller

import (
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"

	"github.com/iqbaludinm/hr-microservice/auth-service/config"
	"github.com/iqbaludinm/hr-microservice/auth-service/exception"
	"github.com/iqbaludinm/hr-microservice/auth-service/service/producers"

	"github.com/iqbaludinm/hr-microservice/auth-service/helper"
	"github.com/iqbaludinm/hr-microservice/auth-service/model/web"
	"github.com/iqbaludinm/hr-microservice/auth-service/service"
)

type AuthController interface {
	// NewAuthRouter(app *fiber.App)
	
	Route(app *fiber.App)
	Register(ctx *fiber.Ctx) error
	ForgetPassword(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type authController struct {
	validate *validator.Validate
	kafkaProducerService producers.KafkaProducerService
	authService service.AuthService
}

func NewAuthController(validate *validator.Validate, kafkaProducerService producers.KafkaProducerService, authService service.AuthService) AuthController {
	return &authController{
		kafkaProducerService: kafkaProducerService,
		validate: validate,
		authService: authService,
	}
}

func (controller *authController) Route(app *fiber.App) {
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
			Code:    fiber.StatusOK,
			Status:  true,
			Message: "ok",
		})
	})
	api := app.Group(config.EndpointPrefixAuth)
	api.Post("/register", controller.Register)
	api.Post("/login", controller.Login)
	api.Post("/forget-password", controller.ForgetPassword)
	api.Post("/reset-password", controller.ResetPassword)
	// api.Post("/check-reset-token", controller.CheckResetToken)
	// api.Post("/refresh", controller.Refresh)
	// api.Post("/refresh-token", controller.RefreshToResponse)
	api.Post("/logout", controller.Logout)
}

func (controller *authController) Register(ctx *fiber.Ctx) error {
	var request web.RegisterRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		log.Println("BodyParser")
		exception.ErrValidateBadRequest(err.Error(), request)
	}
	
	// validate the values of the request body
	err = controller.validate.Struct(&request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}
	
	authResponse, err := controller.authService.Register(ctx, request)
	
	if err != nil {
		log.Println("AuthService")
		return exception.ErrorHandler(ctx, err)
	}

	// action := fmt.Sprintf("register user %s", AuthResponse.ID)
	// data := web.LogCreateRequest{
	// 	Actor:     "",
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }

	return ctx.Status(fiber.StatusCreated).JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "success",
		Data:    authResponse,
	})
}

func (controller *authController) Login(ctx *fiber.Ctx) error {
	var request web.LoginRequest
	_ = ctx.BodyParser(&request)

	cookie, resp, err := controller.authService.Login(ctx, request)

	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	refreshJwt, err := helper.GenerateRefreshJwt(resp.ID, resp.Name, resp.Email, resp.Phone)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	session, _ := strconv.Atoi(config.SessionRefreshToken)

	refreshToken := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshJwt,
		Expires:  time.Now().Add(time.Hour * time.Duration(session)),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)
	ctx.Cookie(&refreshToken)

	// kafka post log
	// user, _ := controller.authService.FindUserNotDeleteByQueryTx(ctx, "email", request.Email)
	// action := fmt.Sprintf("login user %s", user.Name)
	// data := web.LogCreateRequest{
	// 	Actor:     request.Email,
	// 	ActorName: user.Name,
	// 	Category:  config.CategoryService,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(data, "POST.LOG", config.KafkaLogTopic)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    resp,
	})
}

func (controller *authController) Logout(ctx *fiber.Ctx) error {
	cookie := controller.authService.Logout()

	refreshToken := fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)
	ctx.Cookie(&refreshToken)

	//produce to kafka system log
	// cookieData := ctx.Cookies("token")
	// actor, _, _, _ := helper.ParseJwt(cookieData)
	// user, _ := controller.AuthService.FindUserWithNameNotDeleteByQueryTx(ctx, "nik", actor)
	// action := fmt.Sprintf("logout user %s", user.Name)
	// data := web.LogCreateRequest{
	// 	Actor:     actor,
	// 	ActorName: user.Name,
	// 	Category:  config.CategoryService,
	// 	Project:   user.ProjectName,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(data, "POST.LOG", config.KafkaLogTopic)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "logout berhasil",
	})
}

func (controller *authController) ForgetPassword(ctx *fiber.Ctx) error {
	var data web.ForgetPassword
	_ = ctx.BodyParser(&data)

	result, err := controller.authService.ForgetPasswordEmail(ctx, data.Email)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	// user, _ := controller.AuthService.FindUserWithNameNotDeleteByQueryTx(ctx, "email", data.Email)
	// action := fmt.Sprintf("request reset password for %s", user.Name)
	// dataLog := web.LogCreateRequest{
	// 	Actor:     user.NIK,
	// 	ActorName: user.Name,
	// 	Category:  config.CategoryService,
	// 	Project:   user.ProjectName,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(dataLog, "POST.LOG", config.KafkaLogTopic)
	// convertData := json.Unmarshal([]byte(data), &myData)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset password has been sent.",
		Data: result,
	})
}

func (controller *authController) ResetPassword(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	token := ctx.Query("token")

	if len(token) == 0 {
		return exception.ErrorHandler(ctx, exception.ErrBadRequest("Token missing."))
	}
	if len(email) == 0 {
		return exception.ErrorHandler(ctx, exception.ErrBadRequest("Email missing."))
	}

	var data web.ResetPassword
	_ = ctx.BodyParser(&data)

	// validate password field on req body
	err := controller.validate.Struct(&data)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), data)
	}

	err = controller.authService.ResetPassword(ctx, email, token, data)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	// produce to kafka post.log system
	// user, _ := controller.AuthService.FindUserWithNameNotDeleteByQueryTx(ctx, "email", email)
	// action := fmt.Sprintf("reset password for %s", user.Name)
	// dataLog := web.LogCreateRequest{
	// 	Actor:     user.NIK,
	// 	ActorName: user.Name,
	// 	Category:  config.CategoryService,
	// 	Project:   user.ProjectName,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(dataLog, "POST.LOG", config.KafkaLogTopic)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset successfully.",
	})
}

// func (controller *AuthControllerImpl) Refresh(ctx *fiber.Ctx) error {

// 	refreshCookie := ctx.Cookies("refresh_token")

// 	nik, err := helper.ParseRefreshJwt(refreshCookie)
// 	if err != nil {
// 		return exception.ErrorHandler(ctx, err)
// 	}

// 	cookie, resp, err := controller.AuthService.Refresh(ctx, nik)
// 	if err != nil {
// 		return exception.ErrorHandler(ctx, err)
// 	}

// 	refreshJwt, err := helper.GenerateRefreshJwt(nik)
// 	if err != nil {
// 		return exception.ErrorHandler(ctx, err)
// 	}

// 	session, _ := strconv.Atoi(utils.GetEnv("REFRESH_TOKEN"))

// 	refreshToken := fiber.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshJwt,
// 		Expires:  time.Now().Add(time.Hour * time.Duration(session)),
// 		HTTPOnly: true,
// 	}

// 	ctx.ClearCookie()
// 	ctx.Cookie(&cookie)
// 	ctx.Cookie(&refreshToken)

// 	cookieData := ctx.Cookies("token")
// 	actor, _, _, _ := helper.ParseJwt(cookieData)
// 	user, _ := controller.AuthService.FindUserWithNameNotDeleteByQueryTx(ctx, "nik", actor)
// 	action := fmt.Sprintf("refresh token for %s", user.Name)
// 	dataLog := web.LogCreateRequest{
// 		Actor:     actor,
// 		ActorName: user.Name,
// 		Project:   user.ProjectName,
// 		Category:  config.CategoryService,
// 		Action:    action,
// 		Timestamp: time.Now(),
// 	}
// 	helper.ProduceToKafka(dataLog, "POST.LOG", config.KafkaLogTopic)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    resp,
// 	})
// }

// // Token
// func (controller *AuthControllerImpl) CheckResetToken(ctx *fiber.Ctx) error {
// 	email := ctx.Query("email")
// 	token := ctx.Query("token")

// 	err := controller.AuthService.CheckToken(ctx, token, email)
// 	if err != nil {
// 		return exception.ErrorHandler(ctx, err)
// 	}
// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    nil,
// 	})
// }

// func (controller *AuthControllerImpl) RefreshToResponse(ctx *fiber.Ctx) error {

// 	refreshCookie := ctx.Cookies("refresh_token")

// 	nik, err := helper.ParseRefreshJwt(refreshCookie)
// 	if err != nil {
// 		return exception.ErrorHandler(ctx, err)
// 	}

// 	token, resp, err := controller.AuthService.RefreshToResponse(ctx, nik)
// 	if err != nil {
// 		return exception.ErrorHandler(ctx, err)
// 	}

// 	user, _ := controller.AuthService.FindUserWithNameNotDeleteByQueryTx(ctx, "nik", resp.NIK)
// 	action := fmt.Sprintf("refresh token to response for %s", user.Name)
// 	dataLog := web.LogCreateRequest{
// 		Actor:     resp.NIK,
// 		ActorName: user.Name,
// 		Project:   user.ProjectName,
// 		Category:  config.CategoryService,
// 		Action:    action,
// 		Timestamp: time.Now(),
// 	}
// 	helper.ProduceToKafka(dataLog, "POST.LOG", config.KafkaLogTopic)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    token,
// 	})
// }
