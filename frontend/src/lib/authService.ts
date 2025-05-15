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

  isLoading.set(true);
  try {
    const client = await createClient();
    auth0ClientStore.set(client);

    // Check for redirect callback
    if (window.location.search.includes('code=') && window.location.search.includes('state=')) {
      await handleRedirectCallback();
    } else {
      // Check session on initial load if not a redirect
      const currentlyAuthenticated = await client.isAuthenticated();
      if (currentlyAuthenticated) {
        const userProfile = await client.getUser();
        isAuthenticated.set(true);
        user.set(userProfile);
      } else {
        isAuthenticated.set(false);
        user.set(null);
      }
    }
  } catch (e: any) {
    console.error('AuthService Error:', e);
    authError.set(e);
    isAuthenticated.set(false);
    user.set(null);
  } finally {
    isLoading.set(false);
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
    authError.set(new Error('Auth0 client not initialized.'));
    return;
  }

  isLoading.set(true);
  try {
    await client.handleRedirectCallback();
    const userProfile = await client.getUser();
    isAuthenticated.set(true);
    user.set(userProfile);
    authError.set(null);
    // Remove query params from URL
    window.history.replaceState({}, document.title, window.location.pathname);
  } catch (e: any) {
    console.error('Redirect callback error', e);
    authError.set(e);
    isAuthenticated.set(false);
    user.set(null);
  } finally {
    isLoading.set(false);
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
