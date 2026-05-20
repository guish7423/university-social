<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listFriendRequests, acceptFriendRequest, rejectFriendRequest, type UserInfo } from "@/api/social"

const requests = ref<UserInfo[]>([])
const loading = ref(true)

onMounted(() => loadRequests())

async function loadRequests() {
  loading.value = true
  try {
    requests.value = await listFriendRequests()
  } catch {}
  loading.value = false
}

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
  <view class="container">
    <u-loading v-if="loading" mode="circle" size="36" class="loading" />
    <u-empty v-else-if="requests.length === 0" text="暂无好友请求" icon="character" />

    <view v-for="r in requests" :key="r.id" class="request-item">
      <u-avatar :src="r.avatar" size="72" :name="r.nickname?.[0] || '?'" />
      <view class="info">
        <text class="name">{{ r.nickname }}</text>
        <text v-if="r.school" class="school">{{ r.school }}</text>
      </view>
      <view class="actions">
        <u-button size="small" type="primary" class="action-btn" @click="accept(r.id)">接受</u-button>
        <u-button size="small" type="default" class="action-btn" @click="reject(r.id)">拒绝</u-button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; }
.loading { margin-top: 200rpx; }
.request-item {
  display: flex; align-items: center; gap: 20rpx;
  padding: 24rpx 30rpx; background: #fff;
  border-bottom: 1rpx solid #eee;
}
.info { flex: 1; }
.name { font-size: 28rpx; font-weight: 600; display: block; }
.school { font-size: 22rpx; color: #999; }
.actions { display: flex; gap: 12rpx; }
.action-btn { min-width: 100rpx; }
</style>
