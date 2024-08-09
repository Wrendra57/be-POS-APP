package controllers

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CreateBrand(service services.BrandService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.BrandCreateReq{}
		if err := ctx.BodyParser(&request); err != nil {
			return exception.CustomResponse(ctx, 500, "Internal server error", nil)
		}
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		b := service.CreateBrand(ctx, request)
		return exception.SuccessResponse(ctx, "success", b)
	}
}

func ListBrand(service services.BrandService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.BrandGetRequest{}

		request.Params = ctx.Query("params")
		limitStr := ctx.Query("limit")
		offsetStr := ctx.Query("offset")
		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				return exception.CustomResponse(ctx, 400, "The 'limit' field must be number/integer", nil)
			}
			request.Limit = limit
			if limit <= 0 {
				return exception.CustomResponse(ctx, 400, "The 'limit' field must be greater than zero", nil)
			}
		} else {
			request.Limit = 15
		}

		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				return exception.CustomResponse(ctx, 400, "The 'offset' field must be number/integer", nil)
			}
			if offset <= 0 {
				return exception.CustomResponse(ctx, 400, "The 'offset' field must be positive", nil)
			}
			request.Offset = offset
		} else {
			request.Offset = 1
		}
		b := service.ListBrand(ctx, request)

		return exception.SuccessResponse(ctx, "Success get data brand", b)
	}
}
