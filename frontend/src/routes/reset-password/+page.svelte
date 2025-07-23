<script lang="ts">
  import { checkAuthStatus } from '$lib/authService';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import ResetPasswordForm from '$lib/components/auth/ResetPasswordForm.svelte';

  onMount(() => {
    if (checkAuthStatus()) {
      const url = new URL(window.location.href);
      const redirectTo = url.searchParams.get('redirect_to') || '/';
      goto(redirectTo, { replaceState: true });
    }
  });
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div class="text-center">
      <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
        Set New Password
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        Enter your new password below
      </p>
    </div>

    <div class="mt-8">
      <ResetPasswordForm />
    </div>
  </div>
</div>