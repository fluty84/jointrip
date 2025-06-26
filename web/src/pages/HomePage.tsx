import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const HomePage: React.FC = () => {
  const { isAuthenticated, user } = useAuth();

  return (
    <div className="max-w-7xl mx-auto">
      {/* Hero Section */}
      <div className="relative overflow-hidden">
        <div className="absolute inset-0 hero-gradient opacity-5"></div>
        <div className="relative text-center py-24 px-4">
          <div className="animate-fade-in">
            <h1 className="text-6xl md:text-7xl font-bold mb-8 leading-tight">
              Find Your Perfect
              <br />
              <span className="hero-text">Travel Companion</span>
            </h1>
            <p className="text-xl md:text-2xl text-dark-700 mb-12 max-w-4xl mx-auto leading-relaxed">
              Connect with like-minded travelers, share amazing experiences, and split costs
              seamlessly. Your next adventure is just a click away.
            </p>
          </div>

          <div className="animate-slide-up">
            {isAuthenticated ? (
              <div className="space-y-6">
                <div className="bg-white/80 backdrop-blur-sm rounded-2xl p-6 shadow-soft max-w-md mx-auto">
                  <p className="text-xl text-dark-700 mb-4">
                    Welcome back, <span className="font-semibold text-primary-600">{user?.firstName}</span>!
                  </p>
                  <p className="text-lg text-dark-600">
                    Ready for your next adventure?
                  </p>
                </div>
                <div className="flex flex-col sm:flex-row justify-center gap-4">
                  <Link to="/trips" className="btn-primary text-lg px-10 py-4 shadow-medium">
                    🌍 Browse Trips
                  </Link>
                  <Link to="/profile" className="btn-secondary text-lg px-10 py-4">
                    👤 View Profile
                  </Link>
                </div>
              </div>
            ) : (
              <div className="flex flex-col sm:flex-row justify-center gap-4">
                <Link to="/login" className="btn-primary text-lg px-12 py-4 shadow-medium">
                  🚀 Get Started
                </Link>
                <Link to="/about" className="btn-secondary text-lg px-12 py-4">
                  📖 Learn More
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="py-24 px-4">
        <div className="text-center mb-20">
          <h2 className="text-4xl md:text-5xl font-bold text-dark-900 mb-6">
            Why Choose JoinTrip?
          </h2>
          <p className="text-lg text-gray-600">
            Everything you need for amazing group travel experiences
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          <div className="text-center p-6">
            <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Find Travel Companions</h3>
            <p className="text-gray-600">
              Connect with verified travelers who share your interests and destinations.
            </p>
          </div>

          <div className="text-center p-6">
            <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Split Expenses</h3>
            <p className="text-gray-600">
              Easily track and split travel costs with built-in expense management tools.
            </p>
          </div>

          <div className="text-center p-6">
            <div className="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Safe & Secure</h3>
            <p className="text-gray-600">
              Travel with confidence using our verified user system and secure platform.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
