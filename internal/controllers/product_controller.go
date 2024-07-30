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
		utils.PanicIfError(err)
		request.SellPrice = price

		categoryStr := ctx.FormValue("category")
		if categoryStr == "" {
			return exception.CustomResponse(ctx, 400, "category is required", nil)
		}
		if err := json.Unmarshal([]byte(categoryStr), &request.Category); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse Category",
				"err":   err,
			})
		}

		brandStr := ctx.FormValue("brand")
		if brandStr == "" {
			return exception.CustomResponse(ctx, 400, "brand is required", nil)
		}
		type brand struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}
		var b brand

		if err := json.Unmarshal([]byte(brandStr), &b); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse brand",
				"err":   err,
			})
		}
		request.Brand.Name = b.Name
		BrandId, err := strconv.Atoi(b.Id)
		utils.PanicIfError(err)
		request.Brand.Id = BrandId

		supplierStr := ctx.FormValue("supplier")
		if supplierStr == "" {
			return exception.CustomResponse(ctx, 400, "supplier is required", nil)
		}
		if err := json.Unmarshal([]byte(supplierStr), &request.Supplier); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse supplier",
				"err":   err,
			})
		}
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse multipart form",
				"err":   err,
			})
		}

		file := form.File["photo"]
		utils.PanicIfError(err)
		request.Photo = file

		if len(file) == 0 {
			return exception.CustomResponse(ctx, 400, "Photo minimal 1", nil)
		}
		//validate
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}

		fmt.Println("jalan service")
		product, e, ok := service.CreateProduct(ctx, request)
		if !ok {
			return exception.CustomResponse(ctx, e.Code, e.Error, nil)
		}

		return exception.SuccessResponse(ctx, "Success", product)
	}
}
