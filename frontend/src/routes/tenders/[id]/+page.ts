import type { PageLoad } from './$types';
import type { Tender } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { getAccessTokenSilently } from '$lib/authService';
import { isAuthenticated, user as userStore } from '$lib/store';
import { get } from 'svelte/store';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, fetch, depends }: LoadEvent) => {
	const id = params.id;
	depends(`app:tender:${id}`);

	if (!get(isAuthenticated)) {
		throw redirect(307, '/'); // Redirect to home or login if not authenticated
	}

	try {
		const token = await getAccessTokenSilently();
		if (!token) {
			return {
				tender: null,
				bids: null,
				error: 'Authentication token not available. Please log in again.',
				bidsError: null
			};
		}

		const response = await fetch(`${PUBLIC_API_BASE_URL}/api/tenders/${id}`, {
			headers: {
				'Authorization': `Bearer ${token}`
			}
		});

		if (!response.ok) {
			if (response.status === 401 || response.status === 403) {
				return {
					tender: null,
					bids: null,
					error: `You are not authorized to view tender ${id}. Your session might be invalid.`,
					bidsError: null
				};
			}
			const errorText = await response.text();
			throw new Error(`Failed to load tender ${id}: ${response.status} ${errorText}`);
		}

		const tender: Tender = await response.json();
		let bids = null;
		let bidsError = null;

		const currentUser = get(userStore);

		// If the user is a procurement officer, fetch bids for this tender
		if (currentUser && currentUser.role === 'procurement_officer') {
			try {
				const bidsResponse = await fetch(`${PUBLIC_API_BASE_URL}/api/tenders/${id}/bids`, {
					headers: {
						'Authorization': `Bearer ${token}`
					}
				});
				if (!bidsResponse.ok) {
					const errorText = await bidsResponse.text();
					bidsError = `Failed to load bids: ${bidsResponse.status} ${errorText}`;
				} else {
					bids = await bidsResponse.json();
				}
			} catch (err: any) {
				console.error(`Error loading bids for tender ${id}:`, err);
				bidsError = err.message || 'An unknown error occurred while fetching bids';
			}
		}

		return {
			tender,
			bids,
			error: null,
			bidsError
		};
	} catch (error: any) {
		console.error(`Error loading tender ${id}:`, error);
		return {
			tender: null,
			bids: null,
			error: error.message || `An unknown error occurred while fetching tender ${id}`,
			bidsError: null
		};
	}
};
