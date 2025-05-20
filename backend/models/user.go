package models

import "time"

// User represents a user record in the database, matching the Users table schema.
// Note: time.Time could be used for CreatedAt/UpdatedAt if you add those columns.
type User struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	Auth0ID       string    `json:"auth0_id" gorm:"column:auth0_id;uniqueIndex"`         // Corresponds to auth0_id TEXT UNIQUE NOT NULL
	Username      string    `json:"username" gorm:"column:username;uniqueIndex"`             // Corresponds to username TEXT UNIQUE NOT NULL
	Email         string    `json:"email" gorm:"column:email;uniqueIndex"`                   // Corresponds to email TEXT UNIQUE NOT NULL
	PictureURL    string    `json:"picture_url,omitempty" gorm:"column:picture_url"`       // Corresponds to picture_url TEXT
	Role          string    `json:"role" gorm:"column:role;default:'requester'"`             // Corresponds to role TEXT NOT NULL DEFAULT 'requester'
	Department    *string   `json:"department,omitempty" gorm:"column:department"`        // Corresponds to department TEXT
	ContactNumber *string   `json:"contact_number,omitempty" gorm:"column:contactNumber"`  // Corresponds to contactNumber TEXT
	IsActive      bool      `json:"is_active" gorm:"column:isActive;default:true"`         // Corresponds to isActive INTEGER DEFAULT 1 (GORM maps bool to numeric for SQLite)
	CreatedAt     time.Time `json:"created_at"`                                           // GORM will manage this
	UpdatedAt     time.Time `json:"updated_at"`                                           // GORM will manage this
}

// UserSyncPayload is the expected structure of the JSON body from the frontend
// when syncing user data from Auth0.
type UserSyncPayload struct {
	Auth0ID string `json:"auth0_id"` // The 'sub' claim from Auth0
	Name    string `json:"name"`     // From Auth0 'name' or 'nickname', will be used for 'username'
	Email   string `json:"email"`    // From Auth0 'email'
	Picture string `json:"picture"`  // From Auth0 'picture', will be used for 'picture_url'
}
