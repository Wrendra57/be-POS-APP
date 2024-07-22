package routes

import (
	"github.com/Wrendra57/Pos-app-be/internal/controllers"
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func OtpRoutes(app fiber.Router, otpService services.OTPService, validate *validator.Validate) {
	app.Post("/v1/users/otp", middleware.Authenticate(), controllers.ValidateOTP(otpService, validate))
	app.Post("/v1/users/otp/resend", middleware.Authenticate(), controllers.ReSendOtp(otpService, validate))
}
