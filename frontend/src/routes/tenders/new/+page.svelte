<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Tender } from '$lib/types';

	let formData: Partial<Tender> = {
		title: '',
		description: '',
		category: 'goods',
		budget: undefined,
		published_date: new Date().toISOString().split('T')[0], // Default to today
		closing_date: '',
		status: 'draft'
	};

	let errorMessage: string | null = null;
	let successMessage: string | null = null;
	let isLoading: boolean = false;

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
			// Prepare data for submission, ensuring correct types
			const tenderDataToSubmit: Partial<Tender> = {
				...formData,
				budget: formData.budget ? parseFloat(String(formData.budget)) : undefined,
				published_date: formData.published_date ? new Date(formData.published_date).toISOString() : null,
				closing_date: formData.closing_date ? new Date(formData.closing_date).toISOString() : null
			};

			const response = await fetch('/api/tenders', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
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
			formData = { title: '', description: '', category: 'goods', budget: undefined, published_date: new Date().toISOString().split('T')[0], closing_date: '', status: 'draft' }; 

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
