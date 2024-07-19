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

type UserDetail struct {
	User_id      uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	Telp         string    `json:"telp"`
	Birthday     time.Time `json:"birthday"`
	Address      string    `json:"address"`
	Foto_profile string    `json:"foto_profile"`
	Role         string    `json:"role"`
	Created_at   time.Time `json:"created_at"`
}
