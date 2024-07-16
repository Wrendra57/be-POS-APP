package services

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type OTPService interface {
	CreateOTP(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.OTP, exception.CustomEror, error)
	ValidateOtp(ctx *fiber.Ctx, otp string) (exception.CustomEror, bool)
}

type otpServiceImpl struct {
	OTPRepository repositories.OtpRepository
	UserRepo      repositories.UserRepository
	OauthRepo     repositories.OauthRepository
	DB            *pgxpool.Pool
	Validate      *validator.Validate
}

func NewOTPService(OauthRepo repositories.OauthRepository, userRepo repositories.UserRepository, db *pgxpool.Pool,
	validate *validator.Validate, otpRepo repositories.OtpRepository) OTPService {
	return &otpServiceImpl{
		OTPRepository: otpRepo,
		UserRepo:      userRepo,
		OauthRepo:     OauthRepo,
		DB:            db,
		Validate:      validate,
	}
}

func (s *otpServiceImpl) CreateOTP(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.OTP, exception.CustomEror, error) {
	//TODO implement me

	panic("implement me")
}

func (s *otpServiceImpl) ValidateOtp(ctx *fiber.Ctx, o string) (exception.CustomEror, bool) {
	user_id, _ := ctx.Locals("user_id").(uuid.UUID)
	now := time.Now()
	// start database tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	//validate user id
	user, err := s.OauthRepo.FindByUUID(ctx, tx, user_id)
	if err != nil {
		return exception.CustomEror{Code: 404, Error: err.Error()}, false
	}
	if user.Is_enabled == true {
		return exception.CustomEror{Code: 400, Error: "Account is already enabled"}, false
	}

	//	get otp
	otp, err := s.OTPRepository.FindByUUID(ctx, tx, user_id)
	if err != nil {
		return exception.CustomEror{Code: 404, Error: "Code Otp Was Expired"}, false
	}
	if otp.Expired_date.Before(now) {
		return exception.CustomEror{Code: 404, Error: "Code Otp Was Expired"}, false
	}
	//	compare otp
	if otp.Otp != o {
		return exception.CustomEror{Code: 404, Error: "Code Otp was wrong"}, false
	}
	//	update user enabled
	fmt.Println("sini")
	_, err = s.OauthRepo.Update(ctx, tx, domain.Oauth{Is_enabled: true}, user_id)
	utils.PanicIfError(err)

	return exception.CustomEror{}, true
}
