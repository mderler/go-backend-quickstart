package handlers

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type TodoCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,max=1000"`
	CreatorID   int32  `json:"creatorId" validate:"required"`
}

type TodoUpdateRequest struct {
	*TodoCreateRequest
	Completed bool `json:"completed" validate:"required"`
}

type TodoAssignRequest struct {
	UserID int32 `json:"userId" validate:"required"`
}
