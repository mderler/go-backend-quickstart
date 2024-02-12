package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Validate(s interface{}) *ValidationErrorResponse {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errorResponse := &ValidationErrorResponse{Type: "validation-error", Detail: "Your request parameters didn't validate."}

	invalidParams := make(map[string]InvalidParam)
	for _, e := range err.(validator.ValidationErrors) {
		fieldName := e.Field()
		tag := e.Tag()
		message := fmt.Sprintf("The %s %s", fieldName, getErrorMessage(tag))
		invalidParams[fieldName] = InvalidParam{Message: message, Tag: tag}
	}
	errorResponse.InvalidParams = invalidParams

	return errorResponse
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
