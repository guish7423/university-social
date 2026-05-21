CREATE TABLE IF NOT EXISTS reports (
    id              BIGSERIAL PRIMARY KEY,
    reporter_id     BIGINT NOT NULL REFERENCES users(id),
    content_type    VARCHAR(32) NOT NULL,  -- 'post', 'comment', 'whisper', 'whisper_comment'
    content_id      BIGINT NOT NULL,
    reason          VARCHAR(256) NOT NULL,
    detail          TEXT DEFAULT '',
    status          VARCHAR(16) DEFAULT 'pending',  -- 'pending', 'resolved', 'dismissed'
    resolved_by     BIGINT REFERENCES users(id),
    resolved_at     TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_content ON reports(content_type, content_id);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at DESC);
