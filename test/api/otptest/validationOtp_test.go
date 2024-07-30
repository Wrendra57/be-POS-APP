package otptest

import (
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be/cmd"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/userstest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func ValidationOtpTestRequest(t *testing.T, app *fiber.App, method, url, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestValidationOtpSuccess(t *testing.T) {
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
	user, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)

	otp := FindOtpRepo(db, user.User_id)
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": otp.Otp,
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	url := "/api/v1/users/otp/" + token

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)

	}
	assert.Equalf(t, "success", response.Status, "respone status should be equal")
	assert.Equalf(t, "success validate", response.Message, "response message should be equal")
}

func TestValidationOtpWrongOtp(t *testing.T) {
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
	_, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": "888234",
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	url := "/api/v1/users/otp/" + token

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, fiber.StatusNotFound, response.Code, "response code should be equal")
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Code Otp was wrong", response.Message, "response message should be equal")
	assert.Equalf(t, nil, response.Data, "response data should be equal")
}

func TestValidationOtpInvalidToken(t *testing.T) {
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
	user, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)

	otp := FindOtpRepo(db, user.User_id)
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": otp.Otp,
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	url := "/api/v1/users/otp/" + token + "a"

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)

	}
	assert.Equalf(t, fiber.StatusUnauthorized, response.Code, "response code should be equal")
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Unauthorized", response.Message, "response message should be equal")
	assert.Equalf(t, nil, response.Data, "response data should be equal")
}

func TestValidationOtpWrongsToken(t *testing.T) {
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

	_, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)
	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	_, _, _, _, _, _ = userstest.InsertNewUserTest(t, db, req)

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": "231342",
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	//token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIzNTMwNjU4NjksImxldmVsIjoibWVtYmVyIiwidXNlcl9pZCI6IjIxOGViYjg3LTY3MTEtNDFjMy1hMWU1LTJiNzUzOTQ1YjE0NiJ9.BuywcVBgPFwiiUZNfkbN3lrJEIaefVDNEY_BKjRuqhI"
	url := "/api/v1/users/otp/" + token

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)

	}
	assert.Equalf(t, fiber.StatusNotFound, response.Code, "response code should be equal")
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "user not found", response.Message, "response message should be equal")
	assert.Equalf(t, nil, response.Data, "response data should be equal")
}

func TestValidationOtpEnableAccount(t *testing.T) {
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

	user, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)

	_ = UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": "231342",
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	url := "/api/v1/users/otp/" + token

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)

	}
	assert.Equalf(t, fiber.StatusBadRequest, response.Code, "response code should be equal")
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Account is already enabled", response.Message, "response message should be equal")
	assert.Equalf(t, nil, response.Data, "response data should be equal")
}

func TestValidationOtpDeletedOtp(t *testing.T) {
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

	_, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)
	_ = TruncateOtp(db)
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": "231342",
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	url := "/api/v1/users/otp/" + token

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)

	}
	assert.Equalf(t, fiber.StatusNotFound, response.Code, "response code should be equal")
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Code Otp Was Not Found", response.Message, "response message should be equal")
	assert.Equalf(t, nil, response.Data, "response data should be equal")
}

func TestValidationOtpExpiredOtp(t *testing.T) {
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

	user, _, _, _, _, token := userstest.InsertNewUserTest(t, db, req)
	_ = UpdateOtpExpired(db, user.User_id)

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReq := map[string]string{
		"otp": "231342",
	}
	jsonBody, err := json.Marshal(bodyReq)
	if err != nil {
		panic(err)
	}
	url := "/api/v1/users/otp/" + token

	request := ValidationOtpTestRequest(t, app, "POST", url, string(jsonBody))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	bodyResp, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

	var response webrespones.ResponseApi

	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)

	}
	assert.Equalf(t, fiber.StatusNotFound, response.Code, "response code should be equal")
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Code Otp Was Expired", response.Message, "response message should be equal")
	assert.Equalf(t, nil, response.Data, "response data should be equal")
}
