import type { PageLoad } from './$types';
import type { Supplier } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export const load: PageLoad = async ({ fetch }: LoadEvent) => {
  try {
    const response = await fetch('/api/suppliers'); // Ensure this API endpoint exists on your Go backend
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load suppliers: ${response.status} ${errorText}`);
    }
    const suppliers: Supplier[] = await response.json();
    return {
      suppliers,
      error: null
    };
  } catch (error: any) {
    console.error('Error loading suppliers:', error);
    return {
      suppliers: [],
      error: error.message || 'An unknown error occurred'
    };
  }
};
