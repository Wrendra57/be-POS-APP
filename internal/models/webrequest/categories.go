package webrequest

type CategoryCreateReq struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=3"`
}
