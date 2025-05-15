// /home/catvader/repos/procurement/frontend/src/lib/userService.ts
import type { User } from '@auth0/auth0-spa-js';
import { getAccessTokenSilently } from './authService'; // Import the function to get the token

// Define a more specific interface for what you expect to store/send to your backend
interface AppUser {
  auth0_id: string; // Typically the 'sub' claim from Auth0 user
  email: string | undefined;
  name?: string | undefined; // Or nickname, given_name, etc.
  picture?: string | undefined;
  // Add any other fields you want to sync or manage
}

export async function syncUserToDb(auth0User: User | null | undefined): Promise<void> {
  if (!auth0User || !auth0User.sub) {
    console.warn('Auth0 user data (or sub claim) is not available for DB sync.');
    return;
  }

  // Map the Auth0 user to your application's user structure
  const appUser: AppUser = {
    auth0_id: auth0User.sub,
    email: auth0User.email,
    name: auth0User.name || auth0User.nickname, // Use name, fallback to nickname
    picture: auth0User.picture,
  };

  console.log('Attempting to sync user to DB. User details:', appUser);

  try {
    const accessToken = await getAccessTokenSilently();
    if (!accessToken) {
      console.error('Failed to get access token. Cannot sync user to DB.');
      // Potentially inform the user or retry
      return;
    }

    const backendUrl = 'http://localhost:8080'; // Assuming backend runs on 8080
    const response = await fetch(`${backendUrl}/api/users/sync`, { // Your backend endpoint
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
      },
      body: JSON.stringify(appUser),
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error(`API responded with ${response.status}: ${errorText}`);
      // Consider throwing an error or setting an error state for the UI to pick up
      throw new Error(`API responded with ${response.status}: ${errorText}`);
    }

    const result = await response.json(); // Assuming your backend sends back a JSON response
    console.log('User synced successfully with backend:', result);
    // You might want to update a Svelte store here if the backend returns new/updated user info

  } catch (error) {
    console.error('Failed to sync user to DB:', error);
    // Handle error appropriately - maybe set an error state in your Svelte store
    // For now, we'll just re-throw or log. UI could show a generic 'sync failed' message.
    // alert('Failed to sync user profile with the server. Please try again later.');
  }
}
