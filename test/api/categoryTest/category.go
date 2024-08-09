package categoryTests

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertCategoriesTest(db *pgxpool.Pool, category domain.Category) domain.Category {
	categoryRepo := repositories.NewCategoriesRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	category = categoryRepo.Insert(context.Background(), tx, category)
	utils.PanicIfError(err)
	return category
}
