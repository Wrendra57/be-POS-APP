package services

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
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
	ValidateOtpAccount(ctx *fiber.Ctx, otp webrequest.ValidateOtpRequest) (exception.CustomEror, bool)
	ReSendOtp(ctx *fiber.Ctx, token string) (exception.CustomEror, bool)
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
	defer utils.CommitOrRollback(ctx.Context(), tx)

	//create otp
	otp := domain.OTP{
		Otp:          utils.GenerateOTP(),
		User_id:      u,
		Expired_date: now.Add(time.Minute * 2),
		Created_at:   now,
		Updated_at:   now,
	}
	otp, err = s.OTPRepository.Insert(ctx.Context(), tx, otp)
	if err != nil {
		return domain.OTP{}, exception.CustomEror{Code: 500, Error: "Internal Server Error"}, false
	}
	return otp, exception.CustomEror{}, true
}

func (s *otpServiceImpl) ValidateOtpAccount(ctx *fiber.Ctx, o webrequest.ValidateOtpRequest) (exception.CustomEror, bool) {
	parsedToken, err := utils.ParseJWT(o.Token)
	if err != nil {
		fmt.Println(err)
		return exception.CustomEror{Code: 401, Error: "Unauthorized"}, false
	}
	now := time.Now()

	// start database tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	//validate user id
	user, err := s.OauthRepo.FindByUUID(ctx, tx, parsedToken.User_id)

	if err != nil {
		return exception.CustomEror{Code: 404, Error: err.Error()}, false
	}
	if user.Is_enabled == true {
		return exception.CustomEror{Code: 400, Error: "Account is already enabled"}, false
	}
	//get otp
	otp, err := s.OTPRepository.FindByUUID(ctx.Context(), tx, parsedToken.User_id)
	if err != nil {
		return exception.CustomEror{Code: 404, Error: "Code Otp Was Not Found"}, false
	}
	if otp.Expired_date.Before(now) {
		return exception.CustomEror{Code: 404, Error: "Code Otp Was Expired"}, false
	}
	//compare otp
	if otp.Otp != o.Otp {
		return exception.CustomEror{Code: 404, Error: "Code Otp was wrong"}, false
	}

	//update user enabled
	_, err = s.OauthRepo.Update(ctx.Context(), tx, domain.Oauth{Is_enabled: true}, parsedToken.User_id)
	utils.PanicIfError(err)

	return exception.CustomEror{}, true
}

func (s *otpServiceImpl) ReSendOtp(ctx *fiber.Ctx, token string) (exception.CustomEror, bool) {
	parsedToken, err := utils.ParseJWT(token)
	if err != nil {
		return exception.CustomEror{Code: fiber.StatusUnauthorized, Error: "Unauthorized"}, false
	}
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	oauth, err := s.OauthRepo.FindByUUID(ctx, tx, parsedToken.User_id)
	if err != nil {
		fmt.Println(err)
		return exception.CustomEror{Code: fiber.StatusNotFound, Error: "Account not found"}, false
	}
	if oauth.Is_enabled == true {
		return exception.CustomEror{Code: fiber.StatusBadRequest, Error: "Account is already enabled"}, false
	}

	//cek dalam 5 menit resend otp berapa kali
	otps, err := s.OTPRepository.FindAllByUserIdAroundTime(ctx, tx, time.Now().Add(time.Minute*5*-1), time.Now(),
		parsedToken.User_id)
	utils.PanicIfError(err)

	if len(otps) >= 5 {
		msgStr := "Account max resend 5 OTP in 5 minute wait for " + otps[0].Created_at.Add(5*time.Minute).Format("15:04:05 02 Jan 2006")
		return exception.CustomEror{Code: fiber.StatusBadRequest, Error: msgStr}, false
	}

	otp, errS, e := s.CreateOTP(ctx, parsedToken.User_id)
	if e != true {
		fmt.Println(e)
		return errS, false
	}

	//sending email code
	strOTP := "Halo " + oauth.Username + ",,, ini kode Otp kamu " + otp.Otp + " expired in " + otp.Expired_date.String()
	fmt.Println(strOTP)
	//err = utils.SendEmail(oauth.Email, "OTP-ACCOUNT", strOTP)
	//utils.PanicIfError(err)

	return exception.CustomEror{}, true
}
