package middlewares

import "github.com/gofiber/fiber/v2"

func ErrorMiddleware() fiber.Handler{
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil{
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"message": e.Message,
				})
			}
		}

		return nil
	}
}