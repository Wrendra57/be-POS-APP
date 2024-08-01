package main

import (
	"github.com/Wrendra57/Pos-app-be/internal/middleware"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/routes"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/pkg"
	"github.com/Wrendra57/Pos-app-be/pkg/db"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitializeApp() (*fiber.App, func(), error) {
	wire.Build(
		db.NewDatabase,
		pkg.NewValidate,
		repositories.NewUserRepository,
		repositories.NewOauthRepository,
		repositories.NewOtpRepository,
		repositories.NewPhotosRepository,
		repositories.NewRoleRepository,
		repositories.NewCategoriesRepository,
		repositories.NewBrandRepository,
		repositories.NewSupplierRepository,
		repositories.NewProductRepository,
		services.NewUserService,
		services.NewOTPService,
		services.NewPhotosService,
		services.NewCategoryService,
		services.NewBrandService,
		services.NewSupplierService,
		services.NewProductService,
		NewApp,
		db.SetupRedis1,
	)

	return nil, nil, nil
}
func NewApp(
	DB *pgxpool.Pool,
	RDB *redis.Client,
	validate *validator.Validate,
	userService services.UserService,
	otpService services.OTPService,
	photoService services.PhotosService,
	categoryService services.CategoryService,
	brandService services.BrandService,
	supplierService services.SupplierService,
	productService services.ProductService,
) (*fiber.App, func(), error) {

	app := fiber.New(
		fiber.Config{
			BodyLimit: 30 * 1024 * 1024,
		},
	)
	app.Use(cors.New())
	app.Use(recover2.New())
	app.Use(middleware.RecoverMiddleware())
	app.Static("/foto", "./storage/photos")
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to Pos App Be"))
	})

	api := app.Group("/api")
	routes.UserRoutes(api, userService, validate)
	routes.OtpRoutes(api, otpService, validate)
	routes.FileRoutes(api, photoService, validate)
	routes.CategoriesRoutes(api, categoryService, validate)
	routes.BrandRoutes(api, brandService, validate)
	routes.SupplierRoutes(api, supplierService, validate)
	routes.ProductRoutes(api, productService, validate)

	cleanup := func() {
		DB.Close()
		_ = RDB.Close()
	}

	return app, cleanup, nil
}
