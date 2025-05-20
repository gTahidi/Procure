<script lang="ts">
    import { onMount } from 'svelte';
    import { initializeAuth } from '$lib/authService'; // The key function that handles callback
    import { syncUserToDb } from '$lib/userService';
    import { user, isAuthenticated, isLoading, authError } from '$lib/store'; // Import all relevant stores
    import { goto } from '$app/navigation';
    import { page } from '$app/stores'; // To get original target_url if any
    import { get } from 'svelte/store'; // To get current store values synchronously

    let localError: any = null;

    onMount(async () => {
        // initializeAuth() will handle the Auth0 redirect, update stores, and clear URL params.
        await initializeAuth();

        const unsubscribe = user.subscribe(async (currentUser) => {
            // Check if auth processing is done (isLoading is false) before acting
            if (!get(isLoading)) {
                if (currentUser && get(isAuthenticated)) {
                    unsubscribe(); // Unsubscribe to prevent multiple executions

                    console.log("Auth0 user processed by authService, now syncing to DB:", currentUser);
                    try {
                        const syncResult = await syncUserToDb(currentUser);
                        console.log("User sync result from /callback:", syncResult);
                    } catch (syncErr) {
                        console.error("Error syncing user from /callback:", syncErr);
                        localError = syncErr;
                    }

                    const targetUrl = $page.url.searchParams.get('target_url') || '/';
                    goto(targetUrl, { replaceState: true });

                } else if (get(authError)) {
                    unsubscribe();
                    localError = get(authError);
                    console.error("Auth error during callback processing:", localError);
                    goto('/', { replaceState: true }); 
                } else if (!currentUser) {
                    unsubscribe();
                    console.log("Callback page: No user session after initializeAuth. Redirecting home.");
                    goto('/', { replaceState: true });
                }
            }
        });

        // Safety timeout in case something unexpected happens and isLoading stays true
        setTimeout(() => {
            if (get(isLoading)) {
                console.warn("Callback page still loading after timeout, redirecting home.");
                if (typeof unsubscribe === 'function') unsubscribe();
                goto('/', { replaceState: true });
            }
        }, 7000); // 7 seconds timeout
    });
</script>

{#if $isLoading}
    <p>Completing login...</p>
{:else if localError || $authError}
    <p>Error during login: { (localError || $authError)?.message || (localError || $authError) }</p>
    <p><a href="/">Go to Homepage</a></p>
{:else}
    <p>Redirecting...</p>
{/if}
