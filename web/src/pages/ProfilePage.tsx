import React from 'react';
import { useAuth } from '../contexts/AuthContext';

export const ProfilePage: React.FC = () => {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Profile</h1>
        <p className="text-gray-600">Manage your account and travel preferences</p>
      </div>

      <div className="grid md:grid-cols-3 gap-8">
        {/* Profile Card */}
        <div className="md:col-span-1">
          <div className="card text-center">
            {user.picture ? (
              <img
                src={user.picture}
                alt={`${user.firstName} ${user.lastName}`}
                className="w-24 h-24 rounded-full mx-auto mb-4"
              />
            ) : (
              <div className="w-24 h-24 bg-gray-300 rounded-full flex items-center justify-center mx-auto mb-4">
                <span className="text-gray-600 text-2xl font-medium">
                  {user.firstName?.[0]}{user.lastName?.[0]}
                </span>
              </div>
            )}
            
            <h2 className="text-xl font-semibold text-gray-900 mb-1">
              {user.firstName} {user.lastName}
            </h2>
            <p className="text-gray-600 mb-4">{user.email}</p>
            
            <div className="flex items-center justify-center space-x-2 mb-4">
              {user.isVerified ? (
                <>
                  <svg className="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                  <span className="text-green-600 text-sm font-medium">Verified</span>
                </>
              ) : (
                <>
                  <svg className="w-5 h-5 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                  <span className="text-yellow-600 text-sm font-medium">Pending Verification</span>
                </>
              )}
            </div>

            <button className="btn-primary w-full">
              Edit Profile
            </button>
          </div>
        </div>

        {/* Profile Details */}
        <div className="md:col-span-2 space-y-6">
          {/* Basic Information */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Basic Information</h3>
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  First Name
                </label>
                <input
                  type="text"
                  value={user.firstName}
                  readOnly
                  className="input-field bg-gray-50"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Last Name
                </label>
                <input
                  type="text"
                  value={user.lastName}
                  readOnly
                  className="input-field bg-gray-50"
                />
              </div>
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  value={user.email}
                  readOnly
                  className="input-field bg-gray-50"
                />
              </div>
            </div>
          </div>

          {/* Travel Preferences */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Travel Preferences</h3>
            <p className="text-gray-600 text-center py-8">
              Travel preferences will be available in the next update.
              <br />
              <span className="text-sm">Coming soon: Bio, interests, travel style, and more!</span>
            </p>
          </div>

          {/* Account Settings */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Account Settings</h3>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Email Notifications</h4>
                  <p className="text-sm text-gray-600">Receive updates about your trips and messages</p>
                </div>
                <button className="btn-secondary text-sm">
                  Configure
                </button>
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Privacy Settings</h4>
                  <p className="text-sm text-gray-600">Control who can see your profile and trips</p>
                </div>
                <button className="btn-secondary text-sm">
                  Manage
                </button>
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Identity Verification</h4>
                  <p className="text-sm text-gray-600">Verify your identity to build trust</p>
                </div>
                <button className="btn-primary text-sm">
                  {user.isVerified ? 'Verified' : 'Verify Now'}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
