<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listReports, resolveReport, dismissReport } from "@/api/admin"

const reports = ref<any[]>([])
const loading = ref(true)
const offset = ref(0)
const hasMore = ref(true)
const statusFilter = ref("")

async function load() {
  try {
    const data = await listReports(statusFilter.value, offset.value)
    if (data.length < 20) hasMore.value = false
    reports.value = offset.value === 0 ? data : [...reports.value, ...data]
  } catch {}
  loading.value = false
}

onMounted(load)

function filterBy(s: string) {
  statusFilter.value = s
  offset.value = 0
  hasMore.value = true
  load()
}

async function doAction(id: number, action: "resolve" | "dismiss") {
  uni.showLoading({ title: "操作中" })
  try {
    if (action === "resolve") {
      await resolveReport(id)
      uni.showToast({ title: "已处理", icon: "success" })
    } else {
      await dismissReport(id)
      uni.showToast({ title: "已驳回", icon: "success" })
    }
    reports.value = reports.value.filter(r => r.id !== id)
  } catch {
    uni.showToast({ title: "操作失败", icon: "none" })
  }
  uni.hideLoading()
}
</script>

<template>
  <view class="page">
    <view class="page-header">
      <text class="page-title">举报管理</text>
    </view>

    <view class="tabs">
      <view :class="{ tab: true, active: statusFilter === '' }" @click="filterBy('')">全部</view>
      <view :class="{ tab: true, active: statusFilter === 'pending' }" @click="filterBy('pending')">待处理</view>
      <view :class="{ tab: true, active: statusFilter === 'resolved' }" @click="filterBy('resolved')">已处理</view>
      <view :class="{ tab: true, active: statusFilter === 'dismissed' }" @click="filterBy('dismissed')">已驳回</view>
    </view>

    <view class="list">
      <view v-for="r in reports" :key="r.id" class="report-item">
        <view class="report-header">
          <text class="report-type">{{ r.target_type }}</text>
          <text class="report-status" :class="r.status">{{ r.status }}</text>
        </view>
        <text class="report-reason">{{ r.reason }}</text>
        <view class="report-footer">
          <text class="reporter">举报人: {{ r.reporter_nickname || '#' + r.reporter_id }}</text>
          <text class="report-time">{{ r.created_at?.slice(0, 16) }}</text>
        </view>
        <view v-if="r.status === 'pending'" class="report-actions">
          <text class="action-btn resolve" @click="doAction(r.id, 'resolve')">✔ 处理</text>
          <text class="action-btn dismiss" @click="doAction(r.id, 'dismiss')">✘ 驳回</text>
        </view>
      </view>
      <view v-if="reports.length === 0" class="empty">暂无举报</view>
    </view>

    <view v-if="hasMore" class="load-more" @click="offset += 20; load()">
      加载更多
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding: 24rpx 16rpx; }
.page-header { margin-bottom: 16rpx; }
.page-title { font-size: 28rpx; font-weight: 600; color: #1E2A3A; }

.tabs { display: flex; gap: 8rpx; margin-bottom: 24rpx; }
.tab {
  padding: 8rpx 20rpx; border-radius: 20rpx; font-size: 24rpx;
  background: #fff; color: #8E9BAB; cursor: pointer;
}
.tab.active { background: #1E2A3A; color: #fff; }

.list { background: #fff; border-radius: 14rpx; overflow: hidden; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); }
.report-item { padding: 24rpx; border-bottom: 1rpx solid #EAE6E0; }
.report-header { display: flex; align-items: center; gap: 12rpx; margin-bottom: 8rpx; }
.report-type { font-size: 22rpx; color: #C67A6A; font-weight: 500; }
.report-status { font-size: 20rpx; padding: 2rpx 10rpx; border-radius: 4rpx; }
.report-status.pending { background: #FFF3E0; color: #E65100; }
.report-status.resolved { background: #E8F5E9; color: #388E3C; }
.report-status.dismissed { background: #F0EDE8; color: #8E9BAB; }
.report-reason { font-size: 24rpx; color: #333; line-height: 1.5; display: block; margin-bottom: 8rpx; }
.report-footer { display: flex; justify-content: space-between; font-size: 20rpx; color: #B8C2CE; margin-bottom: 12rpx; }
.report-actions { display: flex; gap: 16rpx; }
.action-btn { font-size: 24rpx; cursor: pointer; padding: 8rpx 16rpx; border-radius: 8rpx; }
.action-btn.resolve { background: #E8F5E9; color: #388E3C; }
.action-btn.dismiss { background: #F0EDE8; color: #8E9BAB; }
.empty { text-align: center; padding: 48rpx; color: #B8C2CE; }
.load-more { text-align: center; padding: 32rpx; color: #C67A6A; font-size: 26rpx; cursor: pointer; }
</style>
