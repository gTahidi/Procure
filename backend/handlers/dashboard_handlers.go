package handlers

import (
	"fmt"
	"net/http"
	"time"

	"procurement/database"
	"procurement/models"
)

// CreationRateStats represents the data for the creation rate chart.
type CreationRateStats struct {
	Requisitions []DailyCreationCount `json:"requisitions"`
	Tenders      []DailyCreationCount `json:"tenders"`
}

// DailyCreationCount represents the number of items created on a specific day.
type DailyCreationCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// DashboardStats represents the statistics for the procurement dashboard.
type DashboardStats struct {
	PendingApproval int64 `json:"pendingApproval"`
	ReadyForTender  int64 `json:"readyForTender"`
	ActiveTenders   int64 `json:"activeTenders"`
	RecentlyClosed  int64 `json:"recentlyClosed"`
}

// RecentRequisition represents a summarized requisition for the dashboard.
type RecentRequisition struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	CreationDate time.Time `json:"creationDate"`
}

// LiveTender represents a summarized tender for the dashboard.
type LiveTender struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	ClosingDate time.Time `json:"closingDate"`
}

// GetCreationRateHandler fetches the creation rate of requisitions and tenders over the last 30 days.
func GetCreationRateHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	stats := CreationRateStats{}

	// Fetch requisition creation data
	rows, err := db.Model(&models.Requisition{}).Select("DATE(created_at) as date, count(*) as count").Where("created_at > ?", time.Now().AddDate(0, 0, -30)).Group("DATE(created_at)").Rows()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch requisition creation data")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dailyCount DailyCreationCount
		if err := rows.Scan(&dailyCount.Date, &dailyCount.Count); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to scan requisition creation data")
			return
		}
		stats.Requisitions = append(stats.Requisitions, dailyCount)
	}

	// Fetch tender creation data
	rows, err = db.Model(&models.Tender{}).Select("DATE(created_at) as date, count(*) as count").Where("created_at > ?", time.Now().AddDate(0, 0, -30)).Group("DATE(created_at)").Rows()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch tender creation data")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dailyCount DailyCreationCount
		if err := rows.Scan(&dailyCount.Date, &dailyCount.Count); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to scan tender creation data")
			return
		}
		stats.Tenders = append(stats.Tenders, dailyCount)
	}

	RespondWithJSON(w, http.StatusOK, stats)
}

// GetRequisitionStatsHandler calculates and returns statistics for requisitions.
func GetRequisitionStatsHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	stats := DashboardStats{}

	// Count requisitions pending approval
	db.Model(&models.Requisition{}).Where("status IN (?)", []string{string(models.RequisitionStatusPendingApproval1), string(models.RequisitionStatusPendingApproval2)}).Count(&stats.PendingApproval)

	// Count requisitions ready for tender
	db.Model(&models.Requisition{}).Where("status IN (?)", []string{string(models.RequisitionStatusApproved), string(models.RequisitionStatusPendingTender)}).Count(&stats.ReadyForTender)

	// Count requisitions that have been tendered
	db.Model(&models.Tender{}).Where("status = ?", "published").Count(&stats.ActiveTenders)

	// Count requisitions recently closed or rejected (in the last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	db.Model(&models.Requisition{}).Where("status IN (?) AND updated_at > ?", []string{string(models.RequisitionStatusClosed), string(models.RequisitionStatusRejected)}, thirtyDaysAgo).Count(&stats.RecentlyClosed)

	RespondWithJSON(w, http.StatusOK, stats)
}

// GetRecentRequisitionsHandler fetches the 5 most recently created requisitions.
func GetRecentRequisitionsHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	var requisitions []models.Requisition

	if err := db.Order("created_at desc").Limit(5).Find(&requisitions).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch recent requisitions")
		return
	}

	responseRequisitions := make([]RecentRequisition, len(requisitions))
	for i, req := range requisitions {
		responseRequisitions[i] = RecentRequisition{
			ID:           uint(req.ID),
			Title:        fmt.Sprintf("Requisition - %03d", req.ID),
			Status:       string(req.Status),
			CreationDate: req.CreatedAt,
		}
	}

	RespondWithJSON(w, http.StatusOK, responseRequisitions)
}

// GetLiveTendersHandler fetches all tenders with a 'published' status.
func GetLiveTendersHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	var tenders []models.Tender

	if err := db.Where("status = ?", "published").Order("closing_date asc").Find(&tenders).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch live tenders")
		return
	}

	responseTenders := make([]LiveTender, len(tenders))
	for i, t := range tenders {
		var category string
		if t.Category != nil {
			category = *t.Category
		}

		var closingDate time.Time
		if t.ClosingDate != nil {
			closingDate = *t.ClosingDate
		}

		responseTenders[i] = LiveTender{
			ID:          uint(t.ID),
			Title:       t.Title,
			Category:    category,
			ClosingDate: closingDate,
		}
	}

	RespondWithJSON(w, http.StatusOK, responseTenders)
}
