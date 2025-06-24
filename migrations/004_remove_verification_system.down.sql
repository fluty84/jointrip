-- Revert removal of verification system
-- Add verification columns back to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS verification_status VARCHAR(20) DEFAULT 'pending';
ALTER TABLE users ADD COLUMN IF NOT EXISTS verification_documents JSONB;

-- Create verification index
CREATE INDEX IF NOT EXISTS idx_users_verification_status ON users(verification_status);

-- Recreate user_verification_requests table
CREATE TABLE IF NOT EXISTS user_verification_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(100),
    document_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'pending',
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    reviewed_by UUID REFERENCES users(id),
    rejection_reason TEXT,
    UNIQUE(user_id, status) DEFERRABLE INITIALLY DEFERRED
);

-- Create indexes for user_verification_requests
CREATE INDEX IF NOT EXISTS idx_user_verification_requests_user_id ON user_verification_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_user_verification_requests_status ON user_verification_requests(status);
CREATE INDEX IF NOT EXISTS idx_user_verification_requests_submitted_at ON user_verification_requests(submitted_at);
