<script lang="ts">
	import type { PageData } from './$types';
	import { page } from '$app/stores';
	import { ArrowLeft } from 'lucide-svelte';

	export let data: PageData;

	$: tender = data.tender;
	$: bids = data.bids;
	$: error = data.error;
	$: tenderId = data.tenderId;

	function formatDateTime(dateTimeString: string | null | undefined) {
		if (!dateTimeString) return 'N/A';
		return new Date(dateTimeString).toLocaleString(undefined, {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatCurrency(amount: number | null | undefined, currency: string | null | undefined) {
		if (amount === null || amount === undefined) return 'N/A';
		return new Intl.NumberFormat(undefined, { style: 'currency', currency: currency || 'USD' }).format(amount);
	}
</script>

<svelte:head>
	<title>Bids for {tender ? `Tender: ${tender.title}` : `Tender ID: ${tenderId}`} - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8 px-4">
	<div class="mb-6 flex items-center justify-between">
		<div>
			<a href={`/tenders/${tenderId}`} class="btn btn-ghost btn-sm mb-2">
				<ArrowLeft class="w-4 h-4 mr-1" />
				Back to Tender Details
			</a>
			<h1 class="text-3xl font-semibold text-gray-800">
				Bids for: {tender ? tender.title : `Tender ID ${tenderId}`}
			</h1>
			{#if tender}
				<p class="text-sm text-gray-500">Review and evaluate submitted bids for this tender.</p>
			{/if}
		</div>
		<!-- Placeholder for future actions like "Start Evaluation" -->
		<!-- <button class="btn btn-primary">Start Evaluation Process</button> -->
	</div>

	{#if error}
		<div class="alert alert-error">
			<p>Error loading bids: {error}</p>
		</div>
	{:else if !bids || bids.length === 0}
		<div class="text-center py-10 bg-white shadow rounded-lg">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
				<path vector-effect="non-scaling-stroke" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
			</svg>
			<h3 class="mt-2 text-sm font-medium text-gray-900">No Bids Submitted Yet</h3>
			<p class="mt-1 text-sm text-gray-500">
				There are currently no bids submitted for this tender. Check back later or ensure suppliers have been notified.
			</p>
			<!-- <div class="mt-6">
				<button type="button" class="btn btn-outline">
					Notify Suppliers (Placeholder)
				</button>
			</div> -->
		</div>
	{:else}
		<div class="bg-white shadow-lg rounded-lg overflow-x-auto">
			<table class="min-w-full divide-y divide-gray-200">
				<thead class="bg-gray-50">
					<tr>
						<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Supplier</th>
						<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Submission Date</th>
						<th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Bid Amount</th>
						<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
						<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Validity (Days)</th>
						<th scope="col" class="relative px-6 py-3">
							<span class="sr-only">Actions</span>
						</th>
					</tr>
				</thead>
				<tbody class="bg-white divide-y divide-gray-200">
					{#each bids as bid (bid.id)}
						<tr>
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm font-medium text-gray-900">{bid.supplier_name || `Supplier ID: ${bid.supplier_id}`}</div>
								<div class="text-xs text-gray-500">ID: {bid.id}</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								{formatDateTime(bid.submission_date)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700 text-right">
								{formatCurrency(bid.bid_amount, bid.currency)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class:px-2={true} class:py-0.5={true} class:inline-flex={true} class:text-xs={true} class:leading-5={true} class:font-semibold={true} class:rounded-full={true}
								  class:bg-blue-100={bid.status === 'submitted' || bid.status === 'under_evaluation'}
								  class:text-blue-800={bid.status === 'submitted' || bid.status === 'under_evaluation'}
								  class:bg-green-100={bid.status === 'shortlisted' || bid.status === 'awarded'}
								  class:text-green-800={bid.status === 'shortlisted' || bid.status === 'awarded'}
								  class:bg-red-100={bid.status === 'rejected'}
								  class:text-red-800={bid.status === 'rejected'}
								  class:bg-gray-100={!(bid.status === 'submitted' || bid.status === 'under_evaluation' || bid.status === 'shortlisted' || bid.status === 'awarded' || bid.status === 'rejected')}
								  class:text-gray-800={!(bid.status === 'submitted' || bid.status === 'under_evaluation' || bid.status === 'shortlisted' || bid.status === 'awarded' || bid.status === 'rejected')}
								>
								  {bid.status || 'N/A'}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								{bid.validity_period_days === null || bid.validity_period_days === undefined ? 'N/A' : `${bid.validity_period_days} days`}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
								<!-- Placeholder for actions like "View Bid Details", "Evaluate Bid" -->
								<a href={`/tenders/${tenderId}/bids/${bid.id}`} class="text-indigo-600 hover:text-indigo-900 hover:underline">Details</a>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}

	<div class="mt-8">
		<a href="/tenders" class="btn btn-outline btn-sm">
			<ArrowLeft class="w-4 h-4 mr-1" /> All Tenders
		</a>
	</div>
</div>
