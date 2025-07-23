package models

import "time"

// Session represents a user session record in the database
type Session struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"column:user_id;index"`
	Token     string    `json:"token" gorm:"column:token;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	IPAddress string    `json:"ip_address,omitempty" gorm:"column:ip_address"`
	UserAgent string    `json:"user_agent,omitempty" gorm:"column:user_agent"`
	IsValid   bool      `json:"is_valid" gorm:"column:is_valid;default:true"`
}