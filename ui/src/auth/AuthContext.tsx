import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { jwtDecode } from 'jwt-decode';
import type {LoginRequest, RefreshTokenRequest, ChangeUserPasswordRequest} from '../generated/api/client';
import type { User } from '../generated/api/client';
import apiClient from '../api/apiClient';
import type { AxiosError } from 'axios';
import axios from 'axios';

interface AuthContextType {
  isAuthenticated: boolean;
  user: User | null;
  login: (username: string, password: string) => Promise<void>;
  loginWithTokens: (accessToken: string, refreshToken: string) => Promise<void>;
  verify2FA: (code: string) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
  updateUserData: () => Promise<void>;
  isLoading: boolean;
  error: string | null;
  hasTmpPassword: boolean;
  is2FARequired: boolean;
  changePassword: (oldPassword: string, newPassword: string) => Promise<void>;
  is2FABlocked: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

function isAxiosError(error: unknown): error is AxiosError<{ message?: string; error?: { message?: string; code?: string; session_id?: string } }> {
  return axios.isAxiosError(error);
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasTmpPassword, setHasTmpPassword] = useState<boolean>(false);
  const [is2FARequired, setIs2FARequired] = useState<boolean>(false);
  const [sessionId, setSessionId] = useState<string>("");
  const [is2FABlocked, setIs2FABlocked] = useState(false);

  // Check if token exists and is valid on mount
  useEffect(() => {
    const checkAuth = async () => {
      const accessToken = localStorage.getItem('accessToken');
      if (!accessToken) return;

      try {
        // Check if token is expired
        const decoded = jwtDecode<Partial<User> & { exp?: number }>(accessToken);
        const currentTime = Date.now() / 1000;

        if (decoded.exp && decoded.exp < currentTime) {
          // Token is expired, try to refresh
          await refreshToken();
        } else {
          // Token is valid
          setIsAuthenticated(true);
          // setUser(decoded as User); // Не используем, ждем getCurrentUser

          // Fetch user data
          try {
            const userData = await apiClient.getCurrentUser();
            setUser(userData.data);
          } catch (error) {
            console.error('Failed to fetch user data:', error);
          }
        }
      } catch (error) {
        console.error('Invalid token:', error);
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
      }
    };

    checkAuth();
  }, []);

  const login = async (username: string, password: string) => {
    setIsLoading(true);
    setError(null);
    setIs2FARequired(false);
    setSessionId("");
    try {
      const loginRequest: LoginRequest = {
        username,
        password
      };
      const response = await apiClient.login(loginRequest);
      localStorage.setItem('accessToken', response.data.access_token);
      localStorage.setItem('refreshToken', response.data.refresh_token);
      setIsAuthenticated(true);
      setHasTmpPassword(response.data.is_tmp_password || false);
      try {
        const userData = await apiClient.getCurrentUser();
        setUser(userData.data);
      } catch (error) {
        console.error('Failed to fetch user data:', error);
      }
    } catch (error: any) {
      console.error('Login failed:', error);
      // Проверяем ошибку 2FA
      if (error?.response?.status === 403 && error?.response?.data?.error?.code === '2fa_required') {
        setIs2FARequired(true);
        setSessionId(error.response.data.error.session_id || "");
        setError(null);
      } else if (isAxiosError(error) && error.response?.data?.message) {
        setError(error.response.data.message);
      } else {
        setError('Login failed. Please check your credentials.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const verify2FA = async (code: string) => {
    setIsLoading(true);
    setError(null);
    try {
      // sessionId можно использовать, если потребуется (например, для защиты от перебора)
      const resp = await apiClient.verify2FA({
        code,
        session_id: sessionId || "",
      });
      localStorage.setItem('accessToken', resp.data.access_token);
      localStorage.setItem('refreshToken', resp.data.refresh_token);
      setIsAuthenticated(true);
      setIs2FARequired(false);
      setSessionId("");
      setHasTmpPassword(false);
      try {
        const userData = await apiClient.getCurrentUser();
        setUser(userData.data);
      } catch (error) {
        console.error('Failed to fetch user data:', error);
      }
    } catch (error: any) {
      if (isAxiosError(error) && error.response?.status === 429) {
        setError('Too many attempts. Please try again later.');
        setIs2FABlocked(true);
        setTimeout(() => setIs2FABlocked(false), 60000);
      } else if (isAxiosError(error) && error.response?.data?.error?.message) {
        setError(error.response.data.error.message);
      } else {
        setError('Invalid 2FA code.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const refreshToken = async () => {
    const currentRefreshToken = localStorage.getItem('refreshToken');
    if (!currentRefreshToken) {
      setIsAuthenticated(false);
      setUser(null);
      return;
    }

    try {
      const refreshRequest: RefreshTokenRequest = {
        refresh_token: currentRefreshToken
      };

      const response = await apiClient.refreshToken(refreshRequest);

      localStorage.setItem('accessToken', response.data.access_token);
      localStorage.setItem('refreshToken', response.data.refresh_token);

      setIsAuthenticated(true);

      // Fetch user data
      try {
        const userData = await apiClient.getCurrentUser();
        setUser(userData.data);
      } catch (error) {
        console.error('Failed to fetch user data:', error);
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
      // If refresh fails, log the user out
      logout();
    }
  };

  const loginWithTokens = async (accessToken: string, refreshToken: string) => {
    setIsLoading(true);
    setError(null);
    try {
      // Сохраняем токены
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
      setIsAuthenticated(true);

      try {
        const userData = await apiClient.getCurrentUser();
        setUser(userData.data);
      } catch (error) {
        console.error('Failed to fetch user data:', error);
      }
    } catch (error: any) {
      console.error('Login with tokens failed:', error);
      if (isAxiosError(error) && error.response?.data?.message) {
        setError(error.response.data.message);
      } else {
        setError('Authentication failed. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    setIsAuthenticated(false);
    setUser(null);
    setHasTmpPassword(false);
  };

  const changePassword = async (oldPassword: string, newPassword: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const changePasswordRequest: ChangeUserPasswordRequest = {
        old_password: oldPassword,
        new_password: newPassword
      };

      await apiClient.userChangeMyPassword(changePasswordRequest);
      
      // Update user data to get the latest is_tmp_password status
      await updateUserData();
    } catch (error: unknown) {
      console.error('Password change failed:', error);
      if (isAxiosError(error) && error.response?.data?.message) {
        setError(error.response.data.message);
      } else {
        setError('Password change failed. Please try again.');
      }
      throw error;
    } finally {
      setIsLoading(false);
    }
  };



  const updateUserData = async () => {
    if (!isAuthenticated) return;
    
    try {
      const userData = await apiClient.getCurrentUser();
      setUser(userData.data);
      // Update hasTmpPassword based on user data
      setHasTmpPassword(userData.data.is_tmp_password || false);
      // Don't return userData.data to match the Promise<void> return type
    } catch (error) {
      console.error('Failed to fetch user data:', error);
    }
  };


  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        user,
        login,
        loginWithTokens,
        verify2FA,
        logout,
        refreshToken,
        updateUserData,
        isLoading,
        error,
        hasTmpPassword,
        is2FARequired,
        changePassword,
        is2FABlocked
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
