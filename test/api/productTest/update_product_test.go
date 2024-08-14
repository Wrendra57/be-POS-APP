package productTest

import (
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/brandTest"
	categoryTests "github.com/Wrendra57/Pos-app-be/test/api/categoryTest"
	"github.com/Wrendra57/Pos-app-be/test/api/photoTest"
	suplier "github.com/Wrendra57/Pos-app-be/test/api/suplier_test"
	"github.com/Wrendra57/Pos-app-be/test/api/userstest"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

type ResponseUpdateProductTest struct {
	Code    int                                    `json:"code"`
	Status  string                                 `json:"status"`
	Data    webrespones.ProductFindByIdResponseApi `json:"data"`
	Message string                                 `json:"message"`
}

func updateProductTestReq(t *testing.T, url, body, token string) *http.Request {
	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req

}

func TestUpdateProductSuccess(t *testing.T) {
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

	brand2 := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier2 := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category2 := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	db.Close()

	reqProduct := map[string]any{
		"product_name": "updated2",
		"sell_price":   3000,
		"call_name":    "update update",
		"category":     category2.Id,
		"brand":        brand2.Id,
		"supplier":     supplier2.Id,
	}
	jsonReqUpdate, err := json.Marshal(reqProduct)
	utils.PanicIfError(err)

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}
	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	url := "http://127.0.0.1:8080/api/v1/product/" + product.Id.String()
	request := updateProductTestReq(t, url, string(jsonReqUpdate), generateToken)
	res, err := app.Test(request, 3000)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response ResponseUpdateProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success update product", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestUpdateProductInvalidProductId(t *testing.T) {
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

	brand2 := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier2 := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category2 := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	db.Close()

	reqProduct := map[string]any{
		"product_name": "updated2",
		"sell_price":   3000,
		"call_name":    "update update",
		"category":     category2.Id,
		"brand":        brand2.Id,
		"supplier":     supplier2.Id,
	}
	jsonReqUpdate, err := json.Marshal(reqProduct)
	utils.PanicIfError(err)

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}
	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	url := "http://127.0.0.1:8080/api/v1/product/" + product.Id.String()[2:]
	request := updateProductTestReq(t, url, string(jsonReqUpdate), generateToken)
	res, err := app.Test(request, 3000)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response ResponseUpdateProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "invalid id product", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

func TestUpdateProductFailedRequest(t *testing.T) {
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

	brand2 := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier2 := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category2 := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		nameTest        string
		url             string
		body            map[string]any
		expectedCode    int
		expectedStatus  string
		expectedMessage string
	}{
		{
			nameTest: "Failed validation update product min Length ProductName",
			url:      "http://127.0.0.1:8080/api/v1/product/" + product.Id.String(),
			body: map[string]any{
				"product_name": "u",
				"sell_price":   3000,
				"call_name":    "update update",
				"category":     category2.Id,
				"brand":        brand2.Id,
				"supplier":     supplier2.Id,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "ProductName must be at least 2 characters long",
		}, {
			nameTest: "Failed validation update product max Length ProductName",
			url:      "http://127.0.0.1:8080/api/v1/product/" + product.Id.String(),
			body: map[string]any{
				"product_name": "uHac consequat natoque sodales condimentum mus. Velit vitae lacinia integer finibus interdum laoreet condimentum semper. Primis eu nulla a, egestas elementum enim. Ante conubia class ante ornare quis elit sapien blandit ipsum. Vulputate viverssssssssssssssssss",
				"sell_price":   3000,
				"call_name":    "update update",
				"category":     category2.Id,
				"brand":        brand2.Id,
				"supplier":     supplier2.Id,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "ProductName must be maximum 252 characters long",
		}, {
			nameTest: "Failed validation update product Min SellPrice",
			url:      "http://127.0.0.1:8080/api/v1/product/" + product.Id.String(),
			body: map[string]any{
				"product_name": "updated product",
				"sell_price":   -2,
				"call_name":    "update update",
				"category":     category2.Id,
				"brand":        brand2.Id,
				"supplier":     supplier2.Id,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "SellPrice must be greater than 0",
		}, {
			nameTest: "Failed validation update product Min Length CallName",
			url:      "http://127.0.0.1:8080/api/v1/product/" + product.Id.String(),
			body: map[string]any{
				"product_name": "update updated product name",
				"sell_price":   3000,
				"call_name":    "u",
				"category":     category2.Id,
				"brand":        brand2.Id,
				"supplier":     supplier2.Id,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "CallName must be at least 2 characters long",
		}, {
			nameTest: "Failed validation update product Max Length CallName",
			url:      "http://127.0.0.1:8080/api/v1/product/" + product.Id.String(),
			body: map[string]any{
				"product_name": "updated product name",
				"sell_price":   3000,
				"call_name":    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam suscipit vestibulum leo, sit amet porta justo tincidunt vitae. Sed efficitur vitae magna et consequat. Nam placerat euismod facilisis. In sed diam purus. Suspendisse gravida iaculis purus, a suscipit nibh iaculis eget. Interdum et malesuada fames ac ante ipsum primis in faucibus. Aliquam congue laoreet sagittis. Curabitur nisi elit, pharetra a nulla id, venenatis tincidunt libero. Ut fringilla odio vel maximus finibus. In pulvinar maximus tellus, sed bibendum magna elementum sed. Quisque dictum in magna a viverra. Etiam non molestie diam.\n\nNam eleifend diam ut consectetur tempor. Proin eu posuere enim. Integer at elementum dui, id pulvinar dolor. Donec a sapien pharetra, cursus velit vel, finibus neque. Aliquam ligula sem, pretium vitae pulvinar nec, consectetur sed nunc. Quisque eget magna at lectus mattis dictum. Integer ornare neque lacinia, consectetur nunc nec, tincidunt sapien.\n\nAliquam porta ante ut lorem iaculis, quis vehicula turpis iaculis. Pellentesque placerat fermentum in. ",
				"category":     category2.Id,
				"brand":        brand2.Id,
				"supplier":     supplier2.Id,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "CallName must be maximum 1052 characters long",
		}, {
			nameTest: "Failed validation update product Min Brand",
			url:      "http://127.0.0.1:8080/api/v1/product/" + product.Id.String(),
			body: map[string]any{
				"product_name": "updated product",
				"sell_price":   2000,
				"call_name":    "update update",
				"category":     category2.Id,
				"brand":        -5,
				"supplier":     supplier2.Id,
			},
			expectedCode:    fiber.StatusBadRequest,
			expectedStatus:  "failed",
			expectedMessage: "Brand must be greater than 0",
		},
	}

	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			jsonDataRequest, err := json.Marshal(test.body)
			utils.PanicIfError(err)
			request := updateProductTestReq(t, test.url, string(jsonDataRequest), generateToken)
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
			assert.Equalf(t, test.expectedStatus, response.Status, "response status should be equal")
			assert.Equalf(t, test.expectedMessage, response.Message, "response message should be equal")
			assert.Nil(t, response.Data)
		})
	}
}

func TestUpdateProductNotFoundProduct(t *testing.T) {
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

	brand2 := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
	supplier2 := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
	category2 := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
	db.Close()

	reqProduct := map[string]any{
		"product_name": "updated2",
		"sell_price":   3000,
		"call_name":    "update update",
		"category":     category2.Id,
		"brand":        brand2.Id,
		"supplier":     supplier2.Id,
	}
	jsonReqUpdate, err := json.Marshal(reqProduct)
	utils.PanicIfError(err)

	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
	if err != nil {
		panic(err)
	}
	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	id, err := uuid.NewUUID()
	utils.PanicIfError(err)

	url := "http://127.0.0.1:8080/api/v1/product/" + id.String()
	request := updateProductTestReq(t, url, string(jsonReqUpdate), generateToken)
	res, err := app.Test(request, 3000)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

	var response ResponseUpdateProductTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "Product not found", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

//func TestUpdateProductFailedUpdate(t *testing.T) {
//	test.InitConfigTest()
//	db, _, err := test.SetupDBtest()
//	if err != nil {
//		panic(err)
//	}
//
//	err = test.TruncateDB(db)
//	if err != nil {
//		panic(err)
//	}
//
//	req := webrequest.UserCreateRequest{
//		Name:     "testUser",
//		Gender:   "male",
//		Telp:     "08213243444",
//		Birthday: "2023-07-15",
//		Address:  "solo",
//		Email:    "testUser@gmail.com",
//		Password: "password",
//		Username: "testerrr",
//	}
//	user, _, role, _, _, _ := userstest.InsertNewUserTest(t, db, req)
//	brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
//	supplier := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
//	category := categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
//	product := InsertProductTest(db, domain.Product{ProductName: "test product name", SellPrice: 5000, CallName: "test 1, test 2", AdminId: user.User_id, BrandId: brand.Id, CategoryId: category.Id, SupplierId: supplier.Id})
//	_ = photoTest.InsertPhotosTest(db, domain.Photos{Url: "http://127.0.0.1:8080/foto/roti-20240808_210706-11050747_584238255051431_6429195438397655233_o.jpg", Owner: product.Id})
//
//	brand2 := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})
//	supplier2 := suplier.InsertSupplierTest(db, domain.Supplier{Name: "test supplier ", ContactInfo: "testsupplier@gmail.com", Address: "test , south test"})
//	_ = categoryTests.InsertCategoriesTest(db, domain.Category{Name: "test_category", Description: "test category"})
//	db.Close()
//
//	id, err := uuid.NewUUID()
//	utils.PanicIfError(err)
//
//	reqProduct := map[string]any{
//		"product_name": "updated2",
//		"sell_price":   3000,
//		"call_name":    "update update",
//		"category":     id,
//		"brand":        brand2.Id,
//		"supplier":     supplier2.Id,
//	}
//	jsonReqUpdate, err := json.Marshal(reqProduct)
//	utils.PanicIfError(err)
//
//	generateToken, err := utils.GenerateJWT(user.User_id, role.Role)
//	if err != nil {
//		panic(err)
//	}
//	app, clean, err := be.InitializeApp()
//	if err != nil {
//		panic(err)
//	}
//	defer clean()
//
//	url := "http://127.0.0.1:8080/api/v1/product/" + product.Id.String()
//	request := updateProductTestReq(t, url, string(jsonReqUpdate), generateToken)
//	res, err := app.Test(request, 3000)
//	assert.NoError(t, err)
//
//	body, err := ioutil.ReadAll(res.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
//
//	var response ResponseUpdateProductTest
//	err = json.Unmarshal(body, &response)
//
//	if err != nil {
//		log.Fatalf("Error unmarshalling JSON: %v", err)
//	}
//	assert.Equalf(t, "failed", response.Status, "response status should be ok")
//	assert.Equalf(t, "Product not found", response.Message, "response message should be equal")
//	assert.Empty(t, response.Data)
//}
