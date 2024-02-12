package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mderler/simple-go-backend/internal/db"
)

type TodoHandler struct {
	*chi.Mux
	queries *db.Queries
}

func NewTodoHandler(queries *db.Queries) *TodoHandler {
	todoHandler := &TodoHandler{chi.NewRouter(), queries}

	todoHandler.Post("/", todoHandler.createTodo)
	todoHandler.Get("/", todoHandler.getTodos)
	todoHandler.Put("/{id}", todoHandler.updateTodo)
	todoHandler.Delete("/{id}", todoHandler.deleteTodo)
	todoHandler.Post("/{id}/assign", todoHandler.assignTodo)

	return todoHandler
}

// @Summary Create a new todo
// @Description Create a new todo with the provided todo data.
// @Tags Todo
// @Accept json
// @Produce json
// @Param todo body TodoCreateRequest true "Todo data"
// @Success 201 {object} db.Todo "Created todo"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 422 {object} ValidationErrorResponse "Bad request"
// @Failure 500 {object} InternalErrorResponse "Internal server error"
// @Router /todo [post]
func (t *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	todo := &TodoCreateRequest{}

	if err := decodeAndValidate(w, r, todo); err != nil {
		return
	}

	params := db.CreateTodoParams{
		Title:       todo.Title,
		Description: todo.Description,
		CreatorID:   todo.CreatorID,
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

// @Summary Get all todos
// @Description Get the list of all todos.
// @Tags Todo
// @Produce json
// @Success 200 {array} db.Todo "List of todos"
// @Failure 500 {object} InternalErrorResponse "Internal server error"
// @Router /todo [get]
func (t *TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := t.queries.ListTodos(r.Context())
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

// @Summary Update a todo
// @Description Update an existing todo with the provided todo data.
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body TodoUpdateRequest true "Todo data"
// @Success 200 {object} db.Todo "Updated todo"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Todo not found"
// @Failure 422 {object} ValidationErrorResponse "Validation error"
// @Failure 500 {object} InternalErrorResponse "Internal server error"
// @Router /todo/{id} [put]
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

	todo := &TodoUpdateRequest{}

	if err := decodeAndValidate(w, r, todo); err != nil {
		return
	}

	params := db.UpdateTodoParams{
		ID:          int32(todoId),
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
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

// @Summary Delete a todo
// @Description Delete an existing todo.
// @Tags Todo
// @Param id path int true "Todo ID"
// @Success 204 "No content"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Todo not found"
// @Failure 500 {object} InternalErrorResponse "Internal server error"
// @Router /todo/{id} [delete]
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

// @Summary Assign a user to a todo
// @Description Assign a user to a todo.
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body TodoAssignRequest true "User data"
// @Success 201 "Created todo assignment"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Todo or User not found"
// @Failure 422 {object} ValidationErrorResponse "Validation error"
// @Failure 500 {object} InternalErrorResponse "Internal server error"
// @Router /todo/{id}/assign [post]
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

	assign := &TodoAssignRequest{}

	if err := decodeAndValidate(w, r, assign); err != nil {
		return
	}

	params := db.AssignUserToTodoParams{
		TodoID: int32(todoId),
		UserID: assign.UserID,
	}
	_, err = t.queries.AssignUserToTodo(r.Context(), params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23503" {
			writeInvalidTodoAssignRequestError(w, pgErr, params.TodoID, assign.UserID)
			return
		}
		writeInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
