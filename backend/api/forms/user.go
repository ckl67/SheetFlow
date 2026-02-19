package forms

type ResetPasswordRequest struct {
	PasswordResetId string `form:"passwordResetId" validate:"required,len=40"`
	Password        string `form:"password" validate:"required,min=8,max=100"`
}

type RequestResetPasswordRequest struct {
	Email string `form:"email" validate:"required,email,max=100"`
}
