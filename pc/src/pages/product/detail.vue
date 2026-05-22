<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!product" class="empty-state"><el-empty description="商品不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>
    <div class="product-detail">
      <div class="product-images">
        <img v-for="(img, i) in product.images" :key="i" :src="img" />
      </div>
      <div class="product-info">
        <el-tag v-if="product.status === 1" type="danger" size="small">已售出</el-tag>
        <h1>{{ product.title }}</h1>
        <div class="price-row">
          <span class="price">¥{{ product.price }}</span>
          <span v-if="product.original_price" class="original">¥{{ product.original_price }}</span>
        </div>
        <div class="meta-row">
          <span>{{ product.category }}</span>
          <span>{{ product.condition }}</span>
        </div>
        <p class="description">{{ product.description }}</p>
        <div class="contact" v-if="product.contact">
          <el-icon><Phone /></el-icon> {{ product.contact }}
        </div>
        <div class="actions" v-if="product.is_owner">
          <el-button v-if="product.status === 0" @click="handleMarkSold">标记已售</el-button>
          <el-button type="danger" @click="handleDelete">删除</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getProduct, deleteProduct, markSold } from "@/api/product"
import type { ProductData } from "@/api/product"
import { ArrowLeft } from "@element-plus/icons-vue"
import { ElMessageBox } from "element-plus"

const route = useRoute()
const router = useRouter()
const product = ref<ProductData | null>(null)
const loading = ref(true)

async function handleDelete() {
  if (!product.value) return
  try { await deleteProduct(product.value.id); router.push("/products") }
  catch { /* handled by interceptor */ }
}

async function handleMarkSold() {
  if (!product.value) return
  try { await markSold(product.value.id); product.value.status = 1 }
  catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try { product.value = await getProduct(id) }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.detail-page { max-width: 800px; }
.back-btn { margin-bottom: 16px; }

.product-detail { display: flex; gap: 24px;
  .product-images {
    flex-shrink: 0; width: 300px;
    img { width: 100%; border-radius: $radius-md; margin-bottom: 8px; }
  }
  .product-info { flex: 1;
    h1 { font-size: 20px; font-weight: 700; margin: 8px 0; }
    .price-row {
      .price { font-size: 24px; font-weight: 700; color: $primary; }
      .original { font-size: 14px; color: $text-muted; text-decoration: line-through; margin-left: 8px; }
    }
    .meta-row { display: flex; gap: 12px; font-size: 13px; color: $text-secondary; margin: 8px 0; }
    .description { font-size: 14px; line-height: 1.8; margin: 12px 0; }
    .contact { font-size: 14px; color: $accent-blue; margin: 12px 0; }
    .actions { margin-top: 16px; display: flex; gap: 10px; }
  }
}
</style>
