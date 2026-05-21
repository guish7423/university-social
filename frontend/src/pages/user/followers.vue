<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getFollowers, type FollowUser, followUser, unfollowUser, checkFollow } from "@/api/follow"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const users = ref<FollowUser[]>([])
const followState = ref<Record<number, boolean>>({})

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  const userId = Number(page?.options?.id || userStore.id)
  loadFollowers(userId)
})

async function loadFollowers(userId: number) {
  try {
    const data = await getFollowers(userId)
    users.value = data
    for (const u of data) {
      try {
        const res = await checkFollow(u.id)
        followState.value[u.id] = res.is_following
      } catch {}
    }
  } catch {}
}

async function toggleFollow(userId: number) {
  try {
    if (followState.value[userId]) {
      await unfollowUser(userId)
      followState.value[userId] = false
    } else {
      await followUser(userId)
      followState.value[userId] = true
    }
  } catch {}
}

function goUser(id: number) {
  uni.navigateTo({ url: `/pages/user/detail?id=${id}` })
}
</script>

<template>
  <view class="container">
    <view class="header">
      <text class="title">粉丝</text>
    </view>
    <view v-if="users.length === 0" class="empty-state">
      <text class="empty-text">暂无粉丝</text>
    </view>
    <view v-for="u in users" :key="u.id" class="user-item" @click="goUser(u.id)">
      <u-avatar :src="u.avatar" :text="u.nickname?.[0] || '?'" size="80" shape="circle" fontSize="32" bgColor="#C67A6A" />
      <view class="user-info">
        <text class="user-name">{{ u.nickname }}</text>
        <view class="user-meta">
          <u-tag v-if="u.is_verified" text="已认证" type="success" shape="circle" size="mini" />
          <text v-if="u.school" class="user-school">{{ u.school }}</text>
        </view>
      </view>
      <u-button
        v-if="u.id !== userStore.id"
        :type="followState[u.id] ? 'default' : 'primary'"
        :text="followState[u.id] ? '已关注' : '+ 关注'"
        size="mini"
        shape="circle"
        @click.stop="toggleFollow(u.id)"
      />
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: var(--color-canvas, #F7F4F0); }
.header { padding: 30rpx; background: var(--color-surface, #fff); }
.title { font-size: 34rpx; font-weight: 700; color: var(--ink, #1E2A3A); }
.empty-state { display: flex; justify-content: center; padding: 120rpx 0; }
.empty-text { font-size: 28rpx; color: var(--ink-tertiary, #B8C2CE); }
.user-item {
  display: flex; align-items: center; gap: 20rpx;
  padding: 24rpx 30rpx; background: var(--color-surface, #fff); border-bottom: 1rpx solid var(--hairline-light, #EAE6E0);
}
.user-info { flex: 1; }
.user-name { font-size: 30rpx; font-weight: 600; color: var(--ink, #1E2A3A); display: block; }
.user-meta { display: flex; align-items: center; gap: 8rpx; margin-top: 6rpx; }
.user-school { font-size: 24rpx; color: var(--ink-tertiary, #B8C2CE); }
</style>
