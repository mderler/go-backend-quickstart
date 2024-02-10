package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mderler/simple-go-backend/internal/db"
	"github.com/mderler/simple-go-backend/internal/validation"
)

type userRequest struct {
	Username string `validate:"required,min=3,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=255"`
}

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if msg := validation.Validate(user); msg != nil {
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	params := db.CreateUserParams(*user)
	dbUser, err := u.queries.CreateUser(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (u *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (u *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if msg := validation.Validate(user); msg != nil {
		http.Error(w, string(msg), http.StatusBadRequest)
		return
	}

	params := db.UpdateUserParams{
		ID:       1,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	dbUser, err := u.queries.UpdateUser(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (u *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// ...
}
