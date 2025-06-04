<script lang="ts">
	import DocumentUpload from '$lib/components/DocumentUpload.svelte';
	import { user } from '$lib/store';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getAccessTokenSilently } from '$lib/authService'; // Import getAccessTokenSilently
	import type { RequisitionItem as BaseRequisitionItem } from '../../../lib/types';

	// Local type for form items, extending base for UI specific fields
	interface FormRequisitionItem extends BaseRequisitionItem {
		total: number | null; // UI calculated field
		// specificationSheetFile and itemImageFile are already in BaseRequisitionItem as per last change to types.ts
	}

	let requisition = {
		requesterId: '', // Should be pre-filled or selected (e.g. logged in user)
		department: '', // Could be derived from user or selected
		requisitionType: 'goods', // 'goods', 'services', 'fixed_asset'
		aac: '', // Authority to Incur Commitment
		estimatedCost: null as number | null,
		description: '',
		urgencyLevel: 'medium', // 'low', 'medium', 'high'
		requiredByDate: '',
		deliveryLocation: '',
		items: [
			{
				description: '',
				quantity: 1,
				unit: 'unit',
				estimated_unit_price: null as number | null,
				total: 0 as number | null,
				specification_text: '',
				specificationSheetFile: null as File | null,
				itemImageFile: null as File | null
			}
		] as FormRequisitionItem[],
		notes: '',
		status: 'draft' // Initial status
	};

	let attachedFiles: any[] = [];
	let loading = false; // For API call feedback
	let submissionMessage = ''; // To display success/error messages

	onMount(() => {
		if ($user && $user.id) { 
			requisition.requesterId = String($user.id); 
		} else {
			console.error('User not logged in or user ID not available for PR creation.');
			submissionMessage = 'Error: You must be logged in to create a requisition.';
		}
	});

	function handleFilesAttached(event: CustomEvent) {
		attachedFiles = [...attachedFiles, ...event.detail];
		console.log('Files for requisition:', attachedFiles);
		// In a real app, you'd likely associate these file details with the requisition data
	}

	function addItem() {
		requisition.items = [
			...requisition.items,
			{ 
				description: '', 
				quantity: 1, 
				unit: 'unit', 
				estimated_unit_price: null as number | null, 
				total: 0 as number | null,
				specification_text: '',
				specificationSheetFile: null as File | null,
				itemImageFile: null as File | null
			}
		];
	}

	function removeItem(index: number) {
		requisition.items = requisition.items.filter((_, i) => i !== index);
		calculateTotals();
	}

	function calculateItemTotal(item: FormRequisitionItem) {
		if (item.quantity && typeof item.estimated_unit_price === 'number') {
			item.total = item.quantity * item.estimated_unit_price;
		} else {
			item.total = 0;
		}
		calculateTotals();
	}

	function calculateTotals() {
		requisition.estimatedCost = requisition.items.reduce((sum, item) => sum + (item.total || 0), 0);
	}

	async function handleSubmit(isDraft: boolean) {
		loading = true;
		submissionMessage = '';
		requisition.status = isDraft ? 'draft' : 'submitted_for_approval';

		if (!$user || !$user.id) {
			console.error('Attempted to submit PR without a logged-in user or user ID.');
			submissionMessage = 'Error: You must be logged in with a valid user ID to create a requisition.';
			loading = false;
			return; 
		}

		// Fetch the access token
		let token;
		try {
			token = await getAccessTokenSilently();
			if (!token) {
				submissionMessage = 'Error: Could not retrieve authentication token. Please log in again.';
				loading = false;
				return;
			}
		} catch (e) {
			console.error('Error fetching token:', e);
			submissionMessage = 'Error: Failed to obtain authentication token.';
			loading = false;
			return;
		}

		const payload = {
			user_id: $user.id, 
			type: requisition.requisitionType,
			aac: requisition.aac || null,
			status: requisition.status,
			items: requisition.items.map(item => ({
				description: item.description,
				quantity: Number(item.quantity),
				unit: item.unit,
				estimated_unit_price: typeof item.estimated_unit_price === 'number' ? Number(item.estimated_unit_price) : 0,
				freight_cost: 0,
				insurance_cost: 0,
				installation_cost: 0,
				value: typeof item.estimated_unit_price === 'number' ? Number(item.estimated_unit_price) * item.quantity : 0
			}))
		};

		try {
			const response = await fetch('http://localhost:8080/api/requisitions', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Accept': 'application/json',
					'Authorization': `Bearer ${token}` // Add Authorization header
				},
				body: JSON.stringify(payload)
			});

			const responseData = await response.json();

			if (!response.ok) {
				// Use the error message from backend if available, otherwise default
				const errorMsg = responseData.error || responseData.message || `HTTP error! status: ${response.status}`;
				throw new Error(errorMsg);
			}

			submissionMessage = `Requisition ${isDraft ? 'saved as draft' : 'submitted for approval'} successfully! ID: ${responseData.id}`;
			console.log('Success:', responseData);

			// If successfully submitted (not just saved as draft), redirect to the requisitions list page
			if (!isDraft) {
				setTimeout(() => { // Short delay to allow user to see success message
					goto('/requisitions');
				}, 1500); // 1.5 seconds delay
			} else {
				// Reset form on successful draft save
				requisition = {
					requesterId: $user?.id ? String($user.id) : '', 
					department: '',
					requisitionType: 'goods',
					aac: '',
					estimatedCost: null,
					description: '',
					urgencyLevel: 'medium',
					requiredByDate: '',
					deliveryLocation: '',
					items: [{
						description: '',
						quantity: 1,
						unit: 'unit',
						estimated_unit_price: null as number | null,
						total: 0 as number | null,
						specification_text: '',
						specificationSheetFile: null as File | null,
						itemImageFile: null as File | null
					}] as FormRequisitionItem[],
					notes: '',
					status: 'draft'
				};
				attachedFiles = [];
			}

		} catch (error: any) {
			console.error('Submission error:', error);
			submissionMessage = `Error: ${error.message || 'Failed to submit requisition'}`;
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>New Purchase Requisition - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4 max-w-4xl">
	<h1 class="text-3xl font-semibold mb-8 text-gray-800">Create New Purchase Requisition</h1>

	<form on:submit|preventDefault class="space-y-8 bg-white p-8 rounded-lg shadow-lg">
		
		<!-- Section 1: Basic Information -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Basic Information</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div>
					<label for="requesterId" class="block text-sm font-medium text-gray-700">Requester ID (Auto-filled)</label>
					<input type="text" id="requesterId" bind:value={requisition.requesterId} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm bg-gray-100" readonly placeholder="Logged-in User ID">
				</div>
				<div>
					<label for="department" class="block text-sm font-medium text-gray-700">Department</label>
					<input type="text" id="department" bind:value={requisition.department} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., IT Department">
				</div>
				<div>
					<label for="requisitionType" class="block text-sm font-medium text-gray-700">Requisition Type</label>
					<select id="requisitionType" bind:value={requisition.requisitionType} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						<option value="goods">Goods</option>
						<option value="services">Services</option>
						<option value="fixed_asset">Fixed Asset</option>
					</select>
				</div>
				<div>
					<label for="aac" class="block text-sm font-medium text-gray-700">Authority to Incur Commitment (AAC)</label>
					<input type="text" id="aac" bind:value={requisition.aac} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., AAC-12345">
				</div>
				<div>
					<label for="urgencyLevel" class="block text-sm font-medium text-gray-700">Urgency Level</label>
					<select id="urgencyLevel" bind:value={requisition.urgencyLevel} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						<option value="low">Low</option>
						<option value="medium">Medium</option>
						<option value="high">High</option>
					</select>
				</div>
				<div>
					<label for="requiredByDate" class="block text-sm font-medium text-gray-700">Required By Date</label>
					<input 
						type="date" 
						id="requiredByDate" 
						bind:value={requisition.requiredByDate}
						min={new Date().toISOString().split('T')[0]}
						class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
					>
				</div>
			</div>
			<div class="mt-6">
				<label for="description" class="block text-sm font-medium text-gray-700">Overall Description / Purpose</label>
				<textarea id="description" rows="3" bind:value={requisition.description} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Briefly describe the purpose of this requisition..."></textarea>
			</div>
			<div class="mt-6">
				<label for="deliveryLocation" class="block text-sm font-medium text-gray-700">Delivery Location</label>
				<input type="text" id="deliveryLocation" bind:value={requisition.deliveryLocation} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., Main Office Warehouse">
			</div>
		</section>

		<!-- Section 2: Requisition Items -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Requisition Items</h2>
			{#each requisition.items as item, i}
				<div class="p-4 border border-gray-200 rounded-md mb-4 space-y-4 relative bg-gray-50">
					<p class="font-medium text-gray-600">Item #{i + 1}</p>
					<div>
						<label for={`item_desc_${i}`} class="block text-sm font-medium text-gray-700">Item Description</label>
						<input type="text" id={`item_desc_${i}`} bind:value={item.description} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., Laptop, Model XYZ">
					</div>
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
						<div>
							<label for={`item_qty_${i}`} class="block text-sm font-medium text-gray-700">Quantity</label>
							<input type="number" min="1" id={`item_qty_${i}`} bind:value={item.quantity} on:input={() => calculateItemTotal(item)} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						</div>
						<div>
							<label for={`item_unit_${i}`} class="block text-sm font-medium text-gray-700">Unit</label>
							<input type="text" id={`item_unit_${i}`} bind:value={item.unit} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., pcs, kg, hour">
						</div>
						<div>
							<label for={`item_unit_price_${i}`} class="block text-sm font-medium text-gray-700">Estimated Unit Price</label>
							<input type="number" step="0.01" id={`item_estimated_unit_price_${i}`} bind:value={item.estimated_unit_price} on:input={() => calculateItemTotal(item)} placeholder="e.g., 1500.00" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
					</div>

					<!-- Item Specification Text (Full Width for this item's grid) -->
					<div class="md:col-span-2">
						<label for={`item_spec_text_${i}`} class="block text-sm font-medium text-gray-700">Item Specification</label>
						<textarea id={`item_spec_text_${i}`} bind:value={item.specification_text} rows="3" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Detailed specification of the item..."></textarea>
					</div>

					<!-- Specification Sheet Upload -->
					<div>
						<label for={`item_spec_sheet_${i}`} class="block text-sm font-medium text-gray-700">Spec Sheet (Optional)</label>
						<input type="file" id={`item_spec_sheet_${i}`} on:change={(e) => item.specificationSheetFile = (e.target as HTMLInputElement).files?.[0] ?? null} accept=".pdf,.doc,.docx,.txt,.xls,.xlsx" class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100">
						{#if item.specificationSheetFile}
							<p class="mt-1 text-xs text-gray-500">Selected: {item.specificationSheetFile.name}</p>
						{/if}
					</div>

					<!-- Item Image Upload -->
					<div>
						<label for={`item_image_${i}`} class="block text-sm font-medium text-gray-700">Item Image (Optional)</label>
						<input type="file" id={`item_image_${i}`} on:change={(e) => item.itemImageFile = (e.target as HTMLInputElement).files?.[0] ?? null} accept="image/*" class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100">
						{#if item.itemImageFile}
							<p class="mt-1 text-xs text-gray-500">Selected: {item.itemImageFile.name}</p>
						{/if}
					</div>
					</div>
					<div>
						<label for={`item_total_${i}`} class="block text-sm font-medium text-gray-700">Item Total</label>
						<input type="text" id={`item_total_${i}`} value={item.total != null ? item.total.toFixed(2) : '0.00'} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm bg-gray-100" readonly>
					</div>
					{#if requisition.items.length > 1}
						<button type="button" on:click={() => removeItem(i)} aria-label="Remove item" class="absolute top-2 right-2 text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-100">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>
						</button>
					{/if}
				</div>
			{/each}
			<button type="button" on:click={addItem} class="mt-2 text-sm text-indigo-600 hover:text-indigo-800 font-medium flex items-center">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" /></svg>
				Add Another Item
			</button>
			<div class="mt-6 text-right">
				<p class="block text-lg font-semibold text-gray-700">Total Estimated Cost: KES {requisition.estimatedCost != null ? requisition.estimatedCost.toFixed(2) : '0.00'}</p>
			</div>
		</section>

		<!-- Section 3: Supporting Documents -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Supporting Documents</h2>
			<DocumentUpload on:filesAttached={handleFilesAttached} />
			{#if attachedFiles.length > 0}
				<div class="mt-4">
					<h3 class="text-md font-medium text-gray-700 mb-2">Attached files:</h3>
					<ul class="list-disc list-inside pl-2 text-sm text-gray-600 border border-gray-200 p-3 rounded-md bg-gray-50">
						{#each attachedFiles as fileDetail}
							<li>{fileDetail.name} ({fileDetail.documentType})</li>
						{/each}
					</ul>
				</div>
			{/if}
		</section>

		<!-- Section 4: Notes -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Additional Notes</h2>
			<textarea id="notes" rows="4" bind:value={requisition.notes} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Any additional notes or justifications for this requisition..."></textarea>
		</section>

		<!-- Actions -->
		<div class="pt-6 border-t mt-4 flex justify-end space-x-3">
			<button type="button" on:click={() => handleSubmit(true)} class="bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium py-2 px-4 rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500">
				Save as Draft
			</button>
			<button type="button" on:click={() => handleSubmit(false)} class="bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
				Submit for Approval
			</button>
		</div>
	</form>

	{#if submissionMessage}
		<div 
			class="mt-4 p-4 rounded-md text-sm"
			class:bg-green-100={submissionMessage.startsWith('Requisition')}
			class:text-green-700={submissionMessage.startsWith('Requisition')}
			class:bg-red-100={submissionMessage.startsWith('Error') || submissionMessage.startsWith('Network')}
			class:text-red-700={submissionMessage.startsWith('Error') || submissionMessage.startsWith('Network')}
			role="alert"
			aria-live="polite"
		>
			{submissionMessage}
		</div>
	{/if}

</div>
