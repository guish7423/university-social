<template>
  <dv-border-box-13 class="view-wrapper">
    <div class="page-header">
      <div class="header-left">
        <el-input v-model="searchQuery" placeholder="搜索昵称/OpenID" clearable style="width:240px" @clear="load" @keyup.enter="load" />
        <el-button :icon="Search" @click="load">搜索</el-button>
      </div>
    </div>
    <el-table :data="users" border stripe v-loading="loading" style="width: 100%" @row-click="showDetail">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="头像" width="60">
        <template #default="{ row }">
          <el-avatar :src="row.avatar" size="small">{{ row.nickname?.[0] }}</el-avatar>
        </template>
      </el-table-column>
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="open_id" label="OpenID" min-width="200" show-overflow-tooltip />
      <el-table-column prop="school" label="学校" min-width="160" />
      <el-table-column prop="is_verified" label="认证" width="70">
        <template #default="{ row }">
          <el-tag :type="row.is_verified ? 'success' : 'info'" size="small">{{ row.is_verified ? "是" : "否" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="角色" width="70">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'default'" size="small">{{ row.role || "user" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_banned" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.is_banned ? 'danger' : 'success'" size="small">{{ row.is_banned ? "封禁" : "正常" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="160" />
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button v-if="!row.is_banned" type="warning" size="small" @click.stop="handleBan(row.id)">封禁</el-button>
          <el-button v-else type="success" size="small" @click.stop="handleUnban(row.id)">解封</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="hasMore" class="load-more">
      <el-button :loading="loading" text @click="loadMore">加载更多</el-button>
    </div>

    <el-dialog v-model="detailVisible" title="用户详情" width="500">
      <template v-if="detailUser">
        <div class="detail-header">
          <el-avatar :src="detailUser.avatar" size="large">{{ detailUser.nickname?.[0] }}</el-avatar>
          <div class="detail-meta">
            <h3>{{ detailUser.nickname }}</h3>
            <p class="detail-id">ID: {{ detailUser.id }} | OpenID: {{ detailUser.open_id }}</p>
          </div>
        </div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="学校">{{ detailUser.school || "未设置" }}</el-descriptions-item>
          <el-descriptions-item label="认证">{{ detailUser.is_verified ? "已认证" : "未认证" }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ detailUser.is_banned ? "已封禁" : "正常" }}</el-descriptions-item>
          <el-descriptions-item label="角色">{{ detailUser.role || "user" }}</el-descriptions-item>
          <el-descriptions-item label="注册时间" :span="2">{{ detailUser.created_at }}</el-descriptions-item>
        </el-descriptions>
      </template>
    </el-dialog>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"
import { Search } from "@element-plus/icons-vue"

interface User {
  id: number; nickname: string; avatar: string; open_id: string; school: string
  is_verified: boolean; is_banned: boolean; role: string; created_at: string
}

const users = ref<User[]>([])
const loading = ref(false)
const searchQuery = ref("")
const detailVisible = ref(false)
const detailUser = ref<User | null>(null)
const offset = ref(0)
const hasMore = ref(true)
const PAGE_LIMIT = 50

async function load(reset = true) {
  if (reset) { offset.value = 0; hasMore.value = true }
  loading.value = true
  try {
    const params: any = { offset: offset.value, limit: PAGE_LIMIT }
    if (searchQuery.value.trim()) params.q = searchQuery.value.trim()
    const res = await api.get("/admin/users", { params })
    if (reset) {
      users.value = res.data
    } else {
      users.value = [...users.value, ...res.data]
    }
    hasMore.value = res.data.length >= PAGE_LIMIT
  } catch {} finally { loading.value = false }
}

function loadMore() {
  offset.value += PAGE_LIMIT
  load(false)
}

function showDetail(row: User) {
  detailUser.value = row
  detailVisible.value = true
}

async function handleBan(id: number) {
  try {
    await ElMessageBox.confirm("确定封禁该用户？封禁后用户将无法登录", "提示")
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

onMounted(() => load())
</script>

<style scoped>
.view-wrapper { padding: 8px; }
.page-header { margin-bottom: 12px; display: flex; justify-content: space-between; }
.header-left { display: flex; gap: 8px; }
.load-more { text-align: center; padding: 20px 0; }
.detail-header { display: flex; gap: 16px; margin-bottom: 20px; align-items: center; }
.detail-meta h3 { margin: 0; font-size: 18px; }
.detail-id { margin: 4px 0 0; font-size: 12px; color: #909399; }
</style>
