package models

import (
	"time"
)

// Bid represents a bid submitted by a supplier for a tender.
type Bid struct {
	ID                   int64      `json:"id" gorm:"primaryKey"`
	TenderID             int64      `json:"tender_id" gorm:"index;not null"`
	SupplierID           int64      `json:"supplier_id" gorm:"index;not null"`
	BidAmount            float64    `json:"bid_amount" gorm:"not null"`
	SubmissionDate       time.Time  `json:"submission_date" gorm:"autoCreateTime"`
	TechnicalProposalURL *string    `json:"technical_proposal_url,omitempty"`
	FinancialProposalURL *string    `json:"financial_proposal_url,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	Status               string     `json:"status" gorm:"default:'submitted';not null"` // e.g., 'submitted', 'under_review', 'shortlisted', 'rejected', 'awarded', 'withdrawn'
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Tender   Tender `json:"tender,omitempty" gorm:"foreignKey:TenderID"`
	Supplier User   `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Items    []BidItem `json:"items,omitempty" gorm:"foreignKey:BidID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // A bid comprises multiple items
}
