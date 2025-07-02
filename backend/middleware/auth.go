package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"procurement/database"
	"procurement/handlers"
	"procurement/models"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Validator is a struct that holds the necessary information for JWT validation.
type Validator struct {
	Audience string
	Issuer   string
	jwks     *jose.JSONWebKeySet
}

// NewValidator creates a new Validator and fetches the JWKS from the Auth0 domain.
func NewValidator() (*Validator, error) {
	issuer := os.Getenv("AUTH0_DOMAIN")
	audience := os.Getenv("AUTH0_AUDIENCE")

	jwksURL := "https://" + issuer + "/.well-known/jwks.json"
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &Validator{
		Audience: audience,
		Issuer:   issuer,
		jwks:     &jwks,
	}, nil
}

// TokenMiddleware verifies the JWT token and adds the user's internal ID to the request context.
func TokenMiddleware(next http.Handler) http.Handler {
	validator, err := NewValidator()
	if err != nil {
		log.Fatalf("Failed to create validator: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			return
		}
		tokenString := parts[1]

		token, err := jwt.ParseSigned(tokenString)
		if err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims := jwt.Claims{}
		if err := token.Claims(validator.jwks, &claims); err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		if err := claims.Validate(jwt.Expected{
			Audience: jwt.Audience{validator.Audience},
			Issuer:   validator.Issuer,
			Time:     time.Now(),
		}); err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Token validation failed")
			return
		}

		db := database.GetDB()
		var user models.User
		if err := db.Where("auth0_id = ?", claims.Subject).First(&user).Error; err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
