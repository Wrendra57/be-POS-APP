package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Photos struct {
	Id        int
	Url       string
	Owner     uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
