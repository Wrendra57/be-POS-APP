package webrequest

type ValidateOtpRequest struct {
	Otp string `json:"otp" validate:"required,min=6,max=6"`
}
