ALTER TABLE users ADD COLUMN IF NOT EXISTS points INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_date DATE;

CREATE TABLE IF NOT EXISTS points_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    action      VARCHAR(50) NOT NULL,
    points      INTEGER NOT NULL,
    description TEXT DEFAULT '',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_points_logs_user_id ON points_logs(user_id);
CREATE INDEX idx_points_logs_action ON points_logs(action);
CREATE INDEX idx_points_logs_created_at ON points_logs(created_at);
