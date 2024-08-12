package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			r := recover()

			if r != nil {
				log.Printf("Panic recovered: %v", r)
				// You can customize the error response as per your needs
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "failed",
					"message": "Internal server error",
					"data":    nil,
				})
			}
		}()
		return c.Next()
	}
}
