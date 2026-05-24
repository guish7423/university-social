<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="items" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column label="用户" width="120">
        <template #default="{ row }">{{ row.nickname || "未知" }}</template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="category" label="类型" width="80" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '已找回' : '寻找中' }}
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

interface Item {
  id: number; title: string; category: string; status: number
  created_at: string; nickname: string
}

const items = ref<Item[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/lost-items")
    items.value = res.data
  } catch { ElMessage.error("加载失败") } finally { loading.value = false }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除该失物招领？", "提示")
    await api.delete(`/admin/lost-items/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
</style>
