<template>
  <div>
    <div class="page-header">
      <el-button :type="showFeatured ? 'default' : 'primary'" @click="showFeatured = false; load()">全部</el-button>
      <el-button :type="showFeatured ? 'primary' : 'default'" @click="showFeatured = true; loadFeatured()">精华帖</el-button>
    </div>
    <el-table :data="showFeatured ? featuredPosts : posts" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column label="作者" width="120">
        <template #default="{ row }">{{ row.author?.nickname || "未知" }}</template>
      </el-table-column>
      <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
      <el-table-column prop="like_count" label="点赞" width="70" />
      <el-table-column prop="comment_count" label="评论" width="70" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-if="!showFeatured" size="small" type="success" @click="handleFeature(row.id)">设为精华</el-button>
          <el-button v-if="showFeatured" size="small" type="warning" @click="handleUnfeature(row.id)">取消精华</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
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
const featuredPosts = ref<Post[]>([])
const loading = ref(false)
const showFeatured = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get("/admin/posts")
    posts.value = res.data
  } catch {} finally { loading.value = false }
}

async function loadFeatured() {
  loading.value = true
  try {
    const res = await api.get("/posts/featured")
    featuredPosts.value = res.data
  } catch {} finally { loading.value = false }
}

async function handleFeature(id: number) {
  try {
    await api.post(`/admin/posts/${id}/feature`)
    ElMessage.success("已设为精华")
  } catch { ElMessage.error("操作失败") }
}

async function handleUnfeature(id: number) {
  try {
    await api.post(`/admin/posts/${id}/unfeature`)
    ElMessage.success("已取消精华")
  } catch { ElMessage.error("操作失败") }
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
