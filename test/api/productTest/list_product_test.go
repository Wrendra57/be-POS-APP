package productTest

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/test"
	"github.com/Wrendra57/Pos-app-be/test/api/brandTest"
	"net/http"
	"strings"
	"testing"
)

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

	brand := brandTest.InsertBrandTest(db, domain.Brand{Name: "test_brand", Description: "test brand"})

}
