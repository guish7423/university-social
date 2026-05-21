-- Migration 020: Additional performance indexes
-- Add missing FK indexes for frequently joined columns
-- Created: 2026-05-21

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_from_user ON notifications(from_user_id);
CREATE INDEX IF NOT EXISTS idx_reports_reporter ON reports(reporter_id);
CREATE INDEX IF NOT EXISTS idx_circle_post_comments_post ON circle_post_comments(post_id);
CREATE INDEX IF NOT EXISTS idx_circle_post_comments_user ON circle_post_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_comments_user ON activity_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_product_comments_user ON product_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_whisper_comments_user ON whisper_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
