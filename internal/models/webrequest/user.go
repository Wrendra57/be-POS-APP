package webrequest

import (
	"time"
)

type UserCreateRequest struct {
	Name              string    `json:"name" validate:"required,min=3,max=32"`
	Gender            string    `json:"gender" validate:"required,oneof=male female"`
	Telp              string    `json:"telp" validate:"required,min=3,max=32"`
	Birthday          string    `json:"birthday" validate:"required"`
	BirthdayConversed time.Time `json:"birthdayConversed" validate:"omitempty,required"`
	Address           string    `json:"address" validate:"required,min=1,max=255"`
	Email             string    `json:"email" validate:"required,email"`
	Password          string    `json:"password" validate:"required,min=8,max=32"`
	Username          string    `json:"username" validate:"required,min=3,max=32"`
}

type UserLoginRequest struct {
	UserName string `json:"username" validate:"omitempty,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
