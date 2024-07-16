package services

import (
	"errors"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserService interface {
	CreateUser(ctx *fiber.Ctx, request webrequest.UserCreateRequest) (string, exception.CustomEror,
		error)
}

type userServiceImpl struct {
	UserRepository  repositories.UserRepository
	OauthRepository repositories.OauthRepository
	OtpRepository   repositories.OtpRepository
	DB              *pgxpool.Pool
	Validate        *validator.Validate
}

func NewUserService(db *pgxpool.Pool,
	validate *validator.Validate, userRepo repositories.UserRepository,
	oauthRepo repositories.OauthRepository, otpRepo repositories.OtpRepository) UserService {
	return &userServiceImpl{
		UserRepository:  userRepo,
		OauthRepository: oauthRepo,
		OtpRepository:   otpRepo,
		DB:              db,
		Validate:        validate,
	}
}

func (s userServiceImpl) CreateUser(ctx *fiber.Ctx, request webrequest.UserCreateRequest) (string,
	exception.CustomEror, error) {

	// start database tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	_, err = s.OauthRepository.FindByEmail(ctx, tx, request.Email)
	if err == nil {
		return "", exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "email already exists"}, errors.New("Email already exist")
	}

	_, err = s.OauthRepository.FindByUserName(ctx, tx, request.Username)
	if err == nil {
		return "", exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "Username already exists"}, errors.New("Username already exist")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return "", exception.CustomEror{Code: fiber.StatusInternalServerError,
			Error: "Error hashing password "}, err
	}

	user := domain.User{
		Name:       request.Name,
		Gender:     request.Gender,
		Telp:       request.Telp,
		Birthday:   request.BirthdayConversed,
		Address:    request.Address,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	oauth := domain.Oauth{
		Email:      request.Email,
		Password:   hashedPassword,
		Is_enabled: false,
		Username:   request.Username,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	user, err = s.UserRepository.InsertUser(ctx, tx, user)
	utils.PanicIfError(err)

	oauth.User_id = user.User_id

	oauth, err = s.OauthRepository.InsertOauth(ctx, tx, oauth)
	utils.PanicIfError(err)

	//Creting OTP
	otp := domain.OTP{Otp: utils.GenerateOTP(), User_id: user.User_id, Expired_date: time.Now().Add(time.Minute * 3),
		Created_at: time.Now(), Updated_at: time.Now()}

	otp, err = s.OtpRepository.Insert(ctx, tx, otp)
	utils.PanicIfError(err)

	//strOTP := "ini kode token kamu " + otp.Otp
	//err = utils.SendEmail("wrendra57@gmail.com", "OTP-ACCOUNT", strOTP)
	//utils.PanicIfError(err)

	//GenerateJWT for access validasi otp
	JWTStr, err := utils.GenerateJWT(user.User_id)
	utils.PanicIfError(err)

	return JWTStr, exception.CustomEror{}, nil
}
