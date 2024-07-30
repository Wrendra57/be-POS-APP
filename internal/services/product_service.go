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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type ProductService interface {
	CreateProduct(ctx *fiber.Ctx, request webrequest.ProductCreateRequest) (domain.Product, exception.CustomEror, bool)
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

func (s productServiceImpl) CreateProduct(ctx *fiber.Ctx, request webrequest.ProductCreateRequest) (domain.Product, exception.CustomEror, bool) {
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

	return product, exception.CustomEror{}, true
}
