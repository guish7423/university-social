//go:build integration

package repository

import (
	"testing"

	"github.com/guish/university-social/internal/model"
)

func insertUser(t *testing.T, repo *UserRepository, suffix string) int64 {
	t.Helper()
	id, err := repo.Create(&model.User{OpenID: "test_" + suffix, Nickname: "用户" + suffix})
	if err != nil {
		t.Fatalf("Create user failed: %v", err)
	}
	return id
}

func TestUserRepository_CreateAndFind(t *testing.T) {
	truncateTables(t)
	repo := NewUserRepository(testDB)

	id := insertUser(t, repo, "cf")

	user, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user")
	}
	if user.Nickname != "用户cf" {
		t.Errorf("Expected '用户cf', got '%s'", user.Nickname)
	}
	if user.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", user.Role)
	}
}

func TestUserRepository_FindByOpenID(t *testing.T) {
	truncateTables(t)
	repo := NewUserRepository(testDB)

	id := insertUser(t, repo, "fbid")

	found, err := repo.FindByOpenID("test_fbid")
	if err != nil {
		t.Fatalf("FindByOpenID failed: %v", err)
	}
	if found == nil {
		t.Fatal("Expected user")
	}
	if found.ID != id {
		t.Errorf("ID mismatch: %d vs %d", found.ID, id)
	}

	notFound, err := repo.FindByOpenID("nonexistent")
	if err != nil {
		t.Fatalf("FindByOpenID failed: %v", err)
	}
	if notFound != nil {
		t.Fatal("Expected nil for nonexistent")
	}
}

func TestUserRepository_UpdateProfile(t *testing.T) {
	truncateTables(t)
	repo := NewUserRepository(testDB)

	id := insertUser(t, repo, "up")

	err := repo.UpdateProfile(id, "新昵称", "https://example.com/ava.png", "测试大学")
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}

	updated, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.Nickname != "新昵称" {
		t.Errorf("Expected '新昵称', got '%s'", updated.Nickname)
	}
	if updated.School != "测试大学" {
		t.Errorf("Expected '测试大学', got '%s'", updated.School)
	}
}

func TestUserRepository_NotFound(t *testing.T) {
	truncateTables(t)
	repo := NewUserRepository(testDB)

	user, err := repo.FindByID(999999)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if user != nil {
		t.Fatal("Expected nil for nonexistent")
	}
}
