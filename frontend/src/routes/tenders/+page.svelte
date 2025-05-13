<svelte:head>
	<title>Tenders - Procurement System</title>
	<meta name="description" content="View and manage tenders" />
</svelte:head>

<script lang="ts">
	// Mock data for tenders
	const mockTenders = [
		{ id: 'TEN-001', title: 'Office Cleaning Services 2025', type: 'Service', status: 'Open', closingDate: '2024-06-15', detailLink: '/tenders/TEN-001' },
		{ id: 'TEN-002', title: 'Supply of IT Equipment', type: 'Goods', status: 'Closed', closingDate: '2024-05-01', detailLink: '/tenders/TEN-002' },
		{ id: 'TEN-003', title: 'Construction of New Warehouse Wing', type: 'Works', status: 'Awarded', closingDate: '2024-04-20', detailLink: '/tenders/TEN-003' },
		{ id: 'TEN-004', title: 'Consultancy for Market Research', type: 'Consultancy', status: 'Open', closingDate: '2024-07-01', detailLink: '/tenders/TEN-004' },
		{ id: 'TEN-005', title: 'Security Services Contract', type: 'Service', status: 'Evaluation', closingDate: '2024-05-25', detailLink: '/tenders/TEN-005' },
	];

	// TODO: Add search/filter logic if needed in the future
	let filteredTenders = mockTenders;
</script>

<div class="container mx-auto py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-semibold">Tenders</h1>
		<a href="/tenders/new" class="bg-teal-600 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
			+ New Tender
		</a>
	</div>

	<div class="bg-white shadow-md rounded-lg overflow-x-auto">
		<table class="min-w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Closing Date</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
				</tr>
			</thead>
			<tbody class="bg-white divide-y divide-gray-200">
				{#if filteredTenders.length > 0}
					{#each filteredTenders as tender (tender.id)}
						<tr>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-teal-600 hover:text-teal-900">
								<a href={tender.detailLink}>{tender.id}</a>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{tender.title}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{tender.type}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm">
								<span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
									{tender.status === 'Open' ? 'bg-blue-100 text-blue-800' : 
									tender.status === 'Closed' ? 'bg-gray-100 text-gray-800' : 
									tender.status === 'Awarded' ? 'bg-green-100 text-green-800' : 
									tender.status === 'Evaluation' ? 'bg-yellow-100 text-yellow-800' :
									'bg-purple-100 text-purple-800'}">
									{tender.status}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{tender.closingDate}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
								<a href={tender.detailLink} class="text-teal-600 hover:text-teal-900">View</a>
								<!-- Add other actions like Bid/Manage based on status/permissions later -->
							</td>
						</tr>
					{/each}
				{:else}
					<tr>
						<td colspan="6" class="px-6 py-12 text-center text-sm text-gray-500">
							No tenders found.
						</td>
					</tr>
				{/if}
			</tbody>
		</table>
	</div>
</div>
