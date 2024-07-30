package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id          uuid.UUID `json:"id"`
	ProductName string    `json:"product_name"`
	SellPrice   int       `json:"sell_price"`
	CallName    string
	AdminId     uuid.UUID
	CategoryId  uuid.UUID
	BrandId     int
	SupplierId  uuid.UUID
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}
