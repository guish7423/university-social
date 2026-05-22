<template>
  <div class="edit-page">
    <h1>编辑资料</h1>
    <el-form @submit.prevent="handleSubmit">
      <el-form-item label="昵称">
        <el-input v-model="nickname" maxlength="20" />
      </el-form-item>
      <el-form-item label="学校">
        <el-input v-model="school" placeholder="如：洛阳师范学院" />
      </el-form-item>
      <div class="form-actions">
        <el-button @click="$router.back()">取消</el-button>
        <el-button type="primary" :loading="submitting" native-type="submit">保存</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getProfile, updateProfile } from "@/api/auth"

const router = useRouter()
const userStore = useUserStore()
const nickname = ref("")
const school = ref("")
const submitting = ref(false)

async function handleSubmit() {
  submitting.value = true
  try {
    await updateProfile({ nickname: nickname.value.trim(), school: school.value.trim() })
    const updated = await getProfile()
    userStore.setUser(updated)
    router.back()
  } finally { submitting.value = false }
}

onMounted(async () => {
  try {
    const profile = await getProfile()
    nickname.value = profile.nickname
    school.value = profile.school || ""
    userStore.setUser(profile)
  } catch { /* handled by interceptor */ }
})
</script>

<style scoped lang="scss">
.edit-page { max-width: 480px; }
h1 { font-size: 20px; font-weight: 700; margin-bottom: 20px; }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
</style>
