package database

import (
	"database/sql"
	"log"
)

// RunAuthMigrations executes the database migrations needed for email/password authentication
func RunAuthMigrations(db *sql.DB) error {
	log.Println("Running authentication migrations...")

	// Add password_hash column to Users table if it doesn't exist
	_, err := db.Exec(`
		ALTER TABLE Users ADD COLUMN password_hash TEXT;
	`)
	if err != nil {
		log.Printf("Note: Error adding password_hash column (may already exist): %v", err)
		// Continue execution - SQLite returns error if column already exists
	}

	// Make auth0_id nullable since we're moving away from Auth0
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Users_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			auth0_id TEXT UNIQUE,             -- Changed to nullable
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT,               -- Password hash field
			picture_url TEXT,
			role TEXT NOT NULL DEFAULT 'requester' CHECK(role IN ('admin', 'procurement_officer', 'requester', 'supplier', 'approver', 'evaluator')),
			department TEXT, 
			contactNumber TEXT, 
			isActive INTEGER DEFAULT 1 CHECK(isActive IN (0,1)),
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Printf("Error creating Users_new table: %v", err)
		return err
	}

	// Only attempt migration if the original Users table exists
	var tableExists int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='Users';").Scan(&tableExists)
	if err != nil {
		log.Printf("Error checking if Users table exists: %v", err)
		return err
	}

	if tableExists > 0 {
		// Copy data from old table to new table
		_, err = db.Exec(`
			INSERT INTO Users_new (id, auth0_id, username, email, picture_url, role, department, contactNumber, isActive)
			SELECT id, auth0_id, username, email, picture_url, role, department, contactNumber, isActive FROM Users;
		`)
		if err != nil {
			log.Printf("Error copying data to Users_new: %v", err)
			return err
		}

		// Drop old table
		_, err = db.Exec(`DROP TABLE Users;`)
		if err != nil {
			log.Printf("Error dropping old Users table: %v", err)
			return err
		}

		// Rename new table to Users
		_, err = db.Exec(`ALTER TABLE Users_new RENAME TO Users;`)
		if err != nil {
			log.Printf("Error renaming Users_new to Users: %v", err)
			return err
		}
	} else {
		// If Users table doesn't exist, rename the new table
		_, err = db.Exec(`ALTER TABLE Users_new RENAME TO Users;`)
		if err != nil {
			log.Printf("Error renaming Users_new to Users: %v", err)
			return err
		}
	}

	// Create password_resets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS PasswordResets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at TEXT NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Printf("Error creating PasswordResets table: %v", err)
		return err
	}

	// Create index on token for faster lookups
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_password_resets_token ON PasswordResets(token);
	`)
	if err != nil {
		log.Printf("Error creating index on PasswordResets: %v", err)
		return err
	}

	// Create index on user_id for faster lookups
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_password_resets_user_id ON PasswordResets(user_id);
	`)
	if err != nil {
		log.Printf("Error creating index on PasswordResets: %v", err)
		return err
	}

	// Create sessions table for managing user sessions
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at TEXT NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			ip_address TEXT,
			user_agent TEXT,
			is_valid INTEGER DEFAULT 1 CHECK(is_valid IN (0,1)),
			FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Printf("Error creating Sessions table: %v", err)
		return err
	}

	// Create index on session token for faster lookups
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_sessions_token ON Sessions(token);
	`)
	if err != nil {
		log.Printf("Error creating index on Sessions: %v", err)
		return err
	}

	// Create index on user_id for faster lookups
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON Sessions(user_id);
	`)
	if err != nil {
		log.Printf("Error creating index on Sessions: %v", err)
		return err
	}

	log.Println("Authentication migrations completed successfully")
	return nil
}