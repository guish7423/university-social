<template>
  <view class="page">
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="circle">
      <view class="circle-header">
        <view class="circle-header-bg" />
        <view class="circle-header-content">
          <view class="circle-header-top">
            <view class="circle-avatar">
              <u-avatar :text="circle.icon || circle.name[0] || '?'" size="120" shape="square"
                bg-color="linear-gradient(135deg, #C67A6A, #1E2A3A)" fontSize="48" />
            </view>
            <view class="circle-info">
              <text class="circle-name">{{ circle.name }}</text>
              <text class="circle-desc">{{ circle.description }}</text>
              <view class="circle-stats">
                <text class="circle-stat">{{ circle.member_count }} 成员</text>
                <text class="circle-stat">{{ circle.post_count }} 帖子</text>
                <text v-if="circle.join_type === 'approve'" class="circle-stat join-type">需审核</text>
              </view>
            </view>
          </view>
          <view class="circle-actions">
            <u-button v-if="!circle.is_member" type="primary" shape="circle" size="medium" @click="handleJoin">
              加入圈子
            </u-button>
            <template v-else>
              <u-button type="primary" shape="circle" size="medium" plain @click="handleLeave">
                退出圈子
              </u-button>
              <u-button type="primary" shape="circle" size="medium" @click="showCreatePost = true">
                发帖
              </u-button>
            </template>
          </view>
        </view>
      </view>

      <view class="tab-bar">
        <view v-for="tab in tabs" :key="tab.key"
          :class="['tab-item', activeTab === tab.key && 'tab-active']"
          @click="activeTab = tab.key">
          <text class="tab-text">{{ tab.label }}</text>
          <view v-if="activeTab === tab.key" class="tab-indicator" />
        </view>
        <view v-if="showReviewBtn" class="review-btn" @click="openJoinReview">
          <text class="review-btn-text">审核</text>
        </view>
      </view>

      <view v-if="activeTab === 'posts'" class="posts-section">
        <view v-if="posts.length === 0" class="empty-state">
          <u-icon name="empty-address" size="120" color="#EAE6E0" />
          <text class="empty-text">还没有帖子</text>
          <u-button type="primary" shape="circle" size="small" plain @click="showCreatePost = true">
            发布第一个帖子
          </u-button>
        </view>
        <view v-for="p in posts" :key="p.id" class="post-card" @click="showComments(p)">
          <view class="post-card-header">
            <u-avatar :src="p.author?.avatar" :text="p.author?.nickname?.[0] || '?'" size="48" shape="circle" />
            <view class="post-card-author">
              <text class="post-card-name">{{ p.author?.nickname || '匿名' }}</text>
              <text class="post-card-time">{{ formatTime(p.created_at) }}</text>
            </view>
          </view>
          <text class="post-card-content">{{ p.content }}</text>
          <view class="post-card-footer">
            <view :class="['post-action', p.is_liked && 'liked']" @click.stop="toggleLike(p)">
              <text :class="['action-icon', p.is_liked && 'liked']">{{ p.is_liked ? '❤' : '♡' }}</text>
              <text :class="['action-text', p.is_liked && 'liked']">{{ p.like_count || '' }}</text>
            </view>
            <view class="post-action">
              <text class="action-icon">💬</text>
              <text class="action-text">{{ p.comment_count || '' }}</text>
            </view>
          </view>
        </view>
      </view>

      <view v-if="activeTab === 'members'" class="members-section">
        <view v-for="m in members" :key="m.id" class="member-card">
          <u-avatar :src="m.user?.avatar" :text="m.user?.nickname?.[0] || '?'" size="56" shape="circle" />
          <view class="member-info">
            <text class="member-name">{{ m.user?.nickname || '匿名' }}</text>
            <text class="member-bio">{{ m.user?.school || '' }}</text>
          </view>
          <text v-if="m.role === 'creator'" class="role-badge creator">创建者</text>
          <text v-else class="role-badge">成员</text>
        </view>
      </view>
      <view v-if="activeTab === 'activities'" class="activities-section">
        <view v-if="activities.length === 0" class="empty-state">
          <u-icon name="info-circle" size="100" color="#EAE6E0" />
          <text class="empty-text">暂无动态</text>
        </view>
        <view v-for="a in activities" :key="a.id" class="activity-item">
          <u-avatar :src="a.user?.avatar" :text="a.user?.nickname?.[0] || '?'" size="48" shape="circle" />
          <view class="activity-body">
            <text class="activity-text">
              <text class="activity-name">{{ a.user?.nickname || '匿名' }}</text>
              <text v-if="a.action === 'join'"> 加入了圈子</text>
              <text v-else-if="a.action === 'create_post'"> 发布了帖子</text>
              <text v-else-if="a.action === 'approve_member'"> 通过了成员申请</text>
              <text v-else> {{ a.action }}</text>
            </text>
            <text v-if="a.action === 'create_post' && a.content" class="activity-preview">{{ a.content.slice(0, 50) }}{{ a.content.length > 50 ? '...' : '' }}</text>
            <text class="activity-time">{{ formatTime(a.created_at) }}</text>
          </view>
        </view>
      </view>
    </template>

    <view v-else class="error-state">
      <u-icon name="info-circle" size="100" color="#B8C2CE" />
      <text class="error-text">圈子不存在</text>
    </view>

    <u-safe-bottom />

    <u-popup :show="showCreatePost" mode="bottom" @close="showCreatePost = false" closeable round="16">
      <view class="create-post-panel">
        <text class="create-post-title">发帖</text>
        <textarea v-model="newPostContent" class="create-post-input"
          placeholder="说点什么..." maxlength="500" />
        <u-button type="primary" shape="circle" :disabled="!newPostContent.trim()"
          @click="handleCreatePost">发布</u-button>
      </view>
    </u-popup>

    <u-popup :show="showJoinReview" mode="bottom" @close="showJoinReview = false" closeable round="16">
      <view class="review-panel">
        <text class="review-title">加入申请</text>
        <view v-if="joinRequests.length === 0" class="review-empty">暂无待审核申请</view>
        <view v-for="r in joinRequests" :key="r.id" class="review-item">
          <u-avatar :src="r.user?.avatar" :text="r.user?.nickname?.[0] || '?'" size="56" shape="circle" />
          <view class="review-item-info">
            <text class="review-item-name">{{ r.user?.nickname || '匿名' }}</text>
            <text class="review-item-time">{{ formatTime(r.created_at) }}</text>
          </view>
          <view class="review-item-actions">
            <u-button size="small" type="primary" @click="handleJoinRequestAction(r.id, 'approve')">通过</u-button>
            <u-button size="small" type="warning" plain @click="handleJoinRequestAction(r.id, 'reject')">拒绝</u-button>
          </view>
        </view>
      </view>
    </u-popup>

    <u-popup :show="showCommentsPanel" mode="bottom" @close="showCommentsPanel = false" closeable round="16" :safe-area-inset-bottom="true">
      <view class="comments-panel">
        <view v-if="selectedPost" class="comments-original">
          <text class="comments-original-author">{{ selectedPost.author?.nickname }}: </text>
          <text class="comments-original-content">{{ selectedPost.content }}</text>
        </view>
        <view class="comments-list">
          <view v-for="c in comments" :key="c.id" class="comment-item">
            <text class="comment-author">{{ c.author?.nickname || '匿名' }}:</text>
            <text class="comment-content">{{ c.content }}</text>
          </view>
          <view v-if="comments.length === 0" class="comments-empty">暂无评论</view>
        </view>
        <view class="comments-input-bar">
          <input v-model="newComment" class="comments-input" placeholder="写评论..." />
          <text class="comments-send" @click="handleCreateComment">发送</text>
        </view>
      </view>
    </u-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { getCircle, joinCircle, leaveCircle, listCircleMembers,
         listCirclePosts, createCirclePost, toggleCirclePostLike,
         listCirclePostComments, createCirclePostComment,
         getCircleJoinRequests, handleJoinRequest,
         listCircleActivities, type CircleActivityData } from "@/api/circle"
import type { CircleData, CircleMemberData, CirclePostData } from "@/api/circle"
import type { UserInfo } from "@/api/user"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const circle = ref<CircleData>()
const members = ref<CircleMemberData[]>([])
const posts = ref<CirclePostData[]>([])
const activities = ref<CircleActivityData[]>([])
const loading = ref(true)
const activeTab = ref("posts")
const showCreatePost = ref(false)
const newPostContent = ref("")
const showCommentsPanel = ref(false)
const selectedPost = ref<CirclePostData>()
const comments = ref<{ id: number; user_id: number; content: string; created_at: string; author: UserInfo }[]>([])
const newComment = ref("")
const showJoinReview = ref(false)
const joinRequests = ref<any[]>([])

const tabs = [{ key: "posts", label: "帖子" }, { key: "members", label: "成员" }, { key: "activities", label: "动态" }]
const circleId = ref(0)

const showReviewBtn = computed(() => {
  if (!circle.value) return false
  return circle.value.join_type === "approve" && userStore.id === circle.value.creator_id
})

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  circleId.value = page?.options?.id ? parseInt(page.options.id) : 0
  if (circleId.value) load()
})

async function load() {
  try {
    const [c, m, p, a] = await Promise.all([
      getCircle(circleId.value),
      listCircleMembers(circleId.value),
      listCirclePosts(circleId.value),
      listCircleActivities(circleId.value),
    ])
    circle.value = c
    members.value = m
    posts.value = p
    activities.value = a
  } catch {
    uni.showToast({ title: "加载失败", icon: "none" })
  } finally {
    loading.value = false
  }
}

async function handleJoin() {
  try {
    const res = await joinCircle(circleId.value)
    if (circle.value?.join_type === "approve") {
      uni.showToast({ title: "已发送申请，请等待审核", icon: "success" })
    } else {
      uni.showToast({ title: "已加入", icon: "success" })
      circle.value!.is_member = true
      circle.value!.member_count++
    }
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
}

async function handleLeave() {
  try {
    await leaveCircle(circleId.value)
    uni.showToast({ title: "已退出", icon: "success" })
    circle.value!.is_member = false
    circle.value!.member_count--
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
}

async function openJoinReview() {
  if (!circleId.value) return
  try {
    joinRequests.value = await getCircleJoinRequests(circleId.value)
    showJoinReview.value = true
  } catch {
    uni.showToast({ title: "加载失败", icon: "none" })
  }
}

async function handleJoinRequestAction(requestId: number, action: "approve" | "reject") {
  try {
    await handleJoinRequest(requestId, action)
    joinRequests.value = joinRequests.value.filter((r: any) => r.id !== requestId)
    if (action === "approve") {
      circle.value!.member_count++
    }
    uni.showToast({ title: action === "approve" ? "已通过" : "已拒绝", icon: "success" })
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
}

async function handleCreatePost() {
  if (!newPostContent.value.trim()) return
  try {
    await createCirclePost(circleId.value, { content: newPostContent.value.trim() })
    uni.showToast({ title: "发布成功", icon: "success" })
    newPostContent.value = ""
    showCreatePost.value = false
    posts.value = await listCirclePosts(circleId.value)
  } catch {
    uni.showToast({ title: "发布失败", icon: "none" })
  }
}

async function toggleLike(p: CirclePostData) {
  try {
    const res = await toggleCirclePostLike(p.id)
    p.is_liked = res.liked
    p.like_count += res.liked ? 1 : -1
  } catch {}
}

function showComments(p: CirclePostData) {
  selectedPost.value = p
  newComment.value = ""
  showCommentsPanel.value = true
  loadComments(p.id)
}

async function loadComments(postId: number) {
  try {
    comments.value = await listCirclePostComments(postId)
  } catch {}
}

async function handleCreateComment() {
  if (!newComment.value.trim() || !selectedPost.value) return
  try {
    await createCirclePostComment(selectedPost.value.id, newComment.value.trim())
    newComment.value = ""
    selectedPost.value.comment_count++
    await loadComments(selectedPost.value.id)
  } catch {
    uni.showToast({ title: "评论失败", icon: "none" })
  }
}

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = (now.getTime() - d.getTime()) / 1000
  if (diff < 60) return "刚刚"
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<style lang="scss" scoped>
.page { min-height: 100vh; background: #F7F4F0; }

.loading-state, .error-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 300rpx 0; gap: 20rpx;
}
.loading-text, .error-text { font-size: 26rpx; color: #8E9BAB; }

.circle-header { position: relative; overflow: hidden; }
.circle-header-bg {
  position: absolute; inset: 0;
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 50%, #3A4A5E 100%);
  z-index: 0;
  &::after {
    content: ''; position: absolute;
    top: -200rpx; right: -100rpx;
    width: 500rpx; height: 500rpx;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(198,122,106,0.15) 0%, transparent 70%);
  }
}
.circle-header-content {
  position: relative; z-index: 1; padding: 40rpx 32rpx;
}
.circle-header-top { display: flex; align-items: flex-start; gap: 24rpx; }
.circle-info { flex: 1; padding-top: 8rpx; }
.circle-name { font-size: 36rpx; font-weight: 700; color: #fff; display: block; }
.circle-desc { font-size: 24rpx; color: rgba(255,255,255,0.6); margin-top: 8rpx; display: block; }
.circle-stats { display: flex; gap: 24rpx; margin-top: 12rpx; }
.circle-stat { font-size: 22rpx; color: rgba(255,255,255,0.5); }
.circle-stat.join-type { color: #FFB74D; }

.circle-actions { display: flex; gap: 16rpx; margin-top: 24rpx; }

.activities-section { padding: 20rpx 16rpx; padding-bottom: 80rpx; }
.activity-item { display: flex; gap: 12rpx; padding: 16rpx 0; border-bottom: 1rpx solid #EAE6E0; }
.activity-body { flex: 1; }
.activity-text { font-size: 24rpx; color: #5C6B7E; display: block; }
.activity-name { color: #C67A6A; font-weight: 600; }
.activity-preview { font-size: 22rpx; color: #B8C2CE; margin-top: 4rpx; display: block; }
.activity-time { font-size: 20rpx; color: #B8C2CE; margin-top: 4rpx; display: block; }
.tab-bar {
  display: flex; background: #fff; padding: 0 40rpx;
  border-bottom: 1rpx solid #E0DBD4; position: sticky; top: 0; z-index: 50;
  align-items: center;
}
.tab-item { flex: 1; display: flex; flex-direction: column; align-items: center; padding: 20rpx 0; cursor: pointer; }
.tab-text { font-size: 28rpx; color: #8E9BAB; font-weight: 500; }
.tab-active .tab-text { color: #C67A6A; font-weight: 600; }
.tab-indicator { width: 40rpx; height: 4rpx; border-radius: 2rpx; background: #C67A6A; margin-top: 12rpx; }
.review-btn {
  padding: 8rpx 20rpx; background: #FFF3E0; border-radius: 20rpx;
  cursor: pointer; flex-shrink: 0;
}
.review-btn-text { font-size: 24rpx; color: #E65100; font-weight: 500; }

.posts-section, .members-section { padding: 20rpx 16rpx; display: flex; flex-direction: column; gap: 16rpx; padding-bottom: 80rpx; }

.post-card {
  background: #fff; border-radius: 14rpx; padding: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); cursor: pointer;
  transition: all 0.2s ease;
  &:active { transform: scale(0.99); }
}
.post-card-header { display: flex; align-items: center; gap: 12rpx; margin-bottom: 12rpx; }
.post-card-author { flex: 1; }
.post-card-name { font-size: 26rpx; font-weight: 600; color: #1E2A3A; display: block; }
.post-card-time { font-size: 20rpx; color: #B8C2CE; margin-top: 4rpx; display: block; }
.post-card-content { font-size: 26rpx; color: #5C6B7E; line-height: 1.6; display: block; }
.post-card-footer {
  display: flex; gap: 32rpx; margin-top: 16rpx; padding-top: 16rpx;
  border-top: 1rpx solid #EAE6E0;
}
.post-action { display: flex; align-items: center; gap: 8rpx; &:active { transform: scale(0.9); } }
.action-icon { font-size: 24rpx; color: #B8C2CE; }
.action-icon.liked { color: #E74C3C; }
.action-text { font-size: 22rpx; color: #B8C2CE; }
.action-text.liked { color: #E74C3C; }

.member-card {
  display: flex; align-items: center; gap: 16rpx;
  background: #fff; border-radius: 14rpx; padding: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
}
.member-info { flex: 1; }
.member-name { font-size: 26rpx; font-weight: 600; color: #1E2A3A; display: block; }
.member-bio { font-size: 22rpx; color: #8E9BAB; margin-top: 4rpx; display: block; }
.role-badge {
  font-size: 20rpx; padding: 4rpx 14rpx; border-radius: 8rpx; color: #8E9BAB; background: #F0EDE8;
}
.role-badge.creator { background: #FFF3E0; color: #E65100; }

.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 120rpx 0; gap: 20rpx;
}
.empty-text { font-size: 26rpx; color: #B8C2CE; }

.create-post-panel { padding: 32rpx; }
.create-post-title { font-size: 32rpx; font-weight: 600; color: #1E2A3A; margin-bottom: 24rpx; display: block; }
.create-post-input {
  width: 100%; height: 240rpx; background: #F7F4F0; border-radius: 14rpx;
  padding: 24rpx; font-size: 26rpx; margin-bottom: 24rpx;
  resize: none; border: none; outline: none;
}

.review-panel { padding: 32rpx; }
.review-title { font-size: 32rpx; font-weight: 600; color: #1E2A3A; margin-bottom: 24rpx; display: block; }
.review-empty { text-align: center; padding: 40rpx 0; color: #B8C2CE; font-size: 26rpx; }
.review-item { display: flex; align-items: center; gap: 16rpx; padding: 16rpx 0; border-bottom: 1rpx solid #EAE6E0; }
.review-item-info { flex: 1; }
.review-item-name { font-size: 26rpx; font-weight: 600; color: #1E2A3A; display: block; }
.review-item-time { font-size: 20rpx; color: #B8C2CE; margin-top: 4rpx; display: block; }
.review-item-actions { display: flex; gap: 12rpx; }

.comments-panel { padding: 0; }
.comments-original { padding: 24rpx 32rpx; border-bottom: 1rpx solid #EAE6E0; }
.comments-original-author { font-size: 24rpx; font-weight: 600; color: #1E2A3A; }
.comments-original-content { font-size: 24rpx; color: #5C6B7E; }
.comments-list { max-height: 400rpx; overflow-y: auto; padding: 16rpx 32rpx; }
.comment-item { padding: 12rpx 0; display: flex; gap: 8rpx; }
.comment-author { font-size: 22rpx; color: #C67A6A; font-weight: 500; white-space: nowrap; }
.comment-content { font-size: 22rpx; color: #5C6B7E; }
.comments-empty { text-align: center; padding: 40rpx; color: #B8C2CE; font-size: 24rpx; }
.comments-input-bar {
  display: flex; gap: 16rpx; padding: 16rpx 32rpx 40rpx;
  border-top: 1rpx solid #EAE6E0;
}
.comments-input {
  flex: 1; background: #F7F4F0; border-radius: 20rpx; padding: 12rpx 24rpx;
  font-size: 24rpx; border: none; outline: none;
}
.comments-send { font-size: 26rpx; color: #C67A6A; font-weight: 500; cursor: pointer; padding: 12rpx 0; }
</style>
