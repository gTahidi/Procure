<svelte:head>
	<title>Tenders - Procurement System</title>
	<meta name="description" content="View and manage tenders" />
</svelte:head>

<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;

  $: tenders = data.tenders;
  $: error = data.error;
  $: loading = !tenders && !error; // Simplified loading state
</script>

<div class="container mx-auto py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-semibold">Tenders</h1>
		<a href="/tenders/new" class="bg-teal-600 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
			+ New Tender
		</a>
	</div>

  {#if loading}
    <p class="text-center py-10">Loading tenders...</p>
  {:else if error}
    <div class="alert alert-error">
      <p>Error loading tenders: {error}</p>
    </div>
  {:else if tenders && tenders.length === 0}
    <p class="text-center py-10">No tenders found. <a href="/tenders/new" class="link">Add one?</a></p>
  {:else if tenders}
    <div class="bg-white shadow-md rounded-lg overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Closing Date</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#each tenders as tender (tender.id)}
            <tr>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-teal-600 hover:text-teal-900">
                <a href={`/tenders/${tender.id}`}>{tender.id}</a>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{tender.title}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{tender.category || '-'}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span class:px-2={true} class:inline-flex={true} class:text-xs={true} class:leading-5={true} class:font-semibold={true} class:rounded-full={true} 
                  class:bg-blue-100={tender.status === 'published' || tender.status === 'open'}
                  class:text-blue-800={tender.status === 'published' || tender.status === 'open'}
                  class:bg-gray-100={tender.status === 'draft' || tender.status === 'closed'}
                  class:text-gray-800={tender.status === 'draft' || tender.status === 'closed'}
                  class:bg-green-100={tender.status === 'awarded'}
                  class:text-green-800={tender.status === 'awarded'}
                  class:bg-yellow-100={tender.status === 'evaluation'}
                  class:text-yellow-800={tender.status === 'evaluation'}
                  class:bg-red-100={tender.status === 'cancelled'}
                  class:text-red-800={tender.status === 'cancelled'}
                  class:bg-purple-100={!(tender.status === 'published' || tender.status === 'open' || tender.status === 'draft' || tender.status === 'closed' || tender.status === 'awarded' || tender.status === 'evaluation' || tender.status === 'cancelled')}
                  class:text-purple-800={!(tender.status === 'published' || tender.status === 'open' || tender.status === 'draft' || tender.status === 'closed' || tender.status === 'awarded' || tender.status === 'evaluation' || tender.status === 'cancelled')}
                >
                  {tender.status || 'N/A'}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {#if tender.closing_date}
                  {new Date(tender.closing_date).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                  })}
                {:else}
                  -
                {/if}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <a href={`/tenders/${tender.id}`} class="text-teal-600 hover:text-teal-900">View</a>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
