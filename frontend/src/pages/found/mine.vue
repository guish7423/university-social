<script setup lang="ts">
import { ref, onMounted } from "vue"
import { myLostItems, listLostItems } from "@/api/found"

const tab = ref("all")
const items = ref<any[]>([])
const loading = ref(false)
const loadStatus = ref("loadmore")
const offset = ref(0)
const hasMore = ref(true)
const limit = 20

onMounted(() => loadItems())

function getCategory() {
  if (tab.value === "all") return undefined
  return tab.value
}

async function loadItems() {
  if (loading.value || !hasMore.value) return
  loading.value = true
  loadStatus.value = "loading"
  try {
    const category = getCategory()
    const { data } = category
      ? await listLostItems(category, -1, offset.value, limit)
      : await myLostItems(offset.value, limit)
    if (offset.value === 0) items.value = data
    else items.value = [...items.value, ...data]
    hasMore.value = data.length >= limit
    offset.value += data.length
    loadStatus.value = hasMore.value ? "loadmore" : "nomore"
  } catch { loadStatus.value = "loadmore" }
  finally { loading.value = false }
}

function switchTab(t: string) {
  tab.value = t
  offset.value = 0
  hasMore.value = true
  items.value = []
  loadItems()
}

function goDetail(id: number) {
  uni.navigateTo({ url: "/pages/found/detail?id=" + id })
}
</script>

<template>
  <view class="page">
    <view class="tabs">
      <view class="tab" :class="{ active: tab === 'all' }" @click="switchTab('all')">全部</view>
      <view class="tab" :class="{ active: tab === 'lost' }" @click="switchTab('lost')">◉ 寻物</view>
      <view class="tab" :class="{ active: tab === 'found' }" @click="switchTab('found')">◎ 招领</view>
    </view>

    <view class="list" v-if="items.length">
      <view class="card" v-for="(item, idx) in items" :key="item.id" :class="'stagger-' + ((idx % 8) + 1)" @click="goDetail(item.id)">
        <view class="card-top">
          <view class="tag" :class="item.category === 'lost' ? 'tag-lost' : 'tag-found'">
            {{ item.category === 'lost' ? '寻物' : '招领' }}
          </view>
          <view class="status-tag" v-if="item.status !== 0">
            {{ item.status === 1 ? '已解决' : '已关闭' }}
          </view>
        </view>
        <text class="card-title">{{ item.title }}</text>
        <text class="card-desc" v-if="item.description">{{ item.description }}</text>
        <view class="card-footer">
          <text class="time">{{ item.created_at }}</text>
        </view>
      </view>
    </view>
    <u-empty text="暂无发布" mode="list" v-else-if="!loading"></u-empty>
    <u-loadmore :status="loadStatus" v-if="items.length" />
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 0 24rpx 120rpx; }
.tabs { display: flex; flex-direction: row; padding: 20rpx 0; gap: 20rpx; }
.tab { font-size: 28rpx; color: var(--ink-tertiary, #B8C2CE); padding: 8rpx 24rpx; border-radius: 30rpx; background: var(--color-surface-1, #F0EDE8); }
.tab.active { background: linear-gradient(135deg, #C67A6A, #1E2A3A); color: #fff; }
.list { padding-top: 10rpx; }
.card { background: var(--color-surface, #fff); border-radius: 22rpx; padding: 28rpx; margin-bottom: 20rpx; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.06); }
.card-top { display: flex; flex-direction: row; align-items: center; margin-bottom: 12rpx; }
.tag { font-size: 22rpx; padding: 4rpx 16rpx; border-radius: 8rpx; font-weight: 500; margin-right: 12rpx; }
.tag-lost { background: var(--color-surface-1, #F0EDE8); color: var(--ink-muted, #5C6B7E); }
.tag-found { background: rgba(139,168,137,0.1); color: var(--sage, #8BA889); }
.status-tag { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); background: var(--color-surface-1, #F0EDE8); padding: 4rpx 14rpx; border-radius: 8rpx; }
.card-title { font-size: 30rpx; font-weight: 600; color: var(--ink, #1E2A3A); display: block; margin-bottom: 8rpx; }
.card-desc { font-size: 26rpx; color: var(--ink-muted, #5C6B7E); display: block; margin-bottom: 12rpx; }
.card-footer { display: flex; flex-direction: row; justify-content: flex-end; }
.time { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); }
</style>
