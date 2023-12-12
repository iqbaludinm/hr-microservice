package controller

import (
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/iqbaludinm/hr-microservice/user-service/config"
	"github.com/iqbaludinm/hr-microservice/user-service/exception"
	"github.com/iqbaludinm/hr-microservice/user-service/middleware"
	"github.com/iqbaludinm/hr-microservice/user-service/model/web"
	"github.com/iqbaludinm/hr-microservice/user-service/service"
	"github.com/iqbaludinm/hr-microservice/user-service/service/producers"
)

type UserController interface {
	// This function is used to define the route of the endpoints in this controller to the main application
	Route(app *fiber.App)
	FindAllUser(ctx *fiber.Ctx) error
	ForgetPassword(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error

	// CreateUser(ctx *fiber.Ctx) error
	// UpdateUser(ctx *fiber.Ctx) error
	// DeleteUser(ctx *fiber.Ctx) error
	// FindByID(ctx *fiber.Ctx) error
}

type userController struct {
	validate             *validator.Validate
	kafkaProducerService producers.KafkaProducerService
	userService          service.UserService
}

func NewUserController(validate *validator.Validate, kafkaProducerService producers.KafkaProducerService, userService service.UserService) UserController {
	return &userController{
		kafkaProducerService: kafkaProducerService,
		validate:             validate,
		userService:          userService,
	}
}

func (controller *userController) Route(app *fiber.App) {
	api := app.Group(config.EndpointPrefixUser, middleware.IsAuthenticated)

	// api.Post("/",
	// 	controller.CreateJob,
	// )

	api.Get("/",
		controller.FindAllUser,
	)
	
	api.Post("/forget-password", controller.ForgetPassword)
	api.Post("/reset-password", controller.ResetPassword)

	// api.Get("/:job_id",
	// 	controller.FindByID,
	// )

	// api.Put("/:job_id",
	// 	// middleware.AuthorizePermissions([]int{
	// 	// 	helper.PermissionAllJob,
	// 	// 	helper.PermissionUpdateJob,
	// 	// }),
	// 	controller.UpdateJobDetails,
	// )
	// api.Put("/:job_id/supervisor",
	// 	// middleware.AuthorizePermissions([]int{
	// 	// 	helper.PermissionAllJob,
	// 	// 	helper.PermissionUpdateJob,
	// 	// }),
	// 	controller.UpdateJobSupervisor,
	// )
	// api.Put("/:job_id/site",
	// 	// middleware.AuthorizePermissions([]int{
	// 	// 	helper.PermissionAllJob,
	// 	// 	helper.PermissionUpdateJob,
	// 	// }),
	// 	controller.UpdateJobSite,
	// )
	// api.Put("/:job_id/part",
	// 	// middleware.AuthorizePermissions([]int{
	// 	// 	helper.PermissionAllJob,
	// 	// 	helper.PermissionUpdateJob,
	// 	// }),
	// 	controller.UpdateJobPart,
	// )
	// api.Put("/:job_id/add-member",
	// 	// middleware.AuthorizePermissions([]int{
	// 	// 	helper.PermissionAllJob,
	// 	// 	helper.PermissionUpdateJob,
	// 	// }),
	// 	controller.UpdateMember,
	// )

	// api.Delete("/:job_id",
	// 	// middleware.AuthorizePermissions([]int{
	// 	// 	helper.PermissionAllJob,
	// 	// 	helper.PermissionDeleteJob,
	// 	// }),
	// 	controller.DeleteJob,
	// )
}

func (controller *userController) ForgetPassword(ctx *fiber.Ctx) error {
	var data web.ForgetPassword
	_ = ctx.BodyParser(&data)

	result, err := controller.userService.ForgetPasswordEmail(ctx, data.Email)
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

func (controller *userController) ResetPassword(ctx *fiber.Ctx) error {
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

	err = controller.userService.ResetPassword(ctx, email, token, data)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset successfully.",
	})
}

// func (controller *userController) CreateJob(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse request body
// 	var request web.CreateJobRequest
// 	err := ctx.BodyParser(&request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}
// 	// validate the values of the request body
// 	err = controller.validate.Struct(request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}

// 	// create job
// 	jobResponse, err := controller.jobService.CreateJob(ctx.Context(), request)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("created job with id %s", jobResponse.ID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    jobResponse,
// 	})
// }

// func (controller *jobController) UpdateJobDetails(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse request body
// 	var request web.UpdateJobRequest
// 	err := ctx.BodyParser(&request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}
// 	// validate the values of the request body
// 	err = controller.validate.Struct(request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}

// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	// update job
// 	jobResponse, err := controller.jobService.UpdateJobDetails(ctx.Context(), jobID, request)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("updated job with id %s", jobResponse.ID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    jobResponse,
// 	})
// }

// func (controller *jobController) UpdateJobSupervisor(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse request body
// 	var request web.UpdateJobSupervisorRequest
// 	err := ctx.BodyParser(&request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}
// 	// validate the values of the request body
// 	err = controller.validate.Struct(request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}

// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	// update job
// 	jobResponse, err := controller.jobService.UpdateJobSupervisor(ctx.Context(), jobID, request)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("updated job with id %s", jobResponse.ID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    jobResponse,
// 	})
// }

// func (controller *jobController) UpdateJobSite(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse request body
// 	var request web.UpdateJobSiteRequest
// 	err := ctx.BodyParser(&request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}
// 	// validate the values of the request body
// 	err = controller.validate.Struct(request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}

// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	// update job
// 	jobResponse, err := controller.jobService.UpdateJobSite(ctx.Context(), jobID, request)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("updated job with id %s", jobResponse.ID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    jobResponse,
// 	})
// }

// func (controller *jobController) UpdateJobPart(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse request body
// 	var request web.UpdateJobPartRequest
// 	err := ctx.BodyParser(&request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}
// 	// validate the values of the request body
// 	err = controller.validate.Struct(request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}

// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	// update job
// 	jobResponse, err := controller.jobService.UpdateJobPart(ctx.Context(), jobID, request)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("updated job with id %s", jobResponse.ID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    jobResponse,
// 	})
// }

// func (controller *jobController) UpdateMember(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse request body
// 	var request web.UpdateJobMemberRequest
// 	err := ctx.BodyParser(&request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}
// 	// validate the values of the request body
// 	err = controller.validate.Struct(request)
// 	if err != nil {
// 		return exception.ErrValidateBadRequest(err.Error(), request)
// 	}

// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	// update job
// 	jobResponse, err := controller.jobMemberService.EditMember(ctx.Context(), jobID, request)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("updated job-members with id %s", jobResponse.ID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    jobResponse,
// 	})
// }

// func (controller *jobController) DeleteJob(ctx *fiber.Ctx) error {
// 	// parse jwt claims
// 	// actor, _, _, _ := helper.ParseJwt(ctx.Cookies("token"))

// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	// delete job
// 	err := controller.jobService.DeleteJob(ctx.Context(), jobID)
// 	if err != nil {
// 		return err
// 	}

// 	// kafka produce log message
// 	// kafkaLogMessage := kafkamodel.KafkaLogMessage{
// 	// 	Actor:     actor,
// 	// 	Action:    fmt.Sprintf("deleted job with id %s", jobID),
// 	// 	Timestamp: time.Now(),
// 	// }
// 	// go controller.kafkaProducerService.Produce(kafkaLogMessage, "POST.LOG", config.KafkaTopicLog)

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 	})
// }

// func (controller *jobController) FindByID(ctx *fiber.Ctx) error {
// 	// parse path params
// 	jobID := ctx.Params("job_id")

// 	job, err := controller.jobService.FindByID(ctx.Context(), jobID, web.JobQueryFilter{})
// 	if err != nil {
// 		return err
// 	}
// 	if job.ID == "" || job.CreatedAt.IsZero() {
// 		return exception.ErrNotFound("job not found")
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
// 		Code:    fiber.StatusOK,
// 		Status:  true,
// 		Message: "success",
// 		Data:    job,
// 	})
// }

func (controller *userController) FindAllUser(ctx *fiber.Ctx) error {
	// parse query params
	var filter web.UserQueryFilter
	if err := ctx.QueryParser(&filter); err != nil {
		return exception.ErrValidateBadRequest(err.Error(), filter)
	}

	userResponses, totalData, err := controller.userService.FindAllUser(ctx.Context(), filter)
	if err != nil {
		return err
	}

	if filter.Page != "" || filter.Limit != "" {
		pageInt, _ := strconv.Atoi(filter.Page)
		if filter.Page == "" {
			pageInt = 1
		}
		return ctx.Status(fiber.StatusOK).JSON(web.WebResponsePagination{
			Code:      fiber.StatusOK,
			Status:    true,
			Page:      pageInt,
			Count:     len(userResponses),
			TotalData: totalData,
			Message:   "success",
			Data:      userResponses,
		})
	}

	log.Println("DIBAWAH INI")
	log.Println(userResponses, totalData)
	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    userResponses,
	})
}
