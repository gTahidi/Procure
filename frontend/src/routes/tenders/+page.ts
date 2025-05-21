import type { PageLoad } from './$types';
import type { Tender } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ fetch }: LoadEvent) => {
  try {
    const response = await fetch(`${PUBLIC_API_BASE_URL}/tenders`); // API endpoint for tenders
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load tenders: ${response.status} ${errorText}`);
    }
    const tenders: Tender[] = await response.json();
    return {
      tenders,
      error: null
    };
  } catch (error: any) {
    console.error('Error loading tenders:', error);
    return {
      tenders: [],
      error: error.message || 'An unknown error occurred'
    };
  }
};
