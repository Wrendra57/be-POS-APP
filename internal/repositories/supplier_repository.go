package repositories

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SupplierRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, supplier domain.Supplier) (domain.Supplier, error)
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
