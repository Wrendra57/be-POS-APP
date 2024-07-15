package webrespones

import (
	"github.com/google/uuid"
	"time"
)

type UserDetail struct {
	User_id    uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	Is_enabled bool      `json:"is_enabled"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	Telp       string    `json:"telp"`
	Birthday   time.Time `json:"birthday"`
	Address    string    `json:"address"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
