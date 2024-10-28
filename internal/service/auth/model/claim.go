package model

import "github.com/dgrijalva/jwt-go"

const (
	// ExamplePath ???
	ExamplePath = "/note_v1.NoteV1/Get"
)

// UserClaims ???
type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}
