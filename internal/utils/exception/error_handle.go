package exception

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

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
