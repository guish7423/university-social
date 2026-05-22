<template>
  <div class="image-uploader">
    <div class="upload-area" @click="triggerUpload" @dragover.prevent @drop.prevent="handleDrop">
      <el-icon :size="32"><Plus /></el-icon>
      <span>点击或拖拽上传图片</span>
    </div>
    <input ref="fileInput" type="file" multiple accept="image/*" style="display:none" @change="handleFileChange" />
    <div v-if="modelValue.length" class="preview-list">
      <div v-for="(url, i) in modelValue" :key="i" class="preview-item">
        <img :src="url" />
        <el-icon class="remove" @click="removeImage(i)"><Close /></el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { ElMessage } from "element-plus"
import { useImageUpload } from "@/composables/useImageUpload"

const props = defineProps<{ modelValue: string[] }>()
const emit = defineEmits<{ (e: "update:modelValue", val: string[]): void }>()
const fileInput = ref<HTMLInputElement>()
const { uploading, uploadImages } = useImageUpload()

async function handleFiles(files: FileList) {
  const valid = Array.from(files).filter(f => f.size < 10 * 1024 * 1024)
  if (valid.length !== files.length) ElMessage.warning("部分文件超过10MB限制")
  if (!valid.length) return
  try {
    const urls = await uploadImages(valid)
    emit("update:modelValue", [...props.modelValue, ...urls])
  } catch { /* handled by interceptor */ }
}

function handleFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files?.length) handleFiles(input.files)
  input.value = ""
}

async function handleDrop(e: DragEvent) {
  if (e.dataTransfer?.files?.length) await handleFiles(e.dataTransfer.files)
}

function triggerUpload() { fileInput.value?.click() }
function removeImage(i: number) {
  emit("update:modelValue", props.modelValue.filter((_, idx) => idx !== i))
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
.upload-area {
  border: 2px dashed $border-color; border-radius: $radius-md; padding: 32px;
  text-align: center; cursor: pointer; color: $text-muted; transition: all 0.2s;
  &:hover { border-color: $primary; color: $primary; }
}
.preview-list { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 12px; }
.preview-item {
  position: relative; width: 100px; height: 100px;
  img { width: 100%; height: 100%; object-fit: cover; border-radius: $radius-sm; }
}
.remove {
  position: absolute; top: -6px; right: -6px; background: #e74c3c; color: #fff;
  border-radius: 50%; padding: 2px; cursor: pointer; font-size: 14px;
}
</style>
