package users

import (
	"encoding/json"
	"fmt"
	be "github.com/Wrendra57/Pos-app-be/cmd"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterUserSuccess(t *testing.T) {
	fmt.Println("TestRegisterUserSuccess")
	db, _, err := test.SetupDBTest()
	if err != nil {
		panic(err)
	}
	fmt.Println("otw truncate")
	err = test.TruncateDB(db)
	fmt.Println(" truncate")
	if err != nil {
		panic(err)
	}
	db.Close()
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	config.InitConfig()

	fmt.Println("init server")
	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	fmt.Println("init app")
	req := webrequest.UserCreateRequest{
		Name:     "testUser",
		Gender:   "male",
		Telp:     "08213243444",
		Birthday: "2023-07-15",
		Address:  "solo",
		Email:    "testUser@gmail.com",
		Password: "password",
		Username: "wawawawa",
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
	fmt.Println("input data")
	jsonData, err := json.Marshal(registerData)
	if err != nil {
		panic(err)
	}
	r, _ := http.NewRequest(
		"POST",
		"/api/v1/users/register",
		strings.NewReader(string(jsonData)),
	)
	res, err := app.Test(r, -1)

	assert.Equalf(t, 200, res.StatusCode, "register user should be successful")

	body, err := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equalf(t, 200, int(resBody["code"].(float64)), "register user should be result status code 200 ")
	assert.Equalf(t, "ok", resBody["status"], "register user should be return status ok")
	assert.Equalf(t, "User created successfully", resBody["message"], "register user should be return success message")
	assert.NotNilf(t, resBody["data"], "register user  should not be nil data")
}
