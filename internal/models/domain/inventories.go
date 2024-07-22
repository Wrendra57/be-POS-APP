package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Inventories struct {
	Id            int          `json:"id"`
	ProductId     uuid.UUID    `json:"product_id"`
	TotalQuantity int          `json:"total_quantity"`
	LastUpdate    time.Time    `json:"last_update"`
	DeletedAt     sql.NullTime `json:"deleted_at"`
}
