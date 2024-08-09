package services

import (
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrandService interface {
	CreateBrand(ctx *fiber.Ctx, r webrequest.BrandCreateReq) domain.Brand
	ListBrand(ctx *fiber.Ctx, r webrequest.BrandGetRequest) []domain.Brand
}

type brandServiceImpl struct {
	BrandRepo repositories.BrandRepository
	DB        *pgxpool.Pool
	Validate  *validator.Validate
}

func NewBrandService(brandRepo repositories.BrandRepository, db *pgxpool.Pool, validate *validator.Validate) BrandService {
	return &brandServiceImpl{
		BrandRepo: brandRepo,
		DB:        db,
		Validate:  validate,
	}
}
func (s brandServiceImpl) CreateBrand(ctx *fiber.Ctx, r webrequest.BrandCreateReq) domain.Brand {
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	brand := domain.Brand{
		Name:        r.Name,
		Description: r.Description,
	}
	brand = s.BrandRepo.Insert(ctx.Context(), tx, brand)

	return brand
}
func (s brandServiceImpl) ListBrand(ctx *fiber.Ctx, r webrequest.BrandGetRequest) []domain.Brand {
	//start db tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	brandReq := webrequest.BrandGetRequest{
		Params: r.Params,
		Limit:  r.Limit,
		Offset: (r.Offset - 1) * r.Limit,
	}

	brands := s.BrandRepo.ListAll(ctx, tx, brandReq)
	if len(brands) == 0 {
		brands = []domain.Brand{}
	}

	return brands
}
