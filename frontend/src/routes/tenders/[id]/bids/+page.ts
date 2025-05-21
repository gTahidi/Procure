import type { PageLoad } from './$types';
import type { Bid, Tender } from '$lib/types';

export const load: PageLoad = async ({ fetch, params, parent }) => {
	const tenderId = params.id;
	let tenderDetails: Tender | null = null;

	// Attempt to get tender details from parent layout if available
	try {
		const parentData = await parent() as { tender?: Tender };
		if (parentData && parentData.tender) {
			tenderDetails = parentData.tender;
		} else {
			// Fallback to fetching tender details directly if not in parent
			const tenderRes = await fetch(`/api/tenders/${tenderId}`);
			if (tenderRes.ok) {
				tenderDetails = await tenderRes.json();
			}
		}
	} catch (e) {
		console.error('Error loading tender details for bids page:', e);
		// Continue to try fetching bids even if tender details fail,
		// but it's better to have tender context.
	}

	try {
		const response = await fetch(`/api/tenders/${tenderId}/bids`);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to load bids. Server returned an error.'}));
			return {
				tender: tenderDetails,
				bids: [],
				error: errorData.message || `HTTP error ${response.status}`,
				tenderId
			};
		}

		const bids: Bid[] = await response.json();
		return {
			tender: tenderDetails,
			bids,
			error: null,
			tenderId
		};
	} catch (e: any) {
		console.error('Fetch bids error:', e);
		return {
			tender: tenderDetails,
			bids: [],
			error: e.message || 'An unexpected error occurred while fetching bids.',
			tenderId
		};
	}
};
