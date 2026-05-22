<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCalendar, type CalendarEvent } from "@/api/campus"

const events = ref<CalendarEvent[]>([])
const loading = ref(true)
const currentMonth = ref("")

const eventColors: Record<string, string> = {
  holiday: "#E74C3C", exam: "#E67E22", registration: "#3498DB",
  activity: "#2ECC71", other: "#95A5A6",
}
const eventLabels: Record<string, string> = {
  holiday: "假期", exam: "考试", registration: "报名", activity: "活动", other: "其他",
}

onMounted(async () => {
  try {
    events.value = await listCalendar()
    if (events.value.length) {
      const d = new Date(events.value[0].event_date)
      currentMonth.value = `${d.getFullYear()}年${d.getMonth() + 1}月`
    }
  } catch {}
  loading.value = false
})

function groupByMonth() {
  const groups: Record<string, CalendarEvent[]> = {}
  for (const e of events.value) {
    const m = e.event_date.slice(0, 7)
    if (!groups[m]) groups[m] = []
    groups[m].push(e)
  }
  return groups
}

function dayOfMonth(date: string) {
  return parseInt(date.slice(8, 10))
}

function monthLabel(m: string) {
  const d = new Date(m + "-01")
  return `${d.getFullYear()}年${d.getMonth() + 1}月`
}
</script>

<template>
  <view class="page">
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />

    </view>
    <template v-else-if="events.length">
      <view v-for="(evts, month) in groupByMonth()" :key="month" class="month-section">
        <text class="month-title">{{ monthLabel(month) }}</text>
        <view v-for="e in evts" :key="e.id" class="event-card">
          <view class="event-date-badge">
            <text class="event-day">{{ dayOfMonth(e.event_date) }}</text>
            <text class="event-month">{{ e.event_date.slice(5, 7) }}月</text>
          </view>
          <view class="event-info">
            <text class="event-title">{{ e.title }}</text>
            <text :style="{ color: eventColors[e.event_type] || '#95A5A6' }" class="event-tag">
              {{ eventLabels[e.event_type] || e.event_type }}
            </text>
          </view>
        </view>
      </view>
    </template>
    <view v-else class="empty-state">
      <u-icon name="calendar" size="120" color="#EAE6E0" />
      <text class="empty-text">暂无校历信息</text>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: #F7F4F0; padding: 24rpx; }
.loading-state, .empty-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 0; gap: 16rpx; }
.loading-text, .empty-text { font-size: 26rpx; color: #B8C2CE; }
.month-section { margin-bottom: 24rpx; }
.month-title {
  font-size: 28rpx; font-weight: 700; color: #1E2A3A;
  padding-left: 8rpx; margin-bottom: 12rpx; border-left: 6rpx solid #C67A6A;
}
.event-card {
  display: flex; gap: 16rpx; background: #fff; border-radius: 14rpx;
  padding: 20rpx; margin-bottom: 12rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
}
.event-date-badge {
  width: 80rpx; height: 80rpx; border-radius: 12rpx;
  background: #1E2A3A; display: flex; flex-direction: column;
  align-items: center; justify-content: center; flex-shrink: 0;
}
.event-day { font-size: 32rpx; font-weight: 700; color: #fff; line-height: 1.1; }
.event-month { font-size: 18rpx; color: rgba(255,255,255,0.6); }
.event-info { flex: 1; display: flex; flex-direction: column; gap: 8rpx; }
.event-title { font-size: 26rpx; font-weight: 600; color: #1E2A3A; }
.event-tag { font-size: 20rpx; padding: 2rpx 12rpx; border-radius: 6rpx; align-self: flex-start; font-weight: 500; }
</style>
