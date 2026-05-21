<template>
  <u-popup :show="visible" mode="bottom" @close="close" round="16">
    <view class="report-popup safe-area-bottom">
      <view class="report-header">
        <text class="report-title">举报</text>
        <text class="report-subtitle">请选择举报原因</text>
      </view>

      <view class="report-reasons">
        <view
          v-for="reason in reasons"
          :key="reason"
          :class="['reason-item', selected === reason && 'is-selected']"
          @click="selected = reason"
        >
          <text class="reason-text">{{ reason }}</text>
          <u-icon
            :name="selected === reason ? 'checkbox-mark' : 'circle'"
            :color="selected === reason ? '#6b4ce6' : '#c0c4cc'"
            size="36"
          />
        </view>
      </view>

      <view class="report-actions">
        <u-button
          class="report-btn"
          :disabled="!selected"
          :loading="submitting"
          type="primary"
          shape="circle"
          @click="handleSubmit"
        >
          提交举报
        </u-button>
        <u-button class="cancel-btn" shape="circle" plain @click="close">
          取消
        </u-button>
      </view>
    </view>
  </u-popup>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { createReport, getReportReasons } from "@/api/report"

const props = defineProps<{
  visible: boolean
  contentType: string
  contentId: number
}>()

const emit = defineEmits<{
  (e: "close"): void
  (e: "success"): void
}>()

const selected = ref("")
const submitting = ref(false)
const reasons = getReportReasons()

function close() {
  selected.value = ""
  emit("close")
}

async function handleSubmit() {
  if (!selected.value || submitting.value) return
  submitting.value = true
  try {
    await createReport({
      content_type: props.contentType,
      content_id: props.contentId,
      reason: selected.value,
    })
    uni.showToast({ title: "举报已提交", icon: "success" })
    emit("success")
  } catch {
    uni.showToast({ title: "举报失败", icon: "none" })
  } finally {
    submitting.value = false
    close()
  }
}
</script>

<style scoped>
.report-popup {
  padding: 32rpx;
}

.report-header {
  margin-bottom: 32rpx;
  text-align: center;
}

.report-title {
  font-size: 36rpx;
  font-weight: 600;
  color: var(--ink, #333);
}

.report-subtitle {
  display: block;
  margin-top: 8rpx;
  font-size: 28rpx;
  color: var(--ink-muted, #999);
}

.report-reasons {
  margin-bottom: 32rpx;
}

.reason-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 16rpx;
  border-radius: 12rpx;
  margin-bottom: 8rpx;
  transition: background 0.2s;
}

.reason-item:active {
  background: var(--color-surface, #f5f5f5);
}

.reason-item.is-selected {
  background: rgba(107, 76, 230, 0.06);
}

.reason-text {
  font-size: 30rpx;
  color: var(--ink, #333);
}

.reason-item.is-selected .reason-text {
  color: var(--brand-primary, #6b4ce6);
  font-weight: 500;
}

.report-actions .report-btn {
  margin-bottom: 16rpx;
}

.report-actions .cancel-btn {
  border-color: var(--hairline, #e8e8e8);
  color: var(--ink-muted, #999);
}
</style>
