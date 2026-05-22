<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { getProduct, deleteProduct, markSold, toggleProductLike, listProductComments, createProductComment, type ProductData, type ProductComment } from "@/api/product"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const product = ref<ProductData | null>(null)
const comments = ref<ProductComment[]>([])
const loading = ref(false)
const commentText = ref("")
const sending = ref(false)
const showActions = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const id = parseInt((uni as any).getSystemInfoSync().platform === 'devtools' ? '' : '')
    // get id from URL params
    const pages = getCurrentPages()
    const page = pages[pages.length - 1] as any
    const pid = parseInt(page?.$page?.query?.id || '0')
    if (!pid) { uni.showToast({ title: '参数错误', icon: 'none' }); return }
    const [p, c] = await Promise.all([getProduct(pid), listProductComments(pid)])
    product.value = p
    comments.value = c
  } catch { uni.showToast({ title: '商品不存在', icon: 'none' }) }
  loading.value = false
})

async function handleLike() {
  if (!product.value) return
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  const res = await toggleProductLike(product.value.id)
  product.value.is_liked = res.liked
  product.value.like_count += res.liked ? 1 : -1
}

async function handleDelete() {
  if (!product.value) return
  uni.showModal({
    title: "提示",
    content: "确定删除此商品？",
    success: async (r) => {
      if (r.confirm) {
        await deleteProduct(product.value!.id)
        uni.showToast({ title: "已删除", icon: "success" })
        setTimeout(() => uni.navigateBack(), 1000)
      }
    },
  })
}

async function handleSold() {
  if (!product.value) return
  await markSold(product.value.id)
  product.value.status = 1
  uni.showToast({ title: "已标记为已出", icon: "success" })
}

async function sendComment() {
  if (!product.value || !commentText.value.trim()) return
  sending.value = true
  try {
    await createProductComment(product.value.id, commentText.value.trim())
    commentText.value = ""
    comments.value = await listProductComments(product.value.id)
  } catch {}
  sending.value = false
}

function formatPrice(p: number) { return "¥" + p.toFixed(2) }
function formatTime(t: string) { return t ? t.slice(0, 16).replace("T", " ") : "" }
function originalPrice() {
  if (!product.value?.original_price) return ""
  return "原价 " + formatPrice(product.value.original_price)
}
</script>

<template>
  <view class="container">
    <view v-if="loading" class="loading-wrap"><u-loading mode="flower" size="48"></u-loading></view>

    <template v-else-if="product">
      <scroll-view scroll-y class="detail-scroll">
        <swiper
          v-if="product.images?.length"
          class="image-swiper"
          indicator-dots indicator-color="#ffffff88" indicator-active-color="#C67A6A"
          autoplay circular
        >
          <swiper-item v-for="(img, i) in product.images" :key="i">
            <image :src="img" class="swiper-img" mode="aspectFill" />
          </swiper-item>
        </swiper>

        <view class="info-section">
          <text class="price">{{ formatPrice(product.price) }}</text>
          <text v-if="originalPrice()" class="original-price">{{ originalPrice() }}</text>
          <text class="title">{{ product.title }}</text>
          <view class="tags">
            <u-tag :text="product.category" type="primary" size="mini" shape="circle" />
            <u-tag :text="product.condition" type="warning" size="mini" shape="circle" plain />
            <u-tag v-if="product.status === 1" text="已出" type="error" size="mini" shape="circle" />
          </view>
          <text v-if="product.description" class="desc">{{ product.description }}</text>
        </view>

        <view class="contact-section">
          <text class="section-label">联系方式</text>
          <text class="contact">{{ product.contact || '请私信联系' }}</text>
        </view>

        <view class="stats-section">
          <view class="stat" @click="handleLike">
            <text :class="['stat-icon', { liked: product.is_liked }]">
              {{ product.is_liked ? '❤️' : '🤍' }}
            </text>
            <text class="stat-num">{{ product.like_count }}</text>
          </view>
          <view class="stat">
            <text class="stat-icon">💬</text>
            <text class="stat-num">{{ product.comment_count }}</text>
          </view>
          <view class="stat">
            <text class="stat-icon">📅</text>
            <text class="stat-num">{{ formatTime(product.created_at).slice(5, 10) }}</text>
          </view>
        </view>

        <view class="comments-section">
          <text class="section-title">留言 ({{ comments.length }})</text>
          <view v-if="comments.length === 0" class="no-comments">
            <text>暂无留言</text>
          </view>
          <view v-for="c in comments" :key="c.id" class="comment-item">
            <text class="comment-content">{{ c.content }}</text>
            <text class="comment-time">{{ formatTime(c.created_at) }}</text>
          </view>
        </view>
      </scroll-view>

      <view class="bottom-bar">
        <view class="action-btn" v-if="!product.is_owner" @click="handleLike">
          <text>{{ product.is_liked ? '❤️' : '🤍' }}</text>
        </view>
        <view class="action-btn" @click="showActions = !showActions">
          <text>•••</text>
        </view>
        <view class="comment-input-area">
          <input
            v-model="commentText"
            class="comment-input"
            placeholder="留言..."
            confirm-type="send"
            @confirm="sendComment"
          />
          <text class="send-btn" :class="{ disabled: !commentText.trim() || sending }" @click="sendComment">
            发送
          </text>
        </view>
      </view>

      <u-action-sheet v-model="showActions" :list="[
        ...(product.is_owner ? [
          { text: product.status === 0 ? '标记已出' : '', },
          { text: '删除商品' },
        ] : []),
        { text: '举报' },
        { text: '取消' },
      ].filter(i => i.text)" @click="(i: number) => {
        showActions = false
        const actions = (product.is_owner ? [product.status === 0 ? handleSold : null, handleDelete].filter(Boolean) : []).concat(null as any)
        if (i < actions.length && actions[i]) { (actions[i] as any)() }
      }" />
    </template>
  </view>
</template>

<style lang="scss">
page { background: #F7F4F0; }
.container { height: 100vh; display: flex; flex-direction: column; }

.loading { text-align: center; padding: 100rpx 0; color: #999; }

.detail-scroll { flex: 1; }

.image-swiper { height: 600rpx; }
.swiper-img { width: 100%; height: 100%; }

.info-section {
  padding: 30rpx;
  background: #fff;
  margin-bottom: 16rpx;
}
.price { font-size: 48rpx; font-weight: 700; color: #C67A6A; }
.original-price {
  font-size: 24rpx;
  color: #999;
  text-decoration: line-through;
  margin-left: 16rpx;
}
.title {
  display: block;
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin: 16rpx 0;
}
.tags { display: flex; gap: 12rpx; margin-bottom: 16rpx; }
.desc {
  display: block;
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
}

.contact-section {
  padding: 24rpx 30rpx;
  background: #fff;
  margin-bottom: 16rpx;
  display: flex;
  align-items: center;
  gap: 20rpx;
}
.section-label { font-size: 26rpx; color: #999; }
.contact { font-size: 28rpx; color: #333; }

.stats-section {
  display: flex;
  justify-content: space-around;
  padding: 24rpx;
  background: #fff;
  margin-bottom: 16rpx;
}
.stat { display: flex; align-items: center; gap: 8rpx; }
.stat-icon { font-size: 32rpx; &.liked { animation: likePulse 0.3s; } }
.stat-num { font-size: 26rpx; color: #666; }

.comments-section {
  padding: 30rpx;
  background: #fff;
  min-height: 200rpx;
}
.section-title { font-size: 28rpx; font-weight: 600; color: #333; margin-bottom: 20rpx; display: block; }
.no-comments { text-align: center; padding: 40rpx 0; color: #ccc; font-size: 26rpx; }
.comment-item { padding: 16rpx 0; border-bottom: 1rpx solid #f0f0f0; }
.comment-content { font-size: 26rpx; color: #333; display: block; }
.comment-time { font-size: 22rpx; color: #ccc; margin-top: 6rpx; display: block; }

.bottom-bar {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background: #fff;
  border-top: 1rpx solid #eee;
  gap: 16rpx;
}
.action-btn { padding: 8rpx 16rpx; font-size: 32rpx; }
.comment-input-area {
  flex: 1;
  display: flex;
  align-items: center;
  background: #F7F4F0;
  border-radius: 40rpx;
  padding: 0 24rpx;
}
.comment-input { flex: 1; height: 64rpx; font-size: 26rpx; }
.send-btn {
  padding: 8rpx 20rpx;
  color: #C67A6A;
  font-size: 26rpx;
  font-weight: 500;
  &.disabled { color: #ccc; }
}

@keyframes likePulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.3); }
  100% { transform: scale(1); }
}
</style>
