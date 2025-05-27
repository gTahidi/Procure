package handlers

import (
	"encoding/json"
	"errors" // Added for errors.Is
	"fmt"
	"log"
	"net/http"
	"procurement/database" // Module name 'procurement' then path
	"procurement/models"
	// "time" // No longer explicitly needed for CreatedAt if GORM handles it
	"gorm.io/gorm" // Added for gorm.ErrRecordNotFound or other GORM specific needs
	"strconv"      // Added for strconv.ParseInt
	"strings"      // Added for strings.EqualFold

	"github.com/go-chi/chi/v5" // Added for chi.URLParam
)

// CreateRequisitionHandler handles POST requests to create a new requisition
func CreateRequisitionHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB() // This returns *gorm.DB
	var reqPayload models.Requisition

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to decode request body: %v\n", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Basic Validation (example)
	if reqPayload.UserID == 0 {
		RespondWithError(w, http.StatusBadRequest, "UserID is required")
		return
	}
	if reqPayload.Type == "" {
		RespondWithError(w, http.StatusBadRequest, "Requisition type is required")
		return
	}
	if len(reqPayload.Items) == 0 {
		RespondWithError(w, http.StatusBadRequest, "At least one item is required")
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to begin transaction: %v\n", tx.Error)
		RespondWithError(w, http.StatusInternalServerError, "Failed to begin transaction: "+tx.Error.Error())
		return
	}

	// Defer a rollback in case of panic or an unhandled error path.
	// Successful Commit will make this Rollback a no-op.
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("PANIC: CreateRequisitionHandler: Recovered from panic: %v\n", r)
			// RespondWithError(w, http.StatusInternalServerError, "Internal server error after panic") // Might not be possible if headers sent
			return
		}
		// If tx.Error is set (e.g. by a failed Commit), and we haven't returned yet, rollback.
		if tx.Error != nil {
			log.Printf("INFO: CreateRequisitionHandler: Rolling back transaction due to error: %v\n", tx.Error)
			tx.Rollback()
		}
	}()

	// --- Start Transactional Operations ---

	// 1. Store items temporarily and clear from main payload to prevent GORM cascade on Requisition create.
	itemsToCreate := reqPayload.Items
	reqPayload.Items = nil // Important: Prevent GORM from attempting to cascade-create items here.

	// 2. Create the main Requisition record.
	if err := tx.Create(&reqPayload).Error; err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to insert requisition: %v\n", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to insert requisition: "+err.Error())
		tx.Rollback() // Explicit rollback
		return
	}

	// reqPayload now has the ID of the newly created requisition (e.g., reqPayload.ID)

	// 3. Iterate over the stored items, set their RequisitionID, ensure their own ID is 0 (for new items),
	//    and then create them individually.
	if len(itemsToCreate) > 0 {
		for i := range itemsToCreate {
			itemsToCreate[i].RequisitionID = reqPayload.ID // Set the foreign key
			itemsToCreate[i].ID = 0                       // CRITICAL: Ensure GORM treats item as new for auto-increment ID

			if err := tx.Create(&itemsToCreate[i]).Error; err != nil {
				log.Printf("ERROR: CreateRequisitionHandler: Failed to insert requisition item ('%s'): %v\n", itemsToCreate[i].Description, err)
				RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert requisition item ('%s'): %v", itemsToCreate[i].Description, err.Error()))
				tx.Rollback() // Explicit rollback
				return
			}
		}
	}

	// 4. If all operations were successful, commit the transaction.
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: CreateRequisitionHandler: Failed to commit transaction: %v\n", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
		// The defer func will handle rollback if tx.Error is set by Commit() failure.
		return
	}

	// --- End Transactional Operations ---

	log.Printf("INFO: CreateRequisitionHandler: Successfully created requisition ID %d with %d items\n", reqPayload.ID, len(itemsToCreate))
	
	// To send back the full requisition with its newly created items (and their DB-generated IDs):
	// We need to reload the requisition with its items. The `reqPayload` has the main requisition details,
	// and `itemsToCreate` has the item details with their new IDs.
	finalReqPayload := reqPayload
	finalReqPayload.Items = itemsToCreate // Items now have their DB-assigned IDs

	RespondWithJSON(w, http.StatusCreated, finalReqPayload)
}

// ListRequisitionsHandler handles GET requests to list requisitions for the authenticated user
func ListRequisitionsHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	if db == nil {
		log.Println("ERROR: ListRequisitionsHandler: Database not initialized")
		RespondWithError(w, http.StatusInternalServerError, "Database connection not initialized")
		return
	}

	userIDFromCtx := r.Context().Value("userID")
	if userIDFromCtx == nil {
		log.Println("ERROR: ListRequisitionsHandler: userID not found in context. Unauthorized.")
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User ID not found in request context.")
		return
	}

	userID, ok := userIDFromCtx.(int64) // Changed from uint to int64
	if !ok || userID == 0 {
		log.Printf("ERROR: ListRequisitionsHandler: Invalid userID type in context or userID is 0. userIDFromCtx: %v", userIDFromCtx)
		RespondWithError(w, http.StatusForbidden, "Forbidden: Invalid user identifier.")
		return
	}

	// Fetch the user from the database to check their role
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("ERROR: ListRequisitionsHandler: User with ID %d not found in database.", userID)
			RespondWithError(w, http.StatusForbidden, "Forbidden: User not found.")
			return
		}
		log.Printf("ERROR: ListRequisitionsHandler: Failed to query user %d: %v\n", userID, err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user details: "+err.Error())
		return
	}

	var requisitions []models.Requisition
	query := db.Preload("Items").Order("created_at desc")

	// Check user role
	// TODO: Make "procurement_officer" a constant or configurable value
	// Use case-insensitive comparison for the role
	if strings.EqualFold(user.Role, "procurement_officer") {
		// Procurement officers see all requisitions
		log.Printf("INFO: User %d (Role: %s) is a procurement officer, fetching all requisitions.", userID, user.Role)
	} else {
		// Other users see only their own requisitions
		log.Printf("INFO: User %d (Role: %s) is not a procurement officer, fetching only their requisitions.", userID, user.Role)
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&requisitions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			RespondWithJSON(w, http.StatusOK, []models.Requisition{}) // Send empty array
			return
		}
		log.Printf("ERROR: ListRequisitionsHandler: Failed to query requisitions for user %d: %v\n", userID, err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve requisitions: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, requisitions)
	log.Printf("INFO: Successfully retrieved %d requisitions for user ID: %d", len(requisitions), userID)
}

// GetRequisitionHandler handles GET requests to fetch a single requisition by ID for the authenticated user
func GetRequisitionHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	if db == nil {
		log.Println("ERROR: GetRequisitionHandler: Database not initialized")
		RespondWithError(w, http.StatusInternalServerError, "Database connection not initialized")
		return
	}

	// Get authenticated user ID from context
	userIDFromCtx := r.Context().Value("userID")
	if userIDFromCtx == nil {
		log.Println("ERROR: GetRequisitionHandler: userID not found in context. Unauthorized.")
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User ID not found in request context.")
		return
	}
	userID, ok := userIDFromCtx.(int64)
	if !ok || userID == 0 {
		log.Printf("ERROR: GetRequisitionHandler: Invalid userID type in context or userID is 0. userIDFromCtx: %v", userIDFromCtx)
		RespondWithError(w, http.StatusForbidden, "Forbidden: Invalid user identifier.")
		return
	}

	// Get requisition ID from URL parameter
	requisitionIDStr := chi.URLParam(r, "id")
	if requisitionIDStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Requisition ID is required")
		return
	}
	requisitionID, err := strconv.ParseInt(requisitionIDStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Requisition ID format")
		return
	}

	// Fetch the user from the database to check their role
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: GetRequisitionHandler: User with ID %d not found in database.", userID)
			RespondWithError(w, http.StatusForbidden, "Forbidden: User not found.")
			return
		}
		log.Printf("ERROR: GetRequisitionHandler: Failed to query user %d: %v\n", userID, err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user details: "+err.Error())
		return
	}

	// ADDED: Diagnostic logging for user role
	log.Printf("DIAGNOSTIC: GetRequisitionHandler: Fetched user for role check. UserID: %d, UserRole from DB: '%s'", user.ID, user.Role)

	var requisition models.Requisition
	query := db.Preload("Items")

	// Role-based access control
	// TODO: Make "procurement_officer" a constant or configurable value
	if strings.EqualFold(user.Role, "procurement_officer") {
		// Procurement officers can view any requisition by ID
		log.Printf("INFO: GetRequisitionHandler: User %d (Role: %s) is a procurement officer. Accessing requisition ID %d.", userID, user.Role, requisitionID)
		query = query.Where("id = ?", requisitionID)
	} else {
		// Other users can only view their own requisitions
		log.Printf("INFO: GetRequisitionHandler: User %d (Role: %s) is not a procurement officer. Accessing own requisition ID %d.", userID, user.Role, requisitionID)
		query = query.Where("id = ? AND user_id = ?", requisitionID, userID)
	}

	// Query for the specific requisition
	if err := query.First(&requisition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Log slightly different message based on role access
			if strings.EqualFold(user.Role, "procurement_officer") {
				log.Printf("WARN: GetRequisitionHandler: Requisition ID %d not found (Procurement Officer access).", requisitionID)
				RespondWithError(w, http.StatusNotFound, "Requisition not found.")
			} else {
				log.Printf("WARN: GetRequisitionHandler: Requisition ID %d not found for user ID %d or user lacks permission.", requisitionID, userID)
				RespondWithError(w, http.StatusNotFound, "Requisition not found or you do not have permission to view it.")
			}
		} else {
			log.Printf("ERROR: GetRequisitionHandler: Failed to query requisition ID %d (User ID %d, Role %s): %v\n", requisitionID, userID, user.Role, err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve requisition: "+err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, requisition)
	log.Printf("INFO: Successfully retrieved requisition ID %d for user ID %d (Role: %s)", requisition.ID, userID, user.Role)
}
