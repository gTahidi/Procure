<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { user, isAuthenticated } from '$lib/store';
  import { logout, checkAuthStatus } from '$lib/authService';
  
  let loading = true;
  
  onMount(() => {
    if (!checkAuthStatus()) {
      // If not authenticated, redirect to login with a redirect back here
      const currentPath = window.location.pathname;
      goto(`/login?redirect_to=${encodeURIComponent(currentPath)}`);
    } else {
      loading = false;
    }
  });
</script>

<svelte:head>
  <title>Dashboard - Procurement System</title>
  <meta name="description" content="Main dashboard for the Procurement System" />
</svelte:head>

{#if loading}
  <div class="flex justify-center items-center h-64">
    <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
  </div>
{:else}
  <div class="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <div class="text-center mb-12">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">
        Welcome back, {$user?.name || 'User'}!
      </h1>
      <p class="text-lg text-gray-600">
        Here's what's happening with your procurement activities today.
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
      <!-- Quick Actions -->
      <div class="bg-white overflow-hidden shadow rounded-lg">
        <div class="px-4 py-5 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900">Quick Actions</h3>
          <div class="mt-6 space-y-4">
            <a href="/requisitions/new" class="block w-full text-left px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
              Create New Requisition
            </a>
            <a href="/tenders" class="block w-full text-left px-4 py-2 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
              View Open Tenders
            </a>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-white overflow-hidden shadow rounded-lg md:col-span-2">
        <div class="px-4 py-5 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900">Recent Activity</h3>
          <div class="mt-6 space-y-4">
            <div class="border-b border-gray-200 pb-4">
              <p class="text-sm text-gray-600">No recent activity to display.</p>
              <p class="text-xs text-gray-500 mt-1">Your recent actions will appear here.</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Getting Started -->
    <div class="bg-white overflow-hidden shadow rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium text-gray-900">Getting Started</h3>
        <div class="mt-6 space-y-4">
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <span class="h-6 w-6 rounded-full bg-indigo-100 flex items-center justify-center">
                <span class="text-indigo-600 font-medium text-sm">1</span>
              </span>
            </div>
            <div class="ml-3">
              <p class="text-sm font-medium text-gray-900">Complete your profile</p>
              <p class="text-sm text-gray-500">Add your contact information and preferences.</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <span class="h-6 w-6 rounded-full bg-indigo-100 flex items-center justify-center">
                <span class="text-indigo-600 font-medium text-sm">2</span>
              </span>
            </div>
            <div class="ml-3">
              <p class="text-sm font-medium text-gray-900">Create your first requisition</p>
              <p class="text-sm text-gray-500">Start the procurement process for your needs.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
