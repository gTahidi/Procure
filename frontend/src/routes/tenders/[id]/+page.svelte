<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import type { Tender } from '$lib/types';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
	import { getAccessTokenSilently } from '$lib/authService';

	export let data: PageData;

	let tender: Tender | null = null;
	let editableTender: Partial<Tender> = {}; // For form binding
	let error: string | null = null;
	let isLoading: boolean = true;
	let editMode: boolean = false;

	let successMessage: string | null = null;
	let errorMessage: string | null = null;

	$: {
		if (data && data.tender) {
			tender = data.tender;
			// Initialize editableTender when tender data is available or changes
			if (tender) {
				editableTender = { ...tender };
				// Ensure dates are in yyyy-MM-dd format for input type="date"
				if (editableTender.published_date) {
					editableTender.published_date = new Date(editableTender.published_date).toISOString().split('T')[0];
				}
				if (editableTender.closing_date) {
					// For datetime-local, it needs to be in YYYY-MM-DDTHH:mm format
					const d = new Date(editableTender.closing_date);
					editableTender.closing_date = new Date(d.getTime() - d.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
				}
			}
			error = data.error;
			isLoading = false;
		} else if (data && data.error) {
			error = data.error;
			tender = null;
			isLoading = false;
		}
	}

	function toggleEditMode() {
		editMode = !editMode;
		successMessage = null;
		errorMessage = null;
		if (editMode && tender) {
			// Re-initialize editableTender from current tender state when entering edit mode
			editableTender = { ...tender };
				if (editableTender.published_date) {
					editableTender.published_date = new Date(editableTender.published_date).toISOString().split('T')[0];
				}
				if (editableTender.closing_date) {
					const d = new Date(editableTender.closing_date);
					editableTender.closing_date = new Date(d.getTime() - d.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
				}
		} 
	}

	async function handleSave() {
		if (!tender || !tender.id) return;
		isLoading = true;
		errorMessage = null;
		successMessage = null;

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
				body: JSON.stringify(editableTender)
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to update tender. The server returned an unexpected response.' }));
				throw new Error(errorData.message || `HTTP error ${response.status}`);
			}

			const updatedTender: Tender = await response.json();
			tender = updatedTender; // Update the main tender object with new data
			// Re-initialize editableTender to reflect saved state and correct date formats for display/next edit
			if (tender) {
				editableTender = { ...tender };
				if (editableTender.published_date) {
					editableTender.published_date = new Date(editableTender.published_date).toISOString().split('T')[0];
				}
				if (editableTender.closing_date) {
					const d = new Date(editableTender.closing_date);
					editableTender.closing_date = new Date(d.getTime() - d.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
				}
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
        // Format to YYYY-MM-DDTHH:mm for datetime-local input
        return new Date(d.getTime() - (d.getTimezoneOffset() * 60000)).toISOString().slice(0, 16);
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
					<div class="flex space-x-2">
						<button on:click={toggleEditMode} class="btn btn-sm btn-outline btn-primary" disabled={isLoading}>Edit</button>
						<button on:click={handleDelete} class="btn btn-sm btn-outline btn-error" disabled={isLoading}>Delete</button>
					</div>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
					<div>
						<h3 class="text-lg font-medium text-gray-700 mb-1">Status</h3>
						<span class:px-2={true} class:py-1={true} class:inline-flex={true} class:text-sm={true} class:leading-5={true} class:font-semibold={true} class:rounded-full={true} 
						  class:bg-blue-100={tender.status === 'published' || tender.status === 'open'}
						  class:text-blue-800={tender.status === 'published' || tender.status === 'open'}
						  class:bg-gray-100={tender.status === 'draft' || tender.status === 'closed'}
						  class:text-gray-800={tender.status === 'draft' || tender.status === 'closed'}
						  class:bg-green-100={tender.status === 'awarded'}
						  class:text-green-800={tender.status === 'awarded'}
						  class:bg-yellow-100={tender.status === 'evaluation'}
						  class:text-yellow-800={tender.status === 'evaluation'}
						  class:bg-red-100={tender.status === 'cancelled'}
						  class:text-red-800={tender.status === 'cancelled'}
						  class:bg-purple-100={!(tender.status === 'published' || tender.status === 'open' || tender.status === 'draft' || tender.status === 'closed' || tender.status === 'awarded' || tender.status === 'evaluation' || tender.status === 'cancelled')}
						  class:text-purple-800={!(tender.status === 'published' || tender.status === 'open' || tender.status === 'draft' || tender.status === 'closed' || tender.status === 'awarded' || tender.status === 'evaluation' || tender.status === 'cancelled')}
						>
						  {tender.status || 'N/A'}
						</span>
					</div>
					<div>
						<h3 class="text-lg font-medium text-gray-700 mb-1">Category</h3>
						<p class="text-gray-600">{tender.category || 'N/A'}</p>
					</div>
					<div>
						<h3 class="text-lg font-medium text-gray-700 mb-1">Budget</h3>
						<p class="text-gray-600">{tender.budget ? tender.budget.toLocaleString(undefined, { style: 'currency', currency: 'USD' }) : 'N/A'}</p> <!-- Adjust currency as needed -->
					</div>
				</div>

				<div class="mb-8">
					<h3 class="text-xl font-semibold text-gray-700 mb-3">Description</h3>
					<p class="text-gray-600 leading-relaxed whitespace-pre-line">{tender.description || 'No description provided.'}</p>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
					<div>
						<h3 class="text-lg font-medium text-gray-700 mb-1">Created By User ID</h3>
						<p class="text-gray-600">{tender.created_by_user_id || 'N/A'}</p>
					</div>
					<div>
						<h3 class="text-lg font-medium text-gray-700 mb-1">Requisition ID</h3>
						<p class="text-gray-600">{tender.requisition_id || 'N/A'}</p>
					</div>
				</div>

				<div class="text-sm text-gray-500 mt-8 pt-4 border-t">
					<p>Last Updated: {formatDate(tender.updated_at)}</p>
					<p>Created At: {formatDate(tender.created_at)}</p>
				</div>
				
				<div class="mt-8 pt-6 border-t border-gray-200 flex justify-start">
					<a href="/tenders" class="btn btn-outline"> Back to Tenders List</a>
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
