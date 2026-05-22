<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="users" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="school" label="学校" min-width="160" />
      <el-table-column prop="is_verified" label="认证" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_verified ? 'success' : 'info'" size="small">
            {{ row.is_verified ? "是" : "否" }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_banned" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_banned ? 'danger' : 'success'" size="small">
            {{ row.is_banned ? "封禁" : "正常" }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-if="!row.is_banned" type="warning" size="small" @click="handleBan(row.id)">封禁</el-button>
          <el-button v-else type="success" size="small" @click="handleUnban(row.id)">解封</el-button>
        </template>
      </el-table-column>
    </el-table>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface User {
  id: number; nickname: string; avatar: string; school: string
  is_verified: boolean; is_banned: boolean; created_at: string
}

const users = ref<User[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/users")
    users.value = res.data
  } catch {} finally {
    loading.value = false
  }
}

async function handleBan(id: number) {
  try {
    await ElMessageBox.confirm("确定封禁该用户？", "提示")
    await api.put(`/admin/users/${id}/ban`)
    ElMessage.success("已封禁")
    load()
  } catch {}
}

async function handleUnban(id: number) {
  try {
    await ElMessageBox.confirm("确定解封该用户？", "提示")
    await api.put(`/admin/users/${id}/unban`)
    ElMessage.success("已解封")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
</style>
