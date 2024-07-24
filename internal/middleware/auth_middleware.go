package middleware

import (
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return exception.UnauthorizedRespone(c, "Unauthorized")
		}
		//Check if the Authorization header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return exception.UnauthorizedRespone(c, "Unauthorized")
		}

		//extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//parse token
		result, err := utils.ParseJWT(tokenString)

		if err != nil {
			return exception.UnauthorizedRespone(c, "Unauthorized")
		}
		c.Locals("user_id", result.User_id)
		c.Locals("token", tokenString)

		return c.Next()
	}
}
