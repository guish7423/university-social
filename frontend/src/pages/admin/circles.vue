<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCircles, deleteCircle } from "@/api/admin"

const circles = ref<any[]>([])
const loading = ref(true)
const offset = ref(0)
const hasMore = ref(true)

async function load() {
  try {
    const data = await listCircles(offset.value)
    if (data.length < 50) hasMore.value = false
    circles.value = offset.value === 0 ? data : [...circles.value, ...data]
  } catch (e) { console.error(e) }
  loading.value = false
}

onMounted(load)

async function doDelete(id: number) {
  uni.showModal({
    title: "确认删除",
    content: "确定要删除这个圈子吗？",
    success: async (res) => {
      if (res.confirm) {
        uni.showLoading({ title: "删除中" })
        try {
          await deleteCircle(id)
          circles.value = circles.value.filter(c => c.id !== id)
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
      <text class="page-title">圈子管理 ({{ circles.length }})</text>
    </view>

    <view class="list">
      <view v-for="c in circles" :key="c.id" class="list-item">
        <view class="list-item-body">
          <text class="list-item-title">{{ c.name }}</text>
          <text class="list-item-desc">{{ c.description }}</text>
        </view>
        <view class="list-item-stats">
          <text class="stat">👥 {{ c.member_count }}</text>
          <text class="stat">📝 {{ c.post_count }}</text>
        </view>
        <text class="delete-btn" @click="doDelete(c.id)">删除</text>
      </view>
      <view v-if="circles.length === 0" class="empty">暂无圈子</view>
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
.list-item {
  display: flex; align-items: center; gap: 16rpx;
  padding: 24rpx; border-bottom: 1rpx solid #EAE6E0;
  &:last-child { border-bottom: none; }
}
.list-item-body { flex: 1; min-width: 0; }
.list-item-title { font-size: 26rpx; color: #1E2A3A; font-weight: 500; display: block; }
.list-item-desc { font-size: 22rpx; color: #8E9BAB; margin-top: 4rpx; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.list-item-stats { text-align: center; }
.stat { font-size: 20rpx; color: #8E9BAB; display: block; }
.delete-btn { font-size: 24rpx; color: #E74C3C; cursor: pointer; padding: 8rpx; }
.empty { text-align: center; padding: 48rpx; color: #B8C2CE; }
.load-more { text-align: center; padding: 32rpx; color: #C67A6A; font-size: 26rpx; cursor: pointer; }
</style>
