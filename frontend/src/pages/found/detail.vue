<script setup lang="ts">
import { ref, computed } from "vue"
import { getLostItem, deleteLostItem, updateLostItemStatus } from "@/api/found"
import { useUserStore } from "@/store/user"
import { onLoad } from "@dcloudio/uni-app"

const userStore = useUserStore()
const loading = ref(true)
const item = ref<any>(null)

const isOwner = computed(() => item.value && userStore.userInfo && item.value.user_id === userStore.userInfo.id)

onLoad((opts: any) => {
  if (opts.id) loadItem(parseInt(opts.id))
  else loading.value = false
})

async function loadItem(id: number) {
  loading.value = true
  try {
    const { data } = await getLostItem(id)
    item.value = data
  } catch { uni.showToast({ title: "加载失败", icon: "none" }) }
  loading.value = false
}

async function markSolved() {
  try {
    await updateLostItemStatus(item.value.id, 1)
    uni.showToast({ title: "已更新", icon: "success" })
    loadItem(item.value.id)
  } catch { uni.showToast({ title: "操作失败", icon: "none" }) }
}

async function doDelete() {
  uni.showModal({
    title: "确认", content: "确定删除？",
    success: async (res) => {
      if (res.confirm) {
        try {
          await deleteLostItem(item.value.id)
          uni.showToast({ title: "已删除", icon: "success" })
          setTimeout(() => uni.navigateBack(), 1000)
        } catch { uni.showToast({ title: "删除失败", icon: "none" }) }
      }
    },
  })
}
</script>

<template>
  <view class="page">
    <u-loading mode="flower" size="60" v-if="loading" />
    <template v-else>
    <view class="header" v-if="item">
      <view class="header-top">
        <view class="tag" :class="item.category === 'lost' ? 'tag-lost' : 'tag-found'">
          {{ item.category === 'lost' ? '◉ 寻物' : '◎ 招领' }}
        </view>
        <view class="status-tag" v-if="item.status !== 0">
          {{ item.status === 1 ? '已解决' : '已关闭' }}
        </view>
      </view>
      <text class="title">{{ item.title }}</text>
      <text class="desc" v-if="item.description">{{ item.description }}</text>

      <view class="info-grid">
        <view class="info-item" v-if="item.item_type">
          <text class="info-label">物品类型</text>
          <text class="info-value">{{ item.item_type }}</text>
        </view>
        <view class="info-item" v-if="item.location">
          <text class="info-label">地点</text>
          <text class="info-value">{{ item.location }}</text>
        </view>
        <view class="info-item" v-if="item.contact">
          <text class="info-label">联系方式</text>
          <text class="info-value">{{ item.contact }}</text>
        </view>
      </view>

      <view class="footer">
        <view class="author">
          <u-avatar :src="item.avatar" size="40"></u-avatar>
          <text class="name">{{ item.nickname || '匿名' }}</text>
        </view>
        <text class="time">{{ item.created_at }}</text>
      </view>
    </view>

    <view class="actions" v-if="item && isOwner">
      <u-button type="warning" shape="circle" @click="markSolved" v-if="item.status === 0">
        {{ item.category === 'lost' ? '已找到' : '已完成' }}
      </u-button>
      <u-button type="error" shape="circle" plain @click="doDelete">删除</u-button>
    </view>
    </template>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 24rpx; }
.header { background: var(--color-surface, #fff); border-radius: 22rpx; padding: 32rpx; margin-bottom: 24rpx; }
.header-top { display: flex; flex-direction: row; align-items: center; margin-bottom: 16rpx; }
.tag { font-size: 22rpx; padding: 6rpx 20rpx; border-radius: 8rpx; font-weight: 500; margin-right: 12rpx; }
.tag-lost { background: var(--color-surface-1, #F0EDE8); color: var(--ink-muted, #5C6B7E); }
.tag-found { background: rgba(139,168,137,0.1); color: var(--sage, #8BA889); }
.status-tag { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); background: var(--color-surface-1, #F0EDE8); padding: 4rpx 14rpx; border-radius: 8rpx; }
.title { font-size: 34rpx; font-weight: 700; color: var(--ink, #1E2A3A); display: block; margin-bottom: 16rpx; }
.desc { font-size: 28rpx; color: var(--ink-muted, #5C6B7E); line-height: 1.6; display: block; margin-bottom: 24rpx; }
.info-grid { background: var(--color-surface-1, #F0EDE8); border-radius: 16rpx; padding: 20rpx; margin-bottom: 24rpx; }
.info-item { display: flex; flex-direction: row; justify-content: space-between; padding: 12rpx 0; }
.info-label { font-size: 26rpx; color: var(--ink-tertiary, #B8C2CE); }
.info-value { font-size: 26rpx; color: var(--ink, #1E2A3A); font-weight: 500; }
.footer { display: flex; flex-direction: row; justify-content: space-between; align-items: center; }
.author { display: flex; flex-direction: row; align-items: center; gap: 12rpx; }
.name { font-size: 26rpx; color: var(--brand-primary, #C67A6A); }
.time { font-size: 24rpx; color: var(--ink-tertiary, #B8C2CE); }
.actions { display: flex; flex-direction: row; gap: 20rpx; }
</style>
