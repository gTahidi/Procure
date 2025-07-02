package handlers

import (
	"net/http"
	"procurement/database"
	"procurement/models"
)

// SupplierDashboardData represents the data for the supplier dashboard.
type SupplierDashboardData struct {
	BidsSubmitted int64             `json:"bids_submitted"`
	BidsAwarded   int64             `json:"bids_awarded"`
	ActiveTenders []models.Tender   `json:"active_tenders"`
	MyBids        []models.Bid      `json:"my_bids"`
}

// GetSupplierDashboardDataHandler fetches all the necessary data for the supplier dashboard.
func GetSupplierDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Could not identify user.")
		return
	}

	var data SupplierDashboardData

	// Get stats
	db.Model(&models.Bid{}).Where("supplier_id = ?", userID).Count(&data.BidsSubmitted)
	db.Model(&models.Bid{}).Where("supplier_id = ? AND status = ?", userID, "awarded").Count(&data.BidsAwarded)

	// Get active tenders (published)
	db.Where("status = ?", "published").Find(&data.ActiveTenders)

	// Get my bids
	db.Where("supplier_id = ?", userID).Preload("Tender").Order("submission_date desc").Find(&data.MyBids)

	RespondWithJSON(w, http.StatusOK, data)
}
