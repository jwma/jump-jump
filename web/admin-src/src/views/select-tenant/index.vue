<template>
  <div class="select-tenant-container">
    <div class="select-tenant-card">
      <h2 class="select-tenant-title">选择租户</h2>
      <p class="select-tenant-subtitle">{{ userStore.authUser?.username }}，请选择要进入的租户</p>

      <div class="tenant-list">
        <div
          v-for="tenant in userStore.tenants"
          :key="tenant.tenantId"
          class="tenant-card"
          @click="selectTenant(tenant.tenantId)"
        >
          <div class="tenant-card-header">
            <span class="tenant-name">{{ tenant.tenantName }}</span>
            <el-tag :type="tenant.role === 'admin' ? 'danger' : 'info'" size="small">
              {{ tenant.role === 'admin' ? '管理员' : '成员' }}
            </el-tag>
          </div>
          <div class="tenant-card-footer">
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>

      <div v-if="userStore.isSuper" class="super-admin-entry" @click="enterSuperAdmin">
        <el-icon><Setting /></el-icon>
        <span>进入超管后台</span>
        <el-icon><ArrowRight /></el-icon>
      </div>

      <el-button class="create-tenant-btn" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        创建新租户
      </el-button>

      <div class="logout-row">
        <el-button type="info" link @click="handleLogout">退出登录</el-button>
      </div>
    </div>

    <el-dialog v-model="showCreateDialog" title="创建租户" width="420px" :close-on-click-modal="false">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-position="top">
        <el-form-item label="租户名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入租户名称" />
        </el-form-item>
        <el-form-item label="租户标识" prop="slug">
          <el-input v-model="createForm.slug" placeholder="请输入租户标识（用于 URL）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="handleCreateTenant">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, ArrowRight, Setting } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { createTenant } from '@/api/tenant'

const router = useRouter()
const userStore = useUserStore()
const showCreateDialog = ref(false)
const createLoading = ref(false)
const createFormRef = ref<FormInstance>()

const createForm = reactive({
  name: '',
  slug: '',
})

const createRules: FormRules = {
  name: [{ required: true, message: '请输入租户名称', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入租户标识', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '只允许小写字母、数字和连字符', trigger: 'blur' },
  ],
}

onMounted(async () => {
  if (!userStore.authUser) {
    try {
      await userStore.fetchAuthInfo()
    } catch {
      router.push('/login')
    }
  }
})

function selectTenant(tenantId: string) {
  userStore.selectTenant(tenantId)
  router.push('/')
}

function enterSuperAdmin() {
  userStore.clearTenant()
  router.push('/')
}

async function handleCreateTenant() {
  const valid = await createFormRef.value?.validate().catch(() => false)
  if (!valid) return

  createLoading.value = true
  try {
    const res = await createTenant({ name: createForm.name, slug: createForm.slug })
    await userStore.fetchAuthInfo()
    showCreateDialog.value = false
    createForm.name = ''
    createForm.slug = ''
    ElMessage.success('租户创建成功')
    userStore.selectTenant(res.data.id)
    router.push('/')
  } catch {
    // Error handled by interceptor
  } finally {
    createLoading.value = false
  }
}

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.select-tenant-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.select-tenant-card {
  width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  padding: 40px 36px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.select-tenant-title {
  text-align: center;
  margin: 0 0 4px;
  font-size: 24px;
  color: #303133;
  font-weight: 600;
}

.select-tenant-subtitle {
  text-align: center;
  margin: 0 0 24px;
  color: #909399;
  font-size: 14px;
}

.tenant-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.tenant-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.tenant-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.15);
}

.tenant-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tenant-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.tenant-card-footer {
  color: #c0c4cc;
}

.super-admin-entry {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
  border: 1px dashed #e6a23c;
  border-radius: 8px;
  cursor: pointer;
  color: #e6a23c;
  font-weight: 500;
  margin-bottom: 16px;
  transition: all 0.2s;
}

.super-admin-entry:hover {
  background: #fdf6ec;
}

.super-admin-entry .el-icon:last-child {
  margin-left: auto;
}

.create-tenant-btn {
  width: 100%;
  margin-bottom: 12px;
}

.logout-row {
  text-align: center;
}
</style>
