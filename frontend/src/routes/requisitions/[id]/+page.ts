import { error, redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { isAuthenticated } from '$lib/store';
import { get } from 'svelte/store';
import { getAccessTokenSilently } from '$lib/authService';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageLoad = async ({ params, fetch }) => {
	if (!get(isAuthenticated)) {
		throw redirect(307, '/'); // Redirect to home if not authenticated
	}

	const requisitionId = params.id;
	if (!requisitionId) {
		throw error(400, 'Requisition ID is missing');
	}

	const token = await getAccessTokenSilently();
	if (!token) {
		throw error(401, 'Not authorized. Please log in again.');
	}

	try {
		const response = await fetch(`${PUBLIC_API_BASE_URL}/api/requisitions/${requisitionId}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			}
		});

		if (response.status === 401 || response.status === 403) {
			throw error(response.status, 'Authentication error. Your session might be invalid or you do not have permission.');
		}
		if (response.status === 404) {
			throw error(404, 'Requisition not found.');
		}
		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to load requisition details.' }));
			throw error(response.status, errorData.message || 'Failed to load requisition details.');
		}

		const requisition = await response.json();
		return {
			requisition
		};
	} catch (e: any) {
		console.error('Error loading requisition details:', e);
		if (e.status) { // If it's an error thrown by SvelteKit's error() helper
			throw e;
		}
		throw error(500, e.message || 'An unexpected error occurred while fetching requisition details.');
	}
};
