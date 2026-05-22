<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { listBuildings, listEmptyRooms, type EmptyRoom } from "@/api/campus"

const buildings = ref<string[]>([])
const rooms = ref<EmptyRoom[]>([])
const loading = ref(true)
const selectedBuilding = ref("")
const weekdays = ["周一", "周二", "周三", "周四", "周五", "周六", "周日"]
const selectedWeekday = ref(1)

const grouped = computed(() => {
  const g: Record<string, EmptyRoom[]> = {}
  for (const r of rooms.value) {
    if (!g[r.building]) g[r.building] = []
    g[r.building].push(r)
  }
  return g
})

onMounted(async () => {
  try {
    buildings.value = await listBuildings()
    if (buildings.value.length) selectedBuilding.value = buildings.value[0]
    await search()
  } catch {}
  loading.value = false
})

async function search() {
  loading.value = true
  try {
    rooms.value = await listEmptyRooms(selectedBuilding.value || undefined, selectedWeekday.value)
  } catch { rooms.value = [] }
  loading.value = false
}

function floorSort(a: string, b: string) {
  return parseInt(a) - parseInt(b)
}
</script>

<template>
  <view class="page">
    <view class="filters">
      <scroll-view scroll-x class="filter-scroll" show-scrollbar="false">
        <view v-for="b in buildings" :key="b" :class="['filter-tag', selectedBuilding === b && 'filter-active']" @click="selectedBuilding = b; search()">
          <text>{{ b }}</text>
        </view>
      </scroll-view>
      <scroll-view scroll-x class="filter-scroll" show-scrollbar="false">
        <view v-for="(wd, i) in weekdays" :key="i" :class="['filter-tag', selectedWeekday === i + 1 && 'filter-active']" @click="selectedWeekday = i + 1; search()">
          <text>{{ wd }}</text>
        </view>
      </scroll-view>
    </view>
    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
    </view>
    <view v-else-if="rooms.length === 0" class="empty-state">
      <u-icon name="map" size="120" color="#EAE6E0" />
      <text class="empty-text">该时段暂无空教室</text>
    </view>
    <view v-else class="list">
      <view v-for="(rs, building) in grouped" :key="building" class="building-section">
        <text class="building-title">{{ building }}</text>
        <view class="room-grid">
          <view v-for="r in rs" :key="r.id" class="room-card">
            <text class="room-name">{{ r.room }}</text>
            <text class="room-capacity">{{ r.capacity }}人</text>
            <text class="room-time">{{ r.start_time.slice(0, 5) }}-{{ r.end_time.slice(0, 5) }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: #F7F4F0; padding: 24rpx; }
.loading-state, .empty-state { display: flex; flex-direction: column; align-items: center; padding: 200rpx 0; gap: 16rpx; }
.empty-text { font-size: 26rpx; color: #B8C2CE; }
.filters { display: flex; flex-direction: column; gap: 12rpx; margin-bottom: 20rpx; }
.filter-scroll { white-space: nowrap; }
.filter-tag {
  display: inline-flex; padding: 8rpx 24rpx; margin-right: 12rpx;
  border-radius: 20rpx; background: #fff; font-size: 24rpx; color: #5C6B7E;
  border: 1rpx solid #E0DBD4; cursor: pointer;
}
.filter-active { background: #C67A6A; color: #fff; border-color: #C67A6A; }
.building-section { margin-bottom: 20rpx; }
.building-title {
  font-size: 26rpx; font-weight: 700; color: #1E2A3A;
  padding-left: 8rpx; margin-bottom: 12rpx; border-left: 6rpx solid #C67A6A;
}
.room-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12rpx; }
.room-card {
  background: #fff; border-radius: 12rpx; padding: 16rpx;
  text-align: center; box-shadow: 0 2rpx 8rpx rgba(30,42,58,0.04);
}
.room-name { font-size: 26rpx; font-weight: 700; color: #1E2A3A; display: block; }
.room-capacity { font-size: 20rpx; color: #8E9BAB; margin-top: 4rpx; display: block; }
.room-time { font-size: 18rpx; color: #2ECC71; margin-top: 4rpx; display: block; }
</style>
