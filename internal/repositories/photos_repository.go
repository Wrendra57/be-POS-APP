package repositories

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PhotosRepository interface {
	Insert(ctx *fiber.Ctx, tx pgx.Tx, photos domain.Photos) (domain.Photos, error)
}

type PhotosRepositoryImpl struct {
}

func NewPhotosRepository() PhotosRepository {
	return &PhotosRepositoryImpl{}
}

func (p PhotosRepositoryImpl) Insert(ctx *fiber.Ctx, tx pgx.Tx, photos domain.Photos) (domain.Photos, error) {
	//TODO implement me
	SQL := "INSERT INTO photos(url,owner_id) VALUES($1, $2) RETURNING id"

	var id uuid.UUID

	err := tx.QueryRow(ctx.Context(), SQL, photos.Url, photos.Owner).Scan(&id)
	if err != nil {
		fmt.Println("repo insert user ==>  " + err.Error())
		return photos, err
	}

	photos.Id = id
	return photos, nil
}
