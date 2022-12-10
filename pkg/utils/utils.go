package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	return base64.URLEncoding.EncodeToString(bytes), err
}

func GenerateRandomState() (string, error) {
	return GenerateRandomString(32)
}
