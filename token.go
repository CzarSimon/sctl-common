package sctl

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/CzarSimon/util"
	_ "github.com/mattn/go-sqlite3" // sqlite driver for token-db
)

// Token Struct to keep a token and its timestamp
type Token struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// TokenBundle Struct to keep all project tokens
type TokenBundle struct {
	Auth   Token
	Master string
}

// ToBundle Turns a token to a TokenBundle with a given master token
func (token Token) ToBundle(masterToken string) TokenBundle {
	return TokenBundle{
		Auth:   token,
		Master: masterToken,
	}
}

// Valid checks if a token in valid or out of date
func (token Token) Valid() bool {
	maxAge := 300.0
	now := time.Now().UTC()
	difference := now.Sub(token.Timestamp)
	return difference.Seconds() <= maxAge
}

// NewToken Return a new random token with timestamp set to the current UTC time
func NewToken(length int) Token {
	return Token{
		Data:      GenerateToken(length),
		Timestamp: time.Now().UTC(),
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
		AUTH VARCHAR(64),
		MASTER VARCHAR(260),
		AUTH_TIMESTAMP DATETIME
	)
	`
}

// GetTokenBundle Retrives tokens from the token-db
func GetTokenBundle(path string) TokenBundle {
	db := ConnectTokenDB(path)
	defer db.Close()
	var tokens TokenBundle
	err := db.QueryRow("SELECT AUTH, MASTER, AUTH_TIMESTAMP FROM TOKEN").Scan(
		&tokens.Auth.Data, &tokens.Master, &tokens.Auth.Timestamp)
	util.CheckErrFatal(err)
	return tokens
}

// ConnectTokenDB Connects to a token db and returns the handler
func ConnectTokenDB(path string) *sql.DB {
	dbFile := filepath.Join(path, "token-db")
	db, err := sql.Open("sqlite3", dbFile)
	util.CheckErrFatal(err)
	err = db.Ping()
	util.CheckErrFatal(err)
	return db
}
