package suplier

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertSupplierTest(db *pgxpool.Pool, supplier domain.Supplier) domain.Supplier {
	supplierRepo := repositories.NewSupplierRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	supplier, err = supplierRepo.Insert(context.Background(), tx, supplier)
	utils.PanicIfError(err)
	return supplier
}
