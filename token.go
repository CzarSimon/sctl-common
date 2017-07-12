package sctl

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Token Struct to keep a token and its timestamp
type Token struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// Valid checks if a token in valid or out of date
type (token Token) Valid() bool {
	var maxAge float64 = 300.0
	now := time.Now().UTC()
	difference := now.Sub(token.Timestamp)
	return difference.Seconds() <= maxAge
}

// NewToken Return a new random token with timestamp set to the current UTC time
func NewToken(length int) Token {
	return Token{
		Data:      GenrateToken(length),
		Timestamp: time.Now().UTC().
	}
}

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

// TokenSchema Returns the token schema
func TokenSchema() string {
	return `CREATE TABLE TOKEN(
		AUTH VARCHAR(64)
		MASTER VARCHAR(260)
		AUTH_TIMESTAMP TIMESTAMP
	)		
	`
}
