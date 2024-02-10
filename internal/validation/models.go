package validation

import "github.com/jackc/pgx/v5/pgtype"

type CreateTodo struct {
	UserID      int32       `validate:"required,min=1"`
	Title       string      `validate:"required,min=1,max=255"`
	Description pgtype.Text `validate:"required"`
}

type CreateUser struct {
	Username string `validate:"required,min=3,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=255"`
}
