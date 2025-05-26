package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"procurement/database"
	"procurement/handlers" // For RespondWithError
	"procurement/models"

	"github.com/golang-jwt/jwt/v5"
)

// TokenMiddleware is a simplified authentication middleware for POC/development purposes.
// It parses the JWT token to extract the 'sub' claim (Auth0ID) WITHOUT full signature verification.
// It then looks up the user in the database and adds their internal ID to the request context.
// WARNING: Not for production use due to lack of signature validation.
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("ERROR: TokenMiddleware: Authorization header missing")
			handlers.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("ERROR: TokenMiddleware: Authorization header format must be Bearer {token}")
			handlers.RespondWithError(w, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			return
		}
		tokenString := parts[1]

		// Parse the token without verifying the signature (POC/development only!)
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			log.Printf("ERROR: TokenMiddleware: Error parsing token: %v\n", err)
			handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("ERROR: TokenMiddleware: Could not assert claims to jwt.MapClaims")
			handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			log.Println("ERROR: TokenMiddleware: 'sub' claim missing or invalid in token")
			handlers.RespondWithError(w, http.StatusUnauthorized, "Token missing user identifier")
			return
		}

		db := database.GetDB()
		var user models.User
		if err := db.Where("auth0_id = ?", sub).First(&user).Error; err != nil {
			log.Printf("ERROR: TokenMiddleware: User not found for auth0_id %s: %v\n", sub, err)
			handlers.RespondWithError(w, http.StatusUnauthorized, "User not found or database error")
			return
		}

		// Add user ID to context. user.ID is int64.
		// ListRequisitionsHandler expects userID in context. It currently asserts to uint.
		// We will need to adjust ListRequisitionsHandler to expect int64 or convert here if strictly necessary.
		// For now, we pass user.ID as is (int64).
		ctx := context.WithValue(r.Context(), "userID", user.ID)
		log.Printf("INFO: TokenMiddleware: Authenticated user %s (DB ID: %d), proceeding to handler.\n", sub, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
