<script lang="ts">
  import type { PageData } from './$types';
  import { goto } from '$app/navigation';
  import type { Supplier } from '$lib/types';

  export let data: PageData;

  let supplier: Partial<Supplier>; // Use Partial for the form binding
  let errorMessage: string | null = null;
  let submitting = false;
  let deleting = false;

  // Initialize form data when page data is available
  $: {
    if (data.supplier) {
      supplier = { ...data.supplier }; // Clone to avoid mutating original data directly
    } else {
      supplier = {}; // Handle case where supplier might be null (e.g. on error)
    }
    errorMessage = data.error;
  }

  async function handleUpdate() {
    if (!data.supplier?.id) return;
    submitting = true;
    errorMessage = null;
    try {
      const response = await fetch(`/api/suppliers/${data.supplier.id}`, { // PUT to your backend API
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(supplier)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to update supplier.' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }
      // Optionally, show a success message or refresh data
      // For now, let's assume the +page.ts will re-fetch or SvelteKit re-runs load on navigation
      goto('/suppliers'); // Or stay on page and show success
    } catch (err: any) {
      errorMessage = err.message;
      console.error('Update error:', err);
    }
    submitting = false;
  }

  async function handleDelete() {
    if (!data.supplier?.id || !confirm('Are you sure you want to delete this supplier?')) {
      return;
    }
    deleting = true;
    errorMessage = null;
    try {
      const response = await fetch(`/api/suppliers/${data.supplier.id}`, { // DELETE to your backend API
        method: 'DELETE'
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to delete supplier.' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }
      goto('/suppliers');
    } catch (err: any) {
      errorMessage = err.message;
      console.error('Delete error:', err);
    }
    deleting = false;
  }
</script>

<svelte:head>
  <title>View/Edit Supplier: {data.supplier?.name || 'Supplier'}</title>
</svelte:head>

<div class="container mx-auto p-4">
  {#if data.error && !data.supplier}
    <div class="alert alert-error">
      <p>Error loading supplier: {data.error}</p>
    </div>
    <a href="/suppliers" class="btn btn-link mt-4">Back to Suppliers List</a>
  {:else if !data.supplier}
    <p>Supplier not found.</p> 
    <a href="/suppliers" class="btn btn-link mt-4">Back to Suppliers List</a>
  {:else}
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">View/Edit Supplier: {data.supplier.name}</h1>
      <button class="btn btn-error btn-outline" on:click={handleDelete} disabled={deleting}>
        {#if deleting}Deleting...{:else}Delete Supplier{/if}
      </button>
    </div>

    {#if errorMessage}
      <div class="alert alert-error mb-4">
        <p>{errorMessage}</p>
      </div>
    {/if}

    <form on:submit|preventDefault={handleUpdate} class="space-y-4">
      <div>
        <label for="name" class="label">Name <span class="text-error">*</span></label>
        <input type="text" id="name" bind:value={supplier.name} class="input input-bordered w-full" required />
      </div>
      <div>
        <label for="contact_person" class="label">Contact Person</label>
        <input type="text" id="contact_person" bind:value={supplier.contact_person} class="input input-bordered w-full" />
      </div>
      <div>
        <label for="email" class="label">Email</label>
        <input type="email" id="email" bind:value={supplier.email} class="input input-bordered w-full" />
      </div>
      <div>
        <label for="phone" class="label">Phone</label>
        <input type="tel" id="phone" bind:value={supplier.phone} class="input input-bordered w-full" />
      </div>

      <div class="flex gap-2 mt-6">
        <button type="submit" class="btn btn-primary" disabled={submitting}>
          {#if submitting}Saving...{:else}Save Changes{/if}
        </button>
        <a href="/suppliers" class="btn btn-ghost">Cancel</a>
      </div>
    </form>
  {/if}
</div>

<style>
  .container {
    max-width: 700px;
  }
</style>
