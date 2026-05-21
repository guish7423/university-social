<script setup lang="ts">
import { ref } from "vue"
import { createLostItem } from "@/api/found"

const form = ref({
  title: "",
  description: "",
  category: "lost",
  item_type: "",
  location: "",
  contact: "",
  images: [] as string[],
})
const submitting = ref(false)

function onUploaded(files: any[]) {
  form.value.images = files.map((f: any) => f.response?.data?.url || "").filter(Boolean)
}

async function submit() {
  if (!form.value.title || !form.value.category) return
  submitting.value = true
  try {
    await createLostItem(form.value)
    uni.showToast({ title: "发布成功", icon: "success" })
    setTimeout(() => uni.navigateBack(), 1200)
  } catch {
    uni.showToast({ title: "发布失败", icon: "none" })
  }
  submitting.value = false
}
</script>

<template>
  <view class="page">
    <view class="card">
      <view class="field">
        <text class="label">◈ 类型</text>
        <view class="type-picker">
          <view class="type-option" :class="{ active: form.category === 'lost' }" @click="form.category = 'lost'">
            <text>◉ 寻物</text>
          </view>
          <view class="type-option" :class="{ active: form.category === 'found' }" @click="form.category = 'found'">
            <text>◎ 招领</text>
          </view>
        </view>
      </view>
      <u-input placeholder="物品名称" v-model="form.title" border="none" inputClass="field-input" />
      <u-input type="textarea" placeholder="详细描述（时间、特征等）" v-model="form.description" border="none" height="100" />
      <u-input placeholder="物品类型（如：校园卡、耳机）" v-model="form.item_type" border="none" inputClass="field-input" />
      <u-input placeholder="地点（如：图书馆3楼）" v-model="form.location" border="none" inputClass="field-input" />
      <u-input placeholder="联系方式" v-model="form.contact" border="none" inputClass="field-input" />

      <view class="field">
        <text class="label">▣ 图片</text>
        <u-upload :action="'/api/v1/upload'" :auto-upload="true" max-count="6" @on-uploaded="onUploaded"></u-upload>
      </view>

      <u-button type="primary" shape="circle" :disabled="!form.title || !form.category" :loading="submitting" @click="submit">发布</u-button>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 24rpx; }
.card { background: var(--color-surface, #fff); border-radius: 22rpx; padding: 32rpx; }
.field { margin-bottom: 24rpx; }
.label { font-size: 28rpx; font-weight: 600; color: var(--ink, #1E2A3A); display: block; margin-bottom: 12rpx; }
.type-picker { display: flex; flex-direction: row; gap: 20rpx; }
.type-option {
  flex: 1; padding: 24rpx; border-radius: 16rpx; text-align: center;
  background: var(--color-surface-1, #F0EDE8); font-size: 28rpx; color: var(--ink-muted, #5C6B7E);
  &.active {
    background: linear-gradient(135deg, #C67A6A, #1E2A3A);
    color: #fff;
  }
}
:deep(.field-input) { font-size: 28rpx; padding: 20rpx 0; border-bottom: 2rpx solid var(--hairline-light, #EAE6E0); margin-bottom: 16rpx; }
</style>
