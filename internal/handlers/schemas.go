package handlers

type userRequest struct {
	Username string `validate:"required,min=3,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=255"`
}

type todoCreateRequest struct {
	Title       string `validate:"required,min=1,max=255"`
	Description string `validate:"required,max=1000"`
	CreatorID   int32  `validate:"required"`
}

type todoUpdateRequest struct {
	*todoCreateRequest
	Completed bool `validate:"required"`
}

type todoAssignRequest struct {
	UserID int32 `validate:"required"`
}
