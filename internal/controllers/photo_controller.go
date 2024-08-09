package controllers

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/services"
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
		if owner == "" {
			return exception.CustomResponse(ctx, 400, "owner_id is required", nil)
		}

		uuidObj, err := uuid.Parse(owner)
		if err != nil {
			return exception.CustomResponse(ctx, 400, "owner_id is invalid", nil)
		}

		request.Owner_id = uuidObj

		request.Name = ctx.FormValue("name")

		file, err := ctx.FormFile("foto")
		if err != nil {
			return exception.CustomResponse(ctx, 400, "photo is required", nil)
		}

		request.Foto = file

		//validate
		if err := pkg.ValidateStruct(&request, validate); err != nil {
			errors := exception.FormatValidationError(err)
			return exception.ValidateErrorResponse(ctx, "Validation error", errors)
		}
		foto := service.UploadPhotoService(ctx, request)

		return exception.SuccessResponse(ctx, "success", foto)

	}
}
