<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listSensitive, addSensitive, removeSensitive } from "@/api/admin"

const words = ref<string[]>([])
const newWord = ref("")
const loading = ref(true)

onMounted(async () => {
  try { words.value = await listSensitive() } catch {}
  loading.value = false
})

async function add() {
  const w = newWord.value.trim()
  if (!w) return
  uni.showLoading({ title: "添加中" })
  try {
    await addSensitive(w)
    words.value.push(w)
    newWord.value = ""
    uni.showToast({ title: "已添加", icon: "success" })
  } catch {
    uni.showToast({ title: "添加失败", icon: "none" })
  }
  uni.hideLoading()
}

async function remove(word: string) {
  uni.showModal({
    title: "确认删除",
    content: `删除敏感词「${word}」？`,
    success: async (res) => {
      if (!res.confirm) return
      uni.showLoading({ title: "删除中" })
      try {
        await removeSensitive(word)
        words.value = words.value.filter(w => w !== word)
        uni.showToast({ title: "已删除", icon: "success" })
      } catch {
        uni.showToast({ title: "删除失败", icon: "none" })
      }
      uni.hideLoading()
    }
  })
}
</script>

<template>
  <view class="page">
    <view class="page-header">
      <text class="page-title">敏感词管理</text>
    </view>

    <view class="add-bar">
      <input v-model="newWord" class="input" placeholder="输入敏感词" @confirm="add" />
      <text class="add-btn" @click="add">添加</text>
    </view>

    <view class="list">
      <view v-for="w in words" :key="w" class="word-item">
        <text class="word-text">{{ w }}</text>
        <text class="del-btn" @click="remove(w)">删除</text>
      </view>
      <view v-if="words.length === 0" class="empty">暂无敏感词</view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 24rpx 16rpx; }
.page-header { margin-bottom: 24rpx; }
.page-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }

.add-bar {
  display: flex; gap: 16rpx; margin-bottom: 24rpx;
}
.input {
  flex: 1; background: #fff; border-radius: 14rpx; padding: 20rpx 24rpx;
  font-size: 26rpx; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
}
.add-btn {
  padding: 20rpx 32rpx; border-radius: 14rpx; background: #1E2A3A; color: #fff;
  font-size: 26rpx; font-weight: 500; cursor: pointer;
}

.list { background: #fff; border-radius: 14rpx; overflow: hidden; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); }
.word-item {
  display: flex; align-items: center; justify-content: space-between;
  padding: 24rpx; border-bottom: 1rpx solid #EAE6E0;
  &:last-child { border-bottom: none; }
}
.word-text { font-size: 26rpx; color: #1E2A3A; }
.del-btn { font-size: 24rpx; color: #E74C3C; cursor: pointer; padding: 8rpx; }
.empty { text-align: center; padding: 48rpx; color: #B8C2CE; }
</style>
