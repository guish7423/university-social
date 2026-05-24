<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="items" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="课程名称" min-width="180" show-overflow-tooltip />
      <el-table-column prop="code" label="课程代码" width="120" />
      <el-table-column prop="department" label="院系" width="150" show-overflow-tooltip />
      <el-table-column prop="teacher" label="授课教师" width="120" />
      <el-table-column prop="credit" label="学分" width="70" />
      <el-table-column label="评分" width="80">
        <template #default="{ row }">
          <el-tag size="small">{{ row.avg_rating?.toFixed?.(1) ?? '暂无' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170" />
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
  id: number; name: string; code: string; department: string
  teacher: string; credit: number; avg_rating: number; created_at: string
}

const items = ref<Item[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/courses")
    items.value = res.data
  } catch { ElMessage.error("加载失败") } finally { loading.value = false }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除该课程？", "提示")
    await api.delete(`/admin/courses/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
</style>
