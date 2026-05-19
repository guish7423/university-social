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
    <view v-if="!post" class="loading">加载中...</view>
    <template v-else>
      <view class="post-card">
        <view class="post-header">
          <image v-if="post.author?.avatar" class="avatar" :src="post.author.avatar" mode="aspectFill" />
          <text class="nickname">{{ post.author?.nickname || "匿名" }}</text>
          <text class="time">{{ formatTime(post.created_at) }}</text>
        </view>
        <text class="content">{{ post.content }}</text>
        <view v-if="post.images?.length" class="images">
          <image v-for="(img, i) in post.images" :key="i" class="image-full" :src="img" mode="widthFix" />
        </view>
        <view class="post-stats">
          <text>{{ post.like_count }} 赞 · {{ post.comment_count }} 评论</text>
        </view>
      </view>

      <view class="comment-section">
        <text class="section-title">评论 ({{ post.comment_count }})</text>
        <view v-if="comments.length === 0" class="empty">暂无评论</view>
        <view v-for="c in comments" :key="c.id" class="comment-item">
          <image v-if="c.author?.avatar" class="comment-avatar" :src="c.author.avatar" mode="aspectFill" />
          <view class="comment-body">
            <text class="comment-author">{{ c.author?.nickname || "匿名" }}</text>
            <text class="comment-content">{{ c.content }}</text>
          </view>
        </view>
      </view>

      <view class="comment-input-bar">
        <input v-model="commentText" class="comment-input" placeholder="写下你的评论..." />
        <button class="send-btn" :loading="submitting" @click="handleComment">发送</button>
      </view>
    </template>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; padding-bottom: 120rpx; }
.loading, .empty { text-align: center; padding: 100rpx 0; color: #999; }
.post-card { background: #fff; padding: 30rpx; margin-bottom: 16rpx; }
.post-header { display: flex; align-items: center; gap: 16rpx; margin-bottom: 16rpx; }
.avatar { width: 60rpx; height: 60rpx; border-radius: 50%; }
.nickname { font-size: 28rpx; font-weight: 600; flex: 1; }
.time { font-size: 22rpx; color: #999; }
.content { font-size: 30rpx; line-height: 1.8; margin-bottom: 16rpx; }
.images { margin-bottom: 16rpx; }
.image-full { width: 100%; border-radius: 8rpx; margin-bottom: 8rpx; }
.post-stats { font-size: 24rpx; color: #999; padding-top: 16rpx; border-top: 1rpx solid #eee; }
.comment-section { background: #fff; padding: 30rpx; }
.section-title { font-size: 28rpx; font-weight: 600; margin-bottom: 20rpx; display: block; }
.comment-item { display: flex; gap: 16rpx; margin-bottom: 24rpx; }
.comment-avatar { width: 48rpx; height: 48rpx; border-radius: 50%; flex-shrink: 0; }
.comment-body { flex: 1; }
.comment-author { font-size: 24rpx; color: #667eea; margin-bottom: 4rpx; display: block; }
.comment-content { font-size: 28rpx; line-height: 1.5; }
.comment-input-bar { position: fixed; bottom: 0; left: 0; right: 0; display: flex; align-items: center; gap: 16rpx; padding: 16rpx 30rpx; background: #fff; border-top: 1rpx solid #eee; }
.comment-input { flex: 1; height: 72rpx; background: #f5f5f5; border-radius: 36rpx; padding: 0 24rpx; font-size: 28rpx; }
.send-btn { font-size: 26rpx; padding: 12rpx 30rpx; background: #667eea; color: #fff; border-radius: 30rpx; }
</style>
