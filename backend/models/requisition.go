package models

import "time"

// Requisition corresponds to the Requisitions table
type Requisition struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Type          string    `json:"type"` // 'goods', 'services', 'fixed_asset'
	AAC           *string   `json:"aac,omitempty"` // 'A', 'F', 'P' (nullable)
	MaterialGroup *string   `json:"material_group,omitempty"` // (nullable)
	ExchangeRate  *float64  `json:"exchange_rate,omitempty"` // (nullable)
	Status        string    `json:"status"` // Default 'pending'
	CreatedAt     time.Time `json:"created_at"`
	Items         []RequisitionItem `json:"items"` // For including items in a single request/response
}

// RequisitionItem corresponds to the RequisitionItems table
type RequisitionItem struct {
	ID                 int64    `json:"id"`
	RequisitionID      int64    `json:"requisition_id,omitempty"` // Foreign key, omitempty if not set yet (e.g. before main req is saved)
	Description        string   `json:"description"`
	Quantity           float64  `json:"quantity"`
	Unit               string   `json:"unit"`
	EstimatedUnitPrice *float64 `json:"estimated_unit_price,omitempty"`
	FreightCost        *float64 `json:"freight_cost,omitempty"`      // Default 0
	InsuranceCost      *float64 `json:"insurance_cost,omitempty"`    // Default 0
	InstallationCost   *float64 `json:"installation_cost,omitempty"` // Default 0
	AmrID              *int64   `json:"amr_id,omitempty"`             // Foreign key (nullable)
	Value              *float64 `json:"value,omitempty"`               // (nullable)
}
