<script lang="ts">
  import { enhance } from '$app/forms';
  import type { PageData } from './$types';
  import { onMount } from 'svelte';
  import type { Asset } from '$lib/types';

  // Define the expected shape of the page data
  interface AssetsPageData {
    assets: Asset[];
    error: string | null;
  }

  export let data: PageData & AssetsPageData;

  // Initialize state with proper types
  let assets: Asset[] = data?.assets || [];
  let error: string | null = data?.error || null;
  let loading = !data;
  
  // Update state when data changes
  $: if (data) {
    assets = data.assets || [];
    error = data.error || null;
    loading = false;
  }
  
  // Form state
  let showAddForm = false;
  let isSubmitting = false;
  
  // Form data type
  interface AssetFormData {
    name: string;
    description: string;
    category: string;
    amr_id?: string;
    emr_id?: string;
    status?: string;
    location?: string;
  }
  
  // Initialize form data with default values
  let formData: AssetFormData = {
    name: '',
    description: '',
    category: '',
    amr_id: '',
    emr_id: '',
    status: 'Available',
    location: ''
  };
  
  // Form validation state
  let formErrors: Partial<Record<keyof AssetFormData, string>> = {};
  
  // Toggle the add asset form
  function toggleAddForm() {
    showAddForm = !showAddForm;
    if (!showAddForm) {
      // Reset form when closing
      formData = {
        name: '',
        description: '',
        category: '',
        amr_id: '',
        emr_id: '',
        status: 'Available',
        location: ''
      };
      formErrors = {};
    }
  }
  
  // Validate form fields
  function validateForm(): boolean {
    const errors: typeof formErrors = {};
    
    if (!formData.name?.trim()) {
      errors.name = 'Name is required';
    }
    
    if (!formData.category) {
      errors.category = 'Category is required';
    }
    
    formErrors = errors;
    return Object.keys(errors).length === 0;
  }
  
  // Handle form submission
  async function handleSubmit(event: Event) {
    event.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    isSubmitting = true;
    
    try {
      // TODO: Replace with actual API call when backend is ready
      console.log('Submitting asset:', formData);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Show success message
      alert('Asset created successfully!');
      
      // Reset form and close
      toggleAddForm();
      
      // Refresh assets list
      // await loadAssets();
      
    } catch (error) {
      console.error('Error creating asset:', error);
      alert('Failed to create asset. Please try again.');
    } finally {
      isSubmitting = false;
    }
  }
  
  // Handle input changes
  function handleInput(field: keyof AssetFormData, value: string) {
    formData = { ...formData, [field]: value };
    
    // Clear error when user types
    if (formErrors[field]) {
      formErrors = { ...formErrors, [field]: undefined };
    }
  }
  
  // Format date for display
  function formatDate(dateString: string | null | undefined): string {
    if (!dateString) return '-';
    try {
      return new Date(dateString).toLocaleDateString();
    } catch (e) {
      console.error('Error formatting date:', e);
      return '-';
    }
  }
</script>

<svelte:head>
  <title>Assets</title>
</svelte:head>

<div class="container mx-auto p-4 max-w-6xl">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-2xl font-bold">Assets</h1>
    <button 
      on:click={toggleAddForm} 
      class="btn btn-primary"
    >
      {showAddForm ? 'Cancel' : 'Add New Asset'}
    </button>
  </div>
  
  {#if showAddForm}
    <div class="bg-base-200 p-6 rounded-lg mb-6 shadow-md">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold">Add New Asset</h2>
        <button 
          on:click={toggleAddForm}
          class="btn btn-ghost btn-sm btn-circle"
          aria-label="Close form"
        >
          <span class="text-xl">×</span>
        </button>
      </div>
      
      <form on:submit={handleSubmit} class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Name Field -->
          <div class="form-control">
            <label for="asset-name" class="label">
              <span class="label-text">Name <span class="text-error">*</span></span>
            </label>
            <input 
              id="asset-name"
              type="text" 
              placeholder="Enter asset name" 
              class="input input-bordered w-full {formErrors.name ? 'input-error' : ''}"
              bind:value={formData.name}
              on:input={() => handleInput('name', formData.name)}
              aria-invalid={!!formErrors.name}
              aria-describedby={formErrors.name ? 'name-error' : undefined}
              required
            />
            {#if formErrors.name}
              <div class="label" id="name-error">
                <span class="label-text-alt text-error">{formErrors.name}</span>
              </div>
            {/if}
          </div>
          
          <!-- Category Field -->
          <div class="form-control">
            <label for="asset-category" class="label">
              <span class="label-text">Category <span class="text-error">*</span></span>
            </label>
            <select 
              id="asset-category"
              class="select select-bordered w-full {formErrors.category ? 'select-error' : ''}"
              bind:value={formData.category}
              on:change={() => handleInput('category', formData.category)}
              aria-invalid={!!formErrors.category}
              aria-describedby={formErrors.category ? 'category-error' : undefined}
              required
            >
              <option value="" disabled selected>Select a category</option>
              <option value="IT Equipment">IT Equipment</option>
              <option value="Furniture">Furniture</option>
              <option value="Vehicle">Vehicle</option>
              <option value="Other">Other</option>
            </select>
            {#if formErrors.category}
              <div class="label" id="category-error">
                <span class="label-text-alt text-error">{formErrors.category}</span>
              </div>
            {/if}
          </div>
          
          <!-- Status Field -->
          <div class="form-control">
            <label for="asset-status" class="label">
              <span class="label-text">Status</span>
            </label>
            <select 
              id="asset-status"
              class="select select-bordered w-full"
              bind:value={formData.status}
              on:change={() => handleInput('status', formData.status || '')}
            >
              <option value="Available">Available</option>
              <option value="In Use">In Use</option>
              <option value="Maintenance">Maintenance</option>
              <option value="Retired">Retired</option>
            </select>
          </div>
          
          <!-- Location Field -->
          <div class="form-control">
            <label for="asset-location" class="label">
              <span class="label-text">Location</span>
            </label>
            <input 
              id="asset-location"
              type="text" 
              placeholder="Enter location" 
              class="input input-bordered w-full"
              bind:value={formData.location}
              on:input={() => handleInput('location', formData.location || '')}
            />
          </div>
          
          <!-- AMR ID Field -->
          <div class="form-control">
            <label for="asset-amr-id" class="label">
              <span class="label-text">AMR ID</span>
            </label>
            <input 
              id="asset-amr-id"
              type="text" 
              placeholder="Enter AMR ID" 
              class="input input-bordered w-full"
              bind:value={formData.amr_id}
              on:input={() => handleInput('amr_id', formData.amr_id || '')}
            />
          </div>
          
          <!-- EMR ID Field -->
          <div class="form-control">
            <label for="asset-emr-id" class="label">
              <span class="label-text">EMR ID</span>
            </label>
            <input 
              id="asset-emr-id"
              type="text" 
              placeholder="Enter EMR ID" 
              class="input input-bordered w-full"
              bind:value={formData.emr_id}
              on:input={() => handleInput('emr_id', formData.emr_id || '')}
            />
          </div>
        </div>
        
        <!-- Description Field -->
        <div class="form-control">
          <label for="asset-description" class="label">
            <span class="label-text">Description</span>
          </label>
          <textarea 
            id="asset-description"
            class="textarea textarea-bordered h-24" 
            placeholder="Enter asset description"
            bind:value={formData.description}
            on:input={() => handleInput('description', formData.description)}
            aria-label="Asset description"
          ></textarea>
        </div>
        
        <div class="flex justify-end space-x-3 pt-2 border-t border-base-300 mt-6 pt-4">
          <!-- Cancel Button -->
          <button 
            type="button" 
            on:click={toggleAddForm}
            class="btn btn-ghost"
            disabled={isSubmitting}
            aria-busy={isSubmitting}
            aria-live="polite"
          >
            Cancel
          </button>
          
          <!-- Save Button -->
          <button 
            type="submit" 
            class="btn btn-primary min-w-32"
            class:btn-disabled={isSubmitting}
            disabled={isSubmitting}
            aria-busy={isSubmitting}
            aria-live="polite"
          >
            {#if isSubmitting}
              <span class="loading loading-spinner loading-sm mr-2"></span>
              <span>Creating...</span>
            {:else}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Save Asset</span>
            {/if}
          </button>
        </div>
      </form>
      
      <div class="mt-4 text-sm text-base-content/70">
        <p>Note: Asset management is currently in development. Some features may be limited.</p>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="flex justify-center py-8">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if error}
    <div class="alert alert-error mb-6">
      <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span>Error: {error}</span>
    </div>
  {:else if assets.length === 0}
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body items-center text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-content/30 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <h2 class="card-title">No Assets Yet</h2>
        <p class="mb-4">Get started by adding your first asset</p>
        <button on:click={toggleAddForm} class="btn btn-primary">
          Add Asset
        </button>
      </div>
    </div>
  {:else if assets}
    <div class="bg-base-100 rounded-lg shadow-sm border border-base-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table table-zebra w-full">
          <thead class="bg-base-200">
            <tr>
              <th class="w-16">ID</th>
              <th>Description</th>
              <th class="hidden md:table-cell">Category</th>
              <th class="hidden lg:table-cell">Status</th>
              <th class="hidden xl:table-cell">Location</th>
              <th class="w-24 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each assets as asset (asset.id)}
              <tr class="hover:bg-base-100 transition-colors">
                <td class="font-mono text-sm">#{asset.id}</td>
                <td>
                  <div class="font-medium">
                    <a 
                      href={`/assets/${asset.id}`} 
                      class="link link-hover link-primary"
                      aria-label={`View details for ${asset.description || 'asset'} ${asset.id}`}
                    >
                      {asset.description || 'Untitled Asset'}
                    </a>
                  </div>
                  <div class="text-xs text-gray-500 mt-1 md:hidden">
                    {asset.category} • {asset.status}
                  </div>
                </td>
                <td class="hidden md:table-cell">
                  <span class="badge badge-ghost">
                    {asset.category || 'Uncategorized'}
                  </span>
                </td>
                <td class="hidden lg:table-cell">
                  {#if asset.status === 'Available'}
                    <span class="badge badge-success">Available</span>
                  {:else if asset.status === 'In Use'}
                    <span class="badge badge-info">In Use</span>
                  {:else if asset.status === 'Maintenance'}
                    <span class="badge badge-warning">Maintenance</span>
                  {:else}
                    <span class="badge">{asset.status || 'N/A'}</span>
                  {/if}
                </td>
                <td class="hidden xl:table-cell">
                  {asset.location || 'Not specified'}
                </td>
                <td class="text-right">
                  <div class="dropdown dropdown-end">
                    <button 
                      type="button" 
                      class="btn btn-ghost btn-sm"
                      aria-label="Actions"
                      aria-haspopup="true"
                      aria-expanded="false"
                      on:click|stopPropagation
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z" />
                      </svg>
                    </button>
                    <ul 
                      class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-40 z-50"
                      role="menu"
                      aria-orientation="vertical"
                    >
                      <li>
                        <a href={`/assets/${asset.id}`} class="text-sm">
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                          </svg>
                          View Details
                        </a>
                      </li>
                      <li>
                        <a href={`/assets/${asset.id}/edit`} class="text-sm">
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                          </svg>
                          Edit
                        </a>
                      </li>
                      <li>
                        <button class="text-sm text-error" on:click|preventDefault={() => console.log('Delete', asset.id)}>
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                          </svg>
                          Delete
                        </button>
                      </li>
                    </ul>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      
      {#if assets.length > 10}
        <div class="flex justify-between items-center p-4 border-t border-base-200">
          <div class="text-sm text-gray-500">
            Showing <span class="font-medium">1-{assets.length}</span> of <span class="font-medium">{assets.length}</span> assets
          </div>
          <div class="join">
            <button class="join-item btn btn-sm btn-ghost" disabled>Previous</button>
            <button class="join-item btn btn-sm btn-active">1</button>
            <button class="join-item btn btn-sm btn-ghost">Next</button>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .container {
    max-width: 1200px;
  }
</style>
