<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import type { Tender } from '$lib/types';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
	import { getAccessTokenSilently } from '$lib/authService';
	import { user } from '$lib/store'; // Import the user store

	// Props from +page.ts
	export let data: PageData;

	// Reactive declarations for tender data and errors
	$: tender = data.tender;
	$: error = data.error; // Error from loading tender itself
	$: bids = data.bids; // Bids data for POs
	$: bidsError = data.bidsError; // Error from loading bids

	// Component state
	let editableTender: Partial<Tender> = {}; // For form binding
	let isLoading: boolean = true;
	let editMode: boolean = false;

	// Bid submission form variables
	let bidAmount: number | null = null;
	let bidProposalUrl: string = ''; // Assuming a URL for a proposal document
	let bidSubmissionError: string | null = null;
	let bidSubmissionSuccess: string | null = null;
	let isSubmittingBid: boolean = false;

	let successMessage: string | null = null; // For general success messages like delete/update
	let errorMessage: string | null = null;

	$: {
		if (data && data.tender) {
			const isNewTenderOrDataChanged = !tender || tender.id !== data.tender.id || tender.updated_at !== data.tender.updated_at;
			tender = data.tender;
			
			// Only re-initialize editableTender if not in edit mode OR if the underlying tender data has actually changed
			// (e.g., after a save, or if a different tender was loaded)
			if (tender && (!editMode || isNewTenderOrDataChanged)) { 
				editableTender = { ...tender };
				if (editableTender.published_date) {
					editableTender.published_date = formatDateForDateInput(editableTender.published_date);
				}
				if (editableTender.closing_date) {
					editableTender.closing_date = formatDateForInput(editableTender.closing_date);
				}
				// Add bid_opening_date formatting as well
				if (editableTender.bid_opening_date) {
					editableTender.bid_opening_date = formatDateForInput(editableTender.bid_opening_date);
				}
			}
			error = data.error; // Ensure data.error is assigned to the local error variable
			isLoading = false;
		} else if (data && data.error) {
			error = data.error;
			tender = null;
			isLoading = false;
		}
	}

	// Reactive conditions for UI rendering based on user role and tender state
	$: isSupplier = $user && $user.role === 'supplier';
	$: isTenderOpenForBidding = tender && (tender.status === 'open' || tender.status === 'published') && tender.closing_date && new Date(tender.closing_date) > new Date();
	$: showBidForm = !editMode && isSupplier && isTenderOpenForBidding && tender?.status !== 'closed' && tender?.status !== 'cancelled' && tender?.status !== 'awarded';
	$: canEditTender = !editMode && $user && ($user.role === 'procurement_officer' || $user.role === 'admin');
	$: canViewBids = !editMode && $user && ($user.role === 'procurement_officer' || $user.role === 'admin');

	// Edit mode state

	function toggleEditMode() {
		editMode = !editMode;
		successMessage = null;
		errorMessage = null;
		if (editMode && tender) {
			// Re-initialize editableTender from current tender state when entering edit mode
			editableTender = { ...tender }; // Ensure all fields from tender are copied, especially those not on the form but part of the Tender type

				// Format dates for their respective input fields
				if (editableTender.published_date) {
					editableTender.published_date = formatDateForDateInput(editableTender.published_date);
				}
				if (editableTender.closing_date) {
					editableTender.closing_date = formatDateForInput(editableTender.closing_date);
				}
				if (editableTender.bid_opening_date) {
					editableTender.bid_opening_date = formatDateForInput(editableTender.bid_opening_date);
				}
		} 
	}

	async function handleSave() {
		if (!tender || !tender.id) return;
		isLoading = true;
		errorMessage = null;
		successMessage = null;

		// Construct the payload with only the fields intended for update
		const payload: Partial<Tender> = {
			title: editableTender.title,
			description: editableTender.description === '' ? null : editableTender.description, // Send null if empty
			category: editableTender.category === '' ? null : editableTender.category,
			status: editableTender.status === '' ? null : editableTender.status,
			evaluation_method: editableTender.evaluation_method === '' ? null : editableTender.evaluation_method,
			// Convert budget to number, or null if empty/invalid
			budget: editableTender.budget !== undefined && editableTender.budget !== null && String(editableTender.budget).trim() !== '' ? Number(editableTender.budget) : null,
			// Convert dates to full ISO strings if they exist, otherwise null
			published_date: editableTender.published_date ? new Date(editableTender.published_date).toISOString() : null,
			closing_date: editableTender.closing_date ? new Date(editableTender.closing_date).toISOString() : null,
			bid_opening_date: editableTender.bid_opening_date ? new Date(editableTender.bid_opening_date).toISOString() : null,
			requisition_id: editableTender.requisition_id ? Number(editableTender.requisition_id) : null
		};

		// Remove null fields from payload if backend prefers omitted fields over nulls for some optional fields
		// For now, sending null for empty optional fields is generally fine for GORM *type fields.

		try {
			const token = await getAccessTokenSilently();
			if (!token) {
				errorMessage = 'Authentication token not available. Please log in again.';
				isLoading = false;
				return;
			}

			const response = await fetch(`${PUBLIC_API_BASE_URL}/api/tenders/${tender.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${token}`
				},
				body: JSON.stringify(payload) // Send the cleaned payload
			});

			if (!response.ok) {
				// Try to parse backend error message, otherwise use a generic one
				const errorData = await response.json().catch(() => ({ error: `HTTP error ${response.status}. Failed to parse error response.` }));
				throw new Error(errorData.error || `HTTP error ${response.status}. An unexpected error occurred.`);
			}

			const updatedTenderFromServer: Tender = await response.json();
			tender = updatedTenderFromServer; // Update the main tender object with new data from server
			
			// Re-initialize editableTender for display/next edit, formatting dates for input fields
			editableTender = { ...tender }; 
			if (editableTender.published_date) {
				editableTender.published_date = formatDateForDateInput(editableTender.published_date);
			}
			if (editableTender.closing_date) {
				editableTender.closing_date = formatDateForInput(editableTender.closing_date);
			}
            if (editableTender.bid_opening_date) {
                editableTender.bid_opening_date = formatDateForInput(editableTender.bid_opening_date);
            }

			successMessage = 'Tender updated successfully!';
			editMode = false;
		} catch (err: any) {
			console.error('Save error:', err);
			errorMessage = err.message || 'An unexpected error occurred.';
		} finally {
			isLoading = false;
		}
	}

	async function handleSubmitBid() {
		if (!tender || !tender.id || !bidAmount) {
			bidSubmissionError = 'Bid amount is required.';
			return;
		}
		isSubmittingBid = true;
		bidSubmissionError = null;
		bidSubmissionSuccess = null;

		try {
			const token = await getAccessTokenSilently();
			if (!token) {
				bidSubmissionError = 'Authentication token not available. Please log in again.';
				isSubmittingBid = false;
				return;
			}

			const bidPayload = {
				bid_amount: bidAmount,
				// Add other bid fields here if necessary, e.g.:
				proposal_document_url: bidProposalUrl.trim() === '' ? null : bidProposalUrl.trim(),
				// notes: bidNotes,
			};

			const response = await fetch(`${PUBLIC_API_BASE_URL}/api/tenders/${tender.id}/bids`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${token}`
				},
				body: JSON.stringify(bidPayload)
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: `HTTP error ${response.status}` }));
				throw new Error(errorData.message || errorData.error || `Failed to submit bid: ${response.status}`);
			}

			const createdBid = await response.json();
			bidSubmissionSuccess = `Bid submitted successfully! Bid ID: ${createdBid.id}`;
			// Optionally, clear the form or redirect
			bidAmount = null;
			bidProposalUrl = '';
			// Consider invalidating tender data if bids list is shown on this page
			// import { invalidate } from '$app/navigation';
			// invalidate((url) => url.pathname === `/api/tenders/${tender.id}/bids`);
		} catch (err: any) {
			console.error('Bid submission error:', err);
			bidSubmissionError = err.message || 'An unexpected error occurred while submitting your bid.';
		} finally {
			isSubmittingBid = false;
		}
	}

	async function handleDelete() {
		if (!tender || !tender.id) return;
		if (!confirm('Are you sure you want to delete this tender?')) {
			return;
		}
		isLoading = true;
		errorMessage = null;

		try {
			const token = await getAccessTokenSilently();
			if (!token) {
				errorMessage = 'Authentication token not available. Please log in again.';
				isLoading = false;
				return;
			}

			const response = await fetch(`${PUBLIC_API_BASE_URL}/api/tenders/${tender.id}`, {
				method: 'DELETE',
				headers: {
					'Authorization': `Bearer ${token}`
				}
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to delete tender. The server returned an unexpected response.' }));
				throw new Error(errorData.message || `HTTP error ${response.status}`);
			}
			
			successMessage = 'Tender deleted successfully!';
			setTimeout(() => {
				goto('/tenders');
			}, 1500);
		} catch (err: any) {
			console.error('Delete error:', err);
			errorMessage = err.message || 'An unexpected error occurred during deletion.';
		} finally {
			isLoading = false;
		}
	}

	function formatDate(dateString: string | null | undefined) {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' });
	}

	function formatDateForInput(dateString: string | null | undefined): string {
        if (!dateString) return '';
        const d = new Date(dateString);
        // Format to YYYY-MM-DDTHH:mm for datetime-local input, adjusting for timezone
        return new Date(d.getTime() - d.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
    }

    function formatDateForDateInput(dateString: string | null | undefined): string {
        if (!dateString) return '';
        return new Date(dateString).toISOString().split('T')[0];
    }

</script>

<svelte:head>
	<title>{tender ? `Tender: ${tender.title}` : 'Tender Details'} - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	{#if isLoading}
		<p class="text-center text-xl">Loading tender details...</p>
	{:else if error}
		<div class="alert alert-error">
			<p>Error loading tender: {error}</p>
			<a href="/tenders" class="btn btn-sm btn-outline mt-2">Back to Tenders</a>
		</div>
	{:else if tender}
		<div class="bg-white shadow-lg rounded-lg p-6 md:p-8">
			{#if successMessage}
				<div class="alert alert-success mb-4">
					<p>{successMessage}</p>
				</div>
			{/if}
			{#if errorMessage && !editMode} <!-- Show general errors when not in edit mode -->
				<div class="alert alert-error mb-4">
					<p>{errorMessage}</p>
				</div>
			{/if}

			{#if editMode}
				<!-- Edit Form -->
				<form on:submit|preventDefault={handleSave} class="space-y-6">
					<h2 class="text-2xl font-semibold text-gray-800 mb-4">Edit Tender: {tender.title}</h2>
					
					{#if errorMessage && editMode} <!-- Show form-specific errors when in edit mode -->
						<div class="alert alert-error mb-4">
							<p>{errorMessage}</p>
						</div>
					{/if}

					<div>
						<label for="edit-title" class="block text-sm font-medium text-gray-700">Title <span class="text-red-500">*</span></label>
						<input type="text" id="edit-title" bind:value={editableTender.title} required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
					</div>
					<div>
						<label for="edit-description" class="block text-sm font-medium text-gray-700">Description</label>
						<textarea id="edit-description" rows="4" bind:value={editableTender.description} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"></textarea>
					</div>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<div>
							<label for="edit-category" class="block text-sm font-medium text-gray-700">Category</label>
							<select id="edit-category" bind:value={editableTender.category} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
								<option value="goods">Goods</option>
								<option value="services">Services</option>
								<option value="works">Works</option>
								<option value="consultancy">Consultancy</option>
							</select>
						</div>
						<div>
							<label for="edit-budget" class="block text-sm font-medium text-gray-700">Budget</label>
							<input type="number" step="0.01" id="edit-budget" bind:value={editableTender.budget} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						</div>
					</div>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<div>
							<label for="edit-issue-date" class="block text-sm font-medium text-gray-700">Published Date</label>
							<input type="date" id="edit-issue-date" bind:value={editableTender.published_date} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						</div>
						<div>
							<label for="edit-closing-date" class="block text-sm font-medium text-gray-700">Closing Date <span class="text-red-500">*</span></label>
							<input type="datetime-local" id="edit-closing-date" bind:value={editableTender.closing_date} required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						</div>
					</div>
					<div>
						<label for="edit-status" class="block text-sm font-medium text-gray-700">Status</label>
						<select id="edit-status" bind:value={editableTender.status} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
							<option value="draft">Draft</option>
							<option value="published">Published</option>
							<option value="open">Open</option> 
							<option value="evaluation">Evaluation</option>
							<option value="awarded">Awarded</option>
							<option value="closed">Closed</option>
							<option value="cancelled">Cancelled</option>
						</select>
					</div>
					<div class="flex justify-end space-x-3 pt-4 border-t mt-6">
						<button type="button" on:click={toggleEditMode} class="btn btn-ghost" disabled={isLoading}>Cancel</button>
						<button type="submit" class="btn btn-primary" disabled={isLoading}>
							{#if isLoading}Saving...{:else}Save Changes{/if}
						</button>
					</div>
				</form>
			{:else}
				<!-- View Details -->
				<div class="mb-6 pb-4 border-b border-gray-200 flex justify-between items-start">
					<div>
						<h1 class="text-3xl font-semibold text-gray-800">{tender.title}</h1>
						<p class="text-sm text-gray-500">
							ID: {tender.id} | Published: {formatDate(tender.published_date)} | Closes: {formatDate(tender.closing_date)}
						</p>
					</div>
					{#if canEditTender}
						<div class="flex justify-end space-x-3 mb-6">
							<button on:click={toggleEditMode} class="btn btn-outline btn-primary">
								Edit Tender
							</button>
							<button on:click={handleDelete} class="btn btn-outline btn-error" disabled={isLoading}>
								{#if isLoading}Deleting...{:else}Delete Tender{/if}
							</button>
						</div>
					{/if}
				</div>

				<!-- Tender Details Section -->
				<div class="mt-6 py-4 space-y-4 border-t border-b border-gray-200 mb-6">
					<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-x-4 gap-y-6">
						<div>
							<p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Status</p>
							<p class="mt-1 text-sm text-gray-900"><span class="badge badge-lg {
								tender.status === 'open' ? 'badge-success' :
								tender.status === 'published' ? 'badge-info' :
								tender.status === 'closed' ? 'badge-error' :
								tender.status === 'draft' ? 'badge-ghost' :
								tender.status === 'evaluation' ? 'badge-warning' :
								tender.status === 'awarded' ? 'badge-accent' : 
								tender.status === 'cancelled' ? 'badge-outline badge-error' :
								'badge-secondary'
							}">{tender.status || 'N/A'}</span></p>
						</div>
						<div>
							<p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Category</p>
							<p class="mt-1 text-sm text-gray-900 capitalize">{tender.category || 'N/A'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Budget</p>
							<p class="mt-1 text-sm text-gray-900">{tender.budget != null ? `$${Number(tender.budget).toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}` : 'N/A'}</p>
						</div>
						{#if tender.requisition_id}
							<div>
								<p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Linked Requisition</p>
								<a href={`/requisitions/${tender.requisition_id}`} class="mt-1 text-sm text-blue-600 hover:underline">REQ-{String(tender.requisition_id).padStart(3, '0')}</a>
							</div>
						{/if}
					</div>

					<div>
						<p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Description</p>
						<div class="mt-1 text-sm text-gray-900 prose prose-sm max-w-none whitespace-pre-wrap break-words">{tender.description || 'No description provided.'}</div>
					</div>
					
					<div class="text-xs text-gray-500 pt-3 mt-3 border-t border-gray-100">
						<p>Created by User ID: {tender.created_by_user_id || 'N/A'} on {formatDate(tender.created_at)}</p>
						<p>Last updated on {formatDate(tender.updated_at)}</p>
					</div>
				</div>
				<!-- End Tender Details Section -->

				{#if showBidForm}
					<!-- Bid Form -->
					<form on:submit|preventDefault={handleSubmitBid} class="space-y-6">
						<h2 class="text-2xl font-semibold text-gray-800 mb-4">Submit Bid for {tender.title}</h2>
						
						{#if bidSubmissionError}
							<div class="alert alert-error mb-4">
								<p>{bidSubmissionError}</p>
							</div>
						{/if}
						{#if bidSubmissionSuccess}
							<div class="alert alert-success mb-4">
								<p>{bidSubmissionSuccess}</p>
							</div>
						{/if}

						<div>
							<label for="bid-amount" class="block text-sm font-medium text-gray-700">Bid Amount <span class="text-red-500">*</span></label>
							<input type="number" step="0.01" id="bid-amount" bind:value={bidAmount} required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						</div>
						<div>
							<label for="bid-proposal-url" class="block text-sm font-medium text-gray-700">Proposal Document URL (Optional)</label>
							<input type="url" id="bid-proposal-url" bind:value={bidProposalUrl} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="https://example.com/proposal.pdf">
						</div>

						<div class="flex justify-end space-x-3 pt-4 border-t mt-6">
							<button type="submit" class="btn btn-primary" disabled={isSubmittingBid}>
								{#if isSubmittingBid}Submitting...{:else}Submit Bid{/if}
							</button>
						</div>
					</form>
				{/if}

				<!-- Display Bids for Procurement Officers/Admins -->
				{#if canViewBids && !editMode}
					<div class="mt-10 pt-6 border-t">
						<h3 class="text-xl font-semibold text-gray-800 mb-4">Submitted Bids</h3>
						{#if bidsError}
							<div class="alert alert-error">
								<p>Error loading bids: {bidsError}</p>
							</div>
						{:else if bids && bids.length > 0}
							<div class="overflow-x-auto">
								<table class="table w-full table-zebra table-compact">
									<thead>
										<tr>
											<th>Bid ID</th>
											<th>Supplier ID</th>
											<th>Amount</th>
											<th>Proposal URL</th>
											<th>Submitted At</th>
											<th>Status</th>
										</tr>
									</thead>
									<tbody>
										{#each bids as bid (bid.id)}
											<tr>
												<td>{bid.id}</td>
												<td>{bid.supplier_id}</td>
												<td>{bid.bid_amount?.toFixed(2) || 'N/A'}</td>
												<td>
													{#if bid.proposal_document_url}
														<a href={bid.proposal_document_url} target="_blank" rel="noopener noreferrer" class="link link-primary">View Proposal</a>
													{:else}
														N/A
													{/if}
												</td>
												<td>{formatDate(bid.created_at)}</td>
												<td><span class="badge badge-sm {bid.status === 'submitted' ? 'badge-info' : bid.status === 'awarded' ? 'badge-success' : bid.status === 'rejected' ? 'badge-error' : 'badge-ghost'}">{bid.status || 'N/A'}</span></td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{:else}
							<p class="text-gray-600 italic">No bids have been submitted for this tender yet.</p>
						{/if}
					</div>
				{/if}

				<div class="text-center mt-4">
					<a href="/tenders" class="btn btn-primary"> Back to Tenders List</a>
				</div>
			{/if}
		</div>
	{:else}
		<p class="text-center text-xl text-gray-500">Tender not found.</p>
		<div class="text-center mt-4">
			<a href="/tenders" class="btn btn-primary">Go to Tenders List</a>
		</div>
	{/if}
</div>
