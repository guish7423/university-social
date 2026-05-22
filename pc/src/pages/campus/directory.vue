<template>
  <div class="campus-directory">
    <div class="page-header">
      <h1>通讯录</h1>
      <el-input v-model="searchQ" placeholder="搜索姓名或部门" clearable style="width:240px" @input="handleSearch" />
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
    <div v-else-if="!filtered.length" class="empty-state"><el-empty description="未找到联系信息" /></div>
    <el-table v-else :data="filtered" stripe>
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="department" label="部门" width="150" />
      <el-table-column prop="phone" label="电话" width="150" />
      <el-table-column prop="email" label="邮箱" min-width="200" />
      <el-table-column prop="office" label="办公室" width="120" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { listDirectory } from "@/api/campus"
import type { CampusContact } from "@/api/campus"

const contacts = ref<CampusContact[]>([])
const loading = ref(true)
const searchQ = ref("")

const filtered = computed(() => {
  const q = searchQ.value.toLowerCase()
  if (!q) return contacts.value
  return contacts.value.filter(c =>
    c.name.toLowerCase().includes(q) || c.department.toLowerCase().includes(q)
  )
})

function handleSearch() {}

onMounted(async () => {
  try { contacts.value = await listDirectory() }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
.campus-directory { max-width: 900px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; h1 { font-size: 22px; font-weight: 700; } }
.loading-wrap, .empty-state { padding: 40px 0; }
</style>
