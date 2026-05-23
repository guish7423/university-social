<template>
  <div class="create-page">
    <h1>发布动态</h1>
    <el-input
      v-model="content"
      type="textarea" :rows="6"
      placeholder="分享你的想法..."
      maxlength="2000" show-word-limit
    />
    <ImageUploader v-model="images" />
    <div class="form-actions">
      <el-button @click="$router.back()">取消</el-button>
      <el-button type="primary" :loading="submitting" :disabled="!content.trim()" @click="handleSubmit">
        发布
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { createPost } from "@/api/post"
import ImageUploader from "@/components/ImageUploader.vue"

const router = useRouter()
const content = ref("")
const images = ref<string[]>([])
const submitting = ref(false)

async function handleSubmit() {
  if (!content.value.trim()) return
  submitting.value = true
  try {
    const res = await createPost({
      content: content.value.trim(),
      images: images.value.length ? images.value : undefined,
    })
    router.push("/post/" + res.id)
  } finally { submitting.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.create-page { max-width: 640px; }
h1 { font-size: 20px; font-weight: 700; margin-bottom: 16px; }
.form-actions {
  display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px;
}
</style>
