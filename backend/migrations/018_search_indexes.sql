CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Posts content search
CREATE INDEX IF NOT EXISTS idx_posts_content_gin ON posts USING GIN (content gin_trgm_ops);

-- Circles name & description search
CREATE INDEX IF NOT EXISTS idx_circles_name_gin ON circles USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_circles_desc_gin ON circles USING GIN (description gin_trgm_ops);

-- Users nickname search
CREATE INDEX IF NOT EXISTS idx_users_nickname_gin ON users USING GIN (nickname gin_trgm_ops);

-- Courses search
CREATE INDEX IF NOT EXISTS idx_courses_name_gin ON courses USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_courses_teacher_gin ON courses USING GIN (teacher gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_courses_school_gin ON courses USING GIN (school gin_trgm_ops);

-- Products title & description search
CREATE INDEX IF NOT EXISTS idx_products_title_gin ON products USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_products_desc_gin ON products USING GIN (description gin_trgm_ops);

-- Lost items search
CREATE INDEX IF NOT EXISTS idx_lost_items_title_gin ON lost_items USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_lost_items_desc_gin ON lost_items USING GIN (description gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_lost_items_type_gin ON lost_items USING GIN (item_type gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_lost_items_location_gin ON lost_items USING GIN (location gin_trgm_ops);
