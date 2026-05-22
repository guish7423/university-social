<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from "vue"
import { listPosts, followingPosts, toggleLike, type PostData } from "@/api/post"
import { listFeaturedPosts } from "@/api/social"
import { useUserStore } from "@/store/user"
import ReportPopup from "@/components/ReportPopup.vue"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(true)
const refreshing = ref(false)
const currentTab = ref(0)
const tabs = ["全部", "关注", "精华"]
const showReport = ref(false)
const reportContentId = ref(0)
const reportContentType = ref("post")
const page = ref(1)
const hasMore = ref(true)

const unsub = uni.$on('post-created', () => {
  if (currentTab.value === 0) fetchPosts()
  else if (currentTab.value === 1) fetchFollowing()
  else fetchFeatured()
})

onMounted(() => fetchPosts())

onUnmounted(() => {
  if (typeof unsub === 'function') unsub()
})

async function fetchPosts() {
  loading.value = true
  try {
    const data = await listPosts()
    posts.value = data || []
    hasMore.value = (data?.length || 0) >= 20
  } catch {}
  loading.value = false
}

async function fetchFollowing() {
  loading.value = true
  try {
    const data = await followingPosts()
    posts.value = data || []
    hasMore.value = (data?.length || 0) >= 20
  } catch { posts.value = [] }
  loading.value = false
}

async function fetchFeatured() {
  loading.value = true
  try {
    const data = await listFeaturedPosts()
    posts.value = data || []
    hasMore.value = (data?.length || 0) >= 20
  } catch { posts.value = [] }
  loading.value = false
}

async function onRefresh() {
  refreshing.value = true
  page.value = 1
  if (currentTab.value === 0) {
    await fetchPosts()
  } else if (currentTab.value === 1) {
    await fetchFollowing()
  } else {
    await fetchFeatured()
  }
  refreshing.value = false
  refreshing.value = false
}

async function loadMore() {
  if (!hasMore.value || loading.value) return
  page.value++
  try {
    const more = currentTab.value === 0
      ? await listPosts()
      : currentTab.value === 1
        ? await followingPosts()
        : await listFeaturedPosts()

    if (more?.length) {
      posts.value.push(...more)
      hasMore.value = more.length >= 20
    } else {
      hasMore.value = false
    }
  } catch { hasMore.value = false }
}

function onTabChange(index: number) {
  currentTab.value = index
  if (!userStore.isLogin && (index === 1 || index === 2)) {
    uni.showToast({ title: "请先登录", icon: "none" })
    uni.navigateTo({ url: "/pages/login/index" })
    return
  }
  if (index === 0) fetchPosts()
  else if (index === 1) fetchFollowing()
  else fetchFeatured()
}


function previewImage(src: string) {
  uni.previewImage({ urls: [src], indicator: "default" })
}

async function handleLike(post: PostData) {
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  post._liking = true
  try {
    const res = await toggleLike(post.id)
    post.is_liked = res.liked
    post.like_count += res.liked ? 1 : -1
  } catch {}
  post._liking = false
}

function handleMore(post: PostData) {
  const btns: string[] = ["举报"]
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
  if (!t) return ""
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  const days = Math.floor(diff / 86400000)
  if (days < 7) return `${days}天前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}

const skeletonItems = ref([1, 2, 3])
</script>

<template>
  <view class="container">
    <view class="sticky-top">
      <view class="top-bar">
        <text class="top-bar-title">广场</text>
        <view class="top-bar-actions">
          <u-icon name="search" size="36" color="#fff" @click="uni.navigateTo({ url: '/pages/search/index' })" />
        </view>
      </view>
      <view class="tab-bar">
        <view v-for="(tab, i) in tabs" :key="i"
          :class="['tab-item', currentTab === i && 'tab-active']"
          @click="onTabChange(i)">
          <text class="tab-text">{{ tab }}</text>
          <view v-if="currentTab === i" class="tab-indicator" />
        </view>
      </view>
    </view>

    <scroll-view
      scroll-y
      class="scroll-area"
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="loadMore"
    >
      <view v-if="loading && !refreshing" class="skeleton-area">
        <view v-for="n in skeletonItems" :key="n" class="skeleton-card">
          <view class="skeleton-row" style="width:40%;height:24rpx" />
          <view class="skeleton-row" style="width:25%;height:20rpx" />
          <view class="skeleton-row" style="width:80%;height:28rpx;margin-top:16rpx" />
          <view class="skeleton-row" style="width:60%;height:28rpx" />
          <view class="skeleton-images">
            <view class="skeleton-img" />
            <view class="skeleton-img" />
            <view class="skeleton-img" />
          </view>
          <view class="skeleton-row" style="width:50%;height:24rpx;margin-top:12rpx" />
        </view>
      </view>

      <template v-else-if="posts.length">
        <view class="post-list">
          <view v-for="post in posts" :key="post.id"
            class="post-card" hover-class="post-card-hover"
            @click="goDetail(post.id)">
            <view class="post-card-header">
              <u-avatar :src="post.author?.avatar" :text="post.author?.nickname?.[0] || '?'" size="56" shape="circle" />
              <view class="post-card-author">
                <text class="post-card-name">{{ post.author?.nickname || "匿名" }}</text>
                <text class="post-card-time">{{ formatTime(post.created_at) }}</text>
              </view>
              <u-icon name="more-dot-fill" color="#B8C2CE" size="28" @click.stop="handleMore(post)" />
            </view>

            <text class="post-card-content">{{ post.content }}</text>

            <view v-if="post?.images?.length" class="post-card-images">
              <image v-for="(img, j) in post.images.slice(0, 3)" :key="j"
                :class="['post-card-image', post.images.length === 1 ? 'single' : '']"
                :src="img" mode="aspectFill" @click.stop="previewImage(img)" />
              <view v-if="post.images.length > 3" class="image-more">+{{ post.images.length - 3 }}</view>
            </view>

            <view class="post-card-footer">
              <view :class="['post-card-action', post.is_liked && 'liked']"
                @click.stop="handleLike(post)">
                <view :class="['like-icon-wrap', post._liking && 'liking']">
                  <u-icon :name="post.is_liked ? 'heart-fill' : 'heart'"
                    :color="post.is_liked ? '#C67A6A' : '#B8C2CE'" size="32" />
                </view>
                <text :class="['action-count', post.is_liked && 'liked']">{{ post.like_count || 0 }}</text>
              </view>
              <view class="post-card-action" @click.stop="goDetail(post.id)">
                <u-icon name="chat" color="#B8C2CE" size="32" />
                <text class="action-count">{{ post.comment_count || 0 }}</text>
              </view>
            </view>
          </view>
        </view>

        <view v-if="hasMore" class="load-more-hint">
          <u-loading mode="flower" size="32" />
          <text class="load-more-text">加载更多...</text>
        </view>
        <view v-else-if="posts.length" class="load-more-hint">
          <text class="load-more-text no-more">— 没有更多了 —</text>
        </view>
      </template>

      <template v-else-if="currentTab === 0">
        <view class="empty-state">
          <u-icon name="empty-address" size="160" color="#EAE6E0" />
          <text class="empty-title">还没有动态</text>
          <text class="empty-desc">快来发布第一条吧，分享校园生活</text>
          <u-button type="primary" shape="circle" class="empty-btn" @click="goCreate">发布动态</u-button>
        </view>
      </template>
      <template v-else-if="currentTab === 1">
        <view class="empty-state">
          <u-icon name="attention" size="160" color="#EAE6E0" />
          <text class="empty-title">关注的人还没有动态</text>
          <text class="empty-desc">去广场发现更多有趣的人</text>
          <u-button type="primary" shape="circle" class="empty-btn" plain @click="onTabChange(0)">去广场看看</u-button>
        </view>
      </template>
      <template v-else>
        <view class="empty-state">
          <u-icon name="star" size="160" color="#EAE6E0" />
          <text class="empty-title">还没有精华帖</text>
          <text class="empty-desc">管理员会将优质帖子设为精华</text>
          <u-button type="primary" shape="circle" class="empty-btn" plain @click="onTabChange(0)">去全部看看</u-button>
        </view>
      </template>

      <u-safe-bottom />
    </scroll-view>

    <view class="fab" @click="goCreate"><u-icon name="plus" size="40" color="#fff" /></view>
  </view>

  <ReportPopup :visible="showReport" :content-type="reportContentType" :content-id="reportContentId" @close="showReport = false" />
</template>

<style lang="scss" scoped>
.container { min-height: 100vh; background: #F7F4F0; display: flex; flex-direction: column; }

.sticky-top { position: sticky; top: 0; z-index: 100; }
.top-bar {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 20rpx 32rpx 0;
  display: flex; align-items: center; justify-content: space-between; height: 88rpx;
}
.top-bar-title { font-size: 32rpx; font-weight: 700; color: #fff; }
.top-bar-actions { display: flex; gap: 24rpx; }

.tab-bar {
  display: flex; background: #fff;
  padding: 0 40rpx;
  border-bottom: 1rpx solid #E0DBD4;
}
.tab-item { flex: 1; display: flex; flex-direction: column; align-items: center; padding: 20rpx 0; cursor: pointer; }
.tab-text { font-size: 28rpx; color: #8E9BAB; font-weight: 500; }
.tab-active .tab-text { color: #C67A6A; font-weight: 600; }
.tab-indicator { width: 40rpx; height: 4rpx; border-radius: 2rpx; background: #C67A6A; margin-top: 12rpx; }

.scroll-area { flex: 1; }

.skeleton-area { padding: 20rpx 16rpx; display: flex; flex-direction: column; gap: 16rpx; }
.skeleton-card {
  background: #fff; border-radius: 18rpx; padding: 28rpx;
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
  animation: pulse 1.5s ease-in-out infinite;
}
.skeleton-row {
  background: #EAE6E0; border-radius: 8rpx; margin-bottom: 8rpx;
}
.skeleton-images { display: flex; gap: 8rpx; margin-top: 12rpx; }
.skeleton-img { flex: 1; height: 200rpx; border-radius: 12rpx; background: #EAE6E0; }

.post-list { padding: 20rpx 16rpx; display: flex; flex-direction: column; gap: 16rpx; padding-bottom: 40rpx; }

.post-card {
  background: #fff; border-radius: 18rpx; padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  transition: all 0.2s ease; cursor: pointer;
}
.post-card-hover { transform: scale(0.99); }

.post-card-header { display: flex; align-items: center; gap: 14rpx; margin-bottom: 16rpx; }
.post-card-author { flex: 1; }
.post-card-name { font-size: 26rpx; font-weight: 600; color: #1E2A3A; display: block; }
.post-card-time { font-size: 20rpx; color: #B8C2CE; display: block; margin-top: 4rpx; }

.post-card-content { font-size: 28rpx; line-height: 1.7; color: #5C6B7E; display: block; word-break: break-all; }

.post-card-images { display: flex; gap: 8rpx; margin-top: 16rpx; position: relative; }
.post-card-image { flex: 1; height: 220rpx; border-radius: 12rpx; background: #E8E4DE; }
.post-card-image.single { flex: none; width: 100%; height: 400rpx; }
.image-more {
  position: absolute; right: 8rpx; bottom: 8rpx;
  background: rgba(0,0,0,0.5); color: #fff; font-size: 22rpx;
  padding: 4rpx 12rpx; border-radius: 8rpx;
}

.post-card-footer {
  display: flex; align-items: center; gap: 48rpx;
  padding-top: 16rpx; margin-top: 16rpx;
  border-top: 1rpx solid #EAE6E0;
}
.post-card-action { display: flex; align-items: center; gap: 8rpx; transition: transform 0.2s ease; &:active { transform: scale(0.9); } }
.like-icon-wrap.liking { animation: likeBounce 0.3s ease; }
@keyframes likeBounce { 0% { transform: scale(1); } 50% { transform: scale(1.3); } 100% { transform: scale(1); } }
.action-count { font-size: 24rpx; color: #B8C2CE; }
.action-count.liked { color: #C67A6A; }

.load-more-hint { display: flex; align-items: center; justify-content: center; gap: 12rpx; padding: 32rpx 0; }
.load-more-text { font-size: 24rpx; color: #B8C2CE; }
.no-more { color: #D0CBC4; }

.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 200rpx 60rpx; gap: 16rpx;
}
.empty-title { font-size: 30rpx; font-weight: 600; color: #1E2A3A; }
.empty-desc { font-size: 26rpx; color: #8E9BAB; }
.empty-btn { margin-top: 16rpx; }

.fab {
  position: fixed; right: 40rpx; bottom: 120rpx;
  width: 96rpx; height: 96rpx; border-radius: 50%;
  background: linear-gradient(135deg, #C67A6A, #B06454);
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 6rpx 24rpx rgba(198,122,106,0.35);
  z-index: 100;
  transition: transform 0.2s ease;
  &:active { transform: scale(0.92); }
}
</style>
