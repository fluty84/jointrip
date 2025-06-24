import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
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
  profile_visibility: string;
  email_notifications: boolean;
  push_notifications: boolean;
}

interface FormData {
  first_name: string;
  last_name: string;
  bio: string;
  location: string;
  phone: string;
  website: string;
  languages: string[];
  interests: string[];
  travel_style: string;
  profile_visibility: string;
  email_notifications: boolean;
  push_notifications: boolean;
}

export const EditProfilePage: React.FC = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [formData, setFormData] = useState<FormData>({
    first_name: '',
    last_name: '',
    bio: '',
    location: '',
    phone: '',
    website: '',
    languages: [],
    interests: [],
    travel_style: '',
    profile_visibility: 'public',
    email_notifications: true,
    push_notifications: true,
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const userData = await authService.getProfile();
      setUser(userData);
      
      // Populate form with current data
      setFormData({
        first_name: userData.first_name || '',
        last_name: userData.last_name || '',
        bio: userData.bio || '',
        location: userData.location || '',
        phone: userData.phone || '',
        website: userData.website || '',
        languages: userData.languages || [],
        interests: userData.interests || [],
        travel_style: userData.travel_style || '',
        profile_visibility: userData.profile_visibility || 'public',
        email_notifications: userData.email_notifications ?? true,
        push_notifications: userData.push_notifications ?? true,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load profile');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setSuccess('');

    try {
      await authService.updateProfile(formData);
      setSuccess('Profile updated successfully!');
      setTimeout(() => {
        navigate('/profile');
      }, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update profile');
    } finally {
      setSaving(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    
    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData(prev => ({ ...prev, [name]: checked }));
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
  };

  const handleArrayInput = (name: 'languages' | 'interests', value: string) => {
    const items = value.split(',').map(item => item.trim()).filter(item => item);
    setFormData(prev => ({ ...prev, [name]: items }));
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!user) {
    return <div>No profile data available</div>;
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Edit Profile</h1>
            <p className="text-gray-600">Update your personal information and preferences</p>
          </div>
          <Link
            to="/profile"
            className="btn-secondary"
          >
            Cancel
          </Link>
        </div>
      </div>

      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Error</h3>
              <div className="mt-2 text-sm text-red-700">{error}</div>
            </div>
          </div>
        </div>
      )}

      {success && (
        <div className="mb-6 bg-green-50 border border-green-200 rounded-md p-4">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-green-800">Success</h3>
              <div className="mt-2 text-sm text-green-700">{success}</div>
            </div>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
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
                name="first_name"
                value={formData.first_name}
                onChange={handleInputChange}
                className="input-field"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Last Name
              </label>
              <input
                type="text"
                name="last_name"
                value={formData.last_name}
                onChange={handleInputChange}
                className="input-field"
                required
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Bio
              </label>
              <textarea
                name="bio"
                value={formData.bio}
                onChange={handleInputChange}
                rows={3}
                className="input-field"
                placeholder="Tell others about yourself..."
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Location
              </label>
              <input
                type="text"
                name="location"
                value={formData.location}
                onChange={handleInputChange}
                className="input-field"
                placeholder="City, Country"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Phone
              </label>
              <input
                type="tel"
                name="phone"
                value={formData.phone}
                onChange={handleInputChange}
                className="input-field"
                placeholder="+1 (555) 123-4567"
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Website
              </label>
              <input
                type="url"
                name="website"
                value={formData.website}
                onChange={handleInputChange}
                className="input-field"
                placeholder="https://yourwebsite.com"
              />
            </div>
          </div>
        </div>

        {/* Travel Preferences */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Travel Preferences</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Travel Style
              </label>
              <select
                name="travel_style"
                value={formData.travel_style}
                onChange={handleInputChange}
                className="input-field"
              >
                <option value="">Select travel style</option>
                <option value="budget">Budget</option>
                <option value="mid-range">Mid-range</option>
                <option value="luxury">Luxury</option>
                <option value="backpacker">Backpacker</option>
                <option value="adventure">Adventure</option>
                <option value="cultural">Cultural</option>
                <option value="relaxation">Relaxation</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Languages (comma-separated)
              </label>
              <input
                type="text"
                value={formData.languages.join(', ')}
                onChange={(e) => handleArrayInput('languages', e.target.value)}
                className="input-field"
                placeholder="English, Spanish, French"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Interests (comma-separated)
              </label>
              <input
                type="text"
                value={formData.interests.join(', ')}
                onChange={(e) => handleArrayInput('interests', e.target.value)}
                className="input-field"
                placeholder="Photography, Hiking, Food, Culture"
              />
            </div>
          </div>
        </div>

        {/* Privacy Settings */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Privacy & Notifications</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Profile Visibility
              </label>
              <select
                name="profile_visibility"
                value={formData.profile_visibility}
                onChange={handleInputChange}
                className="input-field"
              >
                <option value="public">Public</option>
                <option value="friends">Friends Only</option>
                <option value="private">Private</option>
              </select>
            </div>
            <div className="space-y-3">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  name="email_notifications"
                  checked={formData.email_notifications}
                  onChange={handleInputChange}
                  className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label className="ml-2 block text-sm text-gray-900">
                  Email notifications
                </label>
              </div>
              <div className="flex items-center">
                <input
                  type="checkbox"
                  name="push_notifications"
                  checked={formData.push_notifications}
                  onChange={handleInputChange}
                  className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label className="ml-2 block text-sm text-gray-900">
                  Push notifications
                </label>
              </div>
            </div>
          </div>
        </div>

        {/* Submit Button */}
        <div className="flex justify-end space-x-4">
          <Link
            to="/profile"
            className="btn-secondary"
          >
            Cancel
          </Link>
          <button
            type="submit"
            disabled={saving}
            className="btn-primary"
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>
    </div>
  );
};
