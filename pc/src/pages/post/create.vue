<template>
  <div class="create-page">
  <PageHeader title="发布动态" />
    <div class="card-wrap stagger-item">
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
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { createPost } from "@/api/post"
import PageHeader from "@/components/PageHeader.vue"

import { useFormDraft } from "@/composables/useFormDraft"

import ImageUploader from "@/components/ImageUploader.vue"

const router = useRouter()


const content = ref("")
const images = ref<string[]>([])
const submitting = ref(false)
const { clearDraft } = useFormDraft("post", { content })


async function handleSubmit() {
  if (!content.value.trim()) return
  submitting.value = true
  try {
    const res = await createPost({
      content: content.value.trim(),
      images: images.value.length ? images.value : undefined,
    })
    clearDraft()

    router.push("/post/" + res.id)
  } finally { submitting.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.create-page { max-width: 640px; }

.card-wrap {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px;
}
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
</style>
