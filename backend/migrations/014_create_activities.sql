CREATE TABLE IF NOT EXISTS activities (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    title           VARCHAR(200) NOT NULL,
    description     TEXT DEFAULT '',
    activity_type   VARCHAR(50) NOT NULL,
    location        VARCHAR(200) DEFAULT '',
    max_participants INTEGER DEFAULT 0,
    start_time      TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time        TIMESTAMP WITH TIME ZONE,
    images          JSONB DEFAULT '[]',
    status          SMALLINT DEFAULT 0,
    participant_count INTEGER DEFAULT 0,
    comment_count   INTEGER DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activities_start_time ON activities(start_time);
CREATE INDEX idx_activities_type ON activities(activity_type);
CREATE INDEX idx_activities_user_id ON activities(user_id);
CREATE INDEX idx_activities_status ON activities(status);

CREATE TABLE IF NOT EXISTS activity_participants (
    id            BIGSERIAL PRIMARY KEY,
    activity_id   BIGINT NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    user_id       BIGINT NOT NULL REFERENCES users(id),
    status        SMALLINT DEFAULT 0,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(activity_id, user_id)
);

CREATE INDEX idx_activity_participants_activity ON activity_participants(activity_id);

CREATE TABLE IF NOT EXISTS activity_comments (
    id            BIGSERIAL PRIMARY KEY,
    activity_id   BIGINT NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    user_id       BIGINT NOT NULL REFERENCES users(id),
    content       TEXT NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activity_comments_activity ON activity_comments(activity_id);
