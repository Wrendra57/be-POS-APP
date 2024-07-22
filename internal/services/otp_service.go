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
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type OTPService interface {
	CreateOTP(ctx *fiber.Ctx, uuid uuid.UUID) (domain.OTP, exception.CustomEror, bool)
	ValidateOtpAccount(ctx *fiber.Ctx, otp string) (exception.CustomEror, bool)
	ReSendOtp(ctx *fiber.Ctx, userId uuid.UUID) (domain.OTP, exception.CustomEror, bool)
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

func (s *otpServiceImpl) CreateOTP(ctx *fiber.Ctx, u uuid.UUID) (domain.OTP, exception.CustomEror, bool) {
	now := time.Now()

	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	//create otp
	otp := domain.OTP{
		Otp:          utils.GenerateOTP(),
		User_id:      u,
		Expired_date: now.Add(time.Minute * 2),
		Created_at:   now,
		Updated_at:   now,
	}
	otp, err = s.OTPRepository.Insert(ctx, tx, otp)
	if err != nil {
		return domain.OTP{}, exception.CustomEror{Code: 500, Error: "Internal Server Error"}, false
	}
	return otp, exception.CustomEror{}, true
}

func (s *otpServiceImpl) ValidateOtpAccount(ctx *fiber.Ctx, o string) (exception.CustomEror, bool) {
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

func (s *otpServiceImpl) ReSendOtp(ctx *fiber.Ctx, userId uuid.UUID) (domain.OTP, exception.CustomEror, bool) {
	//TODO implement me
	// cek user

	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	oauth, err := s.OauthRepo.FindByUUID(ctx, tx, userId)
	if err != nil {
		fmt.Println(err)
		return domain.OTP{}, exception.CustomEror{Code: 400, Error: "Account not found"}, false
	}
	if oauth.Is_enabled == true {
		return domain.OTP{}, exception.CustomEror{Code: 400, Error: "Account is already enabled"}, false
	}
	//cek dalam 5 menit resend otp berapa kali

	otps, err := s.OTPRepository.FindAllByUserIdAroundTime(ctx, tx, time.Now().Add(time.Minute*5*-1), time.Now(), userId)
	utils.PanicIfError(err)
	fmt.Println("len otps", len(otps))
	//fmt.Println(otps)
	if len(otps) >= 5 {
		msgStr := "Account max resend 5 OTP in 5 minute wait for " + otps[0].Created_at.Add(5*time.Minute).Format("15:04:05 02 Jan 2006")
		return domain.OTP{}, exception.CustomEror{Code: 400, Error: msgStr}, false
	}
	fmt.Println("k")
	otp, errS, e := s.CreateOTP(ctx, userId)
	if e != true {
		fmt.Println(e)
		return domain.OTP{}, errS, false
	}

	//sending email code
	strOTP := "Halo " + oauth.Username + ",,, ini kode Otp kamu " + otp.Otp + " expired in " + otp.Expired_date.String()
	fmt.Println(strOTP)
	//err = utils.SendEmail(oauth.Email, "OTP-ACCOUNT", strOTP)
	//utils.PanicIfError(err)

	return domain.OTP{}, exception.CustomEror{}, true

}
