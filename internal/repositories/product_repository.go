package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, product domain.Product) (domain.Product, error)
	FindById(ctx context.Context, tx pgx.Tx, id uuid.UUID) (webrespones.ProductFindDetail, error)
	ListAll(ctx context.Context, tx pgx.Tx, request webrequest.ProductListRequest) []domain.ProductList
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

func (p productRepositoryImpl) FindById(ctx context.Context, tx pgx.Tx, id uuid.UUID) (webrespones.ProductFindDetail, error) {
	SQL := `SELECT p.id           AS product_id,
				   p.product_name,
				   p.sell_price,
				   p.call_name,
				   u.user_id      AS admin_id,
				   u.name         AS admin_name,
				   c.id           AS category_id,
				   c.name         AS category_name,
				   c.description  AS category_description,
				   b.id           AS brand_id,
				   b.name         AS brand_name,
				   b.description  AS brand_description,
				   s.id           AS supplier_id,
				   s.name         AS supplier_name,
				   s.contact_info AS supplier_contact_info,
				   s.address      AS supplier_address,
				   p.created_at,
				   p.updated_at
			FROM products p
					 JOIN categories c ON c.id = p.category_id
					 JOIN brands b ON b.id = p.brand_id
					 JOIN suppliers s ON s.id = p.supplier_id
					 JOIN users u ON u.user_id = p.admin_id
			WHERE p.id = $1
			  AND p.deleted_at IS NULL;`

	rows, err := tx.Query(ctx, SQL, id)
	utils.PanicIfError(err)
	defer rows.Close()

	product := webrespones.ProductFindDetail{}

	if rows.Next() {
		err := rows.Scan(&product.Id, &product.ProductName, &product.SellPrice, &product.CallName, &product.AdminId, &product.AdminName, &product.CategoryId,
			&product.CategoryName, &product.CategoryDescription, &product.BrandId, &product.BrandName, &product.BrandDescription,
			&product.SupplierId, &product.SupplierName, &product.SupplierContactInfo,
			&product.SupplierAddress, &product.CreatedAt, &product.UpdatedAt)
		utils.PanicIfError(err)
		return product, nil
	} else {
		return product, errors.New("Product not found")
	}
}

func (p productRepositoryImpl) ListAll(ctx context.Context, tx pgx.Tx, request webrequest.ProductListRequest) []domain.ProductList {
	SQL := `SELECT p.id   AS product_id,
				   p.product_name,
				   p.sell_price,
				   p.call_name,
				   c.name AS category_name,
				   b.name AS brand_name,
				   p.created_at,
				   p.updated_at
			FROM products p
					 JOIN categories c ON c.id = p.category_id
					 JOIN suppliers s ON s.id = p.supplier_id
					 JOIN brands b ON b.id = p.brand_id
			WHERE (p.product_name ILIKE $1 OR p.call_name ILIKE $2)
			  AND p.deleted_at IS NULL
			ORDER BY p.product_name ASC
			LIMIT $3 OFFSET $4`
	searchParams := "%" + request.Params + "%"

	rows, err := tx.Query(ctx, SQL, searchParams, searchParams, request.Limit, request.Offset)
	utils.PanicIfError(err)
	defer rows.Close()

	var products []domain.ProductList
	for rows.Next() {
		product := domain.ProductList{}
		err := rows.Scan(&product.Id, &product.ProductName, &product.SellPrice, &product.CallName, &product.Category,
			&product.Brand, &product.CreatedAt, &product.UpdatedAt)
		utils.PanicIfError(err)
		products = append(products, product)
	}
	return products
}
