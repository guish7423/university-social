<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getPost, listComments, createComment, type PostData, type CommentData } from "@/api/post"
import { followUser, unfollowUser, checkFollow } from "@/api/follow"

const post = ref<PostData | null>(null)
const comments = ref<CommentData[]>([])
const commentText = ref("")
const submitting = ref(false)
const postId = ref(0)
const isFollowingAuthor = ref(false)
const followLoading = ref(false)

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  postId.value = Number(page?.options?.id || 0)
  if (postId.value) {
    loadPost(postId.value)
    loadComments(postId.value)
  }
})

async function loadPost(id: number) {
  try {
    post.value = await getPost(id)
    if (post.value?.author?.id) {
      checkFollowState(post.value.author.id)
    }
  } catch {}
}

async function checkFollowState(userId: number) {
  try {
    const res = await checkFollow(userId)
    isFollowingAuthor.value = res.is_following
  } catch {}
}

async function loadComments(id: number) {
  try { comments.value = await listComments(id) } catch {}
}

async function handleComment() {
  if (!commentText.value.trim()) return
  submitting.value = true
  try {
    await createComment(postId.value, commentText.value)
    commentText.value = ""
    if (post.value) post.value.comment_count++
    loadComments(postId.value)
  } catch {}
  submitting.value = false
}

async function toggleFollowAuthor() {
  if (!post.value?.author?.id || followLoading.value) return
  followLoading.value = true
  try {
    if (isFollowingAuthor.value) {
      await unfollowUser(post.value.author.id)
      isFollowingAuthor.value = false
    } else {
      await followUser(post.value.author.id)
      isFollowingAuthor.value = true
    }
  } catch {}
  followLoading.value = false
}

function formatTime(t: string) {
  const d = new Date(t)
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`
}
</script>

<template>
  <view class="container">
    <view v-if="!post" class="loading-state">
      <u-loading mode="flower" size="60" />
    </view>
    <template v-else>
      <view class="post-card">
        <view class="post-header">
          <u-avatar
            :src="post.author?.avatar"
            :text="post.author?.nickname?.[0]"
            size="64"
            shape="circle"
          />
          <view class="post-author">
            <text class="author-name">{{ post.author?.nickname || "匿名" }}</text>
            <text class="post-time">{{ formatTime(post.created_at) }}</text>
          </view>
          <u-button
            v-if="post.author?.id"
            size="mini"
            shape="circle"
            :type="isFollowingAuthor ? 'default' : 'primary'"
            :text="isFollowingAuthor ? '已关注' : '+ 关注'"
            :loading="followLoading"
            :customStyle="{ marginLeft: 'auto', minWidth: '100rpx' }"
            @click="toggleFollowAuthor"
          />
        </view>
        <text class="post-content">{{ post.content }}</text>
        <view v-if="post.images?.length" class="post-images">
          <image v-for="(img, i) in post.images" :key="i" class="image-full" :src="img" mode="widthFix" />
        </view>
        <view class="post-stats">
          <text>{{ post.like_count }} 赞 · {{ post.comment_count }} 评论</text>
        </view>
      </view>

      <view class="comment-section">
        <text class="section-title">评论</text>
        <u-empty
          v-if="!comments?.length"
          mode="message"
          text="暂无评论，来写第一条吧"
          :iconSize="80"
        />
        <view v-for="c in (comments || [])" :key="c.id" class="comment-item">
          <u-avatar
            :src="c.author?.avatar"
            :text="c.author?.nickname?.[0]"
            size="48"
            shape="circle"
            fontSize="20"
          />
          <view class="comment-body">
            <text class="comment-author">{{ c.author?.nickname || "匿名" }}</text>
            <text class="comment-content">{{ c.content }}</text>
          </view>
        </view>
      </view>
    </template>

    <view class="comment-input-bar">
      <u-input
        v-model="commentText"
        placeholder="写下你的评论..."
        :border="false"
        shape="circle"
        :customStyle="{ flex: 1, height: '72rpx', background: '#f5f5f5', padding: '0 24rpx', fontSize: '28rpx', borderRadius: '36rpx' }"
      />
      <u-button
        type="primary"
        size="small"
        shape="circle"
        :loading="submitting"
        :disabled="!commentText.trim()"
        :customStyle="{ height: '72rpx', padding: '0 30rpx', fontSize: '26rpx' }"
        @click="handleComment"
      >发送</u-button>
    </view>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background: var(--u-bg-color, #f3f4f6);
  padding-bottom: 120rpx;
}
.loading-state {
  display: flex;
  justify-content: center;
  padding: 200rpx 0;
}
.post-card {
  background: #fff;
  padding: 30rpx;
  margin-bottom: 16rpx;
}
.post-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 20rpx;
}
.post-header .u-button {
  margin-left: auto;
}
.post-author {
  flex: 1;
}
.author-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  display: block;
}
.post-time {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}
.post-content {
  font-size: 30rpx;
  line-height: 1.6;
  color: #303133;
  white-space: pre-wrap;
}
.post-images {
  margin-top: 20rpx;
}
.image-full {
  width: 100%;
  border-radius: 12rpx;
  margin-bottom: 12rpx;
}
.post-stats {
  margin-top: 20rpx;
  font-size: 24rpx;
  color: #999;
}
.comment-section {
  background: #fff;
  padding: 30rpx;
}
.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20rpx;
  display: block;
}
.comment-item {
  display: flex;
  gap: 12rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}
.comment-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.comment-author {
  font-size: 24rpx;
  font-weight: 600;
  color: #606266;
}
.comment-content {
  font-size: 28rpx;
  color: #303133;
  line-height: 1.5;
}
.comment-input-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 16rpx;
  align-items: center;
  background: #fff;
  padding: 16rpx 30rpx;
  border-top: 1rpx solid #eee;
  z-index: 100;
}
</style>

@media (min-width: 1024px) {
.container {
  max-width: 960px;
  margin: 0 auto;
  padding: 32rpx 0 120rpx;
}
.post-card { border-radius: 22rpx; margin: 0 0 16rpx; }
.comment-section { border-radius: 22rpx; }
.comment-input-bar {
  max-width: 960px;
  left: 50%;
  transform: translateX(-50%);
  border-radius: 999rpx;
  bottom: 32rpx;
}
}
