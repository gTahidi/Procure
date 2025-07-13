package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"procurement/database"
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
func NewValidator(domain, audience string) (*Validator, error) {
	// The issuer in the JWT token is the full URL of the Auth0 domain.
	issuerURL := "https://" + domain + "/"

	// The JWKS endpoint is located at a specific path on that domain.
	jwksURL := issuerURL + ".well-known/jwks.json"

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
		Issuer:   issuerURL, // Use the full URL for validation.
		jwks:     &jwks,
	}, nil
}

// TokenMiddleware verifies the JWT token and adds the user's internal ID to the request context.
func TokenMiddleware(validator *Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Auth Error: Authorization header is required")
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Printf("Auth Error: Invalid header format: %s", authHeader)
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 2. Parse and validate the token
			token, err := jwt.ParseSigned(tokenString)
			if err != nil {
				log.Printf("Auth Error: Failed to parse token: %v", err)
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			// 3. Get claims and verify signature with JWKS
			claims := &jwt.Claims{}
			if err := token.Claims(validator.jwks, claims); err != nil {
				log.Printf("Auth Error: Failed to validate token signature with JWKS: %v", err)
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}

			// 4. Validate the claims (audience, issuer, expiry)
			log.Printf("Auth Check: Token Claims: Issuer=[%s], Audience=%s", claims.Issuer, claims.Audience)
			log.Printf("Auth Check: Validator Config: Issuer=[%s], Audience=[%s]", validator.Issuer, validator.Audience)
			expected := jwt.Expected{
				Audience: claims.Audience,
				Issuer:   validator.Issuer,
				Time:     time.Now(),
			}
			if err := claims.Validate(expected); err != nil {
				log.Printf("Auth Error: Token claims validation failed: %v", err)
				http.Error(w, "Token claims validation failed", http.StatusUnauthorized)
				return
			}

			log.Println("Auth Success: Token is valid")

			// 5. Find user in DB and add to context
			db := database.GetDB()
			var user models.User
			if err := db.Where("auth0_id = ?", claims.Subject).First(&user).Error; err != nil {
				log.Printf("Auth Error: User with Auth0ID '%s' not found in DB: %v", claims.Subject, err)
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
