package repository

import (
	"database/sql"

	"github.com/guish/university-social/internal/model"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(userID, postID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO user_favorites (user_id, post_id, created_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (user_id, post_id) DO NOTHING`,
		userID, postID,
	)
	return err
}

func (r *FavoriteRepository) Remove(userID, postID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM user_favorites WHERE user_id=$1 AND post_id=$2`,
		userID, postID,
	)
	return err
}

func (r *FavoriteRepository) IsFavorited(userID, postID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM user_favorites WHERE user_id=$1 AND post_id=$2`,
		userID, postID,
	).Scan(&count)
	return count > 0, err
}

func (r *FavoriteRepository) ListByUser(userID int64, page, limit int) ([]*model.Post, int, error) {
	var total int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM user_favorites uf
		 JOIN posts p ON p.id = uf.post_id
		 WHERE uf.user_id = $1 AND p.deleted_at IS NULL`,
		userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.images, p.topic_id,
		        p.school, p.like_count, p.comment_count, p.is_featured,
		        p.created_at, p.updated_at
		 FROM user_favorites uf
		 JOIN posts p ON p.id = uf.post_id
		 WHERE uf.user_id = $1 AND p.deleted_at IS NULL
		 ORDER BY uf.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Images, &p.TopicID,
			&p.School, &p.LikeCount, &p.CommentCount, &p.IsFeatured, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		posts = append(posts, &p)
	}
	return posts, total, nil
}

func (r *FavoriteRepository) FavoriteIDs(userID int64) ([]int64, error) {
	rows, err := r.db.Query(
		`SELECT post_id FROM user_favorites WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
