package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"procurement/models"
)

// TenderHandler holds dependencies for tender related handlers.
type TenderHandler struct {
	DB *gorm.DB
}

// NewTenderHandler creates a new TenderHandler with the given DB connection.
func NewTenderHandler(db *gorm.DB) *TenderHandler {
	return &TenderHandler{DB: db}
}

// CreateTender handles the creation of a new tender.
// POST /api/tenders
func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tenderInput models.Tender

	// Bind JSON input to the Tender struct
	if err := json.NewDecoder(r.Body).Decode(&tenderInput); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input: " + err.Error()})
		return
	}

	// TODO: Add validation logic here if needed
	// E.g., ensure required fields like Title, ClosingDate are present

	// Set CreatedByUserID from the authenticated user's ID in the request context
	userIDFromContext := r.Context().Value("userID")
	if userIDFromContext == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID not found in context, authentication required"})
		return
	}

	// userID in context is int64 (from User model's ID and auth middleware)
	userIDInt64FromCtx, ok := userIDFromContext.(int64)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID in context is of an unexpected type, expected int64"})
		return
	}

	// CreatedByUserID is already *int64, so assign the address of the int64 from context
	tenderInput.CreatedByUserID = &userIDInt64FromCtx

	// Save the tender to the database
	if err := h.DB.Create(&tenderInput).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create tender: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenderInput)
}

// GetTenders handles listing all tenders.
// GET /api/tenders
func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "User ID not found or invalid in context")
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user details: "+err.Error())
		return
	}

	log.Printf("GetTenders: UserID: %d, Role: %s", userID, user.Role)

	var tenders []models.Tender
	query := h.DB.Model(&models.Tender{}) // Start with a base query

	if strings.EqualFold(user.Role, "supplier") {
		log.Println("GetTenders: Applying supplier-specific filters")
		query = query.Where("status IN (?, ?)", "published", "open").Where("closing_date > ?", time.Now())

		// Filtering by category for suppliers
		category := r.URL.Query().Get("category")
		if category != "" {
			log.Printf("GetTenders: Supplier filtering by category: %s", category)
			query = query.Where("LOWER(category) = LOWER(?) ", category) // Case-insensitive category search
		}
		// TODO: Add other supplier filters like department if needed
		// TODO: Add sorting options for suppliers

	} else if strings.EqualFold(user.Role, "procurement_officer") {
		log.Println("GetTenders: Procurement officer fetching all tenders (or apply specific PO filters here)")
		// Procurement officers can see all tenders, or add specific filters for them if needed.
		// For now, no additional filters beyond the base query, so they get all.
	} else {
		// Other roles (e.g., requester) see only tenders they created
		log.Printf("GetTenders: User role '%s' fetching their created tenders", user.Role)
		query = query.Where("created_by_user_id = ?", userID)
	}

	// Execute the query
	if err := query.Order("created_at DESC").Find(&tenders).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tenders: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

// GetTenderByID handles retrieving a single tender by its ID.
// GET /api/tenders/:id
func (h *TenderHandler) GetTenderByID(w http.ResponseWriter, r *http.Request) {
	tenderID := chi.URLParam(r, "id") // Use chi.URLParam to get the ID
	var tender models.Tender

	// Preload Requisition and its Items. 
	// The Tender model must have a 'Requisition' field, and the Requisition model an 'Items' field.
	if err := h.DB.Preload("Requisition").Preload("Requisition.Items").First(&tender, tenderID).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Tender not found"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve tender: " + err.Error()})
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

// UpdateTender handles updating an existing tender by its ID.
// PUT /api/tenders/:id
func (h *TenderHandler) UpdateTender(w http.ResponseWriter, r *http.Request) {
	tenderID := chi.URLParam(r, "id")
	var tenderInput models.Tender

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&tenderInput); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input: " + err.Error()})
		return
	}

	var existingTender models.Tender
	// Check if the tender exists
	if err := h.DB.First(&existingTender, tenderID).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Tender not found"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve tender: " + err.Error()})
		}
		return
	}

	// Assign the ID from the path to ensure we're updating the correct record.
	tenderInput.ID = existingTender.ID

	// Update the tender record in the database
	if err := h.DB.Model(&existingTender).Updates(tenderInput).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update tender: " + err.Error()})
		return
	}

	// Fetch the updated record to return it
	if err := h.DB.First(&existingTender, tenderID).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve updated tender: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingTender)
}

// TODO: Add DeleteTender handler as needed.
