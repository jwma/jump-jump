<template>
  <div class="change-pwd-container">
    <div class="change-pwd-card">
      <h2 class="change-pwd-title">修改密码</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleSubmit">
        <el-form-item label="旧密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入旧密码" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" placeholder="请输入新密码" show-password />
          <div class="strength-bar" v-if="form.newPassword">
            <div class="strength-segments">
              <div
                v-for="i in 4"
                :key="i"
                class="strength-segment"
                :class="{ active: i <= strengthLevel }"
                :style="{ backgroundColor: i <= strengthLevel ? strengthColor : '#e0e0e0' }"
              />
            </div>
            <span class="strength-text" :style="{ color: strengthColor }">{{ strengthLabel }}</span>
          </div>
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="请再次输入新密码" show-password @keyup.enter="handleSubmit" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="submit-button" @click="handleSubmit">
            确认修改
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { changePassword } from '@/api/user'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  password: '',
  newPassword: '',
  confirmPassword: '',
})

function getPasswordStrength(pwd: string): number {
  if (!pwd) return 0
  let score = 0
  if (pwd.length >= 6) score++
  if (pwd.length >= 10) score++
  if (/[A-Z]/.test(pwd) && /[a-z]/.test(pwd)) score++
  if (/\d/.test(pwd) && /[^A-Za-z0-9]/.test(pwd)) score++
  return Math.min(score, 4)
}

const strengthLevel = computed(() => getPasswordStrength(form.newPassword))

const strengthColor = computed(() => {
  const colors = ['', '#F56C6C', '#E6A23C', '#409EFF', '#67C23A']
  return colors[strengthLevel.value] || ''
})

const strengthLabel = computed(() => {
  const labels = ['', '弱', '一般', '较强', '强']
  return labels[strengthLevel.value] || ''
})

const validateConfirmPassword = (_rule: unknown, value: string, callback: (err?: Error) => void) => {
  if (value !== form.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  password: [
    { required: true, message: '请输入旧密码', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await changePassword({ password: form.password, newPassword: form.newPassword })
    ElMessage.success('密码修改成功，请重新登录')
    userStore.resetState()
    router.push('/login')
  } catch {
    // Error handled by axios interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.change-pwd-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 60px;
}

.change-pwd-card {
  width: 440px;
  padding: 40px 36px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.change-pwd-title {
  text-align: center;
  margin-bottom: 32px;
  font-size: 22px;
  color: #303133;
  font-weight: 600;
}

.submit-button {
  margin-right: 12px;
}

.strength-bar {
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.strength-segments {
  display: flex;
  gap: 4px;
  flex: 1;
}

.strength-segment {
  height: 4px;
  flex: 1;
  border-radius: 2px;
  background-color: #e0e0e0;
  transition: background-color 0.3s;
}

.strength-segment.active {
  transition: background-color 0.3s;
}

.strength-text {
  font-size: 12px;
  white-space: nowrap;
}
</style>
