package userstest

import (
	"context"
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be/cmd"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/otptest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func AuthMeTestRequest(t *testing.T, token string) *http.Request {
	req, err := http.NewRequest("GET", "/api/v1/user", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestAuthMeSuccess(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, role, _, _, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := AuthMeTestRequest(t, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success Get Data", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestAuthMeSuccessViaRedis(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, role, _, photos, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	resp := domain.UserDetail{
		User_id:      user.User_id,
		Email:        req.Email,
		Username:     req.Username,
		Name:         req.Name,
		Gender:       req.Gender,
		Telp:         req.Telp,
		Birthday:     user.Birthday,
		Address:      req.Address,
		Foto_profile: photos.Url,
		Role:         role.Role,
		Created_at:   user.Created_at,
	}
	jsonResp, err := json.Marshal(resp)
	rdb := test.SetupRedisTest()
	err = rdb.Set(context.Background(), user.User_id.String(), jsonResp, 60*time.Second).Err()
	if err != nil {
		panic(err)
	}
	rdb.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := AuthMeTestRequest(t, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success Get Data", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestAuthMeWithoutAuthorization(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, _, _, _, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request, err := http.NewRequest("GET", "/api/v1/user", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Unauthorized", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

func TestAuthMeWithoutBearer(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, _, _, _, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request, err := http.NewRequest("GET", "/api/v1/user", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Unauthorized", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

func TestAuthMeWrongFormatToken(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, role, _, _, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	token := generateToken + "dwfewf"
	request := AuthMeTestRequest(t, token)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Unauthorized", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}
func TestAuthMeInvalidToken(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, _, _, _, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjA1MzY3OTQsImxldmVsMiI6Im1lbWJlciIsInVzZXJfaWQiOiI4NDdmMzRkNi0yN2YwLTRhNzMtOTEzNC1iNjdmM2ViZTNmYzQifQ.tMo9jkkvUuUDWg5xeeeePf4L20T1Nw_p5whdL1Cj0Mo"

	request := AuthMeTestRequest(t, token)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Unauthorized", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

func TestAuthMeUserNotExist(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
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
	user, _, role, _, _, _ := InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := AuthMeTestRequest(t, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "User not found", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}
