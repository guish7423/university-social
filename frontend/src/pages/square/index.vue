<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { listPosts, followingPosts, toggleLike, type PostData } from "@/api/post"
import { useUserStore } from "@/store/user"
import ReportPopup from "@/components/ReportPopup.vue"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(true)
const currentTab = ref(0)
const tabs = ["全部", "关注"]
const showReport = ref(false)
const reportContentId = ref(0)
const reportContentType = ref("post")

const unsub = uni.$on('post-created', () => {
  if (currentTab.value === 0) fetchPosts()
  else fetchFollowing()
})

onMounted(() => {
  fetchPosts()
  if (userStore.isLogin) fetchFollowing()
})

onUnmounted(() => { if (typeof unsub === 'function') unsub() })

async function fetchPosts() {
  loading.value = true
  try { posts.value = await listPosts() } catch { posts.value = [] }
  loading.value = false
}

async function fetchFollowing() {
  loading.value = true
  try { posts.value = await followingPosts() } catch { posts.value = [] }
  loading.value = false
}

function onTabChange(index: number) {
  currentTab.value = index
  if (!userStore.isLogin && index === 1) {
    uni.showToast({ title: "请先登录", icon: "none" })
    uni.navigateTo({ url: "/pages/login/index" })
    return
  }
  if (index === 0) fetchPosts()
  else fetchFollowing()
}

async function handleLike(post: PostData) {
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  const res = await toggleLike(post.id)
  post.is_liked = res.liked
  post.like_count += res.liked ? 1 : -1
}

function handleMore(post: PostData) {
  const btns = ["举报"]
  if (post.user_id === userStore.user?.id) btns.push("删除")
  uni.showActionSheet({
    itemList: btns,
    success(res) {
      if (btns[res.tapIndex] === "举报") {
        reportContentId.value = post.id
        reportContentType.value = "post"
        showReport.value = true
      }
    }
  })
}

function goCreate() { uni.navigateTo({ url: "/pages/post/create" }) }
function goDetail(id: number) { uni.navigateTo({ url: `/pages/post/detail?id=${id}` }) }

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  const days = Math.floor(diff / 86400000)
  if (days < 7) return `${days}天前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<template>
  <view class="page-glow-top" />
  <view class="page-glow-bottom" />
  <view class="container">
    <view class="tab-bar glass-card">
      <view
        v-for="(tab, i) in tabs"
        :key="i"
        :class="['tab-item', currentTab === i && 'tab-active']"
        @click="onTabChange(i)"
      >
        <text class="tab-text">{{ tab }}</text>
        <view v-if="currentTab === i" class="tab-indicator" />
      </view>
    </view>

    <view v-if="loading" class="loading-state">
      <view v-for="n in 3" :key="n" class="skeleton-card glass-card" :style="{ animationDelay: `${(n-1)*0.1}s` }">
        <view class="skeleton-header">
          <view class="skeleton-avatar" />
          <view class="skeleton-lines">
            <view class="skeleton-line w-60" />
            <view class="skeleton-line w-40" />
          </view>
        </view>
        <view class="skeleton-body">
          <view class="skeleton-line w-90" />
          <view class="skeleton-line w-70" />
        </view>
      </view>
    </view>

    <template v-else-if="posts && posts.length">
      <view class="post-list">
        <view
          v-for="(post, idx) in posts"
          :key="post.id"
          :class="['post-card glass-card', 'animate-fade-in-up', 'stagger-' + Math.min(idx + 1, 8)]"
          @click="goDetail(post.id)"
        >
          <view class="post-header">
            <view class="author-avatar-wrap">
              <image class="author-avatar" :src="post.author?.avatar || ''" mode="aspectFill" />
            </view>
            <view class="post-author">
              <text class="author-name">{{ post.author?.nickname || "匿名" }}</text>
              <text class="post-time">{{ formatTime(post.created_at) }}</text>
            </view>
            <view class="more-btn" @click.stop="handleMore(post)">
              <text class="more-dot" />
            </view>
          </view>

          <view class="post-body">
            <text class="post-content">{{ post.content }}</text>
            <view v-if="post?.images?.length" class="post-images">
              <image v-for="(img, i) in post.images.slice(0, 3)" :key="i" class="post-image" :src="img" mode="aspectFill" />
            </view>
          </view>

          <view class="post-actions">
            <view :class="['action-btn', post.is_liked && 'is-liked']" @click.stop="handleLike(post)">
              <view class="like-icon-wrap">
                <text :class="['like-icon', post.is_liked && 'is-liked-icon']">{{ post.is_liked ? '❤️' : '♡' }}</text>
              </view>
              <text :class="['action-count', post.is_liked && 'liked']">{{ post.like_count || 0 }}</text>
            </view>
            <view class="action-btn">
              <text class="action-emoji">💬</text>
              <text class="action-count">{{ post.comment_count || 0 }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="fab glass-card" @click="goCreate">
        <text class="fab-icon">✏️</text>
      </view>
    </template>

    <template v-else-if="currentTab === 0">
      <view class="empty-state animate-fade-in-up">
        <text class="empty-icon">🌸</text>
        <text class="empty-title">还没有动态</text>
        <text class="empty-desc">成为第一个分享的人吧~</text>
        <view class="empty-btn" @click="goCreate">发布第一条动态</view>
      </view>
    </template>

    <template v-else>
      <view class="empty-state animate-fade-in-up">
        <text class="empty-icon">🌿</text>
        <text class="empty-title">关注的人还没有动态</text>
        <text class="empty-desc">去发现更多有趣的人吧~</text>
      </view>
    </template>
  </view>

  <ReportPopup :visible="showReport" :content-type="reportContentType" :content-id="reportContentId" @close="showReport = false" />
</template>

<style lang="scss" scoped>
.page-glow-top {
  position: fixed; top: -300rpx; right: -200rpx;
  width: 700rpx; height: 700rpx;
  background: radial-gradient(circle, rgba(192,132,252,0.12) 0%, rgba(249,168,212,0.08) 50%, transparent 70%);
  pointer-events: none; z-index: 0;
}
.page-glow-bottom {
  position: fixed; bottom: -200rpx; left: -200rpx;
  width: 500rpx; height: 500rpx;
  background: radial-gradient(circle, rgba(103,232,249,0.06) 0%, transparent 70%);
  pointer-events: none; z-index: 0;
}
.container { position: relative; z-index: 1; min-height: 100vh; padding-bottom: 120rpx; }

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  box-shadow: var(--glass-shadow);
}

.tab-bar {
  display: flex; margin: 20rpx 24rpx;
  border-radius: 20rpx; padding: 8rpx;
  position: sticky; top: 16rpx; z-index: 50;
}
.tab-item {
  flex: 1; display: flex; flex-direction: column;
  align-items: center; padding: 18rpx 0;
  cursor: pointer; border-radius: 14rpx;
  transition: all 0.3s ease;
}
.tab-item:hover { background: rgba(192,132,252,0.06); }
.tab-text { font-size: 28rpx; color: var(--ink-muted); font-weight: 500; transition: color 0.3s ease; }
.tab-active .tab-text { color: var(--brand-primary); font-weight: 600; }
.tab-indicator {
  width: 48rpx; height: 4rpx; border-radius: 2rpx;
  background: var(--brand-gradient); margin-top: 10rpx;
  animation: fadeIn 0.3s ease;
}

.loading-state { padding: 24rpx; display: flex; flex-direction: column; gap: 20rpx; }
.skeleton-card { padding: 32rpx; border-radius: 22rpx; }
.skeleton-header { display: flex; align-items: center; gap: 20rpx; margin-bottom: 24rpx; }
.skeleton-avatar { width: 72rpx; height: 72rpx; border-radius: 50%; background: linear-gradient(135deg, #E9D5FF 0%, #F3E8FF 100%); }
.skeleton-lines { flex: 1; display: flex; flex-direction: column; gap: 12rpx; }
.skeleton-line { height: 20rpx; border-radius: 10rpx; background: linear-gradient(135deg, #E9D5FF 0%, #F3E8FF 100%); }
.skeleton-body { display: flex; flex-direction: column; gap: 12rpx; }
.w-60 { width: 60%; } .w-40 { width: 40%; } .w-90 { width: 90%; } .w-70 { width: 70%; }

.post-list { padding: 0 24rpx 24rpx; }
.post-card {
  padding: 32rpx; margin-bottom: 20rpx; border-radius: 22rpx;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  &:active { transform: scale(0.98); }
}
.post-header { display: flex; align-items: center; gap: 20rpx; margin-bottom: 20rpx; }
.author-avatar-wrap {
  width: 72rpx; height: 72rpx; border-radius: 50%;
  padding: 3rpx; background: var(--brand-gradient);
}
.author-avatar { width: 100%; height: 100%; border-radius: 50%; }
.post-author { flex: 1; }
.author-name { font-size: 28rpx; font-weight: 600; color: var(--ink); display: block; }
.post-time { font-size: 22rpx; color: var(--ink-tertiary); display: block; margin-top: 4rpx; }
.more-btn {
  width: 48rpx; height: 48rpx; display: flex; align-items: center; justify-content: center;
  border-radius: 50%; transition: background 0.2s;
  &:hover { background: var(--hairline); }
}
.more-dot {
  width: 6rpx; height: 6rpx; border-radius: 50%; background: var(--ink-tertiary);
  box-shadow: 0 -10rpx 0 var(--ink-tertiary), 0 10rpx 0 var(--ink-tertiary);
}

.post-body { margin-bottom: 20rpx; }
.post-content { font-size: 30rpx; line-height: 1.7; color: var(--ink); display: block; }
.post-images { display: flex; gap: 12rpx; margin-top: 16rpx; }
.post-image { width: 210rpx; height: 210rpx; border-radius: 14rpx; background: var(--color-surface-1); }

.post-actions {
  display: flex; align-items: center; gap: 48rpx;
  padding-top: 20rpx; border-top: 1rpx solid var(--hairline);
}
.action-btn { display: flex; align-items: center; gap: 8rpx; transition: all 0.2s ease; }
.action-btn:active { transform: scale(0.9); }
.action-btn.is-liked { animation: likePulse 0.4s ease; }
.like-icon-wrap { width: 40rpx; height: 40rpx; display: flex; align-items: center; justify-content: center; }
.like-icon { font-size: 32rpx; color: var(--ink-tertiary); transition: all 0.3s ease; }
.like-icon.is-liked-icon { color: #FB7185; }
.action-emoji { font-size: 28rpx; }
.action-count { font-size: 24rpx; color: var(--ink-subtle); }
.action-count.liked { color: #FB7185; }

@keyframes likePulse {
  0% { transform: scale(1); }
  25% { transform: scale(1.3); }
  50% { transform: scale(0.95); }
  70% { transform: scale(1.15); }
  100% { transform: scale(1); }
}

.fab {
  position: fixed; right: 40rpx; bottom: 120rpx;
  width: 100rpx; height: 100rpx; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  background: var(--brand-gradient); z-index: 100;
  transition: all 0.3s ease;
  &:active { transform: scale(0.9); }
}
.fab-icon { font-size: 40rpx; }

.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 200rpx 40rpx; gap: 16rpx;
}
.empty-icon { font-size: 80rpx; }
.empty-title { font-size: 32rpx; font-weight: 600; color: var(--ink); }
.empty-desc { font-size: 26rpx; color: var(--ink-subtle); }
.empty-btn {
  margin-top: 24rpx; padding: 20rpx 48rpx;
  background: var(--brand-gradient); border-radius: 40rpx;
  color: #fff; font-size: 28rpx; font-weight: 500;
  transition: all 0.3s ease;
  &:active { transform: scale(0.95); }
}

@media (min-width: 1024px) {
  .post-list {
    display: grid; grid-template-columns: 1fr 1fr; gap: 24rpx;
    padding: 32rpx;
  }
  .post-card { margin-bottom: 0; }
  .post-image { width: 280rpx; height: 280rpx; }
  .fab { display: none; }
  .container { max-width: 1400rpx; margin: 0 auto; }
}
</style>
