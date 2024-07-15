package main

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/routes"
	"github.com/Wrendra57/Pos-app-be/internal/services"
	"github.com/Wrendra57/Pos-app-be/pkg"
	db2 "github.com/Wrendra57/Pos-app-be/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

func main() {
	fmt.Println("server starting")
	err := godotenv.Load()
	//utils.PanicIfError(ctx, fiber.StatusBadRequest, err)
	if err != nil {

		panic(err)
	}
	config.InitConfig()

	viper.SetConfigFile(".env")
	DB, cleanUp, err := db2.NewDatabase()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	//utils.PanicIfError(err)
	defer DB.Close()
	defer cleanUp()

	validate := pkg.NewValidate()
	userRepo := repositories.NewUserRepository()
	oauthRepo := repositories.NewOauthRepository()
	userService := services.NewUserService(DB, validate, userRepo, oauthRepo)

	app := fiber.New()
	fmt.Println("applying cors")
	app.Use(cors.New())
	app.Use(recover2.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to Pos App Be"))
	})
	api := app.Group("/api")
	routes.UserRoutes(api, userService, validate)

	fmt.Println("Server Ready")

	log.Fatal(app.Listen(":8080"))
}
