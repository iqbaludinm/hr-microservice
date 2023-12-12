package controller

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"

	"github.com/iqbaludinm/hr-microservice/profile-service/config"
	"github.com/iqbaludinm/hr-microservice/profile-service/exception"
	"github.com/iqbaludinm/hr-microservice/profile-service/middleware"
	"github.com/iqbaludinm/hr-microservice/profile-service/service/producers"

	"github.com/iqbaludinm/hr-microservice/profile-service/model/web"
	"github.com/iqbaludinm/hr-microservice/profile-service/service"
)

type ProfileController interface {
	Route(app *fiber.App)
	UpdateMyProfile(ctx *fiber.Ctx) error
	ForgetPassword(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
}

type profileController struct {
	validate *validator.Validate
	kafkaProducerService producers.KafkaProducerService
	profileService service.ProfileService
}

func NewProfileController(validate *validator.Validate, kafkaProducerService producers.KafkaProducerService, profileService service.ProfileService) ProfileController {
	return &profileController{
		kafkaProducerService: kafkaProducerService,
		validate: validate,
		profileService: profileService,
	}
}

func (controller *profileController) Route(app *fiber.App) {
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
			Code:    fiber.StatusOK,
			Status:  true,
			Message: "ok",
		})
	})
	api := app.Group(config.EndpointPrefixProfile, middleware.IsAuthenticated)
	api.Put("/:profile_id", controller.UpdateMyProfile)
	api.Post("/forget-password", controller.ForgetPassword)
	api.Post("/reset-password", controller.ResetPassword)
}

func (controller *profileController) UpdateMyProfile(ctx *fiber.Ctx) error {
	var request web.UpdateProfileRequest
	profileID := ctx.Params("profile_id")
	
	err := ctx.BodyParser(&request)
	if err != nil {
		log.Println("BodyParser")
		exception.ErrValidateBadRequest(err.Error(), request)
	}
	
	// validate the values of the request body
	err = controller.validate.Struct(&request)
	log.Println(err)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}
	
	profileResponse, err := controller.profileService.UpdateMyProfile(ctx, profileID, request)
	
	if err != nil {
		log.Println("ProfileService")
		return exception.ErrorHandler(ctx, err)
	}

	// action := fmt.Sprintf("register user %s", ProfileResponse.ID)
	// data := web.LogCreateRequest{
	// 	Actor:     "",
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }

	return ctx.Status(fiber.StatusCreated).JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "success",
		Data:    profileResponse,
	})
}

func (controller *profileController) ForgetPassword(ctx *fiber.Ctx) error {
	var data web.ForgetPassword
	_ = ctx.BodyParser(&data)

	result, err := controller.profileService.ForgetPasswordEmail(ctx, data.Email)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset password has been sent.",
		Data: result,
	})
}

func (controller *profileController) ResetPassword(ctx *fiber.Ctx) error {
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

	err = controller.profileService.ResetPassword(ctx, email, token, data)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset successfully.",
	})
}