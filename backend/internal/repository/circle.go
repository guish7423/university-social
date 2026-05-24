package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/guish/university-social/internal/model"
)

type CircleRepository struct {
	db *sql.DB
}

func NewCircleRepository(db *sql.DB) *CircleRepository {
	return &CircleRepository{db: db}
}

func (r *CircleRepository) List(offset, limit int, userID int64) ([]*model.Circle, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.name, c.description, c.icon, c.cover, c.creator_id,
		        c.member_count, c.post_count, c.join_type, c.created_at,
		        CASE WHEN cm.id IS NOT NULL THEN true ELSE false END,
		        EXISTS(SELECT 1 FROM circle_join_requests WHERE circle_id = c.id AND user_id = $3 AND status = 'pending')
		 FROM circles c
		 LEFT JOIN circle_members cm ON cm.circle_id = c.id AND cm.user_id = $3
		 ORDER BY c.member_count DESC LIMIT $1 OFFSET $2`,
		limit, offset, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list circles: %w", err)
	}
	defer rows.Close()

	var circles []*model.Circle
	for rows.Next() {
		var c model.Circle
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Icon, &c.Cover,
			&c.CreatorID, &c.MemberCount, &c.PostCount, &c.JoinType, &c.CreatedAt, &c.IsMember, &c.HasPendingRequest); err != nil {
			return nil, err
		}
		circles = append(circles, &c)
	}
	return circles, nil
}

func (r *CircleRepository) GetByID(id, userID int64) (*model.Circle, error) {
	var c model.Circle
	err := r.db.QueryRow(
		`SELECT c.id, c.name, c.description, c.icon, c.cover, c.creator_id,
		        c.member_count, c.post_count, c.join_type, c.created_at,
		        CASE WHEN cm.id IS NOT NULL THEN true ELSE false END,
		        EXISTS(SELECT 1 FROM circle_join_requests WHERE circle_id = c.id AND user_id = $2 AND status = 'pending')
		 FROM circles c
		 LEFT JOIN circle_members cm ON cm.circle_id = c.id AND cm.user_id = $2
		 WHERE c.id = $1`, id, userID,
	).Scan(&c.ID, &c.Name, &c.Description, &c.Icon, &c.Cover,
		&c.CreatorID, &c.MemberCount, &c.PostCount, &c.JoinType, &c.CreatedAt, &c.IsMember, &c.HasPendingRequest)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get circle: %w", err)
	}
	return &c, nil
}

func (r *CircleRepository) Create(userID int64, req *model.CreateCircleRequest) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	joinType := req.JoinType
	if joinType == "" {
		joinType = "free"
	}

	var circleID int64
	err = tx.QueryRow(
		`INSERT INTO circles (name, description, icon, cover, creator_id, join_type)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		req.Name, req.Description, req.Icon, req.Cover, userID, joinType,
	).Scan(&circleID)
	if err != nil {
		return 0, fmt.Errorf("create circle: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO circle_members (circle_id, user_id, role) VALUES ($1, $2, 'creator')`,
		circleID, userID,
	)
	if err != nil {
		return 0, fmt.Errorf("add creator: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return circleID, nil
}

func (r *CircleRepository) Join(circleID, userID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO circle_members (circle_id, user_id, role)
		 VALUES ($1, $2, 'member') ON CONFLICT DO NOTHING`,
		circleID, userID,
	)
	if err != nil {
		return fmt.Errorf("join circle: %w", err)
	}
	r.db.Exec("UPDATE circles SET member_count = (SELECT COUNT(*) FROM circle_members WHERE circle_id = $1) WHERE id = $1", circleID)
	return nil
}

func (r *CircleRepository) Leave(circleID, userID int64) error {
	res, err := r.db.Exec(
		"DELETE FROM circle_members WHERE circle_id=$1 AND user_id=$2 AND role != 'creator'",
		circleID, userID,
	)
	if err != nil {
		return fmt.Errorf("leave circle: %w", err)
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		r.db.Exec("UPDATE circles SET member_count = (SELECT COUNT(*) FROM circle_members WHERE circle_id = $1) WHERE id = $1", circleID)
	}
	return nil
}

func (r *CircleRepository) GetJoinType(circleID int64) (string, error) {
	var joinType string
	err := r.db.QueryRow(`SELECT join_type FROM circles WHERE id = $1`, circleID).Scan(&joinType)
	if err != nil {
		return "", fmt.Errorf("get join type: %w", err)
	}
	return joinType, nil
}

func (r *CircleRepository) IsCreatorOrAdmin(circleID, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM circle_members WHERE circle_id = $1 AND user_id = $2 AND role IN ('creator', 'admin'))`,
		circleID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *CircleRepository) PinPost(postID, circleID, userID int64) error {
	isAdmin, err := r.IsCreatorOrAdmin(circleID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return fmt.Errorf("not authorized")
	}
	_, err = r.db.Exec("UPDATE circle_posts SET is_pinned = true WHERE id = $1 AND circle_id = $2", postID, circleID)
	return err
}

func (r *CircleRepository) UnpinPost(postID, circleID, userID int64) error {
	isAdmin, err := r.IsCreatorOrAdmin(circleID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return fmt.Errorf("not authorized")
	}
	_, err = r.db.Exec("UPDATE circle_posts SET is_pinned = false WHERE id = $1 AND circle_id = $2", postID, circleID)
	return err
}

func (r *CircleRepository) CreateJoinRequest(circleID, userID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO circle_join_requests (circle_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (circle_id, user_id) DO UPDATE SET status = 'pending', handled_at = NULL, handler_id = NULL`,
		circleID, userID,
	)
	if err != nil {
		return fmt.Errorf("create join request: %w", err)
	}
	return nil
}

func (r *CircleRepository) ListJoinRequests(circleID int64) ([]*model.JoinRequest, error) {
	rows, err := r.db.Query(
		`SELECT jr.id, jr.circle_id, jr.user_id, jr.status, jr.created_at,
		        jr.handled_at, jr.handler_id, u.nickname, u.avatar
		 FROM circle_join_requests jr
		 JOIN users u ON u.id = jr.user_id
		 WHERE jr.circle_id = $1 AND jr.status = 'pending'
		 ORDER BY jr.created_at ASC`, circleID,
	)
	if err != nil {
		return nil, fmt.Errorf("list join requests: %w", err)
	}
	defer rows.Close()

	var requests []*model.JoinRequest
	for rows.Next() {
		var r model.JoinRequest
		var u model.User
		if err := rows.Scan(&r.ID, &r.CircleID, &r.UserID, &r.Status, &r.CreatedAt,
			&r.HandledAt, &r.HandlerID, &u.Nickname, &u.Avatar); err != nil {
			return nil, err
		}
		r.User = &u
		requests = append(requests, &r)
	}
	return requests, nil
}

func (r *CircleRepository) HandleJoinRequest(requestID, handlerID int64, action string) error {
	status := "approved"
	if action == "reject" {
		status = "rejected"
	}
	res, err := r.db.Exec(
		`UPDATE circle_join_requests SET status = $1, handled_at = NOW(), handler_id = $2 WHERE id = $3 AND status = 'pending'`,
		status, handlerID, requestID,
	)
	if err != nil {
		return fmt.Errorf("handle join request: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("join request not found or already handled")
	}
	return nil
}

func (r *CircleRepository) ApproveJoinRequest(requestID, circleID, userID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO circle_members (circle_id, user_id, role) VALUES ($1, $2, 'member') ON CONFLICT DO NOTHING`,
		circleID, userID,
	)
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	tx.Exec("UPDATE circles SET member_count = (SELECT COUNT(*) FROM circle_members WHERE circle_id = $1) WHERE id = $1", circleID)

	return tx.Commit()
}

func (r *CircleRepository) ListMembers(circleID int) ([]*model.CircleMember, error) {
	rows, err := r.db.Query(
		`SELECT cm.id, cm.circle_id, cm.user_id, cm.role, cm.created_at,
		        u.nickname, u.avatar
		 FROM circle_members cm
		 JOIN users u ON u.id = cm.user_id
		 WHERE cm.circle_id = $1
		 ORDER BY cm.role = 'creator' DESC, cm.created_at ASC`,
		circleID,
	)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	var members []*model.CircleMember
	for rows.Next() {
		var m model.CircleMember
		var u model.User
		if err := rows.Scan(&m.ID, &m.CircleID, &m.UserID, &m.Role, &m.CreatedAt, &u.Nickname, &u.Avatar); err != nil {
			return nil, err
		}
		m.User = &u
		members = append(members, &m)
	}
	return members, nil
}

func (r *CircleRepository) Update(id, userID int64, req *model.CreateCircleRequest) error {
	joinType := req.JoinType
	if joinType == "" {
		joinType = "free"
	}
	res, err := r.db.Exec(
		`UPDATE circles SET name=$1, description=$2, icon=$3, cover=$4, join_type=$5
		 WHERE id=$6 AND creator_id=$7`,
		req.Name, req.Description, req.Icon, req.Cover, joinType, id, userID,
	)
	if err != nil {
		return fmt.Errorf("update circle: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not authorized or not found")
	}
	return nil
}

func (r *CircleRepository) CreatePost(circleID, userID int64, content string, imageUrls []string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO circle_posts (circle_id, user_id, content, image_urls)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		circleID, userID, content, pq.Array(imageUrls),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create circle post: %w", err)
	}
	r.db.Exec("UPDATE circles SET post_count = post_count + 1 WHERE id = $1", circleID)
	return id, nil
}

func (r *CircleRepository) ListPosts(circleID, offset, limit int, userID int64) ([]*model.CirclePost, error) {
	rows, err := r.db.Query(
		`SELECT cp.id, cp.circle_id, cp.user_id, cp.content, cp.image_urls,
		        cp.like_count, cp.comment_count, cp.created_at, cp.is_pinned,
		        u.nickname, u.avatar,
		        CASE WHEN cpl.id IS NOT NULL THEN true ELSE false END
		 FROM circle_posts cp
		 JOIN users u ON u.id = cp.user_id
		 LEFT JOIN circle_post_likes cpl ON cpl.post_id = cp.id AND cpl.user_id = $4
		 WHERE cp.circle_id = $1
		 ORDER BY cp.is_pinned DESC, cp.created_at DESC LIMIT $2 OFFSET $3`,
		circleID, limit, offset, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list circle posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.CirclePost
	for rows.Next() {
		var p model.CirclePost
		var u model.User
		if err := rows.Scan(&p.ID, &p.CircleID, &p.UserID, &p.Content, pq.Array(&p.ImageUrls),
			&p.LikeCount, &p.CommentCount, &p.CreatedAt, &p.IsPinned, &u.Nickname, &u.Avatar, &p.IsLiked); err != nil {
			return nil, err
		}
		p.Author = &u
		posts = append(posts, &p)
	}
	return posts, nil
}

func (r *CircleRepository) ToggleLike(postID, userID int64) (bool, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO circle_post_likes (post_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (post_id, user_id) DO NOTHING RETURNING id`,
		postID, userID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		r.db.Exec("DELETE FROM circle_post_likes WHERE post_id=$1 AND user_id=$2", postID, userID)
		r.db.Exec("UPDATE circle_posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1", postID)
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("toggle like: %w", err)
	}
	r.db.Exec("UPDATE circle_posts SET like_count = like_count + 1 WHERE id = $1", postID)
	return true, nil
}

func (r *CircleRepository) CreateComment(postID, userID int64, content string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO circle_post_comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id`,
		postID, userID, content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create comment: %w", err)
	}
	r.db.Exec("UPDATE circle_posts SET comment_count = comment_count + 1 WHERE id = $1", postID)
	return id, nil
}

func (r *CircleRepository) ListComments(postID int) ([]*model.CirclePostComment, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname, u.avatar
		 FROM circle_post_comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.post_id = $1
		 ORDER BY c.created_at ASC`, postID,
	)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.CirclePostComment
	for rows.Next() {
		var c model.CirclePostComment
		var u model.User
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &u.Nickname, &u.Avatar); err != nil {
			return nil, err
		}
		c.Author = &u
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *CircleRepository) Search(query string, offset, limit int) ([]*model.Circle, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.name, c.description, c.icon, c.creator_id, c.member_count, c.post_count, c.created_at,
			u.nickname, u.avatar
		FROM circles c JOIN users u ON u.id = c.creator_id
		WHERE c.name ILIKE $1 OR c.description ILIKE $1
		ORDER BY similarity(c.name, $4) DESC, c.member_count DESC LIMIT $2 OFFSET $3`,
		"%"+query+"%", limit, offset, query)
	if err != nil {
		return nil, fmt.Errorf("search circles: %w", err)
	}
	defer rows.Close()

	var circles []*model.Circle
	for rows.Next() {
		var c model.Circle
		var u model.User
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Icon, &c.CreatorID, &c.MemberCount, &c.PostCount, &c.CreatedAt, &u.Nickname, &u.Avatar); err != nil {
			return nil, err
		}
		c.Creator = &u
		circles = append(circles, &c)
	}
	return circles, nil
}

type AvatarWallEntry struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

func (r *CircleRepository) MemberAvatarWall(circleID int64) ([]*AvatarWallEntry, error) {
	rows, err := r.db.Query(
		`SELECT cm.user_id, u.nickname, u.avatar, cm.role
		 FROM circle_members cm
		 JOIN users u ON u.id = cm.user_id
		 WHERE cm.circle_id = $1
		 ORDER BY cm.role = 'creator' DESC, cm.created_at ASC`,
		circleID,
	)
	if err != nil {
		return nil, fmt.Errorf("avatar wall: %w", err)
	}
	defer rows.Close()

	var entries []*AvatarWallEntry
	for rows.Next() {
		var e AvatarWallEntry
		if err := rows.Scan(&e.UserID, &e.Nickname, &e.Avatar, &e.Role); err != nil {
			return nil, err
		}
		entries = append(entries, &e)
	}
	return entries, nil
}

func (r *CircleRepository) CreateActivity(circleID, userID int64, action, content string) error {
	_, err := r.db.Exec(
		`INSERT INTO circle_activities (circle_id, user_id, action, content) VALUES ($1, $2, $3, $4)`,
		circleID, userID, action, content,
	)
	return err
}

func (r *CircleRepository) ListActivities(circleID int64, limit int) ([]*model.CircleActivity, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT ca.id, ca.circle_id, ca.user_id, ca.action, ca.content, ca.created_at,
		        u.nickname, u.avatar
		 FROM circle_activities ca
		 JOIN users u ON u.id = ca.user_id
		 WHERE ca.circle_id = $1
		 ORDER BY ca.created_at DESC LIMIT $2`, circleID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	defer rows.Close()

	var acts []*model.CircleActivity
	for rows.Next() {
		var a model.CircleActivity
		var u model.User
		if err := rows.Scan(&a.ID, &a.CircleID, &a.UserID, &a.Action, &a.Content, &a.CreatedAt,
			&u.Nickname, &u.Avatar); err != nil {
			continue
		}
		a.User = &u
		acts = append(acts, &a)
	}
	if acts == nil {
		acts = []*model.CircleActivity{}
	}
	return acts, nil
}
