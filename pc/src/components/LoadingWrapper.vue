<template>
  <div v-if="loading" class="loading-wrap">
    <SkeletonCard
      v-for="i in (rows || 3)" :key="i"
      :variant="skeletonVariant"
      :has-images="skeletonImages"
    />
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
import SkeletonCard from "./SkeletonCard.vue"

interface EmptyConfig {
  description?: string
  action?: { label: string; handler: () => void }
}

withDefaults(defineProps<{
  loading: boolean
  data?: any
  rows?: number
  skeletonVariant?: 'post-card' | 'circle-card' | 'user-card' | 'detail'
  skeletonImages?: boolean
  empty?: EmptyConfig
}>(), {
  rows: 3,
  skeletonVariant: 'post-card',
  skeletonImages: false,
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.loading-wrap {
  padding: $space-2 0;
}

.empty-container {
  padding: $space-12 0;
}
</style>
