import type { PageLoad } from './$types';
import type { Supplier } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, fetch }: LoadEvent) => {
  try {
    const id = params.id;
    const response = await fetch(`/api/suppliers/${id}`); // API endpoint to get a single supplier

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load supplier ${id}: ${response.status} ${errorText}`);
    }

    const supplier: Supplier = await response.json();
    return {
      supplier,
      error: null
    };
  } catch (error: any) {
    console.error('Error loading supplier:', error);
    return {
      supplier: null, // Or an empty/default supplier object
      error: error.message || 'An unknown error occurred'
    };
  }
};
