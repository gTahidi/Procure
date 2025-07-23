package handlers

import (
	"encoding/json"
	"errors" // Added for errors.Is
	"fmt"
	"log"
	"net/http"
	"procurement/database" // Module name 'procurement' then path
	"procurement/models"
	"strconv" // Added for strconv.ParseInt
	"strings" // Added for strings.EqualFold
	"time"    // Needed for setting approval timestamps

	"gorm.io/gorm" // Added for gorm.ErrRecordNotFound or other GORM specific needs

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

	// Get User ID from context
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		log.Println("ERROR: CreateRequisitionHandler: Could not retrieve userID from context or type assertion failed.")
		RespondWithError(w, http.StatusInternalServerError, "Could not process request: user authentication issue.")
		return
	}

	// Assign the authenticated user's ID to the requisition
	reqPayload.UserID = int64(userID)

	// Basic Validation (example)
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
			itemsToCreate[i].ID = 0                        // CRITICAL: Ensure GORM treats item as new for auto-increment ID

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
	if strings.EqualFold(user.Role, "procurement_officer") || strings.EqualFold(user.Role, "admin") {
		// Procurement officers see all requisitions
		log.Printf("INFO: User %d (Role: %s) is a procurement officer or admin, fetching all requisitions.", userID, user.Role)
	} else {
		// Other users see only their own requisitions
		log.Printf("INFO: User %d (Role: %s) is not a procurement officer or admin, fetching only their requisitions.", userID, user.Role)
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
	if strings.EqualFold(user.Role, "procurement_officer") || strings.EqualFold(user.Role, "admin") {
		// Procurement officers can view any requisition by ID
		log.Printf("INFO: GetRequisitionHandler: User %d (Role: %s) is a procurement officer or admin. Accessing requisition ID %d.", userID, user.Role, requisitionID)
		query = query.Where("id = ?", requisitionID)
	} else {
		// Other users can only view their own requisitions
		log.Printf("INFO: GetRequisitionHandler: User %d (Role: %s) is not a procurement officer or admin. Accessing own requisition ID %d.", userID, user.Role, requisitionID)
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

// RequisitionActionPayload defines the structure for the request body of requisition actions
type RequisitionActionPayload struct {
	Action string `json:"action"`           // "approve" or "reject"
	Reason string `json:"reason,omitempty"` // Required if action is "reject"
}

// HandleRequisitionAction handles POST requests to approve or reject a requisition
// MyRequisitionStats defines the statistics for a requester's personal dashboard.
type MyRequisitionStats struct {
	Pending  int64 `json:"pending"`
	Approved int64 `json:"approved"`
	Rejected int64 `json:"rejected"`
	Draft    int64 `json:"draft"` // Assuming 'draft' is not a formal status, but we can add logic if needed.
}

// GetMyRequisitionStatsHandler calculates and returns statistics for the current user's requisitions.
func GetMyRequisitionStatsHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Could not identify user.")
		return
	}

	var stats MyRequisitionStats
	db.Model(&models.Requisition{}).Where("user_id = ? AND status IN (?)", userID, []string{string(models.RequisitionStatusPendingApproval1), string(models.RequisitionStatusPendingApproval2)}).Count(&stats.Pending)
	db.Model(&models.Requisition{}).Where("user_id = ? AND status IN (?)", userID, []string{string(models.RequisitionStatusApproved), string(models.RequisitionStatusPendingTender), string(models.RequisitionStatusTendered)}).Count(&stats.Approved)
	db.Model(&models.Requisition{}).Where("user_id = ? AND status = ?", userID, models.RequisitionStatusRejected).Count(&stats.Rejected)

	RespondWithJSON(w, http.StatusOK, stats)
}

// GetMyRecentRequisitionsHandler fetches the 5 most recent requisitions for the current user.
func GetMyRecentRequisitionsHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Could not identify user.")
		return
	}

	var requisitions []models.Requisition
	if err := db.Where("user_id = ?", userID).Order("created_at desc").Limit(5).Find(&requisitions).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch recent requisitions")
		return
	}

	RespondWithJSON(w, http.StatusOK, requisitions)
}

func HandleRequisitionAction(w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: HandleRequisitionAction: Entered function.")
	db := database.GetDB()
	if db == nil {
		log.Println("ERROR: HandleRequisitionAction: Database not initialized")
		RespondWithError(w, http.StatusInternalServerError, "Database connection not initialized")
		return
	}

	// Get authenticated user ID from context
	userIDFromCtx := r.Context().Value("userID")
	if userIDFromCtx == nil {
		log.Println("ERROR: HandleRequisitionAction: userID not found in context. Unauthorized.")
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User ID not found in request context.")
		return
	}

	adminID, okUserID := userIDFromCtx.(int64)
	if !okUserID || adminID == 0 {
		log.Printf("ERROR: HandleRequisitionAction: Invalid userID type in context. userIDFromCtx: %v", userIDFromCtx)
		RespondWithError(w, http.StatusForbidden, "Forbidden: Invalid user identifier.")
		return
	}

	// Fetch the admin user from DB to get their role
	var adminUser models.User
	if err := db.First(&adminUser, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: HandleRequisitionAction: Admin user with ID %d not found in DB.", adminID)
			RespondWithError(w, http.StatusUnauthorized, "Unauthorized: Admin user not found.")
		} else {
			log.Printf("ERROR: HandleRequisitionAction: Error fetching admin user %d: %v", adminID, err)
			RespondWithError(w, http.StatusInternalServerError, "Error verifying admin user.")
		}
		return
	}
	adminRole := adminUser.Role

	// Check if the user is an admin
	// TODO: Use a constant for "admin" role
	if !strings.EqualFold(adminRole, "admin") {
		log.Printf("WARN: HandleRequisitionAction: User %d (Role: %s) attempted to perform admin action.", adminID, adminRole)
		RespondWithError(w, http.StatusForbidden, "Forbidden: This action requires admin privileges.")
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
		log.Printf("ERROR: HandleRequisitionAction: Invalid Requisition ID format '%s': %v", requisitionIDStr, err)
		RespondWithError(w, http.StatusBadRequest, "Invalid Requisition ID format")
		return
	}
	log.Printf("DEBUG: HandleRequisitionAction: Requisition ID parsed: %d. Admin ID: %d. Admin Role: %s", requisitionID, adminID, adminRole)

	// Decode the request body
	var payload RequisitionActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("ERROR: HandleRequisitionAction: Failed to decode request body: %v\n", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()
	log.Printf("DEBUG: HandleRequisitionAction: Payload decoded successfully: Action='%s', Reason='%s'", payload.Action, payload.Reason)

	// Validate action
	payload.Action = strings.ToLower(payload.Action)
	if payload.Action != "approve" && payload.Action != "reject" {
		RespondWithError(w, http.StatusBadRequest, "Invalid action specified. Must be 'approve' or 'reject'.")
		return
	}

	if payload.Action == "reject" && strings.TrimSpace(payload.Reason) == "" {
		RespondWithError(w, http.StatusBadRequest, "Rejection reason is required when action is 'reject'.")
		return
	}

	var requisition models.Requisition
	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: HandleRequisitionAction: Failed to begin transaction: %v\n", tx.Error)
		RespondWithError(w, http.StatusInternalServerError, "Failed to process request: "+tx.Error.Error())
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("PANIC: HandleRequisitionAction: Recovered from panic: %v\n", r)
			// RespondWithError might not be possible if headers already sent
			return
		}
		if tx.Error != nil {
			log.Printf("INFO: HandleRequisitionAction: Rolling back transaction due to error: %v\n", tx.Error)
			tx.Rollback()
		}
	}()

	if err := tx.First(&requisition, requisitionID).Error; err != nil { // Fetched within transaction
		if errors.Is(err, gorm.ErrRecordNotFound) {
			RespondWithError(w, http.StatusNotFound, "Requisition not found.")
		} else {
			log.Printf("ERROR: HandleRequisitionAction: Failed to query requisition ID %d: %v\n", requisitionID, err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve requisition: "+err.Error())
		}
		tx.Rollback() // Explicit rollback
		return
	}

	// Perform action based on current status
	log.Printf("DEBUG: HandleRequisitionAction: Validated payload. Action: %s. Requisition ID: %d. Current Requisition Status from DB: %s", payload.Action, requisition.ID, requisition.Status)
	switch payload.Action {
	case "approve":
		switch requisition.Status {
		case models.RequisitionStatusPendingApproval1, "submitted_for_approval": // Accept both for first approval
			// Ensure it's not already approved by this admin if it was somehow in pending_approval_1 but had an approver
			if requisition.ApproverOneID != nil && *requisition.ApproverOneID == adminID && requisition.Status == models.RequisitionStatusPendingApproval2 {
				log.Printf("WARN: HandleRequisitionAction: Admin %d attempted to re-approve requisition %d which they already first-approved.", adminID, requisitionID)
				RespondWithError(w, http.StatusForbidden, "You have already performed the first approval.")
				tx.Rollback()
				return
			}
			requisition.ApproverOneID = &adminID
			now := time.Now()
			requisition.ApprovedOneAt = &now
			requisition.Status = models.RequisitionStatusPendingApproval2
			log.Printf("INFO: HandleRequisitionAction: Requisition %d approved (1st approval) by admin %d. Status -> %s\n", requisitionID, adminID, requisition.Status)
		case models.RequisitionStatusPendingApproval2:
			if requisition.ApproverOneID != nil && *requisition.ApproverOneID == adminID {
				RespondWithError(w, http.StatusForbidden, "Second approval must be by a different admin.")
				tx.Rollback()
				return
			}
			requisition.ApproverTwoID = &adminID
			now := time.Now()
			requisition.ApprovedTwoAt = &now
			requisition.Status = models.RequisitionStatusApproved
			log.Printf("INFO: HandleRequisitionAction: Requisition %d approved (2nd approval) by admin %d. Status -> %s\n", requisitionID, adminID, requisition.Status)
		default:
			log.Printf("ERROR: HandleRequisitionAction: Attempt to approve requisition %d in unexpected status '%s' by admin %d.", requisitionID, requisition.Status, adminID)
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cannot approve requisition in status '%s'.", requisition.Status))
			tx.Rollback() // Ensure rollback
			return
		}
	case "reject":
		if requisition.Status == models.RequisitionStatusApproved || requisition.Status == models.RequisitionStatusTendered || requisition.Status == models.RequisitionStatusClosed {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Requisition cannot be rejected. Current status: %s", requisition.Status))
			tx.Rollback()
			return
		}
		// For simplicity, any admin can reject at pending_approval_1 or pending_approval_2 stage.
		// We could record who rejected it if needed, e.g., by setting ApproverOneID/TwoID with a negative value or a dedicated RejectedByID field.
		requisition.Status = models.RequisitionStatusRejected
		requisition.RejectionReason = &payload.Reason
		// Optionally, clear approval fields if it's rejected after first approval
		// requisition.ApproverOneID = nil
		// requisition.ApprovedOneAt = nil
		log.Printf("INFO: HandleRequisitionAction: Requisition %d rejected by admin %d. Reason: %s. Status -> %s\n", requisitionID, adminID, payload.Reason, requisition.Status)
	}

	if err := tx.Save(&requisition).Error; err != nil {
		log.Printf("ERROR: HandleRequisitionAction: Failed to save requisition ID %d: %v\n", requisitionID, err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to update requisition: "+err.Error())
		tx.Rollback() // Explicit rollback
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: HandleRequisitionAction: Failed to commit transaction for requisition ID %d: %v\n", requisitionID, err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to finalize requisition update: "+err.Error())
		// tx.Rollback() is implicitly handled by defer if Commit fails and sets tx.Error
		return
	}

	RespondWithJSON(w, http.StatusOK, requisition)
	log.Printf("INFO: HandleRequisitionAction: Successfully processed action '%s' for requisition ID %d by admin %d\n", payload.Action, requisitionID, adminID)
}
