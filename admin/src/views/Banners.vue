<template>
  <div>
    <div class="page-header">
      <el-button type="primary" @click="showCreate = true">添加 Banner</el-button>
    </div>
    <el-table :data="banners" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="标题" width="200" />
      <el-table-column label="图片" width="120">
        <template #default="{ row }">
          <el-image :src="row.image_url" style="width:80px;height:50px" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="link_url" label="链接" min-width="200" show-overflow-tooltip />
      <el-table-column prop="sort_order" label="排序" width="70" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">{{ row.is_active ? "启用" : "停用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" :type="row.is_active ? 'warning' : 'success'" @click="handleToggle(row.id)">
            {{ row.is_active ? "停用" : "启用" }}
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreate" title="添加 Banner" width="500">
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="图片URL"><el-input v-model="form.image_url" placeholder="https://..." /></el-form-item>
        <el-form-item label="链接"><el-input v-model="form.link_url" placeholder="可选" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface Banner {
  id: number; title: string; image_url: string; link_url: string
  sort_order: number; is_active: boolean; created_at: string
}

const banners = ref<Banner[]>([])
const loading = ref(false)
const showCreate = ref(false)
const creating = ref(false)
const form = reactive({ title: "", image_url: "", link_url: "", sort_order: 0 })

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/banners")
    banners.value = res.data
  } catch {} finally { loading.value = false }
}

async function handleCreate() {
  if (!form.title || !form.image_url) { ElMessage.warning("请填写标题和图片"); return }
  creating.value = true
  try {
    await api.post("/admin/banners", form)
    ElMessage.success("已添加")
    showCreate.value = false
    Object.assign(form, { title: "", image_url: "", link_url: "", sort_order: 0 })
    load()
  } catch {} finally { creating.value = false }
}

async function handleToggle(id: number) {
  try {
    await api.put(`/admin/banners/${id}/toggle`)
    ElMessage.success("已更新")
    load()
  } catch { ElMessage.error("操作失败") }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除？", "提示")
    await api.delete(`/admin/banners/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>
