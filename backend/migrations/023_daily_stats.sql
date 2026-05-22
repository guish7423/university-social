-- daily_stats: 每日统计聚合，用于管理后台DAU/MAU/留存图表

CREATE TABLE IF NOT EXISTS daily_stats (
    id       SERIAL PRIMARY KEY,
    date     DATE NOT NULL UNIQUE,
    dau      INT NOT NULL DEFAULT 0,
    new_users    INT NOT NULL DEFAULT 0,
    new_posts    INT NOT NULL DEFAULT 0,
    new_comments INT NOT NULL DEFAULT 0,
    new_likes    INT NOT NULL DEFAULT 0,
    new_circles  INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMP DEFAULT NOW()
);

-- 回填历史数据（从已有记录聚合）
INSERT INTO daily_stats (date, new_users, new_posts, new_circles)
SELECT
    d::date AS date,
    (SELECT COUNT(*) FROM users WHERE DATE(created_at) = d::date),
    (SELECT COUNT(*) FROM posts WHERE DATE(created_at) = d::date),
    (SELECT COUNT(*) FROM circles WHERE DATE(created_at) = d::date)
FROM generate_series(
    (SELECT DATE(MIN(created_at)) FROM users),
    CURRENT_DATE,
    '1 day'::interval
) AS d
ON CONFLICT (date) DO NOTHING;
