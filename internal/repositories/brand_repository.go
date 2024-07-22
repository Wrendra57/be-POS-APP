package repositories

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type BrandRepository interface {
	Insert(ctx *fiber.Ctx, tx pgx.Tx, c domain.Brand) (domain.Brand, error)
}

type brandRepositoryImpl struct {
}

func NewBrandRepository() BrandRepository {
	return &brandRepositoryImpl{}
}

func (r brandRepositoryImpl) Insert(ctx *fiber.Ctx, tx pgx.Tx, brand domain.Brand) (domain.Brand, error) {
	//TODO implement me
	SQL := "INSERT INTO brands(name, description) VALUES ($1, $2) RETURNING id"

	var id int
	row := tx.QueryRow(ctx.Context(), SQL, brand.Name, brand.Description)

	err := row.Scan(&id)
	if err != nil {
		return domain.Brand{}, err
	}
	brand.Id = id
	return brand, nil
}
