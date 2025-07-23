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

	// Find user by email, as it's the stable unique identifier
	user := models.User{
		Email: payload.Email, // Condition for FirstOrCreate
	}

	// Data to assign if creating or updating.
	assignData := models.User{
		Username:   payload.Name,
		PictureURL: payload.Picture,
		// The Auth0ID might be empty now, so we ensure it's handled correctly.
		// If payload.Auth0ID is not empty, it will be set.
	}
	if payload.Auth0ID != "" {
		assignData.Auth0ID = &payload.Auth0ID
	}

	// FirstOrCreate will find the user by Email or create a new one if not found.
	// Assign will update the fields in assignData for both found and new records.
	result := dbConn.Where(models.User{Email: payload.Email}).Assign(assignData).FirstOrCreate(&user)

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
