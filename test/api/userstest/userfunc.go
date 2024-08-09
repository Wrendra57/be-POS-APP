package userstest

import (
	"context"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
	"time"
)

func InsertNewUserTest(t *testing.T, db *pgxpool.Pool, request webrequest.UserCreateRequest) (domain.User,
	domain.Oauth, domain.Roles, domain.OTP, domain.Photos, string) {
	userRepo := repositories.NewUserRepository()
	oauthRepo := repositories.NewOauthRepository()
	roleRepo := repositories.NewRoleRepository()
	otpRepo := repositories.NewOtpRepository()
	photoRepo := repositories.NewPhotosRepository()
	tx, err := db.BeginTx(context.Background(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(context.Background(), tx)

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		panic(err)
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

	user, err = userRepo.InsertUser(context.Background(), tx, user)
	utils.PanicIfError(err)

	oauth.User_id = user.User_id
	//insert to db oauths
	oauth.User_id = user.User_id
	oauth, err = oauthRepo.InsertOauth(context.Background(), tx, oauth)
	utils.PanicIfError(err)

	//insert to db roles
	role, err := roleRepo.Insert(context.Background(), tx, domain.Roles{Role: "member", User_id: user.User_id})
	utils.PanicIfError(err)

	//create OTP using random 6 angka
	otp := domain.OTP{Otp: utils.GenerateOTP(), User_id: user.User_id, Expired_date: time.Now().Add(time.Minute * 3),
		Created_at: time.Now(), Updated_at: time.Now()}

	//Insert OTP to db
	otp = otpRepo.Insert(context.Background(), tx, otp)

	photo, _ := photoRepo.Insert(context.Background(), tx, domain.Photos{Url: "http://127.0.0.1:8080/foto/default-photo-picture.png", Owner: user.User_id})

	JWTStr, err := utils.GenerateJWT(user.User_id, role.Role)
	utils.PanicIfError(err)

	return user, oauth, role, otp, photo, JWTStr
}
