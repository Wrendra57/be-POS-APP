package exception

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

//	func ErrorHandler(ctx *fiber.Ctx, err interface{}) {
//		if _, ok := validationErrors(ctx, err); ok {
//			return
//		}
//		if customError(ctx, err) {
//			return
//		}
//		internalServerError(ctx, err)
//	}

//
//func validationErrors(ctx *fiber.Ctx, err interface{}) (error, bool) {
//	_, ok := err.(validator.ValidationErrors)
//
//	if ok {
//		ctx.Set("Content-Type", "application/json")
//		//ctx.Status(fiber.StatusBadRequest)
//
//		for _, e := range err.(validator.ValidationErrors) {
//
//			response := webrespones.ResponseApi{
//				Code:    fiber.StatusBadRequest,
//				Status:  "failed",
//				Message: customErrorMessage(e),
//				Data:    nil,
//			}
//			//send := ctx.Status(fiber.StatusBadRequest).JSON(response)
//			return NotFoundResponse(ctx, response.Message), true
//		}
//		return nil, true
//	} else {
//		return nil, false
//	}
//
//}
//
//func customError(ctx *fiber.Ctx, err interface{}) bool {
//	exception, ok := err.(CustomEror)
//	if ok {
//		ctx.Set("Content-Type", "application/json")
//		ctx.Status(fiber.StatusBadRequest)
//
//		response := webrespones.ResponseApi{
//			Code:    exception.Code,
//			Status:  "failed",
//			Message: exception.Error,
//			Data:    nil,
//		}
//		ctx.Status(fiber.StatusBadRequest).JSON(response)
//		return true
//	} else {
//		return false
//	}
//}
//
//func internalServerError(ctx *fiber.Ctx, err interface{}) {
//	ctx.Set("Content-Type", "application/json")
//	ctx.Status(fiber.StatusInternalServerError)
//	response := webrespones.ResponseApi{
//		Code:    fiber.StatusInternalServerError,
//		Status:  "failed",
//		Message: "Internal Server Error",
//		Data:    nil,
//	}
//	ctx.JSON(response)
//}

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {

	return c.Status(fiber.StatusOK).JSON(JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}
func customErrorMessage(err *ValidationErrorResponse) string {
	switch err.Tag {
	case "required":
		return fmt.Sprintf("%s is required", err.FailedField)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.FailedField, err.Value)
	case "max":
		return fmt.Sprintf("%s must be maximum %s characters long", err.FailedField, err.Value)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.FailedField)
	case "oneof":
		return fmt.Sprintf("%s must be a %s", err.FailedField, err.Value)
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime format YYYY-MM-DD", err.FailedField)

	default:
		return fmt.Sprintf("Validation error on field %s with tag %s", err.FailedField, err.Tag)
	}
}
func ValidateErrorResponse(c *fiber.Ctx, message string, data []*ValidationErrorResponse) error {

	return c.Status(fiber.StatusBadRequest).JSON(JSONResponse{
		Status:  "error",
		Message: customErrorMessage(data[0]),
		Data:    data,
	})

}

func NotFoundResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(JSONResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}

func FormatValidationError(err error) []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return errors
	}

	for _, err := range err.(validator.ValidationErrors) {
		var element ValidationErrorResponse
		element.FailedField = err.Field()
		element.Tag = err.Tag()
		element.Value = err.Param()
		errors = append(errors, &element)
	}
	return errors
}

func CustomResponse(c *fiber.Ctx, code int, message string, data interface{}) error {
	fmt.Println("ew")
	return c.Status(code).JSON(JSONResponse{
		Status:  "failed",
		Message: message,
		Data:    data,
	})
}
