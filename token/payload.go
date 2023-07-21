package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("invalid token signing method")

// Payload is the output of the token creation process
type Payload struct {
	ID uuid.UUID `json:"id"` // Unique ID for the token
	Username string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload {
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}

	if payload.Username == "" {
		return ErrInvalidToken
	}

	return nil
}
