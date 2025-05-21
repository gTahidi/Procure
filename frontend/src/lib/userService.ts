// frontend/src/lib/userService.ts
import { getAccessTokenSilently } from './authService'; // To get Auth0 token
import { user as auth0UserStore } from './store';      // To get the Auth0 user object from the store
import { get } from 'svelte/store';

/**
 * Interface for the payload sent to the backend when syncing user information.
 * This should match the expected structure of your Go backend's /api/users/sync endpoint.
 */
export interface UserSyncPayload {
  auth0_id: string;       // Typically the 'sub' claim from Auth0 user profile
  email: string | undefined;
  name?: string | undefined;
  picture?: string | undefined;
  // Add any other fields from the Auth0 user profile that your backend needs to store.
  // Example: nickname?: string;
}

/**
 * Synchronizes the authenticated Auth0 user's data with the application's backend database.
 * This function should be called after a user successfully logs in via Auth0
 * and their profile is available in the auth0UserStore.
 */
export async function syncUserToDb(): Promise<void> {
  const auth0User = get(auth0UserStore); // Get the current Auth0 user profile from the Svelte store

  // Ensure there is an Auth0 user and, critically, an ID (sub claim) to sync against.
  if (!auth0User || !auth0User.sub) {
    console.warn('syncUserToDb: No Auth0 user data or auth0_id (sub) found in store. Skipping sync.');
    return; // Exit if no user or sub (Auth0 ID) is found
  }

  try {
    // Obtain an access token from Auth0. This token will be used to authenticate the request to your backend.
    const accessToken = await getAccessTokenSilently();
    if (!accessToken) {
      console.error('syncUserToDb: Failed to obtain access token. Cannot sync user.');
      return; // Exit if no access token could be retrieved
    }

    // Prepare the payload with the data to be sent to the backend.
    const payload: UserSyncPayload = {
      auth0_id: auth0User.sub,       // 'sub' is the standard Auth0 unique user identifier
      email: auth0User.email,
      name: auth0User.name || auth0User.nickname, // Use name, fallback to nickname if available
      picture: auth0User.picture,
      // Map other relevant Auth0 user fields to your payload as needed:
      // nickname: auth0User.nickname,
      // given_name: auth0User.given_name,
      // family_name: auth0User.family_name,
    };

    // Make the API call to your backend's sync endpoint.
    const response = await fetch('/api/users/sync', { // Ensure this is your correct backend endpoint URI
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`, // Send the Auth0 access token for backend verification
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
    const responseData = await response.json();
    console.log('syncUserToDb: User synchronized successfully with backend.', responseData);

    // TODO: Optionally, update any local application state based on the backend's response.
    // For example, if your backend assigns an internal ID to the user and returns it,
    // you might store that in a Svelte store or use it to update the UI.

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
