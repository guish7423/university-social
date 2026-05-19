<template>
  <view class="min-h-screen bg-gray-50">
    <view v-if="circle" class="bg-gradient-to-r from-indigo-500 to-purple-600 px-4 pt-6 pb-4 text-white">
      <view class="flex items-center gap-3">
        <view class="w-16 h-16 rounded-full bg-white/20 flex items-center justify-center text-3xl">
          {{ circle.icon || circle.name[0] }}
        </view>
        <view class="flex-1">
          <text class="text-xl font-bold block">{{ circle.name }}</text>
          <text class="text-sm text-white/80 block mt-0.5">{{ circle.description }}</text>
          <view class="flex gap-3 mt-1">
            <text class="text-xs text-white/70">{{ circle.member_count }} 成员</text>
            <text class="text-xs text-white/70">{{ circle.post_count }} 帖子</text>
          </view>
        </view>
      </view>
      <view class="flex gap-3 mt-4">
        <button v-if="!circle.is_member"
                class="flex-1 bg-white text-indigo-600 rounded-full py-1.5 text-sm font-medium"
                @click="handleJoin">加入圈子</button>
        <button v-else
                class="flex-1 bg-white/20 text-white rounded-full py-1.5 text-sm font-medium"
                @click="handleLeave">退出圈子</button>
        <button v-if="circle.is_member"
                class="bg-white text-indigo-600 rounded-full py-1.5 px-4 text-sm"
                @click="showCreatePost = true">发帖</button>
      </view>
    </view>

    <view v-if="circle" class="flex border-b border-gray-200 bg-white">
      <view v-for="tab in tabs" :key="tab.key"
            class="flex-1 text-center py-3 text-sm"
            :class="activeTab === tab.key ? 'text-indigo-600 border-b-2 border-indigo-600 font-medium' : 'text-gray-500'"
            @click="activeTab = tab.key">{{ tab.label }}</view>
    </view>

    <view v-if="activeTab === 'posts'" class="p-4 space-y-3">
      <view v-if="posts.length === 0" class="text-center py-10 text-gray-400">还没有帖子</view>
      <view v-for="p in posts" :key="p.id" class="bg-white rounded-lg p-4 shadow-sm">
        <view class="flex items-center gap-2 mb-2">
          <image :src="p.author.avatar || '/static/default-avatar.png'" class="w-8 h-8 rounded-full bg-gray-200" />
          <text class="text-sm font-medium">{{ p.author.nickname }}</text>
          <text class="text-xs text-gray-400 ml-auto">{{ formatTime(p.created_at) }}</text>
        </view>
        <text class="text-sm leading-relaxed">{{ p.content }}</text>
        <view v-if="p.image_urls?.length" class="flex flex-wrap gap-2 mt-2">
          <image v-for="(img, i) in p.image_urls" :key="i"
                 :src="img" class="w-20 h-20 rounded-lg bg-gray-100" mode="aspectFill" />
        </view>
        <view class="flex gap-4 mt-3 text-xs text-gray-400 border-t border-gray-100 pt-2">
          <view class="flex items-center gap-1" @click="toggleLike(p)">
            <text :class="p.is_liked ? 'text-red-500' : ''">{{ p.is_liked ? '❤️' : '🤍' }}</text>
            <text>{{ p.like_count || '' }}</text>
          </view>
          <view class="flex items-center gap-1" @click="showComments(p)">
            <text>💬</text>
            <text>{{ p.comment_count || '' }}</text>
          </view>
        </view>
      </view>
    </view>

    <view v-if="activeTab === 'members'" class="p-4 space-y-2">
      <view v-for="m in members" :key="m.id"
            class="bg-white rounded-lg p-3 flex items-center gap-3 shadow-sm">
        <image :src="m.user.avatar || '/static/default-avatar.png'" class="w-10 h-10 rounded-full bg-gray-200" />
        <view class="flex-1">
          <text class="text-sm font-medium block">{{ m.user.nickname }}</text>
        </view>
        <text v-if="m.role === 'creator'" class="text-xs text-yellow-500 bg-yellow-50 px-2 py-0.5 rounded">创建者</text>
        <text v-else class="text-xs text-gray-400">成员</text>
      </view>
    </view>

    <view v-if="showCreatePost" class="fixed inset-0 bg-black/30 z-50 flex items-end justify-center" @click="showCreatePost = false">
      <view class="bg-white rounded-t-2xl w-full p-4 pb-8" @click.stop>
        <text class="font-bold text-base mb-3 block">发帖</text>
        <textarea v-model="newPostContent" class="w-full border border-gray-200 rounded-lg p-3 text-sm h-28"
                  placeholder="说点什么..." maxlength="500" />
        <button class="w-full mt-3 py-2.5 rounded-lg text-sm text-white bg-gradient-to-r from-indigo-500 to-purple-600"
                :disabled="!newPostContent.trim()"
                :class="newPostContent.trim() ? '' : 'opacity-50'"
                @click="handleCreatePost">发布</button>
      </view>
    </view>

    <view v-if="showCommentsPanel" class="fixed inset-0 bg-black/30 z-50 flex items-end justify-center" @click="showCommentsPanel = false">
      <view class="bg-white rounded-t-2xl w-full max-h-96 overflow-y-auto" @click.stop>
        <view v-if="selectedPost" class="p-4 border-b border-gray-100">
          <text class="font-medium">{{ selectedPost.author.nickname }}: </text>
          <text class="text-sm text-gray-600">{{ selectedPost.content }}</text>
        </view>
        <view v-if="comments.length === 0" class="px-4 py-6 text-center text-xs text-gray-400">暂无评论</view>
        <view v-for="c in comments" :key="c.id" class="px-4 pb-3 flex gap-2">
          <text class="text-xs text-indigo-500 shrink-0">{{ c.author.nickname }}:</text>
          <text class="text-xs">{{ c.content }}</text>
        </view>
        <view class="flex gap-2 p-3 border-t border-gray-100">
          <input v-model="newComment" class="flex-1 border border-gray-200 rounded-full px-3 py-1.5 text-sm"
                 placeholder="写评论..." />
          <button class="text-sm text-indigo-500 font-medium" @click="handleCreateComment">发送</button>
        </view>
      </view>
    </view>

    <view v-if="!circle && !loading" class="flex flex-col items-center py-20">
      <text class="text-gray-400">圈子不存在</text>
    </view>
    <view v-if="loading" class="flex justify-center py-20">
      <text class="text-gray-400">加载中...</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getCircle, joinCircle, leaveCircle, listCircleMembers,
         listCirclePosts, createCirclePost, toggleCirclePostLike,
         listCirclePostComments, createCirclePostComment } from "@/api/circle"
import type { CircleData, CircleMemberData, CirclePostData } from "@/api/circle"
import type { UserInfo } from "@/api/user"

const circle = ref<CircleData>()
const members = ref<CircleMemberData[]>([])
const posts = ref<CirclePostData[]>([])
const loading = ref(true)
const activeTab = ref("posts")
const showCreatePost = ref(false)
const newPostContent = ref("")
const showCommentsPanel = ref(false)
const selectedPost = ref<CirclePostData>()
const comments = ref<{ id: number; user_id: number; content: string; created_at: string; author: UserInfo }[]>([])
const newComment = ref("")

const tabs = [{ key: "posts", label: "帖子" }, { key: "members", label: "成员" }]
const circleId = ref(0)

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  circleId.value = page?.options?.id ? parseInt(page.options.id) : 0
  if (circleId.value) load()
})

async function load() {
  try {
    const [c, m, p] = await Promise.all([
      getCircle(circleId.value),
      listCircleMembers(circleId.value),
      listCirclePosts(circleId.value),
    ])
    circle.value = c
    members.value = m
    posts.value = p
  } catch {
    uni.showToast({ title: "加载失败", icon: "none" })
  } finally {
    loading.value = false
  }
}

async function handleJoin() {
  try {
    await joinCircle(circleId.value)
    uni.showToast({ title: "已加入", icon: "success" })
    circle.value!.is_member = true
    circle.value!.member_count++
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
}

async function handleLeave() {
  try {
    await leaveCircle(circleId.value)
    uni.showToast({ title: "已退出", icon: "success" })
    circle.value!.is_member = false
    circle.value!.member_count--
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
}

async function handleCreatePost() {
  if (!newPostContent.value.trim()) return
  try {
    await createCirclePost(circleId.value, { content: newPostContent.value.trim() })
    uni.showToast({ title: "发布成功", icon: "success" })
    newPostContent.value = ""
    showCreatePost.value = false
    posts.value = await listCirclePosts(circleId.value)
  } catch {
    uni.showToast({ title: "发布失败", icon: "none" })
  }
}

async function toggleLike(p: CirclePostData) {
  try {
    const res = await toggleCirclePostLike(p.id)
    p.is_liked = res.liked
    p.like_count += res.liked ? 1 : -1
  } catch {}
}

function showComments(p: CirclePostData) {
  selectedPost.value = p
  newComment.value = ""
  showCommentsPanel.value = true
  loadComments(p.id)
}

async function loadComments(postId: number) {
  try {
    comments.value = await listCirclePostComments(postId)
  } catch {}
}

async function handleCreateComment() {
  if (!newComment.value.trim() || !selectedPost.value) return
  try {
    await createCirclePostComment(selectedPost.value.id, newComment.value.trim())
    newComment.value = ""
    selectedPost.value.comment_count++
    await loadComments(selectedPost.value.id)
  } catch {
    uni.showToast({ title: "评论失败", icon: "none" })
  }
}

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = (now.getTime() - d.getTime()) / 1000
  if (diff < 60) return "刚刚"
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>
