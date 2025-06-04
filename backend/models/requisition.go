package models

import "time"

// RequisitionStatus defines the possible statuses for a Requisition.
type RequisitionStatus string

const (
	RequisitionStatusPendingApproval1 RequisitionStatus = "pending_approval_1"
	RequisitionStatusPendingApproval2 RequisitionStatus = "pending_approval_2"
	RequisitionStatusApproved         RequisitionStatus = "Approved"
	RequisitionStatusRejected         RequisitionStatus = "rejected"
	RequisitionStatusPendingTender    RequisitionStatus = "pending_tender" // Or 'awaiting_tender', 'ready_for_tender'
	RequisitionStatusTendered         RequisitionStatus = "tendered"
	RequisitionStatusClosed           RequisitionStatus = "closed" // e.g., after tender awarded or PR cancelled
)

// Requisition corresponds to the Requisitions table
type Requisition struct {
	ID            int64             `json:"id" gorm:"primaryKey"`
	UserID        int64             `json:"user_id" gorm:"index"`     // User who created the PR
	Type          string            `json:"type"`                     // 'goods', 'services', 'fixed_asset'
	AAC           *string           `json:"aac,omitempty"`            // 'A', 'F', 'P' (nullable)
	MaterialGroup *string           `json:"material_group,omitempty"` // (nullable)
	ExchangeRate  *float64          `json:"exchange_rate,omitempty"`  // (nullable)
	Status        RequisitionStatus `json:"status" gorm:"type:varchar(50);default:'pending_approval_1'"`

	// Approval fields
	ApproverOneID   *int64     `json:"approver_one_id,omitempty" gorm:"index"` // ID of the first approver
	ApprovedOneAt   *time.Time `json:"approved_one_at,omitempty"`              // Timestamp of first approval
	ApproverTwoID   *int64     `json:"approver_two_id,omitempty" gorm:"index"` // ID of the second approver
	ApprovedTwoAt   *time.Time `json:"approved_two_at,omitempty"`              // Timestamp of second approval
	RejectionReason *string    `json:"rejection_reason,omitempty"`             // Reason if rejected

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Items []RequisitionItem `json:"items" gorm:"foreignKey:RequisitionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// User  User              `json:"user,omitempty" gorm:"foreignKey:UserID"` // Optional: to preload user details
}

// RequisitionItem corresponds to the RequisitionItems table
type RequisitionItem struct {
	ID                 int64    `json:"id" gorm:"primaryKey"`
	RequisitionID      int64    `json:"requisition_id" gorm:"index;not null"` // Foreign key to Requisition
	Description        string   `json:"description"`
	Quantity           float64  `json:"quantity"`
	Unit               string   `json:"unit"`
	EstimatedUnitPrice *float64 `json:"estimated_unit_price,omitempty"`
	FreightCost        *float64 `json:"freight_cost,omitempty"`
	InsuranceCost      *float64 `json:"insurance_cost,omitempty"`
	InstallationCost   *float64 `json:"installation_cost,omitempty"`
	AmrID              *int64   `json:"amr_id,omitempty" gorm:"index"` // Foreign key (nullable)
	Value              *float64 `json:"value,omitempty"`               // (nullable)

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
