package utils

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/gofiber/fiber/v2"
)

func PanicIfError(ctx *fiber.Ctx, code int, err error) {
	if err != nil {
		fmt.Println(err)

		exception.CustomResponse(ctx, code, err.Error(), nil)
		panic(err)

	}

}
