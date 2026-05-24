<template>
  <div class="security-page">
    <PageHeader title="账号安全" />

    <div class="card form-card">
      <h3 class="card-title">修改密码</h3>
      <p class="card-desc">为你的账号设置一个密码，增强安全性</p>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        class="password-form"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="旧密码" prop="old_password" v-if="hasPassword">
          <el-input
            v-model="form.old_password"
            type="password"
            placeholder="输入当前密码"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item label="新密码" prop="new_password">
          <el-input
            v-model="form.new_password"
            type="password"
            placeholder="至少6位密码"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input
            v-model="form.confirm_password"
            type="password"
            placeholder="再次输入新密码"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item class="form-actions">
          <el-button
            type="primary"
            native-type="submit"
            :loading="submitting"
            class="submit-btn"
          >
            {{ hasPassword ? '更新密码' : '设置密码' }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-alert
        v-if="successMsg"
        type="success"
        :title="successMsg"
        show-icon
        closable
        class="result-alert"
      />
      <el-alert
        v-if="errorMsg"
        type="error"
        :title="errorMsg"
        show-icon
        closable
        class="result-alert"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import { Lock } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"
import PageHeader from "@/components/PageHeader.vue"
import { getProfile, changePassword } from "@/api/auth"

interface FormType {
  old_password: string
  new_password: string
  confirm_password: string
}

const formRef = ref()
const submitting = ref(false)
const hasPassword = ref(false)
const successMsg = ref("")
const errorMsg = ref("")

const form = reactive<FormType>({
  old_password: "",
  new_password: "",
  confirm_password: "",
})

const validateConfirm = (_rule: any, value: string, callback: any) => {
  if (value !== form.new_password) {
    callback(new Error("两次密码不一致"))
  } else {
    callback()
  }
}

const rules = {
  old_password: [
    { required: true, message: "请输入旧密码", trigger: "blur" },
  ],
  new_password: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, message: "密码至少6位", trigger: "blur" },
  ],
  confirm_password: [
    { required: true, message: "请确认新密码", trigger: "blur" },
    { validator: validateConfirm, trigger: "blur" },
  ],
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  successMsg.value = ""
  errorMsg.value = ""

  try {
    const res = await changePassword({
      old_password: hasPassword.value ? form.old_password : "",
      new_password: form.new_password,
    })
    successMsg.value = "密码修改成功，下次登录请使用新密码"
    form.old_password = ""
    form.new_password = ""
    form.confirm_password = ""
    formRef.value?.resetFields()
    hasPassword.value = true
  } catch (err: any) {
    errorMsg.value = err?.response?.data?.error || "修改失败，请重试"
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const user = await getProfile()
    // If user has no password_hash, old_password is optional
    // We'll detect by checking if the user has set a password before
    hasPassword.value = false // default to no password for dev-login users
  } catch {}
})
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.security-page {
  max-width: 520px;
  margin: 0 auto;
  padding: $space-6 $space-4;
}

.form-card {
  background: var(--bg-card);
  border: 1px solid var(--border-light);
  border-radius: $radius-lg;
  padding: $space-8;
}

.card-title {
  font-size: $text-lg;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 $space-1;
}

.card-desc {
  font-size: $text-sm;
  color: var(--text-muted);
  margin: 0 0 $space-6;
}

.password-form {
  :deep(.el-form-item__label) {
    font-weight: 500;
    color: var(--text-secondary);
    padding-bottom: $space-1;
  }
}

.submit-btn {
  width: 100%;
  background: linear-gradient(135deg, $brand-primary, $brand-primary-light) !important;
  border: none !important;
  color: #fff !important;
  border-radius: $radius-md;
  font-weight: 500;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(198, 122, 106, 0.4);
  }
}

.result-alert {
  margin-top: $space-4;
}
</style>
