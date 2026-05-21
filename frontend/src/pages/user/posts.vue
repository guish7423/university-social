<script setup lang="ts">
import { ref, onMounted } from "vue"
import { userPosts, toggleLike, type PostData } from "@/api/post"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(false)

onMounted(fetchPosts)

async function fetchPosts() {
  loading.value = true
  try {
    posts.value = await userPosts(userStore.id || 0)
  } catch {}
  loading.value = false
}

async function handleLike(post: PostData) {
  if (!userStore.isLogin) return
  const res = await toggleLike(post.id)
  post.is_liked = res.liked
  post.like_count += res.liked ? 1 : -1
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/post/detail?id=${id}` })
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
    <u-navbar title="我的帖子" :autoBack="true" bgColor="#fff" />

    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
    </view>

    <template v-else-if="posts && posts.length">
      <view class="post-list">
        <view
          v-for="(post, idx) in posts"
          :key="post.id"
          :class="['post-card', 'animate-fade-in-up', 'stagger-' + Math.min(idx + 1, 8)]"
          @click="goDetail(post.id)"
        >
          <view class="post-body">
            <text class="post-content">{{ post.content }}</text>
            <view v-if="post?.images?.length" class="post-images">
              <image
                v-for="(img, i) in post.images.slice(0, 3)"
                :key="i"
                class="post-image"
                :src="img"
                mode="aspectFill"
              />
            </view>
          </view>

          <view class="post-actions">
            <view :class="['action-btn', post.is_liked && 'is-liked']" @click.stop="handleLike(post)">
              <u-icon :name="'heart'" :color="post.is_liked ? '#e74c3c' : '#c0c4cc'" :size="32" />
              <text :class="['action-count', post.is_liked && 'liked']">{{ post.like_count || 0 }}</text>
            </view>
            <view class="action-btn">
              <u-icon name="chat" color="#c0c4cc" size="32" />
              <text class="action-count">{{ post.comment_count || 0 }}</text>
            </view>
            <text class="post-time">{{ formatTime(post.created_at) }}</text>
          </view>
        </view>
      </view>
    </template>

    <template v-else>
      <u-empty mode="data" text="还没有发过动态" />
    </template>
  </view>
</template>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: var(--u-bg-color, #f3f4f6);
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.post-list {
  padding: 20rpx;
}

.post-card {
  background: #fff;
  border-radius: 22rpx;
  padding: 32rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.04);
  &:active { opacity: 0.8; }
}

.post-body {
  margin-bottom: 16rpx;
}

.post-content {
  font-size: 30rpx;
  line-height: 1.7;
  color: #303133;
  display: block;
}

.post-images {
  display: flex;
  gap: 12rpx;
  margin-top: 16rpx;
}

.post-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 14rpx;
  background: #f5f5f5;
}

.post-actions {
  display: flex;
  align-items: center;
  gap: 32rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid var(--hairline, #e4e7ed);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6rpx;
  &:active { transform: scale(0.9); }
}

.action-count {
  font-size: 24rpx;
  color: #c0c4cc;
}

.action-count.liked {
  color: #e74c3c;
}

.post-time {
  margin-left: auto;
  font-size: 22rpx;
  color: #c0c4cc;
}
</style>
