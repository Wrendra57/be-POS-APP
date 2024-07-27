package exception

type CustomEror struct {
	Code  int
	Error string
}

type ValidationErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}
