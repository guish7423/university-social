<template>
  <div class="edit-page">
    <PageHeader title="编辑资料" description="修改你的个人信息" />

    <div class="edit-card stagger-item">
      <!-- Avatar preview area -->
      <div class="avatar-section">
        <div class="avatar-wrap" @click="triggerUpload">
          <el-avatar :size="72" :src="avatar" class="avatar-img" />
          <div class="avatar-overlay">
            <el-icon :size="20"><Camera /></el-icon>
            <span>更换头像</span>
          </div>
        </div>
      </div>

      <el-form @submit.prevent="handleSubmit" label-position="top">
        <el-form-item label="昵称">
          <el-input v-model="nickname" maxlength="20" placeholder="输入昵称" />
        </el-form-item>
        <el-form-item label="学校">
          <el-input v-model="school" placeholder="如：洛阳师范学院" />
        </el-form-item>
        <el-form-item label="个人简介">
          <el-input
            v-model="bio"
            type="textarea"
            :maxlength="200"
            :rows="3"
            placeholder="介绍一下你自己…"
            show-word-limit
          />
        </el-form-item>
        <div class="form-actions">
          <el-button @click="$router.back()">取消</el-button>
          <el-button type="primary" :loading="submitting" native-type="submit">保存</el-button>
        </div>
      </el-form>
    </div>

    <!-- Hidden file input for avatar (UI only, not implementing real upload yet) -->
    <input ref="fileInput" type="file" accept="image/*" style="display:none" @change="handleFileChange" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getProfile, updateProfile } from "@/api/auth"
import PageHeader from "@/components/PageHeader.vue"
import { Camera } from "@element-plus/icons-vue"

const router = useRouter()
const userStore = useUserStore()
const nickname = ref("")
const school = ref("")
const bio = ref("")
const avatar = ref("")
const submitting = ref(false)
const fileInput = ref<HTMLInputElement>()

function triggerUpload() {
  fileInput.value?.click()
}

function handleFileChange(e: Event) {
  // UI-only preview — real upload not implemented
  const input = e.target as HTMLInputElement
  if (input.files?.length) {
    const url = URL.createObjectURL(input.files[0])
    avatar.value = url
  }
}

async function handleSubmit() {
  submitting.value = true
  try {
    await updateProfile({
      nickname: nickname.value.trim(),
      school: school.value.trim(),
      bio: bio.value.trim(),
    } as any)
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
    bio.value = (profile as any).bio || ""
    avatar.value = profile.avatar
    userStore.setUser(profile)
  } catch { /* handled by interceptor */ }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.edit-page { max-width: 520px; }

.edit-card {
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: 24px;

  .avatar-section {
    display: flex;
    justify-content: center;
    margin-bottom: 24px;

    .avatar-wrap {
      position: relative;
      cursor: pointer;
      border-radius: 50%;
      overflow: hidden;

      .avatar-img {
        display: block;
      }

      .avatar-overlay {
        position: absolute;
        inset: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        background: rgba(0,0,0,0.5);
        color: #fff;
        font-size: 12px;
        gap: 4px;
        opacity: 0;
        transition: opacity 0.2s;
      }

      &:hover .avatar-overlay {
        opacity: 1;
      }
    }
  }

  .page-header {
    margin-bottom: 20px;
    h1 { font-size: 22px; font-weight: 700; }
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 16px;
  }
}
</style>
