<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!item" class="empty-state"><el-empty description="信息不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn stagger-item">返回</el-button>
    <div class="found-detail">
      <div class="header">
        <el-tag :type="item.category === '寻物' ? 'warning' : 'success'">{{ item.category }}</el-tag>
        <span :class="['status', item.status === 0 ? 'active' : 'resolved']">
          {{ item.status === 0 ? '进行中' : '已解决' }}
        </span>
      </div>
      <h1>{{ item.title }}</h1>
      <p class="description">{{ item.description }}</p>
      <div class="meta">
        <div v-if="item.location"><el-icon><Location /></el-icon> {{ item.location }}</div>
        <div v-if="item.contact"><el-icon><Phone /></el-icon> {{ item.contact }}</div>
      </div>
      <div v-if="item.images?.length" class="images">
        <img v-for="(img, i) in item.images" :key="i" :src="img" loading="lazy" @click="lightboxIndex = i" />
      </div>
      <div class="actions" v-if="item.is_owner && item.status === 0">
        <el-button type="success" @click="handleResolve">标记已解决</el-button>
        <el-button type="danger" @click="handleDelete">删除</el-button>
      </div>
    </div>
  </div>
    <ImagePreview :images="item?.images || []" v-model="lightboxIndex" />

</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getLostItem, deleteLostItem, updateLostItemStatus } from "@/api/found"
import type { LostItemData } from "@/api/found"
import { ArrowLeft } from "@element-plus/icons-vue"
import { ElMessageBox } from "element-plus"
import ImagePreview from "@/components/ImagePreview.vue"

const lightboxIndex = ref<number | null>(null)

const route = useRoute()
const router = useRouter()
const item = ref<LostItemData | null>(null)
const loading = ref(true)

async function handleResolve() {
  if (!item.value) return
  try { await updateLostItemStatus(item.value.id, 1); item.value.status = 1 }
  catch { /* handled by interceptor */ }
}

async function handleDelete() {
  if (!item.value) return
  try { await deleteLostItem(item.value.id); router.push("/found") }
  catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try { item.value = await getLostItem(id) }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.detail-page { max-width: 680px; }
.back-btn { margin-bottom: 16px; }
.found-detail {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px;
  .header { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }
  .status.active { color: $accent-green; font-size: 12px; }
  .status.resolved { color: $text-muted; font-size: 12px; }
  h1 { font-size: 20px; font-weight: 700; margin-bottom: 12px; }
  .description { font-size: 14px; line-height: 1.8; margin-bottom: 16px; }
  .meta { font-size: 13px; color: $text-secondary; margin-bottom: 12px; .el-icon { vertical-align: middle; margin-right: 4px; } }
  .images { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 16px;
    img { width: 100px; height: 100px; object-fit: cover; border-radius: $radius-sm; transition: transform 0.3s ease; cursor: zoom-in;
      &:hover { transform: scale(1.08); }
    }
  }
  .actions { padding-top: 16px; border-top: 1px solid $border-color; }
}
</style>
