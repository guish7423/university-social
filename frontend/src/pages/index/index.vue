<script setup lang="ts">
import { ref, onMounted } from "vue"
import { trendingPosts, listTopics, type PostData } from "@/api/post"
import { listBanners, listAnnouncements, type BannerData, type AnnouncementData } from "@/api/social"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const trending = ref<PostData[]>([])
const topics = ref<Topic[]>([])
const loading = ref(true)
const banners = ref<BannerData[]>([])
const announcements = ref<AnnouncementData[]>([])

const features = [
  { icon: "moment", title: "动态广场", desc: "发现校园新鲜事", path: "/pages/square/index", tab: true },
  { icon: "community", title: "校园圈子", desc: "找到志同道合的伙伴", path: "/pages/circle/list", tab: true },
  { icon: "favor", title: "树洞", desc: "匿名倾诉你的心声", path: "/pages/whisper/index", tab: true },
  { icon: "gift", title: "二手市场", desc: "闲置物品交易", path: "/pages/product/list" },
  { icon: "calendar", title: "校园活动", desc: "参与精彩活动", path: "/pages/activity/list" },
  { icon: "star", title: "课程评价", desc: "选课参考指南", path: "/pages/course/search" },
  { icon: "search", title: "失物招领", desc: "找回丢失物品", path: "/pages/found/list" },
  { icon: "integral", title: "积分排行", desc: "赚积分赢奖励", path: "/pages/points/index" },
  { icon: "clock", title: "校历", desc: "学期安排一览", path: "/pages/campus/calendar" },
  { icon: "phone", title: "校园黄页", desc: "机构/老师联系", path: "/pages/campus/directory" },
  { icon: "map", title: "空教室", desc: "自习室查询", path: "/pages/campus/rooms" },
]


onMounted(async () => {
  try {
    const [p, t, b, a] = await Promise.all([
      trendingPosts(), listTopics(),
      listBanners(), listAnnouncements(),
    ])
    trending.value = p || []
    topics.value = t || []
    banners.value = b || []
    announcements.value = a || []
  } catch {}
  loading.value = false
})

function go(path: string, tab = false) {
  if (tab) uni.switchTab({ url: path })
  else uni.navigateTo({ url: path })
}

function formatTime(t: string) {
  if (!t) return ""
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return Math.floor(diff / 60000) + "分钟前"
  if (diff < 86400000) return Math.floor(diff / 3600000) + "小时前"
  const days = Math.floor(diff / 86400000)
  if (days < 7) return days + "天前"
  return (d.getMonth() + 1) + "/" + d.getDate()
}
</script>

<template>
  <view class="page">
    <view class="hero-section">
      <view class="hero-bg" />
      <view class="hero-content">
        <view class="hero-text">
          <text class="hero-title">校园社</text>
          <text class="hero-subtitle">连接校园里的每一个人</text>
        </view>
        <view class="hero-stats">
          <view class="hero-stat-item">
            <text class="hero-stat-num">8</text>
            <text class="hero-stat-label">功能模块</text>
          </view>
          <view class="hero-stat-divider" />
          <view class="hero-stat-item">
            <text class="hero-stat-num">∞</text>
            <text class="hero-stat-label">无限可能</text>
          </view>
        </view>
      </view>
    </view>

    <view class="body">
      <view v-if="banners.length" class="banner-section">
        <u-swiper :list="banners.map(b => b.image_url)" height="320rpx" borderRadius="16" :autoplay="true" />
      </view>
      <view class="section">
        <view class="section-header">
          <text class="section-title">全部功能</text>
        </view>
        <view class="feature-grid">
          <view v-for="(f, i) in features" :key="i"
            class="feature-card" hover-class="feature-card-hover"
            @click="go(f.path, f.tab)"
          >
            <view class="feature-icon-box">
              <u-icon :name="f.icon" size="40" color="#C67A6A" />
            </view>
            <text class="feature-title">{{ f.title }}</text>
            <text class="feature-desc">{{ f.desc }}</text>
          </view>
        </view>
      </view>

      <view v-if="trending.length" class="section">
        <view class="section-header">
          <text class="section-title">热门动态</text>
          <text class="section-more" @click="go('/pages/square/index', true)">查看全部 ›</text>
        </view>
        <view class="post-list">
          <view v-for="post in trending.slice(0, 3)" :key="post.id"
            class="post-card" hover-class="post-card-hover"
            @click="go('/pages/post/detail?id=' + post.id)"
          >
            <view class="post-card-header">
              <u-avatar :src="post.author?.avatar" :text="post.author?.nickname?.[0] || '?'" size="56" shape="circle" />
              <view class="post-card-author">
                <text class="post-card-name">{{ post.author?.nickname || '匿名' }}</text>
                <text class="post-card-time">{{ formatTime(post.created_at) }}</text>
              </view>
            </view>
            <text class="post-card-content">{{ post.content }}</text>
            <view v-if="post?.images?.length" class="post-card-images">
              <image v-for="(img, j) in post.images.slice(0, 3)" :key="j"
                :src="img" mode="aspectFill" class="post-card-image"
              />
            </view>
            <view class="post-card-footer">
              <view class="post-card-stat">
                <u-icon name="heart" size="28" color="#B8C2CE" />
                <text class="stat-text">{{ post.like_count || 0 }}</text>
              </view>
              <view class="post-card-stat">
                <u-icon name="chat" size="28" color="#B8C2CE" />
                <text class="stat-text">{{ post.comment_count || 0 }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <view v-if="trending.length === 0" class="section">
        <view class="section-header">
          <text class="section-title">热门话题</text>
        </view>
        <view class="topic-list">
          <view v-for="t in topics.slice(0, 6)" :key="t.id" class="topic-chip" @click="go('/pages/square/index', true)">
            <u-icon :name="t.icon || 'tags-fill'" size="28" color="#C67A6A" />
            <text class="topic-name">{{ t.name }}</text>
          </view>
          <view v-if="topics.length === 0" class="topic-chip">
            <u-icon name="tags-fill" size="28" color="#C67A6A" />
            <text class="topic-name">暂无话题</text>
          </view>
        </view>
      </view>

      <view v-if="announcements.length" class="section">
        <view class="section-header">
          <text class="section-title">📢 校园公告</text>
          <text class="section-more" @click="go('/pages/announcement/index')">查看全部 ›</text>
        </view>
        <view class="announcement-list">
          <view v-for="a in announcements.slice(0, 3)" :key="a.id" class="announcement-card" @click="uni.showModal({ title: a.title, content: a.content, showCancel: false })">
            <view class="announcement-dot" />
            <text class="announcement-title">{{ a.title }}</text>
          </view>
        </view>
      </view>
      <view v-if="!userStore.isLogin" class="login-banner">
        <view class="login-banner-content">
          <text class="login-banner-title">登录后解锁全部功能</text>
          <text class="login-banner-desc">参与互动，发现更多精彩内容</text>
        </view>
        <u-button type="primary" shape="circle" size="medium" @click="go('/pages/login/index')">立即登录</u-button>
      </view>
    </view>

    <u-safe-bottom />
  </view>
</template>

<style lang="scss">
page { background: #F7F4F0; }
</style>

<style lang="scss" scoped>
.hero-section {
  position: relative;
  padding: 100rpx 40rpx 60rpx;
  overflow: hidden;
}
.hero-bg {
  position: absolute; inset: 0;
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 50%, #3A4A5E 100%);
  z-index: 0;
  &::after {
    content: ''; position: absolute;
    top: -200rpx; right: -100rpx;
    width: 500rpx; height: 500rpx;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(198,122,106,0.12) 0%, transparent 70%);
  }
}
.hero-content { position: relative; z-index: 1; display: flex; flex-direction: column; gap: 36rpx; }
.hero-text { display: flex; flex-direction: column; gap: 8rpx; }
.hero-title { font-size: 44rpx; font-weight: 700; color: #C67A6A; letter-spacing: 4rpx; }
.hero-subtitle { font-size: 26rpx; color: rgba(255,255,255,0.5); font-weight: 300; letter-spacing: 2rpx; }
.hero-stats { display: flex; align-items: center; gap: 32rpx; }
.hero-stat-item { display: flex; flex-direction: column; gap: 4rpx; }
.hero-stat-num { font-size: 36rpx; font-weight: 700; color: #fff; }
.hero-stat-label { font-size: 22rpx; color: rgba(255,255,255,0.4); }
.hero-stat-divider { width: 1rpx; height: 48rpx; background: rgba(255,255,255,0.1); }

.body { padding: 0 24rpx 40rpx; position: relative; z-index: 2; margin-top: -20rpx; }

.section { margin-bottom: 32rpx; }
.section-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 0 8rpx; margin-bottom: 20rpx;
}
.section-title { font-size: 30rpx; font-weight: 700; color: #1E2A3A; }
.section-more { font-size: 24rpx; color: #C67A6A; }

.feature-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16rpx;
}
.feature-card {
  background: #fff; border-radius: 16rpx; padding: 28rpx 12rpx;
  display: flex; flex-direction: column; align-items: center; gap: 12rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  transition: all 0.2s ease; cursor: pointer;
  &:active { transform: scale(0.96); }
}
.feature-card-hover { transform: scale(0.96); }
.feature-icon-box {
  width: 64rpx; height: 64rpx;
  border-radius: 16rpx;
  background: rgba(198,122,106,0.08);
  display: flex; align-items: center; justify-content: center;
}
.feature-title { font-size: 24rpx; font-weight: 600; color: #1E2A3A; }
.feature-desc { font-size: 20rpx; color: #8E9BAB; text-align: center; }

.post-list { display: flex; flex-direction: column; gap: 16rpx; }
.post-card {
  background: #fff; border-radius: 14rpx; padding: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  cursor: pointer; transition: all 0.2s ease;
}
.post-card-hover { transform: scale(0.99); }
.post-card-header { display: flex; align-items: center; gap: 12rpx; margin-bottom: 12rpx; }
.post-card-author { flex: 1; }
.post-card-name { font-size: 26rpx; font-weight: 600; color: #1E2A3A; display: block; }
.post-card-time { font-size: 20rpx; color: #B8C2CE; margin-top: 4rpx; display: block; }
.post-card-content { font-size: 26rpx; color: #5C6B7E; line-height: 1.6; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
.post-card-images { display: flex; gap: 8rpx; margin-top: 12rpx; }
.post-card-image { width: 210rpx; height: 210rpx; border-radius: 10rpx; background: #E8E4DE; }
.post-card-footer { display: flex; gap: 32rpx; margin-top: 16rpx; padding-top: 16rpx; border-top: 1rpx solid #EAE6E0; }
.post-card-stat { display: flex; align-items: center; gap: 6rpx; }
.stat-text { font-size: 22rpx; color: #B8C2CE; }

.topic-list { display: flex; flex-wrap: wrap; gap: 12rpx; }
.topic-chip {
  background: #fff; border-radius: 20rpx; padding: 12rpx 24rpx;
  display: flex; align-items: center; gap: 8rpx;
  box-shadow: 0 2rpx 8rpx rgba(30,42,58,0.03); cursor: pointer;
}

.login-banner {
  background: linear-gradient(135deg, rgba(198,122,106,0.08), rgba(198,122,106,0.04));
  border-radius: 16rpx; padding: 28rpx 24rpx;
  display: flex; align-items: center; justify-content: space-between;
  border: 1rpx solid rgba(198,122,106,0.15);
}
.login-banner-content { display: flex; flex-direction: column; gap: 6rpx; }
.login-banner-title { font-size: 26rpx; font-weight: 600; color: #1E2A3A; }
.login-banner-desc { font-size: 22rpx; color: #8E9BAB; }

@media (min-width: 1024px) {
  .body { max-width: 1200rpx; margin: -20rpx auto 0; }
  .feature-grid { grid-template-columns: repeat(4, 1fr); }
}
.banner-section { margin-bottom: 24rpx; margin-top: -12rpx; }
.announcement-list { display: flex; flex-direction: column; gap: 12rpx; }
.announcement-card {
  display: flex; align-items: center; gap: 12rpx;
  padding: 16rpx 20rpx; background: #fff; border-radius: 12rpx;
  box-shadow: 0 2rpx 8rpx rgba(30,42,58,0.03); cursor: pointer;
  transition: all 0.15s ease;
  &:active { transform: scale(0.99); }
}
.announcement-dot {
  width: 8rpx; height: 8rpx; border-radius: 50%; background: #C67A6A;
  flex-shrink: 0;
}
.announcement-title { font-size: 24rpx; color: #5C6B7E; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
