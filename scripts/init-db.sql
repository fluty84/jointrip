-- Initialize JoinTrip database
-- This script runs when the PostgreSQL container starts for the first time

-- Create database if it doesn't exist (already created by POSTGRES_DB env var)
-- CREATE DATABASE IF NOT EXISTS jointrip;

-- Connect to the jointrip database
\c jointrip;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a function to generate random UUIDs (alternative to uuid-ossp)
CREATE OR REPLACE FUNCTION gen_random_uuid() RETURNS uuid AS $$
BEGIN
    RETURN uuid_generate_v4();
END;
$$ LANGUAGE plpgsql;

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON DATABASE jointrip TO postgres;

-- Create schemas if needed (optional, using public schema for now)
-- CREATE SCHEMA IF NOT EXISTS jointrip_schema;

-- Log successful initialization
DO $$
BEGIN
    RAISE NOTICE 'JoinTrip database initialized successfully';
END $$;
