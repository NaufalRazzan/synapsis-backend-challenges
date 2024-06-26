package middlewares

import (
	"fmt"
	"synapsis-backend-test/configs"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler{
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(configs.GetConfig().JWTSecretKey),
	})
}

func VerifyAuth() fiber.Handler{
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == ""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "missing auth header",
			})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
				return nil, fmt.Errorf("unexpected signin method: %v", t.Method)
			}

			return []byte(configs.GetConfig().JWTSecretKey), nil
		})

		if err != nil || !token.Valid{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Locals("user_id", claims["user_id"].(string))

		return c.Next()
	}
}