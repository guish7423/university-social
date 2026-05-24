<template>
  <dv-border-box-13 class="view-wrapper">
    <div class="page-header">
      <el-button type="primary" @click="showAdd = true">添加敏感词</el-button>
    </div>
    <el-table :data="wordList" border stripe v-loading="loading" style="width: 100%">
      <el-table-column type="index" label="#" width="60" />
      <el-table-column label="敏感词" min-width="300">
        <template #default="{ row }">{{ row }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="handleRemove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="showAdd" title="添加敏感词" width="400">
      <el-form :model="form" label-width="0">
        <el-input v-model="form.word" placeholder="输入敏感词" @keyup.enter="handleAdd" />
      </el-form>
      <template #footer>
        <el-button @click="showAdd = false">取消</el-button>
        <el-button type="primary" @click="handleAdd" :loading="adding">确认</el-button>
      </template>
    </el-dialog>
  </dv-border-box-13>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

const wordList = ref<string[]>([])
const loading = ref(false)
const showAdd = ref(false)
const adding = ref(false)
const form = reactive({ word: "" })

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/sensitive-words")
    wordList.value = res.data || []
  } catch { ElMessage.error("加载失败") } finally { loading.value = false }
}

async function handleAdd() {
  if (!form.word.trim()) { ElMessage.warning("请输入敏感词"); return }
  adding.value = true
  try {
    await api.post("/admin/sensitive-words", { word: form.word.trim() })
    ElMessage.success("已添加")
    showAdd.value = false
    form.word = ""
    load()
  } catch { ElMessage.error("操作失败") } finally { adding.value = false }
}

async function handleRemove(word: string) {
  try {
    await ElMessageBox.confirm(`确定删除敏感词「${word}」？`, "提示")
    await api.delete(`/admin/sensitive-words/${encodeURIComponent(word)}`)
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
