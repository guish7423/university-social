<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { listDirectory, type DirectoryEntry } from "@/api/campus"

const entries = ref<DirectoryEntry[]>([])
const loading = ref(true)
const search = ref("")
const selectedDept = ref("")

const departments = computed(() => {
  const depts = new Set(entries.value.map(e => e.department))
  return ["全部", ...Array.from(depts)]
})

const filtered = computed(() => {
  return entries.value.filter(e => {
    if (selectedDept.value && selectedDept.value !== "全部" && e.department !== selectedDept.value) return false
    if (search.value) {
      const q = search.value.toLowerCase()
      return e.name.includes(q) || e.department.includes(q) || (e.position && e.position.includes(q))
    }
    return true
  })
})

onMounted(async () => {
  try { entries.value = await listDirectory() } catch {}
  loading.value = false
})

function callPhone(phone: string) {
  uni.makePhoneCall({ phoneNumber: phone })
}
</script>

<template>
  <view class="page">
    <view class="search-bar">
      <u-search v-model="search" placeholder="搜索姓名/部门/职位" :show-action="false" />
    </view>
    <scroll-view scroll-x class="dept-scroll" show-scrollbar="false">
      <view v-for="d in departments" :key="d" :class="['dept-tag', selectedDept === d && 'dept-active']" @click="selectedDept = d">
        <text>{{ d === '全部' ? '全部' : d }}</text>
      </view>
    </scroll-view>
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
    </view>
    <view v-else-if="filtered.length === 0" class="empty-state">
      <u-icon name="search" size="120" color="#EAE6E0" />
      <text class="empty-text">未找到</text>
    </view>
    <view v-else class="list">
      <view v-for="e in filtered" :key="e.id" class="entry-card">
        <view class="entry-top">
          <text class="entry-name">{{ e.name }}</text>
          <text class="entry-position">{{ e.position || '' }}</text>
        </view>
        <text class="entry-dept">{{ e.department }}</text>
        <view class="entry-contact">
          <view v-if="e.phone" class="contact-item" @click="callPhone(e.phone)">
            <u-icon name="phone" size="24" color="#C67A6A" />
            <text class="contact-text">{{ e.phone }}</text>
          </view>
          <text v-if="e.office" class="contact-item">
            <text class="office-icon">📍</text>
            <text class="contact-text">{{ e.office }}</text>
          </text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: #F7F4F0; padding: 24rpx; }
.search-bar { margin-bottom: 16rpx; }
.loading-state, .empty-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 0; gap: 16rpx; }
.empty-text { font-size: 26rpx; color: #B8C2CE; }
.dept-scroll { white-space: nowrap; margin-bottom: 16rpx; }
.dept-tag {
  display: inline-flex; padding: 8rpx 24rpx; margin-right: 12rpx;
  border-radius: 20rpx; background: #fff; font-size: 24rpx; color: #5C6B7E;
  border: 1rpx solid #E0DBD4; cursor: pointer;
}
.dept-active { background: #C67A6A; color: #fff; border-color: #C67A6A; }
.list { display: flex; flex-direction: column; gap: 12rpx; }
.entry-card { background: #fff; border-radius: 14rpx; padding: 24rpx; box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04); }
.entry-top { display: flex; align-items: center; gap: 12rpx; margin-bottom: 8rpx; }
.entry-name { font-size: 28rpx; font-weight: 700; color: #1E2A3A; }
.entry-position { font-size: 22rpx; color: #8E9BAB; }
.entry-dept { font-size: 22rpx; color: #C67A6A; margin-bottom: 12rpx; display: block; }
.entry-contact { display: flex; gap: 24rpx; }
.contact-item { display: flex; align-items: center; gap: 6rpx; cursor: pointer; }
.contact-text { font-size: 22rpx; color: #5C6B7E; }
</style>
