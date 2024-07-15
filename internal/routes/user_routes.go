package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router, service services.UserService, validate *validator.Validate) {
	//api := app.Group("/v1/users")
	app.Get("/v1/users", controllers.CreateUser(service, validate))
	app.Post("/v1/users", controllers.CreateUser(service, validate))
	//app.
}
