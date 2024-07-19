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
		services.NewUserService,
		services.NewOTPService,
		services.NewPhotosService,
		NewApp,
	)

	return nil, nil, nil
}
func NewApp(
	DB *pgxpool.Pool,
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
	}

	return app, cleanup, nil
}
