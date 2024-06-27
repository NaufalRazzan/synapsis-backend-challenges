package middlewares

import "github.com/gofiber/fiber/v2"

func UndefinedRoutesMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedpath := []string{
			"/register",
			"/login",
			"/v1/listProduct/",
			"/v1/insertShoppingCart",
			"/v1/listShoppingCart",
			"/v1/deleteShoppingCart",
			"/v1/checkout",
		}

		matchedPath := false
		for _, testpath := range allowedpath {
			if c.Path() == "/" {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"message": "application running smoothly",
				})
			}
			if c.Path() == testpath {
				matchedPath = true
				break
			}
		}

		if !matchedPath {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "resource not found",
			})
		}

		return c.Next()
	}
}
