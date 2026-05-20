<template>
  <view class="container">
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="circles.length === 0">
      <u-empty mode="data" text="还没有圈子，来创建一个吧" />
      <view class="empty-action">
        <u-button type="primary" shape="circle" @click="goCreate">创建圈子</u-button>
      </view>
    </template>

    <view v-else class="circle-list">
      <view
        v-for="c in circles"
        :key="c.id"
        class="circle-card"
        @click="goDetail(c.id)"
      >
        <u-avatar
          :text="c.icon || c.name[0]"
          size="96"
          shape="square"
          fontSize="40"
          bgColor="linear-gradient(135deg, #667eea, #764ba2)"
        >
        </u-avatar>
        <view class="circle-info">
          <text class="circle-name">{{ c.name }}</text>
          <text class="circle-desc">{{ c.description || '暂无简介' }}</text>
          <view class="circle-meta">
            <text class="meta-text">{{ c.member_count }} 成员</text>
            <text class="meta-text">{{ c.post_count }} 帖子</text>
          </view>
        </view>
        <view v-if="c.is_member" class="member-badge">已加入</view>
      </view>
    </view>

    <u-safe-bottom />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCircles } from "@/api/circle"
import type { CircleData } from "@/api/circle"

const circles = ref<CircleData[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    circles.value = await listCircles()
  } catch (e: any) {
    uni.showToast({ title: "加载失败", icon: "none" })
  } finally {
    loading.value = false
  }
})

function goCreate() {
  uni.navigateTo({ url: "/pages/circle/create" })
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/circle/detail?id=${id}` })
}
</script>

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

.circle-list {
  padding: 20rpx;
}

.circle-card {
  display: flex;
  align-items: center;
  gap: 24rpx;
  background: #fff;
  border-radius: 20rpx;
  padding: 28rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.04);
}

.circle-info {
  flex: 1;
  min-width: 0;
}

.circle-name {
  font-size: 30rpx;
  font-weight: 700;
  color: #303133;
  display: block;
}

.circle-desc {
  font-size: 24rpx;
  color: #909399;
  display: block;
  margin-top: 6rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.circle-meta {
  display: flex;
  gap: 24rpx;
  margin-top: 10rpx;
}

.meta-text {
  font-size: 22rpx;
  color: #c0c4cc;
}

.member-badge {
  font-size: 22rpx;
  color: #667eea;
  border: 1rpx solid #667eea;
  border-radius: 30rpx;
  padding: 6rpx 20rpx;
  white-space: nowrap;
}

.empty-action {
  padding: 40rpx 60rpx;
}
</style>
