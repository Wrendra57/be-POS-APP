package controllers

import (
	"fmt"
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
			exception.CustomResponse(ctx, fiber.StatusBadRequest, message, nil)
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
		fmt.Println("9")
		fmt.Println(createUser)
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
