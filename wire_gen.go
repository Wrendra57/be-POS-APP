// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package be

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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Injectors from wire.go:

func InitializeApp() (*fiber.App, func(), error) {
	pool, cleanup, err := db.NewDatabase()
	if err != nil {
		return nil, nil, err
	}
	client := db.SetupRedis1()
	validate := pkg.NewValidate()
	userRepository := repositories.NewUserRepository()
	oauthRepository := repositories.NewOauthRepository()
	otpRepository := repositories.NewOtpRepository()
	roleRepository := repositories.NewRoleRepository()
	photosRepository := repositories.NewPhotosRepository()
	userService := services.NewUserService(pool, validate, client, userRepository, oauthRepository, otpRepository, roleRepository, photosRepository)
	otpService := services.NewOTPService(oauthRepository, userRepository, pool, validate, otpRepository)
	photosService := services.NewPhotosService(photosRepository, pool, validate)
	app, cleanup2, err := NewApp(pool, client, validate, userService, otpService, photosService)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

func NewApp(
	DB *pgxpool.Pool,
	RDB *redis.Client,
	validate *validator.Validate,
	userService services.UserService,
	otpService services.OTPService,
	photoService services.PhotosService,
) (*fiber.App, func(), error) {

	app := fiber.New()
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

	cleanup := func() {
		DB.Close()
		_ = RDB.Close()
	}

	return app, cleanup, nil
}
