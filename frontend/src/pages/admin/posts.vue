<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listPosts, deletePost } from "@/api/admin"

const posts = ref<any[]>([])
const loading = ref(true)
const offset = ref(0)
const hasMore = ref(true)

async function load() {
  try {
    const data = await listPosts(offset.value)
    if (data.length < 50) hasMore.value = false
    posts.value = offset.value === 0 ? data : [...posts.value, ...data]
  } catch {}
  loading.value = false
}

onMounted(load)

async function doDelete(id: number) {
  uni.showModal({
    title: "确认删除",
    content: "确定要删除这条帖子吗？",
    success: async (res) => {
      if (res.confirm) {
        uni.showLoading({ title: "删除中" })
        try {
          await deletePost(id)
          posts.value = posts.value.filter(p => p.id !== id)
          uni.showToast({ title: "已删除", icon: "success" })
        } catch {
          uni.showToast({ title: "删除失败", icon: "none" })
        }
        uni.hideLoading()
      }
    }
  })
}
</script>

<template>
  <view class="page">
    <view class="page-header">
      <text class="page-title">帖子管理 ({{ posts.length }})</text>
    </view>

    <view class="list">
      <view v-for="p in posts" :key="p.id" class="post-item">
        <view class="post-header">
          <u-avatar :src="p.avatar" :text="p.nickname?.[0] || '?'" size="56" shape="circle" />
          <text class="post-author">{{ p.nickname }}</text>
          <text class="post-time">{{ p.created_at?.slice(0, 16) }}</text>
        </view>
        <text class="post-content">{{ p.content }}</text>
        <view v-if="p.images" class="post-images">
          <image v-for="(img, i) in p.images?.split(',')" :key="i" :src="img" mode="aspectFill" class="post-image" />
        </view>
        <view class="post-footer">
          <text class="post-stat">❤ {{ p.like_count }}</text>
          <text class="post-stat">💬 {{ p.comment_count }}</text>
          <text class="delete-btn" @click="doDelete(p.id)">删除</text>
        </view>
      </view>
      <view v-if="posts.length === 0" class="empty">暂无帖子</view>
    </view>

    <view v-if="hasMore" class="load-more" @click="offset += 50; load()">
      加载更多
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 24rpx 16rpx; }
.page-header { margin-bottom: 24rpx; }
.page-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }

.list { background: #fff; border-radius: 14rpx; overflow: hidden; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); }
.post-item { padding: 24rpx; border-bottom: 1rpx solid #EAE6E0; }
.post-header { display: flex; align-items: center; gap: 12rpx; margin-bottom: 12rpx; }
.post-author { font-size: 24rpx; font-weight: 500; color: #1E2A3A; flex: 1; }
.post-time { font-size: 20rpx; color: #B8C2CE; }
.post-content { font-size: 24rpx; color: #333; line-height: 1.6; display: block; margin-bottom: 12rpx; }
.post-images { display: flex; gap: 8rpx; flex-wrap: wrap; margin-bottom: 12rpx; }
.post-image { width: 200rpx; height: 200rpx; border-radius: 8rpx; }
.post-footer { display: flex; align-items: center; gap: 24rpx; }
.post-stat { font-size: 22rpx; color: #8E9BAB; }
.delete-btn { margin-left: auto; font-size: 24rpx; color: #E74C3C; cursor: pointer; }
.empty { text-align: center; padding: 48rpx; color: #B8C2CE; }
.load-more { text-align: center; padding: 32rpx; color: #C67A6A; font-size: 26rpx; cursor: pointer; }
</style>
