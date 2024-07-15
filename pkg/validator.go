package pkg

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func NewValidate() *validator.Validate {
	fmt.Println("init validator")
	validate := validator.New()
	return validate
}

func ValidateStruct(s interface{}, validate *validator.Validate) error {
	return validate.Struct(s)
}
