package services

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type CategoryService interface {
	CreateCategory(ctx *fiber.Ctx, r webrequest.CategoryCreateReq) domain.Category
	FindByParamsCategory(ctx *fiber.Ctx, r webrequest.CategoryFindByParam) []domain.Category
}

type categoryServiceImpl struct {
	CategoryRepo repositories.CategoriesRespository
	DB           *pgxpool.Pool
	Validate     *validator.Validate
}

func NewCategoryService(categoryRepo repositories.CategoriesRespository, db *pgxpool.Pool, validate *validator.Validate) CategoryService {
	return &categoryServiceImpl{
		CategoryRepo: categoryRepo,
		DB:           db,
		Validate:     validate,
	}
}

func (s categoryServiceImpl) CreateCategory(ctx *fiber.Ctx, r webrequest.CategoryCreateReq) domain.Category {

	//start db tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	category := domain.Category{
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	category = s.CategoryRepo.Insert(ctx.Context(), tx, category)
	return category
}

func (s categoryServiceImpl) FindByParamsCategory(ctx *fiber.Ctx, r webrequest.CategoryFindByParam) []domain.Category {
	//TODO implement me
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	c := webrequest.CategoryFindByParam{
		Params: r.Params,
		Limit:  r.Limit,
		Offset: (r.Offset - 1) * r.Limit,
	}

	categories := s.CategoryRepo.FindByParams(ctx, tx, c)
	fmt.Println("asd")
	if len(categories) == 0 {
		categories = []domain.Category{}
	}

	return categories
}
