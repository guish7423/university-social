<template>
  <div class="course-page">
    <h1>课程评价</h1>
    <el-input v-model="query" placeholder="搜索课程、老师..." size="large" clearable @keyup.enter="handleSearch" style="margin:16px 0">
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>

    <div v-if="searchLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="searched && !courses.length" class="empty-state"><el-empty description="未找到课程" /></div>
    <div v-else-if="courses.length" class="course-grid">
      <div v-for="c in courses" class="course-card stagger-item" :key="c.id">
        <h3>{{ c.name }}</h3>
        <div class="meta">
          <span>{{ c.teacher }}</span>
          <span>{{ c.department }}</span>
        </div>
        <div class="rating">
          <el-rate v-model="c.rating" disabled :max="5" />
        </div>
        <div class="extra">
          <span>{{ c.credits }} 学分</span>
          <span>{{ c.enrolled }}/{{ c.capacity }} 人</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { searchCourses } from "@/api/course"
import type { CourseData } from "@/api/course"

const query = ref("")
const courses = ref<CourseData[]>([])
const searchLoading = ref(false)
const searched = ref(false)

async function handleSearch() {
  if (!query.value.trim()) return
  searched.value = true
  searchLoading.value = true
  try { courses.value = await searchCourses(query.value.trim()) }
  catch { /* handled by interceptor */ } finally { searchLoading.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.course-page { max-width: 780px; }
h1 { font-size: 22px; font-weight: 700; }

.course-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 14px; }

.course-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px;
  h3 { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
  .meta { display: flex; gap: 10px; font-size: 13px; color: $text-secondary; margin-bottom: 8px; }
  .rating { margin-bottom: 8px; }
  .extra { display: flex; gap: 12px; font-size: 12px; color: $text-muted; }
}
</style>
