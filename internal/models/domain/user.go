package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	User_id    uuid.UUID
	Name       string
	Gender     string
	Telp       string
	Birthday   time.Time
	Address    string
	Created_at time.Time
	Updated_at time.Time
	Deleted_at sql.NullTime
}
