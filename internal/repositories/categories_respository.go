package repositories

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CategoriesRespository interface {
	Insert(ctx *fiber.Ctx, tx pgx.Tx, c domain.Category) (domain.Category, error)
}

type categoriesRespositoryImpl struct {
}

func NewCategoriesRepository() CategoriesRespository {
	return &categoriesRespositoryImpl{}
}

func (r categoriesRespositoryImpl) Insert(ctx *fiber.Ctx, tx pgx.Tx, c domain.Category) (domain.Category, error) {
	SQL := "INSERT INTO categories(name, description) VALUES ($1, $2) RETURNING id"

	var id uuid.UUID
	row := tx.QueryRow(ctx.Context(), SQL, c.Name, c.Description)

	err := row.Scan(&id)
	if err != nil {
		fmt.Println("repo insert Category ==>  " + err.Error())
		return domain.Category{}, err
	}
	c.Id = id
	return c, nil
}
