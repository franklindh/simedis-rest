package utils

import (
	"errors"
	"regexp"

	"github.com/franklindh/simedis-api/internal/model"
)

func ValidatePetugas(petugas model.Petugas) error {
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
