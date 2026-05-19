<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listPosts, toggleLike, type PostData } from "@/api/post"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    posts.value = await listPosts()
  } catch {}
  loading.value = false
})

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
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<template>
  <view class="container">
    <view class="header-actions">
      <text class="page-title">动态广场</text>
      <button class="create-btn" @click="goCreate">发帖</button>
    </view>

    <view v-if="loading" class="loading">加载中...</view>

    <view v-else-if="posts.length === 0" class="empty">还没有动态，快来发第一条吧</view>

    <view v-for="post in posts" :key="post.id" class="post-card" @click="goDetail(post.id)">
      <view class="post-header">
        <image v-if="post.author?.avatar" class="avatar" :src="post.author.avatar" mode="aspectFill" />
        <text class="nickname">{{ post.author?.nickname || "匿名" }}</text>
        <text class="time">{{ formatTime(post.created_at) }}</text>
      </view>
      <text class="content">{{ post.content }}</text>
      <view v-if="post.images?.length" class="images">
        <image v-for="(img, i) in post.images.slice(0, 3)" :key="i" class="image-thumb" :src="img" mode="aspectFill" />
      </view>
      <view class="post-actions">
        <view class="action-item" @click.stop="handleLike(post)">
          <text :class="['like-icon', post.is_liked && 'liked']">{{ post.is_liked ? "❤" : "♡" }}</text>
          <text class="action-count">{{ post.like_count || "" }}</text>
        </view>
        <view class="action-item">
          <text class="comment-icon">💬</text>
          <text class="action-count">{{ post.comment_count || "" }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; padding: 20rpx; }
.header-actions { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; }
.page-title { font-size: 40rpx; font-weight: 700; }
.create-btn { font-size: 26rpx; padding: 10rpx 30rpx; background: #667eea; color: #fff; border-radius: 30rpx; }
.loading, .empty { text-align: center; padding: 200rpx 0; color: #999; font-size: 28rpx; }
.post-card { background: #fff; border-radius: 16rpx; padding: 30rpx; margin-bottom: 20rpx; }
.post-header { display: flex; align-items: center; gap: 16rpx; margin-bottom: 16rpx; }
.avatar { width: 60rpx; height: 60rpx; border-radius: 50%; }
.nickname { font-size: 28rpx; font-weight: 600; flex: 1; }
.time { font-size: 22rpx; color: #999; }
.content { font-size: 30rpx; line-height: 1.6; margin-bottom: 16rpx; }
.images { display: flex; gap: 10rpx; margin-bottom: 16rpx; }
.image-thumb { width: 200rpx; height: 200rpx; border-radius: 8rpx; }
.post-actions { display: flex; gap: 40rpx; padding-top: 16rpx; border-top: 1rpx solid #eee; }
.action-item { display: flex; align-items: center; gap: 6rpx; }
.like-icon.liked { color: #e74c3c; }
.action-count { font-size: 24rpx; color: #999; }
</style>
