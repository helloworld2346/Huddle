-- Migration: 002_auth_schema.sql
-- Description: Authentication schema for Huddle application

-- Add new fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_attempts INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP;

-- Create password resets table
CREATE TABLE IF NOT EXISTS password_resets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create user activities table
CREATE TABLE IF NOT EXISTS user_activities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL, -- login, logout, password_change, etc.
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token);
CREATE INDEX IF NOT EXISTS idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX IF NOT EXISTS idx_password_resets_expires_at ON password_resets(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_activities_user_id ON user_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activities_activity_type ON user_activities(activity_type);
CREATE INDEX IF NOT EXISTS idx_user_activities_created_at ON user_activities(created_at);
CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login);
CREATE INDEX IF NOT EXISTS idx_users_login_attempts ON users(login_attempts);
