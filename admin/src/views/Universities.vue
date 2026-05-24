<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="items" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="高校名称" min-width="200" show-overflow-tooltip />
      <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
      <el-table-column prop="member_count" label="成员数" width="100" />
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
  id: number; name: string; description: string
  member_count: number; created_at: string
}

const items = ref<Item[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/universities")
    items.value = res.data
  } catch { ElMessage.error("加载失败") } finally { loading.value = false }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除该高校？该操作会同时移除所有成员记录。", "提示")
    await api.delete(`/admin/universities/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
</style>
