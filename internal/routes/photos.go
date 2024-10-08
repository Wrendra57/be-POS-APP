package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func FileRoutes(app fiber.Router, photoService services.PhotosService, validate *validator.Validate) {
	app.Post("/v1/file/upload", middleware.Authenticate(), controllers.UploadPhoto(photoService, validate))
}
