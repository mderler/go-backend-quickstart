package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mderler/simple-go-backend/internal/validation"
)

func decodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeJsonDecodeError(w, err)
		return err
	}

	if msg := validation.Validate(v); msg != nil {
		http.Error(w, string(msg), http.StatusBadRequest)
		return errors.New("validation error")
	}

	return nil
}
