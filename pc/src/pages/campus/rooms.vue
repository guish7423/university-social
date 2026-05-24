<template>
  <div class="campus-rooms">
<PageHeader title="空教室">
  <template #extra>
    <el-select v-model="selectedBuilding" placeholder="选择教学楼" clearable @change="handleFilter" style="width:200px">
      <el-option v-for="b in buildings" class="stagger-item" :key="b.id" :label="b.name" :value="b.name" />
    </el-select>
  </template>
</PageHeader>

    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!rooms.length" class="empty-state"><el-empty description="当前无空教室" /></div>
    <div v-else class="room-grid">
      <div v-for="room in rooms" class="room-card stagger-item" :key="room.id">
        <div class="room-building">{{ room.building }}</div>
        <div class="room-name">{{ room.room }}</div>
        <div class="room-capacity">
          <span class="capacity-badge" :class="capacityLevel(room.capacity)">{{ room.capacity }} 座</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listEmptyRooms, listBuildings } from "@/api/campus"
import PageHeader from "@/components/PageHeader.vue"

import type { EmptyRoom, Building } from "@/api/campus"

const rooms = ref<EmptyRoom[]>([])
const buildings = ref<Building[]>([])
const selectedBuilding = ref("")
const loading = ref(true)

async function loadRooms(building?: string) {
  loading.value = true
  try { rooms.value = await listEmptyRooms(building) }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }
}

function handleFilter() { loadRooms(selectedBuilding.value || undefined) }

onMounted(async () => {
  try { buildings.value = await listBuildings() } catch { /* handled by interceptor */ }
  await loadRooms()
})

function capacityLevel(c: number) {
  if (c >= 100) return "high"
  if (c >= 50) return "medium"
  return "low"
}
</script>

<style scoped lang="scss">
@use "sass:color";
@use "@/styles/variables.scss" as *;
 margin: 0 auto;
.campus-rooms { max-width: 900px; }
.loading-wrap, .empty-state { padding: 40px 0; }
.room-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px,1fr)); gap: 14px; }
.room-card {
  background: $bg-card; border: 1px solid $border-color; border-radius: $radius-md;
  padding: 20px; text-align: center; cursor: default;
  .room-building { font-size: 12px; color: $text-muted; margin-bottom: 6px; }
  .room-name { font-size: 20px; font-weight: 700; color: $text-primary; }
  .room-capacity { margin-top: 8px; }
  .capacity-badge {
    display: inline-block; padding: 2px 12px; border-radius: 999px;
    font-size: 13px; font-weight: 600;
    &.high { background: color.change($color-success, $alpha: 0.15); color: $color-success; }
    &.medium { background: color.change($color-warning, $alpha: 0.15); color: $color-warning; }
    &.low { background: color.change($color-danger, $alpha: 0.15); color: $color-danger; }
  }
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0,0,0,0.15);
  }
  transition: transform 0.25s ease, box-shadow 0.25s ease;
}
</style>

.stagger-item {
  opacity: 0;
  transform: translateY(12px);
  animation: staggerFadeIn 0.4s ease forwards;

  @for $i from 1 through 30 {
    &:nth-child(#{$i}) { animation-delay: #{$i * 0.05}s; }
  }
}

@keyframes staggerFadeIn {
  to { opacity: 1; transform: translateY(0); }
}
