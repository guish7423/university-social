<template>
  <dv-border-box-13 class="view-wrapper">
    <div class="page-header">
      <el-button type="primary" @click="showCreate = true">发布公告</el-button>
    </div>
    <el-table :data="announcements" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="标题" width="250" />
      <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="160" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">{{ row.is_active ? "显示" : "隐藏" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" :type="row.is_active ? 'warning' : 'success'" @click="handleToggle(row.id)">
            {{ row.is_active ? "隐藏" : "显示" }}
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="showCreate" title="发布公告" width="500">
      <el-form :model="form" label-width="60px">
        <el-form-item label="标题"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="内容"><el-input v-model="form.content" type="textarea" :rows="4" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">发布</el-button>
      </template>
    </el-dialog>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface Announcement {
  id: number; title: string; content: string; is_active: boolean; created_at: string
}

const announcements = ref<Announcement[]>([])
const loading = ref(false)
const showCreate = ref(false)
const creating = ref(false)
const form = reactive({ title: "", content: "" })

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/announcements")
    announcements.value = res.data
  } catch {} finally { loading.value = false }
}

async function handleCreate() {
  if (!form.title || !form.content) { ElMessage.warning("请填写标题和内容"); return }
  creating.value = true
  try {
    await api.post("/admin/announcements", form)
    ElMessage.success("已发布")
    showCreate.value = false
    form.title = ""; form.content = ""
    load()
  } catch {} finally { creating.value = false }
}

async function handleToggle(id: number) {
  try {
    await api.put(`/admin/announcements/${id}/toggle`)
    ElMessage.success("已更新")
    load()
  } catch { ElMessage.error("操作失败") }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除？", "提示")
    await api.delete(`/admin/announcements/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>

<style scoped>
.view-wrapper { padding: 8px; }
.page-header { margin-bottom: 12px; }
</style>
