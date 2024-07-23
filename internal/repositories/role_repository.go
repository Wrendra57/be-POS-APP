package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoleRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, roles domain.Roles) (domain.Roles, error)
	FindByUserId(ctx *fiber.Ctx, tx pgx.Tx, userId uuid.UUID) (domain.Roles, error)
}

type roleRepositoryImpl struct {
}

func NewRoleRepository() RoleRepository {
	return &roleRepositoryImpl{}
}

func (r roleRepositoryImpl) Insert(ctx context.Context, tx pgx.Tx, roles domain.Roles) (domain.Roles, error) {
	//TODO implement me
	SQL := "INSERT INTO roles(user_id, role) VALUES($1, $2) returning id"

	var id int
	row := tx.QueryRow(ctx, SQL, roles.User_id, roles.Role)

	err := row.Scan(&id)

	if err != nil {
		fmt.Println("repo insert role ==>  " + err.Error())
		return roles, err
	}
	roles.Id = id
	return roles, nil
}

func (r roleRepositoryImpl) FindByUserId(ctx *fiber.Ctx, tx pgx.Tx, userId uuid.UUID) (domain.Roles, error) {
	//TODO implement me
	SQL := "SELECT id, user_id, role, created_at, updated_at, deleted_at FROM roles WHERE user_id= $1 and deleted_at is null"

	rows, err := tx.Query(ctx.Context(), SQL, userId)
	utils.PanicIfError(err)
	defer rows.Close()

	role := domain.Roles{}

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.User_id, &role.Role, &role.Created_at, &role.Updated_at, &role.Deleted_at)
		utils.PanicIfError(err)

		return role, nil
	} else {
		return role, errors.New("role not found")
	}

	panic("implement me")
}
