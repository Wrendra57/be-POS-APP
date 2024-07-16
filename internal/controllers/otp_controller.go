package controllers

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateOTP(service services.OTPService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Println("ddd")
		request := webrequest.ValidateOtpRequest{}

		if err := ctx.BodyParser(&request); err != nil {
			fmt.Println(err)
			return exception.CustomResponse(ctx, 500, "Internal Server Error", nil)
		}
		//	validasi
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		error, err := service.ValidateOtp(ctx, request.Otp)

		if err == false {
			return exception.CustomResponse(ctx, error.Code, error.Error, nil)
		}
		token, _ := ctx.Locals("token").(string)
		return exception.SuccessResponse(ctx, "success validate", token)
	}
}
