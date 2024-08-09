package brandTest

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertBrandTest(db *pgxpool.Pool, brand domain.Brand) domain.Brand {
	brandRepo := repositories.NewBrandRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	brand = brandRepo.Insert(context.Background(), tx, brand)
	utils.PanicIfError(err)
	return brand
}
