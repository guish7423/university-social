<template>
  <div class="page-header">
    <div class="page-header-left">
      <h1 class="page-title">{{ title }}</h1>
      <p v-if="description" class="page-desc">{{ description }}</p>
    </div>
    <div class="page-header-right">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: $space-6;
  gap: $space-4;
}

.page-title {
  font-size: $text-2xl;
  font-weight: 700;
  line-height: 1.3;
}

.page-desc {
  font-size: $text-sm;
  color: $text-secondary;
  margin-top: $space-1;
}
</style>
