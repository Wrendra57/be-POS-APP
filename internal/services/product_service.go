package services

import (
	"errors"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type ProductService interface {
	CreateProduct(ctx *fiber.Ctx, request webrequest.ProductCreateRequest) domain.Product
	FindProductById(ctx *fiber.Ctx, id uuid.UUID) (webrespones.ProductFindByIdResponseApi, exception.CustomEror, bool)
	ListProduct(ctx *fiber.Ctx, request webrequest.ProductListRequest) []domain.ProductList
	DeleteProduct(ctx *fiber.Ctx, id uuid.UUID) error
}

type productServiceImpl struct {
	ProductRepository repositories.ProductRepository
	PhotosService     PhotosService
	DB                *pgxpool.Pool
	Validate          *validator.Validate
	RedisDB           *redis.Client
}

func NewProductService(db *pgxpool.Pool,
	validate *validator.Validate, rdb *redis.Client, productRepo repositories.ProductRepository, photoService PhotosService) ProductService {
	return &productServiceImpl{
		DB:                db,
		Validate:          validate,
		RedisDB:           rdb,
		ProductRepository: productRepo,
		PhotosService:     photoService,
	}
}

func (s productServiceImpl) CreateProduct(ctx *fiber.Ctx, request webrequest.ProductCreateRequest) domain.Product {
	//TODO implement me
	adminId, _ := ctx.Locals("user_id").(uuid.UUID)
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())

	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	product := domain.Product{
		ProductName: request.ProductName,
		SellPrice:   request.SellPrice,
		CallName:    request.CallName,
		AdminId:     adminId,
		CategoryId:  request.Category.Id,
		BrandId:     request.Brand.Id,
		SupplierId:  request.Supplier.Id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	product, err = s.ProductRepository.Insert(ctx.Context(), tx, product)
	utils.PanicIfError(err)

	var photos []domain.Photos

	for _, photo := range request.Photo {
		p, _ := s.PhotosService.UploadPhoto(ctx, tx, request.ProductName, photo, product.Id)
		photos = append(photos, p)
	}

	return product
}

func (s productServiceImpl) FindProductById(ctx *fiber.Ctx, id uuid.UUID) (webrespones.ProductFindByIdResponseApi, exception.CustomEror, bool) {
	//TODO implement me
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	product := webrespones.ProductFindDetail{}
	product, err = s.ProductRepository.FindByIdDetail(ctx.Context(), tx, id)

	if err != nil {
		return webrespones.ProductFindByIdResponseApi{}, exception.CustomEror{Code: fiber.StatusNotFound, Error: err.Error()}, false
	}

	return webrespones.ProductFindByIdResponseApi{
		Id:          product.Id,
		ProductName: product.ProductName,
		SellPrice:   product.SellPrice,
		CallName:    product.CallName,
		Admin: struct {
			AdminId   uuid.UUID `json:"admin_id"`
			AdminName string    `json:"admin_name"`
		}{AdminId: product.AdminId, AdminName: product.AdminName},
		Category: struct {
			CategoryId          uuid.UUID `json:"category_id"`
			CategoryName        string    `json:"category_name"`
			CategoryDescription string    `json:"category_description"`
		}{CategoryId: product.CategoryId, CategoryName: product.CategoryName, CategoryDescription: product.CategoryDescription},
		Brand: struct {
			BrandId          int    `json:"brand_id"`
			BrandName        string `json:"brand_name"`
			BrandDescription string `json:"brand_description"`
		}{BrandId: product.BrandId, BrandName: product.BrandName, BrandDescription: product.BrandDescription},
		Supplier: struct {
			SupplierId          uuid.UUID `json:"supplier_id"`
			SupplierName        string    `json:"supplier_name"`
			SupplierContactInfo string    `json:"supplier_contact_info"`
			SupplierAddress     string    `json:"supplier_address"`
		}{SupplierId: product.SupplierId, SupplierName: product.SupplierName, SupplierContactInfo: product.SupplierContactInfo, SupplierAddress: product.SupplierAddress},
		Photos:    product.Photos,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, exception.CustomEror{}, true

}

func (s productServiceImpl) ListProduct(ctx *fiber.Ctx, request webrequest.ProductListRequest) []domain.ProductList {
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	request.Offset = (request.Offset - 1) * request.Limit
	products := s.ProductRepository.ListAll(ctx.Context(), tx, request)

	if len(products) == 0 {
		products = []domain.ProductList{}
	}
	return products

}

func (s productServiceImpl) DeleteProduct(ctx *fiber.Ctx, id uuid.UUID) error {
	//TODO implement me
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	product, err := s.ProductRepository.FindById(ctx.Context(), tx, id)
	if err != nil {
		return err
	}
	if product.DeletedAt.Valid == true {
		return errors.New("product not found / was deleted")
	}

	_ = s.ProductRepository.Delete(ctx.Context(), tx, id)

	return nil

}
