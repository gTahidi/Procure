import type { PageLoad } from './$types';
import type { Tender } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';
import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';
import { getAccessTokenSilently } from '$lib/authService';
import { isAuthenticated } from '$lib/store';
import { get } from 'svelte/store';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async ({ fetch, depends }: LoadEvent) => {
	depends('app:tenders'); // Ensure data reloads when this dependency invalidates

	if (!get(isAuthenticated)) {
		throw redirect(307, '/'); // Redirect to home or login if not authenticated
	}

	try {
		const token = await getAccessTokenSilently();
		if (!token) {
			return {
				tenders: [],
				error: 'Authentication token not available. Please log in again.'
			};
		}

		const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders`, {
			headers: {
				'Authorization': `Bearer ${token}`
			}
		}); 

		if (!response.ok) {
			if (response.status === 401 || response.status === 403) {
				return {
					tenders: [],
					error: 'You are not authorized to view tenders. Your session might be invalid.'
				};
			}
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
			error: error.message || 'An unknown error occurred while fetching tenders'
		};
	}
};
