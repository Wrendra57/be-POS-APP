package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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
	DB, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("ds")
	SQL := `TRUNCATE TABLE brands, categories,devices,inventories,inventory_transactions,oauths,order_items,orders,otp,photos,products,roles,suppliers,users RESTART IDENTITY CASCADE`
	_ = DB.QueryRow(context.Background(), SQL)
	fmt.Println("executed")
	DB.Commit(context.Background())
	return nil
}
