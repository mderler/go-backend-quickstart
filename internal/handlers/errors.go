package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func writeJsonDecodeError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(`{"error":"json-decode-error","title":"Your request payload didn't decode","detail":"%s"}`, err.Error())
	log.Println("Error decoding JSON:", err)
	http.Error(w, msg, http.StatusBadRequest)
}

func writeInternalServerError(w http.ResponseWriter, err error) {
	msg := `{"error":"internal-server-error","title":"Something went wrong"}`
	log.Println("Internal server error:", err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func writeUserNotFoundError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"error":"user-not-found","title":"User not found","detail":"User with id %s not found"}`, id)
	log.Println("User not found:", id)
	http.Error(w, msg, http.StatusNotFound)
}

func writeInvalidUserIdError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"error":"invalid-user-id","title":"Invalid user id","detail":"The user id %s is not valid"}`, id)
	log.Println("Invalid user id:", id)
	http.Error(w, msg, http.StatusBadRequest)
}

func writeTodoNotFoundError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"error":"todo-not-found","title":"Todo not found","detail":"Todo with id %s not found"}`, id)
	log.Println("Todo not found:", id)
	http.Error(w, msg, http.StatusNotFound)
}

func writeInvalidTodoIdError(w http.ResponseWriter, id string) {
	msg := fmt.Sprintf(`{"error":"invalid-todo-id","title":"Invalid todo id","detail":"The todo id %s is not valid"}`, id)
	log.Println("Invalid todo id:", id)
	http.Error(w, msg, http.StatusBadRequest)
}
