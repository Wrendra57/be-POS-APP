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
	FindUserDetail(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.UserDetail, error)
}

type userRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) InsertUser(ctx *fiber.Ctx, tx pgx.Tx, user domain.User) (domain.User, error) {
	SQL := "INSERT INTO users(name, gender, telp, birthdate, address) VALUES($1, $2, $3, $4, $5) RETURNING user_id"

	var userID uuid.UUID

	row := tx.QueryRow(ctx.Context(), SQL, user.Name, user.Gender, user.Telp, user.Birthday, user.Address)

	err := row.Scan(&userID)

	if err != nil {
		fmt.Println("repo insert user ==>  " + err.Error())
		return user, err
	}

	user.User_id = userID
	fmt.Println(user)
	return user, nil
}

func (r *userRepositoryImpl) FindByID(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.User, error) {
	SQL := "select user_id, name, gender, telp, birthdate,address,created_at," +
		"updated_at from users where user_id= $1 and deleted_at is null"

	rows, err := tx.Query(ctx.Context(), SQL, uuid)
	utils.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(&user.User_id, &user.Name, &user.Gender, &user.Telp, &user.Birthday, &user.Address,
			&user.Created_at, &user.Updated_at)
		utils.PanicIfError(err)
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}

func (r *userRepositoryImpl) FindUserDetail(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.UserDetail, error) {
	SQL := `SELECT u.user_id   AS user_id,
				   o.email     AS email,
				   o.username  AS username,
				   u.name      AS name,
				   u.gender    AS gender,
				   u.telp      AS telp,
				   u.birthdate AS birthdate,
				   u.address   AS address,
				   p.url       AS foto_profil,
				   r.role      AS role,
				   u.created_at
			FROM users u
					 JOIN oauths o on o.user_id = u.user_id
					 JOIN photos p ON u.user_id = p.owner_id
					 JOIN roles r on u.user_id = r.user_id
			where u.user_id = $1
			  AND o.is_enabled = true
			  AND o.deleted_at is null`

	rows, err := tx.Query(ctx.Context(), SQL, uuid)
	utils.PanicIfError(err)
	defer rows.Close()

	user := domain.UserDetail{}
	if rows.Next() {
		err := rows.Scan(&user.User_id, &user.Email, &user.Username, &user.Name, &user.Gender, &user.Telp, &user.Birthday, &user.Address,
			&user.Foto_profile, &user.Role, &user.Created_at)
		utils.PanicIfError(err)
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}
