package services

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OTPService interface {
	CreateOTP(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.OTP, exception.CustomEror, error)
}

type otpServiceImpl struct {
	OTPRepository repositories.OtpRepository
	DB            *pgxpool.Pool
	Validate      *validator.Validate
}

func NewOTPService(db *pgxpool.Pool, validate *validator.Validate, otpRepo repositories.OtpRepository) OTPService {
	return &otpServiceImpl{
		OTPRepository: otpRepo,
		DB:            db,
		Validate:      validate,
	}
}

func (s *otpServiceImpl) CreateOTP(ctx *fiber.Ctx, tx pgx.Tx, uuid uuid.UUID) (domain.OTP, exception.CustomEror, error) {
	//TODO implement me

	panic("implement me")
}
