<template>
  <div class="circle-card" @click="$router.push('/circles/' + circle.id)">
    <div class="circle-avatar">
      <el-avatar :size="48" :src="circle.avatar" />
    </div>
    <div class="circle-info">
      <div class="circle-name">{{ circle.name }}</div>
      <div class="circle-desc">{{ circle.description }}</div>
      <div class="circle-meta">
        <span>{{ circle.member_count }} 人</span>
        <span>· {{ circle.post_count }} 帖</span>
      </div>
    </div>
    <div class="circle-action">
      <el-button
        v-if="circle.is_joined"
        size="small"
        plain
        @click.stop="emit('leave', circle)"
      >
        已加入
      </el-button>
      <el-button
        v-else
        size="small"
        type="primary"
        @click.stop="emit('join', circle)"
      >
        加入
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface CircleCardData {
  id: number
  name: string
  description: string
  avatar?: string
  member_count: number
  post_count: number
  is_joined: boolean
}

defineProps<{
  circle: CircleCardData
}>()

const emit = defineEmits<{
  join: [circle: CircleCardData]
  leave: [circle: CircleCardData]
}>()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.circle-card {
  display: flex;
  align-items: center;
  gap: $space-4;
  background: $bg-card;
  border: 1px solid $border-default;
  border-radius: $radius-md;
  padding: $space-4;
  cursor: pointer;
  transition: all $duration-fast $ease-out;

  &:hover {
    border-color: rgba($brand-primary-hex, 0.4);
    background: $bg-hover;
    transform: translateY(-1px);
  }
}

.circle-avatar {
  flex-shrink: 0;
}

.circle-info {
  flex: 1;
  min-width: 0;
}

.circle-name {
  font-size: $text-sm;
  font-weight: 600;
  margin-bottom: 2px;
}

.circle-desc {
  font-size: $text-xs;
  color: $text-secondary;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: $space-1;
}

.circle-meta {
  font-size: $text-xs;
  color: $text-muted;
}

.circle-action {
  flex-shrink: 0;
}
</style>
