package models

import "time"

// PasswordReset represents a password reset token record in the database
type PasswordReset struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"column:user_id;index"`
	Token     string    `json:"token" gorm:"column:token;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}