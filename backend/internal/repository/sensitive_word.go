package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type SensitiveWordRepository struct {
	db *sql.DB
}

func NewSensitiveWordRepository(db *sql.DB) *SensitiveWordRepository {
	return &SensitiveWordRepository{db: db}
}

func (r *SensitiveWordRepository) List() ([]string, error) {
	rows, err := r.db.Query(`SELECT word FROM sensitive_words ORDER BY LENGTH(word) DESC`)
	if err != nil {
		return nil, fmt.Errorf("list sensitive words: %w", err)
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var w string
		if err := rows.Scan(&w); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, nil
}

func (r *SensitiveWordRepository) Add(word string) error {
	_, err := r.db.Exec(`INSERT INTO sensitive_words (word) VALUES ($1) ON CONFLICT DO NOTHING`, word)
	return err
}

func (r *SensitiveWordRepository) Remove(word string) error {
	_, err := r.db.Exec(`DELETE FROM sensitive_words WHERE word = $1`, word)
	return err
}

// ReplaceSensitive replaces sensitive words in content with ***
func (r *SensitiveWordRepository) ReplaceSensitive(content string) (string, bool) {
	words, err := r.List()
	if err != nil || len(words) == 0 {
		return content, false
	}

	hit := false
	for _, w := range words {
		if w == "" {
			continue
		}
		if strings.Contains(content, w) {
			hit = true
			placeholder := strings.Repeat("*", len([]rune(w)))
			content = strings.ReplaceAll(content, w, placeholder)
		}
	}
	return content, hit
}
