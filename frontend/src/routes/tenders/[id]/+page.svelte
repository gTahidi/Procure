<script lang="ts">
	import { page } from '$app/stores';

	// Get the tender ID from the route parameters
	const tenderId = $page.params.id;

	// Mock data - in a real app, you'd fetch this based on tenderId
	const mockTenderDetails = {
		id: tenderId,
		title: `Details for Tender ${tenderId}`,
		description: 'This is a detailed description of the tender. It would include scope of work/goods, evaluation criteria, submission guidelines, and important dates.',
		type: 'Service',
		status: 'Open',
		issueDate: '2024-05-01',
		closingDate: '2024-06-15',
		documents: [
			{ name: 'RFP_Document.pdf', url: '#' },
			{ name: 'Technical_Specifications.docx', url: '#' },
			{ name: 'Draft_Contract.pdf', url: '#' }
		],
		clarificationDeadline: '2024-05-30',
		contactPerson: 'proc.officer@example.com'
	};

	const tender = mockTenderDetails;

</script>

<svelte:head>
	<title>Tender {tender.id} - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	<div class="bg-white shadow-lg rounded-lg p-6 md:p-8">
		<div class="mb-6 pb-4 border-b border-gray-200">
			<h1 class="text-3xl font-semibold text-gray-800">{tender.title}</h1>
			<p class="text-sm text-gray-500">Issued on: {tender.issueDate} | Closes on: {tender.closingDate}</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Tender ID</h3>
				<p class="text-gray-600">{tender.id}</p>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Status</h3>
				<span class={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full 
								${tender.status === 'Open' ? 'bg-blue-100 text-blue-800' : 
								 tender.status === 'Closed' ? 'bg-gray-100 text-gray-800' : 
								 tender.status === 'Awarded' ? 'bg-green-100 text-green-800' : 
								 tender.status === 'Evaluation' ? 'bg-yellow-100 text-yellow-800' :
								'bg-purple-100 text-purple-800'}`}>
					{tender.status}
				</span>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Type</h3>
				<p class="text-gray-600">{tender.type}</p>
			</div>
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Description</h3>
			<p class="text-gray-600 leading-relaxed whitespace-pre-line">{tender.description}</p>
		</div>

		<div class="mb-8">
			<h3 class="text-xl font-semibold text-gray-700 mb-3">Tender Documents</h3>
			{#if tender.documents && tender.documents.length > 0}
				<ul class="list-disc list-inside space-y-1">
					{#each tender.documents as doc (doc.name)}
						<li>
							<a href={doc.url} class="text-teal-600 hover:text-teal-800 hover:underline" target="_blank" rel="noopener noreferrer">
								{doc.name}
							</a>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-gray-500 italic">No documents provided for this tender.</p>
			{/if}
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Clarification Deadline</h3>
				<p class="text-gray-600">{tender.clarificationDeadline || 'N/A'}</p>
			</div>
			<div>
				<h3 class="text-lg font-medium text-gray-700 mb-1">Contact Person</h3>
				<p class="text-gray-600"><a href="mailto:{tender.contactPerson}" class="text-teal-600 hover:underline">{tender.contactPerson || 'N/A'}</a></p>
			</div>
		</div>

		<div class="mt-8 pt-6 border-t border-gray-200 flex justify-end space-x-3">
			<!-- Placeholder for action buttons like Submit Bid, Ask Clarification, etc. -->
			<button class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50">
				Back to List
			</button>
			<button class="px-4 py-2 bg-teal-600 text-white rounded-md hover:bg-teal-700 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-opacity-50">
				Submit Bid (Placeholder)
			</button>
		</div>
	</div>
</div>
