<template>
  <view class="container">
    <view class="search-bar">
      <u-icon name="search" size="32" color="#909399"></u-icon>
      <input
        v-model="query"
        class="search-input"
        placeholder="搜索帖子、圈子、用户..."
        :adjust-position="false"
        maxlength="50"
        confirm-type="search"
        @confirm="doSearch"
      />
      <text v-if="query" class="clear-btn" @click="query = ''; results = null">清除</text>
    </view>

    <view class="tabs">
      <view
        v-for="t in tabs" :key="t.key"
        class="tab"
        :class="activeTab === t.key && 'tab-active'"
        @click="switchTab(t.key)"
      >{{ t.label }}</view>
    </view>

    <!-- 搜索提示 -->
    <view v-if="!query && !trendingPosts.length" class="placeholder">
      <text class="placeholder-title">热门趋势</text>
      <text class="placeholder-desc">搜索校园新鲜事</text>
    </view>

    <!-- 搜索结果：帖子 -->
    <view v-if="activeTab === 'posts'" class="result-list">
      <view v-if="loading" class="loading"><u-loading mode="flower" size="40" /></view>
      <view v-else-if="!posts.length" class="empty">没有找到相关帖子</view>
      <view
        v-for="(p, i) in posts" :key="p.id"
        class="card stagger-{{ i + 1 }}"
        @click="goPost(p.id)"
      >
        <view class="card-header">
          <u-avatar :src="p.author?.avatar" size="60" shape="circle" />
          <view class="card-author">
            <text class="card-name">{{ p.author?.nickname || '匿名' }}</text>
            <text class="card-school">{{ p.school }}</text>
          </view>
          <text class="card-time">{{ formatTime(p.created_at) }}</text>
        </view>
        <text class="card-content">{{ p.content }}</text>
        <view v-if="p.images?.length" class="card-images">
          <image v-for="(img, j) in p.images.slice(0, 3)" :key="j" :src="img" mode="aspectFill" class="card-img" />
        </view>
        <view class="card-actions">
          <text>👍 {{ p.like_count }}</text>
          <text>💬 {{ p.comment_count }}</text>
        </view>
      </view>
    </view>

    <!-- 搜索结果：圈子 -->
    <view v-if="activeTab === 'circles'" class="result-list">
      <view v-if="loading" class="loading"><u-loading mode="flower" size="40" /></view>
      <view v-else-if="!circles.length" class="empty">没有找到相关圈子</view>
      <view
        v-for="c in circles" :key="c.id"
        class="card"
        @click="goCircle(c.id)"
      >
        <u-avatar :text="c.icon || c.name[0]" size="72" shape="square" />
        <view class="card-info">
          <text class="card-name">{{ c.name }}</text>
          <text class="card-desc">{{ c.description || '暂无简介' }}</text>
          <view class="card-meta">
            <text>{{ c.member_count }} 成员 · {{ c.post_count }} 帖子</text>
            <text v-if="c.is_member" class="member-badge">已加入</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 搜索结果：用户 -->
    <view v-if="activeTab === 'users'" class="result-list">
      <view v-if="loading" class="loading"><u-loading mode="flower" size="40" /></view>
      <view v-else-if="!users.length" class="empty">没有找到相关用户</view>
      <view v-for="u in users" :key="u.id" class="card user-card">
        <u-avatar :src="u.avatar" size="72" shape="circle" />
        <view class="card-info">
          <text class="card-name">{{ u.nickname }}</text>
          <text class="card-desc">{{ u.school || '未知学校' }}</text>
        </view>
        <text class="follow-btn" :class="u.is_friend && 'followed'" @click="toggleFollow(u)">
          {{ u.is_friend ? '好友' : '加好友' }}
        </text>
      </view>
    </view>

    <!-- 热门帖子 -->
    <view v-if="!query && trendingPosts.length" class="trending-section">
      <view class="section-header">
        <text class="section-title">🔥 今日热议</text>
      </view>
      <view
        v-for="(p, i) in trendingPosts" :key="p.id"
        class="trending-item stagger-{{ i + 1 }}"
        @click="goPost(p.id)"
      >
        <text class="trending-rank">{{ i + 1 }}</text>
        <view class="trending-body">
          <text class="trending-content">{{ p.content }}</text>
          <view class="trending-meta">
            <text>👍 {{ p.like_count }}</text>
            <text>💬 {{ p.comment_count }}</text>
          </view>
        </view>
      </view>
    </view>

    <u-safe-bottom />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { searchPosts, trendingPosts, listPosts } from "@/api/post"
import { searchCircles } from "@/api/circle"
import { searchUsers, sendFriendRequest } from "@/api/social"
import type { PostData } from "@/api/post"
import type { CircleData } from "@/api/circle"
import type { UserInfo } from "@/api/user"

const query = ref("")
const activeTab = ref("posts")
const loading = ref(false)
const posts = ref<PostData[]>([])
const circles = ref<CircleData[]>([])
const users = ref<(UserInfo & { is_friend?: boolean })[]>([])
const trendingPosts = ref<PostData[]>([])

const tabs = [
  { key: "posts", label: "帖子" },
  { key: "circles", label: "圈子" },
  { key: "users", label: "用户" },
]

onMounted(async () => {
  try { trendingPosts.value = await trendingPosts() } catch {}
})

async function doSearch() {
  if (!query.value.trim()) return
  loading.value = true
  try {
    const q = query.value.trim()
    if (activeTab.value === "posts") posts.value = await searchPosts(q)
    else if (activeTab.value === "circles") circles.value = await searchCircles(q)
    else if (activeTab.value === "users") users.value = await searchUsers(q)
  } catch { uni.showToast({ title: "搜索失败", icon: "none" }) }
  finally { loading.value = false }
}

function switchTab(key: string) {
  activeTab.value = key
  if (query.value.trim()) doSearch()
}

function goPost(id: number) { uni.navigateTo({ url: `/pages/post/detail?id=${id}` }) }
function goCircle(id: number) { uni.navigateTo({ url: `/pages/circle/detail?id=${id}` }) }

async function toggleFollow(u: UserInfo & { is_friend?: boolean }) {
  if (u.is_friend) return
  try {
    await sendFriendRequest(u.id)
    u.is_friend = true
    uni.showToast({ title: "已发送好友请求", icon: "success" })
  } catch { uni.showToast({ title: "操作失败", icon: "none" }) }
}

function formatTime(t: string) {
  const d = new Date(t); const now = new Date(); const diff = (now.getTime() - d.getTime()) / 1000
  if (diff < 60) return "刚刚"
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<style scoped>
.container { min-height: 100vh; background: $color-canvas; padding: 20rpx; }

.search-bar {
  display: flex; align-items: center; gap: 16rpx;
  background: $color-surface; border-radius: $rounded-full;
  padding: 16rpx 28rpx; margin-bottom: 20rpx;
  border: 1rpx solid $color-hairline;
}
.search-input { flex: 1; font-size: 28rpx; color: $ink; }
.clear-btn { font-size: 24rpx; color: $ink-muted; }

.tabs { display: flex; gap: 8rpx; margin-bottom: 20rpx; }
.tab {
  flex: 1; text-align: center; padding: 16rpx 0;
  font-size: 26rpx; color: $ink-muted;
  border-radius: $rounded-md; background: $color-surface;
}
.tab-active { color: #fff; background: linear-gradient(135deg, #667eea, #764ba2); font-weight: 600; }

.placeholder { text-align: center; padding: 100rpx 0; }
.placeholder-title { font-size: 32rpx; font-weight: 700; color: $ink; display: block; }
.placeholder-desc { font-size: 26rpx; color: $ink-muted; margin-top: 10rpx; }

.loading { text-align: center; padding: 100rpx 0; }
.empty { text-align: center; padding: 100rpx 0; color: $ink-muted; font-size: 26rpx; }

.result-list { display: flex; flex-direction: column; gap: 16rpx; }

.card {
  background: $color-surface; border-radius: $rounded-lg; padding: 24rpx;
  border: 1rpx solid $color-hairline; transition: transform 150ms ease-out;
}
.card:active { transform: scale(0.98); }

.card-header { display: flex; align-items: center; gap: 16rpx; margin-bottom: 12rpx; }
.card-author { flex: 1; }
.card-name { font-size: 26rpx; font-weight: 600; color: $ink; display: block; }
.card-school { font-size: 22rpx; color: $ink-muted; }
.card-time { font-size: 22rpx; color: $ink-tertiary; white-space: nowrap; }
.card-content { font-size: 28rpx; line-height: 1.6; color: $ink; margin-bottom: 12rpx; }
.card-images { display: flex; gap: 8rpx; margin-bottom: 12rpx; }
.card-img { width: 200rpx; height: 200rpx; border-radius: $rounded-md; background: $color-hairline; }
.card-actions { display: flex; gap: 24rpx; font-size: 24rpx; color: $ink-muted; padding-top: 12rpx; border-top: 1rpx solid $color-hairline; }

.card-info { flex: 1; }
.card-desc { font-size: 24rpx; color: $ink-muted; margin-top: 4rpx; display: block; }
.card-meta { font-size: 22rpx; color: $ink-tertiary; margin-top: 8rpx; display: flex; gap: 12rpx; align-items: center; }
.member-badge { font-size: 20rpx; color: #667eea; border: 1rpx solid #667eea; border-radius: $rounded-full; padding: 4rpx 16rpx; }

.user-card { display: flex; align-items: center; gap: 20rpx; }
.follow-btn {
  font-size: 24rpx; color: #667eea; border: 1rpx solid #667eea;
  border-radius: $rounded-full; padding: 8rpx 24rpx; white-space: nowrap;
}
.follow-btn.followed { color: $ink-muted; border-color: $color-hairline; }

.trending-section { margin-top: 20rpx; }
.section-header { margin-bottom: 16rpx; }
.section-title { font-size: 30rpx; font-weight: 700; color: $ink; }

.trending-item {
  display: flex; gap: 16rpx; padding: 20rpx 0;
  border-bottom: 1rpx solid $color-hairline;
}
.trending-rank {
  width: 48rpx; height: 48rpx; border-radius: $rounded-md;
  background: linear-gradient(135deg, #667eea, #764ba2); color: #fff;
  text-align: center; line-height: 48rpx; font-size: 24rpx; font-weight: 700;
}
.trending-body { flex: 1; }
.trending-content { font-size: 26rpx; color: $ink; display: block; margin-bottom: 6rpx; }
.trending-meta { display: flex; gap: 20rpx; font-size: 22rpx; color: $ink-tertiary; }
</style>
