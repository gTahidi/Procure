package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"procurement/database" // Module name 'procurement' then path
	"procurement/models"
	"time"
)

// CreateRequisitionHandler handles POST requests to create a new requisition
func CreateRequisitionHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
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


	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to start transaction: %v\n", err)
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into Requisitions table
	if reqPayload.Status == "" {
		reqPayload.Status = "pending" // Default status
	}
	// CreatedAt will be set by the database default or here
	// For consistency, let's set it here, though db_init.go also defines a default.
	// If we rely on DB default, we'd need to query it back.
	currentTime := time.Now()


	result, err := tx.Exec(`
		INSERT INTO Requisitions (user_id, type, aac, materialGroup, exchangeRate, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, reqPayload.UserID, reqPayload.Type, reqPayload.AAC, reqPayload.MaterialGroup, reqPayload.ExchangeRate, reqPayload.Status, currentTime)

	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: CreateRequisitionHandler: Failed to insert requisition: %v\n", err)
		http.Error(w, "Failed to insert requisition: "+err.Error(), http.StatusInternalServerError)
		return
	}

	requisitionID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: CreateRequisitionHandler: Failed to get last insert ID for requisition: %v\n", err)
		http.Error(w, "Failed to get last insert ID for requisition: "+err.Error(), http.StatusInternalServerError)
		return
	}
	reqPayload.ID = requisitionID
	reqPayload.CreatedAt = currentTime // Ensure response has the correct time

	// Insert into RequisitionItems table
	stmt, err := tx.Prepare(`
		INSERT INTO RequisitionItems (requisition_id, description, quantity, unit, estimated_unit_price, freight_cost, insurance_cost, installation_cost, amr_id, value)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: CreateRequisitionHandler: Failed to prepare requisition item statement: %v\n", err)
		http.Error(w, "Failed to prepare requisition item statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for i := range reqPayload.Items {
		item := &reqPayload.Items[i]
		item.RequisitionID = requisitionID 
		
		itemResult, itemErr := stmt.Exec(
			item.RequisitionID, item.Description, item.Quantity, item.Unit,
			item.EstimatedUnitPrice, item.FreightCost, item.InsuranceCost,
			item.InstallationCost, item.AmrID, item.Value,
		)
		if itemErr != nil {
			tx.Rollback()
			log.Printf("ERROR: CreateRequisitionHandler: Failed to insert requisition item: %v\n", itemErr)
			http.Error(w, "Failed to insert requisition item: "+itemErr.Error(), http.StatusInternalServerError)
			return
		}
		// Get the ID of the newly inserted item and update the struct
		itemId, idErr := itemResult.LastInsertId()
		if idErr != nil {
			tx.Rollback()
			log.Printf("ERROR: CreateRequisitionHandler: Failed to get last insert ID for requisition item: %v\n", idErr)
			http.Error(w, "Failed to get last insert ID for requisition item: "+idErr.Error(), http.StatusInternalServerError)
			return
		}
		item.ID = itemId
	}

	if err := tx.Commit(); err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to commit transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqPayload) 
	log.Printf("Requisition created successfully with ID: %d", reqPayload.ID)
}
