package test

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func SetupDBTest() (*pgxpool.Pool, func(), error) {
	dbHost := "127.0.0.1"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "password"
	dbName := "Pos_app"

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

func TruncateDB(db *pgxpool.Pool) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()
	SQL := `TRUNCATE TABLE brands, categories, devices, inventories, inventory_transactions, oauths, order_items, 
orders, otp, photos, products, roles, suppliers, users RESTART IDENTITY CASCADE`
	_, err = tx.Exec(context.Background(), SQL)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}
func InitConfigTest() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("../../..")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}
