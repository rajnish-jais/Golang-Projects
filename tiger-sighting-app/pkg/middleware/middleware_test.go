package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiger-sighting-app/pkg/auth"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a new Auth instance with a secret key
	authInstance := auth.NewAuth("my-secret-key")

	// Create a test handler that sets a response header with the username from the context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username, ok := ctx.Value("username").(string)
		if !ok {
			username = "unknown"
		}
		w.Header().Set("X-Username", username)
	})

	// Create a new request with a valid token in the Authorization header
	token, _ := authInstance.GenerateToken("testuser")
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create the auth middleware and wrap it around the testHandler
	authMiddleware := AuthMiddleware(authInstance, testHandler)

	// Call the authMiddleware with the test request
	authMiddleware.ServeHTTP(rr, req)

	// Check the response status code and headers
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")
	assert.Equal(t, "testuser", rr.Header().Get("X-Username"), "Unexpected username in response header")

	// Test unauthorized request (no Authorization header)
	reqWithoutAuth := httptest.NewRequest("GET", "/test", nil)
	rrWithoutAuth := httptest.NewRecorder()

	authMiddleware.ServeHTTP(rrWithoutAuth, reqWithoutAuth)
	assert.Equal(t, http.StatusUnauthorized, rrWithoutAuth.Code, "Unexpected status code for unauthorized request")
}
