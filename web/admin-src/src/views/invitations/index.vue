<template>
  <div class="invitations-container">
    <div class="invitations-card">
      <h2 class="invitations-title">待处理邀请</h2>

      <div v-if="loading" class="loading-state">
        <el-icon class="is-loading" :size="24"><Loading /></el-icon>
      </div>

      <el-empty v-else-if="invitations.length === 0" description="暂无待处理邀请" />

      <div v-else class="invitation-list">
        <div v-for="inv in invitations" :key="inv.id" class="invitation-item">
          <div class="invitation-info">
            <span class="invitation-tenant">{{ inv.tenantName }}</span>
            <span class="invitation-detail">
              {{ inv.inviterUsername }} 邀请你加入
            </span>
            <span class="invitation-time">{{ formatTime(inv.createdAt) }}</span>
          </div>
          <div class="invitation-actions">
            <template v-if="inv.status === 'pending'">
              <el-button type="primary" size="small" :loading="actionLoading[inv.id]" @click="handleAccept(inv.id)">
                接受
              </el-button>
              <el-button size="small" :loading="actionLoading[inv.id]" @click="handleReject(inv.id)">
                拒绝
              </el-button>
            </template>
            <el-tag v-else :type="inv.status === 'accepted' ? 'success' : 'info'" size="small">
              {{ inv.status === 'accepted' ? '已接受' : '已拒绝' }}
            </el-tag>
          </div>
        </div>
      </div>

      <div class="back-row">
        <el-button @click="router.back()">返回</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { listInvitations, acceptInvitation, rejectInvitation, type Invitation } from '@/api/invitation'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const invitations = ref<Invitation[]>([])
const actionLoading = reactive<Record<string, boolean>>({})

onMounted(async () => {
  loading.value = true
  try {
    const res = await listInvitations()
    invitations.value = res.data ?? []
  } catch {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
})

function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN')
}

async function handleAccept(id: string) {
  actionLoading[id] = true
  try {
    await acceptInvitation(id)
    const inv = invitations.value.find((i) => i.id === id)
    if (inv) inv.status = 'accepted'
    await userStore.fetchAuthInfo()
    await userStore.fetchInvitations()
    ElMessage.success('已接受邀请')
  } catch {
    // Error handled by interceptor
  } finally {
    actionLoading[id] = false
  }
}

async function handleReject(id: string) {
  actionLoading[id] = true
  try {
    await rejectInvitation(id)
    const inv = invitations.value.find((i) => i.id === id)
    if (inv) inv.status = 'rejected'
    await userStore.fetchInvitations()
    ElMessage.success('已拒绝邀请')
  } catch {
    // Error handled by interceptor
  } finally {
    actionLoading[id] = false
  }
}
</script>

<style scoped>
.invitations-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
}

.invitations-card {
  width: 560px;
  max-height: 80vh;
  overflow-y: auto;
  padding: 40px 36px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.invitations-title {
  text-align: center;
  margin: 0 0 24px;
  font-size: 22px;
  color: #303133;
  font-weight: 600;
}

.loading-state {
  text-align: center;
  padding: 40px 0;
  color: #409eff;
}

.invitation-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.invitation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
}

.invitation-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.invitation-tenant {
  font-weight: 500;
  font-size: 15px;
  color: #303133;
}

.invitation-detail {
  font-size: 13px;
  color: #606266;
}

.invitation-time {
  font-size: 12px;
  color: #909399;
}

.invitation-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.back-row {
  text-align: center;
  margin-top: 12px;
}
</style>
