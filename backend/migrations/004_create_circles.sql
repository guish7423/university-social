CREATE TABLE IF NOT EXISTS circles (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(64) NOT NULL,
    description  TEXT DEFAULT '',
    icon         VARCHAR(256) DEFAULT '',
    cover        VARCHAR(256) DEFAULT '',
    creator_id   BIGINT NOT NULL REFERENCES users(id),
    member_count INTEGER DEFAULT 1,
    post_count   INTEGER DEFAULT 0,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_circles_creator ON circles(creator_id);

CREATE TABLE IF NOT EXISTS circle_members (
    id         BIGSERIAL PRIMARY KEY,
    circle_id  BIGINT NOT NULL REFERENCES circles(id),
    user_id    BIGINT NOT NULL REFERENCES users(id),
    role       VARCHAR(16) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(circle_id, user_id)
);

CREATE INDEX idx_circle_members_circle ON circle_members(circle_id);
CREATE INDEX idx_circle_members_user ON circle_members(user_id);

CREATE TABLE IF NOT EXISTS circle_posts (
    id         BIGSERIAL PRIMARY KEY,
    circle_id  BIGINT NOT NULL REFERENCES circles(id),
    user_id    BIGINT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    image_urls TEXT[] DEFAULT '{}',
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_circle_posts_circle ON circle_posts(circle_id);

CREATE TABLE IF NOT EXISTS circle_post_likes (
    id        BIGSERIAL PRIMARY KEY,
    post_id   BIGINT NOT NULL REFERENCES circle_posts(id),
    user_id   BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

CREATE TABLE IF NOT EXISTS circle_post_comments (
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT NOT NULL REFERENCES circle_posts(id),
    user_id    BIGINT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
