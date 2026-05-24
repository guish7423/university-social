package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	var userCount, postCount, circleCount int
	h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	h.db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&postCount)
	h.db.QueryRow("SELECT COUNT(*) FROM circles").Scan(&circleCount)
	c.JSON(http.StatusOK, gin.H{
		"user_count":   userCount,
		"post_count":   postCount,
		"circle_count": circleCount,
	})
}

func (h *AdminHandler) DailyStats(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n > 0 && n <= 365 {
			days = n
		}
	}

	rows, err := h.db.Query(
		`SELECT date, dau, new_users, new_posts, new_comments, new_likes, new_circles
		 FROM daily_stats ORDER BY date DESC LIMIT $1`, days,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败"})
		return
	}
	defer rows.Close()

	type DailyStat struct {
		Date        string `json:"date"`
		Dau         int    `json:"dau"`
		NewUsers    int    `json:"new_users"`
		NewPosts    int    `json:"new_posts"`
		NewComments int    `json:"new_comments"`
		NewLikes    int    `json:"new_likes"`
		NewCircles  int    `json:"new_circles"`
	}

	var stats []*DailyStat
	for rows.Next() {
		var s DailyStat
		if err := rows.Scan(&s.Date, &s.Dau, &s.NewUsers, &s.NewPosts, &s.NewComments, &s.NewLikes, &s.NewCircles); err != nil {
			continue
		}
		stats = append(stats, &s)
	}
	if stats == nil {
		stats = []*DailyStat{}
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	q := c.Query("q")
	offset := (page - 1) * limit

	where := ""
	args := []any{limit, offset}
	if q != "" {
		where = " WHERE (nickname ILIKE $3 OR open_id ILIKE $3)"
		args = []any{limit, offset, "%" + q + "%"}
	}

	var total int
	h.db.QueryRow("SELECT COUNT(*) FROM users"+where, args[1:]...).Scan(&total)

	rows, err := h.db.Query(
		"SELECT id, nickname, avatar, open_id, school, is_verified, role, is_banned, created_at"+
		" FROM users"+where+" ORDER BY id DESC LIMIT $1 OFFSET $2", args...,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}
	defer rows.Close()

	type UserItem struct {
		ID         int64  `json:"id"`
		Nickname   string `json:"nickname"`
		Avatar     string `json:"avatar"`
		OpenID     string `json:"open_id"`
		School     string `json:"school"`
		IsVerified bool   `json:"is_verified"`
		Role       string `json:"role"`
		IsBanned   bool   `json:"is_banned"`
		CreatedAt  string `json:"created_at"`
	}

	var users []*UserItem
	for rows.Next() {
		var u UserItem
		var createdAt string
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Avatar, &u.OpenID, &u.School, &u.IsVerified, &u.Role, &u.IsBanned, &createdAt); err != nil {
			continue
		}
		u.CreatedAt = createdAt
		users = append(users, &u)
	}
	if users == nil {
		users = []*UserItem{}
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "total": total})
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("UPDATE users SET is_banned = true WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已封禁"})
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("UPDATE users SET is_banned = false WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已解封"})
}

func (h *AdminHandler) ListPosts(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	rows, err := h.db.Query(
		`SELECT p.id, p.content, p.images, p.like_count, p.comment_count, p.created_at,
		        u.nickname, u.avatar
		 FROM posts p JOIN users u ON u.id = p.user_id
		 ORDER BY p.id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子列表失败"})
		return
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		var u model.User
		if err := rows.Scan(&p.ID, &p.Content, &p.Images, &p.LikeCount, &p.CommentCount, &p.CreatedAt,
			&u.Nickname, &u.Avatar); err != nil {
			continue
		}
		p.Author = &u
		posts = append(posts, &p)
	}
	if posts == nil {
		posts = []*model.Post{}
	}
	c.JSON(http.StatusOK, posts)
}

func (h *AdminHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

func (h *AdminHandler) ListCircles(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	rows, err := h.db.Query(
		`SELECT c.id, c.name, c.description, c.member_count, c.post_count, c.created_at
		 FROM circles c ORDER BY c.id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取圈子列表失败"})
		return
	}
	defer rows.Close()

	type CircleItem struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		MemberCount int    `json:"member_count"`
		PostCount   int    `json:"post_count"`
		CreatedAt   string `json:"created_at"`
	}

	var circles []*CircleItem
	for rows.Next() {
		var c CircleItem
		var createdAt string
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.MemberCount, &c.PostCount, &createdAt); err != nil {
			continue
		}
		c.CreatedAt = createdAt
		circles = append(circles, &c)
	}
	if circles == nil {
		circles = []*CircleItem{}
	}
	c.JSON(http.StatusOK, circles)
}

func (h *AdminHandler) DeleteCircle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM circles WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
