package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) SearchCourses(query string, offset, limit int) ([]*model.Course, error) {
	like := "%" + query + "%"
	rows, err := r.db.Query(
		`SELECT id, name, teacher, school, department, created_at FROM courses
		 WHERE name ILIKE $1 OR teacher ILIKE $1 OR school ILIKE $1
		 ORDER BY name LIMIT $2 OFFSET $3`, like, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search courses: %w", err)
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Teacher, &c.School, &c.Department, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}
	return courses, nil
}

func (r *CourseRepository) FindByID(id int64) (*model.Course, error) {
	var c model.Course
	err := r.db.QueryRow(
		`SELECT id, name, teacher, school, department, created_at FROM courses WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Teacher, &c.School, &c.Department, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course not found")
		}
		return nil, fmt.Errorf("find course: %w", err)
	}
	return &c, nil
}

func (r *CourseRepository) CreateRating(userID, courseID int64, req *model.CreateRatingRequest) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO course_ratings (course_id, user_id, teaching_rating, difficulty, grading, comment, is_anonymous)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (course_id, user_id) DO UPDATE SET
		   teaching_rating = EXCLUDED.teaching_rating, difficulty = EXCLUDED.difficulty,
		   grading = EXCLUDED.grading, comment = EXCLUDED.comment,
		   is_anonymous = EXCLUDED.is_anonymous, updated_at = NOW()
		 RETURNING id`,
		courseID, userID, req.TeachingRating, req.Difficulty, req.Grading, req.Comment, req.IsAnonymous,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create rating: %w", err)
	}
	return id, nil
}

func (r *CourseRepository) GetUserRating(userID, courseID int64) (*model.CourseRating, error) {
	var cr model.CourseRating
	err := r.db.QueryRow(
		`SELECT id, course_id, user_id, teaching_rating, difficulty, grading,
		        COALESCE(comment,''), is_anonymous, helpful_count, status, created_at, updated_at
		 FROM course_ratings WHERE course_id = $1 AND user_id = $2`,
		courseID, userID,
	).Scan(&cr.ID, &cr.CourseID, &cr.UserID, &cr.TeachingRating, &cr.Difficulty,
		&cr.Grading, &cr.Comment, &cr.IsAnonymous, &cr.HelpfulCount, &cr.Status,
		&cr.CreatedAt, &cr.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cr, nil
}

func (r *CourseRepository) ListRatings(courseID int64, order string, offset, limit int) ([]*model.CourseRating, error) {
	if order != "helpful" {
		order = "newest"
	}
	orderClause := "cr.created_at DESC"
	if order == "helpful" {
		orderClause = "cr.helpful_count DESC, cr.created_at DESC"
	}

	rows, err := r.db.Query(
		fmt.Sprintf(`SELECT cr.id, cr.course_id, cr.user_id, cr.teaching_rating, cr.difficulty, cr.grading,
		       COALESCE(cr.comment,''), cr.is_anonymous, cr.helpful_count, cr.status, cr.created_at, cr.updated_at,
		       CASE WHEN cr.is_anonymous THEN '' ELSE u.nickname END,
		       CASE WHEN cr.is_anonymous THEN '' ELSE u.avatar END
		 FROM course_ratings cr
		 JOIN users u ON u.id = cr.user_id
		 WHERE cr.course_id = $1 AND cr.status = 0
		 ORDER BY %s LIMIT $2 OFFSET $3`, orderClause),
		courseID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list ratings: %w", err)
	}
	defer rows.Close()

	var ratings []*model.CourseRating
	for rows.Next() {
		var cr model.CourseRating
		if err := rows.Scan(&cr.ID, &cr.CourseID, &cr.UserID, &cr.TeachingRating, &cr.Difficulty,
			&cr.Grading, &cr.Comment, &cr.IsAnonymous, &cr.HelpfulCount, &cr.Status,
			&cr.CreatedAt, &cr.UpdatedAt, &cr.Nickname, &cr.Avatar); err != nil {
			return nil, err
		}
		ratings = append(ratings, &cr)
	}
	return ratings, nil
}

func (r *CourseRepository) MarkHelpful(ratingID, userID int64) error {
	var exists bool
	r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM course_rating_helpful WHERE rating_id=$1 AND user_id=$2)`,
		ratingID, userID).Scan(&exists)
	if exists {
		return fmt.Errorf("already marked")
	}
	_, err := r.db.Exec(`INSERT INTO course_rating_helpful (rating_id, user_id) VALUES ($1, $2)`, ratingID, userID)
	if err != nil {
		return fmt.Errorf("mark helpful: %w", err)
	}
	_, _ = r.db.Exec(`UPDATE course_ratings SET helpful_count = helpful_count + 1 WHERE id = $1`, ratingID)
	return nil
}

func (r *CourseRepository) GetCourseDetail(courseID, userID int64) (*model.CourseDetail, error) {
	course, err := r.FindByID(courseID)
	if err != nil {
		return nil, err
	}

	detail := &model.CourseDetail{Course: *course}
	r.db.QueryRow(
		`SELECT COALESCE(AVG(teaching_rating),0), COALESCE(AVG(difficulty),0), COALESCE(AVG(grading),0), COUNT(*)
		 FROM course_ratings WHERE course_id = $1 AND status = 0`,
		courseID,
	).Scan(&detail.AvgTeaching, &detail.AvgDifficulty, &detail.AvgGrading, &detail.RatingCount)

	if userID > 0 {
		detail.UserRating, _ = r.GetUserRating(userID, courseID)
	}

	return detail, nil
}
