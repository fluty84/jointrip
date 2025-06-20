-- Create user_sessions table
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    google_access_token TEXT,
    google_refresh_token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for user_sessions table
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_access_token ON user_sessions(access_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_active ON user_sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_user_sessions_created_at ON user_sessions(created_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_last_used_at ON user_sessions(last_used_at);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_active ON user_sessions(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_user_sessions_active_expires ON user_sessions(is_active, expires_at);

-- Create trigger to update last_used_at timestamp
CREATE TRIGGER update_user_sessions_last_used_at 
    BEFORE UPDATE ON user_sessions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
