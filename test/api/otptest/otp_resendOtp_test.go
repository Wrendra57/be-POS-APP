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
	"time"
)

func ResendOtpTestRequest(t *testing.T, app *fiber.App, method, url, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}
func TestResendOtpSuccess(t *testing.T) {
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

	url := "/api/v1/users/otp/resend/" + token
	request := ResendOtpTestRequest(t, app, "POST", url, "")
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
	assert.Equalf(t, "success send otp again", response.Message, "response message should be equal")
}

func TestResendOtpInvalidToken(t *testing.T) {
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

	url := "/api/v1/users/otp/resend/a" + token
	request := ResendOtpTestRequest(t, app, "POST", url, "")
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
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Unauthorized", response.Message, "response message should be equal")
	assert.Empty(t, response.Data, "response data should be empty")
}

func TestResendOtpWrongToken(t *testing.T) {
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

	url := "/api/v1/users/otp/resend/" + token
	request := ResendOtpTestRequest(t, app, "POST", url, "")
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
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Account not found", response.Message, "response message should be equal")
	assert.Empty(t, response.Data, "response data should be empty")

}

func TestResendOtpEnabledAccount(t *testing.T) {
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

	url := "/api/v1/users/otp/resend/" + token
	request := ResendOtpTestRequest(t, app, "POST", url, "")
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
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Equalf(t, "Account is already enabled", response.Message, "response message should be equal")
	assert.Empty(t, response.Data, "response data should be empty")
}

func TestResendOtpMaxAccess(t *testing.T) {
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
	for i := 0; i < 6; i++ {
		_ = InsertOtpTest(db, domain.OTP{Otp: "131313", User_id: user.User_id, Expired_date: time.Now().Add(time.Minute)})
	}
	db.Close()
	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	url := "/api/v1/users/otp/resend/" + token
	request := ResendOtpTestRequest(t, app, "POST", url, "")
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
	assert.Equalf(t, "failed", response.Status, "respone status should be equal")
	assert.Empty(t, response.Data, "response data should be empty")
}
