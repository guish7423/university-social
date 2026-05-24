# Post Editing (帖子编辑)

Allow post authors to edit their own posts after publishing.

## Discovery

No clarifying questions needed — scope is well understood from existing patterns.

**Research:**
- Post model has `Content`, `Images`, `TopicID` fields that can be updated
- Existing `DeletePost` handler verifies ownership via `user_id` check in SQL
- `CreatePostRequest` struct is the pattern for the new `UpdatePostRequest`
- Frontend `detail.vue` already shows "删除" button for own posts — add "编辑" next to it
- API client `post.ts` follows `request.put/get/post/delete` pattern

## Non-Goals
- No version history or edit tracking
- No editing circle posts (separate handler exists)
- No moderation edit flag

## Design Summary

Simple inline edit on the detail page. Author clicks "编辑", content becomes textarea, images become editable (add/remove), save sends PUT to `/posts/:id`. Backend updates `content`, `images`, `topic_id`, and bumps `updated_at`.

## Tasks

### 1. Backend — Repository Update

**Depends on**: none

**Files**: 
- Modify: `backend/internal/model/post.go` — add `UpdatePostRequest`
- Modify: `backend/internal/repository/post.go` — add `Update()` method

**What**:
1. Add `UpdatePostRequest` struct with `Content` (required), `Images`, `TopicID`
2. Add `Update(ctx, id, userID, req)` method that runs `UPDATE posts SET content=$1, images=$2, topic_id=$3, updated_at=NOW() WHERE id=$4 AND user_id=$5`

**Must NOT**:
- Allow updating other users' posts (verified by `user_id=$5` in WHERE clause)
- Expose `IsFeatured`, `LikeCount`, etc as updatable fields

**References**: `post.go:136` (Delete pattern), `post.go:19` (Create pattern)

**Verify**: `go build ./...` in backend/

### 2. Backend — Handler + Route

**Depends on**: 1

**Files**:
- Modify: `backend/internal/handler/post.go` — add `UpdatePost` handler
- Modify: `backend/internal/router/router.go` — add PUT route

**What**:
1. Add `UpdatePost(c *gin.Context)` handler: parse post ID from param, bind JSON to `UpdatePostRequest`, call `repo.Update()`, return OK
2. Add route: `auth.PUT("/posts/:id", postHandler.UpdatePost)` after `auth.GET("/posts/:id", postHandler.GetPost)`

**Verify**: `go build ./...` in backend/

### 3. Frontend — API + Edit Mode

**Depends on**: 2

**Files**:
- Modify: `pc/src/api/post.ts` — add `updatePost` function
- Modify: `pc/src/pages/post/detail.vue` — add edit mode toggle

**What**:
1. Add `updatePost(id, data)` API function using `request.put`
2. Add `editing` ref, `editForm` reactive, `handleEdit`/`handleSave`/`handleCancel` functions
3. Show "编辑" button next to "删除" for own posts
4. In edit mode: content becomes textarea, images show add/remove controls
5. Save calls API, updates local post data, exits edit mode

**Verify**: `vue-tsc --noEmit && vite build` in pc/
