// src/lib/store.ts
import { writable } from 'svelte/store';
import type { Auth0Client } from '@auth0/auth0-spa-js';

// Store for the Auth0 client instance
export const auth0ClientStore = writable<Auth0Client | null>(null);

// Store for authentication status
export const isAuthenticated = writable<boolean>(false);

// Store for user profile data
export const user = writable<any | null>(null); // Replace 'any' with a proper User type if available or defined

// Store for loading state during auth operations
export const isLoading = writable<boolean>(true);

// Store for authentication errors
export const authError = writable<Error | null>(null);
