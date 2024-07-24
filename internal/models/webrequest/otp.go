package webrequest

type ValidateOtpRequest struct {
	Token string `json:"token" validate:"required"`
	Otp   string `json:"otp" validate:"required,min=6,max=6"`
}
