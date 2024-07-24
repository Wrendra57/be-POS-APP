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

func CreateSupplier(service services.SupplierService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.SupplierRequest{}
		if err := ctx.BodyParser(&request); err != nil {
			fmt.Println(err)
			return exception.CustomResponse(ctx, 500, "Internal server error", nil)
		}
		//	validasi
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		s, errs, e := service.CreateSupplier(ctx, request)
		if e == false {
			return exception.CustomResponse(ctx, errs.Code, errs.Error, nil)
		}
		return exception.SuccessResponse(ctx, "success", s)
	}
}

func FindByParamsSupplier(service services.SupplierService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.SupplierListRequest{}
		if err := ctx.BodyParser(&request); err != nil {
			fmt.Println(err)
			return exception.CustomResponse(ctx, 500, "Internal server error", err)
		}
		//	validasi
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		s, _ := service.FindByParamSupplier(ctx, request)
		return exception.SuccessResponse(ctx, "success", s)
	}

}
