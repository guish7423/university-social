//go:build integration

package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"

	"github.com/guish/university-social/internal/config"
	"github.com/guish/university-social/internal/migration"
	"github.com/guish/university-social/internal/repository"
)

var testDB *sql.DB
var testCfg *config.Config
var (
	authHandler   *AuthHandler
	postHandler   *PostHandler
	circleHandler *CircleHandler
	friendHandler *FriendHandler
)
var (
	userRepo   *repository.UserRepository
	postRepo   *repository.PostRepository
	circleRepo *repository.CircleRepository
	friendRepo *repository.FriendRepository
	sensitiveRepo *repository.SensitiveWordRepository
	pointsRepo *repository.PointsRepository
)

func getTestDSN(dbName string) string {
	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		base = "postgres://postgres:postgres@localhost:5432/%s?sslmode=disable"
	}
	return fmt.Sprintf(base, dbName)
}

func ensureTestDB(dbName string) error {
	adminDSN := getTestDSN("postgres")
	db, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return fmt.Errorf("connect admin db: %w", err)
	}
	defer db.Close()

	var exists bool
	db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("create test db: %w", err)
		}
		log.Printf("Created test database: %s", dbName)
	}
	return nil
}

func findMigrationsDir() string {
	candidates := []string{"../migrations", "../../migrations", "migrations"}
	for _, d := range candidates {
		abs, err := filepath.Abs(d)
		if err != nil {
			continue
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			return abs
		}
	}
	return "../migrations"
}

func truncateTables(t *testing.T) {
	t.Helper()
	tables := []string{
		"likes", "comments", "posts", "notifications",
		"friends", "follows", "circle_members", "circle_post_comments",
		"messages", "whisper_comments", "whispers",
		"product_comments", "products", "activity_comments", "activities",
		"course_comments", "courses", "lost_item_comments",
		"reports", "points_logs", "verification_requests",
		"invite_codes", "circle_posts", "circles", "users",
	}
	for _, tbl := range tables {
		testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tbl))
	}
}

func createTestContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c.Request = req
	return c, w
}

func createTestContextWithUser(method, path string, body interface{}, userID int64) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := createTestContext(method, path, body)
	c.Set("user_id", userID)
	return c, w
}

func TestMain(m *testing.M) {
	dbName := "university_social_test"
	if err := ensureTestDB(dbName); err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}

	var err error
	testDB, err = sql.Open("postgres", getTestDSN(dbName))
	if err != nil {
		log.Fatalf("Failed to connect test db: %v", err)
	}
	if err := testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping test db: %v", err)
	}

	migDir := findMigrationsDir()
	log.Printf("Using migrations: %s", migDir)
	if err := migration.Run(testDB, migDir); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	testCfg = &config.Config{
		JWTSecret: "test-secret",
	}

	userRepo = repository.NewUserRepository(testDB)
	postRepo = repository.NewPostRepository(testDB)
	circleRepo = repository.NewCircleRepository(testDB)
	friendRepo = repository.NewFriendRepository(testDB)
	sensitiveRepo = repository.NewSensitiveWordRepository(testDB)
	pointsRepo = repository.NewPointsRepository(testDB)

	authHandler = NewAuthHandler(userRepo, testCfg, pointsRepo)
	postHandler = NewPostHandler(postRepo, sensitiveRepo, pointsRepo)
	circleHandler = NewCircleHandler(circleRepo)
	friendHandler = NewFriendHandler(friendRepo)

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}
