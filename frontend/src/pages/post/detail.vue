<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { getPost, listComments, createComment, toggleLike, updatePost, type PostData, type CommentData } from "@/api/post"
import { followUser, unfollowUser, checkFollow } from "@/api/follow"
import ReportPopup from "@/components/ReportPopup.vue"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const isOwner = computed(() => post.value?.user_id === userStore.id)

const post = ref<PostData | null>(null)
const comments = ref<CommentData[]>([])
const commentText = ref("")
const submitting = ref(false)
const postId = ref(0)
const isFollowingAuthor = ref(false)
const followLoading = ref(false)
const showReport = ref(false)
const editContent = ref("")
const editing = ref(false)
const imageIndex = ref(0)

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  postId.value = Number(page?.options?.id || 0)
  if (postId.value) {
    loadPost(postId.value)
    loadComments(postId.value)
  }
})

async function loadPost(id: number) {
  try {
    post.value = await getPost(id)
    if (post.value?.author?.id) checkFollowState(post.value.author.id)
  } catch (e) { console.error(e) }
}
    editContent.value = post.value?.content || ""

async function checkFollowState(userId: number) {
  try { isFollowingAuthor.value = (await checkFollow(userId)).is_following } catch (e) { console.error(e) }
}

async function loadComments(id: number) {
  try { comments.value = await listComments(id) } catch (e) { console.error(e) }
}

async function handleComment() {
  if (!commentText.value.trim()) return
  submitting.value = true
  try {
    await createComment(postId.value, commentText.value)
    commentText.value = ""
    if (post.value) post.value.comment_count++
    loadComments(postId.value)
  } catch (e) { console.error(e) }
  submitting.value = false
}

async function handleLikePost() {
  if (!post.value) return
  const res = await toggleLike(postId.value)
  post.value.is_liked = res.liked
  post.value.like_count += res.liked ? 1 : -1
}

async function toggleFollowAuthor() {
  if (!post.value?.author?.id || followLoading.value) return
  followLoading.value = true
  try {
    if (isFollowingAuthor.value) {
      await unfollowUser(post.value.author.id)
      isFollowingAuthor.value = false
    } else {
      await followUser(post.value.author.id)
      isFollowingAuthor.value = true
    }
  } catch (e) { console.error(e) }
  followLoading.value = false
}

async function startEdit() {
  editing.value = true
  editContent.value = post.value?.content || ""
}

async function saveEdit() {
  if (!post.value || !editContent.value.trim()) return
  try {
    post.value = await updatePost(postId.value, { content: editContent.value })
    editing.value = false
    uni.showToast({ title: "已编辑", icon: "none" })
  } catch {
    uni.showToast({ title: "编辑失败", icon: "error" })
  }
}

function cancelEdit() {
  editing.value = false
  editContent.value = post.value?.content || ""
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
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<template>
  <view class="container">
    <view v-if="!post" class="loading-state">
      <u-loading mode="flower" size="60" />
    </view>
    <template v-else>
      <view class="post-card">
        <view class="post-header">
          <u-avatar :src="post.author?.avatar" :text="post.author?.nickname?.[0] || '?'" size="64" shape="circle" />
          <view class="post-author">
            <text class="author-name">{{ post.author?.nickname || "匿名" }}</text>
            <text class="post-time">{{ formatTime(post.created_at) }}</text>
          </view>
          <u-button
            v-if="post.author?.id"
            size="mini" shape="circle"
            :type="isFollowingAuthor ? 'default' : 'primary'"
            :text="isFollowingAuthor ? '已关注' : '+ 关注'"
            :loading="followLoading"
            :customStyle="{ marginLeft: 'auto', minWidth: '100rpx' }"
            @click="toggleFollowAuthor"
          />
        </view>
        <textarea v-if="editing" v-model="editContent" class="edit-textarea" :maxlength="2000" />
        <text v-else class="post-content">{{ post.content }}</text>
        <view v-if="post.images?.length" class="post-images">
          <swiper class="image-swiper" :indicator-dots="post.images.length > 1" indicator-color="#EAE6E0" indicator-active-color="#C67A6A" @change="(e: any) => imageIndex = e.detail.current">
            <swiper-item v-for="(img, i) in post.images" :key="i">
              <image :src="img" mode="aspectFill" class="swiper-image" />
            </swiper-item>
          </swiper>
        </view>
        <view class="post-stats">
          <text class="stat-item"><u-icon name="eye" size="28" color="#B8C2CE" /> {{ post.images?.length || 0 }} 张图片</text>
        </view>
      </view>

      <view class="action-bar">
        <view :class="['action-btn', post.is_liked && 'liked']" @click="handleLikePost">
          <u-icon :name="post.is_liked ? 'heart-fill' : 'heart'" :color="post.is_liked ? '#C67A6A' : '#B8C2CE'" size="36" />
          <text :class="['action-text', post.is_liked && 'liked']">{{ post.like_count || 0 }} 赞</text>
        </view>
        <view class="action-btn">
          <u-icon name="chat" color="#B8C2CE" size="36" />
          <text class="action-text">{{ post.comment_count || 0 }} 评论</text>
        </view>
        <view v-if="isOwner" class="action-btn" @click="editing ? saveEdit() : startEdit()">
          <u-icon :name="editing ? 'checkmark' : 'edit-pen'" color="#C67A6A" size="36" />
          <text class="action-text">{{ editing ? '保存' : '编辑' }}</text>
        </view>
        <view class="action-btn" @click="showReport = true">
          <u-icon name="info-circle" color="#B8C2CE" size="36" />
          <text class="action-text">举报</text>
        </view>
      </view>

      <view class="comment-section">
        <view class="comment-header">
          <text class="section-title">评论 ({{ comments?.length || 0 }})</text>
        </view>
        <u-empty v-if="!comments?.length" mode="message" text="暂无评论" :iconSize="80" />
        <view v-for="c in (comments || [])" :key="c.id" class="comment-item">
          <u-avatar :src="c.author?.avatar" :text="c.author?.nickname?.[0] || '?'" size="48" shape="circle" fontSize="20" />
          <view class="comment-body">
            <view class="comment-top">
              <text class="comment-author">{{ c.author?.nickname || "匿名" }}</text>
              <text class="comment-time">{{ formatTime(c.created_at) }}</text>
            </view>
            <text class="comment-content">{{ c.content }}</text>
          </view>
        </view>
      </view>
    </template>

    <view class="comment-input-bar">
      <u-input v-model="commentText" placeholder="写下你的评论..." :border="false" shape="circle"
        :customStyle="{ flex: 1, height: '72rpx', background: '#F0EDE8', padding: '0 24rpx', fontSize: '28rpx', borderRadius: '36rpx' }" />
      <u-button type="primary" size="small" shape="circle" :loading="submitting"
        :disabled="!commentText.trim()"
        :customStyle="{ height: '72rpx', padding: '0 30rpx', fontSize: '26rpx' }"
        @click="handleComment">发送</u-button>
    </view>

    <ReportPopup :visible="showReport" content-type="post" :content-id="postId" @close="showReport = false" />
  </view>
</template>

<style lang="scss">
page { background: #F7F4F0; }
</style>

<style lang="scss" scoped>
.container { min-height: 100vh; background: #F7F4F0; padding-bottom: 140rpx; }
.loading-state { display: flex; justify-content: center; padding: 200rpx 0; }

.post-card { background: #fff; padding: 30rpx; }
.post-header { display: flex; align-items: center; gap: 16rpx; margin-bottom: 20rpx; }
.post-author { flex: 1; }
.author-name { font-size: 28rpx; font-weight: 600; color: #1E2A3A; display: block; }
.post-time { font-size: 22rpx; color: #B8C2CE; margin-top: 4rpx; display: block; }

.post-content { font-size: 30rpx; line-height: 1.7; color: #1E2A3A; white-space: pre-wrap; }

.post-images { margin-top: 20rpx; border-radius: 14rpx; overflow: hidden; }
.image-swiper { width: 100%; height: 500rpx; }
.swiper-image { width: 100%; height: 500rpx; background: #E8E4DE; }

.post-stats { margin-top: 16rpx; }
.stat-item { font-size: 22rpx; color: #B8C2CE; display: flex; align-items: center; gap: 6rpx; }

.action-bar {
  display: flex; align-items: center; justify-content: space-around;
  background: #fff; margin-top: 16rpx; padding: 20rpx 0;
  border-top: 1rpx solid #EAE6E0; border-bottom: 1rpx solid #EAE6E0;
}
.action-btn {
  display: flex; align-items: center; gap: 8rpx;
  transition: transform 0.2s ease; &:active { transform: scale(0.9); }
}
.action-text { font-size: 24rpx; color: #B8C2CE; }
.action-text.liked { color: #C67A6A; }

.comment-section { background: #fff; margin-top: 16rpx; padding: 30rpx; }
.comment-header { margin-bottom: 16rpx; }
.section-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }

.comment-item { display: flex; gap: 12rpx; padding: 20rpx 0; border-bottom: 1rpx solid #EAE6E0; &:last-child { border-bottom: none; } }
.comment-body { flex: 1; }
.comment-top { display: flex; align-items: center; gap: 12rpx; margin-bottom: 6rpx; }
.comment-author { font-size: 24rpx; font-weight: 600; color: #5C6B7E; }
.comment-time { font-size: 20rpx; color: #B8C2CE; }
.comment-content { font-size: 28rpx; color: #1E2A3A; line-height: 1.5; }

.comment-input-bar {
  position: fixed; bottom: 0; left: 0; right: 0;
  display: flex; gap: 16rpx; align-items: center;
  background: #fff; padding: 16rpx 30rpx;
  border-top: 1rpx solid #E0DBD4; z-index: 100;
}

@media (min-width: 1024px) {
  .container { max-width: 960px; margin: 0 auto; padding: 32rpx 0 140rpx; }
  .post-card { border-radius: 22rpx; }
  .action-bar { border-radius: 0 0 22rpx 22rpx; }
  .comment-section { border-radius: 22rpx; }
  .comment-input-bar { max-width: 960px; left: 50%; transform: translateX(-50%); border-radius: 999rpx; bottom: 32rpx; }
}
</style>
