package repository

import (
	"database/sql"

	"github.com/guish/university-social/internal/model"
)

type UniversityRepository struct {
	db *sql.DB
}

func NewUniversityRepository(db *sql.DB) *UniversityRepository {
	return &UniversityRepository{db: db}
}

func (r *UniversityRepository) List(userID int64) ([]*model.University, error) {
	query := `SELECT u.id, u.name, u.domain, u.logo_url, u.created_at,
		(SELECT COUNT(*) FROM university_members WHERE university_id = u.id) as member_count,
		EXISTS(SELECT 1 FROM university_members WHERE university_id = u.id AND user_id = $1) as is_member
		FROM universities u ORDER BY member_count DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var universities []*model.University
	for rows.Next() {
		var u model.University
		if err := rows.Scan(&u.ID, &u.Name, &u.Domain, &u.LogoURL, &u.CreatedAt, &u.MemberCount, &u.IsMember); err != nil {
			return nil, err
		}
		universities = append(universities, &u)
	}
	if universities == nil {
		universities = []*model.University{}
	}
	return universities, nil
}

func (r *UniversityRepository) GetByID(id, userID int64) (*model.University, error) {
	query := `SELECT u.id, u.name, u.domain, u.logo_url, u.created_at,
		(SELECT COUNT(*) FROM university_members WHERE university_id = u.id) as member_count,
		EXISTS(SELECT 1 FROM university_members WHERE university_id = u.id AND user_id = $1) as is_member
		FROM universities u WHERE u.id = $2`
	var u model.University
	err := r.db.QueryRow(query, userID, id).Scan(&u.ID, &u.Name, &u.Domain, &u.LogoURL, &u.CreatedAt, &u.MemberCount, &u.IsMember)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UniversityRepository) Join(userID, universityID int64) error {
	_, err := r.db.Exec(`INSERT INTO university_members (user_id, university_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, universityID)
	return err
}

func (r *UniversityRepository) Leave(userID, universityID int64) error {
	_, err := r.db.Exec(`DELETE FROM university_members WHERE user_id = $1 AND university_id = $2`, userID, universityID)
	return err
}

func (r *UniversityRepository) GetMembers(universityID int64) ([]*model.UniversityMember, error) {
	query := `SELECT um.id, um.user_id, um.university_id, um.is_admin, um.created_at,
		u.id, u.nickname, u.avatar, u.school
		FROM university_members um
		JOIN users u ON u.id = um.user_id
		WHERE um.university_id = $1
		ORDER BY um.is_admin DESC, um.created_at ASC`
	rows, err := r.db.Query(query, universityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*model.UniversityMember
	for rows.Next() {
		var m model.UniversityMember
		var u model.User
		if err := rows.Scan(&m.ID, &m.UserID, &m.UniversityID, &m.IsAdmin, &m.CreatedAt,
			&u.ID, &u.Nickname, &u.Avatar, &u.School); err != nil {
			return nil, err
		}
		m.User = &u
		members = append(members, &m)
	}
	if members == nil {
		members = []*model.UniversityMember{}
	}
	return members, nil
}

func (r *UniversityRepository) GetUserUniversity(userID int64) (*model.University, error) {
	query := `SELECT u.id, u.name, u.domain, u.logo_url, u.created_at,
		(SELECT COUNT(*) FROM university_members WHERE university_id = u.id) as member_count,
		true as is_member
		FROM universities u
		JOIN university_members um ON um.university_id = u.id
		WHERE um.user_id = $1`
	var u model.University
	err := r.db.QueryRow(query, userID).Scan(&u.ID, &u.Name, &u.Domain, &u.LogoURL, &u.CreatedAt, &u.MemberCount, &u.IsMember)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
