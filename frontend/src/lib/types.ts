// This file will store various TypeScript types used across the frontend.

export interface Supplier {
  id: number; // Changed from int64 to number for JS/TS
  name: string;
  contact_person?: string | null;
  email?: string | null;
  phone?: string | null;
  // Add other fields from your Go model if they are exposed to the frontend
}

export interface Asset {
  id: number; // from int64
  amr_id?: number | null; // from *int64
  emr_id?: number | null; // from *int64
  description: string;
  category?: string | null;
  status?: string | null;
  location?: string | null;
  purchase_date?: string | null; // from *time.Time (consider string for date inputs)
  purchase_price?: number | null; // from *float64
  supplier_id?: number | null; // from *int64
  created_by_user_id?: number | null; // from *int64
  // Add other fields from your Go model Asset if they are exposed to the frontend
  name?: string;
  // Add any other fields that might be used in the frontend
}

// Add other types here as needed, for example:
// export interface User {
//   id: number;
//   username: string;
//   email: string;
//   user_type: 'organization' | 'supplier';
//   created_at: string; // Or Date
//   updated_at: string; // Or Date
// }

// You can expand this file with types for Requisition, Tender, etc.

export interface Requisition {
  id: number;
  user_id: number;
  type: string;
  aac: number;
  material_group: string;
  exchange_rate: number;
  status: string;
  created_at: string;
  updated_at: string;
  // Add any additional fields that might be needed
  title?: string;
  description?: string;
  // Add other fields from your backend model as needed
  approver_one_id?: number | null;
  approved_one_at?: string | null;
  approver_two_id?: number | null;
  approved_two_at?: string | null;
  rejection_reason?: string | null;
  items: RequisitionItem[];
}

export interface RequisitionItem {
  id?: number; // Optional if it's a new item not yet saved
  description: string;
  quantity: number;
  unit: string;
  estimated_unit_price?: number | null;
  specification_text?: string | null;
  specificationSheetFile?: File | null; // For frontend handling of the file object
  specification_sheet_url?: string | null; // For backend to store/retrieve URL
  itemImageFile?: File | null; // For frontend handling
  item_image_url?: string | null; // For backend
  // Add any other existing item fields if necessary
}

export interface BidItem {
  id?: number; // Optional if it's a new item not yet saved
  bid_id?: number; // Will be set by backend when bid is created
  requisition_item_id?: number | null; // Link to the original RequisitionItem
  tender_item_id?: number; // Optional, to link to the original tender item if applicable (can be same as requisition_item_id or derived)
  description: string; // Could be pre-filled from tender item
  quantity: number; // Could be pre-filled from tender item
  unit: string; // Could be pre-filled from tender item
  offered_unit_price: number | null; // Supplier's price for this item
  specification_text?: string | null;
  specificationSheetFile?: File | null; // For frontend handling of the file object
  specification_sheet_url?: string | null; // For backend to store/retrieve URL
  itemImageFile?: File | null; // For frontend handling
  item_image_url?: string | null; // For backend
}


export interface Tender {
  id: number; // from int64
  requisition_id?: number | null; // from *int64
  title: string;
  description?: string | null;
  category?: string | null;
  budget?: number | null; // from *float64
  status?: string | null; // e.g., 'draft', 'published', 'evaluation', 'awarded', 'cancelled'
  published_date?: string | null; // from *time.Time
  closing_date?: string | null; // from *time.Time
  evaluation_method?: string | null;
  bid_opening_date?: string | null; // from *time.Time
  created_by_user_id?: number | null; // from *int64
  created_at?: string | null; // from *time.Time
  updated_at?: string | null; // from *time.Time
  requisition?: Requisition; // For preloaded requisition details
  bidders_invited_count?: number; // Number of suppliers who have placed a bid
  // Add other relevant fields from your Go model
}

export interface Bid {
  id: number; // from int64
  tender_id: number; // from int64
  supplier_id: number; // from int64, links to Supplier
  supplier_name?: string; // Denormalized for easier display
  submission_date: string; // from time.Time
  bid_amount?: number | null; // from *float64
  currency?: string | null; // e.g., 'USD', 'EUR'
  status: string; // e.g., 'submitted', 'under_evaluation', 'shortlisted', 'rejected', 'awarded'
  notes?: string | null;
  validity_period_days?: number | null; // from *int
  items?: BidItem[]; // Array of bid items with detailed specifications
  created_at?: string | null; // from *time.Time
  updated_at?: string | null; // from *time.Time
  // Potentially an array of bid documents
  // documents?: Array<{ name: string; url: string; type: string }>;
}

export interface Evaluation {
  id: number; // from int64
  bid_id: number; // from int64
  evaluator_id: number; // from int64 (links to User)
  score?: number | null; // from *float64, e.g., overall score or could be more complex
  comments?: string | null;
  evaluation_date: string; // from time.Time
  status?: string | null; // e.g., 'pending_review', 'approved', 'rejected_for_award'
  // Potentially a more structured criteria object or array
  // criteria?: Array<{ criterion_id: number; score: number; comment?: string }>;
  created_at?: string | null; // from *time.Time
  updated_at?: string | null; // from *time.Time
}

export interface User {
  id: number;
  username: string;
  email: string;
  user_type: 'organization' | 'supplier';
  created_at: string; // Or Date
  updated_at: string; // Or Date
}
