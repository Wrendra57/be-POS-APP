package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type InventoryTransaction struct {
	Id              int64  `json:"id"`
	TransactionType string `json:"transaction_type"`
	PriceModal      int    `json:"price_modal"`
	TotalQuantity   int    `json:"total_quantity"`
	LastUpdated     time.Time
	TransactionDate time.Time
	AdminId         uuid.UUID
	SupplierId      uuid.UUID
	Keterangan      string
	ProductId       uuid.UUID
	DeletedAt       sql.NullTime
}
