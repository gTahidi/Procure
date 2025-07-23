// frontend/src/lib/services/apiService.ts
import { getAccessToken } from '../authService';
import { goto } from '$app/navigation';

const API_BASE_URL = 'http://localhost:8080/api';

interface RequestOptions extends RequestInit {
  headers?: Record<string, string>;
}

/**
 * Makes an API request with authentication
 * @param endpoint - API endpoint path
 * @param options - Request options
 * @returns Response data
 */
async function apiRequest(endpoint: string, options: RequestOptions = {}): Promise<any> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // Set default headers
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };
  
  // Add authentication token if available
  const token = getAccessToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  // Make the request
  const response = await fetch(url, {
    ...options,
    headers,
  });
  
  // Handle authentication errors
  if (response.status === 401) {
    // Redirect to login page if unauthorized
    goto('/login');
    throw new Error('Authentication required');
  }
  
  // Parse response
  if (response.status >= 200 && response.status < 300) {
    // Check if response is empty
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    }
    return await response.text();
  }
  
  // Handle errors
  const error = await response.text();
  throw new Error(error || `Request failed with status ${response.status}`);
}

/**
 * GET request
 * @param endpoint - API endpoint
 * @param options - Additional options
 * @returns Response data
 */
export function get(endpoint: string, options: RequestOptions = {}): Promise<any> {
  return apiRequest(endpoint, {
    method: 'GET',
    ...options,
  });
}

/**
 * POST request
 * @param endpoint - API endpoint
 * @param data - Request body data
 * @param options - Additional options
 * @returns Response data
 */
export function post(endpoint: string, data: any, options: RequestOptions = {}): Promise<any> {
  return apiRequest(endpoint, {
    method: 'POST',
    body: JSON.stringify(data),
    ...options,
  });
}

/**
 * PUT request
 * @param endpoint - API endpoint
 * @param data - Request body data
 * @param options - Additional options
 * @returns Response data
 */
export function put(endpoint: string, data: any, options: RequestOptions = {}): Promise<any> {
  return apiRequest(endpoint, {
    method: 'PUT',
    body: JSON.stringify(data),
    ...options,
  });
}

/**
 * DELETE request
 * @param endpoint - API endpoint
 * @param options - Additional options
 * @returns Response data
 */
export function del(endpoint: string, options: RequestOptions = {}): Promise<any> {
  return apiRequest(endpoint, {
    method: 'DELETE',
    ...options,
  });
}

/**
 * Upload file(s)
 * @param endpoint - API endpoint
 * @param formData - Form data with files
 * @param options - Additional options
 * @returns Response data
 */
export function upload(endpoint: string, formData: FormData, options: RequestOptions = {}): Promise<any> {
  // Don't set Content-Type header for multipart/form-data
  // The browser will set it automatically with the boundary
  const headers: Record<string, string> = {};
  
  // Add authentication token if available
  const token = getAccessToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return apiRequest(endpoint, {
    method: 'POST',
    headers,
    body: formData,
    ...options,
  });
}

export default {
  get,
  post,
  put,
  del,
  upload,
};