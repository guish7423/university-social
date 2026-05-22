<template>
  <div class="campus-rooms">
    <div class="page-header">
      <h1>空教室</h1>
      <el-select v-model="selectedBuilding" placeholder="选择教学楼" clearable @change="handleFilter" style="width:200px">
        <el-option v-for="b in buildings" class="stagger-item" :key="b.id" :label="b.name" :value="b.name" />
      </el-select>
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!rooms.length" class="empty-state"><el-empty description="当前无空教室" /></div>
    <div v-else class="room-grid">
      <div v-for="room in rooms" class="room-card stagger-item" :key="room.id">
        <div class="room-building">{{ room.building }}</div>
        <div class="room-name">{{ room.room }}</div>
        <div class="room-capacity">{{ room.capacity }} 座</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listEmptyRooms, listBuildings } from "@/api/campus"
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
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
.campus-rooms { max-width: 900px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; h1 { font-size: 22px; font-weight: 700; } }
.loading-wrap, .empty-state { padding: 40px 0; }
.room-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px,1fr)); gap: 14px; }
.room-card {
  background: $bg-card; border: 1px solid $border-color; border-radius: $radius-md;
  padding: 20px; text-align: center; cursor: default;
  .room-building { font-size: 12px; color: $text-muted; margin-bottom: 6px; }
  .room-name { font-size: 20px; font-weight: 700; color: $text-primary; }
  .room-capacity { font-size: 12px; color: $text-secondary; margin-top: 6px; }
}
</style>
