package sctl

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateToken Creates a random token consiting of aspecified number of sha256 hashes
func GenerateToken(length int) string {
	tokenParts := make([]string, 0)
	for i := 0; i < length; i++ {
		tokenParts = append(tokenParts, TokenPart())
	}
	return strings.Join(tokenParts, "-")
}

// TokenPart Generates a random sha256 hash and returns as a string
func TokenPart() string {
	randBytes := make([]byte, 64)
	rand.Seed(time.Now().UnixNano())
	rand.Read(randBytes)
	return fmt.Sprintf("%x", sha256.Sum256(randBytes))
}
