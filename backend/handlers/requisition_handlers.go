package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"procurement/database" // Module name 'procurement' then path
	"procurement/models"
	// "time" // No longer explicitly needed for CreatedAt if GORM handles it
	"gorm.io/gorm" // Added for gorm.ErrRecordNotFound or other GORM specific needs
)

// CreateRequisitionHandler handles POST requests to create a new requisition
func CreateRequisitionHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB() // This returns *gorm.DB
	if db == nil {
		log.Println("ERROR: CreateRequisitionHandler: Database not initialized")
		http.Error(w, "Database connection not initialized", http.StatusInternalServerError)
		return
	}

	var reqPayload models.Requisition
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Invalid request payload: %v\n", err)
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Basic Validation (more comprehensive validation should be added)
	if reqPayload.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if reqPayload.Type == "" {
		http.Error(w, "Requisition type is required", http.StatusBadRequest)
		return
	}
	if len(reqPayload.Items) == 0 {
		http.Error(w, "At least one requisition item is required", http.StatusBadRequest)
		return
	}
	for _, item := range reqPayload.Items {
		if item.Description == "" || item.Quantity <= 0 || item.Unit == "" {
			http.Error(w, "All requisition items must have a description, quantity, and unit", http.StatusBadRequest)
			return
		}
	}

	// Start a GORM transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to start GORM transaction: %v\n", tx.Error)
		http.Error(w, "Failed to start transaction: "+tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Defer a rollback in case of panic or error.
	// If tx.Commit() is called successfully, the rollback is a no-op.
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			// Re-panic if you want to propagate the panic after rollback
			// panic(rec)
			log.Printf("PANIC: CreateRequisitionHandler: Rolled back transaction due to panic: %v", rec)
			// Ensure a response is sent if a panic occurs during HTTP handling
			if הודעה, בסדר := rec.(string); בסדר {
				http.Error(w, "Internal server error after panic: "+הודעה, http.StatusInternalServerError)
			} else {
				http.Error(w, "Internal server error after panic", http.StatusInternalServerError)
			}
		} else if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound { // Check tx.Error from operations like Create, Save, Delete
			// If tx.Error is set from an operation like tx.Create, and we haven't committed yet,
			// then we should roll back.
			// gorm.ErrRecordNotFound might be a legitimate case in some non-creation scenarios, 
			// but for creation, any error usually means rollback.
			log.Printf("INFO: CreateRequisitionHandler: Transaction error, rolling back: %v", tx.Error)
			tx.Rollback() // tx.Error should have been set by a failing GORM operation
		}
	}()

	// Set default status if not provided
	if reqPayload.Status == "" {
		reqPayload.Status = "pending" // Default status
	}
	// GORM will handle CreatedAt and UpdatedAt automatically if fields exist in the model (e.g. via gorm.Model embedding or explicit fields)

	// Insert Requisition using GORM
	if err := tx.Create(&reqPayload).Error; err != nil {
		// Rollback is handled by defer, but log and return error immediately
		log.Printf("ERROR: CreateRequisitionHandler: Failed to insert requisition: %v\n", err)
		http.Error(w, "Failed to insert requisition: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// reqPayload.ID is now populated by GORM

	// Insert RequisitionItems using GORM
	for i := range reqPayload.Items {
		item := &reqPayload.Items[i] // Get a pointer to the item in the slice
		item.RequisitionID = reqPayload.ID // Set the foreign key

		if err := tx.Create(item).Error; err != nil {
			// Rollback is handled by defer, but log and return error immediately
			log.Printf("ERROR: CreateRequisitionHandler: Failed to insert requisition item (%s): %v\n", item.Description, err)
			http.Error(w, "Failed to insert requisition item: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// item.ID is now populated by GORM
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		// Rollback is handled by defer if commit fails and sets tx.Error, but explicit check is good.
		log.Printf("ERROR: CreateRequisitionHandler: Failed to commit transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reqPayload); err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to encode response: %v\n", err)
		// The transaction is committed, but we couldn't send the response.
		// http.Error might not work here if headers are already sent.
	}
	log.Printf("INFO: Requisition created successfully with ID: %d", reqPayload.ID)
}
