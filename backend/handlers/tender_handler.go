package handlers

import (
	"encoding/json"
	"net/http"

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
	var tenders []models.Tender
	// TODO: Add pagination, filtering (e.g., by status, category) as needed
	if err := h.DB.Find(&tenders).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve tenders: " + err.Error()})
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

	if err := h.DB.First(&tender, tenderID).Error; err != nil {
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
