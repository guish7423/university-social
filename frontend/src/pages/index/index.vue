<script setup lang="ts">
import { ref, computed } from "vue"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const isLogin = computed(() => userStore.isLogin)

const features = [
  { icon: "📝", label: "发布动态", color: "#C084FC", desc: "分享校园生活", path: "/pages/post/create" },
  { icon: "👥", label: "校园圈子", color: "#F472B6", desc: "找到同好", path: "/pages/circle/list" },
  { icon: "🌿", label: "树洞", color: "#34D399", desc: "匿名倾诉", path: "/pages/whisper/index" },
  { icon: "🛍️", label: "二手市场", color: "#FBBF24", desc: "闲置交易", path: "/pages/product/list" },
  { icon: "🎯", label: "校园活动", color: "#67E8F9", desc: "精彩活动", path: "/pages/activity/list" },
  { icon: "📚", label: "课程评价", color: "#A78BFA", desc: "选课参考", path: "/pages/course/search" },
  { icon: "🔎", label: "失物招领", color: "#FB7185", desc: "互帮互助", path: "/pages/found/list" },
  { icon: "✨", label: "积分排行", color: "#F59E0B", desc: "赚取积分", path: "/pages/points/index" },
]

function go(url: string) {
  const tabs = ["/pages/square/index", "/pages/circle/list", "/pages/whisper/index", "/pages/user/index"]
  if (tabs.includes(url)) uni.switchTab({ url })
  else uni.navigateTo({ url })
}
</script>

<template>
  <view class="page-root">
    <view class="page-glow-1" />
    <view class="page-glow-2" />

    <view class="hero">
      <view class="hero-bg" />
      <view class="hero-content">
        <text class="hero-emoji animate-float">🌸</text>
        <text class="hero-title">梦园社</text>
        <text class="hero-subtitle">发现校园美好，连接有趣灵魂</text>
        <view v-if="!isLogin" class="hero-actions">
          <view class="hero-btn primary" @click="go('/pages/login/index')">开始探索</view>
        </view>
        <view v-else class="hero-greeting">
          <text class="greeting-text">欢迎回来，{{ userStore.nickname }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <view class="section-header">
        <text class="section-title">✨ 功能广场</text>
      </view>
      <view class="feature-grid">
        <view
          v-for="(f, i) in features"
          :key="i"
          :class="['feature-card glass-card', 'animate-fade-in-up', 'stagger-' + (i + 1)]"
          @click="go(f.path)"
        >
          <view class="feature-icon" :style="{ background: f.color + '20', color: f.color }">
            <text>{{ f.icon }}</text>
          </view>
          <text class="feature-label">{{ f.label }}</text>
          <text class="feature-desc">{{ f.desc }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <view class="section-header">
        <text class="section-title">🔥 热门推荐</text>
      </view>
      <view class="recommend-list">
        <view class="recommend-card glass-card animate-fade-in-up" @click="go('/pages/square/index')">
          <text class="recommend-emoji">💬</text>
          <view class="recommend-text">
            <text class="recommend-title">逛逛动态广场</text>
            <text class="recommend-desc">看看大家都在聊什么</text>
          </view>
          <text class="recommend-arrow">→</text>
        </view>
        <view class="recommend-card glass-card animate-fade-in-up stagger-2" @click="go('/pages/product/list')">
          <text class="recommend-emoji">🛍️</text>
          <view class="recommend-text">
            <text class="recommend-title">二手好物</text>
            <text class="recommend-desc">发现学长学姐的宝藏</text>
          </view>
          <text class="recommend-arrow">→</text>
        </view>
        <view class="recommend-card glass-card animate-fade-in-up stagger-3" @click="go('/pages/activity/list')">
          <text class="recommend-emoji">🎯</text>
          <view class="recommend-text">
            <text class="recommend-title">校园活动</text>
            <text class="recommend-desc">不要错过精彩的活动</text>
          </view>
          <text class="recommend-arrow">→</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page-root { min-height: 100vh; position: relative; overflow: hidden; }
.page-glow-1 {
  position: fixed; top: -400rpx; right: -200rpx;
  width: 800rpx; height: 800rpx;
  background: radial-gradient(circle, rgba(192,132,252,0.1) 0%, transparent 60%);
  pointer-events: none; z-index: 0;
}
.page-glow-2 {
  position: fixed; bottom: -300rpx; left: -200rpx;
  width: 600rpx; height: 600rpx;
  background: radial-gradient(circle, rgba(249,168,212,0.08) 0%, transparent 60%);
  pointer-events: none; z-index: 0;
}

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  box-shadow: var(--glass-shadow);
}

.hero {
  position: relative; z-index: 1;
  padding: 100rpx 40rpx 80rpx;
  text-align: center;
  overflow: hidden;
}
.hero-bg {
  position: absolute; top: 0; left: -50%; right: -50%;
  height: 100%;
  background: linear-gradient(180deg, rgba(192,132,252,0.08) 0%, transparent 100%);
  transform: skewY(-3deg);
}
.hero-content { position: relative; z-index: 1; display: flex; flex-direction: column; align-items: center; gap: 16rpx; }
.hero-emoji { font-size: 100rpx; }
.hero-title { font-family: 'ZCOOL KuaiLe', 'PingFang SC', sans-serif; font-size: 64rpx; background: var(--brand-gradient); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.hero-subtitle { font-size: 28rpx; color: var(--ink-muted); }
.hero-btn {
  margin-top: 20rpx; padding: 22rpx 64rpx; border-radius: 50rpx;
  font-size: 30rpx; font-weight: 600; color: #fff;
  transition: all 0.3s ease;
  &.primary { background: var(--brand-gradient); }
  &:active { transform: scale(0.95); }
}
.hero-greeting { margin-top: 16rpx; }
.greeting-text { font-size: 28rpx; color: var(--ink-muted); }

.section { padding: 0 32rpx 40rpx; position: relative; z-index: 1; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24rpx; }
.section-title { font-size: 32rpx; font-weight: 700; color: var(--ink); }

.feature-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20rpx;
}
.feature-card {
  display: flex; flex-direction: column; align-items: center;
  padding: 28rpx 16rpx; border-radius: 20rpx; gap: 12rpx;
  cursor: pointer; transition: all 0.3s ease;
  &:active { transform: scale(0.95); }
}
.feature-icon {
  width: 80rpx; height: 80rpx; border-radius: 20rpx;
  display: flex; align-items: center; justify-content: center;
  font-size: 36rpx;
}
.feature-label { font-size: 24rpx; font-weight: 600; color: var(--ink); }
.feature-desc { font-size: 20rpx; color: var(--ink-subtle); }

.recommend-list { display: flex; flex-direction: column; gap: 16rpx; }
.recommend-card {
  display: flex; align-items: center; gap: 24rpx;
  padding: 28rpx; border-radius: 20rpx; cursor: pointer;
  transition: all 0.3s ease;
  &:active { transform: scale(0.98); }
}
.recommend-emoji { font-size: 44rpx; }
.recommend-text { flex: 1; }
.recommend-title { font-size: 28rpx; font-weight: 600; color: var(--ink); display: block; }
.recommend-desc { font-size: 24rpx; color: var(--ink-subtle); margin-top: 4rpx; }
.recommend-arrow { font-size: 32rpx; color: var(--ink-subtle); }

@media (max-width: 480px) {
  .feature-grid { grid-template-columns: repeat(4, 1fr); }
}
</style>
