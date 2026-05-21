<template>
  <view class="container">
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else>
      <view v-if="circles.length === 0" class="empty-state">
        <u-empty mode="data" text="还没有圈子，来创建一个吧" />
        <view class="empty-action">
          <u-button type="primary" shape="circle" @click="goCreate">创建圈子</u-button>
        </view>
      </view>

      <view
        v-for="(c, i) in circles"
        :key="c.id"
        :class="['circle-card', 'stagger-' + ((i % 8) + 1)]"
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
      <u-safe-bottom />
    </template>
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

.container {
  min-height: 100vh;
  background: $color-canvas;
  padding: 20rpx;
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
  color: $ink-muted;
}

.circle-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.circle-card {
  display: flex;
  align-items: center;
  gap: 24rpx;
  background: $color-surface;
  border-radius: $rounded-lg;
  padding: 28rpx;
  border: 1rpx solid $color-hairline;
  transition: transform 150ms ease-out, box-shadow 150ms ease-out;
}
.circle-card:active {
  transform: scale(0.98);
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
}

.circle-info {
  flex: 1;
  min-width: 0;
}

.circle-name {
  font-size: 30rpx;
  font-weight: 700;
  color: $ink;
  display: block;
}
.circle-desc {
  font-size: 24rpx;
  color: $ink-muted;
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
  color: $ink-tertiary;
}
.member-badge {
  font-size: 22rpx;
  color: #667eea;
  border: 1rpx solid #667eea;
  border-radius: $rounded-full;
  padding: 6rpx 20rpx;
  white-space: nowrap;
}

.empty-action {
  padding: 40rpx 60rpx;
}
.stagger-1 { animation-delay: 0ms; }
.stagger-2 { animation-delay: 60ms; }
.stagger-3 { animation-delay: 120ms; }
.stagger-4 { animation-delay: 180ms; }
.stagger-5 { animation-delay: 240ms; }
.stagger-6 { animation-delay: 300ms; }
.stagger-7 { animation-delay: 360ms; }
.stagger-8 { animation-delay: 420ms; }
