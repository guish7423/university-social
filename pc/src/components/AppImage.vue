<template>
  <div class="app-image" :class="{ loaded, error }" :style="wrapperStyle">
    <div v-if="!loaded && !error" class="app-image-placeholder">
      <el-icon :size="24"><PictureFilled /></el-icon>
    </div>
    <img
      v-show="loaded"
      :src="src"
      :alt="alt"
      :class="imgClass"
      loading="lazy"
      @load="onLoad"
      @error="onError"
    />
    <div v-if="error" class="app-image-error">
      <el-icon :size="20"><PictureFilled /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { PictureFilled } from '@element-plus/icons-vue'

const props = withDefaults(defineProps<{
  src: string
  alt?: string
  imgClass?: string
  aspectRatio?: string
}>(), {
  alt: '',
  imgClass: '',
  aspectRatio: '',
})

const emit = defineEmits<{
  load: []
}>()

const loaded = ref(false)
const error = ref(false)

const wrapperStyle = computed(() => {
  if (props.aspectRatio) {
    return { aspectRatio: props.aspectRatio }
  }
  return {}
})

function onLoad() {
  loaded.value = true
  emit('load')
}

function onError() {
  error.value = true
  loaded.value = true
}
</script>

<style scoped>
.app-image {
  position: relative;
  overflow: hidden;
  border-radius: 8px;
  background: var(--bg-hover, oklch(0.23 0.015 30));
  min-height: 60px;
}

.app-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: opacity 0.3s ease;
}

.app-image-placeholder,
.app-image-error {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted, oklch(0.45 0.01 30));
  background: var(--bg-hover, oklch(0.23 0.015 30));
}
</style>
