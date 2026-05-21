CREATE TABLE IF NOT EXISTS verification_codes (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    email       VARCHAR(256) NOT NULL,
    code        VARCHAR(6) NOT NULL,
    school      VARCHAR(128) NOT NULL DEFAULT '',
    status      VARCHAR(16) NOT NULL DEFAULT 'pending',
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_verification_user ON verification_codes(user_id, status);
