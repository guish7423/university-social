<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listLostItems, searchLostItems, type LostItem } from "@/api/found"

const items = ref<LostItem[]>([])
const loading = ref(true)
const keyword = ref("")

const categories = ["全部", "校园卡", "电子产品", "书籍", "钥匙", "衣物", "其他"]
const activeCat = ref(0)

onMounted(async () => {
  try { items.value = await listLostItems() } catch (e) { console.error(e) }
  loading.value = false
})

async function onSearch() {
  loading.value = true
  try {
    items.value = keyword.value
      ? await searchLostItems(keyword.value)
      : await listLostItems()
  } catch (e) { console.error(e) }
  loading.value = false
}

function goDetail(id: number) { uni.navigateTo({ url: `/pages/found/detail?id=${id}` }) }
function goCreate() { uni.navigateTo({ url: "/pages/found/create" }) }
function goFound() { uni.switchTab({ url: "/pages/found/list" }) }
</script>

<template>
  <view class="page">
    <view class="header">
      <view class="header-bg" />
      <view class="header-content">
        <text class="header-title">失物招领</text>
        <text class="header-sub">捡到物品 · 寻找失主</text>
      </view>
    </view>

    <view class="search-bar">
      <u-icon name="search" size="32" color="#B8C2CE" />
      <input v-model="keyword" class="search-input" placeholder="搜索失物" @confirm="onSearch" />
      <text v-if="keyword" class="search-clear" @click="keyword = ''; onSearch()">清除</text>
    </view>

    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="48" />
    </view>

    <view v-else-if="items.length" class="list">
      <view v-for="(item, i) in items" :key="item.id"
        class="card" hover-class="card-hover"
        :style="{ animationDelay: (i * 0.06) + 's' }"
        @click="goDetail(item.id)">
        <view class="card-left">
          <image :src="item.images?.[0] || '/static/default.png'" mode="aspectFill" class="card-img" />
        </view>
        <view class="card-right">
          <text class="card-title">{{ item.title }}</text>
          <text class="card-desc">{{ item.description?.slice(0, 50) }}{{ item.description?.length > 50 ? '...' : '' }}</text>
          <view class="card-meta">
            <text class="card-location">{{ item.location || '未知位置' }}</text>
            <u-icon name="arrow-right" size="24" color="#B8C2CE" />
          </view>
        </view>
      </view>
    </view>

    <view v-else class="empty-state">
      <u-icon name="empty-address" size="160" color="#EAE6E0" />
      <text class="empty-title">还没有失物信息</text>
      <text class="empty-desc">来发布第一条信息吧</text>
      <u-button type="primary" shape="circle" class="empty-btn" @click="goCreate">发布信息</u-button>
    </view>
  </view>
</template>

<style lang="scss" scoped>
page { background: #F7F4F0; }
.page { min-height: 100vh; }

.header {
  position: relative; padding: 80rpx 32rpx 60rpx; overflow: hidden;
}
.header-bg {
  position: absolute; inset: 0;
  background: linear-gradient(135deg, #1E2A3A, #2A3A4E);
  z-index: 0;
  &::after {
    content: ''; position: absolute; top: -200rpx; right: -80rpx;
    width: 400rpx; height: 400rpx; border-radius: 50%;
    background: radial-gradient(circle, rgba(198,122,106,0.15), transparent 70%);
  }
}
.header-content { position: relative; z-index: 1; }
.header-title { font-size: 36rpx; font-weight: 700; color: #fff; display: block; }
.header-sub { font-size: 24rpx; color: rgba(255,255,255,0.5); margin-top: 8rpx; }

.search-bar {
  display: flex; align-items: center; gap: 12rpx;
  background: #fff; border-radius: 40rpx; padding: 16rpx 28rpx;
  margin: -24rpx 24rpx 24rpx; position: relative; z-index: 2;
  box-shadow: 0 4rpx 16rpx rgba(30,42,58,0.06);
}
.search-input { flex: 1; font-size: 26rpx; color: #1E2A3A; }
.search-clear { font-size: 22rpx; color: #C67A6A; cursor: pointer; }

.loading-state { display: flex; justify-content: center; padding: 200rpx 0; }

.list { padding: 0 20rpx 40rpx; display: flex; flex-direction: column; gap: 16rpx; }

.card {
  background: #fff; border-radius: 16rpx; padding: 20rpx;
  display: flex; gap: 16rpx; cursor: pointer; animation: fadeInUp 0.4s both;
  box-shadow: 0 2rpx 10rpx rgba(30,42,58,0.04);
}
.card-hover { transform: scale(0.98); }
.card-left { width: 180rpx; height: 180rpx; flex-shrink: 0; }
.card-img { width: 100%; height: 100%; border-radius: 12rpx; background: #E8E4DE; }
.card-right { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 8rpx; }
.card-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }
.card-desc { font-size: 24rpx; color: #8E9BAB; line-height: 1.5; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; }
.card-meta { display: flex; align-items: center; justify-content: space-between; margin-top: auto; }
.card-location { font-size: 22rpx; color: #B8C2CE; }

.empty-state {
  display: flex; flex-direction: column; align-items: center; gap: 16rpx;
  padding: 200rpx 60rpx;
}
.empty-title { font-size: 30rpx; font-weight: 600; color: #1E2A3A; }
.empty-desc { font-size: 26rpx; color: #8E9BAB; }
.empty-btn { margin-top: 16rpx; }

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20rpx); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
