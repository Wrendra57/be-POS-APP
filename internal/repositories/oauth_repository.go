package repositories

import (
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type OauthRepository interface {
	InsertOauth(ctx *fiber.Ctx, tx pgx.Tx, oauth domain.Oauth) (domain.Oauth, error)
	FindByEmail(ctx *fiber.Ctx, tx pgx.Tx, email string) (domain.Oauth, error)
	FindByUserName(ctx *fiber.Ctx, tx pgx.Tx, string2 string) (domain.Oauth, error)
	FindByUUID(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.Oauth, error)
	Update(ctx *fiber.Ctx, tx pgx.Tx, oauth domain.Oauth, u uuid.UUID) (domain.Oauth, error)
}

type oauthRepositoryImpl struct {
}

func NewOauthRepository() OauthRepository {
	return &oauthRepositoryImpl{}
}

func (r *oauthRepositoryImpl) InsertOauth(ctx *fiber.Ctx, tx pgx.Tx, oauth domain.Oauth) (domain.Oauth, error) {
	SQL := "INSERT INTO oauths(email, password, username, user_id) VALUES($1, $2, $3, $4) RETURNING id"

	// Eksekusi query dengan QueryRow dan scan hasilnya ke oauth.Id
	err := tx.QueryRow(ctx.Context(), SQL, oauth.Email, oauth.Password, oauth.Username, oauth.User_id).Scan(&oauth.Id)

	// Periksa apakah terjadi kesalahan
	if err != nil {
		fmt.Println("insertoauth ==>  " + err.Error())

		return oauth, err // Mengembalikan kesalahan yang terjadi
	}

	return oauth, nil
}

func (r *oauthRepositoryImpl) FindByEmail(ctx *fiber.Ctx, tx pgx.Tx, email string) (domain.Oauth,
	error) {

	SQL := "SELECT id, email, password, is_enabled, username, user_id, created_at, updated_at FROM oauths WHERE email = $1"
	row := tx.QueryRow(ctx.Context(), SQL, email)

	oauth := domain.Oauth{}

	err := row.Scan(&oauth.Id, &oauth.Email, &oauth.Password, &oauth.Is_enabled, &oauth.Username, &oauth.User_id,
		&oauth.Created_at, &oauth.Updated_at)

	if err != nil {
		fmt.Println("repo find by username ==>  " + err.Error())
		return oauth, errors.New("user not found")
	}

	return oauth, nil
}
func (r *oauthRepositoryImpl) FindByUserName(ctx *fiber.Ctx, tx pgx.Tx, username string) (domain.Oauth,
	error) {

	SQL := "SELECT id, email, password, is_enabled, username, user_id, created_at, updated_at FROM oauths WHERE username = $1"
	row := tx.QueryRow(ctx.Context(), SQL, username)

	oauth := domain.Oauth{}

	err := row.Scan(&oauth.Id, &oauth.Email, &oauth.Password, &oauth.Is_enabled, &oauth.Username, &oauth.User_id,
		&oauth.Created_at, &oauth.Updated_at)
	//fmt.Println(oauth)
	if err != nil {
		fmt.Println("repo find by username ==>  " + err.Error())
		return oauth, errors.New("user not found")
	}

	return oauth, nil
}
func (r *oauthRepositoryImpl) FindByUUID(ctx *fiber.Ctx, tx pgx.Tx, u uuid.UUID) (domain.Oauth, error) {
	SQL := "SELECT id, email, password, is_enabled, username, user_id, created_at, updated_at FROM oauths WHERE user_id = $1"
	row := tx.QueryRow(ctx.Context(), SQL, u)

	oauth := domain.Oauth{}

	err := row.Scan(&oauth.Id, &oauth.Email, &oauth.Password, &oauth.Is_enabled, &oauth.Username, &oauth.User_id,
		&oauth.Created_at, &oauth.Updated_at)

	if err != nil {
		fmt.Println("repo find by uuid ==>  " + err.Error())
		return oauth, errors.New("user not found")
	}
	return oauth, nil
}

func (r *oauthRepositoryImpl) Update(ctx *fiber.Ctx, tx pgx.Tx, o domain.Oauth, u uuid.UUID) (domain.Oauth, error) {
	SQL := "UPDATE oauths SET "
	var args []interface{}
	var index int

	if o.Email != "" {
		index++
		SQL += fmt.Sprintf("email = $%d, ", index)
		args = append(args, o.Email)
	}
	if o.Username != "" {
		index++
		SQL += fmt.Sprintf("username = $%d, ", index)
		args = append(args, o.Username)
	}
	if o.Password != "" {
		index++
		SQL += fmt.Sprintf("password = $%d, ", index)
		args = append(args, o.Password)
	}
	index++
	SQL += fmt.Sprintf("is_enabled = $%d, ", index)
	args = append(args, true)

	index++
	SQL += fmt.Sprintf("updated_at = $%d, ", index)
	args = append(args, time.Now())

	SQL = SQL[:len(SQL)-2]

	//add where clause
	index++
	SQL += fmt.Sprintf(" WHERE user_id = $%d", index)
	args = append(args, u)

	// Execute the update query
	_, err := tx.Exec(ctx.Context(), SQL, args...)
	if err != nil {
		return domain.Oauth{}, fmt.Errorf("failed to update oauth: %w", err)
	}

	// Retrieve the updated row to return it
	row := tx.QueryRow(ctx.Context(), "SELECT id, email, password, is_enabled, username, user_id, created_at, "+
		"updated_at FROM oauths WHERE user_id = $1", u)

	var oauth domain.Oauth
	err = row.Scan(&oauth.Id, &oauth.Email, &oauth.Password, &oauth.Is_enabled, &oauth.Username, &oauth.User_id, &oauth.Created_at, &oauth.Updated_at)
	if err != nil {
		return domain.Oauth{}, fmt.Errorf("failed to retrieve updated oauth: %w", err)
	}

	return o, nil
}
