import type { PageLoad } from './$types';
import type { Tender } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, fetch }: LoadEvent) => {
  try {
    const id = params.id;
    const response = await fetch(`/api/tenders/${id}`); // API endpoint to get a single tender

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load tender ${id}: ${response.status} ${errorText}`);
    }

    const tender: Tender = await response.json();
    return {
      tender,
      error: null
    };
  } catch (error: any) {
    console.error('Error loading tender:', error);
    return {
      tender: null, 
      error: error.message || 'An unknown error occurred'
    };
  }
};
