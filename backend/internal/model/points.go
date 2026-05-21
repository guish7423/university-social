package model

import "time"

type PointsLog struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Action      string    `json:"action"`
	Points      int       `json:"points"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type PointsBalance struct {
	Points int    `json:"points"`
	Level  string `json:"level"`
	LevelID int   `json:"level_id"`
	NextLevel string `json:"next_level,omitempty"`
	NextThreshold int `json:"next_threshold,omitempty"`
}

type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	School   string `json:"school"`
	Avatar   string `json:"avatar"`
	Points   int    `json:"points"`
	Level    string `json:"level"`
	Rank     int    `json:"rank"`
}

var LevelThresholds = []struct {
	Min  int
	Name string
}{
	{0, "Lv.1 萌新"},
	{100, "Lv.2 活跃"},
	{300, "Lv.3 资深"},
	{1000, "Lv.4 大佬"},
	{3000, "Lv.5 传奇"},
}

func CalculateLevel(points int) (levelName string, levelID int, nextName string, nextThreshold int) {
	for i, lt := range LevelThresholds {
		if i == len(LevelThresholds)-1 {
			return lt.Name, i + 1, "", 0
		}
		if points < LevelThresholds[i+1].Min {
			return lt.Name, i + 1, LevelThresholds[i+1].Name, LevelThresholds[i+1].Min
		}
	}
	return LevelThresholds[len(LevelThresholds)-1].Name, len(LevelThresholds), "", 0
}
