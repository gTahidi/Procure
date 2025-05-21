<svelte:head>
	<title>Requisitions - Procurement System</title>
	<meta name="description" content="View and manage requisitions" />
</svelte:head>

<div class="container mx-auto py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-semibold">Requisitions</h1>
		<a href="/requisitions/new" class="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
			+ New Requisition
		</a>
	</div>

	<div class="bg-white shadow-md rounded-lg overflow-x-auto">
		<table class="min-w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Requester</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created On</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
				</tr>
			</thead>
			<tbody class="bg-white divide-y divide-gray-200">
				{#if data.error}
					<tr>
						<td colspan="7" class="px-6 py-12 text-center text-sm text-red-500">
							Error loading requisitions: {data.error}
						</td>
					</tr>
				{:else if data.requisitions && data.requisitions.length > 0}
					{#each data.requisitions as req (req.id)}
						<tr>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-indigo-600 hover:text-indigo-900">
								<a href={req.detailLink}>{req.id}</a>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{req.title}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{req.requester}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{req.type}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm">
								<span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
									{req.status === 'Approved' ? 'bg-green-100 text-green-800' : 
									req.status === 'Pending Approval' ? 'bg-yellow-100 text-yellow-800' : 
									req.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
									req.status === 'Rejected' ? 'bg-red-100 text-red-800' : 
									'bg-gray-100 text-gray-800'}">
									{req.status}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{req.creationDate}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
								<a href={req.detailLink} class="text-indigo-600 hover:text-indigo-900">View</a>
								<!-- Add other actions like Edit/Delete based on status/permissions later -->
							</td>
						</tr>
					{/each}
				{:else}
					<tr>
						<td colspan="7" class="px-6 py-12 text-center text-sm text-gray-500">
							No requisitions found.
						</td>
					</tr>
				{/if}
			</tbody>
		</table>
	</div>
</div>

<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;
</script>
