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
	"time"
)

type OtpRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, otp domain.OTP) (domain.OTP, error)
	FindByUUID(ctx context.Context, tx pgx.Tx, uuid uuid.UUID) (domain.OTP, error)
	FindAllByUserIdAroundTime(ctx *fiber.Ctx, tx pgx.Tx, timeStart, timeEnd time.Time, user_id uuid.UUID) ([]domain.OTP, error)
}

type OtpRepositoryImpl struct {
}

func NewOtpRepository() OtpRepository {
	return &OtpRepositoryImpl{}
}

func (OtpRepositoryImpl) Insert(ctx context.Context, tx pgx.Tx, o domain.OTP) (domain.OTP, error) {
	//TODO implement me
	SQL := "INSERT INTO otp(user_id, otp, expired_date, created_at, updated_at) VALUES($1, $2, $3, $4,$5) returning id"

	var id int

	row := tx.QueryRow(ctx, SQL, o.User_id, o.Otp, o.Expired_date, time.Now(), time.Now())

	err := row.Scan(&id)

	if err != nil {
		fmt.Println("repo insert otp ==>  " + err.Error())
		return o, err // Mengembalikan kesalahan yang terjadi
	}

	o.Id = id
	fmt.Println(o)
	return o, nil
}

func (OtpRepositoryImpl) FindByUUID(ctx context.Context, tx pgx.Tx, uuid uuid.UUID) (domain.OTP, error) {
	//TODO implement me
	SQL := "SELECT id, user_id, otp, expired_date, created_at, updated_at FROM otp WHERE user_id = $1 AND deleted_at is null order by otp.created_at desc LIMIT 1"

	rows, err := tx.Query(ctx, SQL, uuid)
	utils.PanicIfError(err)
	defer rows.Close()

	otp := domain.OTP{}

	if rows.Next() {
		err := rows.Scan(&otp.Id, &otp.User_id, &otp.Otp, &otp.Expired_date, &otp.Created_at, &otp.Updated_at)
		utils.PanicIfError(err)
		return otp, nil
	} else {
		return otp, errors.New("otp not found")
	}
}

func (OtpRepositoryImpl) FindAllByUserIdAroundTime(ctx *fiber.Ctx, tx pgx.Tx, timeStart, timeEnd time.Time, user_id uuid.UUID) ([]domain.OTP, error) {
	//TODO implement me
	SQL := "SELECT id, user_id, otp, expired_date, created_at, " +
		"updated_at FROM otp WHERE created_at >= $1 AND created_at <= $2 AND user_id = $3 AND deleted_at is null order by created_at ASC"

	rows, err := tx.Query(ctx.Context(), SQL, timeStart, timeEnd, user_id)
	utils.PanicIfError(err)
	defer rows.Close()

	var otps []domain.OTP
	for rows.Next() {
		otp := domain.OTP{}
		err := rows.Scan(&otp.Id, &otp.User_id, &otp.Otp, &otp.Expired_date, &otp.Created_at, &otp.Updated_at)
		utils.PanicIfError(err)
		otps = append(otps, otp)
	}
	return otps, nil

}
