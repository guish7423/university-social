<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #e6f7ff; color: #1890ff"><User /></div>
            <div>
              <div class="stat-value">{{ data.user_count }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #fff7e6; color: #fa8c16"><Document /></div>
            <div>
              <div class="stat-value">{{ data.post_count }}</div>
              <div class="stat-label">帖子总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f6ffed; color: #52c41a"><Collection /></div>
            <div>
              <div class="stat-value">{{ data.circle_count }}</div>
              <div class="stat-label">圈子总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"

const data = ref({ user_count: 0, post_count: 0, circle_count: 0 })

onMounted(async () => {
  try {
    const res = await api.get("/admin/dashboard")
    data.value = res.data
  } catch {}
})
</script>

<style scoped>
.stat-card { display: flex; align-items: center; gap: 20px; }
.stat-icon { width: 60px; height: 60px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 28px; }
.stat-value { font-size: 28px; font-weight: 700; }
.stat-label { font-size: 14px; color: #999; margin-top: 4px; }
</style>
