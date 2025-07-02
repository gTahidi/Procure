<script lang="ts">
	import { page } from '$app/stores';
	import type { PageData } from './$types';
	import { user } from '$lib/store';
	import { getAccessTokenSilently } from '$lib/authService';
	import { PUBLIC_VITE_API_BASE_URL } from '$env/static/public';
	import { invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';

	export let data: PageData;

	onMount(() => {
		console.log('MOUNTED: Current user object:', $user);
		console.log('MOUNTED: User role (direct):', $user?.role);
		console.log('MOUNTED: Is admin (direct check)?:', $user?.role === 'admin');
		console.log('MOUNTED: User ID (for 2nd approval):', $user?.id);
		console.log('MOUNTED: Requisition data on mount:', data.requisition);
		console.log('MOUNTED: Requisition status on mount:', data.requisition?.status);
		console.log('MOUNTED: Requisition approver_one_id:', data.requisition?.approver_one_id);
	});

	let rejectionReason = '';
	let showRejectionModal = false;
	let isProcessing = false;
	let apiError = '';

	// Helper to format date string (e.g., YYYY-MM-DD)
	function formatDate(dateString: string | undefined): string {
		if (!dateString) return 'N/A';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-CA'); // YYYY-MM-DD format
	}

	async function handleRequisitionAction(action: 'approve' | 'reject', reason?: string) {
		if (!data.requisition?.id) return;
		isProcessing = true;
		apiError = '';
		const token = await getAccessTokenSilently();
		if (!token) {
			apiError = 'Authentication error. Please log in again.';
			isProcessing = false;
			return;
		}

		try {
			const response = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/requisitions/${data.requisition.id}/action`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({ action, reason })
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || `Failed to ${action} requisition: ${response.statusText}`);
			}

			// alert(`Requisition successfully ${action}d.`);
			showRejectionModal = false; // Close modal if open
			rejectionReason = ''; // Clear reason
			await invalidateAll(); // Refetch all data for the current page
		} catch (err: any) {
			console.error(`Error ${action}ing requisition:`, err);
			apiError = err.message || `An unexpected error occurred while ${action}ing.`;
			// alert(`Error: ${apiError}`);
		} finally {
			isProcessing = false;
		}
	}

	function promptReject() {
		apiError = ''; // Clear previous errors
		rejectionReason = ''; // Clear previous reason
		showRejectionModal = true;
	}

	function submitRejection() {
		if (!rejectionReason.trim()) {
			apiError = 'Rejection reason cannot be empty.';
			return;
		}
		handleRequisitionAction('reject', rejectionReason.trim());
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape' && showRejectionModal) {
			showRejectionModal = false;
		}
	}

	// Reactive statements to determine button visibility
	$: isAdmin = $user?.role === 'admin';
	$: canPerformFirstApproval = isAdmin && (data.requisition?.status === 'pending_approval_1' || data.requisition?.status === 'submitted_for_approval');
	$: canPerformSecondApproval = isAdmin && 
							 data.requisition?.status === 'pending_approval_2' && 
							 $user?.id !== data.requisition?.approver_one_id;
	$: canRejectFirstStep = canPerformFirstApproval;
	$: canRejectSecondStep = canPerformSecondApproval;

	// Display rejection reason if present
	$: displayRejectionReason = data.requisition?.status === 'rejected' && data.requisition?.rejection_reason;
</script>

<svelte:head>
	<title>Requisition {data.requisition?.id || 'Details'} - Procurement System</title>
</svelte:head>

<svelte:window on:keydown={handleKeyDown}/>

{#if data.requisition}
<div class="container mx-auto py-8 px-4">
	<div class="bg-white shadow-lg rounded-lg p-6 md:p-8">
		<div class="mb-6 pb-4 border-b border-gray-200">
			<h1 class="text-3xl font-semibold text-gray-800">Requisition ID: {data.requisition.id}</h1>
			<p class="text-sm text-gray-500">Created on: {formatDate(data.requisition.created_at)} by User ID: {data.requisition.user_id}</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Requisition ID</h3>
				<p class="text-gray-600">{data.requisition.id}</p>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Status</h3>
				<span class={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full 
								${data.requisition.status === 'Approved' ? 'bg-green-100 text-green-800' : 
								 data.requisition.status === 'pending_approval_1' || data.requisition.status === 'pending_approval_2' || data.requisition.status === 'Pending Approval' || data.requisition.status === 'Submitted for Approval' ? 'bg-yellow-100 text-yellow-800' : 
					 data.requisition.status === 'rejected' ? 'bg-red-100 text-red-800' : 
					 data.requisition.status === 'Draft' ? 'bg-blue-100 text-blue-800' :
								'bg-gray-100 text-gray-800'}`}>
					{data.requisition.status.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase())}
				</span>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Requester (User ID)</h3>
				<p class="text-gray-600">{data.requisition.user_id}</p>
			</div>
			{#if data.requisition.type}
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Type</h3>
				<p class="text-gray-600">{data.requisition.type?.replace(/_/g, ' ')?.replace(/\b\w/g, (l: string) => l.toUpperCase()) || 'N/A'}</p>
			</div>
			{/if}
			{#if data.requisition.aac}
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">AAC (Annual Allocation Code)</h3>
				<p class="text-gray-600">{data.requisition.aac}</p>
			</div>
			{/if}
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Description</h3>
			<p class="text-gray-600 leading-relaxed whitespace-pre-line">{data.requisition.description || 'No description provided.'}</p>
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Items</h3>
			{#if data.requisition.items && data.requisition.items.length > 0}
				<ul class="divide-y divide-gray-200">
					{#each data.requisition.items as item (item.id) }
						<li class="py-3 flex justify-between items-center">
							<span class="text-gray-700">{item.description} ({item.unit})</span>
							<span class="text-gray-500">Quantity: {item.quantity} @ ${item.estimated_unit_price?.toFixed(2) || 'N/A'}/unit</span>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-gray-500 italic">No items listed for this requisition.</p>
			{/if}
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Attachments</h3>
			<!-- Currently, backend does not support attachments for requisitions -->
			<p class="text-gray-500 italic">Attachments feature not yet implemented for requisitions.</p>
			<!-- {#if data.requisition.attachments && data.requisition.attachments.length > 0}
				<ul class="list-disc list-inside space-y-1">
					{#each data.requisition.attachments as attachment (attachment.name)}
						<li>
							<a href={attachment.url} class="text-indigo-600 hover:text-indigo-800 hover:underline" target="_blank" rel="noopener noreferrer">
								{attachment.name}
							</a>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-gray-500 italic">No attachments for this requisition.</p>
			{/if} -->
		</div>

		<!-- Display Rejection Reason if applicable -->
		{#if displayRejectionReason}
		<div class="mt-6 p-4 border-l-4 border-red-500 bg-red-50 rounded-md">
			<h4 class="text-md font-semibold text-red-700">Rejection Reason:</h4>
			<p class="text-red-600 whitespace-pre-line">{data.requisition.rejection_reason}</p>
		</div>
		{/if}

		<div class="mt-8 pt-6 border-t border-gray-200 flex flex-wrap justify-end gap-3">
			<!-- Admin Action Buttons -->
			{#if canPerformFirstApproval}
				<button on:click={() => handleRequisitionAction('approve')} 
						class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
						disabled={isProcessing}>Approve (1st Step)
				</button>
			{/if}
			{#if canPerformSecondApproval}
				<button on:click={() => handleRequisitionAction('approve')} 
						class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
						disabled={isProcessing}>Approve (2nd Step)
				</button>
			{/if}
			{#if canRejectFirstStep || canRejectSecondStep}
				<button on:click={promptReject} 
					class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50"
					disabled={isProcessing}>Reject
				</button>
			{/if}

			<a href="/requisitions" class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 no-underline">
				Back to List
			</a>

			{#if data.requisition.status === 'Approved'}
				<a href={`/tenders/new?requisitionId=${data.requisition.id}`}
				   class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50 no-underline">
					Create Tender
				</a>
			{/if}

			{#if data.requisition.status === 'Draft' && !isAdmin}
			<button class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-opacity-50">
				Edit Requisition
			</button>
			{:else if data.requisition.status !== 'Approved' && !isAdmin} <!-- Only show disabled edit if not Draft and not Approved (where Create Tender shows) -->
			<button class="px-4 py-2 bg-gray-400 text-white rounded-md cursor-not-allowed" disabled>
				Edit Requisition
			</button>
			{/if}
		</div>

		<!-- Rejection Reason Modal -->
		{#if showRejectionModal}
		<div class="fixed inset-0 bg-gray-600 bg-opacity-75 overflow-y-auto h-full w-full z-50 flex justify-center items-center p-4" role="dialog" aria-modal="true">
		  <div class="relative p-6 border w-full max-w-lg shadow-xl rounded-lg bg-white" role="document">
			<h3 class="text-xl font-semibold leading-6 text-gray-900 mb-4">Provide Rejection Reason</h3>
			<div class="mt-2">
			  <textarea bind:value={rejectionReason}
						class="w-full px-3 py-2 text-gray-700 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
						rows="4"
						placeholder="Enter reason for rejection (required)..."></textarea>
			</div>
			{#if apiError && !rejectionReason.trim() && showRejectionModal}
			  <p class="text-sm text-red-600 mt-1">{apiError}</p>
			{/if}
			<div class="mt-6 flex justify-end space-x-3">
			  <button type="button" 
					  on:click={() => { showRejectionModal = false; apiError = ''; }}
					  class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400 disabled:opacity-50"
					  disabled={isProcessing}>Cancel</button>
			  <button type="button" 
					  on:click={submitRejection} 
					  class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50"
					  disabled={isProcessing || !rejectionReason.trim()}>{isProcessing ? 'Submitting...' : 'Submit Rejection'}</button>
			</div>
		  </div>
		</div>
		{/if}
		<!-- End Rejection Reason Modal -->

	</div>
</div>


{:else if $page.error}
	<div class="container mx-auto py-8 px-4 text-center">
		<h1 class="text-2xl font-semibold text-red-600">Error Loading Requisition</h1>
		<p class="text-gray-600">{$page.error.message || 'An unknown error occurred.'}</p>
		<a href="/requisitions" class="mt-4 inline-block px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 no-underline">
			Back to Requisitions List
		</a>
	</div>
{:else}
	<div class="container mx-auto py-8 px-4 text-center">
		<p class="text-gray-600">Loading requisition details...</p>
	</div>
{/if}
