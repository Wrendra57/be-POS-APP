package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SupplierRoutes(app fiber.Router, supplierService services.SupplierService, validate *validator.Validate) {
	app.Post("/v1/supplier", middleware.Authenticate(), controllers.CreateSupplier(supplierService, validate))
	app.Get("/v1/supplier", controllers.FindByParamsSupplier(supplierService, validate))
}
