package handlers

type userRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type todoCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,max=1000"`
	CreatorID   int32  `json:"creatorId" validate:"required"`
}

type todoUpdateRequest struct {
	*todoCreateRequest
	Completed bool `json:"completed" validate:"required"`
}

type todoAssignRequest struct {
	UserID int32 `json:"userId" validate:"required"`
}
