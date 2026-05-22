//go:build integration

package repository

import (
	"testing"

	"github.com/guish/university-social/internal/model"
)

func insertCircle(t *testing.T, repo *CircleRepository, name string, creatorID int64) int64 {
	t.Helper()
	id, err := repo.Create(creatorID, &model.CreateCircleRequest{Name: name})
	if err != nil {
		t.Fatalf("Create circle failed: %v", err)
	}
	return id
}

func TestCircleRepository_CreateAndGet(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	circleRepo := NewCircleRepository(testDB)

	uid := insertUser(t, userRepo, "c1")
	cid := insertCircle(t, circleRepo, "测试圈子", uid)

	circle, err := circleRepo.GetByID(cid, uid)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if circle == nil {
		t.Fatal("Expected circle")
	}
	if circle.Name != "测试圈子" {
		t.Errorf("Expected '测试圈子', got '%s'", circle.Name)
	}
	if circle.JoinType != "free" {
		t.Errorf("Default join_type should be 'free', got '%s'", circle.JoinType)
	}
}

func TestCircleRepository_JoinAndLeave(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	circleRepo := NewCircleRepository(testDB)

	creatorID := insertUser(t, userRepo, "c2c")
	memberID := insertUser(t, userRepo, "c2m")
	cid := insertCircle(t, circleRepo, "圈子", creatorID)

	if err := circleRepo.Join(cid, memberID); err != nil {
		t.Fatalf("Join failed: %v", err)
	}

	members, err := circleRepo.ListMembers(int(cid))
	if err != nil {
		t.Fatalf("ListMembers failed: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("Expected 2 members (creator+joiner), got %d", len(members))
	}

	if err := circleRepo.Leave(cid, memberID); err != nil {
		t.Fatalf("Leave failed: %v", err)
	}

	members, _ = circleRepo.ListMembers(int(cid))
	if len(members) != 1 {
		t.Errorf("Expected 1 member (creator) after leave, got %d", len(members))
	}
}

func TestCircleRepository_JoinRequestFlow(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	circleRepo := NewCircleRepository(testDB)

	creatorID := insertUser(t, userRepo, "c3c")
	memberID := insertUser(t, userRepo, "c3m")
	cid := insertCircle(t, circleRepo, "审核圈子", creatorID)

	if err := circleRepo.CreateJoinRequest(cid, memberID); err != nil {
		t.Fatalf("CreateJoinRequest failed: %v", err)
	}

	requests, err := circleRepo.ListJoinRequests(cid)
	if err != nil {
		t.Fatalf("ListJoinRequests failed: %v", err)
	}
	if len(requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requests))
	}
	if requests[0].UserID != memberID {
		t.Errorf("Request user_id mismatch: %d vs %d", requests[0].UserID, memberID)
	}

	if err := circleRepo.HandleJoinRequest(requests[0].ID, creatorID, "approve"); err != nil {
		t.Fatalf("HandleJoinRequest failed: %v", err)
	}

	members, _ := circleRepo.ListMembers(int(cid))
	if len(members) != 1 {
		t.Errorf("Expected 1 member after approval, got %d", len(members))
	}
}

func TestCircleRepository_Activities(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	circleRepo := NewCircleRepository(testDB)

	uid := insertUser(t, userRepo, "c4")
	cid := insertCircle(t, circleRepo, "动态测试", uid)

	if err := circleRepo.CreateActivity(cid, uid, "join", ""); err != nil {
		t.Fatalf("CreateActivity join failed: %v", err)
	}
	if err := circleRepo.CreateActivity(cid, uid, "create_post", "测试内容"); err != nil {
		t.Fatalf("CreateActivity post failed: %v", err)
	}

	acts, err := circleRepo.ListActivities(cid, 50)
	if err != nil {
		t.Fatalf("ListActivities failed: %v", err)
	}
	if len(acts) != 2 {
		t.Fatalf("Expected 2 activities, got %d", len(acts))
	}
	if acts[0].Action != "create_post" {
		t.Errorf("First (latest) should be create_post, got %s", acts[0].Action)
	}
	if acts[1].Action != "join" {
		t.Errorf("Second should be join, got %s", acts[1].Action)
	}
}
