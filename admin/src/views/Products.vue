<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="products" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column label="作者" width="120">
        <template #default="{ row }">{{ row.author_name || "未知" }}</template>
      </el-table-column>
      <el-table-column prop="title" label="商品名称" min-width="200" show-overflow-tooltip />
      <el-table-column prop="price" label="价格(¥)" width="100">
        <template #default="{ row }">{{ row.price?.toFixed?.(2) ?? row.price }}</template>
      </el-table-column>
      <el-table-column prop="category" label="分类" width="120" />
      <el-table-column prop="condition" label="成色" width="80" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '在售' : row.status === 0 ? '下架' : '删除' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="发布时间" width="170" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface ProductItem {
  id: number; title: string; price: number; category: string
  condition: string; status: number; created_at: string
  author_name: string
}

const products = ref<ProductItem[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/products")
    products.value = res.data
  } catch { ElMessage.error("加载失败") } finally { loading.value = false }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除该商品？", "提示")
    await api.delete(`/admin/products/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
</style>
