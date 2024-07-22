package services

import (
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type CategoryService interface {
	CreateCategory(ctx *fiber.Ctx, r webrequest.CategoryCreateReq) (domain.Category, exception.CustomEror, bool)
	FindByParamsCategory(ctx *fiber.Ctx, r webrequest.CategoryFindByParam) ([]domain.Category, exception.CustomEror,
		bool)
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

func (s categoryServiceImpl) CreateCategory(ctx *fiber.Ctx, r webrequest.CategoryCreateReq) (domain.Category, exception.CustomEror, bool) {

	//start db tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	category := domain.Category{
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	category, err = s.CategoryRepo.Insert(ctx, tx, category)
	if err != nil {
		return domain.Category{}, exception.CustomEror{Code: 500, Error: "Invalid Server Error"}, false
	}
	return category, exception.CustomEror{}, true
}

func (s categoryServiceImpl) FindByParamsCategory(ctx *fiber.Ctx, r webrequest.CategoryFindByParam) ([]domain.Category, exception.CustomEror, bool) {
	//TODO implement me
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	c := webrequest.CategoryFindByParam{
		Params: r.Params,
		Limit:  r.Limit,
		Offset: (r.Offset - 1) * r.Limit,
	}
	categories := s.CategoryRepo.FindByParams(ctx, tx, c)

	return categories, exception.CustomEror{}, true
}
