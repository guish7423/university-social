<template>
  <div class="create-page">
    <PageHeader title="发起评价" @back="$router.back()" />

    <div class="form-wrap">
      <el-form @submit.prevent="handleCreate" label-position="top">
        <el-form-item label="评价对象名称" required>
          <el-input v-model="name" placeholder="输入你想评价的对象（课程、老师、食堂菜品...）" maxlength="50" show-word-limit
            @input="handleSearch" clearable>
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <div v-if="searchResults.length" class="search-dropdown">
            <div v-for="t in searchResults" :key="t.id" class="search-item" @click="selectExisting(t)">
              <span class="si-name">{{ t.name }}</span>
              <span class="si-meta">{{ categoryLabel(t.category) }} · {{ t.totalRatings }} 评</span>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="分类" required>
          <el-select v-model="category" placeholder="选择分类" style="width:100%">
            <el-option label="课程" value="course" />
            <el-option label="教师" value="teacher" />
            <el-option label="商品" value="product" />
            <el-option label="活动" value="activity" />
            <el-option label="生活" value="life" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="描述（选填）">
          <el-input v-model="description" type="textarea" :rows="3" maxlength="200" show-word-limit placeholder="简单介绍一下这个对象..." />
        </el-form-item>

        <el-button type="primary" size="large" :loading="submitting" :disabled="!canSubmit" @click="handleCreate" style="width:100%">
          创建评价对象
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue"
import { useRouter } from "vue-router"
import { createReviewTarget, searchReviewTargets } from "@/api/review"
import type { ReviewTargetData } from "@/api/review"
import PageHeader from "@/components/PageHeader.vue"
import { Search } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"

const router = useRouter()
const name = ref("")
const category = ref("")
const description = ref("")
const submitting = ref(false)
const searchResults = ref<ReviewTargetData[]>([])
let searchTimer: ReturnType<typeof setTimeout> | null = null

const canSubmit = computed(() => name.value.trim() && category.value)

function categoryLabel(cat: string) {
  const map: Record<string, string> = { course: "课程", teacher: "教师", product: "商品", activity: "活动", life: "生活", other: "其他" }
  return map[cat] || cat
}

function selectExisting(t: ReviewTargetData) {
  router.push("/reviews/" + t.id)
}

function handleSearch() {
  if (searchTimer) clearTimeout(searchTimer)
  const q = name.value.trim()
  if (!q || q.length < 2) { searchResults.value = []; return }
  searchTimer = setTimeout(async () => {
    try {
      const res = await searchReviewTargets(q)
      searchResults.value = (res.targets || []).slice(0, 5)
    } catch { searchResults.value = [] }
  }, 300)
}

async function handleCreate() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    const target = await createReviewTarget({
      name: name.value.trim(),
      category: category.value,
      description: description.value.trim() || undefined,
    })
    ElMessage.success("创建成功")
    router.push("/reviews/" + target.id)
  } catch { /* */ } finally { submitting.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.create-page { max-width: 600px; }

.form-wrap {
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: $space-6;
}

.search-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  z-index: 10;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-sm;
  margin-top: 4px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
  overflow: hidden;
}

.search-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $space-3;
  cursor: pointer;
  transition: background 0.15s;
  &:hover { background: $bg-hover; }
  &:not(:last-child) { border-bottom: 1px solid $border-color; }
}

.si-name {
  font-size: $text-sm;
  font-weight: 600;
  color: $text-primary;
}

.si-meta {
  font-size: $text-xs;
  color: $text-secondary;
}
</style>
