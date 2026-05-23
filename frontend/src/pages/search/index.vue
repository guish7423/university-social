<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue"
import { searchPosts, trendingPosts, type PostData } from "@/api/post"
import { searchCircles } from "@/api/circle"
import { searchUsers, sendFriendRequest } from "@/api/social"
import type { CircleData } from "@/api/circle"
import type { UserInfo } from "@/api/user"

const query = ref("")
const activeTab = ref("posts")
const loading = ref(false)
const posts = ref<PostData[]>([])
const circles = ref<CircleData[]>([])
const users = ref<(UserInfo & { is_friend?: boolean })[]>([])
const trending = ref<PostData[]>([])
const searchInput = ref<InstanceType<typeof HTMLInputElement> | null>(null)

const tabs = [
  { key: "posts", label: "帖子" },
  { key: "circles", label: "圈子" },
  { key: "users", label: "用户" },
]

onMounted(async () => {
  try { trending.value = await trendingPosts() } catch (e) { console.error(e) }
  nextTick(() => {
    const input = document.querySelector('.search-input') as HTMLInputElement
    if (input) input.focus()
  })
})

async function doSearch() {
  const q = query.value.trim()
  if (!q) return
  loading.value = true
  try {
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
function goUser(id: number) { uni.navigateTo({ url: `/pages/user/detail?id=${id}` }) }

async function toggleFollow(u: UserInfo & { is_friend?: boolean }) {
  if (u.is_friend) return
  try {
    await sendFriendRequest(u.id)
    u.is_friend = true
    uni.showToast({ title: "已发送好友请求", icon: "success" })
  } catch { uni.showToast({ title: "操作失败", icon: "none" }) }
}

function formatTime(t: string) {
  if (!t) return ""
  const d = new Date(t), now = new Date(), diff = (now.getTime() - d.getTime()) / 1000
  if (diff < 60) return "刚刚"
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <view class="search-box" @click.stop>
        <u-icon name="search" size="32" color="#8E9BAB" />
        <input
          v-model="query"
          class="search-input"
          placeholder="搜索帖子、圈子、用户..."
          :adjust-position="false"
          maxlength="50"
          confirm-type="search"
          @confirm="doSearch"
        />
        <text v-if="query" class="clear-btn" @click="query = ''">清除</text>
        <text class="search-btn" @click="doSearch">搜索</text>
      </view>
    </view>

    <view class="tabs">
      <view v-for="t in tabs" :key="t.key"
        :class="['tab', activeTab === t.key && 'tab-active']"
        @click="switchTab(t.key)">
        {{ t.label }}
      </view>
    </view>

    <scroll-view scroll-y class="scroll-area" @scrolltolower="doSearch">
      <view v-if="loading" class="loading"><u-loading mode="flower" size="40" /></view>

      <template v-else-if="query.trim()">
        <view v-if="activeTab === 'posts'" class="result-list">
          <view v-for="p in posts" :key="p.id"
            class="result-card" hover-class="result-card-hover"
            @click="goPost(p.id)">
            <view class="card-header">
              <u-avatar :src="p.author?.avatar" :text="p.author?.nickname?.[0] || '?'" size="56" shape="circle" />
              <view class="card-author">
                <text class="card-name">{{ p.author?.nickname || '匿名' }}</text>
                <text class="card-school">{{ p.school || '' }}</text>
              </view>
              <text class="card-time">{{ formatTime(p.created_at) }}</text>
            </view>
            <text class="card-content">{{ p.content }}</text>
            <view v-if="p.images?.length" class="card-images">
              <image v-for="(img, j) in p.images.slice(0, 3)" :key="j" :src="img" mode="aspectFill" class="card-img" />
            </view>
            <view class="card-actions">
              <u-icon name="heart" size="28" color="#B8C2CE" />
              <text class="action-txt">{{ p.like_count || 0 }}</text>
              <u-icon name="chat" size="28" color="#B8C2CE" />
              <text class="action-txt">{{ p.comment_count || 0 }}</text>
            </view>
          </view>
          <view v-if="posts.length === 0" class="empty">没有找到相关帖子</view>
        </view>

        <view v-if="activeTab === 'circles'" class="result-list">
          <view v-for="c in circles" :key="c.id"
            class="result-card row" hover-class="result-card-hover"
            @click="goCircle(c.id)">
            <u-avatar :text="c.name?.[0] || 'C'" size="72" shape="square" />
            <view class="card-info">
              <text class="card-name">{{ c.name }}</text>
              <text class="card-desc">{{ c.description || '暂无简介' }}</text>
              <view class="card-meta">
                <text>{{ c.member_count }} 成员 · {{ c.post_count }} 帖子</text>
                <text v-if="c.is_member" class="badge">已加入</text>
              </view>
            </view>
          </view>
          <view v-if="circles.length === 0" class="empty">没有找到相关圈子</view>
        </view>

        <view v-if="activeTab === 'users'" class="result-list">
          <view v-for="u in users" :key="u.id" class="result-card row" hover-class="result-card-hover">
            <u-avatar :src="u.avatar" :text="u.nickname?.[0] || '?'" size="72" shape="circle" @click="goUser(u.id)" />
            <view class="card-info" @click="goUser(u.id)">
              <text class="card-name">{{ u.nickname }}</text>
              <text class="card-desc">{{ u.school || '未知学校' }}</text>
            </view>
            <text :class="['follow-btn', u.is_friend && 'followed']" @click="toggleFollow(u)">
              {{ u.is_friend ? '好友' : '加好友' }}
            </text>
          </view>
          <view v-if="users.length === 0" class="empty">没有找到相关用户</view>
        </view>
      </template>

      <template v-else>
        <view v-if="trending.length" class="trending-section">
          <view class="section-header">
            <text class="section-title">今日热议</text>
          </view>
          <view v-for="(p, i) in trending" :key="p.id"
            class="trending-item" hover-class="trending-item-hover"
            @click="goPost(p.id)">
            <view :class="['trending-rank', i < 3 && 'top']">{{ i + 1 }}</view>
            <view class="trending-body">
              <text class="trending-content">{{ p.content }}</text>
              <view class="trending-meta">
                <u-icon name="heart" size="24" color="#B8C2CE" />
                <text>{{ p.like_count || 0 }}</text>
                <u-icon name="chat" size="24" color="#B8C2CE" />
                <text>{{ p.comment_count || 0 }}</text>
              </view>
            </view>
          </view>
        </view>
        <view v-else class="placeholder">
          <u-icon name="search" size="120" color="#EAE6E0" />
          <text class="placeholder-title">输入关键词搜索</text>
          <text class="placeholder-desc">搜索帖子、圈子或校园用户</text>
        </view>
      </template>

      <u-safe-bottom />
    </scroll-view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: #F7F4F0; display: flex; flex-direction: column; }

.top-bar { background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%); padding: 20rpx 24rpx 24rpx; }
.search-box {
  display: flex; align-items: center; gap: 16rpx;
  background: rgba(255,255,255,0.12); border-radius: 40rpx;
  padding: 0 20rpx; height: 72rpx;
}
.search-input { flex: 1; font-size: 28rpx; color: #fff; }
.search-input::placeholder { color: rgba(255,255,255,0.4); }
.clear-btn { font-size: 24rpx; color: rgba(255,255,255,0.6); cursor: pointer; }
.search-btn {
  font-size: 26rpx; color: #C67A6A; font-weight: 600; cursor: pointer;
  padding: 8rpx 16rpx; border-left: 1rpx solid rgba(255,255,255,0.15);
}

.tabs { display: flex; gap: 8rpx; padding: 16rpx 24rpx; background: #fff; border-bottom: 1rpx solid #EAE6E0; }
.tab {
  flex: 1; text-align: center; padding: 14rpx 0;
  font-size: 26rpx; color: #8E9BAB; border-radius: 10rpx;
  background: #F7F4F0; cursor: pointer;
}
.tab-active { color: #fff; background: linear-gradient(135deg, #C67A6A, #B06454); font-weight: 600; }

.scroll-area { flex: 1; }

.loading { text-align: center; padding: 200rpx 0; }
.empty { text-align: center; padding: 80rpx 0; color: #B8C2CE; font-size: 26rpx; }

.result-list { padding: 20rpx 16rpx; display: flex; flex-direction: column; gap: 16rpx; }

.result-card {
  background: #fff; border-radius: 14rpx; padding: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); transition: all 0.2s ease;
}
.result-card.row { display: flex; align-items: center; gap: 20rpx; }
.result-card-hover { transform: scale(0.99); }

.card-header { display: flex; align-items: center; gap: 14rpx; margin-bottom: 12rpx; }
.card-author { flex: 1; }
.card-name { font-size: 26rpx; font-weight: 600; color: #1E2A3A; display: block; }
.card-school { font-size: 22rpx; color: #8E9BAB; margin-top: 4rpx; display: block; }
.card-time { font-size: 22rpx; color: #B8C2CE; white-space: nowrap; }
.card-content { font-size: 26rpx; line-height: 1.6; color: #5C6B7E; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
.card-images { display: flex; gap: 8rpx; margin-top: 12rpx; }
.card-img { flex: 1; height: 180rpx; border-radius: 10rpx; background: #E8E4DE; }
.card-actions { display: flex; align-items: center; gap: 8rpx; margin-top: 12rpx; padding-top: 12rpx; border-top: 1rpx solid #EAE6E0; }
.action-txt { font-size: 22rpx; color: #B8C2CE; margin-right: 24rpx; }

.card-info { flex: 1; min-width: 0; }
.card-desc { font-size: 24rpx; color: #8E9BAB; margin-top: 4rpx; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-meta { font-size: 22rpx; color: #B8C2CE; margin-top: 8rpx; display: flex; align-items: center; gap: 12rpx; }
.badge { font-size: 20rpx; color: #C67A6A; border: 1rpx solid #C67A6A; border-radius: 20rpx; padding: 2rpx 16rpx; }

.follow-btn {
  font-size: 24rpx; color: #C67A6A; border: 1rpx solid #C67A6A;
  border-radius: 20rpx; padding: 8rpx 24rpx; white-space: nowrap; cursor: pointer;
}
.follow-btn.followed { color: #B8C2CE; border-color: #EAE6E0; }

.trending-section { padding: 24rpx 16rpx; }
.section-header { margin-bottom: 20rpx; }
.section-title { font-size: 30rpx; font-weight: 700; color: #1E2A3A; }

.trending-item {
  display: flex; gap: 16rpx; padding: 20rpx 0; border-bottom: 1rpx solid #EAE6E0; cursor: pointer;
  &:last-child { border-bottom: none; }
  transition: transform 0.2s ease;
}
.trending-item-hover { transform: scale(0.99); }
.trending-rank {
  width: 48rpx; height: 48rpx; border-radius: 10rpx;
  background: #EAE6E0; color: #8E9BAB;
  text-align: center; line-height: 48rpx; font-size: 24rpx; font-weight: 700; flex-shrink: 0;
}
.trending-rank.top {
  background: linear-gradient(135deg, #C67A6A, #B06454); color: #fff;
}
.trending-body { flex: 1; min-width: 0; }
.trending-content { font-size: 26rpx; color: #5C6B7E; display: block; margin-bottom: 8rpx; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.trending-meta { display: flex; align-items: center; gap: 8rpx; font-size: 22rpx; color: #B8C2CE; }

.placeholder {
  display: flex; flex-direction: column; align-items: center;
  padding: 200rpx 60rpx; gap: 16rpx;
}
.placeholder-title { font-size: 30rpx; font-weight: 600; color: #1E2A3A; }
.placeholder-desc { font-size: 26rpx; color: #8E9BAB; }
</style>
