package webrequest

type BrandCreateReq struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=3"`
}

type BrandGetRequest struct {
	Params string `json:"params" validate:"omitempty,min=1"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=30"`
	Offset int    `json:"offset" validate:"omitempty,min=1"`
}
