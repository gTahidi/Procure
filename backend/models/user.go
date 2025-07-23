package models

import "time"

// User represents a user record in the database, matching the Users table schema.
type User struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	Auth0ID       *string   `json:"auth0_id,omitempty" gorm:"column:auth0_id;uniqueIndex"` // Now nullable
	Username      string    `json:"username" gorm:"column:username;uniqueIndex"`
	Email         string    `json:"email" gorm:"column:email;uniqueIndex"`
	PasswordHash  string    `json:"-" gorm:"column:password_hash"`                        // Password hash not exposed in JSON
	PictureURL    string    `json:"picture_url,omitempty" gorm:"column:picture_url"`
	Role          string    `json:"role" gorm:"column:role;default:'requester'"`
	Department    *string   `json:"department,omitempty" gorm:"column:department"`
	ContactNumber *string   `json:"contact_number,omitempty" gorm:"column:contactNumber"`
	IsActive      bool      `json:"is_active" gorm:"column:isActive;default:true"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserSyncPayload is the expected structure of the JSON body from the frontend
// when syncing user data from Auth0.
type UserSyncPayload struct {
	Auth0ID string `json:"auth0_id"` // The 'sub' claim from Auth0
	Name    string `json:"name"`     // From Auth0 'name' or 'nickname', will be used for 'username'
	Email   string `json:"email"`    // From Auth0 'email'
	Picture string `json:"picture"`  // From Auth0 'picture', will be used for 'picture_url'
}

// RegisterRequest represents the data needed to register a new user
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest represents the data needed to log in a user
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// PasswordChangeRequest represents the data needed to change a user's password
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// PasswordResetRequest represents the data needed to request a password reset
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest represents the data needed to confirm a password reset
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse represents the response sent after successful authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
