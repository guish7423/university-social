package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type LostItemRepository struct {
	db *sql.DB
}

func NewLostItemRepository(db *sql.DB) *LostItemRepository {
	return &LostItemRepository{db: db}
}

func (r *LostItemRepository) Create(userID int64, req *model.CreateLostItemRequest) (*model.LostItem, error) {
	imagesJSON, _ := json.Marshal(req.Images)
	var item model.LostItem
	err := r.db.QueryRow(
		`INSERT INTO lost_items (user_id, title, description, category, item_type, location, contact, images)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, user_id, title, description, category, item_type, location, contact, images, status, created_at, updated_at`,
		userID, req.Title, req.Description, req.Category, req.ItemType, req.Location, req.Contact, imagesJSON,
	).Scan(&item.ID, &item.UserID, &item.Title, &item.Description, &item.Category, &item.ItemType,
		&item.Location, &item.Contact, &item.Images, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create lost item: %w", err)
	}
	return &item, nil
}

func (r *LostItemRepository) FindByID(id int64, currentUserID int64) (*model.LostItem, error) {
	var item model.LostItem
	var imagesStr string
	err := r.db.QueryRow(
		`SELECT l.id, l.user_id, l.title, l.description, l.category, l.item_type,
		        l.location, l.contact, l.images, l.status, l.created_at, l.updated_at,
		        COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
		 FROM lost_items l
		 LEFT JOIN users u ON l.user_id = u.id
		 WHERE l.id = $1`, id,
	).Scan(&item.ID, &item.UserID, &item.Title, &item.Description, &item.Category, &item.ItemType,
		&item.Location, &item.Contact, &imagesStr, &item.Status, &item.CreatedAt, &item.UpdatedAt,
		&item.Nickname, &item.Avatar)
	if err != nil {
		return nil, fmt.Errorf("find lost item: %w", err)
	}
	json.Unmarshal([]byte(imagesStr), &item.Images)
	return &item, nil
}

func (r *LostItemRepository) List(category string, status int, offset, limit int) ([]*model.LostItem, error) {
	query := `SELECT l.id, l.user_id, l.title, l.description, l.category, l.item_type,
	                 l.location, l.contact, l.images, l.status, l.created_at, l.updated_at,
	                 COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
	          FROM lost_items l
	          LEFT JOIN users u ON l.user_id = u.id
	          WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if category == "lost" || category == "found" {
		query += fmt.Sprintf(" AND l.category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	if status >= 0 {
		query += fmt.Sprintf(" AND l.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query += ` ORDER BY l.created_at DESC`
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list lost items: %w", err)
	}
	defer rows.Close()

	var items []*model.LostItem
	for rows.Next() {
		var item model.LostItem
		var imagesStr string
		err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Description, &item.Category, &item.ItemType,
			&item.Location, &item.Contact, &imagesStr, &item.Status, &item.CreatedAt, &item.UpdatedAt,
			&item.Nickname, &item.Avatar)
		if err != nil {
			return nil, fmt.Errorf("scan lost item: %w", err)
		}
		json.Unmarshal([]byte(imagesStr), &item.Images)
		items = append(items, &item)
	}
	if items == nil {
		return []*model.LostItem{}, nil
	}
	return items, nil
}

func (r *LostItemRepository) Delete(id int64, userID int64) error {
	result, err := r.db.Exec(`DELETE FROM lost_items WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete lost item: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("lost item not found or not owned")
	}
	return nil
}

func (r *LostItemRepository) UpdateStatus(id int64, userID int64, status int) error {
	result, err := r.db.Exec(
		`UPDATE lost_items SET status = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		status, id, userID)
	if err != nil {
		return fmt.Errorf("update lost item status: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("lost item not found or not owned")
	}
	return nil
}

func (r *LostItemRepository) ListByUser(userID int64, offset, limit int) ([]*model.LostItem, error) {
	rows, err := r.db.Query(
		`SELECT l.id, l.user_id, l.title, l.description, l.category, l.item_type,
		        l.location, l.contact, l.images, l.status, l.created_at, l.updated_at,
		        COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
		 FROM lost_items l
		 LEFT JOIN users u ON l.user_id = u.id
		 WHERE l.user_id = $1
		 ORDER BY l.created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list user lost items: %w", err)
	}
	defer rows.Close()

	var items []*model.LostItem
	for rows.Next() {
		var item model.LostItem
		var imagesStr string
		err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Description, &item.Category, &item.ItemType,
			&item.Location, &item.Contact, &imagesStr, &item.Status, &item.CreatedAt, &item.UpdatedAt,
			&item.Nickname, &item.Avatar)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		json.Unmarshal([]byte(imagesStr), &item.Images)
		items = append(items, &item)
	}
	if items == nil {
		return []*model.LostItem{}, nil
	}
	return items, nil
}

func (r *LostItemRepository) Search(query, category string, offset, limit int) ([]*model.LostItem, error) {
	like := "%" + query + "%"
	args := []interface{}{like, limit, offset}
	querySQL := `SELECT l.id, l.user_id, l.title, l.description, l.category, l.item_type,
	                    l.location, l.contact, l.images, l.status, l.created_at, l.updated_at,
	                    COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
	             FROM lost_items l
	             LEFT JOIN users u ON l.user_id = u.id
	             WHERE (l.title ILIKE $1 OR l.description ILIKE $1 OR l.item_type ILIKE $1 OR l.location ILIKE $1)`
	argIdx := 4
	if category == "lost" || category == "found" {
		querySQL += fmt.Sprintf(" AND l.category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	querySQL += fmt.Sprintf(" ORDER BY l.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("search lost items: %w", err)
	}
	defer rows.Close()

	var items []*model.LostItem
	for rows.Next() {
		var item model.LostItem
		var imagesStr string
		err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Description, &item.Category, &item.ItemType,
			&item.Location, &item.Contact, &imagesStr, &item.Status, &item.CreatedAt, &item.UpdatedAt,
			&item.Nickname, &item.Avatar)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		json.Unmarshal([]byte(imagesStr), &item.Images)
		items = append(items, &item)
	}
	if items == nil {
		return []*model.LostItem{}, nil
	}
	return items, nil
}
