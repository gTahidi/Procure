package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"procurement/database"
	"procurement/models"
	"procurement/services"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// AuthController handles authentication-related requests
type AuthController struct {
	DB             *gorm.DB
	TokenService   services.TokenService
	PasswordService services.PasswordService
	EmailService   services.EmailService
}

// NewAuthController creates a new AuthController
func NewAuthController() (*AuthController, error) {
	db := database.GetDB()
	
	tokenService, err := services.NewJWTTokenService()
	if err != nil {
		return nil, err
	}
	
	passwordService := services.NewBCryptPasswordService()
	
	emailService, err := services.GetEmailService()
	if err != nil {
		return nil, err
	}
	
	return &AuthController{
		DB:             db,
		TokenService:   tokenService,
		PasswordService: passwordService,
		EmailService:   emailService,
	}, nil
}

// RegisterRoutes registers the authentication routes
func (c *AuthController) RegisterRoutes(r chi.Router) {
	r.Post("/auth/register", c.Register)
	r.Post("/auth/login", c.Login)
	r.Post("/auth/password/change", c.ChangePassword)
	r.Post("/auth/password/reset/request", c.RequestPasswordReset)
	r.Post("/auth/password/reset/confirm", c.ResetPassword)
	r.Post("/auth/logout", c.Logout)
}

// Register handles user registration
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate email format
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	
	// Validate password
	if err := c.PasswordService.ValidatePasswordStrength(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Check if email already exists
	var existingUser models.User
	result := c.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	
	// Hash password
	hashedPassword, err := c.PasswordService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	
	// Create user
	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         "requester", // Default role
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	if result := c.DB.Create(&user); result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	
	// Generate token
	token, err := c.TokenService.GenerateToken(&user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	// Return response
	response := models.AuthResponse{
		Token: token,
		User:  user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find user by email
	var user models.User
	if result := c.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	
	// Check if user is active
	if !user.IsActive {
		http.Error(w, "Account is inactive", http.StatusForbidden)
		return
	}
	
	// Verify password
	if err := c.PasswordService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	
	// Generate token
	token, err := c.TokenService.GenerateToken(&user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	// Return response
	response := models.AuthResponse{
		Token: token,
		User:  user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ChangePassword handles password change requests
func (c *AuthController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordChangeRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Find user
	var user models.User
	if result := c.DB.First(&user, userID); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Verify current password
	if err := c.PasswordService.VerifyPassword(user.PasswordHash, req.CurrentPassword); err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}
	
	// Validate new password
	if err := c.PasswordService.ValidatePasswordStrength(req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Hash new password
	hashedPassword, err := c.PasswordService.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	
	// Update password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()
	if result := c.DB.Save(&user); result.Error != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	
	// Invalidate existing sessions (optional)
	if err := c.invalidateUserSessions(userID); err != nil {
		// Log error but don't fail the request
		// log.Printf("Failed to invalidate sessions: %v", err)
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password changed successfully"}`))
}

// Helper method to invalidate user sessions
func (c *AuthController) invalidateUserSessions(userID int64) error {
	result := c.DB.Model(&models.Session{}).
		Where("user_id = ? AND is_valid = ?", userID, true).
		Updates(map[string]interface{}{"is_valid": false})
	return result.Error
}

// RequestPasswordReset handles password reset requests
func (c *AuthController) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordResetRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find user by email
	var user models.User
	if result := c.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		// Don't reveal if email exists or not for security reasons
		// Just return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"If your email is registered, you will receive a password reset link"}`))
		return
	}
	
	// Generate reset token
	token, err := c.TokenService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}
	
	// Store reset token in database
	passwordReset := models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}
	
	if result := c.DB.Create(&passwordReset); result.Error != nil {
		http.Error(w, "Failed to store reset token", http.StatusInternalServerError)
		return
	}
	
	// Create reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", c.getBaseURL(r), token)
	
	// Send reset email
	if err := c.EmailService.SendPasswordResetEmail(user.Email, resetLink); err != nil {
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password reset link sent to your email"}`))
}

// Helper method to get base URL
func (c *AuthController) getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

// ResetPassword handles password reset confirmation
func (c *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordResetConfirmRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find reset token in database
	var passwordReset models.PasswordReset
	if result := c.DB.Where("token = ?", req.Token).First(&passwordReset); result.Error != nil {
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}
	
	// Check if token is expired
	if time.Now().After(passwordReset.ExpiresAt) {
		http.Error(w, "Reset token has expired", http.StatusBadRequest)
		return
	}
	
	// Validate new password
	if err := c.PasswordService.ValidatePasswordStrength(req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Find user
	var user models.User
	if result := c.DB.First(&user, passwordReset.UserID); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Hash new password
	hashedPassword, err := c.PasswordService.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	
	// Update password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()
	if result := c.DB.Save(&user); result.Error != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	
	// Invalidate reset token
	if result := c.DB.Delete(&passwordReset); result.Error != nil {
		// Log error but don't fail the request
		// log.Printf("Failed to delete reset token: %v", result.Error)
	}
	
	// Invalidate existing sessions
	if err := c.invalidateUserSessions(user.ID); err != nil {
		// Log error but don't fail the request
		// log.Printf("Failed to invalidate sessions: %v", err)
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password has been reset successfully"}`))
}

// Logout handles user logout
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
		return
	}
	
	// Extract token from "Bearer {token}"
	tokenString := authHeader[7:]
	
	// Invalidate the current session
	session := models.Session{
		UserID:  userID,
		Token:   tokenString,
		IsValid: false,
	}
	
	// Update if exists, create if not
	result := c.DB.Where("token = ?", tokenString).FirstOrCreate(&session)
	if result.Error != nil {
		http.Error(w, "Failed to invalidate session", http.StatusInternalServerError)
		return
	}
	
	// If session existed, mark it as invalid
	if result.RowsAffected > 0 {
		c.DB.Model(&session).Update("is_valid", false)
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}