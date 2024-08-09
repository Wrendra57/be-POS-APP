package repositories

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

type CategoriesRespository interface {
	Insert(ctx context.Context, tx pgx.Tx, c domain.Category) domain.Category
	FindByParams(ctx *fiber.Ctx, tx pgx.Tx, s webrequest.CategoryFindByParam) []domain.Category
}

type categoriesRespositoryImpl struct {
}

func NewCategoriesRepository() CategoriesRespository {
	return &categoriesRespositoryImpl{}
}

func (r categoriesRespositoryImpl) Insert(ctx context.Context, tx pgx.Tx, c domain.Category) domain.Category {
	SQL := "INSERT INTO categories(name, description) VALUES ($1, $2) RETURNING id"

	var id uuid.UUID
	row := tx.QueryRow(ctx, SQL, c.Name, c.Description)

	err := row.Scan(&id)
	utils.PanicIfError(err)
	c.Id = id
	return c
}
func (r categoriesRespositoryImpl) FindByParams(ctx *fiber.Ctx, tx pgx.Tx,
	request webrequest.CategoryFindByParam) []domain.Category {

	SQL := `SELECT id, name, description, created_at, updated_at, deleted_at
			FROM categories
			WHERE (name ILIKE $1 OR description ILIKE $2)
			  AND deleted_at IS NULL
			ORDER BY name ASC
			LIMIT $3 OFFSET $4`

	searchParams := "%" + request.Params + "%"
	rows, err := tx.Query(ctx.Context(), SQL, searchParams, searchParams, request.Limit, request.Offset)
	utils.PanicIfError(err)
	defer rows.Close()

	var categories []domain.Category

	for rows.Next() {
		var c domain.Category
		err := rows.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
		utils.PanicIfError(err)
		categories = append(categories, c)
	}
	return categories
}
