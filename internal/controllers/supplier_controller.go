package controllers

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
		return exception.SuccessResponse(ctx, "Success create supplier", s)
	}
}

func FindByParamsSupplier(service services.SupplierService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.SupplierListRequest{}

		request.Params = ctx.Query("params")
		limitStr := ctx.Query("limit")
		offsetStr := ctx.Query("offset")

		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				return exception.CustomResponse(ctx, 400, "The 'limit' field must be number/integer", nil)
			}
			request.Limit = limit
		} else {
			request.Limit = 15
		}

		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				return exception.CustomResponse(ctx, 400, "The 'offset' field must be number/integer", nil)
			}
			request.Offset = offset
		} else {
			request.Offset = 1
		}

		s, _ := service.FindByParamSupplier(ctx, request)
		return exception.SuccessResponse(ctx, "Success get data supplier", s)
	}

}
