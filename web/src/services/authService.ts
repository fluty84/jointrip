import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    picture?: string;
    isVerified: boolean;
  };
}

interface GoogleAuthUrlResponse {
  auth_url: string;
  state: string;
}

class AuthService {
  private apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  constructor() {
    // Add request interceptor to include auth token
    this.apiClient.interceptors.request.use(
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
    this.apiClient.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const refreshToken = localStorage.getItem('refreshToken');
            if (refreshToken) {
              const response = await this.refreshToken(refreshToken);
              localStorage.setItem('accessToken', response.data.accessToken);
              
              // Retry the original request with new token
              originalRequest.headers.Authorization = `Bearer ${response.data.accessToken}`;
              return this.apiClient(originalRequest);
            }
          } catch (refreshError) {
            // Refresh failed, redirect to login
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
            window.location.href = '/login';
          }
        }

        return Promise.reject(error);
      }
    );
  }

  async getGoogleAuthUrl(state?: string): Promise<string> {
    const response = await this.apiClient.get<GoogleAuthUrlResponse>('/auth/google/url', {
      params: { state },
    });
    return response.data.auth_url;
  }

  async login(code: string, state?: string): Promise<LoginResponse> {
    const response = await this.apiClient.post<LoginResponse>('/auth/google/login', {
      code,
      state,
    });

    // Store tokens in localStorage
    if (response.data.accessToken) {
      localStorage.setItem('accessToken', response.data.accessToken);
    }
    if (response.data.refreshToken) {
      localStorage.setItem('refreshToken', response.data.refreshToken);
    }

    return response.data;
  }

  async logout(): Promise<void> {
    try {
      await this.apiClient.post('/auth/logout');
    } finally {
      // Always clear tokens from localStorage
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    }
  }

  async refreshToken(refreshToken: string): Promise<any> {
    return await this.apiClient.post('/auth/refresh', {
      refreshToken: refreshToken,
    });
  }

  async validateToken(): Promise<boolean> {
    try {
      await this.apiClient.get('/auth/validate');
      return true;
    } catch {
      return false;
    }
  }

  async getProfile(): Promise<any> {
    const response = await this.apiClient.get('/profile');
    return response.data.user;
  }
}

export const authService = new AuthService();
