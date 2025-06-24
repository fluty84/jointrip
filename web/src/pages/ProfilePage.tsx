import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { authService } from '../services/authService';

interface User {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  bio?: string;
  location?: string;
  phone?: string;
  website?: string;
  languages: string[] | null;
  interests: string[] | null;
  travel_style?: string;
  profile_photo_url: string;
  google_photo_url: string;
  reputation_score: number;
  rating_average: number;
  rating_count: number;
  profile_completion_percentage: number;
  created_at: string;
}

export const ProfilePage: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const userData = await authService.getProfile();
      setUser(userData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load profile');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
          <div className="text-center">
            <h3 className="mt-2 text-sm font-medium text-gray-900">Error Loading Profile</h3>
            <p className="mt-1 text-sm text-gray-500">{error}</p>
            <div className="mt-6">
              <button
                onClick={fetchProfile}
                className="btn-primary"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!user) {
    return <div>No profile data available</div>;
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
            <img
              src={user.profile_photo_url || user.google_photo_url || '/default-avatar.png'}
              alt={`${user.first_name} ${user.last_name}`}
              className="w-24 h-24 rounded-full mx-auto mb-4 object-cover"
            />

            <h2 className="text-xl font-semibold text-gray-900 mb-1">
              {user.first_name} {user.last_name}
            </h2>
            <p className="text-gray-600 mb-2">@{user.username}</p>
            <p className="text-gray-600 mb-4">{user.email}</p>

            <div className="flex items-center justify-center space-x-2 mb-4">
              <span className="text-yellow-400">⭐</span>
              <span className="text-sm text-gray-600">
                {user.rating_average.toFixed(1)} ({user.rating_count} reviews)
              </span>
            </div>

            <div className="mb-4">
              <div className="text-sm text-gray-600 mb-2">
                Profile {user.profile_completion_percentage}% complete
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-primary-600 h-2 rounded-full"
                  style={{ width: `${user.profile_completion_percentage}%` }}
                ></div>
              </div>
            </div>

            <Link to="/profile/edit" className="btn-primary w-full block text-center">
              Edit Profile
            </Link>
          </div>
        </div>

        {/* Profile Details */}
        <div className="md:col-span-2 space-y-6">
          {/* About */}
          {user.bio && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">About</h3>
              <p className="text-gray-700">{user.bio}</p>
            </div>
          )}

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
                  value={user.first_name}
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
                  value={user.last_name}
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
              {user.location && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Location
                  </label>
                  <input
                    type="text"
                    value={user.location}
                    readOnly
                    className="input-field bg-gray-50"
                  />
                </div>
              )}
              {user.phone && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Phone
                  </label>
                  <input
                    type="text"
                    value={user.phone}
                    readOnly
                    className="input-field bg-gray-50"
                  />
                </div>
              )}
              {user.website && (
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Website
                  </label>
                  <a
                    href={user.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary-600 hover:text-primary-500"
                  >
                    {user.website}
                  </a>
                </div>
              )}
            </div>
          </div>

          {/* Travel Preferences */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Travel Preferences</h3>
            <div className="space-y-4">
              {user.travel_style && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Travel Style
                  </label>
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-primary-100 text-primary-800">
                    {user.travel_style}
                  </span>
                </div>
              )}

              {user.languages && user.languages.length > 0 && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Languages
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {user.languages.map((language, index) => (
                      <span
                        key={index}
                        className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800"
                      >
                        {language}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {user.interests && user.interests.length > 0 && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Interests
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {user.interests.map((interest, index) => (
                      <span
                        key={index}
                        className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800"
                      >
                        {interest}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {!user.travel_style && (!user.languages || user.languages.length === 0) && (!user.interests || user.interests.length === 0) && (
                <p className="text-gray-600 text-center py-8">
                  No travel preferences set yet.
                  <br />
                  <Link to="/profile/edit" className="text-primary-600 hover:text-primary-500">
                    Add your preferences
                  </Link>
                </p>
              )}
            </div>
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

            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
