<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listFriendRequests, acceptFriendRequest, rejectFriendRequest, type UserInfo } from "@/api/social"

const requests = ref<UserInfo[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    requests.value = await listFriendRequests()
  } catch (e) { console.error(e) }
  loading.value = false
})

async function accept(id: number) {
  try {
    await acceptFriendRequest(id)
    requests.value = requests.value.filter(r => r.id !== id)
    uni.showToast({ title: "已添加为好友", icon: "success" })
  } catch {
    uni.showToast({ title: "操作失败", icon: "error" })
  }
}

async function reject(id: number) {
  try {
    await rejectFriendRequest(id)
    requests.value = requests.value.filter(r => r.id !== id)
    uni.showToast({ title: "已拒绝", icon: "none" })
  } catch {
    uni.showToast({ title: "操作失败", icon: "error" })
  }
}
</script>

<template>
  <view class="page">
    <u-loading v-if="loading" mode="circle" size="36" class="loading" />
    <u-empty v-else-if="requests.length === 0" text="暂无好友请求" icon="character" />

    <view v-for="r in requests" :key="r.id" class="request-item">
      <u-avatar :src="r.avatar" size="72" :name="r.nickname?.[0] || '?'" />
      <view class="info">
        <text class="name">{{ r.nickname }}</text>
        <text v-if="r.school" class="school">{{ r.school }}</text>
      </view>
      <view class="actions">
        <u-button size="small" type="primary" @click="accept(r.id)">接受</u-button>
        <u-button size="small" type="default" @click="reject(r.id)">拒绝</u-button>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); }
.loading { margin-top: 200rpx; }
.request-item {
  display: flex; align-items: center; gap: 20rpx;
  padding: 24rpx 30rpx; background: var(--color-surface, #fff);
  border-bottom: 1rpx solid var(--hairline-light, #EAE6E0);
}
.info { flex: 1; }
.name { font-size: 28rpx; font-weight: 600; color: var(--ink, #1E2A3A); display: block; }
.school { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); }
.actions { display: flex; gap: 12rpx; }
</style>
