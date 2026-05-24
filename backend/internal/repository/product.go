package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/guish/university-social/internal/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(userID int64, req *model.CreateProductRequest) (int64, error) {
	imagesJSON, _ := json.Marshal(req.Images)
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO products (user_id, title, description, price, original_price, category, condition, images, contact, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id`,
		userID, req.Title, req.Description, req.Price, req.OriginalPrice,
		req.Category, req.Condition, imagesJSON, req.Contact,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create product: %w", err)
	}
	return id, nil
}

func (r *ProductRepository) scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}, currentUserID int64) (*model.Product, error) {
	var p model.Product
	var desc, contact, condition sql.NullString
	var originalPrice sql.NullFloat64
	var imagesJSON []byte
	var isLiked bool

	err := scanner.Scan(&p.ID, &p.UserID, &p.Title, &desc, &p.Price,
		&originalPrice, &p.Category, &condition, &imagesJSON, &contact,
		&p.Status, &p.LikeCount, &p.CommentCount, &p.CreatedAt, &p.UpdatedAt, &isLiked)
	if err != nil {
		return nil, err
	}

	p.Description = desc.String
	p.Contact = contact.String
	p.Condition = condition.String
	if originalPrice.Valid {
		p.OriginalPrice = &originalPrice.Float64
	}
	if len(imagesJSON) > 0 {
		json.Unmarshal(imagesJSON, &p.Images)
	}
	if p.Images == nil {
		p.Images = []string{}
	}

	p.IsLiked = isLiked
	p.IsOwner = p.UserID == currentUserID

	return &p, nil
}

const productSelectCols = `p.id, p.user_id, p.title, p.description, p.price,
	p.original_price, p.category, p.condition, p.images, p.contact,
	p.status, p.like_count, p.comment_count, p.created_at, p.updated_at,
	CASE WHEN l.id IS NOT NULL THEN true ELSE false END as is_liked`

const productFromJoins = `FROM products p
	LEFT JOIN product_likes l ON l.product_id = p.id AND l.user_id = $1`

func (r *ProductRepository) List(category, sort string, status, offset, limit int, currentUserID int64) ([]*model.Product, error) {
	args := []interface{}{currentUserID}
	argIdx := 2
	where := []string{"p.status = $3"}
	args = append(args, status)

	if category != "" && category != "全部" {
		where = append(where, fmt.Sprintf("p.category = $%d", argIdx+1))
		args = append(args, category)
		argIdx++
	}

	orderBy := "p.created_at DESC"
	if sort == "price_asc" {
		orderBy = "p.price ASC"
	} else if sort == "price_desc" {
		orderBy = "p.price DESC"
	}

	query := fmt.Sprintf(`SELECT %s %s WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		productSelectCols, productFromJoins,
		strings.Join(where, " AND "), orderBy, argIdx+1, argIdx+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p, err := r.scanProduct(rows, currentUserID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) FindByID(id int64, currentUserID int64) (*model.Product, error) {
	query := fmt.Sprintf(`SELECT %s %s WHERE p.id = $2`, productSelectCols, productFromJoins)
	row := r.db.QueryRow(query, currentUserID, id)
	p, err := r.scanProduct(row, currentUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("find product: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) Delete(id, userID int64) error {
	res, err := r.db.Exec(`DELETE FROM products WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("product not found or not owned")
	}
	return nil
}

func (r *ProductRepository) MarkSold(id, userID int64) error {
	res, err := r.db.Exec(`UPDATE products SET status = 1, updated_at = NOW() WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("mark sold: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("product not found or not owned")
	}
	return nil
}

func (r *ProductRepository) ListByUser(userID, currentUserID int64, offset, limit int) ([]*model.Product, error) {
	query := fmt.Sprintf(`SELECT %s %s WHERE p.user_id = $2 ORDER BY p.created_at DESC LIMIT $3 OFFSET $4`,
		productSelectCols, productFromJoins)
	rows, err := r.db.Query(query, currentUserID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list user products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p, err := r.scanProduct(rows, currentUserID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) SearchProducts(query string, offset, limit int, currentUserID int64) ([]*model.Product, error) {
	q := fmt.Sprintf(`SELECT %s %s WHERE p.status = 0 AND (p.title ILIKE $2 OR p.description ILIKE $2) ORDER BY similarity(p.title, $5) DESC, p.created_at DESC LIMIT $3 OFFSET $4`,
		productSelectCols, productFromJoins)
	pattern := "%" + query + "%"
	rows, err := r.db.Query(q, currentUserID, pattern, limit, offset, query)
	if err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p, err := r.scanProduct(rows, currentUserID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) ToggleLike(userID, productID int64) (bool, error) {
	var id int64
	err := r.db.QueryRow(
		`SELECT id FROM product_likes WHERE user_id = $1 AND product_id = $2`,
		userID, productID,
	).Scan(&id)
	if err == nil {
		_, err = r.db.Exec(`DELETE FROM product_likes WHERE id = $1`, id)
		if err != nil {
			return false, fmt.Errorf("unlike product: %w", err)
		}
		_, _ = r.db.Exec(`UPDATE products SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1`, productID)
		return false, nil
	}
	_, err = r.db.Exec(`INSERT INTO product_likes (user_id, product_id) VALUES ($1, $2)`, userID, productID)
	if err != nil {
		return false, fmt.Errorf("like product: %w", err)
	}
	_, _ = r.db.Exec(`UPDATE products SET like_count = like_count + 1 WHERE id = $1`, productID)
	return true, nil
}

func (r *ProductRepository) CreateComment(userID, productID int64, content string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO product_comments (product_id, user_id, content) VALUES ($1, $2, $3) RETURNING id`,
		productID, userID, content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create product comment: %w", err)
	}
	_, _ = r.db.Exec(`UPDATE products SET comment_count = comment_count + 1 WHERE id = $1`, productID)
	return id, nil
}

func (r *ProductRepository) ListComments(productID int64) ([]*model.ProductComment, error) {
	rows, err := r.db.Query(
		`SELECT pc.id, pc.product_id, pc.user_id, pc.content, pc.created_at
		 FROM product_comments pc ORDER BY pc.created_at ASC`, productID)
	if err != nil {
		return nil, fmt.Errorf("list product comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.ProductComment
	for rows.Next() {
		var c model.ProductComment
		if err := rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *ProductRepository) Categories() []string {
	return []string{"全部", "教材", "电子", "生活", "运动", "其他"}
}

func (r *ProductRepository) Counts(userID int64) (total, active int) {
	r.db.QueryRow(`SELECT COUNT(*) FROM products WHERE user_id = $1`, userID).Scan(&total)
	r.db.QueryRow(`SELECT COUNT(*) FROM products WHERE user_id = $1 AND status = 0`, userID).Scan(&active)
	return
}

func (r *ProductRepository) Trending(limit int, currentUserID int64) ([]*model.Product, error) {
	query := fmt.Sprintf(`SELECT %s %s WHERE p.status = 0 ORDER BY (p.like_count * 3 + p.comment_count * 2) DESC, p.created_at DESC LIMIT $2`,
		productSelectCols, productFromJoins)
	rows, err := r.db.Query(query, currentUserID, limit)
	if err != nil {
		return nil, fmt.Errorf("trending products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p, err := r.scanProduct(rows, currentUserID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if products == nil {
		products = []*model.Product{}
	}
	return products, nil
}

func (r *ProductRepository) scanNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (r *ProductRepository) ToFixed(f float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(f*pow) / pow
}

func (r *ProductRepository) AdminList(offset, limit int) ([]*model.Product, error) {
	col := `p.id, p.user_id, p.title, p.description, p.price,
		p.original_price, p.category, p.condition, p.images, p.contact,
		p.status, p.like_count, p.comment_count, p.created_at, p.updated_at, false as is_liked`
	rows, err := r.db.Query(`SELECT `+col+` FROM products p ORDER BY p.id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("admin list products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p, err := r.scanProduct(rows, 0)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) AdminDelete(id int64) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
