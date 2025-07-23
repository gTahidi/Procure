package middleware

import (
	"context"
	"errors"
	"net/http"
	"procurement/database"
	"procurement/models"
	"procurement/services"
	"strings"

	"gorm.io/gorm"
)

// AuthMiddleware handles JWT authentication and authorization
type AuthMiddleware struct {
	TokenService services.TokenService
	DB           *gorm.DB
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware() (*AuthMiddleware, error) {
	tokenService, err := services.NewJWTTokenService()
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		TokenService: tokenService,
		DB:           database.GetDB(),
	}, nil
}

// Authenticate validates the JWT token and adds the user ID to the request context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// Validate token
		claims, err := m.TokenService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if token is blacklisted/invalidated
		var session models.Session
		result := m.DB.Where("token = ? AND is_valid = ?", tokenString, false).First(&session)
		if result.Error == nil {
			// Token is blacklisted
			http.Error(w, "Token has been invalidated", http.StatusUnauthorized)
			return
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Database error
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Check if user exists and is active
		var user models.User
		if result := m.DB.First(&user, claims.UserID); result.Error != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if !user.IsActive {
			http.Error(w, "User account is inactive", http.StatusForbidden)
			return
		}

		// Add user ID and role to context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks if the user has the required role
func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context
			role, ok := r.Context().Value("userRole").(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has one of the required roles
			hasRole := false
			for _, requiredRole := range roles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}