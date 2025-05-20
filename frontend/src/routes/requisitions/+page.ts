import type { PageLoad, PageLoadEvent } from './$types';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export interface BackendRequisition {
  id: number;          // Matches 'ID' from backend JSON
  user_id: number;     // Matches 'UserID'
  type: string;
  aac: number;
  material_group: string;
  exchange_rate: number;
  status: string;
  created_at: string;  // ISO 8601 date string e.g., "2024-05-20T15:04:05Z"
  updated_at: string;
}

export interface FrontendRequisition {
  id: string; // Displayed as REQ-XXX
  title: string;
  requester: string; // Will be UserID for now
  type: string;
  status: string;
  creationDate: string; // Formatted date
  detailLink: string;
}

export const load: PageLoad = async (event: PageLoadEvent) => {
  try {
    const response = await event.fetch(`${PUBLIC_API_BASE_URL}/requisitions`); // Ensure backend is running and accessible
    if (!response.ok) {
      throw new Error(`HTTP error ${response.status} while fetching requisitions`);
    }
    const backendRequisitions: BackendRequisition[] = await response.json();

    const requisitions: FrontendRequisition[] = backendRequisitions.map(req => ({
      id: `REQ-${String(req.id).padStart(3, '0')}`,
      // Title: Currently not available from this backend endpoint. Using a placeholder.
      title: `Requisition - ${String(req.id).padStart(3, '0')}`,
      requester: `User ID: ${req.user_id}`,
      type: req.type,
      status: req.status,
      creationDate: new Date(req.created_at).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric'
      }),
      detailLink: `/requisitions/REQ-${String(req.id).padStart(3, '0')}`
    }));

    return {
      requisitions
    };
  } catch (error) {
    console.error('Failed to load requisitions:', error);
    return {
      requisitions: [],
      error: error instanceof Error ? error.message : 'Unknown error loading requisitions'
    };
  }
};
