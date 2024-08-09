package productTest

import (
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/brandTest"
	categoryTests "github.com/Wrendra57/Pos-app-be/test/api/categoryTest"
	"github.com/Wrendra57/Pos-app-be/test/api/photoTest"
	suplier "github.com/Wrendra57/Pos-app-be/test/api/suplier_test"
	"github.com/Wrendra57/Pos-app-be/test/api/userstest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type ResponseListTest struct {
	Code    int                  `json:"code"`
	Status  string               `json:"status"`
	Data    []domain.ProductList `json:"data"`
	Message string               `json:"message"`
}

func GetListProduct(t *testing.T, url, body string) *http.Request {
	req, err := http.NewRequest("GET", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req

}

func TestGetListProductSuccess(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/product?" + limit + "&" + offset
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, total, "should lenght "+strconv.Itoa(total))
}

func TestGetListProductWithParamsSuccess(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=test product name 1"
	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/product?" + limit + "&" + offset + "&" + param
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, 1, "should lenght "+strconv.Itoa(1))
}

func TestGetListProductWithParamsSuccessNoResult(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=test product name 1234"
	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/product?" + limit + "&" + offset + "&" + param
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, 0, "should lenght "+strconv.Itoa(1))
}

func TestGetListProductWithoutLimit(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	offset := "offset=1"
	url := "/api/v1/product?" + offset
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, total, "should lenght "+strconv.Itoa(total))
}

func TestGetListProductWithoutOffset(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	limit := "limit=10"

	url := "/api/v1/product?" + limit
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, total, "should lenght "+strconv.Itoa(total))
}

func TestGetListProductFailedLimit(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=test product name 1234"
	limit := "limit=w2"
	offset := "offset=1"
	url := "/api/v1/product?" + limit + "&" + offset + "&" + param
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "The 'limit' field must be number/integer", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, 0, "should lenght "+strconv.Itoa(1))
}

func TestGetListProductFailedNegativeLimit(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=test product name 1234"
	limit := "limit=-20"
	offset := "offset=1"
	url := "/api/v1/product?" + limit + "&" + offset + "&" + param
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "The 'limit' field must be greater than zero", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, 0, "should lenght "+strconv.Itoa(1))
}

func TestGetListProductFailedOffset(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=test product name 1234"
	limit := "limit=10"
	offset := "offset=1ww"
	url := "/api/v1/product?" + limit + "&" + offset + "&" + param
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "The 'offset' field must be number/integer", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, 0, "should lenght "+strconv.Itoa(1))
}

func TestGetListProductFailedNegativeOffset(t *testing.T) {
	test.InitConfigTest()
	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}

	total := 5
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
	for i := 0; i < total; i++ {
		brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand " + strconv.Itoa(i), Description: "test brand" + strconv.Itoa(i)})
		supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier " + strconv.Itoa(i), ContactInfo: strconv.Itoa(i) + "testsupplier@gmail.com", Address: "test , south test" + strconv.Itoa(i)})
		category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category " + strconv.Itoa(i), Description: "test category " + strconv.Itoa(i)})
		product := InsertProductTest(db, domain.Product{ProductName: "test product name " + strconv.Itoa(i), SellPrice: 5000, CallName: "test 1, test 2 " + strconv.Itoa(i), AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
		_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=test product name 1234"
	limit := "limit=20"
	offset := "offset=-1"
	url := "/api/v1/product?" + limit + "&" + offset + "&" + param
	request := GetListProduct(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "The 'offset' field must be positive", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
	assert.Lenf(t, response.Data, 0, "should lenght "+strconv.Itoa(1))
}
