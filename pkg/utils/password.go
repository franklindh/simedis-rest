package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 2
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
	saltBytes    uint32 = 16
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be blank")
	}
	salt := make([]byte, saltBytes)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodeHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)

	return encodeHash, nil
}

func VerifyPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return errors.New("invalid encoded hash format")
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}

	decodedHash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	hashToCompare := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	if subtle.ConstantTimeCompare(hashToCompare, decodedHash) == 1 {
		return nil
	}

	return errors.New("incorrect password")
}
