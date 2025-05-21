// src/lib/authService.ts
import { createAuth0Client, type Auth0Client, type LogoutOptions, type RedirectLoginOptions, type GetTokenSilentlyOptions, type User } from '@auth0/auth0-spa-js';
import { auth0ClientStore, isAuthenticated, user, isLoading, authError } from './store';
import { PUBLIC_VITE_AUTH0_DOMAIN, PUBLIC_VITE_AUTH0_CLIENT_ID, PUBLIC_VITE_AUTH0_CALLBACK_URL, PUBLIC_VITE_AUTH0_AUDIENCE } from '$env/static/public';
import { get } from 'svelte/store';
import { browser } from '$app/environment';

async function createClient(): Promise<Auth0Client> {
  const client = await createAuth0Client({
    domain: PUBLIC_VITE_AUTH0_DOMAIN,
    clientId: PUBLIC_VITE_AUTH0_CLIENT_ID,
    authorizationParams: {
      redirect_uri: PUBLIC_VITE_AUTH0_CALLBACK_URL,
      audience: PUBLIC_VITE_AUTH0_AUDIENCE, // For requesting token for your Go API
    },
  });
  return client;
}

export async function initializeAuth(): Promise<void> {
  if (!browser) return; // Ensure this runs only in the browser

  isLoading.set(true); // Set loading at the very start

  let client = get(auth0ClientStore);

  if (!client) {
    console.log('[authService initializeAuth] No existing client, creating new one...');
    client = await createClient();
    auth0ClientStore.set(client);
    console.log('[authService initializeAuth] New Auth0 client created and stored.');
  } else {
    console.log('[authService initializeAuth] Using existing Auth0 client.');
  }

  // Check if this is a redirect from Auth0
  if (window.location.search.includes('code=') && window.location.search.includes('state=')) {
    console.log('[authService initializeAuth] Detected Auth0 redirect (code and state in URL). Calling handleRedirectCallback...');
    // handleRedirectCallback will manage its own isLoading state and update user/isAuthenticated stores.
    await handleRedirectCallback(); 
    // IMPORTANT: Exit initializeAuth after handling the redirect. 
    // isLoading.set(false) is handled by handleRedirectCallback's finally block.
    return; 
  }

  // If not a redirect, check current authentication state
  console.log('[authService initializeAuth] Not an Auth0 redirect. Checking current session status...');
  try {
    const currentlyAuthenticated = await client.isAuthenticated();
    console.log('[authService initializeAuth] client.isAuthenticated() check result:', currentlyAuthenticated);
    if (currentlyAuthenticated) {
      const userProfile = await client.getUser();
      console.log('[authService initializeAuth] User is authenticated (no redirect). Profile:', userProfile);
      user.set(userProfile);
      isAuthenticated.set(true);
      authError.set(null);
    } else {
      console.log('[authService initializeAuth] User is NOT authenticated (no redirect).');
      isAuthenticated.set(false);
      user.set(null);
      authError.set(null); // Not an error, just not logged in.
    }
  } catch (e: any) {
    console.error('[authService initializeAuth] Error checking/refreshing session (non-redirect):', e);
    authError.set(e);
    isAuthenticated.set(false);
    user.set(null);
  } finally {
    isLoading.set(false); // Ensure loading is false if it's the non-redirect path
    console.log('[authService initializeAuth] initializeAuth finished (non-redirect path). isLoading:', get(isLoading));
  }
}

export async function login(options?: RedirectLoginOptions): Promise<void> {
  const client = get(auth0ClientStore);
  if (!client) {
    console.error('Auth0 client not initialized for login');
    authError.set(new Error('Auth0 client not initialized.'));
    return;
  }
  await client.loginWithRedirect(options);
}

async function handleRedirectCallback(): Promise<void> {
  const client = get(auth0ClientStore);
  if (!client) {
    console.error('Auth0 client not initialized for handleRedirectCallback');
    // Ensure authError is set to trigger UI update
    authError.set(new Error('Auth0 client not initialized.')); 
    isAuthenticated.set(false); // Explicitly set isAuthenticated to false
    user.set(null); // Explicitly set user to null
    return;
  }

  isLoading.set(true);
  try {
    console.log('[authService] Starting client.handleRedirectCallback()...');
    await client.handleRedirectCallback(); // Exchanges code for tokens
    console.log('[authService] client.handleRedirectCallback() completed.');

    // Get User Profile
    const userProfile = await client.getUser();
    console.log('[authService] client.getUser() profile:', userProfile);

    // Get ID Token Claims
    const idTokenClaims = await client.getIdTokenClaims();
    console.log('[authService] client.getIdTokenClaims() result:', idTokenClaims);

    // Attempt to get Access Token
    let accessToken;
    try {
      accessToken = await client.getTokenSilently();
      console.log('[authService] client.getTokenSilently() accessToken retrieved (details omitted for brevity).');
      // For debugging, you could log the token itself but be careful with sharing full tokens.
      // console.log('[authService] Access Token:', accessToken); 
    } catch (tokenError) {
      console.error('[authService] Error calling client.getTokenSilently():', tokenError);
    }

    if (userProfile) { // Or you could check idTokenClaims if that's more reliable for your setup
      isAuthenticated.set(true);
      user.set(userProfile); // Storing the profile from getUser()
      authError.set(null);
      console.log('[authService] User profile loaded and stores updated.');
    } else {
      console.error('[authService] Token exchange likely successful, but no user profile (from getUser) or no ID token claims.');
      isAuthenticated.set(false);
      user.set(null);
      authError.set(new Error('Login successful, but failed to retrieve user details. Check console for token claims.'));
    }
    
    // Remove query params from URL - should be safe to do even if userProfile is null
    // as long as handleRedirectCallback itself didn't throw before this.
    window.history.replaceState({}, document.title, window.location.pathname);
    console.log('[authService] Query params removed from URL.');

  } catch (e: any) {
    console.error('[authService] Error during redirect callback processing (outer try-catch):', e);
    authError.set(e); // This is likely where "Login required" or similar Auth0 errors are caught
    isAuthenticated.set(false);
    user.set(null);
  } finally {
    isLoading.set(false);
    console.log('[authService] handleRedirectCallback finally block. isLoading:', get(isLoading));
  }
}

export async function logout(options?: LogoutOptions): Promise<void> {
  const client = get(auth0ClientStore);
  if (!client) {
    console.error('Auth0 client not initialized for logout');
    authError.set(new Error('Auth0 client not initialized.'));
    return;
  }
  
  // Ensure callback URL is defined for logout, falling back to origin
  const logoutOptions: LogoutOptions = { ...options }; // Use LogoutOptions type
  if (!logoutOptions.logoutParams?.returnTo) {
    logoutOptions.logoutParams = { ...logoutOptions.logoutParams, returnTo: window.location.origin };
  }

  isAuthenticated.set(false);
  user.set(null);
  await client.logout(logoutOptions);
}

// Function to get an access token
export async function getAccessTokenSilently(options?: Omit<GetTokenSilentlyOptions, 'detailedResponse'>): Promise<string | undefined> {
  const client = get(auth0ClientStore);
  if (!client) {
    console.error('Auth0 client not initialized for getAccessTokenSilently');
    authError.set(new Error('Auth0 client not initialized.'));
    return undefined;
  }
  try {
    // Ensure options don't ask for detailedResponse to match return type string
    const token = await client.getTokenSilently(options);
    return token;
  } catch (e: any) {
    console.error('Error getting token silently:', e);
    if (e.error === 'login_required' || e.error === 'consent_required') {
      // Consider triggering interactive login if appropriate
      // await login({ appState: { targetUrl: window.location.pathname } });
    }
    authError.set(e);
    return undefined;
  }
}

/**
 * Checks the current authentication status from the store.
 * Does not make an API call; reflects the last known state.
 * @returns boolean indicating if user is currently marked as authenticated.
 */
export function checkAuthStatus(): boolean {
  if (!browser) return false;
  let currentAuthStatus = false;
  isAuthenticated.subscribe(value => currentAuthStatus = value)(); // Get current value
  return currentAuthStatus;
}
