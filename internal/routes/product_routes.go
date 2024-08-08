package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app fiber.Router, service services.ProductService, validate *validator.Validate) {
	app.Post("/v1/product", middleware.Authenticate(), controllers.CreateProduct(service, validate))
	app.Get("/v1/product/:id", controllers.FindById(service))
	app.Get("/v1/product", controllers.ListProduct(service))

}
