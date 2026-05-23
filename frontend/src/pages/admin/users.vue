<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listUsers, banUser, unbanUser } from "@/api/admin"

const users = ref<any[]>([])
const loading = ref(true)
const offset = ref(0)
const hasMore = ref(true)

async function load() {
  try {
    const data = await listUsers(offset.value)
    if (data.length < 50) hasMore.value = false
    users.value = offset.value === 0 ? data : [...users.value, ...data]
  } catch (e) { console.error(e) }
  loading.value = false
}

onMounted(load)

async function doBan(id: number, banned: boolean) {
  uni.showLoading({ title: "操作中" })
  try {
    if (banned) {
      await unbanUser(id)
      uni.showToast({ title: "已解封", icon: "success" })
    } else {
      await banUser(id)
      uni.showToast({ title: "已封禁", icon: "success" })
    }
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
  uni.hideLoading()
}
</script>

<template>
  <view class="page">
    <view class="search-bar">
      <text class="page-title">用户管理 ({{ users.length }})</text>
    </view>

    <view class="list">
      <view v-for="u in users" :key="u.id" class="list-item">
        <u-avatar :src="u.avatar" :text="u.nickname?.[0] || '?'" size="72" shape="circle" />
        <view class="list-item-body">
          <text class="list-item-title">{{ u.nickname }}</text>
          <text class="list-item-desc">{{ u.school || '未设置' }} · {{ u.created_at?.slice(0,10) }}</text>
        </view>
        <view class="tag" :class="{ verified: u.is_verified }">
          {{ u.is_verified ? '已认证' : '未认证' }}
        </view>
      </view>
    </view>

    <view v-if="hasMore" class="load-more" @click="offset += 50; load()">
      加载更多
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 24rpx 16rpx; }
.search-bar { margin-bottom: 24rpx; }
.page-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }

.list { background: #fff; border-radius: 14rpx; overflow: hidden; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); }
.list-item {
  display: flex; align-items: center; gap: 16rpx;
  padding: 24rpx; border-bottom: 1rpx solid #EAE6E0;
  &:last-child { border-bottom: none; }
}
.list-item-body { flex: 1; min-width: 0; }
.list-item-title { font-size: 26rpx; color: #1E2A3A; font-weight: 500; display: block; }
.list-item-desc { font-size: 22rpx; color: #8E9BAB; margin-top: 4rpx; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.tag {
  font-size: 20rpx; padding: 4rpx 12rpx; border-radius: 6rpx;
  background: #F0EDE8; color: #8E9BAB;
}
.tag.verified { background: #E8F5E9; color: #388E3C; }

.load-more {
  text-align: center; padding: 32rpx; color: #C67A6A; font-size: 26rpx; cursor: pointer;
}
</style>
