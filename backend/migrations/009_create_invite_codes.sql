CREATE TABLE IF NOT EXISTS invite_codes (
    id            BIGSERIAL PRIMARY KEY,
    code          VARCHAR(12) NOT NULL UNIQUE,
    creator_id    BIGINT NOT NULL REFERENCES users(id),
    invited_user_id BIGINT REFERENCES users(id),
    status        VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    used_at       TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_invite_codes_creator ON invite_codes(creator_id);
CREATE INDEX idx_invite_codes_code ON invite_codes(code);
