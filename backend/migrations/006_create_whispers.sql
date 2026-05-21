CREATE TABLE IF NOT EXISTS whispers (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    content         TEXT NOT NULL,
    images          JSONB DEFAULT '[]',
    is_anonymous    BOOLEAN DEFAULT true,
    codename        VARCHAR(64) DEFAULT '',
    whisper_type    SMALLINT DEFAULT 0,
    target_user_id  BIGINT REFERENCES users(id),
    like_count      INTEGER DEFAULT 0,
    comment_count   INTEGER DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_whispers_created_at ON whispers(created_at DESC);
CREATE INDEX idx_whispers_user_id ON whispers(user_id);
CREATE INDEX idx_whispers_type ON whispers(whisper_type);

CREATE TABLE IF NOT EXISTS whisper_comments (
    id          BIGSERIAL PRIMARY KEY,
    whisper_id  BIGINT NOT NULL REFERENCES whispers(id) ON DELETE CASCADE,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    content     TEXT NOT NULL,
    codename    VARCHAR(64) DEFAULT '',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_whisper_comments_whisper_id ON whisper_comments(whisper_id);

CREATE TABLE IF NOT EXISTS whisper_likes (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    whisper_id  BIGINT NOT NULL REFERENCES whispers(id) ON DELETE CASCADE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, whisper_id)
);

CREATE INDEX idx_whisper_likes_whisper_id ON whisper_likes(whisper_id);
