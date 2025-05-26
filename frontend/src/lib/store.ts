// src/lib/store.ts
import { writable } from 'svelte/store';
import type { User as Auth0UserProfile } from '@auth0/auth0-spa-js'; 
import type { Auth0Client } from '@auth0/auth0-spa-js'; 

// Define a new interface for our application's user, extending Auth0's user profile
export interface AppUser extends Auth0UserProfile {
  role?: string; // To store the role from our backend (e.g., 'requester', 'admin')
  id?: number; // Assuming your backend User model has an int64 ID
  picture_url?: string; // Explicitly add picture_url to match backend response
}

// Store for the Auth0 client instance
export const auth0ClientStore = writable<Auth0Client | null>(null);

// Store for authentication status
export const isAuthenticated = writable<boolean>(false);

// Store for user profile data from Auth0
export const user = writable<AppUser | null>(null);

// Store for loading state during auth operations
export const isLoading = writable<boolean>(true);

// Store for authentication errors
export const authError = writable<Error | null>(null);
