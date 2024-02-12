package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func decodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		WriteJsonDecodeError(w, err)
		return false
	}

	if msg := Validate(v); msg != nil {
		writeJson(w, msg, http.StatusUnprocessableEntity)
		return false
	}

	return true
}

func writeJson(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("Error encoding JSON:", err)
		w.Write([]byte(`{"type":"internal-server-error","title":"Something went wrong"}`))
	}
}
