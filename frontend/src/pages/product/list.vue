<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listProducts, trendingProducts, productCategories, toggleProductLike, type ProductData } from "@/api/product"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const products = ref<ProductData[]>([])
const trending = ref<ProductData[]>([])
const categories = ref<string[]>([])
const currentCat = ref("全部")
const loading = ref(false)
const showTrending = ref(true)

onMounted(async () => {
  loading.value = true
  try {
    const [cats, trend] = await Promise.all([productCategories(), trendingProducts()])
    categories.value = cats
    trending.value = trend
    await loadProducts()
  } catch {}
  loading.value = false
})

async function loadProducts() {
  const cat = currentCat.value === "全部" ? "" : currentCat.value
  products.value = await listProducts({ category: cat })
  showTrending.value = currentCat.value === "全部"
}

function switchCat(cat: string) {
  currentCat.value = cat
  loadProducts()
}

async function handleLike(item: ProductData) {
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  const res = await toggleProductLike(item.id)
  item.is_liked = res.liked
  item.like_count += res.liked ? 1 : -1
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

function goCreate() {
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  uni.navigateTo({ url: "/pages/product/create" })
}

function formatPrice(p: number) {
  return "¥" + p.toFixed(2)
}

const conditionLabels: Record<string, string> = {
  "全新": "全新", "几乎全新": "几乎全新", "轻微使用": "轻微使用", "明显使用痕迹": "有使用痕迹",
}
</script>

<template>
  <view class="top-bar">
    <text class="top-bar-title">二手市场</text>
    <u-icon name="search" size="36" color="#fff" />
  </view>
  <view class="container">
    <view class="category-bar">
      <scroll-view scroll-x class="cat-scroll">
        <view
          v-for="cat in categories" :key="cat"
          class="cat-item" :class="{ active: currentCat === cat }"
          @click="switchCat(cat)"
        >
          <text>{{ cat }}</text>
        </view>
      </scroll-view>
    </view>

    <view v-if="showTrending && trending.length > 0" class="trending-section">
      <view class="section-title">
        <text class="title-icon">🔥</text>
        <text>热门闲置</text>
      </view>
      <scroll-view scroll-x class="trending-scroll">
        <view
          v-for="item in trending" :key="item.id"
          class="trending-card" @click="goDetail(item.id)"
        >
          <image v-if="item.images?.[0]" :src="item.images[0]" class="trending-img" mode="aspectFill" />
          <view class="trending-info">
            <text class="trending-title">{{ item.title }}</text>
            <text class="trending-price">{{ formatPrice(item.price) }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <view v-if="loading" class="loading"><text>加载中...</text></view>

    <view v-else class="product-grid">
      <view
        v-for="(item, idx) in products" :key="item.id"
        class="product-card"
        :class="'stagger-' + ((idx % 4) + 1)"
        @click="goDetail(item.id)"
      >
        <image v-if="item.images?.[0]" :src="item.images[0]" class="product-img" mode="aspectFill" />
        <image v-else class="product-img placeholder-img" src="/static/placeholder.png" mode="aspectFill" />
        <view class="product-body">
          <text class="product-title">{{ item.title }}</text>
          <view class="product-meta">
            <text class="product-condition">{{ item.condition }}</text>
          </view>
          <view class="product-footer">
            <text class="product-price">{{ formatPrice(item.price) }}</text>
            <view class="product-actions">
              <text class="like-btn" :class="{ liked: item.is_liked }" @click.stop="handleLike(item)">
                {{ item.is_liked ? '❤️' : '🤍' }} {{ item.like_count }}
              </text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <view v-if="!loading && products.length === 0" class="empty">
      <text>该分类暂无商品</text>
    </view>

    <view class="fab" @click="goCreate">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #f5f7fa; }
.container { padding-bottom: 120rpx; }

.category-bar {
  background: #fff;
  padding: 20rpx 0;
  position: sticky;
  top: 0;
  z-index: 10;
}
.cat-scroll { white-space: nowrap; padding: 0 20rpx; }
.cat-item {
  display: inline-flex;
  padding: 12rpx 28rpx;
  margin-right: 16rpx;
  border-radius: 30rpx;
  background: #f0f2f5;
  font-size: 26rpx;
  color: #666;
  transition: all 0.2s;
  &.active {
    background: linear-gradient(135deg, #C67A6A, #1E2A3A);
    color: #fff;
    font-weight: 600;
  }
}

.trending-section { padding: 24rpx 0; }
.section-title {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 0 30rpx;
  margin-bottom: 16rpx;
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}
.trending-scroll { padding: 0 20rpx; white-space: nowrap; }
.trending-card {
  display: inline-flex;
  flex-direction: column;
  width: 220rpx;
  margin-right: 20rpx;
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
}
.trending-img { width: 220rpx; height: 220rpx; }
.trending-info { padding: 12rpx; }
.trending-title {
  font-size: 24rpx;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}
.trending-price {
  font-size: 28rpx;
  font-weight: 700;
  color: #e74c3c;
}

.loading { text-align: center; padding: 60rpx; color: #999; }

.product-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
  padding: 20rpx;
}
.product-card {
  background: #fff;
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
  transition: all 0.2s;
  &:active { transform: scale(0.98); }
}
.product-img { width: 100%; height: 340rpx; }
.placeholder-img {
  background: linear-gradient(135deg, #C67A6A22, #1E2A3A22);
}
.product-body { padding: 16rpx; }
.product-title {
  font-size: 26rpx;
  font-weight: 500;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}
.product-meta { margin-top: 8rpx; }
.product-condition {
  font-size: 22rpx;
  color: #999;
  background: #f5f5f5;
  padding: 4rpx 12rpx;
  border-radius: 6rpx;
}
.product-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12rpx;
}
.product-price {
  font-size: 32rpx;
  font-weight: 700;
  color: #e74c3c;
}
.like-btn {
  font-size: 22rpx;
  color: #999;
}

.empty { text-align: center; padding: 100rpx 0; color: #999; font-size: 28rpx; }

.fab {
  position: fixed;
  right: 30rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #C67A6A, #1E2A3A);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 20rpx rgba(198,122,106,0.4);
  z-index: 100;
  &:active { transform: scale(0.9); }
}
.fab-icon { color: #fff; font-size: 48rpx; font-weight: 300; line-height: 1; }

.stagger-1 { animation: fadeInUp 0.4s ease-out 0s both; }
.stagger-2 { animation: fadeInUp 0.4s ease-out 0.08s both; }
.stagger-3 { animation: fadeInUp 0.4s ease-out 0.16s both; }
.stagger-4 { animation: fadeInUp 0.4s ease-out 0.24s both; }

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20rpx); }
  to { opacity: 1; transform: translateY(0); }
}

.top-bar {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 20rpx 32rpx; display: flex; align-items: center; justify-content: space-between;
  height: 88rpx;
}
.top-bar-title { font-size: 32rpx; font-weight: 700; color: #fff; }
.empty-title { font-size: 30rpx; font-weight: 600; color: #1E2A3A; margin-top: 12rpx; }
.empty-desc { font-size: 26rpx; color: #8E9BAB; margin-top: 8rpx; }
.empty-btn { margin-top: 16rpx; }

@media (min-width: 1024px) {
.product-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24rpx;
  padding: 32rpx;
}
.product-card { margin-bottom: 0; }
.fab { display: none; }
.container { max-width: 1400rpx; margin: 0 auto; }
}
</style>
