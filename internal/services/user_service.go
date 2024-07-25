package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserService interface {
	CreateUser(ctx *fiber.Ctx, request webrequest.UserCreateRequest) (string, exception.CustomEror,
		error)
	Login(ctx *fiber.Ctx, request webrequest.UserLoginRequest) (webrespones.TokenResp, exception.CustomEror, bool)
	AuthMe(ctx *fiber.Ctx) (domain.UserDetail, exception.CustomEror, bool)
}

type userServiceImpl struct {
	UserRepository  repositories.UserRepository
	OauthRepository repositories.OauthRepository
	OtpRepository   repositories.OtpRepository
	RoleRepository  repositories.RoleRepository
	PhotoRepository repositories.PhotosRepository
	DB              *pgxpool.Pool
	Validate        *validator.Validate
	RedisDB         *redis.Client
}

func NewUserService(db *pgxpool.Pool,
	validate *validator.Validate, rdb *redis.Client, userRepo repositories.UserRepository,
	oauthRepo repositories.OauthRepository, otpRepo repositories.OtpRepository, roleRepo repositories.RoleRepository,
	photoRepo repositories.PhotosRepository) UserService {
	return &userServiceImpl{
		UserRepository:  userRepo,
		OauthRepository: oauthRepo,
		OtpRepository:   otpRepo,
		PhotoRepository: photoRepo,
		RoleRepository:  roleRepo,
		DB:              db,
		RedisDB:         rdb,
		Validate:        validate,
	}
}

func (s userServiceImpl) CreateUser(ctx *fiber.Ctx, request webrequest.UserCreateRequest) (string,
	exception.CustomEror, error) {

	photoTemplate := "http://127.0.0.1:8080/foto/default-photo-picture.png"
	// start database tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	//check email in db
	_, err = s.OauthRepository.FindByEmail(ctx, tx, request.Email)
	if err == nil {
		fmt.Println(err)
		return "", exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "Email already exists"}, errors.New("Email already exist")
	}

	//check username in db
	_, err = s.OauthRepository.FindByUserName(ctx, tx, request.Username)
	if err == nil {
		return "", exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "Username already exists"}, errors.New("Username already exist")
	}

	//hashing password using bycript
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

	//insert to db userstest
	user, err = s.UserRepository.InsertUser(ctx.Context(), tx, user)
	utils.PanicIfError(err)

	//insert to db oauths
	oauth.User_id = user.User_id
	oauth, err = s.OauthRepository.InsertOauth(ctx.Context(), tx, oauth)
	utils.PanicIfError(err)

	//insert to db roles
	role, err := s.RoleRepository.Insert(ctx.Context(), tx, domain.Roles{Role: "member", User_id: user.User_id})
	utils.PanicIfError(err)

	//create OTP using random 6 angka
	otp := domain.OTP{Otp: utils.GenerateOTP(), User_id: user.User_id, Expired_date: time.Now().Add(time.Minute * 3),
		Created_at: time.Now(), Updated_at: time.Now()}

	//Insert OTP to db
	otp, err = s.OtpRepository.Insert(ctx.Context(), tx, otp)
	utils.PanicIfError(err)

	//insert photo default to db
	_, err = s.PhotoRepository.Insert(ctx.Context(), tx, domain.Photos{Url: photoTemplate, Owner: user.User_id})
	utils.PanicIfError(err)

	//sending otp via email
	//strOTP := "ini kode token kamu " + otp.Otp
	//err = utils.SendEmail("wrendra57@gmail.com", "OTP-ACCOUNT", strOTP)
	//utils.PanicIfError(err)

	//GenerateJWT for access validasi otp
	JWTStr, err := utils.GenerateJWT(user.User_id, role.Role)
	utils.PanicIfError(err)

	return JWTStr, exception.CustomEror{}, nil
}

func (s userServiceImpl) Login(ctx *fiber.Ctx, request webrequest.UserLoginRequest) (webrespones.TokenResp, exception.CustomEror,
	bool) {
	// begin database tx

	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	//get data from db
	o, err := s.OauthRepository.FindByUsernameOrEmail(ctx, tx, request.UserName)
	if err != nil {
		return webrespones.TokenResp{}, exception.CustomEror{Code: fiber.StatusNotFound,
			Error: "Account / Password was wrong"}, false
	}
	if !o.Is_enabled {
		return webrespones.TokenResp{}, exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "Account not enabled"}, false
	}
	//compare password
	comparePassword := utils.CheckPasswordHash(request.Password, o.Password)

	if !comparePassword {
		return webrespones.TokenResp{}, exception.CustomEror{Code: fiber.StatusBadRequest,
			Error: "Account / Password was wrong"}, false
	}

	//get role by user id
	role, err := s.RoleRepository.FindByUserId(ctx, tx, o.User_id)
	utils.PanicIfError(err)

	//generate token jwt from user_id, role
	tokenJwt, err := utils.GenerateJWT(o.User_id, role.Role)
	utils.PanicIfError(err)

	return webrespones.TokenResp{Token: tokenJwt}, exception.CustomEror{}, true

}

func (s userServiceImpl) AuthMe(ctx *fiber.Ctx) (domain.UserDetail, exception.CustomEror, bool) {
	//TODO implement me
	userId, _ := ctx.Locals("user_id").(uuid.UUID)
	var user domain.UserDetail

	//check data in redis
	result := s.RedisDB.Get(ctx.Context(), userId.String())
	if len(result.Val()) != 0 {
		err := json.Unmarshal([]byte(result.Val()), &user)
		utils.PanicIfError(err)
		fmt.Println(result.Val())
		return user, exception.CustomEror{}, true
	}

	//begin db tx
	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	//get data from db
	user, err = s.UserRepository.FindUserDetail(ctx, tx, userId)
	if err != nil {
		return user, exception.CustomEror{Code: fiber.StatusNotFound, Error: "User not found"}, false
	}

	//convert to json
	jsonData, err := json.Marshal(user)
	utils.PanicIfError(err)

	//insert to redis
	err = s.RedisDB.Set(ctx.Context(), userId.String(), jsonData, 60*time.Second).Err()
	utils.PanicIfError(err)

	return user, exception.CustomEror{}, true
}
