package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomID(byteLen int) (string, error) {
	if byteLen <= 0 {
		return "", fmt.Errorf("byteLen must be > 0")
	}

	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func MustRandomID() string {
	id, err := GenerateRandomID(8) // 8 bytes = 16 hex chars
	if err != nil {
		panic(err)
	}
	return id
}
