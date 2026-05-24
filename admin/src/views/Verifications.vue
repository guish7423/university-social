<template>
  <dv-border-box-13 class="view-wrapper">
    <el-table :data="items" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column label="用户" width="120">
        <template #default="{ row }">{{ row.nickname || "未知" }}</template>
      </el-table-column>
      <el-table-column prop="code" label="验证码" width="130" />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column label="认证状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'verified' ? 'success' : row.status === 'expired' ? 'info' : 'warning'" size="small">
            {{ row.status === 'verified' ? '已验证' : row.status === 'expired' ? '已过期' : '待验证' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="已认证" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_verified ? 'success' : 'danger'" size="small">
            {{ row.is_verified ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="提交时间" width="170" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status !== 'verified' && row.status !== 'expired'"
            size="small" type="primary" @click="handleApprove(row.id)">通过</el-button>
          <el-button v-if="row.status !== 'verified' && row.status !== 'expired'"
            size="small" type="danger" @click="handleReject(row.id)">拒绝</el-button>
          <span v-else class="text-muted">—</span>
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
  id: number; user_id: number; code: string; type: string
  status: string; created_at: string; nickname: string; is_verified: boolean
}

const items = ref<Item[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/verifications")
    items.value = res.data
  } catch { ElMessage.error("加载失败") } finally { loading.value = false }
}

async function handleApprove(id: number) {
  try {
    await ElMessageBox.confirm("确认通过该认证申请？", "提示")
    await api.post(`/admin/verifications/${id}/approve`)
    ElMessage.success("已通过认证")
    load()
  } catch {}
}

async function handleReject(id: number) {
  try {
    await ElMessageBox.confirm("确认拒绝该认证申请？", "提示")
    await api.post(`/admin/verifications/${id}/reject`)
    ElMessage.success("已拒绝")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
.text-muted { color: var(--el-text-color-placeholder); font-size: 13px; }
</style>
