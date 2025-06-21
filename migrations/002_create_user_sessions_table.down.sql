-- Drop trigger
DROP TRIGGER IF EXISTS update_user_sessions_last_used_at ON user_sessions;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_sessions_active_expires;
DROP INDEX IF EXISTS idx_user_sessions_user_active;
DROP INDEX IF EXISTS idx_user_sessions_last_used_at;
DROP INDEX IF EXISTS idx_user_sessions_created_at;
DROP INDEX IF EXISTS idx_user_sessions_is_active;
DROP INDEX IF EXISTS idx_user_sessions_expires_at;
DROP INDEX IF EXISTS idx_user_sessions_refresh_token;
DROP INDEX IF EXISTS idx_user_sessions_access_token;
DROP INDEX IF EXISTS idx_user_sessions_user_id;

-- Drop user_sessions table
DROP TABLE IF EXISTS user_sessions;
