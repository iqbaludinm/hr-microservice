package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaludinm/hr-microservice/user-service/helper"
	"github.com/iqbaludinm/hr-microservice/user-service/model/web"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("token") // ambil token di cookies, dengan key "token"

	if _, err := helper.ParseJwt(cookie); err != nil {
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
