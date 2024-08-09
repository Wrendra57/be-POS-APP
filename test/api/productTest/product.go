package productTest

import (
	"context"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
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
func DeleteProductTest(db *pgxpool.Pool, id uuid.UUID) error {

	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()
	SQL := `UPDATE products SET deleted_at = $1 WHERE id = $2`
	_, err = tx.Exec(context.Background(), SQL, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete tables: %w", err)
	}

	return nil
}
