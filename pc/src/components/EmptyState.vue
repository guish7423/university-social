<template>
  <div class="empty-state" :class="type">
    <div class="empty-illustration" :style="{ width: imageSize + 'px', height: imageSize + 'px' }">
      <!-- Empty type -->
      <svg v-if="type === 'empty'" viewBox="0 0 120 120" fill="none">
        <rect x="20" y="30" width="80" height="60" rx="8" stroke="#C67A6A" stroke-width="2" fill="rgba(198,122,106,0.08)"/>
        <rect x="35" y="45" width="50" height="4" rx="2" fill="#C67A6A" opacity="0.3"/>
        <rect x="35" y="55" width="35" height="4" rx="2" fill="#C67A6A" opacity="0.2"/>
        <rect x="35" y="65" width="40" height="4" rx="2" fill="#C67A6A" opacity="0.15"/>
        <circle cx="48" cy="38" r="16" stroke="#C67A6A" stroke-width="2" fill="rgba(198,122,106,0.1)"/>
        <circle cx="48" cy="38" r="4" fill="#C67A6A" opacity="0.4"/>
        <line x1="60" y1="44" x2="68" y2="50" stroke="#C67A6A" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <!-- Search type -->
      <svg v-else-if="type === 'search'" viewBox="0 0 120 120" fill="none">
        <circle cx="50" cy="50" r="24" stroke="#C67A6A" stroke-width="2.5" fill="rgba(198,122,106,0.08)"/>
        <circle cx="50" cy="50" r="8" fill="#C67A6A" opacity="0.3"/>
        <line x1="68" y1="68" x2="82" y2="82" stroke="#C67A6A" stroke-width="3" stroke-linecap="round"/>
      </svg>
      <!-- Error type -->
      <svg v-else-if="type === 'error'" viewBox="0 0 120 120" fill="none">
        <circle cx="60" cy="55" r="30" stroke="#E06C75" stroke-width="2" fill="rgba(224,108,117,0.08)"/>
        <line x1="46" y1="43" x2="74" y2="67" stroke="#E06C75" stroke-width="2.5" stroke-linecap="round"/>
        <line x1="74" y1="43" x2="46" y2="67" stroke="#E06C75" stroke-width="2.5" stroke-linecap="round"/>
      </svg>
      <!-- Message type -->
      <svg v-else viewBox="0 0 120 120" fill="none">
        <rect x="20" y="28" width="80" height="56" rx="10" stroke="#C67A6A" stroke-width="2" fill="rgba(198,122,106,0.08)"/>
        <path d="M40 65 L55 52 L70 65" stroke="#C67A6A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <circle cx="48" cy="44" r="4" fill="#C67A6A" opacity="0.4"/>
        <circle cx="60" cy="44" r="4" fill="#C67A6A" opacity="0.3"/>
        <circle cx="72" cy="44" r="4" fill="#C67A6A" opacity="0.2"/>
      </svg>
    </div>
    <p class="empty-title">{{ title }}</p>
    <p v-if="description" class="empty-desc">{{ description }}</p>
    <div v-if="$slots.action" class="empty-action">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  type?: 'empty' | 'search' | 'error' | 'message'
  title?: string
  description?: string
  imageSize?: number
}>(), {
  type: 'empty',
  title: '',
  description: '',
  imageSize: 120,
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;

  .empty-illustration {
    margin-bottom: 20px;
    opacity: 0.8;
    svg { width: 100%; height: 100%; }
  }

  .empty-title {
    font-size: $text-lg;
    color: $text-secondary;
    font-weight: 500;
    margin-bottom: 4px;
  }

  .empty-desc {
    font-size: $text-sm;
    color: $text-muted;
    max-width: 280px;
    line-height: 1.6;
  }

  .empty-action {
    margin-top: 20px;
  }

  &.error {
    .empty-illustration { opacity: 1; }
    .empty-title { color: $semantic-danger; }
  }
}
</style>
