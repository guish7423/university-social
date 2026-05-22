<template>
  <div class="create-page">
    <h1>发起活动</h1>
    <el-form @submit.prevent="handleSubmit">
      <el-form-item label="活动标题">
        <el-input v-model="title" placeholder="活动名称" maxlength="60" />
      </el-form-item>
      <el-form-item label="活动类型">
        <el-input v-model="activityType" placeholder="如：运动、学习、娱乐等" />
      </el-form-item>
      <el-form-item label="活动地点">
        <el-input v-model="location" placeholder="活动地点" />
      </el-form-item>
      <el-form-item label="开始时间">
        <el-date-picker v-model="startTime" type="datetime" placeholder="选择开始时间" />
      </el-form-item>
      <el-form-item label="人数上限">
        <el-input-number v-model="maxParticipants" :min="2" :max="200" />
      </el-form-item>
      <el-form-item label="活动描述">
        <el-input v-model="description" type="textarea" :rows="4" placeholder="描述你的活动..." />
      </el-form-item>
      <div class="form-actions">
        <el-button @click="$router.back()">取消</el-button>
        <el-button type="primary" :loading="submitting" :disabled="!title.trim()" native-type="submit">
          发起
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { createActivity } from "@/api/activity"
import dayjs from "dayjs"

const router = useRouter()
const title = ref("")
const activityType = ref("")
const location = ref("")
const startTime = ref<Date>(new Date())
const maxParticipants = ref(20)
const description = ref("")
const submitting = ref(false)

async function handleSubmit() {
  if (!title.value.trim()) return
  submitting.value = true
  try {
    const res = await createActivity({
      title: title.value.trim(),
      activity_type: activityType.value.trim() || "其他",
      location: location.value.trim(),
      start_time: dayjs(startTime.value).format("YYYY-MM-DD HH:mm:ss"),
      max_participants: maxParticipants.value,
      description: description.value.trim(),
    })
    router.push("/activities/" + res.id)
  } finally { submitting.value = false }
}
</script>

<style scoped lang="scss">
.create-page { max-width: 560px; }
h1 { font-size: 20px; font-weight: 700; margin-bottom: 20px; }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
</style>
