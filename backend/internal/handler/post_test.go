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

func TestCreatePost_Success(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "post_creator")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/posts", model.CreatePostRequest{
		Content: "测试帖子内容",
	}, userID)
	postHandler.CreatePost(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID == 0 {
		t.Error("expected non-zero post ID")
	}
}

func TestCreatePost_EmptyContent(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "empty_post")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/posts", model.CreatePostRequest{}, userID)
	postHandler.CreatePost(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListPosts_Empty(t *testing.T) {
	truncateTables(t)
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/posts", nil, 1)
	postHandler.ListPosts(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var posts []*model.Post
	json.Unmarshal(w.Body.Bytes(), &posts)
	if len(posts) != 0 {
		t.Errorf("expected empty list, got %d items", len(posts))
	}
}

func TestListPosts_WithPosts(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "list_test")
	for i := 0; i < 3; i++ {
		c, _ := createTestContextWithUser(
			http.MethodPost, "/api/v1/posts",
			model.CreatePostRequest{Content: "帖子" + strconv.Itoa(i)},
			userID,
		)
		postHandler.CreatePost(c)
	}
	c, w := createTestContextWithUser(http.MethodGet, "/api/v1/posts", nil, userID)
	postHandler.ListPosts(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var posts []*model.Post
	json.Unmarshal(w.Body.Bytes(), &posts)
	if len(posts) != 3 {
		t.Errorf("expected 3 posts, got %d", len(posts))
	}
}

func TestGetPost_Success(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "get_post")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/posts",
		model.CreatePostRequest{Content: "获取测试"}, userID)
	postHandler.CreatePost(c)
	var created struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &created)

	c2, w2 := createTestContextWithUser(http.MethodGet, "/api/v1/posts/1", nil, userID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	postHandler.GetPost(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
	var post model.Post
	json.Unmarshal(w2.Body.Bytes(), &post)
	if post.Content != "获取测试" {
		t.Errorf("expected content '获取测试', got %s", post.Content)
	}
}

func TestDeletePost_Owner(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "del_owner")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/posts",
		model.CreatePostRequest{Content: "删除测试"}, userID)
	postHandler.CreatePost(c)
	var created struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &created)

	c2, w2 := createTestContextWithUser(http.MethodDelete, "/api/v1/posts/1", nil, userID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	postHandler.DeletePost(c2)
	if w2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestDeletePost_NonOwner(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "owner")
	otherID, _ := createTestUser(t, "other")

	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/posts",
		model.CreatePostRequest{Content: "他人帖子"}, userID)
	postHandler.CreatePost(c)
	var created struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &created)

	c2, w2 := createTestContextWithUser(http.MethodDelete, "/api/v1/posts/1", nil, otherID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	postHandler.DeletePost(c2)
	if w2.Code != http.StatusForbidden {
		t.Errorf("expected 403 for non-owner, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestCreateComment_Success(t *testing.T) {
	truncateTables(t)
	userID, _ := createTestUser(t, "commenter")
	c, w := createTestContextWithUser(http.MethodPost, "/api/v1/posts",
		model.CreatePostRequest{Content: "评论测试帖"}, userID)
	postHandler.CreatePost(c)
	var created struct{ ID int64 }
	json.Unmarshal(w.Body.Bytes(), &created)

	c2, w2 := createTestContextWithUser(http.MethodPost, "/api/v1/posts/1/comments",
		model.CommentRequest{Content: "测试评论"}, userID)
	c2.Params = []gin.Param{{Key: "id", Value: strconv.FormatInt(created.ID, 10)}}
	postHandler.CreateComment(c2)
	if w2.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w2.Code, w2.Body.String())
	}
}
