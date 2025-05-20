<script lang="ts">
  import type { Supplier } from '$lib/types';
  import { goto } from '$app/navigation';

  let supplier: Partial<Supplier> = {
    name: '',
    contact_person: '',
    email: '',
    phone: ''
  };
  let errorMessage: string | null = null;
  let submitting = false;

  async function handleSubmit() {
    submitting = true;
    errorMessage = null;
    try {
      const response = await fetch('/api/suppliers', { // POST to your backend API endpoint
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(supplier)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to submit form. Please try again.' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      // const newSupplier = await response.json();
      // Optionally, navigate to the new supplier's page or back to the list
      goto('/suppliers'); 
    } catch (err: any) {
      errorMessage = err.message;
      console.error('Submission error:', err);
    }
    submitting = false;
  }
</script>

<svelte:head>
  <title>Add New Supplier</title>
</svelte:head>

<div class="container mx-auto p-4">
  <h1 class="text-2xl font-bold mb-6">Add New Supplier</h1>

  {#if errorMessage}
    <div class="alert alert-error mb-4">
      <p>{errorMessage}</p>
    </div>
  {/if}

  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
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
      <input type="tel" id="phone" bind:value={supplier.phone} class="input input-input-bordered w-full" />
    </div>

    <div class="flex gap-2">
      <button type="submit" class="btn btn-primary" disabled={submitting}>
        {#if submitting}Saving...{:else}Save Supplier{/if}
      </button>
      <a href="/suppliers" class="btn btn-ghost">Cancel</a>
    </div>
  </form>
</div>

<style>
  .container {
    max-width: 700px;
  }
</style>
