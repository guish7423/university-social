<template>
  <div class="create-page">
  <PageHeader title="发布失物信息" />
    <div class="card-wrap stagger-item">
      <el-radio-group v-model="category" style="margin-bottom:16px">
        <el-radio-button value="寻物">我丢了东西</el-radio-button>
        <el-radio-button value="招领">我捡到了东西</el-radio-button>
      </el-radio-group>
      <el-form @submit.prevent="handleSubmit">
        <el-form-item label="标题">
          <el-input v-model="title" placeholder="如：丢失校园卡" maxlength="60" />
        </el-form-item>
        <el-form-item label="地点">
          <el-input v-model="location" placeholder="丢失/捡到的地点" />
        </el-form-item>
        <el-form-item label="联系方式">
          <el-input v-model="contact" placeholder="QQ/微信/手机号" />
        </el-form-item>
        <el-form-item label="详细描述">
          <el-input v-model="description" type="textarea" :rows="4" placeholder="描述物品特征..." />
        </el-form-item>
        <div class="form-actions">
          <el-button @click="$router.back()">取消</el-button>
          <el-button type="primary" :loading="submitting" :disabled="!title.trim()" native-type="submit">发布</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { createLostItem } from "@/api/found"
import PageHeader from "@/components/PageHeader.vue"

const router = useRouter()
const title = ref("")
const category = ref("寻物")
const location = ref("")
const contact = ref("")
const description = ref("")
const submitting = ref(false)

async function handleSubmit() {
  if (!title.value.trim()) return
  submitting.value = true
  try {
    const res = await createLostItem({
      title: title.value.trim(),
      category: category.value,
      location: location.value.trim(),
      contact: contact.value.trim(),
      description: description.value.trim(),
    })
    router.push("/found/" + res.id)
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
