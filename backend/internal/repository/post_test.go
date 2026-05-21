//go:build integration

package repository

import (
	"testing"

	"github.com/guish/university-social/internal/model"
)

func TestPostRepository_CreateAndFind(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	postRepo := NewPostRepository(testDB)

	userID := insertUser(t, userRepo, "p_cf_user")

	postID, err := postRepo.Create(userID, &model.CreatePostRequest{
		Content: "测试帖子内容",
	})
	if err != nil {
		t.Fatalf("Create post failed: %v", err)
	}
	if postID == 0 {
		t.Fatal("Expected non-zero post id")
	}

	post, err := postRepo.FindByID(postID, userID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if post == nil {
		t.Fatal("Expected post to exist")
	}
	if post.Content != "测试帖子内容" {
		t.Errorf("Content mismatch: '%s'", post.Content)
	}
}

func TestPostRepository_List(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	postRepo := NewPostRepository(testDB)

	userID := insertUser(t, userRepo, "p_list_user")

	_, err := postRepo.Create(userID, &model.CreatePostRequest{Content: "帖子 A"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	_, err = postRepo.Create(userID, &model.CreatePostRequest{Content: "帖子 B"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	posts, err := postRepo.List("", nil, 0, 10, userID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(posts) < 2 {
		t.Fatalf("Expected >=2 posts, got %d", len(posts))
	}
}

func TestPostRepository_Delete(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	postRepo := NewPostRepository(testDB)

	userID := insertUser(t, userRepo, "p_del_user")
	postID, err := postRepo.Create(userID, &model.CreatePostRequest{Content: "待删除"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = postRepo.Delete(postID, userID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = postRepo.FindByID(postID, userID)
	if err == nil {
		t.Fatal("Expected error after delete")
	}
}

func TestPostRepository_LikeToggle(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	postRepo := NewPostRepository(testDB)

	authorID := insertUser(t, userRepo, "p_like_auth")
	likerID := insertUser(t, userRepo, "p_like_liker")

	postID, err := postRepo.Create(authorID, &model.CreatePostRequest{Content: "点赞测试"})
	if err != nil {
		t.Fatalf("Create post failed: %v", err)
	}

	liked, err := postRepo.ToggleLike(likerID, postID)
	if err != nil {
		t.Fatalf("ToggleLike failed: %v", err)
	}
	if !liked {
		t.Fatal("Expected liked=true")
	}

	liked, err = postRepo.ToggleLike(likerID, postID)
	if err != nil {
		t.Fatalf("ToggleLike unlike failed: %v", err)
	}
	if liked {
		t.Fatal("Expected liked=false after unlike")
	}
}

func TestPostRepository_SearchPosts(t *testing.T) {
	truncateTables(t)
	userRepo := NewUserRepository(testDB)
	postRepo := NewPostRepository(testDB)

	userID := insertUser(t, userRepo, "p_search_user")

	postRepo.Create(userID, &model.CreatePostRequest{Content: "今天天气真不错"})
	postRepo.Create(userID, &model.CreatePostRequest{Content: "期末考试加油"})

	results, err := postRepo.SearchPosts("天气", 0, 10, userID)
	if err != nil {
		t.Fatalf("SearchPosts failed: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Expected results for '天气'")
	}
}
