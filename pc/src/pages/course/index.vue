<template>
  <div class="course-page">
    <PageHeader title="课程评价" />

    <!-- Search -->
    <el-input v-model="query" placeholder="搜索课程、老师、学院..." size="large" clearable @keyup.enter="handleSearch" style="margin:16px 0 12px">
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>

    <!-- Quick Categories -->
    <div class="quick-cats">
      <button v-for="cat in quickCats" :key="cat" :class="['cat-btn', { active: query === cat }]" @click="query = cat; handleSearch()">{{ cat }}</button>
    </div>

    <!-- Popular Courses (pre-search hint) -->
    <div v-if="!searched" class="welcome-hint">
      <p>在上方搜索栏输入课程名称、老师姓名或学院来查找评价</p>
      <div class="popular-hint">
        <span class="hint-label">热门搜索：</span>
        <span v-for="hot in hotSearches" :key="hot" class="hot-tag" @click="query = hot; handleSearch()">{{ hot }}</span>
      </div>
    </div>

    <!-- Results -->
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
import PageHeader from "@/components/PageHeader.vue"

const query = ref("")
const courses = ref<CourseData[]>([])
const searchLoading = ref(false)
const searched = ref(false)

const quickCats = ["计算机", "数学", "英语", "经管", "物理", "化学"]
const hotSearches = ["高等数学", "大学英语", "C语言", "线性代数", "Python"]

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
.quick-cats {
  display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 16px;
  .cat-btn {
    padding: 4px 14px; border-radius: 16px; border: 1px solid $border-color; font-size: 13px; cursor: pointer;
    background: $bg-card; color: $text-secondary; transition: all 0.2s;
    &:hover { border-color: $primary; color: $primary; }
    &.active { background: $primary; color: #fff; border-color: $primary; }
  }
}
.welcome-hint {
  text-align: center; padding: 40px 0;
  p { color: $text-muted; font-size: 14px; margin-bottom: 20px; }
  .popular-hint { font-size: 13px; }
  .hint-label { color: $text-secondary; margin-right: 8px; }
  .hot-tag {
    display: inline-block; padding: 3px 12px; margin: 3px; border-radius: 12px;
    background: color-mix(in oklch, $primary 8%, transparent); color: $primary; cursor: pointer; transition: all 0.2s;
    &:hover { background: $primary; color: #fff; }
  }
}
.course-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 14px; }
.course-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px;
  h3 { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
  .meta { display: flex; gap: 10px; font-size: 13px; color: $text-secondary; margin-bottom: 8px; }
  .rating { margin-bottom: 8px; }
  .extra { display: flex; gap: 12px; font-size: 12px; color: $text-muted; }
:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
}
</style>
