package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CategoriesRoutes(app fiber.Router, category services.CategoryService, validate *validator.Validate) {
	app.Post("v1/categories", middleware.Authenticate(), controllers.CreateCategory(category, validate))
	app.Post("/v1/categories/search", controllers.FindByParams(category, validate))
}
