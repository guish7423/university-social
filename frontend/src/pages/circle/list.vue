<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCircles } from "@/api/circle"
import type { CircleData } from "@/api/circle"

const circles = ref<CircleData[]>([])
const loading = ref(true)

onMounted(async () => {
  try { circles.value = await listCircles() }
  catch { uni.showToast({ title: "加载失败", icon: "none" }) }
  finally { loading.value = false }
})

function goCreate() { uni.navigateTo({ url: "/pages/circle/create" }) }
function goDetail(id: number) { uni.navigateTo({ url: `/pages/circle/detail?id=${id}` }) }
function goSearch() { uni.navigateTo({ url: "/pages/search/index" }) }
</script>

<template>
  <view class="container">
    <view class="top-bar">
      <text class="top-bar-title">校园圈子</text>
      <u-icon name="search" size="36" color="#fff" @click="goSearch" />
    </view>

    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />

    </view>

    <template v-else>
      <view v-if="circles.length === 0" class="empty-state">
        <u-icon name="empty-address" size="160" color="#EAE6E0" />
        <text class="empty-title">还没有圈子</text>
        <text class="empty-desc">创建属于你们的圈子吧</text>
        <u-button type="primary" shape="circle" class="empty-btn" @click="goCreate">创建圈子</u-button>
      </view>

      <view class="circle-grid">
        <view v-for="(c, i) in circles" :key="c.id"
          class="circle-card" hover-class="circle-card-hover"
          :style="{ animationDelay: i * 0.06 + 's' }"
          @click="goDetail(c.id)">
          <view class="circle-card-header">
            <u-avatar :text="c.name[0]" size="80" shape="square"
              fontSize="36" bgColor="#C67A6A" />
            <view class="circle-card-info">
              <text class="circle-name">{{ c.name }}</text>
              <text class="circle-desc">{{ c.description || '暂无简介' }}</text>
            </view>
            <view v-if="c.is_member" class="member-badge">已加入</view>
          </view>
          <view class="circle-card-footer">
            <text class="circle-stat"><u-icon name="person" size="24" color="#B8C2CE" /> {{ c.member_count }}</text>
            <text class="circle-stat"><u-icon name="edit-pen" size="24" color="#B8C2CE" /> {{ c.post_count }}</text>
          </view>
        </view>
      </view>

      <view class="fab" @click="goCreate"><u-icon name="plus" size="40" color="#fff" /></view>
    </template>
  </view>
</template>

<style lang="scss" scoped>
.container { min-height: 100vh; background: #F7F4F0; }

.top-bar {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 20rpx 32rpx; display: flex; align-items: center; justify-content: space-between;
  height: 88rpx;
}
.top-bar-title { font-size: 32rpx; font-weight: 700; color: #fff; }

.loading-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 0; gap: 20rpx; }
.loading-text { font-size: 26rpx; color: #8E9BAB; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 60rpx; gap: 16rpx; }
.empty-title { font-size: 30rpx; font-weight: 600; color: #1E2A3A; }
.empty-desc { font-size: 26rpx; color: #8E9BAB; }
.empty-btn { margin-top: 16rpx; }

.circle-grid { display: flex; flex-direction: column; gap: 16rpx; padding: 20rpx 16rpx 160rpx; }

.circle-card {
  background: #fff; border-radius: 18rpx; padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  transition: all 0.2s ease; cursor: pointer;
  animation: fadeSlideUp 0.4s ease both;
}
.circle-card-hover { transform: scale(0.99); opacity: 0.7; }

.circle-card-header { display: flex; align-items: center; gap: 20rpx; }
.circle-card-info { flex: 1; min-width: 0; }
.circle-name { font-size: 30rpx; font-weight: 700; color: #1E2A3A; display: block; }
.circle-desc { font-size: 24rpx; color: #8E9BAB; margin-top: 8rpx; display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.member-badge { font-size: 22rpx; color: #C67A6A; border: 1rpx solid #C67A6A; border-radius: 20rpx; padding: 6rpx 20rpx; white-space: nowrap; }

.circle-card-footer { display: flex; gap: 32rpx; margin-top: 20rpx; padding-top: 20rpx; border-top: 1rpx solid #EAE6E0; }
.circle-stat { font-size: 22rpx; color: #8E9BAB; display: flex; align-items: center; gap: 6rpx; }

.fab {
  position: fixed; right: 40rpx; bottom: 120rpx;
  width: 96rpx; height: 96rpx; border-radius: 50%;
  background: linear-gradient(135deg, #C67A6A, #B06454);
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 6rpx 24rpx rgba(198,122,106,0.35);
  z-index: 100;
  &:active { transform: scale(0.92); }
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(20rpx); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
