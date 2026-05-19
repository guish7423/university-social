<template>
  <div>
    <el-table :data="posts" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="作者" width="120">
        <template #default="{ row }">{{ row.author?.nickname || "未知" }}</template>
      </el-table-column>
      <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
      <el-table-column prop="like_count" label="点赞" width="70" />
      <el-table-column prop="comment_count" label="评论" width="70" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"
import { ElMessage, ElMessageBox } from "element-plus"

interface Post {
  id: number; content: string; like_count: number; comment_count: number
  author: { nickname: string }
}

const posts = ref<Post[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/posts")
    posts.value = res.data
  } catch {} finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确定删除该帖子？", "提示")
    await api.delete(`/admin/posts/${id}`)
    ElMessage.success("已删除")
    load()
  } catch {}
}

onMounted(load)
</script>
