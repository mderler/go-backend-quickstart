package validation

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Validate(s interface{}) []byte {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errorResponse := make(map[string]interface{})
	errorResponse["type"] = "validation-error"
	errorResponse["detail"] = "Your request parameters didn't validate."

	invalidParams := make(map[string]interface{})
	for _, e := range err.(validator.ValidationErrors) {
		fieldName := e.Field()
		tag := e.Tag()
		message := fmt.Sprintf("The %s %s", fieldName, getErrorMessage(tag))
		invalidParams[fieldName] = map[string]interface{}{"detail": message, "tag": tag}
	}
	errorResponse["invalid-params"] = invalidParams

	msg, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println("Error marshalling error response:", err)
		return []byte(`{"type": "validation-error", "detail": "Your request parameters didn't validate."}`)
	}

	log.Println("Validation error:", string(msg))

	return msg
}

func getErrorMessage(tag string) string {
	switch tag {
	case "required":
		return "field is required"
	case "min":
		return "value is too short"
	case "email":
		return "field must be a valid email"
	default:
		return "invalid value"
	}
}
