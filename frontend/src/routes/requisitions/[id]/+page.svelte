<script lang="ts">
	import { page } from '$app/stores';

	// Get the requisition ID from the route parameters
	const requisitionId = $page.params.id;

	// Mock data - in a real app, you'd fetch this based on requisitionId
	const mockRequisitionDetails = {
		id: requisitionId,
		title: `Details for Requisition ${requisitionId}`,
		description: 'This is a detailed description of the requisition. It would include items, quantities, justifications, budget codes, and approval history.',
		requester: 'Mock User',
		status: 'Pending Approval',
		creationDate: '2024-05-13',
		items: [
			{ name: 'Item A', quantity: 10, unitPrice: 50 },
			{ name: 'Item B', quantity: 5, unitPrice: 120 },
		],
		attachments: [
			{ name: 'Quote_SupplierX.pdf', url: '#' },
			{ name: 'SpecSheet_ItemA.docx', url: '#' },
		]
	};

	// If you want to simulate finding a specific requisition from a list (like the one on the list page)
	// You could import that list or have a shared store, then find by ID.
	// For this placeholder, we'll just use the generated mock details directly.
	const requisition = mockRequisitionDetails;

</script>

<svelte:head>
	<title>Requisition {requisition.id} - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	<div class="bg-white shadow-lg rounded-lg p-6 md:p-8">
		<div class="mb-6 pb-4 border-b border-gray-200">
			<h1 class="text-3xl font-semibold text-gray-800">{requisition.title}</h1>
			<p class="text-sm text-gray-500">Created on: {requisition.creationDate} by {requisition.requester}</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Requisition ID</h3>
				<p class="text-gray-600">{requisition.id}</p>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Status</h3>
				<span class={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full 
								${requisition.status === 'Approved' ? 'bg-green-100 text-green-800' : 
								 requisition.status === 'Pending Approval' ? 'bg-yellow-100 text-yellow-800' : 
								 requisition.status === 'Rejected' ? 'bg-red-100 text-red-800' : 
								'bg-gray-100 text-gray-800'}`}>
					{requisition.status}
				</span>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Requester</h3>
				<p class="text-gray-600">{requisition.requester}</p>
			</div>
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Description</h3>
			<p class="text-gray-600 leading-relaxed whitespace-pre-line">{requisition.description}</p>
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Items</h3>
			{#if requisition.items && requisition.items.length > 0}
				<ul class="divide-y divide-gray-200">
					{#each requisition.items as item (item.name)}
						<li class="py-3 flex justify-between items-center">
							<span class="text-gray-700">{item.name}</span>
							<span class="text-gray-500">Quantity: {item.quantity} @ ${item.unitPrice}/unit</span>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-gray-500 italic">No items listed for this requisition.</p>
			{/if}
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Attachments</h3>
			{#if requisition.attachments && requisition.attachments.length > 0}
				<ul class="list-disc list-inside space-y-1">
					{#each requisition.attachments as attachment (attachment.name)}
						<li>
							<a href={attachment.url} class="text-indigo-600 hover:text-indigo-800 hover:underline" target="_blank" rel="noopener noreferrer">
								{attachment.name}
							</a>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-gray-500 italic">No attachments for this requisition.</p>
			{/if}
		</div>

		<div class="mt-8 pt-6 border-t border-gray-200 flex justify-end space-x-3">
			<!-- Placeholder for action buttons like Approve, Reject, Edit, etc. -->
			<button class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50">
				Back to List
			</button>
			<button class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-opacity-50">
				Edit Requisition (Placeholder)
			</button>
		</div>
	</div>
</div>
