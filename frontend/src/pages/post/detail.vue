<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getPost, listComments, createComment, type PostData, type CommentData } from "@/api/post"

const post = ref<PostData | null>(null)
const comments = ref<CommentData[]>([])
const commentText = ref("")
const submitting = ref(false)
const postId = ref(0)

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
  try { post.value = await getPost(id) } catch {}
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

.post-author {
  flex: 1;
}

.author-name {
  font-size: 28rpx;
  font-weight: 600;
  display: block;
}

.post-time {
  font-size: 22rpx;
  color: #c0c4cc;
  display: block;
  margin-top: 4rpx;
}

.post-content {
  font-size: 30rpx;
  line-height: 1.8;
  display: block;
  margin-bottom: 20rpx;
}

.post-images {
  margin-bottom: 20rpx;
}

.image-full {
  width: 100%;
  border-radius: 12rpx;
  margin-bottom: 8rpx;
}

.post-stats {
  font-size: 24rpx;
  color: #c0c4cc;
  padding-top: 16rpx;
  border-top: 1rpx solid var(--u-divider-color, #e4e7ed);
}

.comment-section {
  background: #fff;
  padding: 30rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #303133;
  margin-bottom: 24rpx;
  display: block;
}

.comment-item {
  display: flex;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.comment-body {
  flex: 1;
}

.comment-author {
  font-size: 24rpx;
  color: #667eea;
  margin-bottom: 4rpx;
  display: block;
}

.comment-content {
  font-size: 28rpx;
  line-height: 1.5;
  color: #303133;
}

.comment-input-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 16rpx 30rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  background: #fff;
  border-top: 1rpx solid var(--u-divider-color, #e4e7ed);
}
</style>
