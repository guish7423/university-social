<template>
  <div class="campus-calendar">
    <div class="page-header">
      <h1>校园日历</h1>
    </div>

    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!events.length" class="empty-state"><el-empty description="暂无校园活动" /></div>

    <template v-else>
      <div v-for="group in groupedEvents" :key="group.month" class="month-section stagger-item">
        <div class="month-header">
          <span class="month-label">{{ group.month }}</span>
          <span class="event-count">{{ group.events.length }} 个活动</span>
        </div>
        <div v-for="e in group.events" :key="e.id" class="event-card">
          <div class="event-date-badge">
            <span class="date-day">{{ extractDay(e.date) }}</span>
            <span class="date-month">{{ extractShortMonth(e.date) }}</span>
          </div>
          <div class="event-body">
            <h3 class="event-title">{{ e.title }}</h3>
            <p class="event-desc">{{ e.description }}</p>
            <div class="event-meta">
              <el-tag size="small">{{ e.location }}</el-tag>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { listCalendar } from "@/api/campus"
import type { CampusEvent } from "@/api/campus"

const events = ref<CampusEvent[]>([])
const loading = ref(true)

const groupedEvents = computed(() => {
  const groups: Record<string, CampusEvent[]> = {}
  const months = ["1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"]

  for (const e of events.value) {
    const d = new Date(e.date)
    const key = `${d.getFullYear()}年${months[d.getMonth()]}`
    if (!groups[key]) groups[key] = []
    groups[key].push(e)
  }

  return Object.entries(groups)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([month, events]) => ({ month, events }))
})

function extractDay(dateStr: string): string {
  return String(new Date(dateStr).getDate()).padStart(2, "0")
}

function extractShortMonth(dateStr: string): string {
  const months = ["1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"]
  return months[new Date(dateStr).getMonth()]
}

onMounted(async () => {
  try { events.value = await listCalendar() }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
@use "sass:color";

.campus-calendar {
  max-width: 720px;
}

.page-header {
  margin-bottom: 24px;
  h1 { font-size: 22px; font-weight: 700; }
}

.loading-wrap, .empty-state { padding: 60px 0; }

.month-section {
  margin-bottom: 32px;
}

.month-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid $border-light;

  .month-label {
    font-size: 18px;
    font-weight: 700;
    color: $text-primary;
  }

  .event-count {
    font-size: 13px;
    color: $text-muted;
  }
}

.event-card {
  display: flex;
  gap: 16px;
  padding: 16px;
  margin-bottom: 12px;
  background: $bg-card;
  border-radius: $radius-md;
  border: 1px solid $border-light;
  transition: all 0.25s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
    border-color: $brand-primary;
  }
}

.event-date-badge {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 56px;
  height: 56px;
  background: linear-gradient(135deg, color.change($brand-primary, $alpha: 0.15), color.change($brand-primary, $alpha: 0.08));
  border-radius: $radius-sm;
  flex-shrink: 0;

  .date-day {
    font-size: 22px;
    font-weight: 800;
    line-height: 1;
    color: $brand-primary;
  }

  .date-month {
    font-size: 11px;
    font-weight: 600;
    color: $text-muted;
    text-transform: uppercase;
    margin-top: 2px;
  }
}

.event-body {
  flex: 1;
  min-width: 0;
}

.event-title {
  font-size: 16px;
  font-weight: 600;
  color: $text-primary;
  margin: 0 0 6px;
  line-height: 1.4;
}

.event-desc {
  font-size: 14px;
  color: $text-secondary;
  margin: 0 0 10px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.event-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
</style>
