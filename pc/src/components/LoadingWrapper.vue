<template>
  <div v-if="loading" class="loading-wrap">
    <el-skeleton :rows="rows" animated :loading="true">
      <template v-slot:template>
        <div v-for="i in rows" :key="i" class="skeleton-row">
          <el-skeleton-item variant="text" :style="{ width: (80 - i * 8) + '%' }" />
        </div>
      </template>
    </el-skeleton>
  </div>
  <div v-else-if="!data && empty" class="empty-container">
    <slot name="empty">
      <el-empty :description="empty.description || '暂无内容'" :image-size="120">
        <template v-if="empty.action">
          <el-button type="primary" @click="empty.action.handler">
            {{ empty.action.label }}
          </el-button>
        </template>
      </el-empty>
    </slot>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
interface EmptyConfig {
  description?: string
  action?: { label: string; handler: () => void }
}

defineProps<{
  loading: boolean
  data?: any
  rows?: number
  empty?: EmptyConfig
}>()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.loading-wrap {
  padding: $space-6 0;
}

.skeleton-row {
  margin-bottom: $space-4;
}

.empty-container {
  padding: $space-12 0;
}
</style>
