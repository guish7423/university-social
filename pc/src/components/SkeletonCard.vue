<template>
  <div class="skeleton-card" :class="variant">
    <div v-if="variant === 'post-card'" class="sk-post">
      <div class="sk-post-header">
        <el-skeleton-item variant="circle" class="sk-avatar" />
        <div class="sk-author">
          <el-skeleton-item variant="text" class="sk-name" />
          <el-skeleton-item variant="text" class="sk-meta" />
        </div>
      </div>
      <el-skeleton-item variant="text" class="sk-line w-100" />
      <el-skeleton-item variant="text" class="sk-line w-80" />
      <el-skeleton-item variant="text" class="sk-line w-60" />
      <div class="sk-images" v-if="hasImages">
        <el-skeleton-item variant="rect" class="sk-img" />
        <el-skeleton-item variant="rect" class="sk-img" />
        <el-skeleton-item variant="rect" class="sk-img" />
      </div>
      <div class="sk-footer">
        <el-skeleton-item variant="text" class="sk-action" />
        <el-skeleton-item variant="text" class="sk-action" />
      </div>
    </div>

    <div v-else-if="variant === 'circle-card'" class="sk-circle">
      <el-skeleton-item variant="circle" class="sk-icon" />
      <div class="sk-info">
        <el-skeleton-item variant="text" class="sk-name" />
        <el-skeleton-item variant="text" class="sk-count" />
        <el-skeleton-item variant="text" class="sk-desc" />
      </div>
    </div>

    <div v-else-if="variant === 'user-card'" class="sk-user">
      <el-skeleton-item variant="circle" class="sk-avatar-lg" />
      <div class="sk-user-info">
        <el-skeleton-item variant="text" class="sk-name" />
        <el-skeleton-item variant="text" class="sk-bio" />
      </div>
    </div>

    <div v-else-if="variant === 'detail'" class="sk-detail">
      <el-skeleton-item variant="text" class="sk-detail-title" />
      <el-skeleton-item variant="text" class="sk-line w-100" />
      <el-skeleton-item variant="text" class="sk-line w-100" />
      <el-skeleton-item variant="text" class="sk-line w-90" />
      <el-skeleton-item variant="text" class="sk-line w-75" />
      <el-skeleton-item variant="rect" class="sk-detail-img" />
      <el-skeleton-item variant="text" class="sk-line w-100" />
      <el-skeleton-item variant="text" class="sk-line w-60" />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  variant?: 'post-card' | 'circle-card' | 'user-card' | 'detail'
  hasImages?: boolean
}>(), {
  variant: 'post-card',
  hasImages: false,
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.skeleton-card {
  background: $bg-card;
  border: 1px solid $border-default;
  border-radius: $radius-md;
  padding: $space-5;
  margin-bottom: $space-4;
}

.el-skeleton-item {
  background: linear-gradient(90deg, $bg-hover 25%, $bg-raised 37%, $bg-hover 63%);
  background-size: 400% 100%;
  animation: el-skeleton-loading 1.4s ease infinite;
  border-radius: 4px;
}

@keyframes el-skeleton-loading {
  0% { background-position: 100% 50%; }
  100% { background-position: 0 50%; }
}

// ── Shared ──
.sk-avatar { width: 36px; height: 36px; flex-shrink: 0; }
.sk-avatar-lg { width: 48px; height: 48px; flex-shrink: 0; }
.sk-name { height: 14px; width: 120px; display: block; margin-bottom: 8px; }
.sk-meta { height: 12px; width: 80px; display: block; }
.sk-line { height: 13px; display: block; margin-bottom: 10px; }
.w-100 { width: 100%; }
.w-90 { width: 90%; }
.w-80 { width: 80%; }
.w-75 { width: 75%; }
.w-70 { width: 70%; }
.w-60 { width: 60%; }

// ── Post Card ──
.sk-post-header {
  display: flex;
  align-items: center;
  gap: $space-3;
  margin-bottom: $space-4;
}
.sk-author {
  flex: 1;
}
.sk-images {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 4px;
  margin: $space-3 0;
}
.sk-img {
  aspect-ratio: 1;
  border-radius: $radius-sm;
}
.sk-footer {
  display: flex;
  gap: $space-6;
  padding-top: $space-3;
  border-top: 1px solid $border-default;
  margin-top: $space-2;
}
.sk-action {
  height: 14px;
  width: 60px;
}

// ── Circle Card ──
.sk-circle {
  display: flex;
  gap: $space-4;
  align-items: flex-start;
}
.sk-icon { width: 48px; height: 48px; flex-shrink: 0; }
.sk-info { flex: 1; }
.sk-count { height: 12px; width: 100px; display: block; margin-bottom: 8px; }
.sk-desc { height: 12px; width: 70%; display: block; }

// ── User Card ──
.sk-user {
  display: flex;
  gap: $space-4;
  align-items: center;
}
.sk-user-info { flex: 1; }
.sk-bio { height: 12px; width: 60%; display: block; }

// ── Detail ──
.sk-detail-title {
  height: 24px;
  width: 65%;
  display: block;
  margin-bottom: $space-4;
}
.sk-detail-img {
  width: 100%;
  height: 200px;
  display: block;
  border-radius: $radius-md;
  margin: $space-3 0;
}
</style>
