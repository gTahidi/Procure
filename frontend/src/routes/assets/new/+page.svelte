<script lang="ts">
  import type { Asset } from '$lib/types';
  import { goto } from '$app/navigation';

  let asset: Partial<Asset> = {
    description: '',
    amr_id: undefined,
    emr_id: undefined,
    category: '',
    status: 'active', // Default status
    location: '',
    purchase_date: undefined,
    purchase_price: undefined,
    supplier_id: undefined
  };
  let errorMessage: string | null = null;
  let submitting = false;

  async function handleSubmit() {
    submitting = true;
    errorMessage = null;
    try {
      // Basic client-side validation (more can be added)
      if (!asset.description) {
        throw new Error('Description is required.');
      }

      const response = await fetch('/api/assets', { // POST to your asset API endpoint
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(asset)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to submit form. Please try again.' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      goto('/assets'); 
    } catch (err: any) {
      errorMessage = err.message;
      console.error('Submission error:', err);
    }
    submitting = false;
  }
</script>

<svelte:head>
  <title>Add New Asset</title>
</svelte:head>

<div class="container mx-auto p-4">
  <h1 class="text-2xl font-bold mb-6">Add New Asset</h1>

  {#if errorMessage}
    <div class="alert alert-error mb-4">
      <p>{errorMessage}</p>
    </div>
  {/if}

  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
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
        {#if submitting}Saving...{:else}Save Asset{/if}
      </button>
      <a href="/assets" class="btn btn-ghost">Cancel</a>
    </div>
  </form>
</div>

<style>
  .container {
    max-width: 700px;
  }
</style>
