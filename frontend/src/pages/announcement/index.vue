<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listAnnouncements, type AnnouncementData } from "@/api/social"

const announcements = ref<AnnouncementData[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    announcements.value = await listAnnouncements()
  } catch (e) { console.error(e) }
  loading.value = false
})

function viewDetail(a: AnnouncementData) {
  uni.showModal({ title: a.title, content: a.content, showCancel: false })
}

function formatTime(t: string) {
  if (!t) return ""
  const d = new Date(t)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`
}
</script>

<template>
  <view class="container">
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />

    </view>
    <template v-else-if="announcements.length">
      <view v-for="a in announcements" :key="a.id" class="announcement-card" @click="viewDetail(a)">
        <view class="card-header">
          <view class="card-tag">公告</view>
          <text class="card-title">{{ a.title }}</text>
        </view>
        <text class="card-summary">{{ a.content?.slice(0, 100) }}{{ a.content?.length > 100 ? "..." : "" }}</text>
        <text class="card-time">{{ formatTime(a.created_at) }}</text>
      </view>
    </template>
    <view v-else class="empty-state">
      <u-icon name="info-circle" size="120" color="#EAE6E0" />
      <text class="empty-text">暂无公告</text>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.container { min-height: 100vh; background: #F7F4F0; padding: 24rpx; }
.loading-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 0; gap: 16rpx; }
.loading-text { font-size: 26rpx; color: #B8C2CE; }
.announcement-card {
  background: #fff; border-radius: 16rpx; padding: 28rpx; margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); cursor: pointer;
  transition: all 0.15s ease;
  &:active { transform: scale(0.99); }
}
.card-header { display: flex; align-items: center; gap: 16rpx; margin-bottom: 12rpx; }
.card-tag {
  font-size: 20rpx; color: #fff; background: #C67A6A;
  padding: 4rpx 14rpx; border-radius: 6rpx; flex-shrink: 0;
}
.card-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; flex: 1; }
.card-summary { font-size: 26rpx; color: #5C6B7E; line-height: 1.6; display: block; margin-bottom: 12rpx; }
.card-time { font-size: 22rpx; color: #B8C2CE; }
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 0; gap: 16rpx; }
.empty-text { font-size: 26rpx; color: #B8C2CE; }
</style>
