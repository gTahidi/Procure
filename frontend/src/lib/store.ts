// src/lib/store.ts
import { writable } from 'svelte/store';

// Define our application's user interface
export interface AppUser {
  id?: number;
  username: string;
  email: string;
  role?: string;
  picture_url?: string;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
  // Add any other fields from your User model
}

// Store for authentication status
export const isAuthenticated = writable<boolean>(false);

// Store for user profile data
export const user = writable<AppUser | null>(null);

// Store for loading state during auth operations
export const isLoading = writable<boolean>(true);

// Store for authentication errors
export const authError = writable<Error | null>(null);
