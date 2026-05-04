<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">Jump Jump Admin</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" :prefix-icon="User" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password :prefix-icon="Lock" @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="login-button" @click="handleLogin">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { getRememberedUsername, setRememberedUsername, removeRememberedUsername } from '@/utils/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)

const form = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 50, message: '用户名长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' },
  ],
}

onMounted(() => {
  const remembered = getRememberedUsername()
  if (remembered) {
    form.username = remembered
    rememberMe.value = true
  }
})

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login({ username: form.username, password: form.password })

    if (rememberMe.value) {
      setRememberedUsername(form.username)
    } else {
      removeRememberedUsername()
    }

    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch {
    // Error handled by axios interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  padding: 40px 36px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.login-title {
  text-align: center;
  margin-bottom: 32px;
  font-size: 24px;
  color: #303133;
  font-weight: 600;
}

.login-button {
  width: 100%;
}
</style>
