-- Extend users table with profile fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS location VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth DATE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS website VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS languages TEXT[]; -- Array of language codes
ALTER TABLE users ADD COLUMN IF NOT EXISTS interests TEXT[]; -- Array of interest tags
ALTER TABLE users ADD COLUMN IF NOT EXISTS travel_style VARCHAR(50); -- budget, mid-range, luxury, backpacker, etc.
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_visibility VARCHAR(20) DEFAULT 'public'; -- public, friends, private
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_notifications BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS push_notifications BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS verification_status VARCHAR(20) DEFAULT 'pending'; -- pending, verified, rejected
ALTER TABLE users ADD COLUMN IF NOT EXISTS verification_documents JSONB; -- Store verification document info
ALTER TABLE users ADD COLUMN IF NOT EXISTS rating_average DECIMAL(3,2) DEFAULT 0.00; -- Average rating from other users
ALTER TABLE users ADD COLUMN IF NOT EXISTS rating_count INTEGER DEFAULT 0; -- Number of ratings received
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_completion_percentage INTEGER DEFAULT 0; -- Profile completion status

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_location ON users(location);
CREATE INDEX IF NOT EXISTS idx_users_travel_style ON users(travel_style);
CREATE INDEX IF NOT EXISTS idx_users_verification_status ON users(verification_status);
CREATE INDEX IF NOT EXISTS idx_users_rating_average ON users(rating_average);
CREATE INDEX IF NOT EXISTS idx_users_interests ON users USING GIN(interests);
CREATE INDEX IF NOT EXISTS idx_users_languages ON users USING GIN(languages);