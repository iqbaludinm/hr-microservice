package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaludinm/hr-microservice/profile-service/helper"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/web"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("token") // ambil token di cookies, dengan key "token"

	issuer, err := helper.ParseJwt(cookie)
	if err != nil {
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

	c.Locals("issuer", issuer)
	return c.Next()
}
