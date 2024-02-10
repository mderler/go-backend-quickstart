package handlers

type userRequest struct {
	Username string `validate:"required,min=3,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=255"`
}
