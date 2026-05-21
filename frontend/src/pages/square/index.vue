<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { listPosts, followingPosts, toggleLike, type PostData } from "@/api/post"
import { useUserStore } from "@/store/user"
import ReportPopup from "@/components/ReportPopup.vue"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(false)
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

onUnmounted(() => {
  if (typeof unsub === 'function') unsub()
})

async function fetchPosts() {
  loading.value = true
  try { posts.value = await listPosts() } catch {}
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
  <view class="container">
    <view class="tab-bar">
      <view v-for="(tab, i) in tabs" :key="i"
        :class="['tab-item', currentTab === i && 'tab-active']"
        @click="onTabChange(i)">
        <text class="tab-text">{{ tab }}</text>
        <view v-if="currentTab === i" class="tab-indicator" />
      </view>
    </view>

    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="48" />
    </view>

    <template v-else-if="posts && posts.length">
      <view class="post-list">
        <view v-for="(post, idx) in posts" :key="post.id"
          :class="['post-card', 'animate-fade-in-up', 'stagger-' + Math.min(idx + 1, 8)]"
          @click="goDetail(post.id)">
          <view class="post-header">
            <u-avatar :src="post.author?.avatar" :text="post.author?.nickname?.[0]" size="64" shape="circle" />
            <view class="post-author">
              <text class="author-name">{{ post.author?.nickname || "匿名" }}</text>
              <text class="post-time">{{ formatTime(post.created_at) }}</text>
            </view>
            <u-icon name="more" color="#B8C2CE" size="28" @click.stop="handleMore(post)" />
          </view>

          <view class="post-body">
            <text class="post-content">{{ post.content }}</text>
            <view v-if="post?.images?.length" class="post-images">
              <image v-for="(img, i) in post.images.slice(0, 3)" :key="i"
                class="post-image" :src="img" mode="aspectFill" />
            </view>
          </view>

          <view class="post-actions">
            <view :class="['action-btn', post.is_liked && 'is-liked']" @click.stop="handleLike(post)">
              <u-icon :name="post.is_liked ? 'heart-fill' : 'heart'" :color="post.is_liked ? '#C67A6A' : '#B8C2CE'" :size="32" />
              <text :class="['action-count', post.is_liked && 'liked']">{{ post.like_count || 0 }}</text>
            </view>
            <view class="action-btn">
              <u-icon name="chat" color="#B8C2CE" size="32" />
              <text class="action-count">{{ post.comment_count || 0 }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="fab" @click="goCreate"><text class="fab-icon">+</text></view>
    </template>

    <template v-else-if="currentTab === 0">
      <u-empty mode="data" text="还没有动态，快来发第一条吧" />
      <view class="empty-create">
        <u-button type="primary" shape="circle" @click="goCreate">发布第一条动态</u-button>
      </view>
    </template>

    <template v-else>
      <u-empty mode="data" text="关注的人还没有动态" />
    </template>
  </view>
</template>

<ReportPopup :visible="showReport" :content-type="reportContentType" :content-id="reportContentId" @close="showReport = false" />

<style lang="scss" scoped>
.container { min-height: 100vh; background: var(--color-canvas, #F7F4F0); }

.tab-bar {
  display: flex; background: var(--color-surface, #fff);
  padding: 0 40rpx; position: sticky; top: 0; z-index: 50;
  border-bottom: 1rpx solid var(--hairline, #E0DBD4);
}
.tab-item { flex: 1; display: flex; flex-direction: column; align-items: center; padding: 24rpx 0; cursor: pointer; }
.tab-text { font-size: 28rpx; color: var(--ink-subtle, #8E9BAB); font-weight: 500; transition: color 0.3s ease; }
.tab-active .tab-text { color: var(--brand-primary, #C67A6A); font-weight: 600; }
.tab-indicator { width: 40rpx; height: 4rpx; border-radius: 2rpx; background: #C67A6A; margin-top: 12rpx; animation: fadeIn 0.3s ease; }

.loading-state { display: flex; justify-content: center; padding: 200rpx 0; }

.post-list { padding: 20rpx; }
.post-card {
  background: var(--color-surface, #fff);
  border-radius: 18rpx;
  padding: 28rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  transition: all 0.2s ease;
  &:active { transform: scale(0.99); }
}

.post-header { display: flex; align-items: center; gap: 16rpx; margin-bottom: 16rpx; }
.post-author { flex: 1; }
.author-name { font-size: 26rpx; font-weight: 600; color: var(--ink, #1E2A3A); display: block; }
.post-time { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); display: block; margin-top: 4rpx; }

.post-body { margin-bottom: 16rpx; }
.post-content { font-size: 28rpx; line-height: 1.7; color: var(--ink-muted, #5C6B7E); display: block; }
.post-images { display: flex; gap: 12rpx; margin-top: 16rpx; }
.post-image { width: 210rpx; height: 210rpx; border-radius: 12rpx; background: var(--color-surface-2, #E8E4DE); }

.post-actions { display: flex; align-items: center; gap: 48rpx; padding-top: 16rpx; border-top: 1rpx solid var(--hairline, #E0DBD4); }
.action-btn { display: flex; align-items: center; gap: 8rpx; transition: transform 0.2s ease; &:active { transform: scale(0.9); } }
.action-count { font-size: 24rpx; color: var(--ink-tertiary, #B8C2CE); }
.action-count.liked { color: #C67A6A; }

.fab {
  position: fixed; right: 40rpx; bottom: 120rpx;
  width: 88rpx; height: 88rpx;
  border-radius: 50%;
  background: #C67A6A;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 4rpx 20rpx rgba(198,122,106,0.3);
  z-index: 100;
  transition: transform 0.2s ease;
  &:active { transform: scale(0.92); }
}
.fab-icon { font-size: 40rpx; color: #fff; font-weight: 300; }

.empty-create { padding: 40rpx 60rpx; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }

@media (min-width: 1024px) {
  .post-list { display: grid; grid-template-columns: 1fr 1fr; gap: 24rpx; padding: 32rpx; }
  .post-card { margin-bottom: 0; }
  .fab { display: none; }
  .container { max-width: 1400rpx; margin: 0 auto; }
}
</style>
