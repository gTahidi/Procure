package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"procurement/models"
)

// TokenClaims represents the claims in the JWT token
type TokenClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenService defines the interface for token-related operations
type TokenService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	GeneratePasswordResetToken(userID int64) (string, error)
	ValidatePasswordResetToken(tokenString string) (int64, error)
}

// JWTTokenService implements TokenService using JWT
type JWTTokenService struct {
	secretKey        []byte
	accessTokenTTL   time.Duration
	resetTokenTTL    time.Duration
	refreshTokenTTL  time.Duration
	issuer           string
}

// NewJWTTokenService creates a new JWTTokenService
func NewJWTTokenService() (*JWTTokenService, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY environment variable is not set")
	}

	return &JWTTokenService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  24 * time.Hour,    // 24 hours for access tokens
		resetTokenTTL:   1 * time.Hour,     // 1 hour for password reset tokens
		refreshTokenTTL: 7 * 24 * time.Hour, // 7 days for refresh tokens
		issuer:          "procurement-app",
	}, nil
}

// GenerateToken generates a new JWT token for the given user
func (s *JWTTokenService) GenerateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates the JWT token and returns the claims
func (s *JWTTokenService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GeneratePasswordResetToken generates a token for password reset
func (s *JWTTokenService) GeneratePasswordResetToken(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(s.resetTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    s.issuer,
		Subject:   fmt.Sprintf("%d", userID),
		ID:        fmt.Sprintf("reset-%d-%d", userID, now.Unix()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidatePasswordResetToken validates a password reset token and returns the user ID
func (s *JWTTokenService) ValidatePasswordResetToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Parse the user ID from the subject claim
		var userID int64
		_, err := fmt.Sscanf(claims.Subject, "%d", &userID)
		if err != nil {
			return 0, errors.New("invalid subject claim in token")
		}
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// GenerateRefreshToken generates a refresh token for the given user
func (s *JWTTokenService) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    s.issuer,
		Subject:   fmt.Sprintf("%d", userID),
		ID:        fmt.Sprintf("refresh-%d-%d", userID, now.Unix()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}