<template>
  <div class="create-page">
  <PageHeader title="发布闲置" />
    <div class="card-wrap stagger-item">
      <el-form @submit.prevent="handleSubmit">
        <el-form-item label="标题">
          <el-input v-model="title" placeholder="商品名称" maxlength="60" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="category" placeholder="如：教材、数码、生活等" />
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="price" :min="0" :precision="2" :step="10" />
        </el-form-item>
        <el-form-item label="成色">
          <el-select v-model="conditionText" placeholder="选择成色">
            <el-option label="全新" value="全新" />
            <el-option label="几乎全新" value="几乎全新" />
            <el-option label="轻微使用痕迹" value="轻微使用痕迹" />
            <el-option label="有明显使用痕迹" value="有明显使用痕迹" />
            <el-option label="老旧" value="老旧" />
          </el-select>
        </el-form-item>
        <el-form-item label="联系方式">
          <el-input v-model="contact" placeholder="QQ/微信/手机号" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="description" type="textarea" :rows="4" placeholder="描述商品情况..." />
        </el-form-item>
        <div class="form-actions">
          <el-button @click="$router.back()">取消</el-button>
          <el-button type="primary" :loading="submitting" :disabled="!title.trim() || !price" native-type="submit">
            发布
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { createProduct } from "@/api/product"
import PageHeader from "@/components/PageHeader.vue"

const router = useRouter()
const title = ref("")
const category = ref("")
const price = ref(0)
const conditionText = ref("")
const contact = ref("")
const description = ref("")
const submitting = ref(false)

async function handleSubmit() {
  if (!title.value.trim() || !price.value) return
  submitting.value = true
  try {
    const res = await createProduct({
      title: title.value.trim(),
      category: category.value.trim() || "其他",
      price: price.value,
      condition: conditionText.value,
      contact: contact.value.trim(),
      description: description.value.trim(),
    })
    router.push("/products/" + res.id)
  } finally { submitting.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.create-page { max-width: 560px; margin: 0 auto; }
.card-wrap {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px;
}
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
</style>
