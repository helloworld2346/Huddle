-- Migration 007: File System
-- Add file sharing capabilities

-- Files table for storing file metadata
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
    message_id INTEGER REFERENCES messages(id) ON DELETE CASCADE,
    
    -- File Information
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_extension VARCHAR(20),
    
    -- Storage Information
    bucket_name VARCHAR(100) NOT NULL DEFAULT 'huddle-files',
    object_key VARCHAR(500) NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    
    -- File Processing
    is_processed BOOLEAN DEFAULT FALSE,
    thumbnail_url VARCHAR(500),
    preview_url VARCHAR(500),
    
    -- Security & Access
    is_public BOOLEAN DEFAULT FALSE,
    access_token VARCHAR(255),
    expires_at TIMESTAMP,
    
    -- Metadata
    width INTEGER,  -- For images/videos
    height INTEGER, -- For images/videos
    duration INTEGER, -- For videos/audio (seconds)
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- File shares table for sharing files
CREATE TABLE file_shares (
    id SERIAL PRIMARY KEY,
    file_id INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    shared_by INTEGER NOT NULL REFERENCES users(id),
    shared_with INTEGER REFERENCES users(id),
    conversation_id INTEGER REFERENCES conversations(id),
    
    -- Share Settings
    can_download BOOLEAN DEFAULT TRUE,
    can_edit BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_conversation_id ON files(conversation_id);
CREATE INDEX idx_files_message_id ON files(message_id);
CREATE INDEX idx_files_created_at ON files(created_at);
CREATE INDEX idx_files_object_key ON files(object_key);

CREATE INDEX idx_file_shares_file_id ON file_shares(file_id);
CREATE INDEX idx_file_shares_shared_by ON file_shares(shared_by);
CREATE INDEX idx_file_shares_shared_with ON file_shares(shared_with);
CREATE INDEX idx_file_shares_conversation_id ON file_shares(conversation_id);

-- Add file_url, file_name, file_size columns to messages table if not exists
-- (These are already in the message model, but ensure they exist in DB)
ALTER TABLE messages ADD COLUMN IF NOT EXISTS file_url VARCHAR(500);
ALTER TABLE messages ADD COLUMN IF NOT EXISTS file_name VARCHAR(255);
ALTER TABLE messages ADD COLUMN IF NOT EXISTS file_size BIGINT;

-- Add indexes for message file fields
CREATE INDEX IF NOT EXISTS idx_messages_file_url ON messages(file_url);
CREATE INDEX IF NOT EXISTS idx_messages_file_name ON messages(file_name);
