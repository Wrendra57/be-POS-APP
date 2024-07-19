package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router, service services.UserService, validate *validator.Validate) {
	app.Post("/v1/users/register", controllers.CreateUser(service, validate))
	app.Post("/v1/users/login", controllers.LoginUser(service, validate))
	app.Get("/v1/user", middleware.Authenticate(), controllers.AuthMe(service))
}
