package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/guish/university-social/internal/model"
)

type ActivityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(userID int64, req *model.CreateActivityRequest) (int64, error) {
	imagesJSON, _ := json.Marshal(req.Images)
	var endTime *time.Time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil { endTime = &t }
	}
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO activities (user_id, title, description, activity_type, location, max_participants, start_time, end_time, images, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id`,
		userID, req.Title, req.Description, req.ActivityType, req.Location,
		req.MaxParticipants, startTime, endTime, imagesJSON,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create activity: %w", err)
	}
	return id, nil
}

func (r *ActivityRepository) scanActivity(scanner interface {
	Scan(dest ...interface{}) error
}, currentUserID int64) (*model.Activity, error) {
	var a model.Activity
	var desc, loc sql.NullString
	var endTime sql.NullTime
	var imagesJSON []byte
	var isParticipant bool

	err := scanner.Scan(&a.ID, &a.UserID, &a.Title, &desc, &a.ActivityType,
		&loc, &a.MaxParticipants, &a.StartTime, &endTime, &imagesJSON,
		&a.Status, &a.ParticipantCount, &a.CommentCount, &a.CreatedAt, &a.UpdatedAt, &isParticipant)
	if err != nil {
		return nil, err
	}
	a.Description = desc.String
	a.Location = loc.String
	if endTime.Valid { a.EndTime = &endTime.Time }
	if len(imagesJSON) > 0 { json.Unmarshal(imagesJSON, &a.Images) }
	if a.Images == nil { a.Images = []string{} }
	a.IsParticipant = isParticipant
	a.IsOwner = a.UserID == currentUserID
	return &a, nil
}

const activitySelectCols = `a.id, a.user_id, a.title, a.description, a.activity_type,
	a.location, a.max_participants, a.start_time, a.end_time, a.images,
	a.status, a.participant_count, a.comment_count, a.created_at, a.updated_at,
	CASE WHEN ap.id IS NOT NULL THEN true ELSE false END as is_participant`

const activityFromJoins = `FROM activities a
	LEFT JOIN activity_participants ap ON ap.activity_id = a.id AND ap.user_id = $1`

func (r *ActivityRepository) List(activityType string, offset, limit int, currentUserID int64) ([]*model.Activity, error) {
	args := []interface{}{currentUserID}
	where := "a.status = 0"
	if activityType != "" && activityType != "全部" {
		where += fmt.Sprintf(" AND a.activity_type = $%d", len(args)+1)
		args = append(args, activityType)
	}

	query := fmt.Sprintf(`SELECT %s %s WHERE %s ORDER BY a.start_time ASC LIMIT $%d OFFSET $%d`,
		activitySelectCols, activityFromJoins, where, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	defer rows.Close()

	var activities []*model.Activity
	for rows.Next() {
		a, err := r.scanActivity(rows, currentUserID)
		if err != nil { return nil, err }
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *ActivityRepository) FindByID(id int64, currentUserID int64) (*model.Activity, error) {
	query := fmt.Sprintf(`SELECT %s %s WHERE a.id = $2`, activitySelectCols, activityFromJoins)
	row := r.db.QueryRow(query, currentUserID, id)
	a, err := r.scanActivity(row, currentUserID)
	if err != nil {
		if err == sql.ErrNoRows { return nil, fmt.Errorf("activity not found") }
		return nil, fmt.Errorf("find activity: %w", err)
	}
	return a, nil
}

func (r *ActivityRepository) Delete(id, userID int64) error {
	res, err := r.db.Exec(`DELETE FROM activities WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete activity: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("activity not found or not owned")
	}
	return nil
}

func (r *ActivityRepository) Join(activityID, userID int64) error {
	var maxP, count int
	r.db.QueryRow(`SELECT max_participants, participant_count FROM activities WHERE id = $1`, activityID).Scan(&maxP, &count)
	if maxP > 0 && count >= maxP {
		return fmt.Errorf("活动已满")
	}

	_, err := r.db.Exec(`INSERT INTO activity_participants (activity_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, activityID, userID)
	if err != nil { return fmt.Errorf("join activity: %w", err) }
	_, _ = r.db.Exec(`UPDATE activities SET participant_count = (SELECT COUNT(*) FROM activity_participants WHERE activity_id = $1) WHERE id = $1`, activityID)
	return nil
}

func (r *ActivityRepository) Leave(activityID, userID int64) error {
	_, err := r.db.Exec(`DELETE FROM activity_participants WHERE activity_id = $1 AND user_id = $2`, activityID, userID)
	if err != nil { return fmt.Errorf("leave activity: %w", err) }
	_, _ = r.db.Exec(`UPDATE activities SET participant_count = (SELECT COUNT(*) FROM activity_participants WHERE activity_id = $1) WHERE id = $1`, activityID)
	return nil
}

func (r *ActivityRepository) ListParticipants(activityID int64) ([]*model.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.nickname, u.avatar FROM activity_participants ap
		 JOIN users u ON u.id = ap.user_id WHERE ap.activity_id = $1`, activityID)
	if err != nil { return nil, err }
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Avatar); err != nil { return nil, err }
		users = append(users, &u)
	}
	return users, nil
}

func (r *ActivityRepository) CreateComment(userID, activityID int64, content string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO activity_comments (activity_id, user_id, content) VALUES ($1, $2, $3) RETURNING id`,
		activityID, userID, content,
	).Scan(&id)
	if err != nil { return 0, fmt.Errorf("create comment: %w", err) }
	_, _ = r.db.Exec(`UPDATE activities SET comment_count = comment_count + 1 WHERE id = $1`, activityID)
	return id, nil
}

func (r *ActivityRepository) ListComments(activityID int64) ([]*model.ActivityComment, error) {
	rows, err := r.db.Query(
		`SELECT ac.id, ac.activity_id, ac.user_id, ac.content, ac.created_at
		 FROM activity_comments ac WHERE ac.activity_id = $1 ORDER BY ac.created_at ASC`, activityID)
	if err != nil { return nil, err }
	defer rows.Close()

	var comments []*model.ActivityComment
	for rows.Next() {
		var c model.ActivityComment
		if err := rows.Scan(&c.ID, &c.ActivityID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *ActivityRepository) ListByUser(userID, currentUserID int64, offset, limit int) ([]*model.Activity, error) {
	query := fmt.Sprintf(`SELECT %s %s WHERE a.user_id = $2 ORDER BY a.start_time DESC LIMIT $3 OFFSET $4`,
		activitySelectCols, activityFromJoins)
	rows, err := r.db.Query(query, currentUserID, userID, limit, offset)
	if err != nil { return nil, err }
	defer rows.Close()

	var activities []*model.Activity
	for rows.Next() {
		a, err := r.scanActivity(rows, currentUserID)
		if err != nil { return nil, err }
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *ActivityRepository) ListJoined(userID int64, offset, limit int) ([]*model.Activity, error) {
	query := fmt.Sprintf(`SELECT %s %s WHERE ap.user_id = $2 ORDER BY a.start_time DESC LIMIT $3 OFFSET $4`,
		activitySelectCols, activityFromJoins)
	rows, err := r.db.Query(query, userID, userID, limit, offset)
	if err != nil { return nil, err }
	defer rows.Close()

	var activities []*model.Activity
	for rows.Next() {
		a, err := r.scanActivity(rows, userID)
		if err != nil { return nil, err }
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *ActivityRepository) Types() []string {
	return []string{"全部", "竞赛组队", "约球", "学习局", "桌游局", "其他"}
}
