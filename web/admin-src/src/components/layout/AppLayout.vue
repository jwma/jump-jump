<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="header-left">
        <h2 class="app-title">Jump Jump Admin</h2>
      </div>
      <div class="header-right">
        <el-dropdown trigger="click" @command="handleTenantSwitch">
          <span class="tenant-selector">
            <span v-if="userStore.currentTenant" class="tenant-selector-text">
              {{ userStore.currentTenant.tenantName }}
            </span>
            <span v-else-if="userStore.isSuper" class="tenant-selector-text">超管模式</span>
            <span v-else class="tenant-selector-text">未选择租户</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                v-for="tenant in userStore.tenants"
                :key="tenant.tenantId"
                :command="tenant.tenantId"
                :class="{ 'is-active': userStore.currentTenantId === tenant.tenantId }"
              >
                {{ tenant.tenantName }}
                <el-tag :type="tenant.role === 'admin' ? 'danger' : 'info'" size="small" class="tenant-role-tag">
                  {{ tenant.role === 'admin' ? '管理员' : '成员' }}
                </el-tag>
              </el-dropdown-item>
              <el-dropdown-item divided command="__manage__">
                <el-icon><Setting /></el-icon>
                管理租户
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <el-badge :value="userStore.pendingInvitationCount" :hidden="userStore.pendingInvitationCount === 0" :max="99">
          <el-button :icon="Bell" circle @click="router.push('/invitations')" />
        </el-badge>

        <span class="username">{{ userStore.authUser?.username }}</span>

        <el-button type="primary" link @click="router.push('/change-password')">修改密码</el-button>
        <el-button type="danger" link @click="handleLogout">退出登录</el-button>
      </div>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowDown, Setting, Bell } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

let refreshTimer: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  if (!userStore.authUser) {
    try {
      await userStore.fetchAuthInfo()
    } catch {
      router.push('/login')
      return
    }
  }
  if (userStore.currentTenantId && !userStore.userInfo) {
    userStore.fetchUserInfo().catch(() => {})
  }
  await userStore.fetchInvitations()

  refreshTimer = setInterval(() => {
    userStore.fetchInvitations().catch(() => {})
  }, 60000)
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})

function handleTenantSwitch(command: string) {
  if (command === '__manage__') {
    router.push('/select-tenant')
    return
  }
  if (command !== userStore.currentTenantId) {
    userStore.selectTenant(command)
    userStore.fetchUserInfo().catch(() => {})
    router.push('/')
  }
}

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: #f0f2f5;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.app-title {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.tenant-selector {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background 0.2s;
  user-select: none;
}

.tenant-selector:hover {
  background: #f5f7fa;
}

.tenant-selector-text {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.tenant-role-tag {
  margin-left: 8px;
}

.username {
  color: #606266;
  font-size: 14px;
}

.app-main {
  padding: 24px;
}
</style>
