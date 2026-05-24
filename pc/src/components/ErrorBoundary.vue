<template>
  <div v-if="hasError" class="error-boundary">
    <div class="eb-content">
      <el-icon :size="48" class="eb-icon"><WarningFilled /></el-icon>
      <h2>页面出错了</h2>
      <p>{{ errorMessage }}</p>
      <el-button type="primary" @click="handleRetry">重新加载</el-button>
    </div>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from "vue"
import { WarningFilled } from "@element-plus/icons-vue"

const hasError = ref(false)
const errorMessage = ref("")

onErrorCaptured((err) => {
  hasError.value = true
  errorMessage.value = (err as Error)?.message || "发生了未知错误"
  return false // prevent propagation
})

function handleRetry() {
  hasError.value = false
  errorMessage.value = ""
  window.location.reload()
}
</script>





<style scoped lang="scss">

@use "@/styles/variables.scss" as *;

.error-boundary {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}

.eb-content {
  text-align: center;
  padding: 40px;
}

.eb-icon { color: #e74c3c; margin-bottom: 16px; }

h2 { font-size: 20px; font-weight: 700; margin-bottom: 8px; color: $text-primary; }

p { font-size: 14px; color: $text-secondary; margin-bottom: 24px; }
</style>
