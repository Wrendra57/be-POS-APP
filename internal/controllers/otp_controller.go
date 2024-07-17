package controllers

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func ValidateOTP(service services.OTPService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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
		error, err := service.ValidateOtpAccount(ctx, request.Otp)

		if err == false {
			return exception.CustomResponse(ctx, error.Code, error.Error, nil)
		}
		token, _ := ctx.Locals("token").(string)
		return exception.SuccessResponse(ctx, "success validate", token)
	}
}
func ReSendOtp(service services.OTPService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		now := time.Now()
		userId, _ := ctx.Locals("user_id").(uuid.UUID)

		_, errs, e := service.ReSendOtp(ctx, userId)
		token, _ := ctx.Locals("token").(string)
		fmt.Println(time.Now().Sub(now))
		if e == false {
			return exception.CustomResponse(ctx, errs.Code, errs.Error, token)
		}
		return exception.SuccessResponse(ctx, "success send otp again", token)
	}
}
