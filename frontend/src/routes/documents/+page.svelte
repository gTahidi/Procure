<script lang="ts">
	// Mock data for demonstration
	// In a real application, this data would come from an API
	const mockDocuments = [
		{ id: 'doc1', filename: 'Quotation_ABC_Corp.pdf', documentType: 'Quotation', uploadDate: '2024-05-10', associatedEntity: 'Requisition #R2024-001', entityLink: '/requisitions/R2024-001' },
		{ id: 'doc2', filename: 'Technical_Specs_ProjectX.docx', documentType: 'Specification Sheet', uploadDate: '2024-05-09', associatedEntity: 'Tender #T2024-005', entityLink: '/tenders/T2024-005' },
		{ id: 'doc3', filename: 'Contract_Draft_SupplierZ.pdf', documentType: 'Contract Draft', uploadDate: '2024-05-11', associatedEntity: 'Purchase Order #PO2024-015', entityLink: '/purchase-orders/PO2024-015' },
		{ id: 'doc4', filename: 'User_Manual_EquipmentA.pdf', documentType: 'Supporting Document', uploadDate: '2024-05-08', associatedEntity: 'Asset #ASSET-007', entityLink: '/assets/ASSET-007' },
		{ id: 'doc5', filename: 'Invoice_INV-00123.pdf', documentType: 'Other', uploadDate: '2024-05-12', associatedEntity: 'Requisition #R2024-002', entityLink: '/requisitions/R2024-002' },
	];

	let searchTerm = '';
	let filterType = '';

	// Placeholder for filtered documents logic
	$: filteredDocuments = mockDocuments.filter(doc => {
		const matchesSearch = doc.filename.toLowerCase().includes(searchTerm.toLowerCase()) || 
						  doc.associatedEntity.toLowerCase().includes(searchTerm.toLowerCase());
		const matchesType = filterType ? doc.documentType === filterType : true;
		return matchesSearch && matchesType;
	});

	// Unique document types for filter dropdown
	const uniqueDocumentTypes = ['', ...new Set(mockDocuments.map(doc => doc.documentType))];

</script>

<svelte:head>
	<title>Document Repository - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	<h1 class="text-3xl font-semibold mb-8 text-gray-800">Document Repository</h1>

	<!-- Search and Filter Placeholder -->
	<div class="mb-6 p-4 bg-gray-50 rounded-lg shadow">
		<h2 class="text-xl font-medium text-gray-700 mb-3">Filter & Search</h2>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div>
				<label for="searchTerm" class="block text-sm font-medium text-gray-700">Search Documents</label>
				<input type="text" id="searchTerm" bind:value={searchTerm} placeholder="Search by filename or entity..." class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2">
			</div>
			<div>
				<label for="filterType" class="block text-sm font-medium text-gray-700">Filter by Type</label>
				<select id="filterType" bind:value={filterType} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2">
					{#each uniqueDocumentTypes as type}
						<option value={type}>{type === '' ? 'All Types' : type}</option>
					{/each}
				</select>
			</div>
		</div>
	</div>

	<!-- Documents Table -->
	<div class="bg-white shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Filename</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Document Type</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Upload Date</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Associated Entity</th>
					<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
				</tr>
			</thead>
			<tbody class="bg-white divide-y divide-gray-200">
				{#if filteredDocuments.length > 0}
					{#each filteredDocuments as doc (doc.id)}
						<tr>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
								<!-- Placeholder for file icon -->
								<span class="mr-2">📄</span> {doc.filename}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{doc.documentType}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{doc.uploadDate}</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								<a href={doc.entityLink || '#'} class="text-indigo-600 hover:text-indigo-900 hover:underline">{doc.associatedEntity}</a>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
								<a href="#" class="text-indigo-600 hover:text-indigo-900 mr-3">View</a>
								<a href="#" class="text-red-600 hover:text-red-900">Delete</a> 
							</td>
						</tr>
					{/each}
				{:else}
					<tr>
						<td colspan="5" class="px-6 py-12 text-center text-sm text-gray-500">
							No documents found matching your criteria.
						</td>
					</tr>
				{/if}
			</tbody>
		</table>
	</div>

	<!-- Pagination Placeholder -->
	<div class="mt-6 text-center">
		<p class="text-sm text-gray-500 italic">(Pagination placeholder)</p>
	</div>
</div>
