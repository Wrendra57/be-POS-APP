package webrespones

import (
	"github.com/google/uuid"
	"time"
)

type photo struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}
type ProductFindDetail struct {
	Id                  uuid.UUID `json:"id"`
	ProductName         string    `json:"product_name"`
	SellPrice           int       `json:"sell_price"`
	CallName            string    `json:"call_name"`
	AdminId             uuid.UUID `json:"admin_id"`
	AdminName           string    `json:"admin_name"`
	CategoryId          uuid.UUID `json:"category_id"`
	CategoryName        string    `json:"category_name"`
	CategoryDescription string    `json:"category_description"`
	BrandId             int       `json:"brand_id"`
	BrandName           string    `json:"brand_name"`
	BrandDescription    string    `json:"brand_description"`
	SupplierId          uuid.UUID `json:"supplier_id"`
	SupplierName        string    `json:"supplier_name"`
	SupplierContactInfo string    `json:"supplier_contact_info"`
	SupplierAddress     string    `json:"supplier_address"`
	Photos              []photo   `json:"photos"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type user struct {
	AdminId   uuid.UUID `json:"admin_id"`
	AdminName string    `json:"admin_name"`
}
type category struct {
	CategoryId          uuid.UUID `json:"category_id"`
	CategoryName        string    `json:"category_name"`
	CategoryDescription string    `json:"category_description"`
}
type brand struct {
	BrandId          int    `json:"brand_id"`
	BrandName        string `json:"brand_name"`
	BrandDescription string `json:"brand_description"`
}
type supplier struct {
	SupplierId          uuid.UUID `json:"supplier_id"`
	SupplierName        string    `json:"supplier_name"`
	SupplierContactInfo string    `json:"supplier_contact_info"`
	SupplierAddress     string    `json:"supplier_address"`
}

type ProductFindByIdResponseApi struct {
	Id          uuid.UUID `json:"id"`
	ProductName string    `json:"product_name"`
	SellPrice   int       `json:"sell_price"`
	CallName    string    `json:"call_name"`
	Admin       user      `json:"admin"`
	Category    category  `json:"category"`
	Brand       brand     `json:"brand"`
	Supplier    supplier  `json:"supplier"`
	Photos      []photo   `json:"photos"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
