package models

import "time"

// Tender corresponds to the Tenders table in the database.
// It represents a formal invitation for suppliers to submit a bid to supply goods or services.
type Tender struct {
	ID                 int64      `json:"id" gorm:"primaryKey"`
	RequisitionID      *int64     `json:"requisition_id,omitempty" gorm:"index"` // Link to the original Requisition
	Title              string     `json:"title" gorm:"not null"`
	Description        *string    `json:"description,omitempty"`
	Category           *string    `json:"category,omitempty"`         // E.g., 'goods', 'services', 'works', 'consultancy'
	Budget             *float64   `json:"budget,omitempty"`           // Estimated budget for the tender
	Status             *string    `json:"status,omitempty" gorm:"default:'draft'"` // E.g., 'draft', 'published', 'evaluation', 'awarded', 'cancelled'
	PublishedDate      *time.Time `json:"published_date,omitempty"`   // Date when the tender is made public
	ClosingDate        *time.Time `json:"closing_date,omitempty"`     // Deadline for bid submissions
	EvaluationMethod   *string    `json:"evaluation_method,omitempty"`// E.g., 'least_cost', 'quality_cost_based'
	BidOpeningDate     *time.Time `json:"bid_opening_date,omitempty"` // Date when bids will be opened
	CreatedByUserID    *int64     `json:"created_by_user_id,omitempty"` // User who created the tender
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Requisition *Requisition `json:"requisition,omitempty" gorm:"foreignKey:RequisitionID;references:ID"` // Embed Requisition details. Pointer to allow omitempty.
	// User        User        `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedByUserID"` // If you want to embed User details

	// Has-many relationship: A Tender can have multiple Bids
	Bids []Bid `json:"bids,omitempty" gorm:"foreignKey:TenderID"`
}

// TenderItem might be needed if tenders have their own line items distinct from requisition items
// For now, we assume the tender refers to the overall requisition or its items implicitly.
// If a tender can consolidate multiple requisitions or specify items differently,
// a TenderItem struct and a has-many relationship would be necessary.
// type TenderItem struct {
// 	ID          int64   `json:"id" gorm:"primaryKey"`
// 	TenderID    int64   `json:"tender_id" gorm:"index;not null"`
// 	Description string  `json:"description" gorm:"not null"`
// 	Quantity    float64 `json:"quantity"`
// 	Unit        string  `json:"unit"`
// 	// ... other fields like estimated cost if different from requisition
// }
