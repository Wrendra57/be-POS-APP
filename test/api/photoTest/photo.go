package photoTest

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertPhotosTest(db *pgxpool.Pool, photo domain.Photos) domain.Photos {
	photoRepo := repositories.NewPhotosRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	photo, err = photoRepo.Insert(context.Background(), tx, photo)
	utils.PanicIfError(err)
	return photo
}
