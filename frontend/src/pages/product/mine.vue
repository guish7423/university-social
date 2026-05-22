<script setup lang="ts">
import { ref, onMounted } from "vue"
import { myProducts, deleteProduct, markSold, type ProductData } from "@/api/product"

const products = ref<ProductData[]>([])
const total = ref(0)
const active = ref(0)
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await myProducts()
    products.value = res.products
    total.value = res.total
    active.value = res.active
  } catch {}
  loading.value = false
})

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

async function handleSold(id: number) {
  await markSold(id)
  const res = await myProducts()
  products.value = res.products
  uni.showToast({ title: "已标记已出", icon: "success" })
}

async function handleDelete(id: number) {
  uni.showModal({
    title: "提示",
    content: "确定删除？",
    success: async (r) => {
      if (r.confirm) {
        await deleteProduct(id)
        const res = await myProducts()
        products.value = res.products
        uni.showToast({ title: "已删除", icon: "success" })
      }
    },
  })
}

function formatPrice(p: number) { return "¥" + p.toFixed(2) }
</script>

<template>
  <view class="container">
    <view class="stats-bar">
      <view class="stat-item"><text class="stat-num">{{ total }}</text><text class="stat-label">全部</text></view>
      <view class="stat-item"><text class="stat-num">{{ active }}</text><text class="stat-label">在售</text></view>
      <view class="stat-item"><text class="stat-num">{{ total - active }}</text><text class="stat-label">已出</text></view>
    </view>

    <view v-if="loading" class="loading-wrap"><u-loading mode="flower" size="48"></u-loading></view>

    <view v-else class="product-list">
      <view v-for="item in products" :key="item.id" class="product-item" @click="goDetail(item.id)">
        <image v-if="item.images?.[0]" :src="item.images[0]" class="item-img" mode="aspectFill" />
        <view class="item-info">
          <text class="item-title">{{ item.title }}</text>
          <text class="item-price">{{ formatPrice(item.price) }}</text>
          <view class="item-tags">
            <u-tag :text="item.status === 0 ? '在售' : '已出'" :type="item.status === 0 ? 'success' : 'error'" size="mini" shape="circle" />
            <u-tag :text="item.category" type="primary" size="mini" shape="circle" plain />
          </view>
        </view>
        <view class="item-actions" @click.stop>
          <text v-if="item.status === 0" class="action-link" @click.stop="handleSold(item.id)">标记已出</text>
          <text class="action-link danger" @click.stop="handleDelete(item.id)">删除</text>
        </view>
      </view>
    </view>

    <view v-if="!loading && products.length === 0" class="empty">
      <text>还没有发布过商品</text>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #f5f7fa; }
.container { padding-bottom: 30rpx; }

.stats-bar {
  display: flex;
  justify-content: space-around;
  padding: 30rpx;
  background: #fff;
  margin-bottom: 16rpx;
}
.stat-item { display: flex; flex-direction: column; align-items: center; }
.stat-num { font-size: 36rpx; font-weight: 700; color: #333; }
.stat-label { font-size: 24rpx; color: #999; margin-top: 4rpx; }

.loading { text-align: center; padding: 60rpx; color: #999; }

.product-item {
  display: flex;
  padding: 20rpx;
  margin: 0 20rpx 16rpx;
  background: #fff;
  border-radius: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
  &:active { transform: scale(0.98); }
}
.item-img { width: 160rpx; height: 160rpx; border-radius: 12rpx; margin-right: 16rpx; }
.item-info { flex: 1; display: flex; flex-direction: column; justify-content: center; gap: 8rpx; }
.item-title { font-size: 28rpx; font-weight: 500; color: #333; }
.item-price { font-size: 30rpx; font-weight: 700; color: #e74c3c; }
.item-tags { display: flex; gap: 8rpx; }
.item-actions {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 12rpx;
  padding-left: 16rpx;
}
.action-link {
  font-size: 24rpx;
  color: #C67A6A;
  padding: 8rpx 16rpx;
  border: 1rpx solid #C67A6A;
  border-radius: 8rpx;
  text-align: center;
  &.danger { color: #e74c3c; border-color: #e74c3c; }
}

.empty { text-align: center; padding: 100rpx 0; color: #999; font-size: 28rpx; }
</style>
