package webrequest

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type category struct {
	Id   uuid.UUID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required"`
}
type brand struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
type supplier struct {
	Id   uuid.UUID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required"`
}
type ProductCreateRequest struct {
	ProductName string                  `json:"product_name" validate:"required,min=2,max=252"`
	SellPrice   int                     `json:"sell_price" validate:"required,gt=0"`
	CallName    string                  `json:"call_name" validate:"required,min=2,max=1052"`
	Category    category                `json:"category" validate:"required"`
	Brand       brand                   `json:"brand" validate:"required"`
	Supplier    supplier                `json:"supplier" validate:"required"`
	Photo       []*multipart.FileHeader `json:"photo" validate:"required"`
}
