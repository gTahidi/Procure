<script lang="ts">
	import type { PageData } from './$types';
	import type { Evaluation } from '$lib/types'; 
	import { page } from '$app/stores';
	import { ArrowLeft, Edit3, CheckCircle, XCircle, Info, Star, FileText } from 'lucide-svelte';
	import { onMount } from 'svelte'; 

	export let data: PageData;

	$: tender = data.tender;
	$: bid = data.bid;
	$: initialEvaluation = data.evaluation; 
	$: error = data.error;
	$: tenderId = data.tenderId;
	$: bidId = data.bidId;

	let evaluationFormData: Partial<Evaluation> = {};
	let editEvaluationMode = false;

	onMount(() => {
		if (initialEvaluation) {
			evaluationFormData = { ...initialEvaluation };
		} else {
			evaluationFormData = {
				bid_id: bid ? bid.id : undefined, 
				score: 0,
				comments: '',
				status: '', 
				evaluation_date: new Date().toISOString()
			};
		}
	});

	function formatDateTime(dateTimeString: string | null | undefined) {
		if (!dateTimeString) return 'N/A';
		return new Date(dateTimeString).toLocaleString(undefined, {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatCurrency(amount: number | null | undefined, currency: string | null | undefined) {
		if (amount === null || amount === undefined) return 'N/A';
		return new Intl.NumberFormat(undefined, { style: 'currency', currency: currency || 'USD' }).format(amount);
	}

	async function handleSaveEvaluation() {
		if (!bid) {
			alert('Bid details are not available. Cannot save evaluation.');
			return;
		}

		if (!evaluationFormData.bid_id) evaluationFormData.bid_id = bid.id;
		if (!evaluationFormData.evaluation_date) evaluationFormData.evaluation_date = new Date().toISOString();

		console.log('Saving evaluation:', evaluationFormData);

		try {
			const method = initialEvaluation && initialEvaluation.id ? 'PUT' : 'POST';
			const endpoint = initialEvaluation && initialEvaluation.id 
				? `/api/evaluations/${initialEvaluation.id}` 
				: `/api/bids/${bid.id}/evaluations`;

			const response = await fetch(endpoint, {
				method: method,
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(evaluationFormData)
			});

			if (response.ok) {
				const savedEvaluation: Evaluation = await response.json();
				initialEvaluation = savedEvaluation; 
				evaluationFormData = { ...savedEvaluation };
				editEvaluationMode = false;
				alert('Evaluation saved successfully!');
			} else {
				const errorResult = await response.json().catch(() => ({ message: 'Failed to save evaluation.' }));
				alert(`Error saving evaluation: ${errorResult.message || response.statusText}`);
			}
		} catch (err: any) {
			console.error('Error in handleSaveEvaluation:', err);
			alert(`An unexpected error occurred: ${err.message}`);
		}
	}

	function startEditMode() {
		if (initialEvaluation) {
			evaluationFormData = { ...initialEvaluation };
		} else {
			evaluationFormData = {
				bid_id: bid ? bid.id : undefined,
				score: 0,
				comments: '',
				status: '',
				evaluation_date: new Date().toISOString()
			};
		}
		editEvaluationMode = true;
	}

</script>

<svelte:head>
	<title>
		{bid ? `Bid Details: ${bid.supplier_name || 'Supplier ' + bid.supplier_id}` : 'Bid Details'}
		{tender ? ` for Tender: ${tender.title}` : ''}
		- Procurement System
	</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	<div class="mb-6">
		<a href={`/tenders/${tenderId}/bids`} class="btn btn-ghost btn-sm mb-2">
			<ArrowLeft class="w-4 h-4 mr-1" />
			Back to All Bids for this Tender
		</a>
		{#if tender}
			<p class="text-sm text-gray-500">Tender: <a href={`/tenders/${tenderId}`} class="text-indigo-600 hover:underline">{tender.title}</a> (ID: {tenderId})</p>
		{/if}
	</div>

	{#if error}
		<div class="alert alert-error">
			<p>Error loading bid details: {error}</p>
		</div>
	{:else if !bid}
		<div class="text-center py-10">
			<Info class="mx-auto h-12 w-12 text-gray-400" />
			<h3 class="mt-2 text-sm font-medium text-gray-900">Bid Not Found</h3>
			<p class="mt-1 text-sm text-gray-500">The requested bid (ID: {bidId}) could not be found.</p>
		</div>
	{:else}
		<div class="bg-white shadow-xl rounded-lg p-6 md:p-8">
			<div class="md:flex md:justify-between md:items-start mb-6 pb-6 border-b">
				<div>
					<h1 class="text-2xl md:text-3xl font-bold text-gray-800">Bid from: {bid.supplier_name || `Supplier ID: ${bid.supplier_id}`}</h1>
					<p class="text-sm text-gray-500">Bid ID: {bid.id} | Submitted: {formatDateTime(bid.submission_date)}</p>
				</div>
				<div class="mt-4 md:mt-0">
					<span class:px-3={true} class:py-1={true} class:text-sm={true} class:font-semibold={true} class:rounded-full={true}
					  class:bg-blue-100={bid.status === 'submitted' || bid.status === 'under_evaluation'}
					  class:text-blue-800={bid.status === 'submitted' || bid.status === 'under_evaluation'}
					  class:bg-green-100={bid.status === 'shortlisted' || bid.status === 'awarded'}
					  class:text-green-800={bid.status === 'shortlisted' || bid.status === 'awarded'}
					  class:bg-red-100={bid.status === 'rejected'}
					  class:text-red-800={bid.status === 'rejected'}
					  class:bg-gray-200={!(bid.status === 'submitted' || bid.status === 'under_evaluation' || bid.status === 'shortlisted' || bid.status === 'awarded' || bid.status === 'rejected')}
					  class:text-gray-800={!(bid.status === 'submitted' || bid.status === 'under_evaluation' || bid.status === 'shortlisted' || bid.status === 'awarded' || bid.status === 'rejected')}
					>
					  Status: {bid.status || 'N/A'}
					</span>
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
				<div>
					<h3 class="text-md font-semibold text-gray-700">Bid Amount</h3>
					<p class="text-2xl font-light text-gray-900">{formatCurrency(bid.bid_amount, bid.currency)}</p>
				</div>
				<div>
					<h3 class="text-md font-semibold text-gray-700">Validity Period</h3>
					<p class="text-lg text-gray-600">{bid.validity_period_days === null || bid.validity_period_days === undefined ? 'N/A' : `${bid.validity_period_days} days`}</p>
				</div>
				<div>
					<h3 class="text-md font-semibold text-gray-700">Last Updated</h3>
					<p class="text-lg text-gray-600">{formatDateTime(bid.updated_at)}</p>
				</div>
			</div>

			{#if bid.notes}
				<div class="mb-8 p-4 bg-gray-50 rounded-md">
					<h3 class="text-lg font-semibold text-gray-700 mb-2">Supplier Notes</h3>
					<p class="text-gray-600 whitespace-pre-line">{bid.notes}</p>
				</div>
			{/if}

			<!-- Bid Documents Section (Placeholder) -->
			<div class="mb-8">
				<h3 class="text-xl font-semibold text-gray-700 mb-3"><FileText class="inline-block w-5 h-5 mr-2 align-text-bottom" />Bid Documents</h3>
				<p class="text-gray-500 italic">Document listing functionality to be implemented. (e.g., links to download submitted documents)</p>
			</div>

			<!-- Evaluation Section -->
			<div class="mt-10 pt-8 border-t-2 border-indigo-500">
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-2xl font-semibold text-indigo-700">Bid Evaluation</h2>
					{#if !editEvaluationMode}
						<button on:click={startEditMode} class="btn btn-outline btn-primary btn-sm">
							<Edit3 class="w-4 h-4 mr-2" /> {initialEvaluation ? 'Edit Evaluation' : 'Evaluate Bid'}
						</button>
					{/if}
				</div>

				{#if editEvaluationMode}
					<form on:submit|preventDefault={handleSaveEvaluation} class="space-y-6 bg-indigo-50 p-6 rounded-lg shadow">
						<div>
							<label for="eval-score" class="block text-sm font-medium text-gray-700 mb-1">Overall Score (0-5)</label>
							<div class="flex items-center">
								{#each Array(5) as _, i}
									<button type="button" on:click={() => evaluationFormData.score = i + 1} class="p-1 focus:outline-none">
										<Star class={`w-6 h-6 ${evaluationFormData.score && evaluationFormData.score > i ? 'text-yellow-400 fill-current' : 'text-gray-300'}`} />
									</button>
								{/each}
							</div>
							<input type="hidden" bind:value={evaluationFormData.score} />
						</div>

						<div>
							<label for="eval-comments" class="block text-sm font-medium text-gray-700">Evaluation Comments</label>
							<textarea id="eval-comments" rows="4" bind:value={evaluationFormData.comments} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 sm:text-sm" placeholder="Strengths, weaknesses, compliance, etc."></textarea>
						</div>
						
						<div>
							<label for="eval-status" class="block text-sm font-medium text-gray-700">Evaluation Status / Recommendation</label>
							<select id="eval-status" bind:value={evaluationFormData.status} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 sm:text-sm">
								<option value="">Select Status...</option>
								<option value="under_review">Under Review</option>
								<option value="shortlisted">Shortlist</option>
								<option value="rejected_technical">Reject (Technical)</option>
								<option value="rejected_commercial">Reject (Commercial)</option>
								<option value="needs_clarification">Needs Further Clarification</option>
								<option value="recommended_for_award">Recommend for Award</option>
								<option value="not_recommended">Not Recommended</option>
							</select>
						</div>

						<div class="flex justify-end space-x-3 pt-4">
							<button type="button" on:click={() => editEvaluationMode = false} class="btn btn-ghost">Cancel</button>
							<button type="submit" class="btn btn-primary">
								<CheckCircle class="w-4 h-4 mr-2" /> Save Evaluation
							</button>
						</div>
					</form>
				{:else if initialEvaluation && initialEvaluation.id } 
					<div class="bg-gray-50 p-6 rounded-md shadow-sm border border-gray-200">
						<h4 class="text-lg font-semibold text-gray-700 mb-3">Current Evaluation</h4>
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
							<p class="text-sm text-gray-600"><strong>Status:</strong> <span class="font-medium text-gray-800">{initialEvaluation.status || 'N/A'}</span></p>
							<p class="text-sm text-gray-600"><strong>Score:</strong> <span class="font-medium text-gray-800">{initialEvaluation.score !== null && initialEvaluation.score !== undefined ? initialEvaluation.score + ' / 5' : 'Not scored'}</span></p>
							<p class="text-sm text-gray-600 col-span-1 sm:col-span-2"><strong>Evaluator ID:</strong> <span class="font-medium text-gray-800">{initialEvaluation.evaluator_id || 'N/A'}</span></p>
							<p class="text-sm text-gray-600 col-span-1 sm:col-span-2"><strong>Evaluation Date:</strong> <span class="font-medium text-gray-800">{formatDateTime(initialEvaluation.evaluation_date)}</span></p>
						</div>
						{#if initialEvaluation.comments}
							<div class="mt-4 pt-4 border-t border-gray-200">
								<p class="text-sm font-semibold text-gray-700 mb-1">Comments:</p>
								<p class="text-sm text-gray-600 whitespace-pre-line bg-white p-3 rounded">{initialEvaluation.comments}</p>
							</div>
						{/if}
					</div>
				{:else}
					<p class="text-center text-gray-500 py-6 italic">This bid has not been evaluated yet.</p>
				{/if}
			</div>
		</div>
	{/if}
</div>
