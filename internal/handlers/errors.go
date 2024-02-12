package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorType string

const (
	JSONDecodeError    ErrorType = "json-decode-error"
	UserNotFoundError  ErrorType = "user-not-found"
	InvalidUserIdError ErrorType = "invalid-user-id"
	TodoNotFoundError  ErrorType = "todo-not-found"
	InvalidTodoIdError ErrorType = "invalid-todo-id"
	InvalidQueryError  ErrorType = "invalid-query"
	TodoAssignError    ErrorType = "todo-assign-error"
)

type InternalErrorResponse struct {
	Type  string `json:"type" enums:"internal-server-error"`
	Title string `json:"title"`
}

type ErrorResponse struct {
	Title  string    `json:"title"`
	Type   ErrorType `json:"type"`
	Detail string    `json:"detail"`
}

type ValidationErrorResponse struct {
	Type          string                  `json:"type" enums:"validation-error"`
	InvalidParams map[string]InvalidParam `json:"invalid_params"`
	Detail        string                  `json:"detail"`
}

type InvalidParam struct {
	Message string `json:"message"`
	Tag     string `json:"tag"`
}

func WriteJsonDecodeError(w http.ResponseWriter, err error) {
	errResponse := ErrorResponse{
		Type:   JSONDecodeError,
		Title:  "Your request payload didn't decode",
		Detail: err.Error(),
	}
	log.Println("Error decoding JSON:", err)
	writeJson(w, errResponse, http.StatusBadRequest)
}

func writeInternalServerError(w http.ResponseWriter, err error) {
	errResponse := InternalErrorResponse{
		Title: "Something went wrong",
		Type:  "internal-server-error",
	}
	log.Println("Internal server error:", err)
	writeJson(w, errResponse, http.StatusInternalServerError)
}

func writeUserNotFoundError(w http.ResponseWriter, id int32) {
	errResponse := ErrorResponse{
		Type:   UserNotFoundError,
		Title:  "User not found",
		Detail: fmt.Sprintf("User with id %d not found", id),
	}
	log.Println("User not found:", id)
	writeJson(w, errResponse, http.StatusNotFound)
}

func writeInvalidUserIdError(w http.ResponseWriter, id string) {
	errResponse := ErrorResponse{
		Type:   InvalidUserIdError,
		Title:  "Invalid user id",
		Detail: fmt.Sprintf("The user id %s is not valid", id),
	}
	log.Println("Invalid user id:", id)
	writeJson(w, errResponse, http.StatusBadRequest)
}

func writeTodoNotFoundError(w http.ResponseWriter, id int32) {
	errResponse := ErrorResponse{
		Type:   TodoNotFoundError,
		Title:  "Todo not found",
		Detail: fmt.Sprintf("Todo with id %d not found", id),
	}
	log.Println("Todo not found:", id)
	writeJson(w, errResponse, http.StatusNotFound)
}

func writeInvalidTodoIdError(w http.ResponseWriter, id string) {
	errResponse := ErrorResponse{
		Type:   InvalidTodoIdError,
		Title:  "Invalid todo id",
		Detail: fmt.Sprintf("The todo id %s is not valid", id),
	}
	log.Println("Invalid todo id:", id)
	writeJson(w, errResponse, http.StatusBadRequest)
}

func writeInvalidTodoAssignRequestError(w http.ResponseWriter, err *pgconn.PgError, todoId int32, userID int32) {
	log.Println("Invalid todo assign request:", err)
	if err.ConstraintName == "todo_user_user_id_fkey" {
		writeUserNotFoundError(w, userID)
		return
	} else if err.ConstraintName == "todo_user_todo_id_fkey" {
		writeTodoNotFoundError(w, todoId)
		return
	}
	writeInternalServerError(w, err)
}

func writeDuplicateTodoAssignRequestError(w http.ResponseWriter) {
	errResponse := ErrorResponse{
		Type:   TodoAssignError,
		Title:  "User already assigned",
		Detail: "The user is already assigned to the todo",
	}
	writeJson(w, errResponse, http.StatusConflict)
}

func writeInvalidQueryError(w http.ResponseWriter, actual string, options []string) {
	errResponse := ErrorResponse{
		Type:   InvalidQueryError,
		Title:  "Invalid query",
		Detail: fmt.Sprintf("The query parameter %s is not valid. Valid options are %v", actual, options),
	}
	log.Printf("Invalid query: actual=%s options=%v\n", actual, options)
	writeJson(w, errResponse, http.StatusBadRequest)
}
