package middlewares

import (
	"synapsis-backend-test/configs"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(configs.GetConfig().JWTSecretKey),
		SigningMethod: "HS256",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		},
	})
}
