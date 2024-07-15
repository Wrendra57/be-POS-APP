package domain

import (
	"github.com/google/uuid"
	"time"
)

type Oauth struct {
	Id         int
	Email      string
	Password   string
	Is_enabled bool
	Username   string
	User_id    uuid.UUID
	Created_at time.Time
	Updated_at time.Time
}
