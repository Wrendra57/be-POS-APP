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
)

type SupplierService interface {
	CreateSupplier(ctx *fiber.Ctx, r webrequest.SupplierRequest) (domain.Supplier, exception.CustomEror, bool)
	FindByParamSupplier(ctx *fiber.Ctx, request webrequest.SupplierListRequest) ([]domain.Supplier, bool)
}

type supplierServiceImpl struct {
	SupplierRepo repositories.SupplierRepository
	DB           *pgxpool.Pool
	Validate     *validator.Validate
}

func NewSupplierService(supplierRepo repositories.SupplierRepository, db *pgxpool.Pool, validate *validator.Validate) SupplierService {
	return &supplierServiceImpl{
		SupplierRepo: supplierRepo,
		DB:           db,
		Validate:     validate,
	}
}
func (s supplierServiceImpl) CreateSupplier(ctx *fiber.Ctx, r webrequest.SupplierRequest) (domain.Supplier, exception.CustomEror, bool) {
	//start db tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	supplier := domain.Supplier{
		Name:        r.Name,
		ContactInfo: r.ContactInfo,
		Address:     r.Address,
	}
	supplier, err = s.SupplierRepo.Insert(ctx.Context(), tx, supplier)
	if err != nil {
		return domain.Supplier{}, exception.CustomEror{Code: 400, Error: err.Error()}, false
	}
	return supplier, exception.CustomEror{}, true

}
func (s supplierServiceImpl) FindByParamSupplier(ctx *fiber.Ctx, request webrequest.SupplierListRequest) ([]domain.Supplier, bool) {
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	supplier := webrequest.SupplierListRequest{
		Params: request.Params,
		Limit:  request.Limit,
		Offset: (request.Offset - 1) * request.Limit,
	}
	suppliers := s.SupplierRepo.ListAll(ctx.Context(), tx, supplier)
	if len(suppliers) == 0 {
		suppliers = []domain.Supplier{}
	}
	return suppliers, true
}
