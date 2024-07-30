package repositories

import (
	"context"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, product domain.Product) (domain.Product, error)
}

type productRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &productRepositoryImpl{}
}
func (p productRepositoryImpl) Insert(ctx context.Context, tx pgx.Tx, product domain.Product) (domain.Product, error) {
	//TODO implement me
	SQL := `INSERT INTO products(product_name, sell_price, call_name, admin_id, category_id, brand_id, supplier_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`

	var id uuid.UUID
	row := tx.QueryRow(ctx, SQL, product.ProductName, product.SellPrice, product.CallName, product.AdminId, product.CategoryId, product.BrandId, product.SupplierId)

	err := row.Scan(&id)
	if err != nil {
		fmt.Println("repo insert product ==>  " + err.Error())
		return product, err
	}
	product.Id = id
	return product, nil
}
