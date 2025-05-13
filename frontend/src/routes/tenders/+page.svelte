<svelte:head>
	<title>Tenders - Procurement System</title>
	<meta name="description" content="View and manage tenders" />
</svelte:head>

<div class="container mx-auto py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-semibold">Tenders</h1>
		<a href="/tenders/new" class="bg-teal-600 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
			+ New Tender
		</a>
	</div>

	<div class="bg-white shadow-md rounded-lg p-6">
		<p class="text-gray-700">
			Tender list will be displayed here. This will include a table with columns like Tender ID, Title, Status, Type, and Closing Date.
			Users will be able to view tender details, and suppliers (once logged in) might be able to view and submit bids.
		</p>
		<p class="mt-4 text-sm text-gray-500">
			(Placeholder content - actual data and table component to be implemented)
		</p>
	</div>
</div>
