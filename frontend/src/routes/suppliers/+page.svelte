<script lang="ts">
  import type { PageData } from './$types'; // Import PageData

  export let data: PageData; // Data from +page.ts

  // Reactive declarations to get suppliers and error from data
  $: suppliers = data.suppliers;
  $: error = data.error;
  $: loading = !suppliers && !error; // Simplified loading state
</script>

<svelte:head>
  <title>Suppliers</title>
</svelte:head>

<div class="container mx-auto p-4">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-2xl font-bold">Suppliers</h1>
    <a href="/suppliers/new" class="btn btn-primary">Add New Supplier</a>
  </div>

  {#if loading}
    <p>Loading suppliers...</p>
  {:else if error}
    <div class="alert alert-error">
      <p>Error loading suppliers: {error}</p>
    </div>
  {:else if suppliers && suppliers.length === 0}
    <p>No suppliers found. <a href="/suppliers/new" class="link">Add one?</a></p>
  {:else if suppliers}
    <div class="overflow-x-auto">
      <table class="table w-full">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Contact Person</th>
            <th>Email</th>
            <th>Phone</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each suppliers as supplier (supplier.id)}
            <tr>
              <td>{supplier.id}</td>
              <td><a href={`/suppliers/${supplier.id}`} class="link link-hover">{supplier.name}</a></td>
              <td>{supplier.contact_person || '-'}</td>
              <td>{supplier.email || '-'}</td>
              <td>{supplier.phone || '-'}</td>
              <td>
                <a href={`/suppliers/${supplier.id}`} class="btn btn-sm btn-outline">View</a>
                <!-- Edit/Delete buttons can be added here -->
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  /* Basic styling for DaisyUI or Tailwind if used */
  .container {
    max-width: 1200px;
  }
</style>
