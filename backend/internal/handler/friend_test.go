//go:build integration

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchUsers(t *testing.T) {
	truncateTables(t)
	createTestUser(t, "alice")
	createTestUser(t, "bob")
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/users/search?q=alice", nil, 3)
	c.Request.URL.RawQuery = "q=alice"
	friendHandler.SearchUsers(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var users []interface{}
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestSearchUsers_EmptyQuery(t *testing.T) {
	truncateTables(t)
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/users/search", nil, 1)
	friendHandler.SearchUsers(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty query, got %d", w.Code)
	}
}

func TestSendAndAcceptFriendRequest(t *testing.T) {
	truncateTables(t)
	aliceID, _ := createTestUser(t, "alice")
	bobID, _ := createTestUser(t, "bob")

	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/friends/request/1", nil, aliceID)
	c.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(bobID, 10)}}
	friendHandler.SendRequest(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on send request, got %d: %s", w.Code, w.Body.String())
	}

	c2, w2 := createTestContextWithUser(http.MethodGet, "/api/v1/friends/requests", nil, bobID)
	friendHandler.ListRequests(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
	var reqs []interface{}
	json.Unmarshal(w2.Body.Bytes(), &reqs)
	if len(reqs) != 1 {
		t.Errorf("expected 1 request, got %d", len(reqs))
	}

	c3, w3 := createTestContextWithUser(http.MethodPost, "/api/v1/friends/accept/1", nil, bobID)
	c3.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(aliceID, 10)}}
	friendHandler.AcceptRequest(c3)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200 on accept, got %d: %s", w3.Code, w3.Body.String())
	}

	c4, w4 := createTestContextWithUser(http.MethodGet, "/api/v1/friends", nil, aliceID)
	friendHandler.ListFriends(c4)
	if w4.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w4.Code, w4.Body.String())
	}
	var friends []interface{}
	json.Unmarshal(w4.Body.Bytes(), &friends)
	if len(friends) != 1 {
		t.Errorf("expected 1 friend, got %d", len(friends))
	}
}

func TestRejectFriendRequest(t *testing.T) {
	truncateTables(t)
	aliceID, _ := createTestUser(t, "alice2")
	bobID, _ := createTestUser(t, "bob2")

	c, _ := createTestContextWithUser(http.MethodPost, "/api/v1/friends/request/1", nil, aliceID)
	c.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(bobID, 10)}}
	friendHandler.SendRequest(c)

	c2, w2 := createTestContextWithUser(http.MethodPost, "/api/v1/friends/reject/1", nil, bobID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(aliceID, 10)}}
	friendHandler.RejectRequest(c2)
	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 on reject, got %d: %s", w2.Code, w2.Body.String())
	}

	c3, w3 := createTestContextWithUser(http.MethodGet, "/api/v1/friends", nil, aliceID)
	friendHandler.ListFriends(c3)
	var friends []interface{}
	json.Unmarshal(w3.Body.Bytes(), &friends)
	if len(friends) != 0 {
		t.Errorf("expected 0 friends after reject, got %d", len(friends))
	}
}

func TestRemoveFriend(t *testing.T) {
	truncateTables(t)
	aliceID, _ := createTestUser(t, "alice3")
	bobID, _ := createTestUser(t, "bob3")

	c, _ := createTestContextWithUser(http.MethodPost, "/api/v1/friends/request/1", nil, aliceID)
	c.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(bobID, 10)}}
	friendHandler.SendRequest(c)
	c2, _ := createTestContextWithUser(http.MethodPost, "/api/v1/friends/accept/1", nil, bobID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(aliceID, 10)}}
	friendHandler.AcceptRequest(c2)

	c3, w3 := createTestContextWithUser(http.MethodDelete, "/api/v1/friends/1", nil, aliceID)
	c3.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(bobID, 10)}}
	friendHandler.RemoveFriend(c3)
	if w3.Code != http.StatusOK {
		t.Errorf("expected 200 on remove, got %d: %s", w3.Code, w3.Body.String())
	}
}

func TestListNotifications(t *testing.T) {
	truncateTables(t)
	aliceID, _ := createTestUser(t, "alice_n")
	bobID, _ := createTestUser(t, "bob_n")

	c, _ := createTestContextWithUser(http.MethodPost, "/api/v1/friends/request/1", nil, aliceID)
	c.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(bobID, 10)}}
	friendHandler.SendRequest(c)

	c2, w2 := createTestContextWithUser(http.MethodGet, "/api/v1/notifications", nil, bobID)
	friendHandler.ListNotifications(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
	var notifs []interface{}
	json.Unmarshal(w2.Body.Bytes(), &notifs)
	if len(notifs) < 1 {
		t.Errorf("expected at least 1 notification, got %d", len(notifs))
	}
}
