// frontend/src/lib/userService.ts
import { getAccessToken } from './authService'; // Updated to use our new function
import { user } from './store'; // Import the user store directly
import type { AppUser } from './store'; // Import the AppUser type
import { get } from 'svelte/store';

/**
 * Interface for the payload sent to the backend when syncing user information.
 * This should match the expected structure of your Go backend's /api/users/sync endpoint.
 */
export interface UserSyncPayload {
  username: string;
  email: string;
  first_name?: string;
  last_name?: string;
  // Add any other fields that your backend needs to store
}

/**
 * Synchronizes the authenticated user's data with the application's backend database.
 * This function should be called after a user successfully logs in
 * and their profile is available in the user store.
 */
export async function syncUserToDb(): Promise<void> {
  const currentUser = get(user); // Get the current user profile from the Svelte store

  // Ensure there is a user to sync
  if (!currentUser || !currentUser.email) {
    console.warn('syncUserToDb: No user data found in store. Skipping sync.');
    return; // Exit if no user is found
  }

  try {
    // Obtain an access token. This token will be used to authenticate the request to your backend.
    const accessToken = getAccessToken();
    if (!accessToken) {
      console.error('syncUserToDb: Failed to obtain access token. Cannot sync user.');
      return; // Exit if no access token could be retrieved
    }

    // Prepare the payload with the data to be sent to the backend.
    const payload: UserSyncPayload = {
      username: currentUser.username,
      email: currentUser.email,
      // Add any other fields that your backend needs
    };

    // Make the API call to your backend's sync endpoint.
    const response = await fetch('http://localhost:8080/api/users/sync', { // Ensure this is your correct backend endpoint URI
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`, // Send the access token for backend verification
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      // Handle non-successful HTTP responses (e.g., 4xx, 5xx errors)
      const errorData = await response.text(); // Attempt to get error details from response body
      console.error('syncUserToDb: Failed to sync user with backend.', {
        status: response.status,
        statusText: response.statusText,
        errorData,
      });
      // Optionally, throw an error to be caught by the caller, allowing for more specific error handling UI-side.
      // throw new Error(`Backend sync failed: ${response.status} - ${errorData || response.statusText}`);
      return;
    }

    // Assuming the backend returns JSON, parse it.
    // This could be the created/updated user record from your database or a success message.
    const responseData: AppUser = await response.json(); // Explicitly type responseData
    console.log('syncUserToDb: User synchronized successfully with backend.', responseData);

    // Update the user store with the full user details from the backend (including role and id)
    user.set(responseData);

  } catch (error) {
    // Handle network errors or other issues during the fetch operation or token retrieval.
    console.error('syncUserToDb: An unexpected error occurred during user synchronization:', error);
    // Optionally, re-throw the error or handle it by updating UI state (e.g., show a notification to the user).
    // throw error;
  }
}

// If you have other user-specific (but not authentication-related) functions,
// such as fetching user preferences from your backend (once they are synced and have an internal ID),
// those could also reside in this service file.
