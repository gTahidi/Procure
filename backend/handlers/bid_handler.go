package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	// Define max upload size (e.g., 10MB per file, overall 50MB)
	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	// r.ParseMultipartForm needs to be called before accessing form data
	if err := r.ParseMultipartForm(50 * 1024 * 1024); err != nil { // 50MB total max size
		RespondWithError(w, http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
		return
	}

	var bidInput models.Bid
	// General bid fields (non-item related)
	bidInput.Notes = getFormValuePointer(r, "notes")
	// Add other general fields like validity period if they are added to the form and model

	// Bid items are expected as a JSON string in a field named 'items_json'
	itemsJSON := r.FormValue("items_json")
	if itemsJSON == "" {
		RespondWithError(w, http.StatusBadRequest, "Bid items (items_json) are required.")
		return
	}

	var bidItems []models.BidItem
	if err := json.Unmarshal([]byte(itemsJSON), &bidItems); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid format for bid items (items_json): "+err.Error())
		return
	}

	if len(bidItems) == 0 {
		RespondWithError(w, http.StatusBadRequest, "At least one bid item is required.")
		return
	}

	// Populate bid details
	bidInput.TenderID = tenderID
	bidInput.SupplierID = currentUser.ID
	// bidInput.Status is defaulted by model
	// bidInput.BidAmount might be calculated or set based on items or a general field
	// For now, let's assume it might still be a general field or we'll calculate it later.
	bidAmountStr := r.FormValue("bid_amount") // If still sending overall bid amount
	if bidAmountStr != "" {
		bidAmountFloat, err := strconv.ParseFloat(bidAmountStr, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid bid_amount format: "+err.Error())
			return
		}
		bidInput.BidAmount = bidAmountFloat
	} else if len(bidItems) > 0 {
		// If no overall bid_amount, calculate from items
		var totalCalculatedAmount float64
		for _, item := range bidItems {
			totalCalculatedAmount += item.OfferedUnitPrice * item.Quantity
		}
		bidInput.BidAmount = totalCalculatedAmount
	}

	if bidInput.BidAmount <= 0 {
		RespondWithError(w, http.StatusBadRequest, "Total bid amount must be greater than zero.")
		return
	}

	// Start a transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to start database transaction: "+tx.Error.Error())
		return
	}

	// Save the main Bid record first to get its ID
	if err := tx.Create(&bidInput).Error; err != nil {
		tx.Rollback()
		RespondWithError(w, http.StatusInternalServerError, "Failed to create bid (main record): "+err.Error())
		return
	}
	log.Printf("CreateBid: Successfully created BidID: %d for TenderID: %d by SupplierID: %d (pre-items)", bidInput.ID, tenderID, currentUser.ID)

	// Process and save BidItems and their files
	for i := range bidItems {
		bidItems[i].BidID = bidInput.ID // Link item to the created Bid

		// Handle Specification Sheet File
		specSheetKey := fmt.Sprintf("item_spec_sheet_%d", i)
		file, header, err := r.FormFile(specSheetKey)
		if err == nil { // File is present
			defer file.Close()
			if header.Size > maxFileSize {
				tx.Rollback()
				RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Specification sheet for item %d (%s) exceeds max size of %dMB", i+1, header.Filename, maxFileSize/1024/1024))
				return
			}
			// Ensure directory exists: ./uploads/bids/{bid_id}/items/{item_index}/specs/
			// Using item.ID might be problematic if it's not set yet, so use index or generate UUID for filename
			filePath := filepath.Join(".", "uploads", "bids", strconv.FormatInt(bidInput.ID, 10), "items", strconv.Itoa(i), "specs", SanitizeFilename(header.Filename))
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				tx.Rollback()
				RespondWithError(w, http.StatusInternalServerError, "Failed to create directory for spec sheet: "+err.Error())
				return
			}
			dst, err := os.Create(filePath)
			if err != nil {
				tx.Rollback()
				RespondWithError(w, http.StatusInternalServerError, "Failed to create file for spec sheet: "+err.Error())
				return
			}
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				tx.Rollback()
				RespondWithError(w, http.StatusInternalServerError, "Failed to save spec sheet: "+err.Error())
				return
			}
			bidItems[i].SpecificationSheetURL = &filePath
		} else if err != http.ErrMissingFile {
			tx.Rollback()
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error processing spec sheet for item %d: %s", i+1, err.Error()))
			return
		}

		// Handle Item Image File
		itemImageKey := fmt.Sprintf("item_image_%d", i)
		file, header, err = r.FormFile(itemImageKey)
		if err == nil { // File is present
			defer file.Close()
			if header.Size > maxFileSize {
				tx.Rollback()
				RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Item image for item %d (%s) exceeds max size of %dMB", i+1, header.Filename, maxFileSize/1024/1024))
				return
			}
			filePath := filepath.Join(".", "uploads", "bids", strconv.FormatInt(bidInput.ID, 10), "items", strconv.Itoa(i), "images", SanitizeFilename(header.Filename))
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				tx.Rollback()
				RespondWithError(w, http.StatusInternalServerError, "Failed to create directory for item image: "+err.Error())
				return
			}
			dst, err := os.Create(filePath)
			if err != nil {
				tx.Rollback()
				RespondWithError(w, http.StatusInternalServerError, "Failed to create file for item image: "+err.Error())
				return
			}
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				tx.Rollback()
				RespondWithError(w, http.StatusInternalServerError, "Failed to save item image: "+err.Error())
				return
			}
			bidItems[i].ItemImageURL = &filePath
		} else if err != http.ErrMissingFile {
			tx.Rollback()
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error processing item image for item %d: %s", i+1, err.Error()))
			return
		}

		// Save the BidItem
		if err := tx.Create(&bidItems[i]).Error; err != nil {
			tx.Rollback()
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save bid item %d: %s", i+1, err.Error()))
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
		return
	}

	log.Printf("CreateBid: Successfully created BidID: %d with %d items for TenderID: %d by SupplierID: %d", bidInput.ID, len(bidItems), tenderID, currentUser.ID)

	// Reload the bid with its items to return the full object
	var finalBid models.Bid
	if err := h.DB.Preload("Items").First(&finalBid, bidInput.ID).Error; err != nil {
	    log.Printf("Error reloading bid with items: %v. Returning bid without items.", err)
	    RespondWithJSON(w, http.StatusCreated, bidInput) // Fallback to returning bidInput without items if reload fails
	    return
	}

	RespondWithJSON(w, http.StatusCreated, finalBid)
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
