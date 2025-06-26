import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { authService } from '../services/authService';

interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  picture?: string;
  isVerified: boolean;
}

// Helper function to map API user data to AuthContext User format
const mapApiUserToUser = (apiUser: any): User => {
  return {
    id: apiUser.id,
    email: apiUser.email,
    firstName: apiUser.first_name || apiUser.firstName,
    lastName: apiUser.last_name || apiUser.lastName,
    picture: apiUser.profile_photo_url || apiUser.google_photo_url || apiUser.picture,
    isVerified: apiUser.is_verified || apiUser.isVerified || false
  };
};

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (code: string, state?: string) => Promise<void>;
  logout: () => Promise<void>;
  getGoogleAuthUrl: () => Promise<string>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user;

  useEffect(() => {
    // Check if user is already authenticated on app load
    const checkAuth = async () => {
      try {
        const token = localStorage.getItem('accessToken');
        if (token) {
          const isValid = await authService.validateToken();
          if (isValid) {
            const profile = await authService.getProfile();
            setUser(mapApiUserToUser(profile));
          } else {
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
          }
        }
      } catch (error) {
        console.error('Auth check failed:', error);
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (code: string, state?: string) => {
    try {
      setIsLoading(true);
      const response = await authService.login(code, state);
      
      localStorage.setItem('accessToken', response.accessToken);
      localStorage.setItem('refreshToken', response.refreshToken);

      setUser(mapApiUserToUser(response.user));
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await authService.logout();
    } catch (error) {
      console.error('Logout failed:', error);
    } finally {
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      setUser(null);
    }
  };

  const getGoogleAuthUrl = async (): Promise<string> => {
    return await authService.getGoogleAuthUrl();
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    logout,
    getGoogleAuthUrl,
  };

  return (
    <AuthContext.Provider value={value}>
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
