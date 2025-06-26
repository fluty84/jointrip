-- Revert extended user profiles
DROP TRIGGER IF EXISTS trigger_update_profile_completion ON users;
DROP TRIGGER IF EXISTS trigger_update_user_rating_average ON user_ratings;

DROP FUNCTION IF EXISTS update_profile_completion();
DROP FUNCTION IF EXISTS calculate_profile_completion(users);
DROP FUNCTION IF EXISTS update_user_rating_average();

DROP TABLE IF EXISTS user_verification_requests;
DROP TABLE IF EXISTS user_ratings;

-- Remove indexes
DROP INDEX IF EXISTS idx_users_location;
DROP INDEX IF EXISTS idx_users_travel_style;
DROP INDEX IF EXISTS idx_users_verification_status;
DROP INDEX IF EXISTS idx_users_rating_average;
DROP INDEX IF EXISTS idx_users_interests;
DROP INDEX IF EXISTS idx_users_languages;

-- Remove columns
ALTER TABLE users DROP COLUMN IF EXISTS profile_completion_percentage;
ALTER TABLE users DROP COLUMN IF EXISTS rating_count;
ALTER TABLE users DROP COLUMN IF EXISTS rating_average;
ALTER TABLE users DROP COLUMN IF EXISTS verification_documents;
ALTER TABLE users DROP COLUMN IF EXISTS verification_status;
ALTER TABLE users DROP COLUMN IF EXISTS push_notifications;
ALTER TABLE users DROP COLUMN IF EXISTS email_notifications;
ALTER TABLE users DROP COLUMN IF EXISTS profile_visibility;
ALTER TABLE users DROP COLUMN IF EXISTS travel_style;
ALTER TABLE users DROP COLUMN IF EXISTS interests;
ALTER TABLE users DROP COLUMN IF EXISTS languages;
ALTER TABLE users DROP COLUMN IF EXISTS website;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users DROP COLUMN IF EXISTS date_of_birth;
ALTER TABLE users DROP COLUMN IF EXISTS location;
ALTER TABLE users DROP COLUMN IF EXISTS bio;
