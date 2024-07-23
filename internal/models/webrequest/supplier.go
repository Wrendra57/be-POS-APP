package webrequest

type SupplierRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=32"`
	ContactInfo string `json:"contact_info" validate:"required,min=3,max=32"`
	Address     string `json:"address" validate:"required,min=3,max=232"`
}
