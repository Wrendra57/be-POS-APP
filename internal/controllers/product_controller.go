package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

func CreateProduct(service services.ProductService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.ProductCreateRequest{}

		request.ProductName = ctx.FormValue("product_name")
		request.CallName = ctx.FormValue("call_name")

		priceStr := ctx.FormValue("sell_price")
		if priceStr == "" {
			return exception.CustomResponse(ctx, 400, "sell price is required", nil)
		}
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			return exception.CustomResponse(ctx, 400, "sell price must be integer/number", nil)

		}
		request.SellPrice = price

		categoryStr := ctx.FormValue("category")
		if categoryStr == "" {
			return exception.CustomResponse(ctx, 400, "category is required", nil)
		}

		if err := json.Unmarshal([]byte(categoryStr), &request.Category); err != nil {
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, "invalid parse category", nil)
		}

		brandStr := ctx.FormValue("brand")
		if brandStr == "" {
			return exception.CustomResponse(ctx, 400, "brand is required", nil)
		}
		if err := json.Unmarshal([]byte(brandStr), &request.Brand); err != nil {
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, "invalid parse brand", nil)
		}

		supplierStr := ctx.FormValue("supplier")
		if supplierStr == "" {
			return exception.CustomResponse(ctx, 400, "supplier is required", nil)
		}
		if err := json.Unmarshal([]byte(supplierStr), &request.Supplier); err != nil {
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, "invalid parse supplier", nil)
		}

		form, err := ctx.MultipartForm()
		utils.PanicIfError(err)

		file := form.File["photo"]
		utils.PanicIfError(err)
		request.Photo = file

		if len(file) == 0 || len(file) >= 15 {
			return exception.CustomResponse(ctx, 400, "Photo minimal 1 or maximal 15", nil)
		}
		//validate
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}

		product := service.CreateProduct(ctx, request)
		return exception.SuccessResponse(ctx, "Success", product)
	}
}
func FindById(service services.ProductService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, "invalid id product", nil)
		}
		product, e, ok := service.FindProductById(ctx, parsedId)
		if !ok {
			return exception.CustomResponse(ctx, e.Code, e.Error, nil)
		}

		return exception.SuccessResponse(ctx, "Success get data", product)
	}
}

func ListProduct(service services.ProductService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.ProductListRequest{}

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

		p := service.ListProduct(ctx, request)
		return exception.SuccessResponse(ctx, "Success get data", p)
	}
}

func DeleteProduct(service services.ProductService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, "invalid id product", nil)
		}
		err = service.DeleteProduct(ctx, parsedId)

		if err != nil {
			fmt.Println("err", err.Error())
			return exception.CustomResponse(ctx, 400, err.Error(), nil)
		}

		return exception.SuccessResponse(ctx, "Success delete product", id)
	}
}
