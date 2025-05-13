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

	<div class="bg-white shadow-md rounded-lg p-6">
		<p class="text-gray-700">
			Requisition list will be displayed here. This will include a table with columns like ID, Requester, Type, Status, and Creation Date. 
			Users will be able to click on a requisition to view its details or initiate actions based on their role and the requisition's status.
		</p>
		<p class="mt-4 text-sm text-gray-500">
			(Placeholder content - actual data and table component to be implemented)
		</p>
	</div>
</div>
