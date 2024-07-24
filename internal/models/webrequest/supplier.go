package webrequest

type SupplierRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=32"`
	ContactInfo string `json:"contact_info" validate:"required,min=3,max=32"`
	Address     string `json:"address" validate:"required,min=3,max=232"`
}
type SupplierListRequest struct {
	Params string `json:"params" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"required,min=1,max=30"`
	Offset int    `json:"offset" validate:"required,min=1"`
}
