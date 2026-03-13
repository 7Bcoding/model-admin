<template>
  <div class="page-container">
    <div class="header">
      <div class="header-left">
        <div class="title-section">
          <i class="el-icon-cpu title-icon"></i>
          <h2 class="page-title">Fusion 管理</h2>
          <el-select 
            v-model="selectedBackend" 
            @change="handleBackendChange"
            placeholder="选择后端"
            style="width: 100px; margin-left: 15px;"
          >
            <el-option
              v-for="option in backendOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            >
              <div class="backend-option">
                <span class="backend-name">{{ option.label }}</span>
                <span class="backend-desc">{{ option.description }}</span>
              </div>
            </el-option>
          </el-select>
        </div>
      </div>
    </div>

    <div class="content-tabs">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="模型管理" name="models">
          <FusionModels />
        </el-tab-pane>
        <el-tab-pane label="API 账号管理" name="accounts">
          <FusionAccounts />
        </el-tab-pane>
        <el-tab-pane label="供应商账号管理" name="vendors">
          <FusionVendors />
        </el-tab-pane>
        <el-tab-pane label="模型价格管理" name="model-prices">
          <FusionModelPrices />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useFusionStore } from '@/stores/fusion'
import FusionAccounts from './FusionAccounts.vue'
import FusionModels from './FusionModels.vue'
import FusionVendors from './FusionVendors.vue'
import FusionModelPrices from './FusionModelPrices.vue'

const route = useRoute()
const router = useRouter()
const fusionStore = useFusionStore()

// Initialize activeTab from URL query parameter or default to 'models'
const activeTab = ref(route.query.tab || 'models')

const selectedBackend = computed({
  get: () => fusionStore.selectedBackend,
  set: (value) => fusionStore.setBackend(value)
})

const backendOptions = computed(() => fusionStore.getBackendOptions())

const handleBackendChange = (backend) => {
  console.log('Backend change triggered:', backend)
  fusionStore.setBackend(backend)
  // Update URL with backend parameter
  updateUrlParams()
  ElMessage.success(`已切换到 ${backend === 'beta' ? 'beta' : 'alpha'} 后端`)
}

// Watch for activeTab changes to update URL
watch(activeTab, (newTab) => {
  updateUrlParams()
})

const updateUrlParams = () => {
  const query = {
    tab: activeTab.value,
    backend: fusionStore.selectedBackend
  }
  router.push({ query })
}

// Watch for route changes to sync state
watch(() => route.query, (newQuery) => {
  if (newQuery.tab && newQuery.tab !== activeTab.value) {
    activeTab.value = newQuery.tab
  }
  if (newQuery.backend && newQuery.backend !== fusionStore.selectedBackend) {
    fusionStore.setBackend(newQuery.backend)
  }
}, { immediate: true })

onMounted(() => {
  // Initialize from URL parameters or defaults
  const urlBackend = route.query.backend
  const urlTab = route.query.tab
  
  if (urlBackend && ['beta', 'alpha'].includes(urlBackend)) {
    fusionStore.setBackend(urlBackend)
  } else if (!fusionStore.selectedBackend) {
    fusionStore.setBackend('beta')
  }
  
  if (urlTab && ['models', 'accounts', 'vendors', 'model-prices'].includes(urlTab)) {
    activeTab.value = urlTab
  }
  
  // Update URL if no parameters are set
  if (!route.query.tab || !route.query.backend) {
    updateUrlParams()
  }
})
</script>

<style scoped>
.page-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.title-section {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  font-size: 24px;
  color: #409eff;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.content-tabs {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.backend-option {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.backend-name {
  font-weight: 500;
  color: #303133;
}

.backend-desc {
  font-size: 12px;
  color: #909399;
}

/* Fix tab padding consistency and make font larger */
.content-tabs :deep(.el-tabs__item) {
  padding: 0 20px !important;
  font-size: 18px !important;
  font-weight: 500 !important;
}
</style>