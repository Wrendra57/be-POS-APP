package domain

import (
	"github.com/google/uuid"
	"time"
)

type Orders struct {
	Id        uuid.UUID
	AdminId   uuid.UUID
	BuyerId   uuid.UUID
	Total     int
	Tax       int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
