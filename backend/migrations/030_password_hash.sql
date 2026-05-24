-- Add password_hash column for account security feature
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(256) DEFAULT '';
