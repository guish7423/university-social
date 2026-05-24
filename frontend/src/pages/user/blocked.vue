<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getBlockedUsers, unblockUser } from "@/api/block"
import type { UserInfo } from "@/api/user"

const users = ref<UserInfo[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    users.value = await getBlockedUsers()
  } catch (e) { console.error(e) }
  loading.value = false
})

async function handleUnblock(user: UserInfo) {
  try {
    await unblockUser(user.id)
    users.value = users.value.filter(u => u.id !== user.id)
    uni.showToast({ title: "已取消屏蔽", icon: "none" })
  } catch (e) { console.error(e) }
}

function goUserDetail(id: number) {
  uni.navigateTo({ url: "/pages/user/detail?id=" + id })
}
</script>

<template>
  <view class="page">
    <u-loading mode="flower" size="60" v-if="loading" />
    <template v-else>
      <view v-if="users.length === 0" class="empty">
        <text class="empty-text">暂无屏蔽用户</text>
      </view>
      <view v-else class="list">
        <view
          v-for="user in users"
          :key="user.id"
          class="user-card"
          @click="goUserDetail(user.id)"
        >
          <u-avatar :src="user.avatar" :text="user.nickname?.[0] || '?'" size="80" shape="circle" />
          <view class="user-info">
            <text class="user-name">{{ user.nickname }}</text>
            <text class="user-school">{{ user.school || '' }}</text>
          </view>
          <view class="unblock-btn" @click.stop="handleUnblock(user)">取消屏蔽</view>
        </view>
      </view>
    </template>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 20rpx; }

.empty {
  display: flex;
  justify-content: center;
  padding-top: 200rpx;
}
.empty-text { font-size: 28rpx; color: #8E9BAB; }

.list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.user-name { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }
.user-school { font-size: 24rpx; color: #8E9BAB; }

.unblock-btn {
  padding: 12rpx 28rpx;
  border-radius: 30rpx;
  font-size: 24rpx;
  color: #C67A6A;
  background: rgba(198,122,106,0.1);
  white-space: nowrap;
  transition: all 0.2s;
  &:active { transform: scale(0.97); }
}
</style>
