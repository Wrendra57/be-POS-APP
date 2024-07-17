package controllers

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadPhoto(service services.PhotosService, validate *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := webrequest.PhotoUploadRequest{}

		owner := ctx.FormValue("owner_id")
		uuidObj, err := uuid.Parse(owner)
		utils.PanicIfError(err)

		request.Owner_id = uuidObj

		request.Name = ctx.FormValue("name")

		file, err := ctx.FormFile("foto")
		utils.PanicIfError(err)
		request.Foto = file

		//validate
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		foto, e, erro := service.UploadPhotos(ctx, request)
		if erro != true {
			fmt.Println(erro)
			return exception.CustomResponse(ctx, e.Code, e.Error, nil)
		}

		return exception.SuccessResponse(ctx, "success", foto)

	}
}
