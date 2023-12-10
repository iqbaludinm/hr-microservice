package exception

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaludinm/hr-microservice/user-service/model/web"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	response := web.WebResponse{
		Code:    code,
		Status:  false,
		Message: err.Error(),
		Data:    err.Error(),
	}

	return c.Status(code).JSON(response)
}
