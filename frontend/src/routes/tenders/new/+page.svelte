<script lang="ts">
	import DocumentUpload from '$lib/components/DocumentUpload.svelte';

	let tender = {
		tenderNumber: '', // Could be auto-generated or manually entered
		title: '',
		description: '',
		tenderType: 'open', // e.g., 'open', 'restricted', 'direct'
		tenderCategory: 'goods', // e.g., 'goods', 'services', 'works', 'consultancy'
		issueDate: '',
		closingDate: '',
		clarificationEndDate: '',
		bidOpeningDate: '',
		bidValidityPeriod: 90, // days
		procurementMethodDetails: '', // Further details if needed, e.g. for restricted/direct
		evaluationCriteria: [{ criterion: '', weight: null as number | null, details: '' }],
		bidSubmissionGuidelines: '',
		contactPerson: '',
		contactEmail: '',
		contactPhone: '',
		notes: '',
		status: 'draft' // Initial status: 'draft', 'published', 'closed', 'evaluating', 'awarded', 'cancelled'
	};

	let attachedFiles: any[] = [];

	function handleFilesAttached(event: CustomEvent) {
		attachedFiles = [...attachedFiles, ...event.detail];
		console.log('Files for tender:', attachedFiles);
	}

	function addCriterion() {
		tender.evaluationCriteria = [
			...tender.evaluationCriteria,
			{ criterion: '', weight: null as number | null, details: '' }
		];
	}

	function removeCriterion(index: number) {
		tender.evaluationCriteria = tender.evaluationCriteria.filter((_, i) => i !== index);
	}

	function handleSubmit(isDraft: boolean) {
		tender.status = isDraft ? 'draft' : 'published';
		console.log('Submitting Tender:', tender, 'Attached Files:', attachedFiles);
		// API call to save/publish tender would go here
		alert(`Tender ${isDraft ? 'saved as draft' : 'published'}. Check console.`);
		// Potentially navigate away or reset form
	}

</script>

<svelte:head>
	<title>New Tender Notice - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4 max-w-4xl">
	<h1 class="text-3xl font-semibold mb-8 text-gray-800">Create New Tender Notice</h1>

	<form on:submit|preventDefault class="space-y-8 bg-white p-8 rounded-lg shadow-lg">

		<!-- Section 1: Tender Details -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Tender Details</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div>
					<label for="tenderNumber" class="block text-sm font-medium text-gray-700">Tender Number</label>
					<input type="text" id="tenderNumber" bind:value={tender.tenderNumber} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., PRQS/TN/2024/001">
				</div>
				<div>
					<label for="title" class="block text-sm font-medium text-gray-700">Tender Title</label>
					<input type="text" id="title" bind:value={tender.title} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., Supply of Stationery">
				</div>
				<div>
					<label for="tenderType" class="block text-sm font-medium text-gray-700">Tender Type</label>
					<select id="tenderType" bind:value={tender.tenderType} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						<option value="open">Open Tender</option>
						<option value="restricted">Restricted Tender</option>
						<option value="direct">Direct Procurement</option>
						<option value="quotation">Request for Quotation</option>
					</select>
				</div>
				<div>
					<label for="tenderCategory" class="block text-sm font-medium text-gray-700">Category</label>
					<select id="tenderCategory" bind:value={tender.tenderCategory} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
						<option value="goods">Goods</option>
						<option value="services">Services</option>
						<option value="works">Works</option>
						<option value="consultancy">Consultancy</option>
					</select>
				</div>
				<div>
					<label for="issueDate" class="block text-sm font-medium text-gray-700">Issue Date</label>
					<input type="date" id="issueDate" bind:value={tender.issueDate} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
				</div>
				<div>
					<label for="closingDate" class="block text-sm font-medium text-gray-700">Closing Date</label>
					<input type="datetime-local" id="closingDate" bind:value={tender.closingDate} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
				</div>
                <div>
					<label for="clarificationEndDate" class="block text-sm font-medium text-gray-700">Clarification End Date</label>
					<input type="datetime-local" id="clarificationEndDate" bind:value={tender.clarificationEndDate} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
				</div>
                <div>
					<label for="bidOpeningDate" class="block text-sm font-medium text-gray-700">Bid Opening Date</label>
					<input type="datetime-local" id="bidOpeningDate" bind:value={tender.bidOpeningDate} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
				</div>
                <div>
                    <label for="bidValidityPeriod" class="block text-sm font-medium text-gray-700">Bid Validity Period (days)</label>
                    <input type="number" min="1" id="bidValidityPeriod" bind:value={tender.bidValidityPeriod} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., 90">
                </div>
			</div>
			<div class="mt-6">
				<label for="description" class="block text-sm font-medium text-gray-700">Detailed Description</label>
				<textarea id="description" rows="4" bind:value={tender.description} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Provide a detailed description of the tender requirements..."></textarea>
			</div>
            <div class="mt-6">
				<label for="procurementMethodDetails" class="block text-sm font-medium text-gray-700">Procurement Method Details (if restricted/direct)</label>
				<textarea id="procurementMethodDetails" rows="2" bind:value={tender.procurementMethodDetails} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Justification or list of prequalified suppliers..."></textarea>
			</div>
		</section>

		<!-- Section 2: Evaluation Criteria -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Evaluation Criteria</h2>
			{#each tender.evaluationCriteria as criterionItem, i}
				<div class="p-4 border border-gray-200 rounded-md mb-4 space-y-3 relative bg-gray-50">
					<p class="font-medium text-gray-600">Criterion #{i + 1}</p>
					<div>
						<label for={`criterion_name_${i}`} class="block text-sm font-medium text-gray-700">Criterion</label>
						<input type="text" id={`criterion_name_${i}`} bind:value={criterionItem.criterion} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., Technical Compliance">
					</div>
					<div>
						<label for={`criterion_weight_${i}`} class="block text-sm font-medium text-gray-700">Weight (%)</label>
						<input type="number" min="0" max="100" step="0.01" id={`criterion_weight_${i}`} bind:value={criterionItem.weight} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., 40">
					</div>
                    <div>
						<label for={`criterion_details_${i}`} class="block text-sm font-medium text-gray-700">Details / Scoring Method</label>
						<textarea id={`criterion_details_${i}`} rows="2" bind:value={criterionItem.details} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Describe how this criterion will be evaluated..."></textarea>
					</div>
					{#if tender.evaluationCriteria.length > 1}
						<button type="button" on:click={() => removeCriterion(i)} aria-label="Remove criterion" class="absolute top-2 right-2 text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-100">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>
						</button>
					{/if}
				</div>
			{/each}
			<button type="button" on:click={addCriterion} class="mt-2 text-sm text-indigo-600 hover:text-indigo-800 font-medium flex items-center">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" /></svg>
				Add Evaluation Criterion
			</button>
		</section>

		<!-- Section 3: Bid Submission Guidelines -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Bid Submission Guidelines</h2>
			<textarea id="bidSubmissionGuidelines" rows="5" bind:value={tender.bidSubmissionGuidelines} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Provide clear instructions for bidders on how to prepare and submit their bids, including format, deadlines, and submission portal/address..."></textarea>
		</section>

		<!-- Section 4: Contact Information -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Contact Information</h2>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div>
					<label for="contactPerson" class="block text-sm font-medium text-gray-700">Contact Person</label>
					<input type="text" id="contactPerson" bind:value={tender.contactPerson} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., Jane Doe">
				</div>
				<div>
					<label for="contactEmail" class="block text-sm font-medium text-gray-700">Contact Email</label>
					<input type="email" id="contactEmail" bind:value={tender.contactEmail} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., tenders@example.com">
				</div>
				<div>
					<label for="contactPhone" class="block text-sm font-medium text-gray-700">Contact Phone</label>
					<input type="tel" id="contactPhone" bind:value={tender.contactPhone} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g., +254 700 000 000">
				</div>
			</div>
		</section>

		<!-- Section 5: Supporting Tender Documents -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Tender Documents & Addenda</h2>
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

		<!-- Section 6: Additional Notes -->
		<section>
			<h2 class="text-xl font-semibold text-gray-700 mb-4 border-b pb-2">Additional Notes</h2>
			<textarea id="notes" rows="3" bind:value={tender.notes} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="Any other relevant information or notes for this tender..."></textarea>
		</section>

		<!-- Actions -->
		<div class="pt-6 border-t mt-4 flex justify-end space-x-3">
			<button type="button" on:click={() => handleSubmit(true)} class="bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium py-2 px-4 rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500">
				Save as Draft
			</button>
			<button type="button" on:click={() => handleSubmit(false)} class="bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
				Publish Tender
			</button>
		</div>
	</form>
</div>
