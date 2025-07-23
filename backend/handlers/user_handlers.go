package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"procurement/database"
	"procurement/models"

	"github.com/go-chi/chi/v5"
)

// SyncUserHandler handles requests to /api/users/sync
func SyncUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.UserSyncPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("ERROR: SyncUserHandler: Failed to decode payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("INFO: SyncUserHandler: Received payload: %+v", payload)

	dbConn := database.GetDB()
	if dbConn == nil {
		log.Println("ERROR: SyncUserHandler: Database not initialized")
		http.Error(w, "Database connection not initialized", http.StatusInternalServerError)
		return
	}

	// Convert string to *string for Auth0ID
	auth0ID := payload.Auth0ID
	user := models.User{
		Auth0ID: &auth0ID, // Condition for FirstOrCreate
	}

	// Data to assign if creating or updating. GORM will only update these fields.
	assignData := models.User{
		Username:   payload.Name, // Assuming payload.Name maps to models.User.Username
		Email:      payload.Email,
		PictureURL: payload.Picture,
		// Role is set by default in the model or can be managed separately
		// Department and ContactNumber are pointers and will be nil if not set here
		// IsActive is true by default in the model
	}

	// FirstOrCreate will find the user by Auth0ID or create a new one if not found.
	// Assign will update the fields specified in assignData for both found and new records.
	// If you only want to update on create, use .Attrs() instead of .Assign() for creation-only fields.
	auth0IDForWhere := payload.Auth0ID
	result := dbConn.Where(models.User{Auth0ID: &auth0IDForWhere}).Assign(assignData).FirstOrCreate(&user)

	if result.Error != nil {
		log.Printf("ERROR: SyncUserHandler: Database error during FirstOrCreate/Assign: %v", result.Error)
		http.Error(w, "Failed to sync user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected > 0 {
		if user.ID == 0 { // Should not happen if FirstOrCreate works as expected and ID is auto-incrementing
			log.Printf("INFO: SyncUserHandler: User with Auth0ID %s might have been created but ID is 0. Re-fetching.", payload.Auth0ID)
			// Re-fetch to be sure, though GORM should populate 'user' correctly
			dbConn.Where("auth0_id = ?", payload.Auth0ID).First(&user)
		}
		log.Printf("INFO: SyncUserHandler: User with Auth0ID %s (ID: %d) synced successfully. Rows affected: %d", payload.Auth0ID, user.ID, result.RowsAffected)
	} else {
		log.Printf("INFO: SyncUserHandler: User with Auth0ID %s (ID: %d) already up-to-date. Rows affected: %d", payload.Auth0ID, user.ID, result.RowsAffected)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Or http.StatusOK if it was an update
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("ERROR: SyncUserHandler: Failed to encode response: %v", err)
		// Header already sent, can't send http.Error
	}
}

func RegisterUserRoutes(router chi.Router) {
	// Assuming db is accessible, perhaps passed in or a global/package variable
	// For simplicity, if db is a global var in main package, it might need to be passed or accessed carefully.
	// Here, we'll assume SyncUserHandler can access the db connection it needs.
	router.Post("/users/sync", SyncUserHandler)

	// Example: router.Get("/users/{id}", GetUserHandler)
	// Example: router.Post("/users", CreateUserHandler)
}
