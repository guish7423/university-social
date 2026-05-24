<template>
  <div class="create-page">
  <PageHeader title="创建圈子" />
    <div class="card-wrap stagger-item">
      <el-form @submit.prevent="handleSubmit">
        <el-form-item label="圈子名称">
          <el-input v-model="name" placeholder="输入圈子名称" maxlength="30" />
        </el-form-item>
        <el-form-item label="圈子简介">
          <el-input v-model="description" type="textarea" :rows="4" placeholder="介绍一下这个圈子..." maxlength="500" />
        </el-form-item>
        <div class="form-actions">
          <el-button @click="$router.back()">取消</el-button>
          <el-button type="primary" :loading="submitting" :disabled="!name.trim()" native-type="submit">创建</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { createCircle } from "@/api/circle"
import PageHeader from "@/components/PageHeader.vue"

const router = useRouter()
const name = ref("")
const description = ref("")
const submitting = ref(false)

async function handleSubmit() {
  if (!name.value.trim()) return
  submitting.value = true
  try {
    const res = await createCircle({ name: name.value.trim(), description: description.value.trim() })
    router.push("/circles/" + res.id)
  } finally { submitting.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.create-page { max-width: 560px; }
.card-wrap {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px;
}
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
</style>
