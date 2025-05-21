// src/lib/store.ts
import { writable } from 'svelte/store';
import type { Auth0Client } from '@auth0/auth0-spa-js'; // For Auth0Client type

// Store for the Auth0 client instance
export const auth0ClientStore = writable<Auth0Client | null>(null);

// Store for authentication status
export const isAuthenticated = writable<boolean>(false);

// Store for user profile data from Auth0
// Using 'any' for now, as Auth0's user profile is a flexible object.
// You might want to define a more specific interface for it based on your needs.
export const user = writable<any | null>(null);

// Store for loading state during auth operations
export const isLoading = writable<boolean>(true);

// Store for authentication errors
export const authError = writable<Error | null>(null);
