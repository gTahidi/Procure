import type { PageLoad } from './$types';
import type { Tender, Bid } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';
import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';
import { getAccessTokenSilently } from '$lib/authService';
import { isAuthenticated, user as userStore } from '$lib/store';
import { get } from 'svelte/store';
import { redirect, error as svelteKitError } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, fetch, depends }: LoadEvent) => {
	const id = params.id;
	depends(`app:tender:${id}:evaluation`);

	if (!get(isAuthenticated)) {
		throw redirect(307, '/');
	}

	const currentUser = get(userStore);
	if (currentUser?.role !== 'procurement_officer') {
		throw svelteKitError(403, 'You are not authorized to view this page.');
	}

	try {
		const token = await getAccessTokenSilently();
		if (!token) {
			throw svelteKitError(401, 'Authentication token not available. Please log in again.');
		}

		// Fetch both tender and bids details in parallel
		const [tenderResponse, bidsResponse] = await Promise.all([
			fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders/${id}`, {
				headers: { 'Authorization': `Bearer ${token}` }
			}),
			fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders/${id}/bids`, {
				headers: { 'Authorization': `Bearer ${token}` }
			})
		]);

		if (!tenderResponse.ok) {
			throw svelteKitError(tenderResponse.status, `Failed to load tender: ${await tenderResponse.text()}`);
		}

		if (!bidsResponse.ok) {
			throw svelteKitError(bidsResponse.status, `Failed to load bids: ${await bidsResponse.text()}`);
		}

		const tender: Tender = await tenderResponse.json();
		const bids: Bid[] = await bidsResponse.json();

		return {
			tender,
			bids
		};
	} catch (error: any) {
		console.error(`Error loading tender evaluation page for tender ${id}:`, error);
		// Re-throw SvelteKit errors, otherwise create a new one
		if (error.status) throw error;
		throw svelteKitError(500, error.message || 'An unknown error occurred.');
	}
};
