<script lang="ts">
	import { page } from '$app/stores';

	// Get the tender ID from the route parameters
	const tenderId = $page.params.id;

	// In a real app, you would fetch tender details and existing criteria/bids here
	const mockTender = {
		id: tenderId,
		title: `Evaluation for Tender ${tenderId}`,
	};

	// Placeholder for evaluation criteria and bid scoring data
	let evaluationCriteria = [
		{ id: 'crit1', name: 'Technical Compliance', weight: 40, description: 'Adherence to technical specifications.' },
		{ id: 'crit2', name: 'Price', weight: 30, description: 'Competitiveness of the bid price.' },
		{ id: 'crit3', name: 'Experience and Past Performance', weight: 20, description: 'Proven track record and relevant experience.' },
		{ id: 'crit4', name: 'Delivery Timeline', weight: 10, description: 'Proposed schedule for delivery of goods/services.' },
	];

	// --- New variables for adding criteria ---
	let newCriterionName: string = '';
	let newCriterionDescription: string = '';
	let newCriterionWeight: number | null = null;
	let showAddCriteriaForm = false;
	// --- End new variables ---

	// Placeholder for bids submitted to this tender
	let submittedBids = [
		{ id: 'bid1', supplierName: 'Supplier Alpha', submissionDate: '2024-06-10', status: 'Under Evaluation'},
		{ id: 'bid2', supplierName: 'Supplier Beta', submissionDate: '2024-06-12', status: 'Under Evaluation'},
		{ id: 'bid3', supplierName: 'Supplier Gamma', submissionDate: '2024-06-11', status: 'Under Evaluation'},
	];

	// --- New variables for bid evaluation ---
	let evaluatingBidId: string | null = null;
	type ScoreEntry = { score: number | null; comment: string };
	type BidScores = Record<string, Record<string, ScoreEntry>>;
	let bidScores: BidScores = {}; 
	// Initialize scores for existing bids and criteria
	submittedBids.forEach(bid => {
		if (!bidScores[bid.id]) bidScores[bid.id] = {};
		evaluationCriteria.forEach(criterion => {
			if (!bidScores[bid.id][criterion.id]) {
				bidScores[bid.id][criterion.id] = { score: null, comment: '' };
			}
		});
	});
	// --- End new variables for bid evaluation ---

	// --- New functions for criteria management ---
	function addCriterion() {
		if (newCriterionName && newCriterionDescription && newCriterionWeight !== null && newCriterionWeight > 0 && newCriterionWeight <= 100) {
			const newId = `crit${Date.now()}`;
			const newCrit = { id: newId, name: newCriterionName, description: newCriterionDescription, weight: newCriterionWeight };
			evaluationCriteria = [
				...evaluationCriteria,
				newCrit
			];
			// Add this new criterion to existing bid scores structure
			submittedBids.forEach(bid => {
				if (!bidScores[bid.id][newCrit.id]) {
					bidScores[bid.id][newCrit.id] = { score: null, comment: '' };
				}
			});
			// Reset form
			newCriterionName = '';
			newCriterionDescription = '';
			newCriterionWeight = null;
			showAddCriteriaForm = false; // Optionally hide form after adding
		} else {
			alert('Please fill all fields correctly. Weight must be between 1 and 100.');
		}
	}

	$: totalWeight = evaluationCriteria.reduce((sum, criterion) => sum + (criterion.weight || 0), 0);
	// --- End new functions ---

	// --- New functions for bid evaluation ---
	function toggleEvaluationForm(bidId: string) {
		if (evaluatingBidId === bidId) {
			evaluatingBidId = null; // Hide if already showing
		} else {
			evaluatingBidId = bidId;
			// Ensure score structure exists for this bid (might be redundant if initialized earlier, but safe)
			if (!bidScores[bidId]) bidScores[bidId] = {};
			evaluationCriteria.forEach(criterion => {
				if (!bidScores[bidId][criterion.id]) {
					bidScores[bidId][criterion.id] = { score: null, comment: '' };
				}
			});
		}
	}

	function saveBidScores(bidId: string) {
		// In a real app, this would send data to a backend
		console.log(`Scores for bid ${bidId}:`, JSON.parse(JSON.stringify(bidScores[bidId])));
		alert(`Scores for bid ${bidId} saved to console (mock).`);
		// Optionally close the form after saving
		// evaluatingBidId = null;
	}
	// --- End new functions for bid evaluation ---

</script>

<svelte:head>
	<title>Tender Evaluation: {mockTender.id} - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	<div class="mb-6 pb-4 border-b border-gray-200">
		<h1 class="text-3xl font-semibold text-gray-800">{mockTender.title}</h1>
		<p class="text-sm text-gray-500">Use this page to define evaluation criteria, score bids, and record evaluation outcomes.</p>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
		<!-- Left Column: Evaluation Criteria Setup -->
		<div class="lg:col-span-1 bg-white shadow-lg rounded-lg p-6">
			<div class="flex justify-between items-center mb-4">
				<h2 class="text-2xl font-semibold text-gray-700">Evaluation Criteria</h2>
				<button 
					on:click={() => showAddCriteriaForm = !showAddCriteriaForm} 
					class="text-sm font-medium rounded-md px-3 py-1.5 
								{showAddCriteriaForm ? 'bg-red-100 text-red-700 hover:bg-red-200' : 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200'}">
					{showAddCriteriaForm ? 'Cancel' : '+ Add Criterion'}
				</button>
			</div>

			{#if showAddCriteriaForm}
				<div class="mb-6 p-4 border border-gray-200 rounded-md bg-gray-50">
					<h3 class="text-lg font-medium text-gray-700 mb-3">Add New Criterion</h3>
					<div class="space-y-3">
						<div>
							<label for="critName" class="block text-sm font-medium text-gray-700">Name</label>
							<input type="text" id="critName" bind:value={newCriterionName} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" placeholder="e.g., Technical Score">
						</div>
						<div>
							<label for="critDesc" class="block text-sm font-medium text-gray-700">Description</label>
							<textarea id="critDesc" bind:value={newCriterionDescription} rows="3" class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" placeholder="Detailed explanation of the criterion"></textarea>
						</div>
						<div>
							<label for="critWeight" class="block text-sm font-medium text-gray-700">Weight (%)</label>
							<input type="number" id="critWeight" bind:value={newCriterionWeight} min="1" max="100" class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" placeholder="e.g., 40">
						</div>
						<button on:click={addCriterion} class="w-full px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-opacity-50">
							Save Criterion
						</button>
					</div>
				</div>
			{/if}

			<p class="text-gray-600 mb-2 text-sm">
				Current criteria for this tender. Ensure total weight equals 100%.
			</p>
			<div class="mb-3 p-2 bg-blue-50 border border-blue-200 rounded-md text-blue-700 font-medium text-center">
				Total Weight: {totalWeight}%
				{#if totalWeight !== 100}
					<span class="text-red-600 font-semibold ml-2">(Warning: Should be 100%)</span>
				{/if}
			</div>

			{#if evaluationCriteria.length > 0}
				<ul class="space-y-3 mb-4">
					{#each evaluationCriteria as criterion (criterion.id)}
						<li class="p-3 bg-gray-50 rounded-md border border-gray-200 hover:shadow-sm transition-shadow">
							<div class="flex justify-between items-start">
								<div>
									<h4 class="font-medium text-gray-700">{criterion.name} (Weight: {criterion.weight}%)</h4>
									<p class="text-sm text-gray-500">{criterion.description}</p>
								</div>
								<!-- Placeholder for Edit/Delete buttons -->
								<div class="text-xs">
									<button class="text-indigo-600 hover:text-indigo-800 mr-1">Edit</button>
									<button class="text-red-600 hover:text-red-800">Delete</button>
								</div>
							</div>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-gray-500 italic text-center py-3">No evaluation criteria defined yet. Click '+ Add Criterion' to start.</p>
			{/if}
			<!-- Removed old Manage Criteria button as it's replaced by Add Criterion toggle -->
		</div>

		<!-- Right Column: Bid Evaluation Area -->
		<div class="lg:col-span-2 bg-white shadow-lg rounded-lg p-6">
			<h2 class="text-2xl font-semibold text-gray-700 mb-4">Bid Evaluation</h2>
			<p class="text-gray-600 mb-6">
				Select a bid to score it against the defined criteria.
			</p>
			
			{#if submittedBids.length > 0}
				<h3 class="text-xl font-semibold text-gray-700 mb-3">Submitted Bids ({submittedBids.length})</h3>
				<div class="space-y-4">
					{#each submittedBids as bid (bid.id)}
						<div class="p-4 border border-gray-200 rounded-md hover:shadow-md transition-shadow 
									{evaluatingBidId === bid.id ? 'bg-indigo-50 ring-2 ring-indigo-500' : ''}">
							<div class="flex justify-between items-center">
								<h4 class="text-lg font-medium text-gray-800">{bid.supplierName}</h4>
								<span class="text-sm text-gray-500">Submitted: {bid.submissionDate}</span>
							</div>
							<p class="text-sm text-gray-600 mb-2">Status: {bid.status}</p>
							<button 
								on:click={() => toggleEvaluationForm(bid.id)}
								class="mt-1 px-3 py-1.5 text-sm rounded-md focus:outline-none focus:ring-2 focus:ring-opacity-50 
											{evaluatingBidId === bid.id ? 
												'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500' : 
												'bg-teal-500 text-white hover:bg-teal-600 focus:ring-teal-500'}">
								{evaluatingBidId === bid.id ? 'Close Evaluation' : 'Evaluate Bid'}
							</button>

							{#if evaluatingBidId === bid.id}
								<div class="mt-4 pt-4 border-t border-gray-200">
									<h5 class="text-md font-semibold text-gray-700 mb-3">Scoring for: {bid.supplierName}</h5>
									{#if evaluationCriteria.length > 0}
										<div class="space-y-4">
											{#each evaluationCriteria as criterion (criterion.id)}
												<div class="p-3 bg-gray-50 rounded-md border border-gray-100">
													<label for={`score-${bid.id}-${criterion.id}`} class="block text-sm font-medium text-gray-700">
														{criterion.name} (Weight: {criterion.weight}%)
													</label>
													<p class="text-xs text-gray-500 mb-1">{criterion.description}</p>
													<input 
														type="number" 
														id={`score-${bid.id}-${criterion.id}`} 
														bind:value={bidScores[bid.id][criterion.id].score} 
														placeholder="Score (e.g., 1-10)" 
														class="mt-1 block w-full sm:w-1/2 px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm mb-2">
													<textarea 
														id={`comment-${bid.id}-${criterion.id}`} 
														bind:value={bidScores[bid.id][criterion.id].comment} 
														rows="2" 
														placeholder="Comments/Justification for this score..." 
														class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"></textarea>
												</div>
											{/each}
											<button 
												on:click={() => saveBidScores(bid.id)}
												class="mt-3 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50">
												Save Scores for {bid.supplierName}
											</button>
										</div>
									{:else}
										<p class="text-gray-500 italic">No evaluation criteria defined. Please add criteria first.</p>
									{/if}
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{:else}
				<p class="text-gray-500 italic py-4 text-center">No bids submitted or selected for evaluation yet.</p>
			{/if}
		</div>
	</div>

	<div class="mt-8 pt-6 border-t border-gray-200 flex justify-end">
		<a href="/tenders/{mockTender.id}" class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50">
			Back to Tender Details
		</a>
	</div>
</div>
