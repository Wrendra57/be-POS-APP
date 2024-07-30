package categoryTests

import (
	"encoding/json"
	"github.com/Wrendra57/Pos-app-be/cmd"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type responeListTest struct {
	Code    int               `json:"code"`
	Status  string            `json:"status"`
	Data    []domain.Category `json:"data"`
	Message string            `json:"message"`
}

func GetLisCategoryTest(t *testing.T, url, body string) *http.Request {
	req, err := http.NewRequest("GET", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestGetListCategoriesSuccess(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		_ = InsertCategoriesTest(db, category)
	}

	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/categories?" + limit + "&" + offset
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data categories", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestGetListBrandWitParamsSucces(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		name := category.Name + strconv.Itoa(i)
		_ = InsertCategoriesTest(db, domain.Category{Name: name, Description: "test brand"})
	}
	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=1"
	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/categories?" + limit + "&" + offset + "&" + param
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data categories", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	assert.Lenf(t, response.Data, 1, "should have one brand")
}

func TestGetListCategoriesWitParamsSuccessNoResult(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		name := category.Name + strconv.Itoa(i)
		_ = InsertCategoriesTest(db, domain.Category{Name: name, Description: "test brand"})
	}
	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=1dwfw"
	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/categories?" + limit + "&" + offset + "&" + param
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data categories", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	assert.Lenf(t, response.Data, 0, "should have one brand")
}

func TestGetListCategoriesWithoutLimit(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		name := category.Name + strconv.Itoa(i)
		_ = InsertCategoriesTest(db, domain.Category{Name: name, Description: "test brand"})
	}
	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	offset := "offset=1"
	url := "/api/v1/categories?" + offset
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data categories", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	assert.Lenf(t, response.Data, 5, "should have one brand")
}

func TestGetListCategoriesWithoutOffset(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		name := category.Name + strconv.Itoa(i)
		_ = InsertCategoriesTest(db, domain.Category{Name: name, Description: "test brand"})
	}
	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	url := "/api/v1/categories?"
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "success", response.Status, "response status should be ok")
	assert.Equalf(t, "Success get data categories", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	assert.Lenf(t, response.Data, 5, "should have one brand")
}

func TestGetListCategoriesFailedLimit(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		name := category.Name + strconv.Itoa(i)
		_ = InsertCategoriesTest(db, domain.Category{Name: name, Description: "test brand"})
	}
	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	limit := "limit=w2"
	url := "/api/v1/categories?" + limit
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "The 'limit' field must be number/integer", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

func TestGetListCategoriesFailedOffset(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	category := domain.Category{Name: "test", Description: "test category"}

	for i := 0; i < 5; i++ {
		name := category.Name + strconv.Itoa(i)
		_ = InsertCategoriesTest(db, domain.Category{Name: name, Description: "test brand"})
	}
	db.Close()

	app, clean, err := main.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	offset := "offset=w2"
	url := "/api/v1/categories?" + offset
	request := GetLisCategoryTest(t, url, "")
	res, err := app.Test(request, 3000)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	var response responeListTest
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	assert.Equalf(t, "failed", response.Status, "response status should be ok")
	assert.Equalf(t, "The 'offset' field must be number/integer", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}
