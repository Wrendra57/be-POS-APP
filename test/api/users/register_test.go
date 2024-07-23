package users__test

import (
	"context"
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func registerTestRequest(t *testing.T, app *fiber.App, method, url, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func insertNewUserTest(t *testing.T, db *pgxpool.Pool, request webrequest.UserCreateRequest) (domain.User,
	domain.Oauth, domain.Roles, domain.OTP, domain.Photos) {
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
	otp, err = otpRepo.Insert(context.Background(), tx, otp)
	utils.PanicIfError(err)

	photo, _ := photoRepo.Insert(context.Background(), tx, domain.Photos{Url: "http://127.0.0.1:8080/foto/default-photo-picture.png", Owner: user.User_id})

	return user, oauth, role, otp, photo
}
func TestRegisterUserSuccess(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"birthday": req.Birthday,
		"address":  req.Address,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "ok", "response status should be ok")
	assert.Equalf(t, response.Message, "User created successfully", "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestRegisterUserEmailExist(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	_, _, _, _, _ = insertNewUserTest(t, db, req)

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"birthday": req.Birthday,
		"address":  req.Address,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Email already exists", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserFailedBirtdateRequire(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"address":  req.Address,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Birthdate is required", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserBirtdateFormatWrong(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"address":  req.Address,
		"birthday": "20-02-2000",
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Birthdate must be format YYYY-MM-DD", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserFailedValidationNameRequire(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"gender":   req.Gender,
		"telp":     req.Telp,
		"address":  req.Address,
		"birthday": req.Birthday,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Name is required", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserFailedValidationNameMinLength(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "te",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"address":  req.Address,
		"birthday": req.Birthday,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Name must be at least 3 characters long", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserFailedValidationNameMaxLength(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "Lorem ipsum odor amet, consectetuer adipiscing elit.",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"address":  req.Address,
		"birthday": req.Birthday,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Name must be maximum 32 characters long", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserFailedValidationGenderRequire(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "Lorem ipsum",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name": req.Name,
		//"gender":   req.Gender,
		"telp":     req.Telp,
		"address":  req.Address,
		"birthday": req.Birthday,
		"email":    req.Email,
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Gender is required", "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestRegisterUserFailedValidation(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "Lorem ipsum",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	tests := []struct {
		nameTest string
		//request string
		body            map[string]string
		expectedCode    int
		expectedStatus  string
		expectedMessage string
	}{
		{
			nameTest: "Failed validation wrong input of gender",
			body: map[string]string{
				"name":     req.Name,
				"gender":   "binary",
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Gender must be a male female",
		},
		{
			nameTest: "Failed validation required 'telp' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Telp is required",
		},
		{
			nameTest: "Failed validation min length 'telp' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     "09",
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Telp must be at least 3 characters long",
		},
		{
			nameTest: "Failed validation max length 'telp' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     "0922222222222222222",
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Telp must be maximum 15 characters long",
		}, {
			nameTest: "Failed validation required 'Address' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     "083244",
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address is required",
		}, {
			nameTest: "Failed validation min length 'Address' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  "ds",
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address must be at least 3 characters long",
		}, {
			nameTest: "Failed validation max length 'Address' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  "Hac consequat natoque sodales condimentum mus. Velit vitae lacinia integer finibus interdum laoreet condimentum semper. Primis eu nulla a, egestas elementum enim. Ante conubia class ante ornare quis elit sapien blandit ipsum. Vulputate viverssssssssssssssssss",
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address must be maximum 255 characters long",
		}, {
			nameTest: "Failed validation required 'Email' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Email is required",
		}, {
			nameTest: "Failed validation email format 'Email' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    "wahyu",
				"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Email must be a valid email address",
		},
		{
			nameTest: "Failed validation required 'Password' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				//"password": req.Password,
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password is required",
		}, {
			nameTest: "Failed validation min length 'Password' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": "wdwd",
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password must be at least 8 characters long",
		}, {
			nameTest: "Failed validation max length 'Password' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": "Hac consequat natoque sodaledawdasdaw",
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password must be maximum 32 characters long",
		},
		{
			nameTest: "Failed validation required 'Username' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": "",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Username is required",
		}, {
			nameTest: "Failed validation min length 'Username' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": "sa",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Username must be at least 3 characters long",
		}, {
			nameTest: "Failed validation max length 'Username' field",
			body: map[string]string{
				"name":     req.Name,
				"gender":   req.Gender,
				"telp":     req.Telp,
				"address":  req.Address,
				"birthday": req.Birthday,
				"email":    req.Email,
				"password": req.Password,
				"username": "Hac consequat natoque sodaledawdasdaw",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Username must be maximum 32 characters long",
		},
	}

	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			jsonData, err := json.Marshal(test.body)
			if err != nil {
				panic(err)
			}
			requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
			res, err := app.Test(requ, 3000)
			assert.Nil(t, err)

			body, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, test.expectedCode, res.StatusCode)

			var response webrespones.ResponseApi
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)
			}
			assert.Equalf(t, response.Status, test.expectedStatus, "response status should be equal")
			assert.Equalf(t, response.Message, test.expectedMessage, "response message should be equal")
			assert.Nil(t, response.Data)
		})
	}

}

func TestRegisterUserUsernameExist(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "testerrr",
	}
	_, _, _, _, _ = insertNewUserTest(t, db, req)

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	registerData := map[string]string{
		"name":     req.Name,
		"gender":   req.Gender,
		"telp":     req.Telp,
		"birthday": req.Birthday,
		"address":  req.Address,
		"email":    "testttt@gmail.com",
		"password": req.Password,
		"username": req.Username,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	requ := registerTestRequest(t, app, "POST", "/api/v1/users/register", string(jsonData))
	res, err := app.Test(requ, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, response.Status, "failed", "response status should be failed")
	assert.Equalf(t, response.Message, "Username already exists", "response message should be equal")
	assert.Nil(t, response.Data)
}
