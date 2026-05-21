CREATE TABLE IF NOT EXISTS courses (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    teacher     VARCHAR(100) NOT NULL DEFAULT '',
    school      VARCHAR(100) NOT NULL DEFAULT '',
    department  VARCHAR(100) NOT NULL DEFAULT '',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_ratings (
    id              BIGSERIAL PRIMARY KEY,
    course_id       BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    teaching_rating SMALLINT NOT NULL DEFAULT 5 CHECK (teaching_rating >= 1 AND teaching_rating <= 5),
    difficulty      SMALLINT NOT NULL DEFAULT 3 CHECK (difficulty >= 1 AND difficulty <= 5),
    grading         SMALLINT NOT NULL DEFAULT 3 CHECK (grading >= 1 AND grading <= 5),
    comment         TEXT DEFAULT '',
    is_anonymous    BOOLEAN DEFAULT true,
    helpful_count   INTEGER DEFAULT 0,
    status          INTEGER DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(course_id, user_id)
);

CREATE INDEX idx_courses_name ON courses(name);
CREATE INDEX idx_courses_teacher ON courses(teacher);
CREATE INDEX idx_courses_school ON courses(school);
CREATE INDEX idx_course_ratings_course_id ON course_ratings(course_id);
CREATE INDEX idx_course_ratings_created_at ON course_ratings(created_at);


CREATE TABLE IF NOT EXISTS course_rating_helpful (
    id       BIGSERIAL PRIMARY KEY,
    rating_id BIGINT NOT NULL REFERENCES course_ratings(id) ON DELETE CASCADE,
    user_id   BIGINT NOT NULL REFERENCES users(id),
    UNIQUE(rating_id, user_id)
);
CREATE INDEX idx_course_helpful_rating ON course_rating_helpful(rating_id);