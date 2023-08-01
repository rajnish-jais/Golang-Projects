package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth_GenerateToken(t *testing.T) {
	// Create a new Auth instance with a secret key
	auth := NewAuth("my-secret-key")

	// Generate a token for a test username
	token, err := auth.GenerateToken("testuser")
	assert.NoError(t, err, "GenerateToken should not return an error")
	assert.NotEmpty(t, token, "Generated token should not be empty")
}

func TestAuth_VerifyToken(t *testing.T) {
	// Create a new Auth instance with a secret key
	auth := NewAuth("my-secret-key")

	// Generate a token for a test username
	token, _ := auth.GenerateToken("testuser")

	// Verify the generated token
	username, err := auth.VerifyToken(token)
	assert.NoError(t, err, "VerifyToken should not return an error")
	assert.Equal(t, "testuser", username, "Verified username should match the original username")

	// Test invalid token
	invalidToken := "invalid-token"
	_, err = auth.VerifyToken(invalidToken)
	assert.Error(t, err, "VerifyToken should return an error for invalid token")

	// Test token expiration
	//expiredToken, _ := auth.GenerateToken("testuser")
	//time.Sleep(2 * time.Second) // Sleep to make the token expire
	//_, err = auth.VerifyToken(expiredToken)
	//assert.Error(t, err, "VerifyToken should return an error for expired token")
}
