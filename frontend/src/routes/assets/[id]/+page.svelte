<script lang="ts">
  import type { PageData } from './$types';
  import { goto } from '$app/navigation';
  import type { Asset } from '$lib/types';

  export let data: PageData;

  let asset: Partial<Asset>;
  let errorMessage: string | null = null;
  let submitting = false;
  let deleting = false;

  $: {
    if (data.asset) {
      asset = { ...data.asset };
      // Convert purchase_date from potential ISO string to yyyy-MM-dd for date input
      if (asset.purchase_date && typeof asset.purchase_date === 'string') {
        asset.purchase_date = asset.purchase_date.split('T')[0];
      }
    } else {
      asset = {};
    }
    errorMessage = data.error;
  }

  async function handleUpdate() {
    if (!data.asset?.id) return;
    submitting = true;
    errorMessage = null;
    try {
      if (!asset.description) {
        throw new Error('Description is required.');
      }
      const response = await fetch(`/api/assets/${data.asset.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(asset)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to update asset.' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }
      goto('/assets');
    } catch (err: any) {
      errorMessage = err.message;
      console.error('Update error:', err);
    }
    submitting = false;
  }

  async function handleDelete() {
    if (!data.asset?.id || !confirm('Are you sure you want to delete this asset?')) {
      return;
    }
    deleting = true;
    errorMessage = null;
    try {
      const response = await fetch(`/api/assets/${data.asset.id}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to delete asset.' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }
      goto('/assets');
    } catch (err: any) {
      errorMessage = err.message;
      console.error('Delete error:', err);
    }
    deleting = false;
  }
</script>

<svelte:head>
  <title>View/Edit Asset: {data.asset?.description || 'Asset'}</title>
</svelte:head>

<div class="container mx-auto p-4">
  {#if data.error && !data.asset}
    <div class="alert alert-error">
      <p>Error loading asset: {data.error}</p>
    </div>
    <a href="/assets" class="btn btn-link mt-4">Back to Assets List</a>
  {:else if !data.asset}
    <p>Asset not found.</p> 
    <a href="/assets" class="btn btn-link mt-4">Back to Assets List</a>
  {:else}
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">View/Edit Asset: {data.asset.description}</h1>
      <button class="btn btn-error btn-outline" on:click={handleDelete} disabled={deleting}>
        {#if deleting}Deleting...{:else}Delete Asset{/if}
      </button>
    </div>

    {#if errorMessage}
      <div class="alert alert-error mb-4">
        <p>{errorMessage}</p>
      </div>
    {/if}

    <form on:submit|preventDefault={handleUpdate} class="space-y-4">
      <div>
        <label for="description" class="label">Description <span class="text-error">*</span></label>
        <textarea id="description" bind:value={asset.description} class="textarea textarea-bordered w-full" required rows="3"></textarea>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
            <label for="amr_id" class="label">AMR ID</label>
            <input type="number" id="amr_id" bind:value={asset.amr_id} class="input input-bordered w-full" />
        </div>
        <div>
            <label for="emr_id" class="label">EMR ID</label>
            <input type="number" id="emr_id" bind:value={asset.emr_id} class="input input-bordered w-full" />
        </div>
      </div>
      <div>
        <label for="category" class="label">Category</label>
        <input type="text" id="category" bind:value={asset.category} class="input input-bordered w-full" />
      </div>
      <div>
        <label for="status" class="label">Status</label>
        <select id="status" bind:value={asset.status} class="select select-bordered w-full">
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
          <option value="disposed">Disposed</option>
          <option value="maintenance">Under Maintenance</option>
        </select>
      </div>
      <div>
        <label for="location" class="label">Location</label>
        <input type="text" id="location" bind:value={asset.location} class="input input-bordered w-full" />
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
              <label for="purchase_date" class="label">Purchase Date</label>
              <input type="date" id="purchase_date" bind:value={asset.purchase_date} class="input input-bordered w-full" />
          </div>
          <div>
              <label for="purchase_price" class="label">Purchase Price</label>
              <input type="number" step="0.01" id="purchase_price" bind:value={asset.purchase_price} class="input input-bordered w-full" />
          </div>
      </div>
      <div>
        <label for="supplier_id" class="label">Supplier ID (Optional)</label>
        <input type="number" id="supplier_id" bind:value={asset.supplier_id} class="input input-bordered w-full" />
      </div>

      <div class="flex gap-2 pt-4">
        <button type="submit" class="btn btn-primary" disabled={submitting}>
          {#if submitting}Saving...{:else}Save Changes{/if}
        </button>
        <a href="/assets" class="btn btn-ghost">Cancel</a>
      </div>
    </form>
  {/if}
</div>

<style>
  .container {
    max-width: 700px;
  }
</style>
