CREATE TABLE IF NOT EXISTS universities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) DEFAULT '',
    logo_url VARCHAR(512) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS university_members (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    university_id BIGINT NOT NULL REFERENCES universities(id),
    is_admin BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, university_id)
);

CREATE INDEX IF NOT EXISTS idx_university_members_user ON university_members(user_id);
CREATE INDEX IF NOT EXISTS idx_university_members_university ON university_members(university_id);
