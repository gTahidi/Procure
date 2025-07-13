<script lang="ts">
	import type { PageData } from './$types';
	import { page } from '$app/stores';
	import { ArrowLeft } from 'lucide-svelte';

	export let data: PageData;

	$: ({ tender, bids } = data);

	// Helper function to format dates
	const formatDate = (dateString: string | null | undefined) => {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-GB', {
			day: '2-digit',
			month: 'short',
			year: 'numeric'
		});
	};

	// Helper to format currency
	const formatCurrency = (amount: number | null | undefined) => {
		if (amount === null || typeof amount === 'undefined') return 'N/A';
		return new Intl.NumberFormat('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(amount);
	};

	// Prepare data for the evaluation table
	$: items = tender?.requisition?.items || [];
	$: vendors = bids?.map(bid => ({ id: bid.supplier_id, name: `Vendor ${bid.supplier_id}`, bid_id: bid.id })) || [];

	// A map for quick lookup: vendorId -> bid
	$: bidsByVendor = new Map(bids?.map(bid => [bid.supplier_id, bid]));

	// A map for quick lookup: requisition_item_id -> { vendorId: bidItem }
	$: bidItemsByReqItem = items.reduce((acc, item) => {
		acc.set(item.id, new Map());
		return acc;
	}, new Map());

	// Calculate the total for each vendor's bid by summing up item totals
	$: vendorTotals = vendors.reduce((acc, vendor) => {
		const total = items.reduce((sum, item) => {
			const bidItem = bidItemsByReqItem.get(item.id)?.get(vendor.id);
			const price = bidItem?.offered_unit_price || 0;
			const quantity = item.quantity;
			return sum + price * quantity;
		}, 0);
		acc.set(vendor.id, total);
		return acc;
	}, new Map());

	$: {
		if (bids) {
			bids.forEach(bid => {
				bid.items?.forEach(bidItem => {
					if (bidItemsByReqItem.has(bidItem.requisition_item_id)) {
						bidItemsByReqItem.get(bidItem.requisition_item_id).set(bid.supplier_id, bidItem);
					}
				});
			});
		}
	}

</script>

<svelte:head>
    <title>Tender Evaluation - {tender?.title || 'Tender'}</title>
</svelte:head>

<div class="container mx-auto p-4 sm:p-6 lg:p-8 bg-base-200/50">
	<a href={`/tenders/${tender?.id}`} class="btn btn-ghost mb-6">
		<ArrowLeft class="mr-2 h-4 w-4" />
		Back to Tender
	</a>

    {#if tender && bids}
	<div class="bg-white shadow-lg rounded-lg p-6 print:shadow-none">
		<div class="text-center mb-4">
			<h1 class="text-2xl font-bold text-primary">PROCUREMENT EVALUATION SHEET</h1>
			<h2 class="text-lg font-semibold">TENDER NO: {tender?.id} - {tender?.title}</h2>
		</div>

		<!-- Tender Information -->
		<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6 text-sm border p-4 rounded-md">
			<div><strong>Invitation Issued:</strong> {formatDate(tender?.published_date)}</div>
			<div><strong>Closing Date:</strong> {formatDate(tender?.closing_date)}</div>
			<div><strong>Invited Bidders:</strong> {tender?.bidders_invited_count || 'N/A'}</div>
			<div><strong>Received Bids:</strong> {bids?.length || 0}</div>
		</div>

		<!-- Evaluation of Offers -->
		<h3 class="text-xl font-bold text-center my-4">Evaluation of the Offers</h3>
		<div class="overflow-x-auto">
			<table class="table table-bordered table-compact w-full text-xs">
				<thead>
					<tr class="text-center bg-base-200">
						<th class="w-8 align-bottom">No.</th>
						<th class="text-left align-bottom">Item Description</th>
						<th class="align-bottom">Quantity</th>
						{#each vendors as vendor}
							<th colspan="2" class="border-l-2 border-base-300">{vendor.name}</th>
						{/each}
					</tr>
					<tr class="text-center bg-base-200">
						<th colspan="3"></th>
						{#each vendors as _}
							<th class="border-l-2 border-base-300">Unit Price (TZS)</th>
							<th>Total Price (TZS)</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each items as item, i}
						<tr>
							<td class="text-center">{i + 1}</td>
							<td>{item.description}</td>
							<td class="text-center">{item.quantity}</td>
							{#each vendors as vendor}
								{@const bidItem = bidItemsByReqItem.get(item.id)?.get(vendor.id)}
								<td class="text-right border-l-2 border-base-300">{formatCurrency(bidItem?.offered_unit_price)}</td>
								<td class="text-right font-semibold">{formatCurrency((bidItem?.offered_unit_price || 0) * item.quantity)}</td>
							{/each}
						</tr>
					{/each}
				</tbody>
				<tfoot>
					<tr class="font-bold bg-base-200">
						<td colspan="3" class="text-right">Sub-Total before VAT (TZS)</td>
						{#each vendors as vendor}
							<td colspan="2" class="text-right border-l-2 border-base-300">{formatCurrency(vendorTotals.get(vendor.id))}</td>
						{/each}
					</tr>
					<!-- Add rows for VAT, Grand Total, etc. as needed -->
				</tfoot>
			</table>
		</div>

		<!-- Procurement Recommendations -->
		<div class="mt-8">
			<h3 class="text-xl font-bold mb-2">PROCUREMENT RECOMMENDATIONS</h3>
			<div class="prose max-w-none text-sm border p-4 rounded-md bg-base-100">
				<p>Procurement received a request from IT unit for purchase of items as per above listed (PR No. {tender?.requisition?.id}).</p>
				<p>Out of the {tender?.bidders_invited_count || 'N/A'} invited bidders, {bids?.length || 0} vendors submitted their quotations by the closure date.</p>
				<p>After due review of all the offers received as per specs of the tender, TZO Procurement Unit recommends tender award to the best competitive offer.</p>
			</div>
		</div>

		<!-- Recommendation for Award -->
		<div class="mt-8">
			<h3 class="text-xl font-bold mb-2">Recommendation for Award</h3>
			<div class="overflow-x-auto">
				<table class="table table-bordered w-full">
					<thead class="bg-base-200">
						<tr>
							<th>No.</th>
							<th>Vendor Name</th>
							<th>PO Amount (TZS)</th>
						</tr>
					</thead>
					<tbody>
						<!-- This would be populated based on evaluation logic -->
						<tr>
							<td>1</td>
							<td>(Recommended Vendor)</td>
							<td>(Recommended Amount)</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

	</div>
    {:else}
        <div class="text-center">
            <p class="text-xl">Loading evaluation data or data not available...</p>
        </div>
    {/if}
</div>

<style>
	.table-bordered th,
	.table-bordered td {
		border: 1px solid hsl(var(--b2, var(--b1)) / 0.2);
	}
    .table-compact th, .table-compact td {
        padding: 0.5rem;
    }
    @media print {
        .container, .container * {
            visibility: visible;
        }
        .container {
            position: absolute;
            left: 0;
            top: 0;
            width: 100%;
        }
        a[href]:after {
            content: none !important;
        }
        .btn {
            display: none;
        }
    }
</style>
