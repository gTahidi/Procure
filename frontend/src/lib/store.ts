// src/lib/store.ts
import { writable } from 'svelte/store';
import type { User } from './userService';

// Store for authentication status
export const isAuthenticated = writable<boolean>(false);

// Store for user profile data
export const user = writable<User | null>(null);

// Store for loading state during auth operations
export const isLoading = writable<boolean>(true);

// Store for authentication errors
export const authError = writable<Error | null>(null);
