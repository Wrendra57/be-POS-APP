package repositories

import (
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type OauthRepository interface {
	InsertOauth(ctx *fiber.Ctx, tx pgx.Tx, oauth domain.Oauth) (domain.Oauth, error)
	FindByEmail(ctx *fiber.Ctx, tx pgx.Tx, email string) (domain.Oauth, error)
	FindByUserName(ctx *fiber.Ctx, tx pgx.Tx, string2 string) (domain.Oauth, error)
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
		return oauth, err // Mengembalikan kesalahan yang terjadi
	}

	return oauth, nil
}

func (r *oauthRepositoryImpl) FindByEmail(ctx *fiber.Ctx, tx pgx.Tx, email string) (domain.Oauth,
	error) {
	fmt.Println("repo find by email")
	SQL := "SELECT id, email, password, is_enabled, username, user_id, created_at, updated_at FROM oauths WHERE email = $1"
	row := tx.QueryRow(ctx.Context(), SQL, email)

	oauth := domain.Oauth{}

	err := row.Scan(&oauth.Id, &oauth.Email, &oauth.Password, &oauth.Is_enabled, &oauth.Username, &oauth.User_id,
		&oauth.Created_at, &oauth.Updated_at)
	fmt.Println(err)
	if err != nil {
		return oauth, errors.New("user not found")
	}

	return oauth, nil
}
func (r *oauthRepositoryImpl) FindByUserName(ctx *fiber.Ctx, tx pgx.Tx, username string) (domain.Oauth,
	error) {
	fmt.Println("repo find by username")
	SQL := "SELECT id, email, password, is_enabled, username, user_id, created_at, updated_at FROM oauths WHERE username = $1"
	row := tx.QueryRow(ctx.Context(), SQL, username)

	oauth := domain.Oauth{}

	err := row.Scan(&oauth.Id, &oauth.Email, &oauth.Password, &oauth.Is_enabled, &oauth.Username, &oauth.User_id,
		&oauth.Created_at, &oauth.Updated_at)
	//fmt.Println(oauth)
	if err != nil {
		return oauth, errors.New("user not found")

		//if err != nil {
		//	return oauth, errors.New("user not found")
		//} else {
		//	return oauth, err // Mengembalikan kesalahan lain jika terjadi
		//}
	}

	return oauth, nil
}
