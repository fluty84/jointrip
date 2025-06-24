import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const AuthCallbackPage: React.FC = () => {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [error, setError] = useState<string>('');
  const [isProcessing, setIsProcessing] = useState(true);

  useEffect(() => {
    const handleCallback = async () => {
      try {
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        const error = searchParams.get('error');

        if (error) {
          setError('Authentication was cancelled or failed. Please try again.');
          setIsProcessing(false);
          return;
        }

        if (!code) {
          setError('No authorization code received. Please try again.');
          setIsProcessing(false);
          return;
        }

        // Perform login with the received code
        await login(code, state || undefined);
        
        // Redirect to home page on successful login
        navigate('/', { replace: true });
      } catch (error) {
        console.error('Login failed:', error);
        setError('Login failed. Please try again.');
        setIsProcessing(false);
      }
    };

    handleCallback();
  }, [searchParams, login, navigate]);

  if (isProcessing) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <h2 className="text-xl font-semibold text-gray-900 mb-2">Completing your login...</h2>
          <p className="text-gray-600">Please wait while we authenticate you with Google.</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="max-w-md mx-auto text-center">
          <div className="card">
            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">Authentication Failed</h2>
            <p className="text-gray-600 mb-6">{error}</p>
            <button
              onClick={() => navigate('/login')}
              className="btn-primary w-full"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    );
  }

  return null;
};
