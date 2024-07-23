package db

import (
	"context"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func NewDatabase() (*pgxpool.Pool, func(), error) {
	// Initialize database connection

	config.InitConfig()

	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")

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
