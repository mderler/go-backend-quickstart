package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/mderler/simple-go-backend/internal/db"
)

type UserHandler struct {
	*chi.Mux
	queries *db.Queries
}

func NewUserHandler(queries *db.Queries) *UserHandler {
	userHandler := &UserHandler{chi.NewRouter(), queries}

	userHandler.Post("/", userHandler.createUser)
	userHandler.Get("/", userHandler.getUsers)
	userHandler.Put("/{id}", userHandler.updateUser)
	userHandler.Delete("/{id}", userHandler.deleteUser)
	userHandler.Get("/{id}/todos", userHandler.getUserTodos)

	return userHandler
}

// createUser creates a new user.
//
// @Summary Create a new user
// @Description Create a new user with the provided user data.
// @Tags User
// @Accept json
// @Produce json
// @Param user body userRequest true "User data"
// @Success 201 {object} db.User "Created user"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /user [post]
func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}

	if err := decodeAndValidate(w, r, user); err != nil {
		return
	}

	params := db.CreateUserParams(*user)
	dbUser, err := u.queries.CreateUser(r.Context(), params)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	resp, err := json.Marshal(dbUser)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// getUsers returns the list of all users.
//
// @Summary Get all users
// @Description Get the list of all users.
// @Tags User
// @Produce json
// @Success 200 {array} db.User "List of users"
// @Failure 500 "Internal server error"
// @Router /user [get]
func (u *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.queries.ListUsers(r.Context())
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.Write(resp)
}

// updateUser updates an existing user.
//
// @Summary Update an existing user
// @Description Update an existing user with the provided user data.
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body userRequest true "User data"
// @Success 200 {object} db.User "Updated user"
// @Failure 400 "Bad request"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /user/{id} [put]
func (u *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "id")
	if userIdParam == "" {
		writeUserNotFoundError(w, userIdParam)
		return
	}

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		writeInvalidUserIdError(w, userIdParam)
		return
	}

	user := &userRequest{}

	if err := decodeAndValidate(w, r, user); err != nil {
		return
	}

	params := db.UpdateUserParams{
		ID:       int32(userId),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	dbUser, err := u.queries.UpdateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeUserNotFoundError(w, userIdParam)
			return
		}
		writeInternalServerError(w, err)
		return
	}

	resp, err := json.Marshal(dbUser)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.Write(resp)
}

// deleteUser deletes an existing user.
//
// @Summary Delete an existing user
// @Description Delete an existing user with the provided user ID.
// @Tags User
// @Param id path int true "User ID"
// @Success 204 "No content"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /user/{id} [delete]
func (u *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "id")
	if userIdParam == "" {
		writeUserNotFoundError(w, userIdParam)
		return
	}

	userID, err := strconv.Atoi(userIdParam)
	if err != nil {
		writeInvalidUserIdError(w, userIdParam)
		return
	}

	affectedRows, err := u.queries.DeleteUser(r.Context(), int32(userID))
	if err != nil {
		writeInternalServerError(w, err)
		return
	}
	if affectedRows == 0 {
		writeUserNotFoundError(w, userIdParam)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getUserTodos returns the list of all todos of a user.
//
// @Summary Get all todos of a user
// @Description Get the list of all todos of a user with the provided user ID.
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} db.Todo "List of todos"
// @Failure 400 "Bad request"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /user/{id}/todos [get]
func (u *UserHandler) getUserTodos(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "id")
	if userIdParam == "" {
		writeUserNotFoundError(w, userIdParam)
		return
	}

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		writeInvalidUserIdError(w, userIdParam)
		return
	}

	todos, err := u.queries.GetTodosOfUser(r.Context(), int32(userId))
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	resp, err := json.Marshal(todos)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.Write(resp)
}
