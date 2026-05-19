<template>
  <view class="min-h-screen bg-gray-50">
    <view class="p-4 flex justify-between items-center">
      <text class="text-lg font-bold">校园社团</text>
      <button class="text-sm text-white bg-gradient-to-r from-indigo-500 to-purple-600 px-4 py-1 rounded-full"
              @click="goCreate" v-if="!loading">创建圈子</button>
    </view>

    <view v-if="loading" class="flex justify-center py-20">
      <text class="text-gray-400">加载中...</text>
    </view>

    <view v-else-if="circles.length === 0" class="flex flex-col items-center py-20">
      <text class="text-5xl mb-4">🏛️</text>
      <text class="text-gray-400">还没有圈子，来创建一个吧</text>
    </view>

    <view v-else class="px-4 space-y-3">
      <view v-for="c in circles" :key="c.id"
            class="bg-white rounded-xl p-4 shadow-sm flex items-center gap-3 active:opacity-70"
            @click="goDetail(c.id)">
        <view class="w-14 h-14 rounded-full bg-gradient-to-br from-indigo-400 to-purple-500 flex items-center justify-center text-white text-xl shrink-0">
          {{ c.icon || c.name[0] }}
        </view>
        <view class="flex-1 min-w-0">
          <text class="font-bold text-base block truncate">{{ c.name }}</text>
          <text class="text-xs text-gray-400 block truncate mt-0.5">{{ c.description || '暂无简介' }}</text>
          <view class="flex gap-4 mt-1">
            <text class="text-xs text-gray-400">{{ c.member_count }} 成员</text>
            <text class="text-xs text-gray-400">{{ c.post_count }} 帖子</text>
          </view>
        </view>
        <view v-if="c.is_member"
              class="text-xs text-indigo-500 border border-indigo-500 rounded-full px-3 py-0.5 shrink-0">
          已加入
        </view>
      </view>
    </view>

    <view class="h-safe-area" />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCircles } from "@/api/circle"
import type { CircleData } from "@/api/circle"

const circles = ref<CircleData[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    circles.value = await listCircles()
  } catch (e: any) {
    uni.showToast({ title: "加载失败", icon: "none" })
  } finally {
    loading.value = false
  }
})

function goCreate() {
  uni.navigateTo({ url: "/pages/circle/create" })
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/circle/detail?id=${id}` })
}
</script>
