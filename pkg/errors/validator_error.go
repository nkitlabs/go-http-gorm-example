package errors

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func ConvertValidatorErrorsToError(err error) *Error {
	if err == nil {
		return nil
	}

	validatorErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrInternal.WithMessage(err.Error())
	}

	errMaps := make(map[string]string)
	for _, e := range validatorErrs {
		errMaps[e.Field()] = validatorTagToMsg(e.Tag())
	}

	jsonString, err := json.Marshal(errMaps)
	if err != nil {
		return ErrInternal.WithMessage(err.Error()) // unlikely to get this error.
	}

	return ErrInvalidInput.WithMessage(string(jsonString))
}

func validatorTagToMsg(tag string) string {
	switch tag {
	case "required":
		return "It is required"
	case "email":
		return "It is not a valid email address"
	default:
		return "It is invalid"
	}
}
