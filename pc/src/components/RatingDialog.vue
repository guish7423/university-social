<template>
  <el-dialog :model-value="visible" title="写评价" width="50vw" :top="'15vh'" @update:model-value="$emit('close')" destroy-on-close>
    <div class="rating-form">
      <div class="form-section">
        <label class="form-label">评分</label>
        <el-rate v-model="score" :max="5" show-score allow-half size="large" />
      </div>
      <div class="form-section">
        <label class="form-label">评价内容</label>
        <el-input
          v-model="comment"
          type="textarea"
          :rows="4"
          maxlength="500"
          show-word-limit
          placeholder="说说你的看法..."
        />
      </div>
    </div>
    <template #footer>
      <el-button @click="$emit('close')">取消</el-button>
      <el-button type="primary" :loading="submitting" :disabled="score === 0" @click="handleSubmit">
        提交评价
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue"

defineProps<{
  visible: boolean
  title?: string
}>()

const emit = defineEmits<{
  (e: "close"): void
  (e: "submit", data: { score: number; comment: string }): void
}>()

const score = ref(0)
const comment = ref("")
const submitting = ref(false)

async function handleSubmit() {
  if (score.value === 0) return
  submitting.value = true
  emit("submit", { score: score.value, comment: comment.value })
  // Reset
  score.value = 0
  comment.value = ""
  submitting.value = false
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.rating-form { padding: 8px 0; }
.form-section { margin-bottom: 20px;
  &:last-child { margin-bottom: 0; }
}
.form-label {
  display: block; font-size: 14px; font-weight: 600; margin-bottom: 10px;
}
</style>
