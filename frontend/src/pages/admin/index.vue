<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getDashboard, listUsers, listPosts, listCircles, listReports } from "@/api/admin"

const stats = ref({ user_count: 0, post_count: 0, circle_count: 0 })
const recentUsers = ref<any[]>([])
const recentPosts = ref<any[]>([])
const recentReports = ref<any[]>([])

const loading = ref(true)

onMounted(async () => {
  try {
    const [s, u, p, r] = await Promise.all([
      getDashboard(),
      listUsers(0, 5),
      listPosts(0, 5),
      listReports('pending', 0, 5)
    ])
    stats.value = s
    recentUsers.value = u || []
    recentPosts.value = p || []
    recentReports.value = r || []
  } catch {}
  loading.value = false
})
</script>

<template>
  <view class="page">
    <view class="page-header">
      <text class="page-title">管理后台</text>
      <text class="page-subtitle">校园社运营概览</text>
    </view>

    <view class="stats-grid">
      <view class="stat-card">
        <text class="stat-value">{{ stats.user_count }}</text>
        <text class="stat-label">用户数</text>
      </view>
      <view class="stat-card">
        <text class="stat-value">{{ stats.post_count }}</text>
        <text class="stat-label">帖子数</text>
      </view>
      <view class="stat-card">
        <text class="stat-value">{{ stats.circle_count }}</text>
        <text class="stat-label">圈子数</text>
      </view>
      <view class="stat-card" @click="uni.navigateTo({url:'/pages/admin/reports'})">
        <text class="stat-value warn">{{ recentReports.length }}</text>
        <text class="stat-label">待处理举报</text>
      </view>
    </view>

    <view class="section">
      <view class="section-header" @click="uni.navigateTo({url:'/pages/admin/users'})">
        <text class="section-title">最新用户</text>
        <text class="section-more">查看全部 ›</text>
      </view>
      <view class="list">
        <view v-for="u in recentUsers" :key="u.id" class="list-item" @click="uni.navigateTo({url:'/pages/user/detail?id='+u.id})">
          <u-avatar :src="u.avatar" :text="u.nickname?.[0] || '?'" size="72" shape="circle" />
          <view class="list-item-body">
            <text class="list-item-title">{{ u.nickname }}</text>
            <text class="list-item-desc">{{ u.school || '未设置' }}</text>
          </view>
          <text class="list-item-meta">{{ u.created_at?.slice(0,10) }}</text>
        </view>
        <view v-if="recentUsers.length === 0" class="empty">暂无用户</view>
      </view>
    </view>

    <view class="section">
      <view class="section-header">
        <text class="section-title">最新帖子</text>
      </view>
      <view class="list">
        <view v-for="p in recentPosts" :key="p.id" class="list-item">
          <view class="list-item-body">
            <text class="list-item-title">{{ p.author?.nickname || '#' + p.user_id }}</text>
            <text class="list-item-desc">{{ p.content?.slice(0,60) }}</text>
          </view>
        </view>
        <view v-if="recentPosts.length === 0" class="empty">暂无帖子</view>
      </view>
    </view>

    <view class="nav-grid">
      <view class="nav-item" @click="uni.navigateTo({url:'/pages/admin/users'})">
        <text class="nav-icon">👥</text>
        <text class="nav-label">用户管理</text>
      </view>
      <view class="nav-item" @click="uni.navigateTo({url:'/pages/admin/posts'})">
        <text class="nav-icon">📝</text>
        <text class="nav-label">帖子管理</text>
      </view>
      <view class="nav-item" @click="uni.navigateTo({url:'/pages/admin/circles'})">
        <text class="nav-icon">🔄</text>
        <text class="nav-label">圈子管理</text>
      </view>
      <view class="nav-item" @click="uni.navigateTo({url:'/pages/admin/reports'})">
        <text class="nav-icon">🚩</text>
        <text class="nav-label">举报管理</text>
      </view>
      <view class="nav-item" @click="uni.navigateTo({url:'/pages/admin/sensitive'})">
        <text class="nav-icon">🔤</text>
        <text class="nav-label">敏感词</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding-bottom: 40rpx; }
.page-header {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 60rpx 32rpx 40rpx;
}
.page-title { font-size: 36rpx; font-weight: 600; color: #fff; }
.page-subtitle { font-size: 24rpx; color: rgba(255,255,255,0.5); margin-top: 8rpx; display: block; }

.stats-grid { display: flex; gap: 16rpx; padding: 24rpx 16rpx; }
.stat-card {
  flex: 1; background: #fff; border-radius: 14rpx; padding: 20rpx;
  text-align: center; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
}
.stat-value { font-size: 40rpx; font-weight: 700; color: #1E2A3A; display: block; }
.stat-value.warn { color: #E67E22; }
.stat-label { font-size: 22rpx; color: #8E9BAB; margin-top: 8rpx; display: block; }

.section { margin: 0 16rpx 24rpx; }
.section-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 0 8rpx 16rpx; cursor: pointer;
}
.section-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }
.section-more { font-size: 24rpx; color: #C67A6A; }

.list { background: #fff; border-radius: 14rpx; overflow: hidden; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); }
.list-item {
  display: flex; align-items: center; gap: 16rpx;
  padding: 24rpx; border-bottom: 1rpx solid #EAE6E0; cursor: pointer;
  &:last-child { border-bottom: none; }
}
.list-item-body { flex: 1; min-width: 0; }
.list-item-title { font-size: 26rpx; color: #1E2A3A; font-weight: 500; display: block; }
.list-item-desc { font-size: 22rpx; color: #8E9BAB; margin-top: 4rpx; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.list-item-meta { font-size: 20rpx; color: #B8C2CE; }
.empty { text-align: center; padding: 32rpx; color: #B8C2CE; font-size: 24rpx; }

.nav-grid {
  display: grid; grid-template-columns: repeat(5, 1fr); gap: 16rpx;
  padding: 0 16rpx;
}
.nav-item {
  background: #fff; border-radius: 14rpx; padding: 20rpx 8rpx;
  text-align: center; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); cursor: pointer;
}
.nav-icon { font-size: 40rpx; display: block; margin-bottom: 8rpx; }
.nav-label { font-size: 22rpx; color: #1E2A3A; }
</style>
