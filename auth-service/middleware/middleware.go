package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	// "github.com/iqbaludinm/hr-microservice/auth-service/config"
	"github.com/iqbaludinm/hr-microservice/auth-service/helper"
	"github.com/iqbaludinm/hr-microservice/auth-service/model/web"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("token")

	if _, _, err := helper.ParseJwt(cookie); err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return c.Status(401).JSON(web.WebResponse{
				Code:    99281,
				Status:  false,
				Message: "token expired",
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(web.WebResponse{
			Code:    fiber.StatusUnauthorized,
			Status:  false,
			Message: "unauthorized",
		})
	}

	return c.Next()
}

// func IsAdmin(c *fiber.Ctx) bool {
// 	cookie := c.Cookies("token")

// 	if _, level, _, err := helper.ParseJwt(cookie); err != nil || level != config.VariableRoleAdmin {
// 		return false
// 	}

// 	return true
// }
