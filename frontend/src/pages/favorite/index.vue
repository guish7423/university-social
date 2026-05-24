<template>
  <view class="container">
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
    </view>
    <template v-else>
      <u-empty v-if="!list.length" mode="favor" text="还没有收藏" :iconSize="120" />
      <view v-for="item in list" :key="item.id" class="favorite-item" @click="goPost(item.post_id)">
        <view class="item-top">
          <text class="item-title">{{ item.post?.content?.slice(0, 60) || '帖子' }}</text>
          <u-icon name="trash" size="32" color="#B8C2CE" @click.stop="handleRemove(item.post_id)" />
        </view>
        <text class="item-meta">{{ formatTime(item.created_at) }} · {{ item.post?.like_count || 0 }} 赞</text>
      </view>
      <u-loadmore v-if="hasMore" :status="loadStatus" @loadmore="loadMore" />
    </template>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listFavorites, removeFavorite, type FavoriteData } from "@/api/favorite"

const list = ref<FavoriteData[]>([])
const loading = ref(true)
const hasMore = ref(true)
const offset = ref(0)
const loadStatus = ref<"loadmore" | "loading" | "nomore">("loadmore")

onMounted(() => loadMore())

async function loadMore() {
  if (loadStatus.value === "loading" || !hasMore.value) return
  loadStatus.value = "loading"
  try {
    const res = await listFavorites(offset.value)
    list.value.push(...res)
    offset.value += res.length
    hasMore.value = res.length >= 20
    loadStatus.value = res.length >= 20 ? "loadmore" : "nomore"
  } catch { loadStatus.value = "loadmore" }
  loading.value = false
}

async function handleRemove(postId: number) {
  try {
    await removeFavorite(postId)
    list.value = list.value.filter(i => i.post_id !== postId)
    uni.showToast({ title: "已取消收藏", icon: "none" })
  } catch { uni.showToast({ title: "操作失败", icon: "error" }) }
}

function goPost(id: number) {
  uni.navigateTo({ url: `/pages/post/detail?id=${id}` })
}

function formatTime(t: string) {
  if (!t) return ""
  const d = new Date(t)
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`
}
</script>

<style lang="scss" scoped>
.container { min-height: 100vh; background: #F7F4F0; padding: 20rpx; }
.loading-state { display: flex; justify-content: center; padding: 200rpx 0; }
.favorite-item {
  background: #fff; border-radius: 16rpx; padding: 28rpx; margin-bottom: 20rpx;
  .item-top { display: flex; justify-content: space-between; align-items: flex-start; }
  .item-title { font-size: 28rpx; color: #2C3A4B; flex: 1; line-height: 1.5; }
  .item-meta { font-size: 24rpx; color: #B8C2CE; margin-top: 12rpx; display: block; }
}
</style>
