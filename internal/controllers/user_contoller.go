package controllers

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CreateUser(service services.UserService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.UserCreateRequest{}
		err := ctx.BodyParser(&request)
		utils.PanicIfError(err)

		// parsing date
		layout := "2006-01-02"
		if request.Birthday == "" {
			message := "Birthdate is required"
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, message, nil)
		}
		parsedTime, err := time.Parse(layout, request.Birthday)
		if err != nil {
			message := "Birthdate must be format YYYY-MM-DD"
			return exception.CustomResponse(ctx, fiber.StatusBadRequest, message, nil)
		}

		request.BirthdayConversed = parsedTime

		//validasi
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		createUser, error, err := service.CreateUser(ctx, request)

		if err != nil {
			return exception.CustomResponse(ctx, error.Code, error.Error, nil)
		}

		responseApi := webrespones.ResponseApi{
			Code:    fiber.StatusOK,
			Status:  "ok",
			Message: "User created successfully",
			Data:    createUser,
		}
		return ctx.Status(fiber.StatusOK).JSON(responseApi)
	}
}

func LoginUser(service services.UserService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.UserLoginRequest{}
		err := ctx.BodyParser(&request)
		utils.PanicIfError(err)

		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}

		resp, er, ok := service.Login(ctx, request)

		if !ok {
			return exception.CustomResponse(ctx, er.Code, er.Error, nil)
		}
		return exception.SuccessResponse(ctx, "Success login", resp)

	}
}

func AuthMe(service services.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		resp, er, ok := service.AuthMe(ctx)
		if !ok {
			return exception.CustomResponse(ctx, er.Code, er.Error, nil)
		}
		return exception.SuccessResponse(ctx, "Success Get Data", resp)

	}
}
