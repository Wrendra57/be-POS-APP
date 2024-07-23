package repositories

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type BrandRepository interface {
	Insert(ctx *fiber.Ctx, tx pgx.Tx, c domain.Brand) (domain.Brand, error)
	ListAll(ctx *fiber.Ctx, tx pgx.Tx, request webrequest.BrandGetRequest) []domain.Brand
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

func (r brandRepositoryImpl) ListAll(ctx *fiber.Ctx, tx pgx.Tx, request webrequest.BrandGetRequest) []domain.Brand {
	SQL := `SELECT id, name, description, created_at, updated_at, deleted_at
			FROM brands
			WHERE (name ILIKE $1 OR description ILIKE $2)
			  AND deleted_at IS NULL
			ORDER BY name ASC
			LIMIT $3 OFFSET $4`

	searchParams := "%" + request.Params + "%"
	rows, err := tx.Query(ctx.Context(), SQL, searchParams, searchParams, request.Limit, request.Offset)
	utils.PanicIfError(err)
	defer rows.Close()

	var brands []domain.Brand

	for rows.Next() {
		var b domain.Brand
		err := rows.Scan(&b.Id, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)
		utils.PanicIfError(err)
		brands = append(brands, b)
	}
	return brands
}
