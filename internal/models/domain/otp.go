package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type OTP struct {
	Id           int
	User_id      uuid.UUID
	Otp          string
	Expired_date time.Time
	Created_at   time.Time
	Updated_at   time.Time
	Deleted_at   sql.NullTime
}
