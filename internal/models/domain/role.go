package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Roles struct {
	Id         int
	User_id    uuid.UUID
	Role       string
	Created_at time.Time
	Updated_at time.Time
	Deleted_at sql.NullTime
}
