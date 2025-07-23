import type { PageLoad, PageLoadEvent } from './$types';
import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';
import { getAccessToken } from '$lib/authService';
import { isAuthenticated } from '$lib/store';
import { get } from 'svelte/store';
import { redirect } from '@sveltejs/kit';

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
  const authenticated = get(isAuthenticated);

  if (!authenticated) {
    console.log('[+page.ts /requisitions] User not authenticated. Redirecting to login.');
    throw redirect(307, '/'); // Redirect to home, which should trigger login if needed via +layout.svelte
  }

  let token: string | undefined;
  try {
    token = getAccessToken() || undefined;
    if (!token) {
      console.error('[+page.ts /requisitions] Authenticated, but failed to retrieve access token.');
      return {
        requisitions: [],
        error: 'Failed to obtain authentication token. Please try logging out and back in.'
      };
    }
  } catch (e) {
    console.error('[+page.ts /requisitions] Error calling getAccessTokenSilently:', e);
    return {
      requisitions: [],
      error: 'Error retrieving authentication token. Please try logging out and back in.'
    };
  }

  try {
    const fetchOptions: RequestInit = {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    };

    const response = await event.fetch(`${PUBLIC_VITE_API_BASE_URL}/api/requisitions`, fetchOptions);
    if (!response.ok) {
      if (response.status === 401 || response.status === 403) {
        console.error(`[+page.ts /requisitions] Auth error from backend: ${response.status}`);
        return {
          requisitions: [],
          error: `Authentication error from server (${response.status}). Your session might be invalid.`
        };
      }
      throw new Error(`HTTP error ${response.status} while fetching requisitions`);
    }
    const backendRequisitions: BackendRequisition[] = await response.json();

    const requisitions: FrontendRequisition[] = backendRequisitions.map(req => ({
      id: `REQ-${String(req.id).padStart(3, '0')}`,
      title: `Requisition - ${String(req.id).padStart(3, '0')}`,
      requester: `User ID: ${req.user_id}`,
      type: req.type,
      status: req.status,
      creationDate: new Date(req.created_at).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric'
      }),
      detailLink: `/requisitions/${req.id}`
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
