CREATE TABLE IF NOT EXISTS lost_items (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    title           VARCHAR(200) NOT NULL,
    description     TEXT DEFAULT '',
    category        VARCHAR(20) NOT NULL DEFAULT 'lost',
    item_type       VARCHAR(100) DEFAULT '',
    location        VARCHAR(200) DEFAULT '',
    contact         VARCHAR(200) DEFAULT '',
    images          JSONB DEFAULT '[]',
    status          SMALLINT DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_lost_items_user_id ON lost_items(user_id);
CREATE INDEX idx_lost_items_category ON lost_items(category);
CREATE INDEX idx_lost_items_status ON lost_items(status);
CREATE INDEX idx_lost_items_created_at ON lost_items(created_at);
