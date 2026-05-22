//go:build integration

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateCircle_Success(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "circle_creator")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/circles", map[string]interface{}{
		"name":        "测试圈子",
		"description": "圈子描述",
		"avatar":      "https://example.com/avatar.png",
	}, userID)
	circleHandler.Create(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID == 0 {
		t.Error("expected non-zero circle ID")
	}
}

func TestCreateCircle_NoName(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "no_name")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/circles", map[string]interface{}{
		"description": "无名称圈子",
	}, userID)
	circleHandler.Create(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetCircleByID(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "circle_owner")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/circles", map[string]interface{}{
		"name": "查询测试圈",
	}, userID)
	circleHandler.Create(c)
	var created struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &created)

	c2, w2 := createTestContextWithUser(http.MethodGet, "/api/v1/circles/1", nil, userID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	circleHandler.GetByID(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	if resp["name"] != "查询测试圈" {
		t.Errorf("expected name '查询测试圈', got %v", resp["name"])
	}
}

func TestJoinLeaveCircle(t *testing.T) {
	truncateTables(t)
	ownerID, _ := createTestUser(t, "owner")
	memberID, _ := createTestUser(t, "member")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/circles", map[string]interface{}{
		"name": "加入测试圈",
	}, ownerID)
	circleHandler.Create(c)
	var created struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &created)

	c2, w2 := createTestContextWithUser(http.MethodPost, "/api/v1/circles/1/join", nil, memberID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	circleHandler.Join(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 on join, got %d: %s", w2.Code, w2.Body.String())
	}

	c3, w3 := createTestContextWithUser(http.MethodPost, "/api/v1/circles/1/leave", nil, memberID)
	c3.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	circleHandler.Leave(c3)
	if w3.Code != http.StatusOK {
		t.Errorf("expected 200 on leave, got %d: %s", w3.Code, w3.Body.String())
	}
}

func TestListCircles(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "list_circles")
	for i := 0; i < 3; i++ {
		c, _ := createTestContextWithUser(http.MethodPost, "/api/v1/circles",
			map[string]interface{}{"name": "圈子" + strconv.Itoa(i)},
			userID)
		circleHandler.Create(c)
	}
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/circles", nil, userID)
	circleHandler.List(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var circles []interface{}
	json.Unmarshal(w.Body.Bytes(), &circles)
	if len(circles) < 3 {
		t.Errorf("expected at least 3 circles, got %d", len(circles))
	}
}
