<template>
  <dv-border-box-13 class="view-wrapper">
    <div class="page-header">
      <el-radio-group v-model="statusFilter" @change="load">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button value="pending">待处理</el-radio-button>
        <el-radio-button value="resolved">已处理</el-radio-button>
        <el-radio-button value="dismissed">已驳回</el-radio-button>
      </el-radio-group>
    </div>
    <el-table :data="reports" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="content_type" label="类型" width="80">
        <template #default="{ row }">{{ typeMap[row.content_type] || row.content_type }}</template>
      </el-table-column>
      <el-table-column prop="content_id" label="内容ID" width="80" />
      <el-table-column prop="reason" label="原因" width="120" />
      <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="举报人" width="100">
        <template #default="{ row }">{{ row.reporter_id }}</template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="160" />
      <el-table-column label="操作" width="160" fixed="right" v-if="statusFilter === 'pending' || !statusFilter">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" size="small" type="success" @click="handleResolve(row.id)">处理</el-button>
          <el-button v-if="row.status === 'pending'" size="small" type="warning" @click="handleDismiss(row.id)">驳回</el-button>
        </template>
      </el-table-column>
    </el-table>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface Report {
  id: number; reporter_id: number; content_type: string; content_id: number
  reason: string; detail: string; status: string; created_at: string
}

const typeMap: Record<string, string> = {
  post: "帖子", comment: "评论", whisper: "树洞", whisper_comment: "树洞评论",
  product: "商品", activity: "活动", user: "用户",
}

const reports = ref<Report[]>([])
const loading = ref(false)
const statusFilter = ref("pending")

async function load() {
  loading.value = true
  try {
    const params: any = { offset: 0, limit: 50 }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get("/admin/reports", { params })
    reports.value = res.data
  } catch {} finally { loading.value = false }
}

function statusTag(s: string) {
  return s === "pending" ? "danger" : s === "resolved" ? "success" : "info"
}
function statusLabel(s: string) {
  return s === "pending" ? "待处理" : s === "resolved" ? "已处理" : "已驳回"
}

async function handleResolve(id: number) {
  try {
    await api.put(`/admin/reports/${id}/resolve`)
    ElMessage.success("已处理")
    load()
  } catch { ElMessage.error("操作失败") }
}

async function handleDismiss(id: number) {
  try {
    await api.put(`/admin/reports/${id}/dismiss`)
    ElMessage.success("已驳回")
    load()
  } catch { ElMessage.error("操作失败") }
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
.page-header { margin-bottom: 12px; }
</style>
