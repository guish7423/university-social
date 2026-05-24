<template>
  <div class="university-page">
    <PageHeader title="高校社区" subtitle="找到你的学校，加入校友圈" />

    <div v-if="loading" class="loading-wrap">
      <el-skeleton :rows="4" animated />
    </div>

    <div v-else-if="universities.length === 0" class="empty-state">
      <el-empty description="暂无高校" />
    </div>

    <div v-else class="university-grid">
      <div
        v-for="uni in universities"
        :key="uni.id"
        class="university-card"
        @click="$router.push(`/university/${uni.id}`)"
      >
        <div class="card-header">
          <div class="uni-icon" :style="{ backgroundColor: stringToColor(uni.name) }">
            {{ uni.name.charAt(0) }}
          </div>
          <div class="uni-info">
            <h3>{{ uni.name }}</h3>
            <span class="uni-domain" v-if="uni.domain">{{ uni.domain }}</span>
          </div>
        </div>
        <div class="card-stats">
          <span class="stat">
            <el-icon><UserFilled /></el-icon>
            {{ uni.member_count }} 人加入
          </span>
        </div>
        <div class="card-actions">
          <el-button
            v-if="uni.is_member"
            type="success"
            size="small"
            round
            @click.stop="$router.push(`/university/${uni.id}`)"
          >
            已加入
          </el-button>
          <el-button
            v-else
            type="primary"
            size="small"
            round
            @click.stop="handleJoin(uni)"
            :loading="joiningId === uni.id"
          >
            加入
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { UserFilled } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"
import PageHeader from "@/components/PageHeader.vue"

interface University {
  id: number
  name: string
  domain: string
  logo_url: string
  member_count: number
  is_member: boolean
}

const router = useRouter()
const universities = ref<University[]>([])
const loading = ref(true)
const joiningId = ref<number | null>(null)

function stringToColor(str: string): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  const hue = Math.abs(hash) % 360
  return `hsl(${hue}, 50%, 45%)`
}

async function fetchUniversities() {
  try {
    const token = localStorage.getItem("token")
    const res = await fetch("/api/v1/universities", {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      universities.value = await res.json()
    }
  } catch { /* ignore */ }
  finally { loading.value = false }
}

async function handleJoin(uni: University) {
  joiningId.value = uni.id
  try {
    const token = localStorage.getItem("token")
    const res = await fetch(`/api/v1/universities/${uni.id}/join`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      uni.is_member = true
      uni.member_count++
      ElMessage.success(`已加入 ${uni.name}`)
    } else {
      const data = await res.json()
      ElMessage.error(data.error || "加入失败")
    }
  } catch { ElMessage.error("网络错误") }
  finally { joiningId.value = null }
}

onMounted(fetchUniversities)
</script>

<style scoped lang="scss">
.university-page {
  max-width: 680px;
  margin: 0 auto;
  padding: 24px;
  width: 100%;
}

.loading-wrap {
  padding: 40px 0;
}

.empty-state {
  padding: 60px 0;
}

.university-grid {
  display: grid;
  gap: 16px;
  margin-top: 20px;
}

.university-card {
  background: var(--bg-card, #fff);
  border: 1px solid var(--border-color, #e4e7ed);
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
    border-color: var(--el-color-primary);
  }
}

.card-header {
  display: flex;
  align-items: center;
  gap: 14px;
}

.uni-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 20px;
  font-weight: 700;
  flex-shrink: 0;
}

.uni-info {
  flex: 1;
  min-width: 0;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    line-height: 1.4;
  }

  .uni-domain {
    font-size: 13px;
    color: var(--text-secondary, #909399);
    margin-top: 2px;
    display: inline-block;
  }
}

.card-stats {
  margin-top: 12px;
  font-size: 13px;
  color: var(--text-secondary, #909399);

  .stat {
    display: inline-flex;
    align-items: center;
    gap: 4px;

    .el-icon {
      font-size: 14px;
    }
  }
}

.card-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
</style>
