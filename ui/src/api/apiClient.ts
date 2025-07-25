import axios from 'axios';
import { Configuration, DefaultApi } from '../generated/api/client';

// Get runtime configuration
const getConfig = () => {
  // Try to get from window.WARDEN_CONFIG first (runtime)
  if (typeof window !== 'undefined' && window.WARDEN_CONFIG) {
    return window.WARDEN_CONFIG;
  }
  
  // Fallback to build-time environment variables
  return {
    API_BASE_URL: import.meta.env.VITE_API_BASE_URL || '/',
    VERSION: import.meta.env.VITE_VERSION || 'dev',
    BUILD_TIME: import.meta.env.VITE_BUILD_TIME || new Date().toISOString(),
  };
};

const config = getConfig();

// Create a base axios instance
const axiosInstance = axios.create({
  baseURL: config.API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to add auth token to requests
axiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle token refresh
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // If the error is 401 (Unauthorized) and we haven't already tried to refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = localStorage.getItem('refreshToken');
        
        if (!refreshToken) {
          // No refresh token available, redirect to login
          window.location.href = '/login';
          return Promise.reject(error);
        }
        
        // Create a new instance without interceptors to avoid infinite loop
        const refreshApi = new DefaultApi(
          new Configuration(),
          '',
          axios.create()
        );
        
        // Try to refresh the token
        const response = await refreshApi.refreshToken({
          refresh_token: refreshToken,
        });
        
        // Store the new tokens
        localStorage.setItem('accessToken', response.data.access_token);
        localStorage.setItem('refreshToken', response.data.refresh_token);
        
        // Update the original request with the new token
        originalRequest.headers.Authorization = `Bearer ${response.data.access_token}`;
        
        // Retry the original request
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        // If refresh fails, redirect to login
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);

// Create API client with the configured axios instance
const apiConfiguration = new Configuration({
  basePath: config.API_BASE_URL,
});

export const apiClient = new DefaultApi(
  apiConfiguration,
  apiConfiguration.basePath,
  axiosInstance
);

export { apiConfiguration, axiosInstance };


export default apiClient;