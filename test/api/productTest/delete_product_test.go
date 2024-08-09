package productTest

import (
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
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
	"strings"
	"testing"
)

type ResponseDeleteProductTest struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func deleteProductTest(t *testing.T, url, token string) *http.Request {
	req, err := http.NewRequest("DELETE", url, strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return req

}
func TestDeleteProductSuccess(t *testing.T) {
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
	brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	product := InsertProductTest(db, domain.Product{ProductName: "test product name", SellPrice: 5000, CallName: "test 1, test 2", AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
	_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
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

	url := "/api/v1/product/" + product.Id.String()
	request := deleteProductTest(t, url, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseDeleteProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success delete product", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	t.Log(response.Data)
}

func TestDeleteProductInvalidId(t *testing.T) {
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
	brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	product := InsertProductTest(db, domain.Product{ProductName: "test product name", SellPrice: 5000, CallName: "test 1, test 2", AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
	_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
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

	url := "/api/v1/product/" + product.Id.String()[2:]
	request := deleteProductTest(t, url, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseDeleteProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "invalid id product", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
}

func TestDeleteProductNotFoundProduct(t *testing.T) {
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
	brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	product := InsertProductTest(db, domain.Product{ProductName: "test product name", SellPrice: 5000, CallName: "test 1, test 2", AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
	_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})

	test.TruncateDB(db)
	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
	brand = brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier = suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category = categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	product2 := InsertProductTest(db, domain.Product{ProductName: "test product name", SellPrice: 5000, CallName: "test 1, test 2", AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
	_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product2.Id})
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

	url := "/api/v1/product/" + product.Id.String()
	request := deleteProductTest(t, url, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseDeleteProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Product not found", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
}

func TestDeleteProductWasDeletedProduct(t *testing.T) {
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
	brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	product := InsertProductTest(db, domain.Product{ProductName: "test product name", SellPrice: 5000, CallName: "test 1, test 2", AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
	_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})

	err = DeleteProductTest(db, product.Id)
	if err != nil {
		t.Fatal(err.Error())
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

	url := "/api/v1/product/" + product.Id.String()
	request := deleteProductTest(t, url, generateToken)
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseDeleteProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "product was deleted", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	t.Log(response.Data)
}
