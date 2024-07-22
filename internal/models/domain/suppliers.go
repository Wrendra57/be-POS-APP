package domain

import (
	"github.com/google/uuid"
	"time"
)

type Supplier struct {
	Id          uuid.UUID
	Name        string
	ContactInfo string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
