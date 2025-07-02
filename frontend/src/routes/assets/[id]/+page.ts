import type { PageLoad } from './$types';
import type { Asset } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';
import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ params, fetch }: LoadEvent) => {
  try {
    const id = params.id;
        const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/assets/${id}`); // API endpoint to get a single asset

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load asset ${id}: ${response.status} ${errorText}`);
    }

    const asset: Asset = await response.json();
    return {
      asset,
      error: null
    };
  } catch (error: any) {
    console.error('Error loading asset:', error);
    return {
      asset: null, 
      error: error.message || 'An unknown error occurred'
    };
  }
};
