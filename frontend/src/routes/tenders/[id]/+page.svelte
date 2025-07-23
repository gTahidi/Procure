<script lang="ts">
	import { page } from '$app/stores';
	import { goto, invalidate } from '$app/navigation';
	import type { PageData } from './$types';
	import type { Tender, Bid, BidItem, RequisitionItem } from '$lib/types'; // Added Bid, BidItem, RequisitionItem
	import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';
	import { getAccessToken } from '$lib/authService';
	import { user } from '$lib/store';

	export let data: PageData;

	$: tender = data.tender;
	$: error = data.error;
	$: bids = data.bids;
	$: bidsError = data.bidsError;

	let editableTender: Partial<Tender> = {};
	let isLoading: boolean = true;
	let editMode: boolean = false;

	// New Bid Submission State
	let newBid: Partial<Bid> & { items: BidItem[] } = {
		notes: '',
		items: []
		// Add other general bid fields here if needed, e.g., validity_period: null
	};
	// To store File objects separately, as they can't be deeply cloned or easily stored in newBid.items directly with full reactivity for binding
	let bidItemFiles: Array<{ specSheetFile: File | null; itemImageFile: File | null }> = [];

	let bidSubmissionError: string | null = null;
	let bidSubmissionSuccess: string | null = null;
	let isSubmittingBid: boolean = false;

	let successMessage: string | null = null;
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

	// Initialize newBid.items when tender data is available and form should be shown
	$: {
		console.log('[Debug] Tender data before bid form init:', JSON.parse(JSON.stringify(tender)));
		console.log('[Debug] showBidForm:', showBidForm);
		if (tender && tender.requisition) {
			console.log('[Debug] tender.requisition.items:', JSON.parse(JSON.stringify(tender.requisition.items)));
		}
	}
	$: if (showBidForm && tender && tender.requisition && tender.requisition.items && newBid.items.length === 0) {
		newBid.items = tender.requisition.items.map((reqItem: RequisitionItem) => ({
			id: 0, // Placeholder, will be set by backend
			bid_id: 0, // Placeholder
			requisition_item_id: reqItem.id,
			description: reqItem.description,
			quantity: reqItem.quantity,
			unit: reqItem.unit,
			offered_unit_price: 0,
			specification_text: '',
			// File URLs will be set by backend, files handled separately for upload
		}));
		bidItemFiles = tender.requisition.items.map(() => ({ specSheetFile: null, itemImageFile: null }));
		console.log('[Debug] newBid.items initialized:', JSON.parse(JSON.stringify(newBid.items)));
	} else if (!showBidForm && newBid.items.length > 0) {
		// Reset if form is hidden to avoid stale data if user navigates or conditions change
		newBid.items = [];
		bidItemFiles = [];
		newBid.notes = '';
	}
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
			const token = getAccessToken();
			if (!token) {
				errorMessage = 'Authentication token not available. Please log in again.';
				isLoading = false;
				return;
			}

			const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders/${tender.id}`, {
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
		if (!tender || !tender.id) {
			bidSubmissionError = 'Tender information is missing.';
			return;
		}
		if (newBid.items.some(item => typeof item.offered_unit_price !== 'number' || item.offered_unit_price <= 0)) {
			bidSubmissionError = 'Offered unit price must be greater than zero for all items.';
			return;
		}

		isSubmittingBid = true;
		bidSubmissionError = null;
		bidSubmissionSuccess = null;

		try {
			const token = getAccessToken();
			if (!token) {
				bidSubmissionError = 'Authentication token not available. Please log in again.';
				isSubmittingBid = false;
				return;
			}

			const formData = new FormData();

			// Append general bid info
			if (newBid.notes) formData.append('notes', newBid.notes);
			// Calculate total bid amount from items for the 'bid_amount' field if backend still expects it
			let totalBidAmount = 0;
			newBid.items.forEach(item => {
				totalBidAmount += (item.offered_unit_price || 0) * (item.quantity || 0);
			});
			formData.append('bid_amount', totalBidAmount.toFixed(2));

			// Append items JSON (excluding files, as they are handled separately)
			const itemsForJson = newBid.items.map(item => ({
				requisition_item_id: item.requisition_item_id,
				description: item.description,
				quantity: item.quantity,
				unit: item.unit,
				offered_unit_price: item.offered_unit_price,
				specification_text: item.specification_text,
				// URLs will be set by backend
			}));
			formData.append('items_json', JSON.stringify(itemsForJson));

			// Append files
			bidItemFiles.forEach((filePair, index) => {
				if (filePair.specSheetFile) {
					formData.append(`item_spec_sheet_${index}`, filePair.specSheetFile);
				}
				if (filePair.itemImageFile) {
					formData.append(`item_image_${index}`, filePair.itemImageFile);
				}
			});

			const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders/${tender.id}/bids`, {
				method: 'POST',
				headers: {
					// 'Content-Type': 'multipart/form-data' is set automatically by browser when using FormData
					'Authorization': `Bearer ${token}`
				},
				body: formData
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ error: 'Failed to submit bid. Please try again.' }));
				throw new Error(errorData.error || `HTTP error ${response.status}`);
			}

			const createdBid = await response.json();
			bidSubmissionSuccess = `Bid (ID: ${createdBid.id}) submitted successfully with ${createdBid.items?.length || 0} items!`;
			
			// Reset form fields
			newBid.notes = '';
			newBid.items = []; // This will trigger re-initialization by the reactive block if still on page
			bidItemFiles = [];
			
			// Invalidate bids data if displayed on this page (for procurement officers)
			invalidate((url) => url.href.startsWith(`${PUBLIC_VITE_API_BASE_URL}/api/tenders/${tender?.id}/bids`));
      // Also invalidate the tender itself in case its status or bid count changes
      invalidate((url) => url.href === `${PUBLIC_VITE_API_BASE_URL}/api/tenders/${tender?.id}`);
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
			const token = getAccessToken();
			if (!token) {
				errorMessage = 'Authentication token not available. Please log in again.';
				isLoading = false;
				return;
			}

			const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders/${tender.id}`, {
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

	function handleSpecSheetFileChange(event: Event, index: number) {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			bidItemFiles[index].specSheetFile = input.files[0];
		} else {
			bidItemFiles[index].specSheetFile = null;
		}
		bidItemFiles = [...bidItemFiles]; // Trigger reactivity for UI updates
	}

	function handleItemImageFileChange(event: Event, index: number) {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			bidItemFiles[index].itemImageFile = input.files[0];
		} else {
			bidItemFiles[index].itemImageFile = null;
		}
		bidItemFiles = [...bidItemFiles]; // Trigger reactivity for UI updates
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
					<!-- Bid Form with Item Specifications -->
					<form on:submit|preventDefault={handleSubmitBid} class="space-y-8 p-6 bg-white shadow rounded-lg">
						<h2 class="text-2xl font-semibold text-gray-800 mb-6 border-b pb-3">Submit Your Bid for: {tender.title}</h2>

						{#if bidSubmissionError}
							<div class="alert alert-error mb-4"><p>{bidSubmissionError}</p></div>
						{/if}
						{#if bidSubmissionSuccess}
							<div class="alert alert-success mb-4"><p>{bidSubmissionSuccess}</p></div>
						{/if}

						<!-- General Bid Information -->
						<div class="form-control">
							<label for="bid-notes" class="label"><span class="label-text">Additional Notes (Optional)</span></label>
							<textarea id="bid-notes" bind:value={newBid.notes} class="textarea textarea-bordered h-24" placeholder="Any additional comments or information for your bid..."></textarea>
						</div>
						<!-- Add other general bid fields here if needed, e.g., validity period -->

						<h3 class="text-xl font-semibold text-gray-700 mt-6 mb-4">Item Specifications</h3>
						{#if newBid.items && newBid.items.length > 0}
							<div class="space-y-8">
								{#each newBid.items as bidItem, i (bidItem.requisition_item_id)}
									{@const originalReqItem = tender?.requisition?.items?.find(item => item.id === bidItem.requisition_item_id)}
									<div class="card card-bordered bg-base-100 shadow-md p-2 sm:p-4 hover:shadow-lg transition-shadow duration-300 ease-in-out">
										<h4 class="text-lg font-semibold text-primary mb-3 border-b pb-2">Item {i + 1}: {originalReqItem?.description || 'N/A'}</h4>
										<div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
											<!-- Left Column: Original Requisition Item Details -->
											<div class="space-y-3 pr-0 md:pr-4 md:border-r">
												<div>
													<p class="text-sm font-medium text-gray-700">Required Quantity:</p>
													<p class="text-sm text-gray-600">{originalReqItem?.quantity} {originalReqItem?.unit}</p>
												</div>
												{#if originalReqItem?.specification_text || originalReqItem?.specification_sheet_url || originalReqItem?.item_image_url}
												<div class="mt-2 p-3 bg-gray-50 rounded-md border border-gray-200">
													<p class="text-xs font-semibold text-gray-700 mb-1">Original Specifications (from Requisition):</p>
													{#if originalReqItem.specification_text}
														<p class="text-xs text-gray-600 whitespace-pre-wrap"><strong>Details:</strong> {originalReqItem.specification_text}</p>
													{/if}
													{#if originalReqItem.specification_sheet_url}
														<p class="text-xs text-gray-600 mt-1"><strong>Sheet:</strong> <a href={PUBLIC_VITE_API_BASE_URL + originalReqItem.specification_sheet_url} target="_blank" rel="noopener noreferrer" class="link link-hover link-primary">View Original Spec Sheet</a></p>
													{/if}
													{#if originalReqItem.item_image_url}
														<p class="text-xs text-gray-600 mt-1"><strong>Image:</strong> <a href={PUBLIC_VITE_API_BASE_URL + originalReqItem.item_image_url} target="_blank" rel="noopener noreferrer" class="link link-hover link-primary">View Original Image</a></p>
													{/if}
												</div>
												{/if}
											</div>

											<!-- Right Column: Supplier Inputs -->
											<div class="space-y-4">
												<div class="form-control">
													<label for={`item-price-${i}`} class="label pb-1">
														<span class="label-text font-medium">Your Offered Unit Price <span class="text-error">*</span></span>
													</label>
													<input type="number" step="0.01" min="0.01" id={`item-price-${i}`} bind:value={bidItem.offered_unit_price} required class="input input-bordered input-primary w-full" />
												</div>

												<div class="form-control">
													<label for={`item-spec-text-${i}`} class="label pb-1">
														<span class="label-text font-medium">Your Item Specification Text (Optional)</span>
													</label>
													<textarea id={`item-spec-text-${i}`} bind:value={bidItem.specification_text} class="textarea textarea-bordered textarea-primary h-24 w-full" placeholder="Describe your offered item, deviations, or compliance..."></textarea>
												</div>

												<div class="form-control">
													<label for={`item-spec-sheet-${i}`} class="label pb-1">
														<span class="label-text font-medium">Upload Specification Sheet (Optional)</span>
														<span class="label-text-alt">PDF/DOCX, max 10MB</span>
													</label>
													<input type="file" id={`item-spec-sheet-${i}`} on:change={(e) => handleSpecSheetFileChange(e, i)} accept=".pdf,.doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document" class="file-input file-input-bordered file-input-primary w-full" />
													{#if bidItemFiles[i]?.specSheetFile}
														<p class="text-xs text-base-content/70 mt-1">Selected: {bidItemFiles[i].specSheetFile?.name}</p>
													{/if}
												</div>

												<div class="form-control">
													<label for={`item-image-${i}`} class="label pb-1">
														<span class="label-text font-medium">Upload Item Image (Optional)</span>
														<span class="label-text-alt">JPG/PNG, max 10MB</span>
													</label>
													<input type="file" id={`item-image-${i}`} on:change={(e) => handleItemImageFileChange(e, i)} accept="image/jpeg,image/png" class="file-input file-input-bordered file-input-primary w-full" />
													{#if bidItemFiles[i]?.itemImageFile}
														<p class="text-xs text-base-content/70 mt-1">Selected: {bidItemFiles[i].itemImageFile?.name}</p>
													{/if}
												</div>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<div class="alert alert-info shadow-lg">
								<div>
									<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current flex-shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
									<span>No items found for this tender, or the bid form is not ready. Please check back later or contact support if you believe this is an error.</span>
								</div>
							</div>
						{/if}

						<div class="flex justify-end space-x-3 pt-6 border-t mt-8">
							<button type="submit" class="btn btn-primary btn-md" disabled={isSubmittingBid || newBid.items.length === 0}>
								{#if isSubmittingBid}Submitting Bid...{:else}Submit Full Bid{/if}
							</button>
						</div>
					</form>
				{/if}

				<!-- Bid Evaluation Section -->
				{#if canViewBids && !editMode}
					<div class="mt-10 pt-6 border-t">
						<h3 class="text-xl font-semibold text-gray-800 mb-4">Bid Evaluation</h3>
						{#if bidsError}
							<div class="alert alert-error">
								<p>Error loading bid data: {bidsError}</p>
							</div>
						{:else if bids && bids.length > 0}
							<p class="mb-4">{bids.length} bid(s) have been submitted for this tender.</p>
							<a href="/tenders/{tender.id}/evaluation" class="btn btn-primary">
								View Bid Comparison
							</a>
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
