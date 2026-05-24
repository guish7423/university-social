ALTER TABLE circle_posts ADD COLUMN is_pinned BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX idx_circle_posts_pinned ON circle_posts (circle_id, is_pinned DESC, created_at DESC);
