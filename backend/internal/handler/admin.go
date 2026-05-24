package handler

import (
"database/sql"
"encoding/json"
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

func (h *AdminHandler) logOp(c *gin.Context, action, targetType string, targetID int64, detail any) {
	adminID := c.GetInt64("user_id")
	var adminName string
	h.db.QueryRow("SELECT nickname FROM users WHERE id = $1", adminID).Scan(&adminName)
	var detailJSON []byte
	if detail != nil {
		detailJSON, _ = json.Marshal(detail)
	}
	h.db.Exec(
		"INSERT INTO operation_logs (admin_id, admin_name, action, target_type, target_id, detail) VALUES ($1,$2,$3,$4,$5,$6)",
		adminID, adminName, action, targetType, targetID, detailJSON,
	)
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
	h.logOp(c, "ban_user", "user", id, nil)
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
	h.logOp(c, "unban_user", "user", id, nil)
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
	h.logOp(c, "delete_post", "post", id, nil)
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
	h.logOp(c, "delete_circle", "circle", id, nil)
}

func (h *AdminHandler) ListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int
	h.db.QueryRow("SELECT COUNT(*) FROM operation_logs").Scan(&total)

	type LogItem struct {
		ID         int64  `json:"id"`
		AdminID    int64  `json:"admin_id"`
		AdminName  string `json:"admin_name"`
		Action     string `json:"action"`
		TargetType string `json:"target_type"`
		TargetID   int64  `json:"target_id"`
		Detail     any    `json:"detail,omitempty"`
		CreatedAt  string `json:"created_at"`
	}

	rows, err := h.db.Query(
		"SELECT id, admin_id, admin_name, action, target_type, target_id, COALESCE(detail::text, ''), created_at"+
		" FROM operation_logs ORDER BY id DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取操作日志失败"})
		return
	}
	defer rows.Close()

	var logs []*LogItem
	for rows.Next() {
		var l LogItem
		var detailStr, createdAt string
		if err := rows.Scan(&l.ID, &l.AdminID, &l.AdminName, &l.Action, &l.TargetType, &l.TargetID, &detailStr, &createdAt); err != nil {
			continue
		}
		l.CreatedAt = createdAt
		if detailStr != "" {
			var detail any
			json.Unmarshal([]byte(detailStr), &detail)
			l.Detail = detail
		}
		logs = append(logs, &l)
	}
	if logs == nil {
		logs = []*LogItem{}
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs, "total": total})
}

func (h *AdminHandler) ListProducts(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	col := `p.id, p.user_id, p.title, p.description, p.price,
		p.original_price, p.category, p.condition, p.images, p.contact,
		p.status, p.like_count, p.comment_count, p.created_at, p.updated_at, false as is_liked`
	rows, err := h.db.Query(`SELECT `+col+` FROM products p ORDER BY p.id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品列表失败"})
		return
	}
	defer rows.Close()

	type ProductItem struct {
		ID            int64   `json:"id"`
		UserID        int64   `json:"user_id"`
		Title         string  `json:"title"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		OriginalPrice *float64 `json:"original_price,omitempty"`
		Category      string  `json:"category"`
		Condition     string  `json:"condition"`
		Status        int     `json:"status"`
		LikeCount     int     `json:"like_count"`
		CommentCount  int     `json:"comment_count"`
		CreatedAt     string  `json:"created_at"`
		AuthorName    string  `json:"author_name"`
	}

	var products []*ProductItem
	for rows.Next() {
		var p ProductItem
		var desc, cond sql.NullString
		var origP sql.NullFloat64
		var createdAt string
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &desc, &p.Price,
			&origP, &p.Category, &cond, nil, nil,
			&p.Status, &p.LikeCount, &p.CommentCount, &createdAt, nil, nil); err != nil {
			continue
		}
		p.Description = desc.String
		if origP.Valid { p.OriginalPrice = &origP.Float64 }
		p.Condition = cond.String
		p.CreatedAt = createdAt
		// Get author name
		h.db.QueryRow("SELECT nickname FROM users WHERE id = $1", p.UserID).Scan(&p.AuthorName)
		products = append(products, &p)
	}
	if products == nil {
		products = []*ProductItem{}
	}
	c.JSON(http.StatusOK, products)
}

func (h *AdminHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
	h.logOp(c, "delete_product", "product", id, nil)
}


func (h *AdminHandler) ListFoundItems(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	type Item struct {
		ID        int64  `json:"id"`
		UserID    int64  `json:"user_id"`
		Title     string `json:"title"`
		Category  string `json:"category"`
		Status    int    `json:"status"`
		CreatedAt string `json:"created_at"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
	}

	rows, err := h.db.Query(
		`SELECT l.id, l.user_id, l.title, l.category, l.status, l.created_at,
		        COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
		 FROM lost_items l
		 LEFT JOIN users u ON l.user_id = u.id
		 ORDER BY l.id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失物招领列表失败"})
		return
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var i Item
		var createdAt string
		if err := rows.Scan(&i.ID, &i.UserID, &i.Title, &i.Category, &i.Status, &createdAt, &i.Nickname, &i.Avatar); err != nil {
			continue
		}
		i.CreatedAt = createdAt
		items = append(items, &i)
	}
	if items == nil {
		items = []*Item{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminHandler) DeleteFoundItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM lost_items WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
	h.logOp(c, "delete_found_item", "lost_item", id, nil)
}

func (h *AdminHandler) ListUniversities(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	type Item struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		MemberCount int    `json:"member_count"`
		CreatedAt   string `json:"created_at"`
	}

	rows, err := h.db.Query(
		`SELECT id, name, COALESCE(description, ''), COALESCE(member_count, 0), created_at
		 FROM universities ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取高校列表失败"})
		return
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var i Item
		var createdAt string
		if err := rows.Scan(&i.ID, &i.Name, &i.Description, &i.MemberCount, &createdAt); err != nil {
			continue
		}
		i.CreatedAt = createdAt
		items = append(items, &i)
	}
	if items == nil {
		items = []*Item{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminHandler) DeleteUniversity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM universities WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
	h.logOp(c, "delete_university", "university", id, nil)
}

func (h *AdminHandler) ListAdminActivities(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	type Item struct {
		ID        int64  `json:"id"`
		Title     string `json:"title"`
		Type      string `json:"type"`
		Status    int    `json:"status"`
		CreatedAt string `json:"created_at"`
		Nickname  string `json:"nickname"`
	}

	rows, err := h.db.Query(
		`SELECT a.id, a.title, a.type, a.status, a.created_at,
		        COALESCE(u.nickname, '')
		 FROM activities a
		 LEFT JOIN users u ON a.user_id = u.id
		 ORDER BY a.id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取活动列表失败"})
		return
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var i Item
		var createdAt string
		if err := rows.Scan(&i.ID, &i.Title, &i.Type, &i.Status, &createdAt, &i.Nickname); err != nil {
			continue
		}
		i.CreatedAt = createdAt
		items = append(items, &i)
	}
	if items == nil {
		items = []*Item{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminHandler) DeleteActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM activities WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
	h.logOp(c, "delete_activity", "activity", id, nil)
}

func (h *AdminHandler) ListAdminCourses(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	type Item struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Code      string  `json:"code"`
		Department string `json:"department"`
		Teacher   string  `json:"teacher"`
		Credit    float64 `json:"credit"`
		AvgRating float64 `json:"avg_rating"`
		CreatedAt string  `json:"created_at"`
	}

	rows, err := h.db.Query(
		`SELECT id, name, COALESCE(code, ''), COALESCE(department, ''), COALESCE(teacher, ''),
		        COALESCE(credit, 0), COALESCE(avg_rating, 0), created_at
		 FROM courses ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程列表失败"})
		return
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var i Item
		var createdAt string
		if err := rows.Scan(&i.ID, &i.Name, &i.Code, &i.Department, &i.Teacher, &i.Credit, &i.AvgRating, &createdAt); err != nil {
			continue
		}
		i.CreatedAt = createdAt
		items = append(items, &i)
	}
	if items == nil {
		items = []*Item{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err = h.db.Exec("DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
	h.logOp(c, "delete_course", "course", id, nil)
}

func (h *AdminHandler) ListVerifications(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	type Item struct {
		ID        int64  `json:"id"`
		UserID    int64  `json:"user_id"`
		Code      string `json:"code"`
		Type      string `json:"type"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
		Nickname  string `json:"nickname"`
		IsVerified bool  `json:"is_verified"`
	}

	rows, err := h.db.Query(
		`SELECT v.id, v.user_id, v.code, COALESCE(v.type, ''),
		        CASE WHEN v.used_at IS NOT NULL THEN 'verified' WHEN v.expires_at < NOW() THEN 'expired' ELSE 'pending' END,
		        v.created_at, COALESCE(u.nickname, ''), COALESCE(u.is_verified, false)
		 FROM verification_codes v
		 LEFT JOIN users u ON v.user_id = u.id
		 ORDER BY v.id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取认证记录失败"})
		return
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var i Item
		var createdAt string
		if err := rows.Scan(&i.ID, &i.UserID, &i.Code, &i.Type, &i.Status, &createdAt, &i.Nickname, &i.IsVerified); err != nil {
			continue
		}
		i.CreatedAt = createdAt
		items = append(items, &i)
	}
	if items == nil {
		items = []*Item{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminHandler) ApproveVerification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	// Get user_id from the verification code
	var userID int64
	err = h.db.QueryRow("SELECT user_id FROM verification_codes WHERE id = $1", id).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "认证记录不存在"})
		return
	}
	// Mark code as used
	_, err = h.db.Exec("UPDATE verification_codes SET used_at = NOW() WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	// Set user as verified
	_, err = h.db.Exec("UPDATE users SET is_verified = true WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户状态失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已通过认证"})
	h.logOp(c, "approve_verification", "verification", id, gin.H{"user_id": userID})
}

func (h *AdminHandler) RejectVerification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	// Delete the verification code record
	_, err = h.db.Exec("DELETE FROM verification_codes WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已拒绝"})
	h.logOp(c, "reject_verification", "verification", id, nil)
}