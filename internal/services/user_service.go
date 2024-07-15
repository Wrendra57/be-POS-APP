package services

import (
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/Wrendra57/Pos-app-be/internal/utils/template_response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserService interface {
	CreateUser(ctx *fiber.Ctx, request webrequest.UserCreateRequest) (webrespones.UserDetail, exception.CustomEror,
		error)
}

type userServiceImpl struct {
	UserRepository  repositories.UserRepository
	OauthRepository repositories.OauthRepository
	DB              *pgxpool.Pool
	Validate        *validator.Validate
}

func NewUserService(db *pgxpool.Pool,
	validate *validator.Validate, userRepo repositories.UserRepository,
	oauthRepo repositories.OauthRepository) UserService {
	return &userServiceImpl{
		UserRepository:  userRepo,
		OauthRepository: oauthRepo,
		DB:              db,
		Validate:        validate,
	}
}

func (s userServiceImpl) CreateUser(ctx *fiber.Ctx, request webrequest.UserCreateRequest) (webrespones.UserDetail,
	exception.CustomEror, error) {
	//fmt.Println(request)

	// start database
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(ctx, fiber.StatusBadRequest, err)
	defer utils.CommitOrRollback(ctx, tx)

	_, err = s.OauthRepository.FindByEmail(ctx, tx, request.Email)
	//fmt.Println(err)
	if err == nil {
		fmt.Println("err FindByEmail")
		return webrespones.UserDetail{}, exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "email already exists"}, errors.New("Email already exist")
	}

	_, err = s.OauthRepository.FindByUserName(ctx, tx, request.Username)

	if err == nil {
		fmt.Println("err2")
		fmt.Println("err23")
		fmt.Println(err)

		return webrespones.UserDetail{}, exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "Username already exists"}, errors.New("Username already exist")
	}
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		//utils.PanicIfError(ctx, fiber.StatusBadRequest, errors.New("Error hashing password"))
		return webrespones.UserDetail{}, exception.CustomEror{Code: fiber.StatusInternalServerError,
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
	//fmt.Println(user)
	//fmt.Println(oauth)
	user, err = s.UserRepository.InsertUser(ctx, tx, user)
	if err != nil {
		//utils.PanicIfError(ctx, fiber.StatusBadRequest, err)
		return webrespones.UserDetail{}, exception.CustomEror{Code: fiber.StatusInternalServerError,
			Error: "Internal server error"}, err
	}
	//utils.PanicIfError(err)
	oauth.User_id = user.User_id

	oauth, err = s.OauthRepository.InsertOauth(ctx, tx, oauth)
	if err != nil {
		//utils.PanicIfError(ctx, fiber.StatusBadRequest, err)
		return webrespones.UserDetail{}, exception.CustomEror{Code: fiber.StatusInternalServerError,
			Error: "Internal server error"}, err
	}
	fmt.Println(user)
	fmt.Println(oauth)
	fmt.Println("oaut")

	return template_response.ToUserRespone(user, oauth), exception.CustomEror{}, nil
}
