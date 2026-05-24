<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="circles" border stripe v-loading="loading" style="width: 100%" empty-text="暂无圈子数据">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="description" label="简介" min-width="240" show-overflow-tooltip />
      <el-table-column prop="member_count" label="成员数" width="80" />
      <el-table-column prop="post_count" label="帖子数" width="80" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="danger" size="small" class="op-btn" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface CircleItem {
  id: number; name: string; description: string; member_count: number
  post_count: number; created_at: string
}

const circles = ref<CircleItem[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/circles")
    circles.value = res.data
  } catch { ElMessage.error("加载失败") } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除该圈子及其所有内容？", "提示")
    await api.delete(`/admin/circles/${id}`)
    ElMessage.success("已删除")
    load()
  } catch { ElMessage.error("删除失败") }
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
</style>
