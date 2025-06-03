package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"procurement/models"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// BidHandler handles HTTP requests for bids.
type BidHandler struct {
	DB *gorm.DB
}

// NewBidHandler creates a new BidHandler with the given database connection.
func NewBidHandler(db *gorm.DB) *BidHandler {
	return &BidHandler{DB: db}
}

// CreateBid handles the submission of a new bid for a tender.
// POST /api/tenders/{tenderId}/bids
func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user details from context
	userIDFromContext, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "User ID not found or invalid in context")
		return
	}

	var currentUser models.User
	if err := h.DB.First(&currentUser, userIDFromContext).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user details: "+err.Error())
		return
	}

	// Role check: Only suppliers can create bids
	if !strings.EqualFold(currentUser.Role, "supplier") {
		RespondWithError(w, http.StatusForbidden, "Forbidden: Only suppliers can submit bids.")
		return
	}

	// Get Tender ID from URL path parameter
	tenderIDStr := chi.URLParam(r, "tenderId")
	tenderID, err := strconv.ParseInt(tenderIDStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid tender ID format: "+err.Error())
		return
	}
	log.Printf("CreateBid: Attempting to create bid for TenderID: %d by SupplierID: %d", tenderID, currentUser.ID)

	// Fetch the target tender to check its status
	var tender models.Tender
	if err := h.DB.First(&tender, tenderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			RespondWithError(w, http.StatusNotFound, "Tender not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tender: "+err.Error())
		}
		return
	}

	// Check if tender is open for bidding
	// Ensure Status is non-nil before dereferencing
	var tenderStatus string
	if tender.Status != nil {
		tenderStatus = *tender.Status
	} else {
		RespondWithError(w, http.StatusBadRequest, "Tender status is not set.")
		return
	}

	if !strings.EqualFold(tenderStatus, "published") {
		RespondWithError(w, http.StatusBadRequest, "Tender is not published and thus not open for bidding.")
		return
	}
	if tender.ClosingDate == nil || !tender.ClosingDate.After(time.Now()) {
		RespondWithError(w, http.StatusBadRequest, "Tender is past its closing date or closing date not set.")
		return
	}
	log.Printf("CreateBid: TenderID %d is open for bidding.", tenderID)

	// Parse the request body for bid details
	var bidInput models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bidInput); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Validate bid input (e.g., BidAmount)
	if bidInput.BidAmount <= 0 {
		RespondWithError(w, http.StatusBadRequest, "Bid amount must be greater than zero.")
		return
	}

	// Populate bid details
	bidInput.TenderID = tenderID
	bidInput.SupplierID = currentUser.ID // This is the ID of the authenticated supplier
	// bidInput.Status is defaulted to 'submitted' by the model
	// bidInput.SubmissionDate is defaulted by autoCreateTime

	// Save the bid to the database
	if err := h.DB.Create(&bidInput).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create bid: "+err.Error())
		return
	}
	log.Printf("CreateBid: Successfully created BidID: %d for TenderID: %d by SupplierID: %d", bidInput.ID, tenderID, currentUser.ID)

	// Respond with the created bid
	RespondWithJSON(w, http.StatusCreated, bidInput)
}

// ListTenderBids handles listing all bids for a specific tender.
// GET /api/tenders/{tenderId}/bids
// Accessible by procurement officers.
func (h *BidHandler) ListTenderBids(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user details from context
	userIDFromContext, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "User ID not found or invalid in context")
		return
	}

	var currentUser models.User
	if err := h.DB.First(&currentUser, userIDFromContext).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user details: "+err.Error())
		return
	}

	// Role check: Only procurement officers can list all bids for a tender
	if !strings.EqualFold(currentUser.Role, "procurement_officer") {
		RespondWithError(w, http.StatusForbidden, "Forbidden: Only procurement officers can view all bids for a tender.")
		return
	}

	// Get Tender ID from URL path parameter
	tenderIDStr := chi.URLParam(r, "tenderId")
	tenderID, err := strconv.ParseInt(tenderIDStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid tender ID format: "+err.Error())
		return
	}
	log.Printf("ListTenderBids: ProcurementOfficerID: %d attempting to list bids for TenderID: %d", currentUser.ID, tenderID)

	// Optional: Verify tender exists (though listing bids for a non-existent tender will just return empty)
	var tender models.Tender
	if err := h.DB.First(&tender, tenderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			RespondWithError(w, http.StatusNotFound, "Tender not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Error checking tender: "+err.Error())
		}
		return
	}

	// Fetch bids for the tender, preloading supplier information
	var bids []models.Bid
	if err := h.DB.Preload("Supplier").Where("tender_id = ?", tenderID).Order("submission_date ASC").Find(&bids).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bids: "+err.Error())
		return
	}
	log.Printf("ListTenderBids: Found %d bids for TenderID: %d", len(bids), tenderID)

	RespondWithJSON(w, http.StatusOK, bids)
}

// ListMyBids handles listing all bids submitted by the authenticated supplier.
// GET /api/my-bids
// Accessible by suppliers.
func (h *BidHandler) ListMyBids(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user details from context
	userIDFromContext, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "User ID not found or invalid in context")
		return
	}

	var currentUser models.User
	if err := h.DB.First(&currentUser, userIDFromContext).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user details: "+err.Error())
		return
	}

	// Role check: Only suppliers can list their own bids
	if !strings.EqualFold(currentUser.Role, "supplier") {
		RespondWithError(w, http.StatusForbidden, "Forbidden: Only suppliers can view their submitted bids.")
		return
	}
	log.Printf("ListMyBids: SupplierID: %d attempting to list their bids", currentUser.ID)

	// Fetch bids submitted by the current supplier, preloading Tender information
	var myBids []models.Bid
	if err := h.DB.Preload("Tender").Where("supplier_id = ?", currentUser.ID).Order("submission_date DESC").Find(&myBids).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve your bids: "+err.Error())
		return
	}
	log.Printf("ListMyBids: Found %d bids for SupplierID: %d", len(myBids), currentUser.ID)

	RespondWithJSON(w, http.StatusOK, myBids)
}
