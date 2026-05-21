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

<template>
  <view class="page">
    <view class="card">
      <text class="card-title">创建圈子</text>

      <view class="field">
        <text class="label">圈子名称</text>
        <input v-model="name" class="input" placeholder="给圈子起个名字" maxlength="32" />
      </view>

      <view class="field">
        <text class="label">圈子简介</text>
        <textarea v-model="description" class="textarea" placeholder="介绍一下这个圈子" maxlength="200" />
      </view>

      <view class="field">
        <text class="label">圈子图标</text>
        <input v-model="icon" class="input" placeholder="输入 emoji 或文字（可选）" maxlength="4" />
      </view>

      <view class="actions">
        <view class="btn btn-cancel" @click="uni.navigateBack()">取消</view>
        <view class="btn btn-primary" :class="{ disabled: !name }" @click="handleCreate">创建</view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--color-canvas, #F7F4F0);
  padding: 40rpx 32rpx;
}
.card {
  background: var(--color-surface, #fff);
  border-radius: 20rpx;
  padding: 40rpx;
  box-shadow: 0 4rpx 20rpx rgba(30,42,58,0.04);
}
.card-title {
  font-size: 36rpx;
  font-weight: 600;
  color: var(--ink, #1E2A3A);
  margin-bottom: 32rpx;
}
.field { margin-bottom: 28rpx; }
.label {
  font-size: 26rpx;
  color: var(--ink-muted, #5C6B7E);
  display: block;
  margin-bottom: 12rpx;
  font-weight: 500;
}
.input, .textarea {
  width: 100%;
  border: 1.5rpx solid var(--hairline, #E0DBD4);
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  color: var(--ink, #1E2A3A);
  background: var(--color-surface-1, #F0EDE8);
  box-sizing: border-box;
}
.input:focus, .textarea:focus {
  border-color: var(--brand-primary, #C67A6A);
}
.textarea { height: 200rpx; }
.actions { display: flex; gap: 20rpx; margin-top: 40rpx; }
.btn {
  flex: 1;
  padding: 24rpx;
  border-radius: 14rpx;
  text-align: center;
  font-size: 28rpx;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  &:active { transform: scale(0.97); }
}
.btn-cancel {
  background: var(--color-surface-1, #F0EDE8);
  color: var(--ink-muted, #5C6B7E);
}
.btn-primary {
  background: linear-gradient(135deg, #C67A6A, #1E2A3A);
  color: #fff;
  &.disabled { opacity: 0.5; }
}
</style>
