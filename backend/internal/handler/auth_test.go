//go:build integration

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guish/university-social/internal/model"
)

func createTestUser(t *testing.T, nickname string) (int64, string) {
	t.Helper()
	c, w := createTestContext(http.MethodPost, "/api/v1/auth/dev-login", map[string]string{
		"nickname": nickname,
	})
	authHandler.DevLogin(c)
	var resp model.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("DevLogin unmarshal: %v", err)
	}
	return resp.User.ID, resp.Token
}

func TestDevLogin_CreatesNewUser(t *testing.T) {
	truncateTables(t)
	c, w := createTestContext(http.MethodPost, "/api/v1/auth/dev-login", map[string]string{
		"nickname": "test_dev",
	})
	authHandler.DevLogin(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp model.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.User.Nickname != "test_dev" {
		t.Errorf("expected nickname test_dev, got %s", resp.User.Nickname)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestDevLogin_ReturnsExistingUser(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "existing")
	c, w := createTestContext(http.MethodPost, "/api/v1/auth/dev-login", map[string]string{
		"nickname": "existing",
	})
	authHandler.DevLogin(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp model.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.User.ID != userID {
		t.Errorf("expected same user ID %d, got %d", userID, resp.User.ID)
	}
}

func TestDevLogin_DefaultNickname(t *testing.T) {
	truncateTables(t)
	c, w := createTestContext(http.MethodPost, "/api/v1/auth/dev-login", map[string]string{})
	authHandler.DevLogin(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp model.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.User.Nickname != "测试用户" {
		t.Errorf("expected default nickname 测试用户, got %s", resp.User.Nickname)
	}
}

func TestGetProfile(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "profile_test")
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/user/profile", nil, userID)
	authHandler.GetProfile(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var user model.User
	json.Unmarshal(w.Body.Bytes(), &user)
	if user.ID != userID {
		t.Errorf("expected user ID %d, got %d", userID, user.ID)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	truncateTables(t)
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/user/profile", nil, 99999)
	authHandler.GetProfile(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateProfile(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "update_test")
	c, w := createTestContextWithUser(http.MethodPut, "/api/v1/user/profile", map[string]string{
		"nickname": "新昵称",
		"school":   "测试大学",
	}, userID)
	authHandler.UpdateProfile(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	c2, w2 := createTestContextWithUser(http.MethodGet, "/api/v1/user/profile", nil, userID)
	authHandler.GetProfile(c2)
	var user model.User
	json.Unmarshal(w2.Body.Bytes(), &user)
	if user.Nickname != "新昵称" {
		t.Errorf("expected nickname '新昵称', got %s", user.Nickname)
	}
	if user.School != "测试大学" {
		t.Errorf("expected school '测试大学', got %s", user.School)
	}
}

func TestGetUser(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "getuser_test")
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/users/1", nil, 1)
	c.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(99999, 10)}}
	authHandler.GetUser(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
	c2, w2 := createTestContextWithUser(http.MethodGet, "/api/v1/users/1", nil, 1)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(userID, 10)}}
	authHandler.GetUser(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}
