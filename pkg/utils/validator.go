package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {

		var errorMessages []string

		for _, e := range validationErrors {
			var message string
			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("Field '%s' is required", e.Field())
			case "min":
				message = fmt.Sprintf("Field '%s' must be at least %s characters long", e.Field(), e.Param())
			case "max":
				message = fmt.Sprintf("Field '%s' must not exceed %s characters", e.Field(), e.Param())
			case "oneof":
				message = fmt.Sprintf("Field '%s' must be one of [%s]", e.Field(), e.Param())
			default:
				message = fmt.Sprintf("Field '%s' is not valid", e.Field())
			}

			errorMessages = append(errorMessages, message)
		}

		return strings.Join(errorMessages, ", ")
	}
	return err.Error()
}

func ValidatePetugasUsername(petugas model.Petugas) error {
	if petugas.Username == "" || petugas.Name == "" {
		return errors.New("username and name are required fields")
	}

	if len(petugas.Username) < 5 || len(petugas.Username) > 20 {
		return errors.New("username must be between 5 and 20 characters")
	}

	isValidUsername, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", petugas.Username)
	if !isValidUsername {
		return errors.New("username can only contain letters, numbers, underscores, and dashes")
	}

	return nil
}
