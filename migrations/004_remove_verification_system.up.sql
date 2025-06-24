-- Remove verification system from MVP
-- Drop verification-related tables
DROP TABLE IF EXISTS user_verification_requests CASCADE;

-- Drop verification-related indexes
DROP INDEX IF EXISTS idx_users_verification_status;

-- Remove verification columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS verification_status;
ALTER TABLE users DROP COLUMN IF EXISTS verification_documents;

-- Update profile completion function to remove verification logic
CREATE OR REPLACE FUNCTION calculate_profile_completion(user_row users)
RETURNS INTEGER AS $$
DECLARE
    completion_score INTEGER := 0;
    total_fields INTEGER := 13; -- Total number of profile fields we're checking (reduced from 15)
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
