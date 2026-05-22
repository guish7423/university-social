package repository

import (
	"database/sql"

	"github.com/guish/university-social/internal/model"
)

type CampusRepository struct {
	db *sql.DB
}

func NewCampusRepository(db *sql.DB) *CampusRepository {
	return &CampusRepository{db: db}
}

func (r *CampusRepository) ListCalendar() ([]*model.CampusCalendar, error) {
	rows, err := r.db.Query(
		`SELECT id, title, event_date, event_type, COALESCE(description,''), created_at
		 FROM campus_calendar ORDER BY event_date ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*model.CampusCalendar
	for rows.Next() {
		var c model.CampusCalendar
		if err := rows.Scan(&c.ID, &c.Title, &c.EventDate, &c.EventType, &c.Description, &c.CreatedAt); err != nil {
			continue
		}
		items = append(items, &c)
	}
	if items == nil {
		items = []*model.CampusCalendar{}
	}
	return items, nil
}

func (r *CampusRepository) ListDirectory(dept string) ([]*model.CampusDirectory, error) {
	var rows *sql.Rows
	var err error
	if dept != "" {
		rows, err = r.db.Query(
			`SELECT id, name, department, COALESCE(position,''), COALESCE(phone,''),
			        COALESCE(email,''), COALESCE(office,''), category, created_at
			 FROM campus_directory WHERE department ILIKE $1 ORDER BY department, name`,
			"%"+dept+"%",
		)
	} else {
		rows, err = r.db.Query(
			`SELECT id, name, department, COALESCE(position,''), COALESCE(phone,''),
			        COALESCE(email,''), COALESCE(office,''), category, created_at
			 FROM campus_directory ORDER BY department, name`,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*model.CampusDirectory
	for rows.Next() {
		var c model.CampusDirectory
		if err := rows.Scan(&c.ID, &c.Name, &c.Department, &c.Position, &c.Phone, &c.Email, &c.Office, &c.Category, &c.CreatedAt); err != nil {
			continue
		}
		items = append(items, &c)
	}
	if items == nil {
		items = []*model.CampusDirectory{}
	}
	return items, nil
}

func (r *CampusRepository) ListEmptyRooms(building string, weekday int) ([]*model.EmptyClassroom, error) {
	var rows *sql.Rows
	var err error
	if building != "" && weekday > 0 {
		rows, err = r.db.Query(
			`SELECT id, building, room, capacity, weekday, start_time, end_time, semester, created_at
			 FROM empty_classrooms WHERE building = $1 AND weekday = $2 ORDER BY start_time`,
			building, weekday,
		)
	} else if building != "" {
		rows, err = r.db.Query(
			`SELECT id, building, room, capacity, weekday, start_time, end_time, semester, created_at
			 FROM empty_classrooms WHERE building = $1 ORDER BY weekday, start_time`,
			building,
		)
	} else {
		rows, err = r.db.Query(
			`SELECT id, building, room, capacity, weekday, start_time, end_time, semester, created_at
			 FROM empty_classrooms ORDER BY weekday, start_time`,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*model.EmptyClassroom
	for rows.Next() {
		var c model.EmptyClassroom
		var startStr, endStr string
		if err := rows.Scan(&c.ID, &c.Building, &c.Room, &c.Capacity, &c.Weekday, &startStr, &endStr, &c.Semester, &c.CreatedAt); err != nil {
			continue
		}
		c.StartTime = startStr[:5]
		c.EndTime = endStr[:5]
		items = append(items, &c)
	}
	if items == nil {
		items = []*model.EmptyClassroom{}
	}
	return items, nil
}

func (r *CampusRepository) ListBuildings() ([]string, error) {
	rows, err := r.db.Query(`SELECT DISTINCT building FROM empty_classrooms ORDER BY building`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var buildings []string
	for rows.Next() {
		var b string
		rows.Scan(&b)
		buildings = append(buildings, b)
	}
	if buildings == nil {
		buildings = []string{}
	}
	return buildings, nil
}
