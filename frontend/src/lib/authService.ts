// src/lib/authService.ts
import { browser } from '$app/environment';
import { request } from './services/apiService'; // Assuming apiService handles sending cookies
import type { User, LoginCredentials } from './userService'; // User and LoginCredentials from userService
import { getCurrentUser as fetchCurrentUser } from './userService'; // Alias to avoid conflict if any
import { isAuthenticated, user, isLoading, authError } from './store';

/**
 * Initializes the authentication state when the app loads.
 * Checks if a valid session exists by calling /auth/me.
 */
export async function initializeAuth(): Promise<void> {
  if (!browser) return;

  isLoading.set(true);
  authError.set(null);
  try {
    const currentUser = await fetchCurrentUser(); // Relies on http-only cookie being sent by browser
    if (currentUser) {
      user.set(currentUser);
      isAuthenticated.set(true);
    } else {
      user.set(null);
      isAuthenticated.set(false);
    }
  } catch (e: unknown) {
    console.error('AuthService InitializeAuth Error:', e);
    user.set(null);
    isAuthenticated.set(false);
    if (e instanceof Error) {
      authError.set(e);
    } else {
      authError.set(new Error('An unknown error occurred during auth initialization.'));
    }
  } finally {
    isLoading.set(false);
  }
}

/**
 * Logs in the user with the given credentials.
 * Backend sets an http-only session cookie on success and returns user data.
 * @param credentials Email and password.
 * @returns The user object if login is successful.
 */
export async function login(credentials: LoginCredentials): Promise<User> {
  if (!browser) {
    throw new Error('Login can only be performed in the browser');
  }

  isLoading.set(true);
  authError.set(null);
  try {
    // Backend sends user data directly in response body for /auth/login
    const loggedInUser = await request('/auth/login', 'POST', credentials) as User;
    
    if (!loggedInUser || !loggedInUser.id) {
      // This case should ideally be handled by request throwing an error for non-2xx status
      throw new Error('Login failed: Invalid response from server.');
    }
    
    user.set(loggedInUser);
    isAuthenticated.set(true);
    return loggedInUser;
  } catch (e: unknown) {
    console.error('AuthService Login Error:', e);
    user.set(null);
    isAuthenticated.set(false);
    if (e instanceof Error) {
        authError.set(e);
        throw e; // Re-throw original error
    } else {
        const err = new Error('An unknown error occurred during login.');
        authError.set(err);
        throw err;
    }
  } finally {
    isLoading.set(false);
  }
}

/**
 * Logs out the current user.
 * Calls the backend /auth/logout endpoint to invalidate the session.
 */
export async function logout(): Promise<void> {
  if (!browser) return;

  isLoading.set(true);
  authError.set(null);
  try {
    await request('/auth/logout', 'POST'); // Backend handles cookie invalidation
  } catch (e: unknown) {
    console.error('AuthService Logout Error:', e);
    // Even if logout API call fails, clear frontend state as a fallback.
    if (e instanceof Error) {
      authError.set(e);
    } else {
      authError.set(new Error('An unknown error occurred during logout.'));
    }
    // Potentially re-throw or handle more gracefully depending on requirements
  } finally {
    user.set(null);
    isAuthenticated.set(false);
    isLoading.set(false);
    // No need to remove 'authToken' from localStorage as we're not using it.
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
