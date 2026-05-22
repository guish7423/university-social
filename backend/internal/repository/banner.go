package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type BannerRepository struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) ListActive() ([]*model.Banner, error) {
	rows, err := r.db.Query(
		`SELECT id, title, image_url, link_url, sort_order, is_active, created_at
		 FROM banners WHERE is_active = true ORDER BY sort_order ASC, id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list banners: %w", err)
	}
	defer rows.Close()

	var banners []*model.Banner
	for rows.Next() {
		var b model.Banner
		if err := rows.Scan(&b.ID, &b.Title, &b.ImageURL, &b.LinkURL, &b.SortOrder, &b.IsActive, &b.CreatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, &b)
	}
	return banners, nil
}

func (r *BannerRepository) ListAll(offset, limit int) ([]*model.Banner, error) {
	rows, err := r.db.Query(
		`SELECT id, title, image_url, link_url, sort_order, is_active, created_at
		 FROM banners ORDER BY sort_order ASC, id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list all banners: %w", err)
	}
	defer rows.Close()

	var banners []*model.Banner
	for rows.Next() {
		var b model.Banner
		if err := rows.Scan(&b.ID, &b.Title, &b.ImageURL, &b.LinkURL, &b.SortOrder, &b.IsActive, &b.CreatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, &b)
	}
	return banners, nil
}

func (r *BannerRepository) Create(req *model.CreateBannerRequest) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO banners (title, image_url, link_url, sort_order) VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Title, req.ImageURL, req.LinkURL, req.SortOrder,
	).Scan(&id)
	return id, err
}

func (r *BannerRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM banners WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("banner not found")
	}
	return nil
}

func (r *BannerRepository) ToggleActive(id int64) error {
	_, err := r.db.Exec(`UPDATE banners SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *BannerRepository) ListAnnouncementsActive() ([]*model.Announcement, error) {
	rows, err := r.db.Query(
		`SELECT id, title, content, priority, is_active, created_at
		 FROM announcements WHERE is_active = true ORDER BY priority DESC, id DESC LIMIT 10`,
	)
	if err != nil {
		return nil, fmt.Errorf("list announcements: %w", err)
	}
	defer rows.Close()

	var anns []*model.Announcement
	for rows.Next() {
		var a model.Announcement
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Priority, &a.IsActive, &a.CreatedAt); err != nil {
			return nil, err
		}
		anns = append(anns, &a)
	}
	return anns, nil
}

func (r *BannerRepository) ListAllAnnouncements(offset, limit int) ([]*model.Announcement, error) {
	rows, err := r.db.Query(
		`SELECT id, title, content, priority, is_active, created_at
		 FROM announcements ORDER BY priority DESC, id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list all announcements: %w", err)
	}
	defer rows.Close()

	var anns []*model.Announcement
	for rows.Next() {
		var a model.Announcement
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Priority, &a.IsActive, &a.CreatedAt); err != nil {
			return nil, err
		}
		anns = append(anns, &a)
	}
	return anns, nil
}

func (r *BannerRepository) CreateAnnouncement(req *model.CreateAnnouncementRequest) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO announcements (title, content, priority) VALUES ($1, $2, $3) RETURNING id`,
		req.Title, req.Content, req.Priority,
	).Scan(&id)
	return id, err
}

func (r *BannerRepository) DeleteAnnouncement(id int64) error {
	res, err := r.db.Exec(`DELETE FROM announcements WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("announcement not found")
	}
	return nil
}

func (r *BannerRepository) ToggleAnnouncementActive(id int64) error {
	_, err := r.db.Exec(`UPDATE announcements SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *BannerRepository) SetPostFeatured(postID int64, featured bool) error {
	res, err := r.db.Exec(`UPDATE posts SET is_featured = $1 WHERE id = $2`, featured, postID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

func (r *BannerRepository) ListFeaturedPosts(offset, limit int) ([]*model.Post, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.images, p.like_count, p.comment_count, p.is_featured, p.created_at, u.nickname, u.avatar
		 FROM posts p JOIN users u ON u.id = p.user_id
		 WHERE p.is_featured = true
		 ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list featured posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		var u model.User
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Images, &p.LikeCount, &p.CommentCount, &p.IsFeatured, &p.CreatedAt, &u.Nickname, &u.Avatar); err != nil {
			return nil, err
		}
		p.Author = &u
		posts = append(posts, &p)
	}
	return posts, nil
}
