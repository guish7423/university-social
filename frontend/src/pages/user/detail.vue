<script setup lang="ts">
import { ref, computed } from "vue"
import { onLoad } from "@dcloudio/uni-app"
import { useUserStore } from "@/store/user"
import { getUserInfo } from "@/api/user"
import { getFollowCounts, checkFollow, followUser, unfollowUser } from "@/api/follow"

const userStore = useUserStore()
const userId = ref(0)
const user = ref<any>(null)
const isFollowing = ref(false)
const counts = ref({ following: 0, followers: 0, posts: 0 })
const loading = ref(true)

const isSelf = computed(() => userId.value === userStore.id)

onLoad((opts: any) => {
  if (opts.id) {
    userId.value = parseInt(opts.id)
    Promise.all([loadUser(), loadCounts()]).finally(() => { loading.value = false })
    if (!isSelf.value) checkFollowing()
  }
})

async function loadUser() {
  try {
    user.value = await getUserInfo(userId.value)
  } catch {}
}

async function loadCounts() {
  try {
    counts.value = await getFollowCounts(userId.value)
  } catch {}
}

async function checkFollowing() {
  try {
    const res = await checkFollow(userId.value)
    isFollowing.value = res.is_following
  } catch {}
}

async function toggleFollow() {
  try {
    if (isFollowing.value) {
      await unfollowUser(userId.value)
      isFollowing.value = false
      counts.value.followers--
    } else {
      await followUser(userId.value)
      isFollowing.value = true
      counts.value.followers++
    }
  } catch {}
}

function goChat() {
  uni.navigateTo({ url: `/pages/chat/detail?userId=${userId.value}&nickname=${encodeURIComponent(user.value?.nickname || '')}` })
}
</script>

<template>
  <view class="page">
    <u-loading mode="flower" size="60" v-if="loading" />
    <template v-else>
      <view class="profile-header">
        <u-avatar :src="user?.avatar" :text="user?.nickname?.[0] || '?'" size="140" shape="circle" />
        <text class="profile-name">{{ user?.nickname }}</text>
        <text class="profile-school">{{ user?.school || '' }}</text>

        <view class="counts">
          <view class="count-item">
            <text class="count-num">{{ counts.posts }}</text>
            <text class="count-label">帖子</text>
          </view>
          <view class="count-item">
            <text class="count-num">{{ counts.following }}</text>
            <text class="count-label">关注</text>
          </view>
          <view class="count-item">
            <text class="count-num">{{ counts.followers }}</text>
            <text class="count-label">粉丝</text>
          </view>
        </view>

        <view v-if="!isSelf" class="follow-btn" :class="{ following: isFollowing }" @click="toggleFollow">
          {{ isFollowing ? '已关注' : '+ 关注' }}
        </view>
        <view class="msg-btn" @click="goChat">发消息</view>
      </view>
    </template>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); }
.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  padding: 80rpx 32rpx 48rpx;
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
}
.profile-name { font-size: 36rpx; font-weight: 600; color: #fff; }
.profile-school { font-size: 26rpx; color: rgba(255,255,255,0.5); font-weight: 300; }
.counts { display: flex; gap: 60rpx; margin-top: 16rpx; }
.count-item { display: flex; flex-direction: column; align-items: center; gap: 4rpx; }
.count-num { font-size: 32rpx; font-weight: 600; color: #fff; }
.count-label { font-size: 22rpx; color: rgba(255,255,255,0.5); }
.follow-btn {
  padding: 12rpx 40rpx;
  border-radius: 30rpx;
  font-size: 26rpx;
  font-weight: 500;
  background: #C67A6A;
  color: #fff;
  margin-top: 8rpx;
  transition: all 0.2s;
  &:active { transform: scale(0.97); }
  &.following { background: rgba(255,255,255,0.15); color: rgba(255,255,255,0.75); }
}

.msg-btn {
  padding: 12rpx 40rpx;
  border-radius: 30rpx;
  font-size: 26rpx;
  font-weight: 500;
  background: rgba(255,255,255,0.15);
  color: #fff;
  margin-top: 8rpx;
  transition: all 0.2s;
  &:active { transform: scale(0.97); }
}
</style>
