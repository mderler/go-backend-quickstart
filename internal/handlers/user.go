package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/mderler/simple-go-backend/internal/db"
	"github.com/mderler/simple-go-backend/internal/validation"
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

	return userHandler
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		writeJsonDecodeError(w, err)
		return
	}

	if msg := validation.Validate(user); msg != nil {
		http.Error(w, string(msg), http.StatusBadRequest)
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

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		writeJsonDecodeError(w, err)
		return
	}

	if msg := validation.Validate(user); msg != nil {
		http.Error(w, string(msg), http.StatusBadRequest)
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
