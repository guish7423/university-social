<template>
  <div class="uni-detail-page">
    <template v-if="loading">
      <div class="loading-wrap">
        <el-skeleton :rows="6" animated />
      </div>
    </template>

    <template v-else-if="university">
      <div class="uni-header">
        <div class="uni-icon" :style="{ backgroundColor: iconColor }">
          {{ university.name.charAt(0) }}
        </div>
        <div class="uni-meta">
          <h1>{{ university.name }}</h1>
          <p class="uni-domain" v-if="university.domain">{{ university.domain }}</p>
          <p class="uni-count">{{ university.member_count }} 人已加入</p>
        </div>
        <div class="uni-actions">
          <el-button
            v-if="university.is_member"
            type="danger"
            plain
            size="small"
            @click="handleLeave"
            :loading="leaving"
          >
            退出高校
          </el-button>
          <el-button
            v-else
            type="primary"
            size="small"
            @click="handleJoin"
            :loading="joining"
          >
            加入
          </el-button>
        </div>
      </div>

      <el-divider />

      <div class="section">
        <h2 class="section-title">成员 ({{ members.length }})</h2>
        <div class="members-grid">
          <div
            v-for="member in members"
            :key="member.id"
            class="member-card"
            @click="$router.push(`/user/${member.user_id}`)"
          >
            <el-avatar :size="40" :src="member.user?.avatar">
              {{ member.user?.nickname?.charAt(0) }}
            </el-avatar>
            <div class="member-info">
              <span class="member-name">
                {{ member.user?.nickname }}
                <el-tag v-if="member.is_admin" size="small" type="warning">管理员</el-tag>
              </span>
            </div>
          </div>
        </div>
        <div v-if="members.length === 0" class="empty-members">
          <el-empty description="暂无成员" :image-size="80" />
        </div>
      </div>
    </template>

    <div v-else class="error-state">
      <el-result icon="error" title="高校不存在" sub-title="请检查链接是否正确">
        <template #extra>
          <el-button type="primary" @click="$router.push('/university')">
            返回高校列表
          </el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { ElMessage, ElMessageBox } from "element-plus"
import PageHeader from "@/components/PageHeader.vue"

interface User {
  id: number
  nickname: string
  avatar: string
}

interface University {
  id: number
  name: string
  domain: string
  logo_url: string
  member_count: number
  is_member: boolean
}

interface UniversityMember {
  id: number
  user_id: number
  university_id: number
  is_admin: boolean
  user?: User
}

const route = useRoute()
const router = useRouter()
const university = ref<University | null>(null)
const members = ref<UniversityMember[]>([])
const loading = ref(true)
const joining = ref(false)
const leaving = ref(false)

const iconColor = computed(() => {
  if (!university.value) return "#409eff"
  let hash = 0
  const str = university.value.name
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  return `hsl(${Math.abs(hash) % 360}, 50%, 45%)`
})

async function fetchDetail() {
  const id = route.params.id
  try {
    const token = localStorage.getItem("token")
    const headers = { Authorization: `Bearer ${token}` }

    const [uniRes, membersRes] = await Promise.all([
      fetch(`/api/v1/universities/${id}`, { headers }),
      fetch(`/api/v1/universities/${id}/members`, { headers }),
    ])

    if (uniRes.ok) university.value = await uniRes.json()
    if (membersRes.ok) members.value = await membersRes.json()
  } catch { /* ignore */ }
  finally { loading.value = false }
}

async function handleJoin() {
  if (!university.value) return
  joining.value = true
  try {
    const token = localStorage.getItem("token")
    const res = await fetch(`/api/v1/universities/${university.value.id}/join`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      university.value.is_member = true
      university.value.member_count++
      ElMessage.success("加入成功")
      members.value = [] // refresh
      fetchDetail()
    } else {
      const data = await res.json()
      ElMessage.error(data.error || "加入失败")
    }
  } catch { ElMessage.error("网络错误") }
  finally { joining.value = false }
}

async function handleLeave() {
  if (!university.value) return
  try {
    await ElMessageBox.confirm("确定退出该高校？", "提示", { type: "warning" })
  } catch { return }

  leaving.value = true
  try {
    const token = localStorage.getItem("token")
    const res = await fetch(`/api/v1/universities/${university.value.id}/leave`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      university.value.is_member = false
      university.value.member_count--
      ElMessage.success("已退出")
      members.value = []
      fetchDetail()
    } else {
      const data = await res.json()
      ElMessage.error(data.error || "退出失败")
    }
  } catch { ElMessage.error("网络错误") }
  finally { leaving.value = false }
}

onMounted(fetchDetail)
</script>

<style scoped lang="scss">
.uni-detail-page {
  max-width: 680px;
  margin: 0 auto;
  padding: 24px;
  width: 100%;
}

.loading-wrap {
  padding: 40px 0;
}

.error-state {
  padding: 60px 0;
}

.uni-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 0;
}

.uni-icon {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 26px;
  font-weight: 700;
  flex-shrink: 0;
}

.uni-meta {
  flex: 1;

  h1 {
    margin: 0;
    font-size: 22px;
    font-weight: 600;
  }

  .uni-domain {
    color: var(--text-secondary, #909399);
    font-size: 13px;
    margin: 4px 0;
  }

  .uni-count {
    color: var(--text-secondary, #909399);
    font-size: 13px;
    margin: 0;
  }
}

.section {
  margin-top: 8px;

  .section-title {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 16px;
  }
}

.members-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.member-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  border: 1px solid var(--border-color, #e4e7ed);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: var(--el-fill-color-light);
    transform: translateY(-1px);
  }
}

.member-info {
  flex: 1;
  min-width: 0;

  .member-name {
    font-size: 14px;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 6px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.empty-members {
  padding: 40px 0;
}
</style>
