package repositories

import (
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	InsertUser(ctx *fiber.Ctx, tx pgx.Tx, user domain.User) (domain.User, error)
	FindByID(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.User, error)
}

type userRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) InsertUser(ctx *fiber.Ctx, tx pgx.Tx, user domain.User) (domain.User, error) {
	SQL := "INSERT INTO users(name, gender, telp, birthdate, address) VALUES($1, $2, $3, $4, $5) RETURNING user_id"

	var userID uuid.UUID // Variabel untuk menyimpan user_id yang dikembalikan

	// Eksekusi query dengan menggunakan QueryRow
	err := tx.QueryRow(ctx.Context(), SQL, user.Name, user.Gender, user.Telp, user.Birthday, user.Address).Scan(&userID)

	// Periksa apakah terjadi kesalahan
	if err != nil {
		return user, err // Mengembalikan kesalahan yang terjadi
	}

	// Set user_id yang telah didapat dari hasil query ke objek user
	user.User_id = userID
	fmt.Println(user)
	return user, nil
}

func (r *userRepositoryImpl) FindByID(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.User, error) {
	SQL := "select user_id, name, gender, telp, birthdate,address,created_at,updated_at from users where user_id= ?"

	rows, err := tx.Query(ctx.Context(), SQL, uuid)
	utils.PanicIfError(ctx, fiber.StatusInternalServerError, err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(&user.User_id, &user.Name, &user.Gender, &user.Telp, &user.Birthday, &user.Address,
			&user.Created_at, &user.Updated_at)
		utils.PanicIfError(ctx, fiber.StatusInternalServerError, err)
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}
