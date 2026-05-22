CREATE TABLE IF NOT EXISTS circle_activities (
    id         SERIAL PRIMARY KEY,
    circle_id  INT NOT NULL REFERENCES circles(id) ON DELETE CASCADE,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action     VARCHAR(50) NOT NULL, -- 'join', 'create_post', 'approve_member'
    content    TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_circle_activities_circle_id ON circle_activities(circle_id, created_at DESC);
