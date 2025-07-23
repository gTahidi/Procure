<script lang="ts">
  import { checkAuthStatus } from '$lib/authService';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import ChangePasswordForm from '$lib/components/auth/ChangePasswordForm.svelte';

  onMount(() => {
    // Redirect to login if not authenticated
    if (!checkAuthStatus()) {
      goto('/login?redirect_to=/change-password', { replaceState: true });
    }
  });
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div class="text-center">
      <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
        Change Password
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        Update your password to keep your account secure
      </p>
    </div>

    <div class="mt-8">
      <ChangePasswordForm />
    </div>
  </div>
</div>