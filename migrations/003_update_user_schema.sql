-- Migration: 003_update_user_schema.sql
-- Description: Update user table to match the User model

-- Drop existing indexes first
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- Rename columns to match model
ALTER TABLE users RENAME COLUMN password_hash TO password;
ALTER TABLE users RENAME COLUMN full_name TO display_name;
ALTER TABLE users RENAME COLUMN avatar_url TO avatar;

-- Add new columns
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_attempts INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Update column sizes to match model
ALTER TABLE users ALTER COLUMN username TYPE VARCHAR(20);
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(254);
ALTER TABLE users ALTER COLUMN display_name TYPE VARCHAR(50);
ALTER TABLE users ALTER COLUMN avatar TYPE VARCHAR(255);

-- Remove old columns that are not in model
ALTER TABLE users DROP COLUMN IF EXISTS status;
ALTER TABLE users DROP COLUMN IF EXISTS last_seen;

-- Recreate indexes with correct names
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login);
CREATE INDEX IF NOT EXISTS idx_users_login_attempts ON users(login_attempts);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
