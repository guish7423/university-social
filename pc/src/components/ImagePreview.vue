<template>
  <Teleport to="body">
    <Transition name="preview-fade">
      <div v-if="modelValue !== null" class="preview-overlay" @keydown="onKeydown" tabindex="-1" ref="overlayRef" @click="emit('update:modelValue', null)">
        <div v-if="imgLoading" class="preview-placeholder"><el-icon :size="48"><PictureFilled /></el-icon></div>
        <img
          v-show="!imgLoading"
          :src="images[modelValue]"
          class="preview-img"
          :class="{ zoomed: zoom === 2 }"
          :style="{ transform: `scale(${zoom})` }"
          @click.stop="toggleZoom"
          @load="imgLoading = false"
          @error="imgLoading = false"
        />
        <div class="preview-nav">
          <span v-if="modelValue! > 0" class="nav-btn prev" @click.stop="go(-1)">‹</span>
          <span v-if="modelValue! < images.length - 1" class="nav-btn next" @click.stop="go(1)">›</span>
        </div>
        <span class="preview-counter" :key="modelValue">{{ modelValue! + 1 }} / {{ images.length }}</span>
        <span class="preview-close" @click.stop="emit('update:modelValue', null)"><el-icon :size="24"><Close /></el-icon></span>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { PictureFilled, Close } from '@element-plus/icons-vue'

const props = defineProps<{ images: string[]; modelValue: number | null }>()
const emit = defineEmits<{ 'update:modelValue': [value: number | null] }>()

const zoom = ref(1)
const imgLoading = ref(true)
const overlayRef = ref<HTMLElement | null>(null)

function go(dir: number) {
  if (props.modelValue === null) return
  const next = props.modelValue + dir
  if (next >= 0 && next < props.images.length) {
    zoom.value = 1
    imgLoading.value = true
    emit('update:modelValue', next)
  }
}

function toggleZoom() {
  zoom.value = zoom.value === 1 ? 2 : 1
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowLeft') go(-1)
  else if (e.key === 'ArrowRight') go(1)
  else if (e.key === 'Escape') emit('update:modelValue', null)
}

watch(() => props.modelValue, (val) => {
  if (val !== null) {
    zoom.value = 1
    imgLoading.value = true
    nextTick(() => overlayRef.value?.focus())
  }
})
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
  outline: none;
}

.preview-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255,255,255,0.3);
}

.preview-img {
  max-width: 92vw;
  max-height: 92vh;
  object-fit: contain;
  border-radius: 6px;
  cursor: zoom-in;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.5);
  transition: transform 0.25s cubic-bezier(0.22,1,0.36,1);
  &.zoomed { cursor: zoom-out; }
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
  animation: counter-pop 0.2s ease-out;
}

.preview-close {
  position: fixed;
  top: 16px;
  right: 16px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255,255,255,0.7);
  background: rgba(255,255,255,0.1);
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  &:hover { background: rgba(255,255,255,0.25); color: #fff; }
}

.preview-fade-enter-active { transition: opacity 0.2s ease; }
.preview-fade-leave-active { transition: opacity 0.15s ease; }
.preview-fade-enter-from, .preview-fade-leave-to { opacity: 0; }

@keyframes counter-pop {
  from { opacity: 0; transform: translateX(-50%) translateY(8px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}
</style>
