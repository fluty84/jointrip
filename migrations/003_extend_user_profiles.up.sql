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

-- Create user_ratings table for peer reviews
CREATE TABLE IF NOT EXISTS user_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rater_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rated_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    trip_id UUID, -- Optional reference to a specific trip
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one rating per user pair per trip (or one general rating if no trip)
    UNIQUE(rater_id, rated_id, trip_id)
);

-- Create indexes for user_ratings
CREATE INDEX IF NOT EXISTS idx_user_ratings_rater_id ON user_ratings(rater_id);
CREATE INDEX IF NOT EXISTS idx_user_ratings_rated_id ON user_ratings(rated_id);
CREATE INDEX IF NOT EXISTS idx_user_ratings_trip_id ON user_ratings(trip_id);
CREATE INDEX IF NOT EXISTS idx_user_ratings_rating ON user_ratings(rating);
CREATE INDEX IF NOT EXISTS idx_user_ratings_created_at ON user_ratings(created_at);

-- Create user_verification_requests table
CREATE TABLE IF NOT EXISTS user_verification_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- passport, driver_license, national_id, etc.
    document_number VARCHAR(100),
    document_url VARCHAR(500), -- URL to uploaded document image
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    reviewed_by UUID REFERENCES users(id), -- Admin who reviewed
    rejection_reason TEXT,
    
    -- Only one active verification request per user
    UNIQUE(user_id, status) DEFERRABLE INITIALLY DEFERRED
);

-- Create indexes for user_verification_requests
CREATE INDEX IF NOT EXISTS idx_user_verification_requests_user_id ON user_verification_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_user_verification_requests_status ON user_verification_requests(status);
CREATE INDEX IF NOT EXISTS idx_user_verification_requests_submitted_at ON user_verification_requests(submitted_at);

-- Create function to update user rating average
CREATE OR REPLACE FUNCTION update_user_rating_average()
RETURNS TRIGGER AS $$
BEGIN
    -- Update the rated user's average rating and count
    UPDATE users 
    SET 
        rating_average = (
            SELECT COALESCE(AVG(rating), 0) 
            FROM user_ratings 
            WHERE rated_id = COALESCE(NEW.rated_id, OLD.rated_id)
        ),
        rating_count = (
            SELECT COUNT(*) 
            FROM user_ratings 
            WHERE rated_id = COALESCE(NEW.rated_id, OLD.rated_id)
        )
    WHERE id = COALESCE(NEW.rated_id, OLD.rated_id);
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update rating averages
DROP TRIGGER IF EXISTS trigger_update_user_rating_average ON user_ratings;
CREATE TRIGGER trigger_update_user_rating_average
    AFTER INSERT OR UPDATE OR DELETE ON user_ratings
    FOR EACH ROW
    EXECUTE FUNCTION update_user_rating_average();

-- Create function to calculate profile completion percentage
CREATE OR REPLACE FUNCTION calculate_profile_completion(user_row users)
RETURNS INTEGER AS $$
DECLARE
    completion_score INTEGER := 0;
    total_fields INTEGER := 15; -- Total number of profile fields we're checking
BEGIN
    -- Basic required fields (already exist)
    IF user_row.email IS NOT NULL AND user_row.email != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.first_name IS NOT NULL AND user_row.first_name != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.last_name IS NOT NULL AND user_row.last_name != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.profile_photo_url IS NOT NULL AND user_row.profile_photo_url != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    -- Extended profile fields
    IF user_row.bio IS NOT NULL AND user_row.bio != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.location IS NOT NULL AND user_row.location != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.date_of_birth IS NOT NULL THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.phone IS NOT NULL AND user_row.phone != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.languages IS NOT NULL AND array_length(user_row.languages, 1) > 0 THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.interests IS NOT NULL AND array_length(user_row.interests, 1) > 0 THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.travel_style IS NOT NULL AND user_row.travel_style != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.verification_status = 'verified' THEN
        completion_score := completion_score + 2; -- Verification is worth 2 points
    END IF;
    
    IF user_row.rating_count > 0 THEN
        completion_score := completion_score + 1;
    END IF;
    
    IF user_row.website IS NOT NULL AND user_row.website != '' THEN
        completion_score := completion_score + 1;
    END IF;
    
    -- Calculate percentage
    RETURN ROUND((completion_score::DECIMAL / total_fields) * 100);
END;
$$ LANGUAGE plpgsql;

-- Create function to update profile completion percentage
CREATE OR REPLACE FUNCTION update_profile_completion()
RETURNS TRIGGER AS $$
BEGIN
    NEW.profile_completion_percentage := calculate_profile_completion(NEW);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update profile completion
DROP TRIGGER IF EXISTS trigger_update_profile_completion ON users;
CREATE TRIGGER trigger_update_profile_completion
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_profile_completion();
