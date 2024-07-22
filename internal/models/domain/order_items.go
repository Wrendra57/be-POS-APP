package domain

import (
	"github.com/google/uuid"
	"time"
)

type OrderItems struct {
	Id        int
	OrderId   uuid.UUID
	ProductId uuid.UUID
	Price     int
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
