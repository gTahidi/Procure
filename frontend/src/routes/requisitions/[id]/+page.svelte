<script lang="ts">
	import { page } from '$app/stores';
	import type { PageData } from './$types';

	export let data: PageData;

	// Helper to format date string (e.g., YYYY-MM-DD)
	function formatDate(dateString: string | undefined): string {
		if (!dateString) return 'N/A';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-CA'); // YYYY-MM-DD format
	}

	// The requisition object is now directly from data prop
	// $: requisition = data.requisition; // This is reactive, good if data can change after load
	// For initial load, direct access is fine too, but reactive is safer for future updates.
	// Let's use direct access for now and ensure +page.ts handles errors by throwing SvelteKit errors.

</script>

<svelte:head>
	<title>Requisition {data.requisition?.id || 'Details'} - Procurement System</title>
</svelte:head>

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
								 data.requisition.status === 'Pending Approval' || data.requisition.status === 'Submitted for Approval' ? 'bg-yellow-100 text-yellow-800' : 
								 data.requisition.status === 'Rejected' ? 'bg-red-100 text-red-800' : 
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
				<p class="text-gray-600">{data.requisition.type.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase())}</p>
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

		<div class="mt-8 pt-6 border-t border-gray-200 flex justify-end space-x-3">
			<a href="/requisitions" class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 no-underline">
				Back to List
			</a>

			{#if data.requisition.status === 'Approved'}
				<a href={`/tenders/new?requisitionId=${data.requisition.id}`}
				   class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50 no-underline">
					Create Tender
				</a>
			{/if}

			{#if data.requisition.status === 'Draft'}
			<button class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-opacity-50">
				Edit Requisition
			</button>
			{:else if data.requisition.status !== 'Approved'} <!-- Only show disabled edit if not Draft and not Approved (where Create Tender shows) -->
			<button class="px-4 py-2 bg-gray-400 text-white rounded-md cursor-not-allowed" disabled>
				Edit Requisition
			</button>
			{/if}
		</div>
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
