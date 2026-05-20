<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { listPosts, toggleLike, type PostData } from "@/api/post"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(false)

onMounted(fetchPosts)

const unsub = uni.$on('post-created', () => {
  fetchPosts()
})

onUnmounted(() => {
  if (typeof unsub === 'function') unsub()
})

async function fetchPosts() {
  loading.value = true
  try {
    posts.value = await listPosts()
  } catch {}
  loading.value = false
}

async function handleLike(post: PostData) {
  if (!userStore.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" })
    return
  }
  const res = await toggleLike(post.id)
  post.is_liked = res.liked
  post.like_count += res.liked ? 1 : -1
}

function goCreate() {
  uni.navigateTo({ url: "/pages/post/create" })
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
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="posts && posts.length">
      <view class="post-list">
        <view
          v-for="post in posts"
          :key="post.id"
          class="post-card"
          @click="goDetail(post.id)"
        >
          <view class="post-header">
            <u-avatar
              :src="post.author?.avatar"
              :text="post.author?.nickname?.[0]"
              size="72"
              shape="circle"
            />
            <view class="post-author">
              <text class="author-name">{{ post.author?.nickname || "匿名" }}</text>
              <text class="post-time">{{ formatTime(post.created_at) }}</text>
            </view>
            <u-icon name="more" color="#c0c4cc" size="32" />
          </view>

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
            <view class="action-btn" @click.stop="handleLike(post)">
              <u-icon
                :name="post.is_liked ? 'heart' : 'heart'"
                :color="post.is_liked ? '#e74c3c' : '#c0c4cc'"
                :size="36"
              />
              <text :class="['action-count', post.is_liked && 'liked']">{{ post.like_count || 0 }}</text>
            </view>
            <view class="action-btn">
              <u-icon name="chat" color="#c0c4cc" size="36" />
              <text class="action-count">{{ post.comment_count || 0 }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="fab" @click="goCreate">
        <u-icon name="plus" color="#fff" size="44" />
      </view>
    </template>

    <template v-else>
      <u-empty mode="data" text="还没有动态，快来发第一条吧" />
      <view class="empty-create">
        <u-button type="primary" shape="circle" @click="goCreate">发布第一条动态</u-button>
      </view>
    </template>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background: var(--u-bg-color, #f3f4f6);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 200rpx 0;
  gap: 20rpx;
}

.loading-text {
  font-size: 26rpx;
  color: #909399;
}

.post-list {
  padding: 20rpx;
}

.post-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.04);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 20rpx;
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
  color: #c0c4cc;
  display: block;
  margin-top: 4rpx;
}

.post-body {
  margin-bottom: 20rpx;
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
  width: 210rpx;
  height: 210rpx;
  border-radius: 12rpx;
  background: #f5f5f5;
}

.post-actions {
  display: flex;
  align-items: center;
  gap: 48rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid var(--u-divider-color, #e4e7ed);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.action-count {
  font-size: 24rpx;
  color: #c0c4cc;
}

.action-count.liked {
  color: #e74c3c;
}

.fab {
  position: fixed;
  right: 40rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 32rpx rgba(102, 126, 234, 0.5);
  z-index: 100;
}

.empty-create {
  padding: 40rpx 60rpx;
}
</style>
