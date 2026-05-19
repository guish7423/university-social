<template>
  <view class="min-h-screen bg-gray-50 p-4">
    <view class="bg-white rounded-xl p-5 shadow-sm">
      <text class="text-lg font-bold block mb-6">创建圈子</text>

      <view class="mb-4">
        <text class="text-sm text-gray-500 mb-1 block">圈子名称 *</text>
        <input v-model="name" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm"
               placeholder="给圈子起个名字" maxlength="32" />
      </view>

      <view class="mb-4">
        <text class="text-sm text-gray-500 mb-1 block">圈子简介</text>
        <textarea v-model="description" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm h-24"
                  placeholder="介绍一下这个圈子" maxlength="200" />
      </view>

      <view class="mb-4">
        <text class="text-sm text-gray-500 mb-1 block">圈子图标</text>
        <input v-model="icon" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm"
               placeholder="输入 emoji 或文字（可选）" maxlength="4" />
      </view>

      <view class="flex gap-3 mt-6">
        <button class="flex-1 py-2.5 rounded-lg text-sm border border-gray-300 text-gray-600"
                @click="uni.navigateBack()">取消</button>
        <button class="flex-1 py-2.5 rounded-lg text-sm text-white bg-gradient-to-r from-indigo-500 to-purple-600"
                :disabled="!name" :class="name ? '' : 'opacity-50'" @click="handleCreate">创建</button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { createCircle } from "@/api/circle"

const name = ref("")
const description = ref("")
const icon = ref("")

async function handleCreate() {
  if (!name.value.trim()) {
    uni.showToast({ title: "请输入圈子名称", icon: "none" })
    return
  }
  try {
    const res = await createCircle({
      name: name.value.trim(),
      description: description.value.trim(),
      icon: icon.value.trim(),
    })
    uni.showToast({ title: "创建成功", icon: "success" })
    setTimeout(() => uni.redirectTo({ url: `/pages/circle/detail?id=${res.id}` }), 500)
  } catch {
    uni.showToast({ title: "创建失败", icon: "none" })
  }
}
</script>
