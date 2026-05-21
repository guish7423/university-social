package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"


	"github.com/guish/university-social/internal/model"
)

type WhisperRepository struct {
	db *sql.DB
}

func NewWhisperRepository(db *sql.DB) *WhisperRepository {
	return &WhisperRepository{db: db}
}

func pickCodename(userID, whisperID int64) string {
	idx := int(userID*31 + whisperID*7)
	if idx < 0 {
		idx = -idx
	}
	return model.CodenamePool[idx%len(model.CodenamePool)]
}

func (r *WhisperRepository) Create(userID int64, req *model.CreateWhisperRequest) (int64, error) {
	imagesJSON, _ := json.Marshal(req.Images)
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO whispers (user_id, content, images, is_anonymous, whisper_type, target_user_id, codename, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, '', NOW()) RETURNING id`,
		userID, req.Content, imagesJSON, req.IsAnonymous, req.WhisperType, req.TargetUserID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create whisper: %w", err)
	}

	if req.IsAnonymous {
		codename := pickCodename(userID, id)
		_, _ = r.db.Exec(`UPDATE whispers SET codename = $1 WHERE id = $2`, codename, id)
	}

	return id, nil
}

func (r *WhisperRepository) List(offset, limit int, currentUserID int64) ([]*model.Whisper, error) {
	query := `SELECT w.id, w.user_id, w.content, w.images, w.is_anonymous,
		w.codename, w.whisper_type, w.target_user_id,
		w.like_count, w.comment_count, w.created_at,
		CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM whispers w
		LEFT JOIN whisper_likes l ON l.whisper_id = w.id AND l.user_id = $1
		ORDER BY w.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, currentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list whispers: %w", err)
	}
	defer rows.Close()

	var whispers []*model.Whisper
	for rows.Next() {
		var w model.Whisper
		var imagesJSON []byte
		err := rows.Scan(&w.ID, &w.UserID, &w.Content, &imagesJSON,
			&w.IsAnonymous, &w.Codename, &w.WhisperType, &w.TargetUserID,
			&w.LikeCount, &w.CommentCount, &w.CreatedAt, &w.IsLiked)
		if err != nil {
			return nil, fmt.Errorf("scan whisper: %w", err)
		}
		json.Unmarshal(imagesJSON, &w.Images)
		if w.Images == nil {
			w.Images = []string{}
		}
		whispers = append(whispers, &w)
	}
	return whispers, nil
}

func (r *WhisperRepository) FindByID(id int64, currentUserID int64) (*model.Whisper, error) {
	query := `SELECT w.id, w.user_id, w.content, w.images, w.is_anonymous,
		w.codename, w.whisper_type, w.target_user_id,
		w.like_count, w.comment_count, w.created_at,
		CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM whispers w
		LEFT JOIN whisper_likes l ON l.whisper_id = w.id AND l.user_id = $1
		WHERE w.id = $2`

	var w model.Whisper
	var imagesJSON []byte
	err := r.db.QueryRow(query, currentUserID, id).Scan(&w.ID, &w.UserID, &w.Content, &imagesJSON,
		&w.IsAnonymous, &w.Codename, &w.WhisperType, &w.TargetUserID,
		&w.LikeCount, &w.CommentCount, &w.CreatedAt, &w.IsLiked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("whisper not found")
		}
		return nil, fmt.Errorf("find whisper: %w", err)
	}
	json.Unmarshal(imagesJSON, &w.Images)
	if w.Images == nil {
		w.Images = []string{}
	}
	return &w, nil
}

func (r *WhisperRepository) Delete(id int64, userID int64) error {
	res, err := r.db.Exec(`DELETE FROM whispers WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete whisper: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("whisper not found or not owned")
	}
	return nil
}

func (r *WhisperRepository) ToggleLike(whisperID, userID int64) (bool, error) {
	var exists bool
	r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM whisper_likes WHERE whisper_id = $1 AND user_id = $2)`,
		whisperID, userID).Scan(&exists)

	if exists {
		_, err := r.db.Exec(`DELETE FROM whisper_likes WHERE whisper_id = $1 AND user_id = $2`, whisperID, userID)
		if err != nil {
			return false, fmt.Errorf("unlike whisper: %w", err)
		}
		_, _ = r.db.Exec(`UPDATE whispers SET like_count = like_count - 1 WHERE id = $1 AND like_count > 0`, whisperID)
		return false, nil
	}

	_, err := r.db.Exec(`INSERT INTO whisper_likes (whisper_id, user_id) VALUES ($1, $2)`, whisperID, userID)
	if err != nil {
		return false, fmt.Errorf("like whisper: %w", err)
	}
	_, _ = r.db.Exec(`UPDATE whispers SET like_count = like_count + 1 WHERE id = $1`, whisperID)
	return true, nil
}

func (r *WhisperRepository) CreateComment(whisperID, userID int64, content string) (int64, error) {
	var codename string
	r.db.QueryRow(`SELECT codename FROM whispers WHERE id = $1`, whisperID).Scan(&codename)

	if codename != "" {
		userCodename := pickCodename(userID, whisperID)
		var id int64
		err := r.db.QueryRow(
			`INSERT INTO whisper_comments (whisper_id, user_id, content, codename) VALUES ($1, $2, $3, $4) RETURNING id`,
			whisperID, userID, content, userCodename,
		).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("create whisper comment: %w", err)
		}
		_, _ = r.db.Exec(`UPDATE whispers SET comment_count = comment_count + 1 WHERE id = $1`, whisperID)
		return id, nil
	}

	var id int64
	err := r.db.QueryRow(
		`INSERT INTO whisper_comments (whisper_id, user_id, content) VALUES ($1, $2, $3) RETURNING id`,
		whisperID, userID, content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create whisper comment: %w", err)
	}
	_, _ = r.db.Exec(`UPDATE whispers SET comment_count = comment_count + 1 WHERE id = $1`, whisperID)
	return id, nil
}

func (r *WhisperRepository) ListComments(whisperID int, offset, limit int) ([]*model.WhisperComment, error) {
	query := `SELECT id, whisper_id, user_id, content, codename, created_at
		FROM whisper_comments WHERE whisper_id = $1
		ORDER BY created_at ASC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, whisperID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list whisper comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.WhisperComment
	for rows.Next() {
		var c model.WhisperComment
		if err := rows.Scan(&c.ID, &c.WhisperID, &c.UserID, &c.Content, &c.Codename, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan whisper comment: %w", err)
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

type WhisperSummary struct {
	Total int `json:"total"`
	Today int `json:"today"`
}

func (r *WhisperRepository) Summary() (*WhisperSummary, error) {
	var s WhisperSummary
	r.db.QueryRow(`SELECT COUNT(*) FROM whispers`).Scan(&s.Total)
	r.db.QueryRow(`SELECT COUNT(*) FROM whispers WHERE created_at >= CURRENT_DATE`).Scan(&s.Today)
	return &s, nil
}
