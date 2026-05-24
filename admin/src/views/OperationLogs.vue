<template>
  <dv-border-box-13 class="view-wrapper">
    <div class="page-header">
      <span class="page-title">操作日志</span>
    </div>
    <el-table :data="logs" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="管理员" width="120">
        <template #default="{ row }">{{ row.admin_name || `#${row.admin_id}` }}</template>
      </el-table-column>
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-tag :type="tagType(row.action)" size="small">{{ actionLabel(row.action) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="对象" width="160">
        <template #default="{ row }">{{ row.target_type }}#{{ row.target_id }}</template>
      </el-table-column>
      <el-table-column label="详情" min-width="200">
        <template #default="{ row }">
          <span class="detail-text">{{ formatDetail(row.detail) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="时间" width="180">
        <template #default="{ row }">{{ row.created_at?.slice(0, 19)?.replace('T', ' ') }}</template>
      </el-table-column>
    </el-table>
    <div class="pagination-wrapper" v-if="total > limit">
      <el-pagination
        v-model:current-page="page"
        :page-size="limit"
        :total="total"
        layout="prev, pager, next"
        @current-change="load"
        background
        small
      />
    </div>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"

interface OperationLog {
  id: number
  admin_id: number
  admin_name: string
  action: string
  target_type: string
  target_id: number
  detail: any
  created_at: string
}

const logs = ref<OperationLog[]>([])
const page = ref(1)
const limit = 50
const total = ref(0)
const loading = ref(false)

const actionLabels: Record<string, string> = {
  ban_user: "封禁用户",
  unban_user: "解封用户",
  delete_post: "删除帖子",
  delete_circle: "删除圈子",
}

const actionTags: Record<string, string> = {
  ban_user: "danger",
  unban_user: "success",
  delete_post: "danger",
  delete_circle: "danger",
}

function actionLabel(action: string): string {
  return actionLabels[action] || action
}

function tagType(action: string): string {
  return actionTags[action] || "info"
}

function formatDetail(detail: any): string {
  if (!detail) return "-"
  if (typeof detail === "string") return detail
  return JSON.stringify(detail)
}

async function load(p?: number) {
  if (p !== undefined) page.value = p
  loading.value = true
  try {
    const res = await api.get("/admin/operation-logs", {
      params: { page: page.value, limit: limit },
    })
    logs.value = res.data.logs || []
    total.value = res.data.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => load())
</script>

<style scoped>
.view-wrapper { padding: 8px; }
.page-header { margin-bottom: 12px; }
.page-title { font-size: 16px; font-weight: 600; }
.pagination-wrapper { margin-top: 12px; display: flex; justify-content: center; }
.detail-text { 
  max-width: 300px; 
  overflow: hidden; 
  text-overflow: ellipsis; 
  white-space: nowrap; 
  display: inline-block; 
  vertical-align: middle;
}
</style>
