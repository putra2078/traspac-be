package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenerateNIPWithPrefix(prefix string) string {
	now := time.Now()
	datePart := now.Format("20060102")
	n, _ := rand.Int(rand.Reader, big.NewInt(99999))
	randomPart := n.Int64() + 1
	return fmt.Sprintf("%s-%s-%05d", prefix, datePart, randomPart)
}
