package categoryTests

import (
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
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func CreateCategoriesTestRequest(t *testing.T, body, token string) *http.Request {
	req, err := http.NewRequest("POST", "/api/v1/categories", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestCategoriesCreateSuccess(t *testing.T) {
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

	category := domain.Category{
		Name:        "testCategory",
		Description: "testCategory description",
	}
	bodyReq := map[string]string{
		"name":        category.Name,
		"description": category.Description,
	}
	jsonReq, _ := json.Marshal(bodyReq)

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateCategoriesTestRequest(t, string(jsonReq), generateToken)
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
	assert.Equalf(t, "Success create category", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)

}

func TestCategoriesCreateValidateFailed(t *testing.T) {
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

	category := domain.Category{
		Name:        "testCategory",
		Description: "testCategory description",
	}

	tests := []struct {
		nameTest        string
		body            map[string]string
		expectedCode    int
		expectedStatus  string
		expectedMessage string
	}{
		{
			nameTest: "Failed validation required 'name' field not exist",
			body: map[string]string{
				"description": category.Description,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name is required",
		}, {
			nameTest: "Failed validation required 'name' field empty",
			body: map[string]string{
				"name":        "",
				"description": category.Description,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name is required",
		}, {
			nameTest: "Failed validation min length 'name' field",
			body: map[string]string{
				"name":        "as",
				"description": category.Description,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name must be at least 3 characters long",
		}, {
			nameTest: "Failed validation required 'description' field not exist",
			body: map[string]string{
				"name": category.Name,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Description is required",
		}, {
			nameTest: "Failed validation required 'description' field empty",
			body: map[string]string{
				"name":        category.Name,
				"description": "",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Description is required",
		},
		{
			nameTest: "Failed validation Min Length 'description' field",
			body: map[string]string{
				"name":        category.Name,
				"description": "sd",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Description must be at least 3 characters long",
		},
	}

	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			jsonReq, _ := json.Marshal(test.body)

			app, clean, err := be.InitializeApp()
			if err != nil {
				panic(err)
			}
			defer clean()

			request := CreateCategoriesTestRequest(t, string(jsonReq), generateToken)
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
			assert.Equalf(t, test.expectedStatus, response.Status, "response status should be ok")
			assert.Equalf(t, test.expectedMessage, response.Message, "response message should be equal")
			assert.Empty(t, response.Data)
		})
	}

}

func TestBrandCreateWithoutToken(t *testing.T) {
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
	user, _, _, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	_ = otptest.UpdateOauthTest(db, domain.Oauth{User_id: user.User_id, Is_enabled: true})
	db.Close()

	category := domain.Category{
		Name:        "testCategory",
		Description: "testCategory description",
	}
	bodyReq := map[string]string{
		"name":        category.Name,
		"description": category.Description,
	}
	jsonReq, _ := json.Marshal(bodyReq)

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request, err := http.NewRequest("POST", "/api/v1/categories", strings.NewReader(string(jsonReq)))
	request.Header.Set("Content-Type", "application/json")

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
