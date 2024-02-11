package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mderler/simple-go-backend/internal/db"
)

type TodoHandler struct {
	*chi.Mux
	queries *db.Queries
}

func NewTodoHandler(queries *db.Queries) *TodoHandler {
	todoHandler := &TodoHandler{chi.NewRouter(), queries}

	todoHandler.Post("/", todoHandler.createTodo)
	todoHandler.Put("/{id}", todoHandler.updateTodo)
	todoHandler.Delete("/{id}", todoHandler.deleteTodo)
	todoHandler.Post("/{id}/assign", todoHandler.assignTodo)

	return todoHandler
}

func (t *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	todo := &todoCreateRequest{}

	if err := decodeAndValidate(w, r, todo); err != nil {
		return
	}

	params := db.CreateTodoParams{
		Title: todo.Title,
		Description: pgtype.Text{
			String: todo.Description,
		},
	}
	dbTodo, err := t.queries.CreateTodo(r.Context(), params)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	resp, err := json.Marshal(dbTodo)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (t *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	todoIdParam := chi.URLParam(r, "id")
	if todoIdParam == "" {
		writeTodoNotFoundError(w, todoIdParam)
		return
	}

	todoId, err := strconv.Atoi(todoIdParam)
	if err != nil {
		writeInvalidTodoIdError(w, todoIdParam)
		return
	}

	todo := &todoUpdateRequest{}

	if err := decodeAndValidate(w, r, todo); err != nil {
		return
	}

	params := db.UpdateTodoParams{
		ID:    int32(todoId),
		Title: todo.Title,
		Description: pgtype.Text{
			String: todo.Description,
		},
		Completed: pgtype.Bool{
			Bool: todo.Completed,
		},
	}
	dbTodo, err := t.queries.UpdateTodo(r.Context(), params)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	resp, err := json.Marshal(dbTodo)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (t *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todoIdParam := chi.URLParam(r, "id")
	if todoIdParam == "" {
		writeTodoNotFoundError(w, todoIdParam)
		return
	}

	todoId, err := strconv.Atoi(todoIdParam)
	if err != nil {
		writeInvalidTodoIdError(w, todoIdParam)
		return
	}

	affectedRows, err := t.queries.DeleteTodo(r.Context(), int32(todoId))
	if err != nil {
		writeInternalServerError(w, err)
		return
	}
	if affectedRows == 0 {
		writeTodoNotFoundError(w, todoIdParam)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TodoHandler) assignTodo(w http.ResponseWriter, r *http.Request) {
	todoIdParam := chi.URLParam(r, "id")
	if todoIdParam == "" {
		writeTodoNotFoundError(w, todoIdParam)
		return
	}

	todoId, err := strconv.Atoi(todoIdParam)
	if err != nil {
		writeInvalidTodoIdError(w, todoIdParam)
		return
	}

	assign := &todoAssignRequest{}

	if err := decodeAndValidate(w, r, assign); err != nil {
		return
	}

	params := db.AssignUserToTodoParams{
		TodoID: int32(todoId),
		UserID: assign.UserID,
	}
	_, err = t.queries.AssignUserToTodo(r.Context(), params)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
