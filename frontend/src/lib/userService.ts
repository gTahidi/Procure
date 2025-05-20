// User service for handling user-related operations
import { request } from './services/apiService';

// Re-export the request function for direct use if needed
export { request } from './services/apiService';

// User type matching the backend model
export interface User {
  id: number; 
  username: string;
  email: string;
  name?: string;  
  phone_number?: string; 
  user_type: 'organization' | 'supplier'; 
  created_at?: string; 
  updated_at?: string; 
}

// Login request type
export interface LoginCredentials {
  email: string;
  password: string;
}

// Login response type
interface LoginResponse {
  user: User;
  token: string;
  expiresIn: number;
}

// User registration type
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone_number?: string;
  user_type: User['user_type'];
}

/**
 * Logs in a user with username and password
 * @param credentials The login credentials
 * @returns The User object or relevant response.
 */
export async function login(credentials: LoginCredentials): Promise<User | null> {
  try {
    // This request might be intended for a different backend endpoint or behavior
    // than the one used in authService. For now, assuming it mirrors /auth/login response.
    const response = await request('/auth/login', 'POST', credentials) as User;
    
    if (!response || !response.id) {
      console.error('userService.login: Invalid response from server');
      return null;
    }
    return response; // Returns the user object
  } catch (error) {
    console.error('userService.login error:', error);
    throw error; // Re-throw to be handled by caller or logged
  }
}

/**
 * Registers a new user
 * @param userData The user registration data
 * @returns The created user data
 */
export async function register(userData: RegisterRequest): Promise<User> {
  try {
    return await request('/auth/register', 'POST', userData);
  } catch (error) {
    console.error('Registration failed:', error);
    throw error;
  }
}

/**
 * Logs out the current user
 */
export function logout(): void {
  // Remove auth data from localStorage
  localStorage.removeItem('authToken');
  localStorage.removeItem('tokenExpiry');
  localStorage.removeItem('user');
  
  // Redirect to login page or home page
  window.location.href = '/login';
}

/**
 * Checks if the user is authenticated
 * @returns boolean indicating if the user is authenticated
 */
export function isAuthenticated(): boolean {
  const token = localStorage.getItem('authToken');
  if (!token) return false;
  
  // Check token expiration if expiry is stored
  const expiry = localStorage.getItem('tokenExpiry');
  if (expiry) {
    const now = new Date().getTime();
    if (now > parseInt(expiry)) {
      logout(); // Token expired
      return false;
    }
  }
  
  return true;
}

/**
 * Gets the current user's data by calling /auth/me
 * @returns The current user or null if not authenticated
 */
export async function getCurrentUser(): Promise<User | null> {
  // Session cookie is handled by the browser, no token to check in localStorage directly
  try {
    // Try to get user data from the API (relies on http-only session cookie)
    const userDataFromApi = await request('/auth/me', 'GET');
    
    if (userDataFromApi && userDataFromApi.id) {
      // Directly return the user data, assuming it matches the User interface
      // Add any necessary transformations if backend structure differs slightly from frontend User interface
      return {
        id: userDataFromApi.id,
        username: userDataFromApi.username,
        email: userDataFromApi.email,
        name: userDataFromApi.name || userDataFromApi.username, // Add name derivation if needed
        phone_number: userDataFromApi.phone_number,
        user_type: userDataFromApi.user_type,
        created_at: userDataFromApi.created_at,
        updated_at: userDataFromApi.updated_at,
      } as User;
    }
    return null; // No user data or invalid data
  } catch (error: any) {
    console.error('Failed to get current user (userService.getCurrentUser):', error);
    // Don't remove 'authToken' as we are not using it. Error implies session is invalid or network issue.
    return null;
  }
}

/**
 * Updates the current user's data
 * @param userData The updated user data
 * @returns The updated user
 */
export async function updateUser(userData: Partial<User>): Promise<User> {
  try {
    const response = await request('/auth/me', 'PUT', userData);
    // Update the stored user data
    if (response.user) {
      localStorage.setItem('user', JSON.stringify(response.user));
    }
    return response.user || response;
  } catch (error) {
    console.error('Failed to update user:', error);
    throw error;
  }
}

/**
 * Changes the current user's password
 * @param currentPassword The current password
 * @param newPassword The new password
 */
export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  try {
    await request('/auth/change-password', 'POST', {
      currentPassword,
      newPassword
    });
  } catch (error) {
    console.error('Failed to change password:', error);
    throw error;
  }
}

/**
 * Syncs user data to the database
 * @param user The user data to sync
 */
export async function syncUserToDb(user: User): Promise<void> {
  try {
    if (!user?.id) {
      console.warn('Cannot sync user: Missing user ID');
      return;
    }
    
    await request(`/users/${user.id}/sync`, 'POST', user);
  } catch (error) {
    console.error('Failed to sync user to database:', error);
    // Don't throw the error to prevent blocking the auth flow
  }
}
