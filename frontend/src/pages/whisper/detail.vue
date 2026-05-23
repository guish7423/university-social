<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getWhisper, toggleWhisperLike, createWhisperComment, listWhisperComments, type WhisperData, type WhisperComment } from "@/api/whisper"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const whisper = ref<WhisperData | null>(null)
const comments = ref<WhisperComment[]>([])
const commentText = ref("")
const submitting = ref(false)

const whisperId = ref(0)

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  whisperId.value = Number(page?.options?.id || 0)
  if (!whisperId.value) {
    uni.showToast({ title: "参数错误", icon: "error" })
    return
  }
  loadWhisper(whisperId.value)
  loadComments(whisperId.value)
})

async function loadWhisper(id: number) {
  try {
    whisper.value = await getWhisper(id)
  } catch (e) { console.error(e) }
}

async function loadComments(id: number) {
  try {
    comments.value = await listWhisperComments(id)
  } catch (e) { console.error(e) }
}

async function handleLike() {
  if (!userStore.isLogin || !whisper.value) return
  const res = await toggleWhisperLike(whisper.value.id)
  whisper.value.is_liked = res.liked
  whisper.value.like_count += res.liked ? 1 : -1
}

async function submitComment() {
  if (!commentText.value.trim() || !whisper.value) return
  submitting.value = true
  try {
    await createWhisperComment(whisper.value.id, commentText.value.trim())
    commentText.value = ""
    comments.value = await listWhisperComments(whisper.value.id)
    uni.showToast({ title: "评论成功", icon: "success" })
  } catch (e) { console.error(e) }
  submitting.value = false
}

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  const days = Math.floor(diff / 86400000)
  if (days < 7) return `${days}天前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<template>
  <view class="container">
    <view v-if="whisper" class="detail">
      <view class="detail-card">
        <view class="card-top">
          <view class="codename-badge">
            <text class="codename-icon">🫥</text>
            <text class="codename-text">{{ whisper.codename || "匿名用户" }}</text>
          </view>
          <text v-if="whisper.whisper_type === 1" class="confess-tag">💕 表白</text>
          <text class="card-time">{{ formatTime(whisper.created_at) }}</text>
        </view>

        <view class="card-body">
          <text class="card-content">{{ whisper.content }}</text>
          <view v-if="whisper?.images?.length" class="card-images">
            <image
              v-for="(img, i) in whisper.images"
              :key="i"
              class="card-image"
              :src="img"
              mode="aspectFill"
            />
          </view>
        </view>

        <view class="card-actions">
          <view class="action-btn" @click="handleLike">
            <u-icon :name="whisper.is_liked ? 'heart-fill' : 'heart'" :color="whisper.is_liked ? '#f56c6c' : '#c0c4cc'" size="36" />
            <text :class="['action-text', whisper.is_liked && 'liked']">{{ whisper.like_count }}</text>
          </view>
          <view class="action-btn">
            <u-icon name="chat" color="#c0c4cc" size="36" />
            <text class="action-text">{{ whisper.comment_count }}</text>
          </view>
        </view>
      </view>

      <view class="comment-section">
        <text class="section-title">评论 ({{ comments.length }})</text>

        <view v-if="!comments.length" class="no-comments">
          <text class="no-comment-text">暂无评论，来说点什么吧</text>
        </view>

        <view v-for="c in comments" :key="c.id" class="comment-item">
          <u-avatar
            :src="c.author?.avatar"
            :text="c.codename?.[0] || '匿'"
            size="56"
            shape="circle"
          />
          <view class="comment-body">
            <text class="comment-author">{{ c.codename || "匿名" }}</text>
            <text class="comment-content">{{ c.content }}</text>
            <text class="comment-time">{{ formatTime(c.created_at) }}</text>
          </view>
        </view>
      </view>
    </view>

    <view v-if="userStore.isLogin && whisper" class="comment-input">
        <input
          v-model="commentText"
          placeholder="写下你的评论..."
          class="comment-field"
        />
      <u-button
        type="primary"
        size="small"
        :loading="submitting"
        :disabled="!commentText.trim()"
        @click="submitComment"
      >
        发送
      </u-button>
    </view>
  </view>
</template>

<style scoped lang="scss">
.container {
  min-height: 100vh;
  background: #f8f9fa;
  padding-bottom: 120rpx;
}

.detail { padding: 32rpx; }

.detail-card {
  background: #fff;
  border-radius: 22rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
}

.card-top {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.codename-badge {
  display: flex;
  align-items: center;
  gap: 6rpx;
  background: linear-gradient(135deg, #C67A6A15, #1E2A3A15);
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
}

.codename-icon { font-size: 24rpx; }

.codename-text {
  font-size: 24rpx;
  color: #C67A6A;
  font-weight: 600;
}

.confess-tag {
  font-size: 22rpx;
  background: #fef0f0;
  color: #f56c6c;
  padding: 4rpx 12rpx;
  border-radius: 12rpx;
}

.card-time {
  margin-left: auto;
  font-size: 22rpx;
  color: #c0c4cc;
}

.card-body { margin-bottom: 20rpx; }

.card-content {
  font-size: 30rpx;
  color: #1a1a2e;
  line-height: 1.6;
}

.card-images {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
  margin-top: 12rpx;
}

.card-image {
  width: 220rpx;
  height: 220rpx;
  border-radius: 12rpx;
}

.card-actions {
  display: flex;
  gap: 32rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f0f0f0;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.action-text {
  font-size: 26rpx;
  color: #909399;
}

.liked { color: #f56c6c; }

.comment-section {
  margin-top: 32rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1a1a2e;
  display: block;
  margin-bottom: 20rpx;
}

.no-comments { text-align: center; padding: 60rpx 0; }

.no-comment-text {
  font-size: 26rpx;
  color: #c0c4cc;
}

.comment-item {
  display: flex;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.comment-body { flex: 1; }

.comment-author {
  font-size: 24rpx;
  font-weight: 600;
  color: #333;
  display: block;
}

.comment-content {
  font-size: 26rpx;
  color: #555;
  line-height: 1.5;
  margin-top: 4rpx;
  display: block;
}

.comment-time {
  font-size: 20rpx;
  color: #c0c4cc;
  margin-top: 4rpx;
  display: block;
}

.comment-input {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 16rpx 32rpx;
  display: flex;
  align-items: center;
  gap: 16rpx;
  border-top: 1rpx solid #f0f0f0;
}
</style>
