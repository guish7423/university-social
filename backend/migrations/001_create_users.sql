CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    open_id     VARCHAR(64) NOT NULL UNIQUE,
    union_id    VARCHAR(64) DEFAULT '',
    nickname    VARCHAR(64) DEFAULT '',
    avatar      VARCHAR(512) DEFAULT '',
    phone       VARCHAR(20) DEFAULT '',
    school      VARCHAR(128) DEFAULT '',
    is_verified BOOLEAN DEFAULT FALSE,
    role        VARCHAR(16) DEFAULT 'user',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_open_id ON users(open_id);
