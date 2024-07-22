package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func BrandRoutes(app fiber.Router, brand services.BrandService, validate *validator.Validate) {
	app.Post("v1/brands", middleware.Authenticate(), controllers.CreateBrand(brand, validate))
}
