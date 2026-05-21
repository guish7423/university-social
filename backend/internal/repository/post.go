package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(userID int64, req *model.CreatePostRequest) (int64, error) {
	var imagesJSON interface{}
	if len(req.Images) > 0 {
		j, _ := json.Marshal(req.Images)
		imagesJSON = j
	}
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO posts (user_id, content, images, topic_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
		userID, req.Content, imagesJSON, req.TopicID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err)
	}
	return id, nil
}

type PostRow struct {
	model.Post
	AuthorNickname string `db:"author_nickname"`
	AuthorAvatar   string `db:"author_avatar"`
}

func scanPost(scanner interface {
	Scan(dest ...interface{}) error
}, includeAuthor bool) (*model.Post, error) {
	var imagesBytes []byte
	if includeAuthor {
		var p PostRow
		err := scanner.Scan(&p.ID, &p.UserID, &p.Content, &imagesBytes,
			&p.TopicID, &p.School, &p.LikeCount, &p.CommentCount,
			&p.CreatedAt, &p.UpdatedAt,
			&p.AuthorNickname, &p.AuthorAvatar,
			&p.IsLiked)
		if err != nil {
			return nil, err
		}
		if len(imagesBytes) > 0 {
			json.Unmarshal(imagesBytes, &p.Images)
		}
		p.Author = &model.User{Nickname: p.AuthorNickname, Avatar: p.AuthorAvatar}
		return &p.Post, nil
	}
	var p model.Post
	err := scanner.Scan(&p.ID, &p.UserID, &p.Content, &imagesBytes,
		&p.TopicID, &p.School, &p.LikeCount, &p.CommentCount,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if len(imagesBytes) > 0 {
		json.Unmarshal(imagesBytes, &p.Images)
	}
	return &p, nil
}

func (r *PostRepository) FindByID(id int64, currentUserID int64) (*model.Post, error) {
	query := `SELECT p.id, p.user_id, p.content, p.images, p.topic_id, p.school,
		p.like_count, p.comment_count, p.created_at, p.updated_at,
		u.nickname, u.avatar,
		CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = $2
		WHERE p.id = $1`
	return scanPost(r.db.QueryRow(query, id, currentUserID), true)
}

func (r *PostRepository) List(school string, topicID *int64, offset, limit int, currentUserID int64) ([]*model.Post, error) {
	args := []interface{}{currentUserID, limit, offset}
	where := ""
	argIdx := 4

	if school != "" {
		where = fmt.Sprintf("AND p.school = $%d", argIdx)
		args = append(args, school)
		argIdx++
	}
	if topicID != nil {
		where += fmt.Sprintf(" AND p.topic_id = $%d", argIdx)
		args = append(args, *topicID)
	}

	query := fmt.Sprintf(
		`SELECT p.id, p.user_id, p.content, p.images, p.topic_id, p.school,
			p.like_count, p.comment_count, p.created_at, p.updated_at,
			u.nickname, u.avatar,
			CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = $1
		WHERE 1=1 %s
		ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`, where)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		p, err := scanPost(rows, true)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) Delete(id, userID int64) error {
	res, err := r.db.Exec("DELETE FROM posts WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("post not found or not owner")
	}
	return nil
}

func (r *PostRepository) CreateComment(postID, userID int64, content string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO comments (post_id, user_id, content, created_at)
		 VALUES ($1, $2, $3, NOW()) RETURNING id`,
		postID, userID, content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create comment: %w", err)
	}
	r.db.Exec("UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1", postID)

	var authorID int64
	r.db.QueryRow("SELECT user_id FROM posts WHERE id = $1", postID).Scan(&authorID)
	if authorID != userID {
		r.CreateNotification(authorID, userID, "comment", content, postID)
	}

	return id, nil
}

func (r *PostRepository) ListComments(postID, offset, limit int) ([]*model.Comment, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname, u.avatar
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		var author model.User
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &author.Nickname, &author.Avatar); err != nil {
			return nil, err
		}
		c.Author = &author
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *PostRepository) ToggleLike(userID, postID int64) (bool, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO likes (user_id, post_id, created_at) VALUES ($1, $2, NOW())
		 ON CONFLICT (user_id, post_id) DO NOTHING RETURNING id`,
		userID, postID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		r.db.Exec("DELETE FROM likes WHERE user_id=$1 AND post_id=$2", userID, postID)
		r.db.Exec("UPDATE posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1", postID)
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("toggle like: %w", err)
	}
	r.db.Exec("UPDATE posts SET like_count = like_count + 1 WHERE id = $1", postID)

	var authorID int64
	r.db.QueryRow("SELECT user_id FROM posts WHERE id = $1", postID).Scan(&authorID)
	if authorID != userID {
		r.CreateNotification(authorID, userID, "like", "赞了你的动态", postID)
	}

	return true, nil
}

func (r *PostRepository) ListTopics() ([]*model.Topic, error) {
	rows, err := r.db.Query(
		"SELECT id, name, icon, post_count, created_at FROM topics ORDER BY post_count DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}
	defer rows.Close()

	var topics []*model.Topic
	for rows.Next() {
		var t model.Topic
		if err := rows.Scan(&t.ID, &t.Name, &t.Icon, &t.PostCount, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, &t)
	}
	return topics, nil
}

func (r *PostRepository) CreateNotification(userID, fromUserID int64, ntype, content string, refID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO notifications (user_id, from_user_id, type, content, ref_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW())`,
		userID, fromUserID, ntype, content, refID)
	return err
}

func (r *PostRepository) SearchPosts(query string, offset, limit int, currentUserID int64) ([]*model.Post, error) {
	args := []interface{}{currentUserID, limit, offset, "%" + query + "%", query}
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.images, p.topic_id, p.school,
			p.like_count, p.comment_count, p.created_at, p.updated_at,
			u.nickname, u.avatar,
			CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = $1
		WHERE p.content ILIKE $4
		ORDER BY similarity(p.content, $5) DESC, p.like_count DESC, p.created_at DESC LIMIT $2 OFFSET $3`, args...)
	if err != nil {
		return nil, fmt.Errorf("search posts: %w", err)
	}
	defer rows.Close()
	var posts []*model.Post
	for rows.Next() {
		p, err := scanPost(rows, true)
		if err != nil { return nil, err }
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) TrendingPosts(offset, limit int, currentUserID int64) ([]*model.Post, error) {
	args := []interface{}{currentUserID, limit, offset}
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.images, p.topic_id, p.school,
			p.like_count, p.comment_count, p.created_at, p.updated_at,
			u.nickname, u.avatar,
			CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = $1
		WHERE p.created_at >= CURRENT_DATE - INTERVAL '7 days'
		ORDER BY (p.like_count * 3 + p.comment_count * 2) DESC, p.created_at DESC
		LIMIT $2 OFFSET $3`, args...)
	if err != nil {
		return nil, fmt.Errorf("trending posts: %w", err)
	}
	defer rows.Close()
	var posts []*model.Post
	for rows.Next() {
		p, err := scanPost(rows, true)
		if err != nil { return nil, err }
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) FollowingPosts(userID int64, offset, limit int) ([]*model.Post, error) {
	args := []interface{}{userID, limit, offset}
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.images, p.topic_id, p.school,
			p.like_count, p.comment_count, p.created_at, p.updated_at,
			u.nickname, u.avatar,
			CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN follows f ON f.following_id = p.user_id
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = $1
		WHERE f.follower_id = $1
		ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`, args...)
	if err != nil {
		return nil, fmt.Errorf("following posts: %w", err)
	}
	defer rows.Close()
	var posts []*model.Post
	for rows.Next() {
		p, err := scanPost(rows, true)
		if err != nil { return nil, err }
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) ListUserPosts(userID, currentUserID int64, offset, limit int) ([]*model.Post, error) {
	args := []interface{}{userID, currentUserID, limit, offset}
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.images, p.topic_id, p.school,
			p.like_count, p.comment_count, p.created_at, p.updated_at,
			u.nickname, u.avatar,
			CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = $2
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC LIMIT $3 OFFSET $4`, args...)
	if err != nil {
		return nil, fmt.Errorf("list user posts: %w", err)
	}
	defer rows.Close()
	var posts []*model.Post
	for rows.Next() {
		p, err := scanPost(rows, true)
		if err != nil { return nil, err }
		posts = append(posts, p)
	}
	return posts, nil
}
