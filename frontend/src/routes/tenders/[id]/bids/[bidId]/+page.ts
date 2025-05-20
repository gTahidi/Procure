import type { PageLoad } from './$types';
import type { Bid, Tender, Evaluation } from '$lib/types';

export const load: PageLoad = async ({ fetch, params, parent }) => {
	const tenderId = params.id;
	const bidId = params.bidId;
	let tenderDetails: Tender | null = null;
	let bidDetails: Bid | null = null;
	let evaluationDetails: Evaluation | null = null;

	// Attempt to get tender details from parent layout if available
	// This helps in having tender context on the bid detail page
	try {
		const parentData = await parent() as { tender?: Tender }; // Type assertion
		// A more specific parent() call might be needed if layouts are nested deeper or data isn't passed down as expected.
		// For now, we assume a relevant parent might provide tender data.
		if (parentData && parentData.tender) { 
			tenderDetails = parentData.tender;
		}
		// If not found via parent, could fetch /api/tenders/${tenderId} as a fallback
	} catch (e) {
		console.warn('Could not load tender details from parent for bid page:', e);
	}

	try {
		const bidResponse = await fetch(`/api/bids/${bidId}`); // Assuming endpoint like /api/bids/:bidId

		if (!bidResponse.ok) {
			const errorData = await bidResponse.json().catch(() => ({ message: 'Failed to load bid details.' }));
			return {
				tender: tenderDetails,
				bid: null,
				evaluation: null,
				error: errorData.message || `HTTP error ${bidResponse.status}`,
				tenderId,
				bidId
			};
		}

		bidDetails = await bidResponse.json();

		// If tender details were not found from parent, and bid has tender_id, fetch tender details
		if (!tenderDetails && bidDetails && bidDetails.tender_id && bidDetails.tender_id.toString() === tenderId) {
			const tenderRes = await fetch(`/api/tenders/${bidDetails.tender_id}`);
			if (tenderRes.ok) {
				tenderDetails = await tenderRes.json();
			}
		}

		// Fetch evaluation details for the bid
		if (bidDetails) {
			const evalResponse = await fetch(`/api/bids/${bidId}/evaluations`); // Assuming this endpoint
			if (evalResponse.ok) {
				evaluationDetails = await evalResponse.json();
			} else if (evalResponse.status !== 404) { // Don't error for 404 (no evaluation yet)
				console.warn(`Failed to load evaluation for bid ${bidId}: ${evalResponse.status}`);
			}
		}

		return {
			tender: tenderDetails,
			bid: bidDetails,
			evaluation: evaluationDetails,
			error: null,
			tenderId,
			bidId
		};
	} catch (e: any) {
		console.error('Fetch bid details error:', e);
		return {
			tender: tenderDetails,
			bid: bidDetails, // ensure bidDetails is passed even on error if fetched
			evaluation: null,
			error: e.message || 'An unexpected error occurred while fetching bid details.',
			tenderId,
			bidId
		};
	}
};
