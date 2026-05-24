<template>
  <div class="campus-directory">
    <div class="page-header">
      <h1>通讯录</h1>
      <el-input
      v-model="searchQ"
      placeholder="搜索姓名或部门"
      clearable
      style="width:240px"
>
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>
    </div>

    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!filtered.length" class="empty-state"><el-empty :description="searchQ ? '未找到匹配的联系人' : '暂无可用的联系信息'" /></div>

    <div v-else class="contacts-grid">
      <div v-for="c in filtered" :key="c.id" class="contact-card stagger-item">
        <div class="card-top">
          <div class="contact-avatar" :style="{ background: avatarGradient(c.name) }">
            {{ nameInitial(c.name) }}
          </div>
          <div class="contact-info">
            <div class="contact-name">{{ c.name }}</div>
            <el-tag size="small" :type="deptTagType(c.department)">{{ c.department }}</el-tag>
          </div>
        </div>
        <div class="card-details">
          <div class="detail-row">
            <span class="detail-icon">📞</span>
            <span class="detail-text">{{ c.phone }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">✉️</span>
            <span class="detail-text">{{ c.email }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">📍</span>
            <span class="detail-text">{{ c.office }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { Search } from "@element-plus/icons-vue"
import { listDirectory } from "@/api/campus"
import type { CampusContact } from "@/api/campus"

const contacts = ref<CampusContact[]>([])
const loading = ref(true)
const searchQ = ref("")

const filtered = computed(() => {
  const q = searchQ.value.toLowerCase().trim()
  if (!q) return contacts.value
  return contacts.value.filter(c =>
    c.name.toLowerCase().includes(q) || c.department.toLowerCase().includes(q)
  )
})

// function handleSearch() {}
function nameInitial(name: string): string {
  return name.charAt(0).toUpperCase()
}

const avatarColors = [
  "linear-gradient(135deg, #e74c3c, #c0392b)",
  "linear-gradient(135deg, #e67e22, #d35400)",
  "linear-gradient(135deg, #2ecc71, #27ae60)",
  "linear-gradient(135deg, #3498db, #2980b9)",
  "linear-gradient(135deg, #9b59b6, #8e44ad)",
  "linear-gradient(135deg, #1abc9c, #16a085)",
  "linear-gradient(135deg, #e84393, #c2185b)",
  "linear-gradient(135deg, #00cec9, #00b894)",
]

function avatarGradient(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = ((hash << 5) - hash) + name.charCodeAt(i)
  }
  return avatarColors[Math.abs(hash) % avatarColors.length]
}

const deptTags = ["教学", "行政", "学生", "后勤", "科研", "信息"]
function deptTagType(dept: string): string {
  if (dept.includes("教学") || dept.includes("教务")) return "primary"
  if (dept.includes("行政") || dept.includes("办公")) return ""
  if (dept.includes("学生") || dept.includes("团委")) return "success"
  if (dept.includes("后勤") || dept.includes("物业")) return "warning"
  if (dept.includes("科研") || dept.includes("研究")) return "info"
  return "primary"
}

onMounted(async () => {
  try { contacts.value = await listDirectory() }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.campus-directory {
  max-width: 760px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  h1 { font-size: 22px; font-weight: 700; }
}

.loading-wrap, .empty-state { padding: 60px 0; }

.contacts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.contact-card {
  background: $bg-card;
  border: 1px solid $border-light;
  border-radius: $radius-md;
  padding: 20px;
  transition: all 0.25s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
    border-color: $brand-primary;
  }
}

.card-top {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 16px;
  padding-bottom: 14px;
  border-bottom: 1px solid $border-light;
}

.contact-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.contact-info {
  flex: 1;
  min-width: 0;
}

.contact-name {
  font-size: 16px;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: 4px;
}

.card-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: $text-secondary;
}

.detail-icon {
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.detail-text {
  word-break: break-all;
}
</style>
