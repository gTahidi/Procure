package middleware

import (
	"errors"
	"net/http"
)

// Validator is a struct that holds the necessary information for JWT validation.
type Validator struct {
	Audience string
	Issuer   string
}

// NewValidator creates a new Validator for JWT validation.
// This is a placeholder for the Auth0 validator that we're replacing.
func NewValidator(domain, audience string) (*Validator, error) {
	if domain == "" || audience == "" {
		return nil, errors.New("domain and audience are required")
	}

	return &Validator{
		Audience: audience,
		Issuer:   "https://" + domain + "/",
	}, nil
}

// TokenMiddleware is a placeholder for the Auth0 token middleware.
// We're keeping this function signature for backward compatibility,
// but we'll use our own authentication middleware instead.
func TokenMiddleware(validator *Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This function should not be called anymore
			// We're using our own authentication middleware instead
			http.Error(w, "Auth0 authentication is no longer supported", http.StatusUnauthorized)
		})
	}
}