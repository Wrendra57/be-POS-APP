package suplier_test

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

func CreateSupplierTestRequest(t *testing.T, body, token string) *http.Request {
	req, err := http.NewRequest("POST", "/api/v1/supplier", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestSupplierCreateSuccess(t *testing.T) {
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

	supplier := domain.Supplier{
		Name:        "testSupplier",
		ContactInfo: "0821324532",
		Address:     "test address",
	}
	bodyReq := map[string]string{
		"name":         supplier.Name,
		"contact_info": supplier.ContactInfo,
		"address":      supplier.Address,
	}
	jsonReq, _ := json.Marshal(bodyReq)

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateSupplierTestRequest(t, string(jsonReq), generateToken)
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
	assert.Equalf(t, "Success create supplier", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}
func TestSupplierCreateValidationFailed(t *testing.T) {
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
	supplier := domain.Supplier{
		Name:        "testSupplier",
		ContactInfo: "0821324532",
		Address:     "test address",
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
				"contact_info": supplier.ContactInfo,
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name is required",
		}, {
			nameTest: "Failed validation required 'name' field empty",
			body: map[string]string{
				"name":         "",
				"contact_info": supplier.ContactInfo,
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name is required",
		}, {
			nameTest: "Failed validation Min Length 'name' field",
			body: map[string]string{
				"name":         "sd",
				"contact_info": supplier.ContactInfo,
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name must be at least 3 characters long",
		}, {
			nameTest: "Failed validation Max Length 'name' field",
			body: map[string]string{
				"name":         "sdwdwdwdwdsdwdwdwdwdsdwdwdwdwddwdw",
				"contact_info": supplier.ContactInfo,
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Name must be maximum 32 characters long",
		}, {
			nameTest: "Failed validation required 'contact_info' field not exist",
			body: map[string]string{
				"name":    supplier.Name,
				"address": supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "ContactInfo is required",
		}, {
			nameTest: "Failed validation required 'contact_info' field empty",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": "",
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "ContactInfo is required",
		}, {
			nameTest: "Failed validation Min Length 'contact_info' field",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": "dw",
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "ContactInfo must be at least 3 characters long",
		}, {
			nameTest: "Failed validation Max Length 'contact_info' field",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi volutpat neque tortor, quis tristique dolor fermentum sit amet. Vivamus augue justo, ultricies sit amet elit nec, fermentum mollis metus. Sed imperdiet orci eget dui consectetur lobortis. Duis tempor urna eget porta ultrices. Vivamus semper accumsan commodo. Nunc quis nibh eu mi rutrum alique",
				"address":      supplier.Address,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "ContactInfo must be maximum 32 characters long",
		}, {
			nameTest: "Failed validation required 'address' field not exist",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": supplier.ContactInfo,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address is required",
		}, {
			nameTest: "Failed validation required 'address' field empty",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": supplier.ContactInfo,
				"address":      "",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address is required",
		}, {
			nameTest: "Failed validation Min Length 'Address' field",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": supplier.ContactInfo,
				"address":      "dw",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address must be at least 3 characters long",
		}, {
			nameTest: "Failed validation Max Length 'Address' field",
			body: map[string]string{
				"name":         supplier.Name,
				"contact_info": supplier.ContactInfo,
				"address":      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi volutpat neque tortor, quis tristique dolor fermentum sit amet. Vivamus augue justo, ultricies sit amet elit nec, fermentum mollis metus. Sed imperdiet orci eget dui consectetur lobortis. Duis tempor urna eget porta ultrices. Vivamus semper accumsan commodo. Nunc quis nibh eu mi rutrum alique",
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Address must be maximum 232 characters long",
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

			request := CreateSupplierTestRequest(t, string(jsonReq), generateToken)
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

func TestSupplierCreateWithoutToken(t *testing.T) {
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

	supplier := domain.Supplier{
		Name:        "testSupplier",
		ContactInfo: "0821324532",
		Address:     "test address",
	}
	bodyReq := map[string]string{
		"name":         supplier.Name,
		"contact_info": supplier.ContactInfo,
		"address":      supplier.Address,
	}
	jsonReq, _ := json.Marshal(bodyReq)

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request, err := http.NewRequest("POST", "/api/v1/supplier", strings.NewReader(string(jsonReq)))
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

func TestSupplierCreateWrongBodyReq(t *testing.T) {
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

	_ = domain.Supplier{
		Name:        "testSupplier",
		ContactInfo: "0821324532",
		Address:     "test address",
	}
	bodyReq := map[string]int{
		"name":         12,
		"contact_info": 12,
		"address":      12,
	}
	jsonReq, _ := json.Marshal(bodyReq)

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateSupplierTestRequest(t, string(jsonReq), generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

	var response webrespones.ResponseApi
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Internal server error", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}
