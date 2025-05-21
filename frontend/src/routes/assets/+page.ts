import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { LoadEvent } from '@sveltejs/kit';
import type { Asset } from '$lib/types';

interface PageData {
  assets: Asset[];
  error: string | null;
}

export const load: PageLoad<PageData> = async ({ fetch }: LoadEvent) => {
  try {
    // For now, return an empty array since we don't have an assets endpoint
    // Once the backend implements the /api/assets endpoint, you can uncomment the fetch below
    // and remove the mock data
    /*
    const response = await fetch('/api/assets');
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load assets: ${response.status} ${errorText}`);
    }
    const assets: Asset[] = await response.json();
    */
    
    // For now, return an empty array of assets
    const assets: Asset[] = [];
    
    return {
      assets,
      error: null
    };
  } catch (err: any) {
    console.error('Error in assets page:', err);
    return {
      assets: [],
      error: 'Failed to load assets. Please try again later.'
    };
  }
};
