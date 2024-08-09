package productTest

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertProductTest(db *pgxpool.Pool, product domain.Product) domain.Product {
	productRepo := repositories.NewProductRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	product, err = productRepo.Insert(context.Background(), tx, product)
	utils.PanicIfError(err)
	return product
}
