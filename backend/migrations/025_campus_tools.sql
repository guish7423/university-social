-- campus_calendar: 校历
CREATE TABLE IF NOT EXISTS campus_calendar (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    event_date  DATE NOT NULL,
    event_type  VARCHAR(50) NOT NULL DEFAULT 'academic', -- academic, holiday, exam, activity
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_campus_calendar_date ON campus_calendar(event_date);

-- campus_directory: 黄页
CREATE TABLE IF NOT EXISTS campus_directory (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    department  VARCHAR(100) NOT NULL,
    position    VARCHAR(100),
    phone       VARCHAR(30),
    email       VARCHAR(100),
    office      VARCHAR(100),
    category    VARCHAR(50) NOT NULL DEFAULT 'teacher', -- teacher, office, service
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_campus_directory_dept ON campus_directory(department);

-- empty_classrooms: 空教室
CREATE TABLE IF NOT EXISTS empty_classrooms (
    id          SERIAL PRIMARY KEY,
    building    VARCHAR(100) NOT NULL,
    room        VARCHAR(50) NOT NULL,
    capacity    INT DEFAULT 0,
    weekday     INT NOT NULL,        -- 1-7
    start_time  TIME NOT NULL,       -- 节次开始
    end_time    TIME NOT NULL,       -- 节次结束
    semester    VARCHAR(50) NOT NULL DEFAULT '',
    created_at  TIMESTAMP DEFAULT NOW(),
    UNIQUE (building, room, weekday, start_time, semester)
);

CREATE INDEX IF NOT EXISTS idx_empty_classrooms_time ON empty_classrooms(building, weekday, start_time);

-- 插入示例数据
INSERT INTO campus_calendar (title, event_date, event_type) VALUES
('开学日', '2026-09-01', 'academic'),
('期中考试', '2026-10-26', 'exam'),
('期末考试', '2027-01-04', 'exam'),
('寒假开始', '2027-01-18', 'holiday'),
('运动会', '2026-11-02', 'activity')
ON CONFLICT DO NOTHING;
