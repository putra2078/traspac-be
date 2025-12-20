package utils

import (
	"crypto/rand"
	"math/big"
)

func GeneratePassCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	passCode := make([]byte, length)
	for i := range passCode {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		passCode[i] = charset[n.Int64()]
	}
	return string(passCode)
}

func GenerateShortPasscode() string {
	return GeneratePassCode(6)
}

func GenerateLinkJoin() string {
	return "https://traspac.com/join/" + GenerateShortPasscode()
}


