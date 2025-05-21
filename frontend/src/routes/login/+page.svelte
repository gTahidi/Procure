<script lang="ts">
  import { login as auth0Login, checkAuthStatus } from '$lib/authService';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let loading = false;
  let error = '';

  onMount(() => {
    if (checkAuthStatus()) {
      const url = new URL(window.location.href);
      const redirectTo = url.searchParams.get('redirect_to') || '/';
      goto(redirectTo, { replaceState: true });
    }
  });

  async function handleLogin() {
    loading = true;
    error = '';
    try {
      // Call the Auth0 login function which handles redirection
      // No credentials are passed here for the redirect-based flow
      await auth0Login();
      // The browser will redirect to Auth0; code below this point might not execute
      // if the redirect happens immediately and successfully.
    } catch (e: any) {
      console.error('Auth0 login initiation failed:', e);
      error = e.message || 'Login failed. Please try again later.';
      loading = false;
    }
    // loading should ideally be handled by redirect or page change
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8 text-center">
    <div>
      <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
        Sign In
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        You will be redirected to our secure login provider.
      </p>
    </div>

    {#if error}
      <div class="bg-red-50 border-l-4 border-red-400 p-4 mt-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm text-red-700">{error}</p>
          </div>
        </div>
      </div>
    {/if}

    <div class="mt-8">
      <button
        type="button" 
        on:click={handleLogin}
        disabled={loading}
        class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {#if loading}
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Redirecting to Login...
        {:else}
          Login with Secure Provider
        {/if}
      </button>
    </div>
  </div>
</div>
