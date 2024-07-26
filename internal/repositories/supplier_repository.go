package repositories

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SupplierRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, supplier domain.Supplier) (domain.Supplier, error)
	ListAll(ctx context.Context, tx pgx.Tx, request webrequest.SupplierListRequest) []domain.Supplier
}

type supplierRepositoryImpl struct {
}

func NewSupplierRepository() SupplierRepository {
	return &supplierRepositoryImpl{}
}
func (s supplierRepositoryImpl) Insert(ctx context.Context, tx pgx.Tx, supplier domain.Supplier) (domain.Supplier, error) {
	//TODO implement me
	SQL := "INSERT INTO suppliers(name, contact_info, address) VALUES ($1, $2, $3) RETURNING id"

	var id uuid.UUID
	row := tx.QueryRow(ctx, SQL, supplier.Name, supplier.ContactInfo, supplier.Address)

	err := row.Scan(&id)
	if err != nil {
		return domain.Supplier{}, err
	}
	supplier.Id = id
	return supplier, nil
}
func (s supplierRepositoryImpl) ListAll(ctx context.Context, tx pgx.Tx, request webrequest.SupplierListRequest) []domain.Supplier {

	SQL := `SELECT id, name, contact_info, address, created_at, updated_at, deleted_at
			FROM suppliers
			WHERE (name ILIKE $1 OR contact_info ILIKE $2 OR address ILIKE $3)
			  AND deleted_at IS NULL
			ORDER BY name ASC
			LIMIT $4 OFFSET $5`
	searchParams := "%" + request.Params + "%"
	rows, err := tx.Query(ctx, SQL, searchParams, searchParams, searchParams, request.Limit, request.Offset)
	utils.PanicIfError(err)
	defer rows.Close()

	var suppliers []domain.Supplier
	for rows.Next() {
		var sup domain.Supplier
		err := rows.Scan(&sup.Id, &sup.Name, &sup.ContactInfo, &sup.Address, &sup.CreatedAt, &sup.UpdatedAt, &sup.DeletedAt)
		utils.PanicIfError(err)
		suppliers = append(suppliers, sup)
	}
	return suppliers
}
