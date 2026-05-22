<template>
  <div class="campus-calendar">
    <div class="page-header">
      <h1>校园日历</h1>
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!events.length" class="empty-state"><el-empty description="暂无校园活动" /></div>
    <el-timeline v-else>
      <el-timeline-item v-for="e in events" class="stagger-item" :key="e.id" :timestamp="formatDate(e.date)" placement="top">
        <el-card shadow="never">
          <h3>{{ e.title }}</h3>
          <p>{{ e.description }}</p>
          <el-tag>{{ e.location }}</el-tag>
        </el-card>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCalendar } from "@/api/campus"
import type { CampusEvent } from "@/api/campus"
import { useTimeFormat } from "@/composables/useTimeFormat"

const { formatDate } = useTimeFormat()
const events = ref<CampusEvent[]>([])
const loading = ref(true)

onMounted(async () => {
  try { events.value = await listCalendar() }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
.campus-calendar { max-width: 680px; }
.page-header { margin-bottom: 20px; h1 { font-size: 22px; font-weight: 700; } }
.loading-wrap, .empty-state { padding: 40px 0; }
:deep(.el-timeline-item__timestamp) { color: $text-muted; }
</style>
