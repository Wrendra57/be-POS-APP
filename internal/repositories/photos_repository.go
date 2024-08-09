package repositories

import (
	"context"
	"errors"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PhotosRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, photos domain.Photos) domain.Photos
	FindByUUID(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.Photos, error)
}

type PhotosRepositoryImpl struct {
}

func NewPhotosRepository() PhotosRepository {
	return &PhotosRepositoryImpl{}
}

func (p PhotosRepositoryImpl) Insert(ctx context.Context, tx pgx.Tx, photos domain.Photos) domain.Photos {
	//TODO implement me
	SQL := "INSERT INTO photos(url,owner_id) VALUES($1, $2) RETURNING id"

	var id int

	row := tx.QueryRow(ctx, SQL, photos.Url, photos.Owner)

	err := row.Scan(&id)

	utils.PanicIfError(err)

	photos.Id = id
	return photos
}

func (p PhotosRepositoryImpl) FindByUUID(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.Photos, error) {
	//TODO implement me
	SQL := "SELECT id, url, owner_id FROM photos WHERE owner_id = $1 AND deleted_at is null"

	rows, err := tx.Query(ctx.Context(), SQL, uuid)
	utils.PanicIfError(err)
	defer rows.Close()

	photos := domain.Photos{}
	if rows.Next() {
		err := rows.Scan(&photos.Id, &photos.Url, &photos.Owner)
		utils.PanicIfError(err)
		return photos, nil
	} else {
		return photos, errors.New("user not found")
	}
}
