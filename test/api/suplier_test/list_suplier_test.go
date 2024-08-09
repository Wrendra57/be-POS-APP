package suplier

import (
	"encoding/json"
	be "github.com/Wrendra57/Pos-app-be"
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
	Data    []domain.Supplier `json:"data"`
	Message string            `json:"message"`
}

func GetListSupplier(t *testing.T, url, body string) *http.Request {
	req, err := http.NewRequest("GET", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestGetListBrandSuccess(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		_ = InsertSupplierTest(db, supplier)
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/supplier?" + limit + "&" + offset
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "Success get data supplier", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
}

func TestGetListSupplierWitParamsSuccess(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=1"
	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/supplier?" + limit + "&" + offset + "&" + param
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "Success get data supplier", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	assert.Lenf(t, response.Data, 1, "should have one brand")
}

func TestGetListSupplierWitParamsSuccessNoResult(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	param := "params=1dwfw"
	limit := "limit=10"
	offset := "offset=1"
	url := "/api/v1/supplier?" + limit + "&" + offset + "&" + param
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "Success get data supplier", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
	assert.Lenf(t, response.Data, 0, "should have one brand")
}

func TestGetListSupplierWithoutLimit(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	offset := "offset=1"
	url := "/api/v1/supplier?" + offset
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "Success get data supplier", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	assert.Lenf(t, response.Data, 5, "should have one brand")
}

func TestGetListSupplierWithoutOffset(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()

	url := "/api/v1/supplier?"
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "Success get data supplier", response.Message, "response message should be equal")
	assert.NotEmpty(t, response.Data)
	assert.Lenf(t, response.Data, 5, "should have one brand")
}

func TestGetListSupplierFailedLimit(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	limit := "limit=w2"
	url := "/api/v1/supplier?" + limit
	request := GetListSupplier(t, url, "")
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

func TestGetListSupplierFailedNegativeLimit(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	limit := "limit=-2"
	url := "/api/v1/supplier?" + limit
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "The 'limit' field must be greater than zero", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}

func TestGetListSupplierFailedOffset(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	offset := "offset=w2"
	url := "/api/v1/supplier?" + offset
	request := GetListSupplier(t, url, "")
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
func TestGetListSupplierFailedNegativeOffset(t *testing.T) {
	test.InitConfigTest()

	db, _, err := test.SetupDBtest()
	if err != nil {
		panic(err)
	}

	err = test.TruncateDB(db)
	if err != nil {
		panic(err)
	}
	supplier := domain.Supplier{Name: "test", ContactInfo: "test supplier", Address: "test supplier"}

	for i := 0; i < 5; i++ {
		name := supplier.Name + strconv.Itoa(i)
		_ = InsertSupplierTest(db, domain.Supplier{Name: name, ContactInfo: "test supplier", Address: "test supplier"})
	}
	db.Close()

	app, clean, err := be.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer clean()
	offset := "offset=-2"
	url := "/api/v1/supplier?" + offset
	request := GetListSupplier(t, url, "")
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
	assert.Equalf(t, "The 'offset' field must be positive", response.Message, "response message should be equal")
	assert.Empty(t, response.Data)
}
