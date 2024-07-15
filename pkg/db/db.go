package db

import (
	"context"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func NewDB(ctx *fiber.Ctx) *pgxpool.Pool {
	config := viper.GetViper()

	// Construct DSN
	dsn := "postgres://" + config.GetString("DB_USER") + ":" + config.GetString("DB_PASSWORD") + ":" + config.GetString("DB_HOST") + ":" + config.GetString("DB_PORT") + "/" + config.GetString("DB_NAME") + "?sslmode=" + config.GetString("DB_SSL_MODE")
	db, err := pgxpool.New(context.Background(), dsn)
	utils.PanicIfError(ctx, fiber.StatusInternalServerError, err)

	defer db.Close()
	return db
}

func NewDatabase() (*pgxpool.Pool, func(), error) {
	// Initialize database connection
	fmt.Println("database starting")

	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	//dbSslMode := viper.GetString("DB_SSLMODE")

	// Membuat URL koneksi PostgreSQL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, nil, err
	}

	// Return a cleanup function to close the pool
	cleanup := func() {
		db.Close()
	}

	return db, cleanup, nil
}
