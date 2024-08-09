package productTest

import (
	"bytes"
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/brandTest"
	categoryTests "github.com/Wrendra57/Pos-app-be/test/api/categoryTest"
	"github.com/Wrendra57/Pos-app-be/test/api/otptest"
	suplier "github.com/Wrendra57/Pos-app-be/test/api/suplier_test"
	"github.com/Wrendra57/Pos-app-be/test/api/userstest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func CreateProductTestReq(t *testing.T, body *bytes.Buffer, token string, writer *multipart.Writer) *http.Request {
	req, err := http.NewRequest("POST", "/api/v1/product", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestCreateProductSuccess(t *testing.T) {
	test.InitConfigTest()
	defer os.RemoveAll("./storage")
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "Success", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}
func TestCreateProductRequiredSellPrice(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "sell price is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}
func TestCreateProductEmptySellPrice(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "sell price is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductNotIntegerSellPrice(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "120o3")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "sell price must be integer/number", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductRequiredCategory(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "1000")
	writer.WriteField("call_name", "test, ting, name")

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "category is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductEmptyCategory(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "1000")
	writer.WriteField("call_name", "test, ting, name")

	writer.WriteField("category", "")

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "category is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}
func TestCreateProductFailedParsingCategory(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "1000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String()[1:],
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "invalid parse category", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductRequiredBrand(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "brand is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductEmptyBrand(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	writer.WriteField("brand", "")

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "brand is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductFailedParseBrand(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   strconv.Itoa(brand.Id),
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "invalid parse brand", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductRequiredSupplier(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	//supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	//jsonSupplierReq, err := json.Marshal(supplierReq)
	//writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "supplier is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductEmptySupplier(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	//supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	//jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", "")

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "supplier is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductInvalidParseSupplier(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]interface{}{"id": supplier.Id[1:], "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "invalid parse supplier", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductMinimalPhoto(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]interface{}{"id": supplier.Id, "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "Photo minimal 1 or maximal 15", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}
func TestCreateProductMaximalPhoto(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]interface{}{"id": supplier.Id, "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png",
		"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png", "../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "Photo minimal 1 or maximal 15", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductRequiredProductName(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	//writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "ProductName is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductMinLengthProductName(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "t")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "ProductName must be at least 2 characters long", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductMinValueSellPrice(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing")
	writer.WriteField("sell_price", "-1")
	writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "SellPrice must be greater than 0", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductRequiredCallName(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "testing name")
	writer.WriteField("sell_price", "5000")
	//writer.WriteField("call_name", "test, ting, name")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "CallName is required", response.Message, "response message should be equal")
	assert.Nil(t, response.Data)

}

func TestCreateProductMinLengthCallName(t *testing.T) {
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

	brand := domain.Brand{Name: "danone", Description: "test_test"}
	brand = brandTest.InsertBrandTest(db, brand)

	category := domain.Category{Name: "test_primary", Description: "test_test"}
	category = categoryTests.InsertCategoriesTest(db, category)

	supplier := domain.Supplier{Name: "test_agus", ContactInfo: "test@test.com", Address: "solo_test"}
	supplier = suplier.InsertSupplierTest(db, supplier)

	db.Close()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	bodyReq := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyReq)
	writer.WriteField("product_name", "test")
	writer.WriteField("sell_price", "5000")
	writer.WriteField("call_name", "t")
	categoryReq := map[string]string{
		"id":   category.Id.String(),
		"name": category.Name,
	}
	jsonCategoryReq, err := json.Marshal(categoryReq)
	utils.PanicIfError(err)
	writer.WriteField("category", string(jsonCategoryReq))

	brandReq := map[string]interface{}{
		"id":   brand.Id,
		"name": brand.Name,
	}
	jsonBrandReq, err := json.Marshal(brandReq)
	writer.WriteField("brand", string(jsonBrandReq))

	supplierReq := map[string]string{"id": supplier.Id.String(), "name": supplier.Name}
	jsonSupplierReq, err := json.Marshal(supplierReq)
	writer.WriteField("supplier", string(jsonSupplierReq))

	files := []string{"../../../storage/photos/anonim-picture.png", "../../../storage/photos/anonim-picture-2.png"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		utils.PanicIfError(err)
		defer file.Close()

		part, err := writer.CreateFormFile("photo", filePath)
		utils.PanicIfError(err)

		_, err = io.Copy(part, file)
		utils.PanicIfError(err)
	}
	writer.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	request := CreateProductTestReq(t, bodyReq, generateToken, writer)

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
	assert.Equalf(t, "CallName must be at least 2 characters long", response.Message,
		"response message should be equal")
	assert.Nil(t, response.Data)

}
