package otptest

import (
	"context"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func FindOtpRepo(db *pgxpool.Pool, userId uuid.UUID) domain.OTP {
	otpRepo := repositories.NewOtpRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	otp, err := otpRepo.FindByUUID(context.Background(), tx, userId)
	utils.PanicIfError(err)
	return otp
}
func UpdateOauthTest(db *pgxpool.Pool, oauth domain.Oauth) domain.Oauth {
	oauthRepo := repositories.NewOauthRepository()
	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	oauth, err = oauthRepo.Update(context.Background(), tx, oauth, oauth.User_id)
	utils.PanicIfError(err)
	return oauth
}
func TruncateOtp(db *pgxpool.Pool) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()
	SQL := `TRUNCATE TABLE otp RESTART IDENTITY CASCADE`
	_, err = tx.Exec(context.Background(), SQL)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}

func UpdateOtpExpired(db *pgxpool.Pool, userId uuid.UUID) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()
	SQL := `UPDATE otp SET expired_date = $1 WHERE user_id = $2`
	_, err = tx.Exec(context.Background(), SQL, time.Now().Add(time.Hour*-10), userId)
	if err != nil {
		return fmt.Errorf("failed to update tables: %w", err)
	}

	return nil
}
func InsertOtpTest(db *pgxpool.Pool, otp domain.OTP) error {
	otpRepo := repositories.NewOtpRepository()

	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	_ = otpRepo.Insert(context.Background(), tx, otp)
	utils.PanicIfError(err)
	return nil
}
