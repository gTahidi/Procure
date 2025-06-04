package models

import "time"

// BidItem represents a specific item within a supplier's bid, corresponding to an item in the tender.
// It includes the supplier's offered price and specifications for that item.
// This table will store the details for each item in a bid.
// The combination of BidID and RequisitionItemID (or a similar unique identifier for the tender item)
// should ideally be unique if suppliers bid on pre-defined tender items.
// If suppliers can add ad-hoc items, RequisitionItemID might be nullable.
type BidItem struct {
	ID                  int64      `json:"id" gorm:"primaryKey"`
	BidID               int64      `json:"bid_id" gorm:"index;not null"` // Foreign key to the Bid
	RequisitionItemID   *int64     `json:"requisition_item_id,omitempty" gorm:"index"` // Foreign key to the original RequisitionItem, if applicable
	Description         string     `json:"description" gorm:"not null"` // Can be copied from RequisitionItem or provided by supplier
	Quantity            float64    `json:"quantity" gorm:"not null"`      // Typically copied from RequisitionItem
	Unit                string     `json:"unit" gorm:"not null"`          // Typically copied from RequisitionItem
	OfferedUnitPrice    float64    `json:"offered_unit_price" gorm:"not null"`
	SpecificationText   *string    `json:"specification_text,omitempty"`
	SpecificationSheetURL *string  `json:"specification_sheet_url,omitempty"` // URL/path to uploaded specification sheet
	ItemImageURL        *string    `json:"item_image_url,omitempty"`          // URL/path to uploaded item image
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Optional: Associations if needed, e.g., back to RequisitionItem
	// RequisitionItem   *RequisitionItem `json:"requisition_item,omitempty" gorm:"foreignKey:RequisitionItemID"`
}
