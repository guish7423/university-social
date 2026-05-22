ALTER TABLE circles ADD COLUMN IF NOT EXISTS join_type VARCHAR(16) NOT NULL DEFAULT 'free';

COMMENT ON COLUMN circles.join_type IS 'free: 自由加入, approve: 需审批';

CREATE TABLE IF NOT EXISTS circle_join_requests (
    id         BIGSERIAL PRIMARY KEY,
    circle_id  BIGINT NOT NULL REFERENCES circles(id),
    user_id    BIGINT NOT NULL REFERENCES users(id),
    status     VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    handled_at TIMESTAMP WITH TIME ZONE,
    handler_id BIGINT REFERENCES users(id),
    UNIQUE(circle_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_join_requests_circle ON circle_join_requests(circle_id);
CREATE INDEX IF NOT EXISTS idx_join_requests_user ON circle_join_requests(user_id);
