package photoTest

import (
	"bytes"
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/otptest"
	"github.com/Wrendra57/Pos-app-be/test/api/userstest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func UploadPhotoTestRequest(t *testing.T, body *bytes.Buffer, token string, writer *multipart.Writer) *http.Request {
	req, err := http.NewRequest("POST", "/api/v1/file/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestUploadPhotoSucccess(t *testing.T) {
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
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("owner_id", user.User_id.String())
	writer.WriteField("name", "testing upload")

	file, err := os.Open("../../../storage/photos/anonim-picture.png")
	utils.PanicIfError(err)
	defer file.Close()

	part, err := writer.CreateFormFile("foto", "anonim-picture.png")
	utils.PanicIfError(err)

	_, err = io.Copy(part, file)
	utils.PanicIfError(err)

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := UploadPhotoTestRequest(t, bodyReq, generateToken, writer)

	res, err := app.Test(request)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be success")
	assert.Equalf(t, "success", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestUploadPhotoValidationOwnerIdRequired(t *testing.T) {
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
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	//writer.WriteField("owner_id", user.User_id.String())
	writer.WriteField("name", "testing upload")

	file, err := os.Open("../../../storage/photos/anonim-picture.png")
	utils.PanicIfError(err)
	defer file.Close()

	part, err := writer.CreateFormFile("foto", "anonim-picture.png")
	utils.PanicIfError(err)

	_, err = io.Copy(part, file)
	utils.PanicIfError(err)

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := UploadPhotoTestRequest(t, bodyReq, generateToken, writer)

	res, err := app.Test(request)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be success")
	assert.Equalf(t, "owner_id is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestUploadPhotoValidationOwnerIdEmpty(t *testing.T) {
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
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("owner_id", "")
	writer.WriteField("name", "testing upload")

	file, err := os.Open("../../../storage/photos/anonim-picture.png")
	utils.PanicIfError(err)
	defer file.Close()

	part, err := writer.CreateFormFile("foto", "anonim-picture.png")
	utils.PanicIfError(err)

	_, err = io.Copy(part, file)
	utils.PanicIfError(err)

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := UploadPhotoTestRequest(t, bodyReq, generateToken, writer)

	res, err := app.Test(request)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be success")
	assert.Equalf(t, "owner_id is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestUploadPhotoValidationNameRequired(t *testing.T) {
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
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("owner_id", user.User_id.String())
	//writer.WriteField("name", "testing upload")

	file, err := os.Open("../../../storage/photos/anonim-picture.png")
	utils.PanicIfError(err)
	defer file.Close()

	part, err := writer.CreateFormFile("foto", "anonim-picture.png")
	utils.PanicIfError(err)

	_, err = io.Copy(part, file)
	utils.PanicIfError(err)

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := UploadPhotoTestRequest(t, bodyReq, generateToken, writer)

	res, err := app.Test(request)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be success")
	assert.Equalf(t, "Name is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)
}

func TestUploadPhotoValidationNameEmpty(t *testing.T) {
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
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("owner_id", user.User_id.String())
	writer.WriteField("name", "")

	file, err := os.Open("../../../storage/photos/anonim-picture.png")
	utils.PanicIfError(err)
	defer file.Close()

	part, err := writer.CreateFormFile("foto", "anonim-picture.png")
	utils.PanicIfError(err)

	_, err = io.Copy(part, file)
	utils.PanicIfError(err)

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := UploadPhotoTestRequest(t, bodyReq, generateToken, writer)

	res, err := app.Test(request)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be success")
	assert.Equalf(t, "Name is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestUploadPhotoValidationPhotoRequired(t *testing.T) {
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
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("owner_id", user.User_id.String())
	writer.WriteField("name", "testing")

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := UploadPhotoTestRequest(t, bodyReq, generateToken, writer)

	res, err := app.Test(request)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be success")
	assert.Equalf(t, "photo is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}
