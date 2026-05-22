<template>
  <div v-if="breadcrumbs.length > 0" class="breadcrumb-bar">
    <el-breadcrumb separator="›">
      <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item
        v-for="crumb in breadcrumbs"
        :key="crumb.label"
        :to="crumb.path ? { path: crumb.path } : undefined"
      >
        {{ crumb.label }}
      </el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  breadcrumbs: { label: string; path?: string }[]
}>()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.breadcrumb-bar {
  padding: $space-3 $space-6;
  border-bottom: 1px solid $border-default;
  flex-shrink: 0;

  :deep(.el-breadcrumb__inner) { color: $text-muted !important; font-size: 13px; }
  :deep(.el-breadcrumb__inner.is-link) { color: $text-secondary !important; }
  :deep(.el-breadcrumb__inner.is-link:hover) { color: $brand-primary !important; }
  :deep(.el-breadcrumb__separator) { color: $text-muted; font-weight: 600; }
}
</style>
