<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Tender } from '$lib/types';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getAccessTokenSilently } from '$lib/authService';
	import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';

	// Simplified interfaces for what we expect from the requisition endpoint
	interface RequisitionItem {
		id: number;
		description: string;
		quantity: number;
		estimated_unit_price?: number | null;
	}

	interface BackendRequisition {
		id: number;
		type: string; // 'goods', 'services', 'fixed_asset'
		items?: RequisitionItem[];
		// Add other fields if needed for pre-filling, e.g., a top-level description or title
	}

	let formData: Partial<Tender> & { requisition_id?: number | null } = {
		title: '',
		description: '',
		category: 'goods', // Default category
		budget: undefined,
		published_date: new Date().toISOString().split('T')[0], // Default to today
		closing_date: '',
		status: 'draft',
		requisition_id: null
	};

	let linkedRequisitionId: number | null = null; // For display
	let errorMessage: string | null = null;
	let successMessage: string | null = null;
	let isLoading: boolean = false;

	function mapRequisitionTypeToTenderCategory(reqType: string): string {
		switch (reqType) {
			case 'goods':
			case 'fixed_asset': // Assuming fixed assets are procured as goods or works
				return 'goods'; // Or 'works' if more appropriate for fixed_asset
			case 'services':
				return 'services';
			default:
				return 'goods'; // Default category
		}
	}

	onMount(async () => {
		const unsubscribePage = page.subscribe(async ($page) => {
			const requisitionIdParam = $page.url.searchParams.get('requisitionId');
			if (requisitionIdParam) {
				isLoading = true;
				errorMessage = null;
				try {
					const token = await getAccessTokenSilently();
					if (!token) {
						throw new Error('Authentication token not available. Please log in.');
					}
					const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/requisitions/${requisitionIdParam}`, {
						headers: {
							'Authorization': `Bearer ${token}`,
							'Accept': 'application/json'
						}
					});
					if (!response.ok) {
						const errorData = await response.json().catch(() => ({ message: `Failed to fetch requisition ${requisitionIdParam}`}));
						throw new Error(errorData.message || `HTTP error ${response.status}`);
					}
					const requisitionData: BackendRequisition = await response.json();

					formData.requisition_id = requisitionData.id;
					linkedRequisitionId = requisitionData.id; // For display
					formData.title = `Tender based on Requisition REQ-${String(requisitionData.id).padStart(3, '0')}`;
					
					let desc = 'Derived from Requisition. Items include:';
					if (requisitionData.items && requisitionData.items.length > 0) {
						desc += requisitionData.items.map(item => `\n- ${item.description} (Qty: ${item.quantity})`).join('');
					} else {
						desc = 'Refer to linked Requisition for detailed item requirements.';
					}
					formData.description = desc;

					formData.category = mapRequisitionTypeToTenderCategory(requisitionData.type);

					const calculatedBudget = requisitionData.items?.reduce((sum, item) => {
						return sum + (item.quantity * (item.estimated_unit_price || 0));
					}, 0);
					formData.budget = calculatedBudget && calculatedBudget > 0 ? parseFloat(calculatedBudget.toFixed(2)) : undefined;

				} catch (err: any) {
					errorMessage = err.message || 'Failed to load requisition data for pre-filling.';
					console.error('Error pre-filling tender from requisition:', err);
				} finally {
					isLoading = false;
				}
			}
		});
	});

	async function handleSubmit() {
		isLoading = true;
		errorMessage = null;
		successMessage = null;

		// Basic validation (can be expanded)
		if (!formData.title || !formData.closing_date) {
			errorMessage = 'Title and Closing Date are required.';
			isLoading = false;
			return;
		}

		try {
			const token = await getAccessTokenSilently();
			if (!token) {
				errorMessage = 'Authentication token not available. Please log in.';
				isLoading = false;
				return;
			}

			// Prepare data for submission, ensuring correct types
			const tenderDataToSubmit: Partial<Tender> = {
				...formData,
				budget: formData.budget ? parseFloat(String(formData.budget)) : undefined,
				published_date: formData.published_date ? new Date(formData.published_date).toISOString() : null,
				closing_date: formData.closing_date ? new Date(formData.closing_date).toISOString() : null
			};

			const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/tenders`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${token}`
				},
				body: JSON.stringify(tenderDataToSubmit)
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to submit tender. Please try again.' }));
				throw new Error(errorData.message || `HTTP error ${response.status}`);
			}

			const newTender: Tender = await response.json();
			successMessage = 'Tender created successfully!';
			
			// Clear form or navigate
			formData = { title: '', description: '', category: 'goods', budget: undefined, published_date: new Date().toISOString().split('T')[0], closing_date: '', status: 'draft', requisition_id: null }; 

			// Optional: Navigate to the new tender's page or the list page after a short delay
			setTimeout(() => {
				goto(`/tenders/${newTender.id}`); // Or goto('/tenders');
			}, 1500);

		} catch (err: any) {
			console.error('Submission error:', err);
			errorMessage = err.message || 'An unexpected error occurred.';
		} finally {
			isLoading = false;
		}
	}

</script>

<svelte:head>
	<title>New Tender - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4 max-w-2xl">
	<h1 class="text-3xl font-semibold mb-8 text-gray-800">Create New Tender</h1>

	{#if linkedRequisitionId}
		<div class="mb-4 p-3 bg-indigo-50 border border-indigo-200 rounded-md text-sm text-indigo-700">
			Pre-filling form based on Requisition ID: <span class="font-semibold">REQ-{String(linkedRequisitionId).padStart(3, '0')}</span>.
			Please review and complete all required fields.
		</div>
	{/if}

	<form on:submit|preventDefault={handleSubmit} class="space-y-6 bg-white p-8 rounded-lg shadow-lg">

		{#if successMessage}
			<div class="alert alert-success">
				<p>{successMessage}</p>
			</div>
		{/if}
		{#if errorMessage}
			<div class="alert alert-error">
				<p>{errorMessage}</p>
			</div>
		{/if}

		<div>
			<label for="title" class="block text-sm font-medium text-gray-700">Tender Title <span class="text-red-500">*</span></label>
			<input type="text" id="title" bind:value={formData.title} required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., Supply of Office Stationery">
		</div>

		<div>
			<label for="description" class="block text-sm font-medium text-gray-700">Detailed Description</label>
			<textarea id="description" rows="4" bind:value={formData.description} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Provide a detailed description of the tender requirements..."></textarea>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div>
				<label for="category" class="block text-sm font-medium text-gray-700">Category</label>
				<select id="category" bind:value={formData.category} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
					<option value="goods">Goods</option>
					<option value="services">Services</option>
					<option value="works">Works</option>
					<option value="consultancy">Consultancy</option>
				</select>
			</div>
			<div>
				<label for="budget" class="block text-sm font-medium text-gray-700">Budget (Optional)</label>
				<input type="number" step="0.01" id="budget" bind:value={formData.budget} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., 50000.00">
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div>
				<label for="published_date" class="block text-sm font-medium text-gray-700">Published Date</label>
				<input
					type="date"
					id="published_date"
					bind:value={formData.published_date}
					class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
				/>
			</div>
			<div>
				<label for="closingDate" class="block text-sm font-medium text-gray-700">Closing Date <span class="text-red-500">*</span></label>
				<input type="datetime-local" id="closingDate" bind:value={formData.closing_date} required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
			</div>
		</div>
		
		<div>
			<label for="status" class="block text-sm font-medium text-gray-700">Initial Status</label>
			<select id="status" bind:value={formData.status} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
				<option value="draft">Draft</option>
				<option value="published">Published</option>
			</select>
		</div>

		<div class="flex justify-end space-x-3 pt-4">
			<a href="/tenders" class="btn btn-ghost">Cancel</a>
			<button type="submit" class="btn btn-primary" disabled={isLoading}>
				{#if isLoading}Saving...{:else}Save Tender{/if}
			</button>
		</div>
	</form>
</div>
