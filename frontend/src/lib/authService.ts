// src/lib/authService.ts
import { isAuthenticated, user, isLoading, authError } from './store';
import type { AppUser } from './store';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';

// Define interfaces for authentication requests and responses
interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  first_name?: string;
  last_name?: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface AuthResponse {
  token: string;
  user: AppUser;
}

interface PasswordChangeRequest {
  current_password: string;
  new_password: string;
}

interface PasswordResetRequest {
  email: string;
}

interface PasswordResetConfirmRequest {
  token: string;
  new_password: string;
}

// API base URL
const API_BASE_URL = '${PUBLIC_VITE_API_BASE_URL}/api';

// Token storage key
const TOKEN_STORAGE_KEY = 'auth_token';
const USER_STORAGE_KEY = 'auth_user';

// Helper function to get stored token
export function getStoredToken(): string | null {
  if (!browser) return null;
  return localStorage.getItem(TOKEN_STORAGE_KEY);
}

// Helper function to get stored user
function getStoredUser(): AppUser | null {
  if (!browser) return null;
  const userJson = localStorage.getItem(USER_STORAGE_KEY);
  return userJson ? JSON.parse(userJson) : null;
}

// Helper function to store authentication data
function storeAuthData(token: string, userData: AppUser): void {
  if (!browser) return;
  localStorage.setItem(TOKEN_STORAGE_KEY, token);
  localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(userData));
}

// Helper function to clear authentication data
function clearAuthData(): void {
  if (!browser) return;
  localStorage.removeItem(TOKEN_STORAGE_KEY);
  localStorage.removeItem(USER_STORAGE_KEY);
}

// Initialize authentication state from local storage
export async function initializeAuth(): Promise<void> {
  if (!browser) return;

  isLoading.set(true);

  try {
    const token = getStoredToken();
    const storedUser = getStoredUser();

    if (token && storedUser) {
      user.set(storedUser);
      isAuthenticated.set(true);
      authError.set(null);
    } else {
      isAuthenticated.set(false);
      user.set(null);
      authError.set(null);
    }
  } catch (error: any) {
    console.error('Error initializing auth:', error);
    authError.set(error);
    isAuthenticated.set(false);
    user.set(null);
    clearAuthData();
  } finally {
    isLoading.set(false);
  }
}

// Register a new user
export async function register(email: string, username: string, password: string, firstName?: string, lastName?: string): Promise<void> {
  isLoading.set(true);
  authError.set(null);

  try {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        username,
        password,
        first_name: firstName,
        last_name: lastName,
      }),
    });

    if (!response.ok) {
      const errorData = await response.text();
      throw new Error(errorData || 'Registration failed');
    }

    const data: AuthResponse = await response.json();

    // Store token and user data
    storeAuthData(data.token, data.user);

    // Update stores
    user.set(data.user);
    isAuthenticated.set(true);

    // Navigate to home page
    await goto('/');
  } catch (error: any) {
    console.error('Registration error:', error);
    authError.set(error);
    throw error;
  } finally {
    isLoading.set(false);
  }
}

// Login user
export async function login(email: string, password: string): Promise<void> {
  isLoading.set(true);
  authError.set(null);

  try {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const errorData = await response.text();
      throw new Error(errorData || 'Login failed');
    }

    const data: AuthResponse = await response.json();

    // Store token and user data
    storeAuthData(data.token, data.user);

    // Update stores
    user.set(data.user);
    isAuthenticated.set(true);

    // Navigate to home page
    await goto('/');
  } catch (error: any) {
    console.error('Login error:', error);
    authError.set(error);
    throw error;
  } finally {
    isLoading.set(false);
  }
}

// Logout user
export async function logout(): Promise<void> {
  isLoading.set(true);

  try {
    const token = getStoredToken();

    if (token) {
      // Call logout endpoint to invalidate token
      await fetch(`${API_BASE_URL}/auth/logout`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
    }
  } catch (error) {
    console.error('Logout error:', error);
    // Continue with logout even if server request fails
  } finally {
    // Clear local storage and update stores
    clearAuthData();
    user.set(null);
    isAuthenticated.set(false);
    isLoading.set(false);

    // Navigate to login page
    await goto('/login');
  }
}

// Change password
export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  const token = getStoredToken();
  if (!token) {
    throw new Error('Not authenticated');
  }

  const response = await fetch(`${API_BASE_URL}/auth/password/change`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      current_password: currentPassword,
      new_password: newPassword
    }),
  });

  if (!response.ok) {
    const errorData = await response.text();
    throw new Error(errorData || 'Password change failed');
  }
}

// Request password reset
export async function requestPasswordReset(email: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/auth/password/reset/request`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email }),
  });

  if (!response.ok) {
    const errorData = await response.text();
    throw new Error(errorData || 'Password reset request failed');
  }
}

// Reset password with token
export async function resetPassword(token: string, newPassword: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/auth/password/reset/confirm`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      token,
      new_password: newPassword
    }),
  });

  if (!response.ok) {
    const errorData = await response.text();
    throw new Error(errorData || 'Password reset failed');
  }
}

// Check if user is authenticated
export function checkAuthStatus(): boolean {
  if (!browser) return false;

  const token = getStoredToken();
  return !!token;
}

// Get access token for API calls
export function getAccessToken(): string | null {
  return getStoredToken();
}

// Get current user
export function getCurrentUser(): AppUser | null {
  return getStoredUser();
}