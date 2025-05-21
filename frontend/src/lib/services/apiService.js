// src/lib/services/apiService.js

// Use the environment variable with a fallback for development
const BASE_URL = import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8080/api';

// Add auth token to requests if available
function getAuthHeader() {
  const token = localStorage.getItem('auth_token');
  return token ? { 'Authorization': `Bearer ${token}` } : { 'Authorization': '' };
}

/**
 * Helper function for making API requests.
 * @param {string} endpoint The API endpoint (e.g., '/users').
 * @param {string} [method='GET'] The HTTP method.
 * @param {Record<string, any> | null} [body=null] The request body for POST/PUT requests.
 * @param {Record<string, string>} [additionalHeaders={}] Additional headers.
 * @returns {Promise<any>} A promise that resolves with the JSON response.
 */
export async function request(endpoint, method = 'GET', body = null, additionalHeaders = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...getAuthHeader(),
        ...additionalHeaders
    };

    /** @type {RequestInit} */
    const config = {
        method: method,
        headers: headers,
    };

    if (body) {
        config.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${BASE_URL}${endpoint}`, config);
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: response.statusText }));
            throw new Error(errorData.error || errorData.message || `HTTP error! status: ${response.status}`);
        }
        if (response.status === 204) { // No Content
            return null;
        }
        return await response.json();
    } catch (error) {
        console.error(`API request error to ${method} ${endpoint}:`, error);
        throw error; // Re-throw to be caught by the caller
    }
}

// --- Tender Management ---
export const tenderService = {
    /**
     * Get all tenders
     * @returns {Promise<Array<import('$lib/types').Tender>>}
     */
    getAllTenders: () => request('/tenders', 'GET'),
    
    /**
     * Get a single tender by ID
     * @param {string | number} id
     * @returns {Promise<import('$lib/types').Tender>}
     */
    getTenderById: (id) => request(`/tenders/${id}`, 'GET'),
    /** @param {Record<string, any>} tenderData */
    createTender: (tenderData) => request('/tenders', 'POST', tenderData),
    /**
     * @param {string | number} id
     * @param {Record<string, any>} tenderData
     */
    updateTender: (id, tenderData) => request(`/tenders/${id}`, 'PUT', tenderData),
    /** @param {string | number} id */
    deleteTender: (id) => request(`/tenders/${id}`, 'DELETE'),
};

// --- Bid Management ---
export const bidService = {
    /**
     * @param {string | number} tenderId
     * @returns {Promise<Record<string, any>[]>} A promise that resolves with the JSON response.
     */
    getBidsForTender: (tenderId) => request(`/tenders/${tenderId}/bids`, 'GET'),
    /**
     * @param {string | number} tenderId
     * @param {Record<string, any>} bidData
     * @returns {Promise<Record<string, any>>} A promise that resolves with the JSON response.
     */
    createBidForTender: (tenderId, bidData) => request(`/tenders/${tenderId}/bids`, 'POST', bidData),
    /**
     * @param {string | number} bidId
     * @returns {Promise<Record<string, any>>} A promise that resolves with the JSON response.
     */
    getBidById: (bidId) => request(`/bids/${bidId}`, 'GET'),
    /**
     * @param {string | number} bidId
     * @param {Record<string, any>} bidData
     * @returns {Promise<Record<string, any>>} A promise that resolves with the JSON response.
     */
    updateBid: (bidId, bidData) => request(`/bids/${bidId}`, 'PUT', bidData),
    /**
     * @param {string | number} bidId
     * @returns {Promise<null>} A promise that resolves with null.
     */
    deleteBid: (bidId) => request(`/bids/${bidId}`, 'DELETE'),
};

// --- Requisition Management ---
export const requisitionService = {
    /**
     * Get all requisitions
     * @returns {Promise<Array<import('$lib/types').Requisition>>}
     */
    getAllRequisitions: () => request('/requisitions', 'GET'),
    
    /**
     * Get a single requisition by ID
     * @param {string | number} id
     * @returns {Promise<import('$lib/types').Requisition>}
     */
    getRequisitionById: (id) => request(`/requisitions/${id}`, 'GET'),
    
    /**
     * Create a new requisition
     * @param {Omit<import('$lib/types').Requisition, 'id' | 'created_at' | 'updated_at'>} requisitionData
     * @returns {Promise<import('$lib/types').Requisition>}
     */
    createRequisition: (requisitionData) => request('/requisitions', 'POST', requisitionData),
    
    /**
     * Update a requisition
     * @param {string | number} id
     * @param {Partial<import('$lib/types').Requisition>} requisitionData
     * @returns {Promise<import('$lib/types').Requisition>}
     */
    updateRequisition: (id, requisitionData) => request(`/requisitions/${id}`, 'PUT', requisitionData),
    
    /**
     * Delete a requisition
     * @param {string | number} id
     * @returns {Promise<void>}
     */
    deleteRequisition: (id) => request(`/requisitions/${id}`, 'DELETE')
};

// --- Asset Management ---
export const assetService = {
    /**
     * Get all assets
     * @returns {Promise<Array<import('$lib/types').Asset>>}
     */
    getAllAssets: () => request('/assets', 'GET'),
    
    /**
     * Get a single asset by ID
     * @param {string | number} id
     * @returns {Promise<import('$lib/types').Asset>}
     */
    getAssetById: (id) => request(`/assets/${id}`, 'GET'),
    
    /**
     * Create a new asset
     * @param {Omit<import('$lib/types').Asset, 'id'>} assetData
     * @returns {Promise<import('$lib/types').Asset>}
     */
    createAsset: (assetData) => request('/assets', 'POST', assetData),
    
    /**
     * Update an asset
     * @param {string | number} id
     * @param {Partial<import('$lib/types').Asset>} assetData
     * @returns {Promise<import('$lib/types').Asset>}
     */
    updateAsset: (id, assetData) => request(`/assets/${id}`, 'PUT', assetData),
    
    /**
     * Delete an asset
     * @param {string | number} id
     * @returns {Promise<void>}
     */
    deleteAsset: (id) => request(`/assets/${id}`, 'DELETE')
};

// --- Document Management ---
// Note: File uploads require FormData instead of JSON
export const documentService = {
    // ...
};

// --- Evaluation Management ---
// We'll add these later
export const evaluationService = {
    // ...
};
