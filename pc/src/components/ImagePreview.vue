<template>
  <Teleport to="body">
    <Transition name="preview-fade">
      <div v-if="modelValue !== null" class="preview-overlay" @click="emit('update:modelValue', null)">
        <img :src="images[modelValue]" class="preview-img" @click.stop />
        <div class="preview-nav">
          <span v-if="modelValue! > 0" class="nav-btn prev" @click.stop="emit('update:modelValue', modelValue! - 1)">‹</span>
          <span v-if="modelValue! < images.length - 1" class="nav-btn next" @click.stop="emit('update:modelValue', modelValue! + 1)">›</span>
        </div>
        <span class="preview-counter">{{ modelValue! + 1 }} / {{ images.length }}</span>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
defineProps<{ images: string[]; modelValue: number | null }>()
const emit = defineEmits<{ 'update:modelValue': [value: number | null] }>()
</script>

<style scoped lang="scss">
.preview-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.92);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.preview-img {
  max-width: 92vw;
  max-height: 92vh;
  object-fit: contain;
  border-radius: 6px;
  cursor: default;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.5);
}

.preview-nav {
  .nav-btn {
    position: fixed;
    top: 50%;
    transform: translateY(-50%);
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 32px;
    color: #fff;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.2s;
    user-select: none;
    &:hover { background: rgba(255, 255, 255, 0.25); transform: translateY(-50%) scale(1.1); }
    &.prev { left: 24px; }
    &.next { right: 24px; }
  }
}

.preview-counter {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  background: rgba(0, 0, 0, 0.5);
  padding: 4px 14px;
  border-radius: 999px;
}
</style>
