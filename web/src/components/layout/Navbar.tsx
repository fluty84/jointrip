import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

export const Navbar: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <nav className="bg-white/95 backdrop-blur-md shadow-soft border-b border-primary-100 sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-18">
          {/* Logo */}
          <Link to="/" className="flex items-center space-x-3 group">
            <div className="relative">
              <img
                src="/img/logo_s.png"
                alt="JoinTrip Logo"
                className="h-10 w-auto transition-transform duration-300 group-hover:scale-110"
              />
              <div className="absolute inset-0 bg-gradient-to-r from-primary-500 to-accent-500 rounded-lg opacity-0 group-hover:opacity-20 transition-opacity duration-300"></div>
            </div>
            <span className="text-2xl font-bold hero-text">JoinTrip</span>
          </Link>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-8">
            <Link
              to="/"
              className="nav-link relative group"
            >
              Home
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-primary-500 to-accent-500 transition-all duration-300 group-hover:w-full"></span>
            </Link>
            <Link
              to="/trips"
              className="nav-link relative group"
            >
              Trips
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-primary-500 to-accent-500 transition-all duration-300 group-hover:w-full"></span>
            </Link>
            <Link
              to="/about"
              className="nav-link relative group"
            >
              About
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-primary-500 to-accent-500 transition-all duration-300 group-hover:w-full"></span>
            </Link>
          </div>

          {/* User Menu */}
          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <div className="flex items-center space-x-4">
                <Link
                  to="/profile"
                  className="flex items-center space-x-3 nav-link group"
                >
                  {user?.picture ? (
                    <div className="relative">
                      <img
                        src={user.picture}
                        alt={`${user.firstName} ${user.lastName}`}
                        className="w-10 h-10 rounded-full border-2 border-primary-200 transition-all duration-300 group-hover:border-primary-400 group-hover:scale-105"
                      />
                      <div className="absolute inset-0 rounded-full bg-gradient-to-r from-primary-500 to-accent-500 opacity-0 group-hover:opacity-20 transition-opacity duration-300"></div>
                    </div>
                  ) : (
                    <div className="w-10 h-10 bg-gradient-to-r from-primary-100 to-accent-100 rounded-full flex items-center justify-center border-2 border-primary-200 transition-all duration-300 group-hover:border-primary-400 group-hover:scale-105">
                      <span className="text-primary-700 text-sm font-semibold">
                        {user?.firstName?.[0]}{user?.lastName?.[0]}
                      </span>
                    </div>
                  )}
                  <span className="hidden md:block font-medium">{user?.firstName}</span>
                </Link>
                <button
                  onClick={handleLogout}
                  className="btn-secondary text-sm px-4 py-2"
                >
                  Logout
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-4">
                <Link
                  to="/login"
                  className="nav-link"
                >
                  Login
                </Link>
                <Link
                  to="/login"
                  className="btn-primary text-sm px-6 py-2"
                >
                  Get Started
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};
