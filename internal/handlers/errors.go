package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type InternalErrorResponse struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type ErrorResponse struct {
	*InternalErrorResponse
	Detail string `json:"detail"`
}

type ValidationErrorResponse struct {
	*ErrorResponse
	InvalidParams map[string]InvalidParam `json:"invalid_params"`
}

type InvalidParam struct {
	Detail string `json:"detail"`
	Tag    string `json:"tag"`
}

func writeJsonDecodeError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(`{"type":"json-decode-error","title":"Your request payload didn't decode","detail":"%s"}`, err.Error())
	log.Println("Error decoding JSON:", err)
	http.Error(w, msg, http.StatusBadRequest)
}

func writeInternalServerError(w http.ResponseWriter, err error) {
	msg := `{"type":"internal-server-error","title":"Something went wrong"}`
	log.Println("Internal server error:", err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func writeUserNotFoundError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"type":"user-not-found","title":"User not found","detail":"User with id %s not found"}`, id)
	log.Println("User not found:", id)
	http.Error(w, msg, http.StatusNotFound)
}

func writeInvalidUserIdError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"type":"invalid-user-id","title":"Invalid user id","detail":"The user id %s is not valid"}`, id)
	log.Println("Invalid user id:", id)
	http.Error(w, msg, http.StatusBadRequest)
}

func writeTodoNotFoundError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"type":"todo-not-found","title":"Todo not found","detail":"Todo with id %s not found"}`, id)
	log.Println("Todo not found:", id)
	http.Error(w, msg, http.StatusNotFound)
}

func writeInvalidTodoIdError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"type":"invalid-todo-id","title":"Invalid todo id","detail":"The todo id %s is not valid"}`, id)
	log.Println("Invalid todo id:", id)
	http.Error(w, msg, http.StatusBadRequest)
}

func writeInvalidTodoAssignRequestError(w http.ResponseWriter, err *pgconn.PgError, todoId int32, userId int32) {
	log.Println("Invalid todo assign request:", err)
	if err.ConstraintName == "todo_user_user_id_fkey" {
		writeUserNotFoundError(w, fmt.Sprint(userId))
		return
	} else if err.ConstraintName == "todo_user_todo_id_fkey" {
		writeTodoNotFoundError(w, fmt.Sprint(todoId))
		return
	}
	writeInternalServerError(w, err)
}

func writeInvalidQueryError(w http.ResponseWriter, actual string, options []string) {
	log.Printf("Invalid query: actual=%s options=%v\n", actual, options)
	msg := fmt.Sprintf(`{"type":"invalid-query","title":"Invalid query","detail":"The query parameter %s is not valid. Valid options are %v"}`, actual, options)
	http.Error(w, msg, http.StatusBadRequest)
}
