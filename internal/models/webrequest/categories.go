package webrequest

type CategoryCreateReq struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=3"`
}

type CategoryFindByParam struct {
	Params string `json:"params" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"required,min=1,max=30"`
	Offset int    `json:"offset" validate:"required,min=1"`
}
