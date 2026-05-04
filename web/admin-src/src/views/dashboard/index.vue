<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>Jump Jump Admin</h2>
      <div class="header-actions">
        <span class="username">{{ userStore.userInfo?.username }}</span>
        <el-button type="primary" link @click="router.push('/change-password')">修改密码</el-button>
        <el-button type="danger" link @click="handleLogout">退出登录</el-button>
      </div>
    </div>
    <div class="dashboard-content">
      <el-empty description="欢迎使用 Jump Jump 管理后台" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

onMounted(() => {
  if (!userStore.userInfo) {
    userStore.fetchUserInfo()
  }
})

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: #f0f2f5;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.dashboard-header h2 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  color: #606266;
  font-size: 14px;
}

.dashboard-content {
  padding: 24px;
}
</style>
