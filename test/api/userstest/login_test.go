package userstest

import (
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/otptest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func LoginUserTestRequest(t *testing.T, method, url, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}
func TestLoginUsingEmailSuccess(t *testing.T) {
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

	bodyReqData := map[string]string{
		"username": req.Email,
		"password": req.Password,
	}
	jsonRequest, err := json.Marshal(bodyReqData)
	if err != nil {
		panic(err)

	}

	request := LoginUserTestRequest(t, "POST", "/api/v1/users/login", string(jsonRequest))
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
	assert.Equalf(t, "Success login", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)

}

func TestLoginUsingUsernameSuccess(t *testing.T) {
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

	bodyReqData := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	jsonRequest, err := json.Marshal(bodyReqData)
	if err != nil {
		panic(err)

	}

	request := LoginUserTestRequest(t, "POST", "/api/v1/users/login", string(jsonRequest))
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
	assert.Equalf(t, "Success login", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)

}

func TestLoginFailedValidation(t *testing.T) {
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
		nameTest        string
		body            map[string]string
		expectedCode    int
		expectedStatus  string
		expectedMessage string
	}{
		{
			nameTest: "Failed validation required 'username' field not exist",
			body: map[string]string{
				"password": req.Password,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "UserName is required",
		}, {
			nameTest: "Failed validation required 'username' field empty string",
			body: map[string]string{
				"username": "",
				"password": req.Password,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "UserName is required",
		}, {
			nameTest: "Failed validation min Length 'username' field  ",
			body: map[string]string{
				"username": "ws",
				"password": req.Password,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "UserName must be at least 3 characters long",
		}, {
			nameTest: "Failed validation max Length 'username' field  ",
			body: map[string]string{
				"username": "wswswswswswswswswswswswswswswswswswswswswswswswswswswswswswsqsqsqs",
				"password": req.Password,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "UserName must be maximum 32 characters long",
		}, {
			nameTest: "Failed validation required 'password' field not exist",
			body: map[string]string{
				"username": req.Username,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password is required",
		}, {
			nameTest: "Failed validation required 'password' field empty string",
			body: map[string]string{
				"username": req.Username,
				"password": "",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password is required",
		}, {
			nameTest: "Failed validation min length 'password' field",
			body: map[string]string{
				"username": req.Username,
				"password": "wswsw",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password must be at least 8 characters long",
		}, {
			nameTest: "Failed validation max length 'password' field",
			body: map[string]string{
				"username": req.Username,
				"password": "wswswswswswswswswswswswswswswswswswswswswswswswswswswswswswsqsqsqs",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Password must be maximum 32 characters long",
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			jsonDataRequest, err := json.Marshal(test.body)
			if err != nil {
				panic(err)
			}

			request := LoginUserTestRequest(t, "POST", "/api/v1/users/login", string(jsonDataRequest))
			res, err := app.Test(request, 3000)
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

func TestLoginUsernameNotFound(t *testing.T) {
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
		Username: "wawawaew",
	}

	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReqData := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	jsonRequest, err := json.Marshal(bodyReqData)
	if err != nil {
		panic(err)

	}

	request := LoginUserTestRequest(t, "POST", "/api/v1/users/login", string(jsonRequest))
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
	assert.Equalf(t, "Account / Password was wrong", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestLoginAccountNotEnabled(t *testing.T) {
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
	_, _, _, _, _, _ = InsertNewUserTest(t, db, req)

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	bodyReqData := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	jsonRequest, err := json.Marshal(bodyReqData)
	if err != nil {
		panic(err)

	}

	request := LoginUserTestRequest(t, "POST", "/api/v1/users/login", string(jsonRequest))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Account not enabled", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)

}

func TestLoginWrongPassword(t *testing.T) {
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

	bodyReqData := map[string]string{
		"username": req.Username,
		"password": req.Password + req.Password,
	}
	jsonRequest, err := json.Marshal(bodyReqData)
	if err != nil {
		panic(err)

	}

	request := LoginUserTestRequest(t, "POST", "/api/v1/users/login", string(jsonRequest))
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Account / Password was wrong", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)

}
