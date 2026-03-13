<template>
  <div class="ndeployment-container">
    <!-- 数据总结 -->
    <div class="data-summary">
      <div class="summary-item">
        <span class="summary-label">当前显示：</span>
        <span class="summary-value">{{ filteredNDeployments.length }}</span>
        <span class="summary-unit">条数据</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">热副本数：</span>
        <span class="summary-value">{{ deploymentSummary.hotReplicas }}</span>
        <span class="summary-unit">个</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">冷副本数：</span>
        <span class="summary-value">{{ deploymentSummary.coldReplicas }}</span>
        <span class="summary-unit">个</span>
      </div>
    </div>

    <div class="table-header">
      <div class="table-filters">
        <el-input
          v-model="filters.name"
          placeholder="搜索 Deployment 名称" 
          prefix-icon="el-icon-search"
          clearable
          class="filter-item search-input"
          style="width: 300px;"
        />
        <el-select
          v-model="filters.status"
          placeholder="选择状态"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="status in statusOptions"
            :key="status"
            :label="status"
            :value="status"
          />
        </el-select>
        <el-select
          v-model="filters.customerRegion"
          placeholder="选择客户区域"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="region in customerRegionOptions"
            :key="region"
            :label="region"
            :value="region"
          />
        </el-select>
        <el-select
          v-model="filters.gpuRequests"
          placeholder="选择GPU请求数量"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="gpu in gpuRequestsOptions"
            :key="gpu"
            :label="gpu"
            :value="gpu"
          />
        </el-select>
        <el-select
          v-model="filters.image"
          placeholder="选择镜像"
          clearable
          filterable
          allow-create
          class="filter-item"
          style="width: 300px;"
        >
          <el-option
            v-for="image in imageOptions"
            :key="image"
            :label="image"
            :value="image"
          />
        </el-select>
        <el-button 
          type="primary" 
          size="small" 
          @click="clearFilters"
          class="clear-filters-btn"
        >
          清除过滤
        </el-button>
      </div>
      <div class="action-buttons">
        <el-button type="info" size="small" @click="showColumnSettings = true">
          <i class="el-icon-setting"></i> 列设置
        </el-button>
        <el-button type="primary" size="small" @click="refreshData">
          <i class="el-icon-refresh"></i> 刷新
        </el-button>
      </div>
    </div>

    <el-table 
      v-loading="loading" 
      :data="filteredNDeployments" 
      style="width: 100%"
      border
    >
      <el-table-column 
        label="名称" 
        min-width="200"
        prop="metadata.name"
        sortable
        v-if="getColumnVisible('name')"
      >
        <template #default="scope">
          <span>{{ scope.row.metadata?.name || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="命名空间" 
        width="120"
        prop="metadata.namespace"
        sortable
        v-if="getColumnVisible('namespace')"
      >
        <template #default="scope">
          <span>{{ scope.row.metadata?.namespace || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="热副本" 
        width="80" 
        align="center"
        prop="hot"
        sortable
        v-if="getColumnVisible('hot')"
      >
        <template #default="scope">
          <el-link type="primary" @click="jumpToWorker(scope.row, 'hot')">
            {{ scope.row.hot || 0 }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column 
        label="冷副本" 
        width="80" 
        align="center"
        prop="cold"
        sortable
        v-if="getColumnVisible('cold')"
      >
        <template #default="scope">
          <el-link type="primary" @click="jumpToWorker(scope.row, 'cold')">
            {{ scope.row.cold || 0 }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column 
        label="状态" 
        width="100" 
        align="center"
        prop="status"
        sortable
        v-if="getColumnVisible('status')"
      >
        <template #default="scope">
          <el-tag 
            :type="getDeploymentStatusType(scope.row.status)"
            size="small"
          >
            {{ scope.row.status || 'Unknown' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column 
        label="客户区域" 
        width="120"
        prop="customerRegion"
        sortable
        v-if="getColumnVisible('customerRegion')"
      >
        <template #default="scope">
          <span>{{ scope.row.customerRegion || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="GPU Request" 
        width="100" 
        align="center"
        v-if="getColumnVisible('gpuRequests')"
      >
        <template #default="scope">
          <span>{{ scope.row.resources?.gpuRequests?.number || 0 }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="关联模型" 
        width="280" 
        align="center"
        v-if="getColumnVisible('models')"
      >
        <template #default="scope">
          <div class="model-tags">
            <template v-if="getDeploymentModels(scope.row.metadata?.name).length">
              <el-tag
                v-for="model in getDeploymentModels(scope.row.metadata?.name)"
                :key="model.name + '_' + model.source"
                size="small"
                :type="isbetaModel(model) ? 'warning' : 'success'"
                effect="plain"
                class="model-tag"
                @click="openModelDetail(model)"
              >
                {{ model.name }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column 
        label="镜像" 
        min-width="200"
        prop="image"
        sortable
        v-if="getColumnVisible('image')"
      >
        <template #default="scope">
          <span class="image-text">{{ scope.row.image || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="创建时间" 
        min-width="180"
        prop="age"
        sortable
        v-if="getColumnVisible('age')"
      >
        <template #default="scope">
          <span>{{ formatAge(scope.row.age) }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="操作" 
        width="180" 
        align="center"
        v-if="getColumnVisible('actions')"
      >
        <template #default="scope">
          <el-button
            type="primary"
            size="small"
            @click="viewDetails(scope.row)"
          >
            详情
          </el-button>
          <el-button
            type="info"
            size="small"
            @click="viewConfig(scope.row)"
          >
            配置
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 列设置对话框 -->
    <el-dialog
      v-model="showColumnSettings"
      title="列设置"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="column-settings">
        <div 
          v-for="column in columnConfig" 
          :key="column.key"
          class="column-item"
        >
          <el-checkbox 
            v-model="column.visible" 
            :disabled="column.key === 'name'"
            @change="saveColumnConfig"
          >
            {{ column.label }}
          </el-checkbox>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showColumnSettings = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailsDialogVisible"
      title="NDeployment 详情"
      width="60%"
      :close-on-click-modal="false"
    >
      <div v-if="selectedNDeployment">
        <pre>{{ JSON.stringify(selectedNDeployment, null, 2) }}</pre>
      </div>
    </el-dialog>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      title="Deployment 配置"
      width="600px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedConfigDeployment" class="config-content">
        <!-- 并发配置 Section -->
        <div class="config-section">
          <h3 class="section-title">
            <i class="el-icon-cpu"></i>
            并发配置
          </h3>
          <div class="section-content">
            <div class="config-item">
              <span class="config-label">单worker最大并发：</span>
              <span class="config-value">{{ selectedConfigDeployment.maxBatchSize || '-' }}</span>
            </div>
          </div>
        </div>

        <!-- 扩缩容配置 Section -->
        <div class="config-section" v-if="selectedConfigDeployment.simpleNSP">
          <h3 class="section-title">
            <i class="el-icon-s-operation"></i>
            扩缩容配置
          </h3>
          <div class="section-content">
            <div class="config-item">
              <span class="config-label">最大副本数：</span>
              <span class="config-value">{{ selectedConfigDeployment.simpleNSP.maxReplicas || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="config-label">最小副本数：</span>
              <span class="config-value">{{ selectedConfigDeployment.simpleNSP.minReplicas || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="config-label">扩容阈值：</span>
              <span class="config-value">{{ selectedConfigDeployment.simpleNSP.batchPerWorker || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="config-label">最大缩容速率：</span>
              <span class="config-value">{{ selectedConfigDeployment.simpleNSP.maxScaleDownRate || '-' }}</span>
            </div>
          </div>
        </div>

        <!-- 调试信息 Section -->
        <div class="config-section">
          <h3 class="section-title">
            <i class="el-icon-info"></i>
            调试信息
          </h3>
          <div class="section-content">
            <div class="config-item">
              <span class="config-label">MaxBatchSize字段：</span>
              <span class="config-value">{{ selectedConfigDeployment.maxBatchSize || 'undefined' }}</span>
            </div>
            <div class="config-item">
              <span class="config-label">SimpleNSP对象：</span>
              <span class="config-value">{{ selectedConfigDeployment.simpleNSP ? '存在' : '不存在' }}</span>
            </div>
            <div class="config-item" v-if="selectedConfigDeployment.simpleNSP">
              <span class="config-label">SimpleNSP内容：</span>
              <span class="config-value">{{ JSON.stringify(selectedConfigDeployment.simpleNSP) }}</span>
            </div>
            <div class="config-item">
              <span class="config-label">完整数据：</span>
              <span class="config-value">
                <el-button size="small" @click="showRawData = !showRawData">
                  {{ showRawData ? '隐藏' : '显示' }}原始数据
                </el-button>
              </span>
            </div>
            <div v-if="showRawData" class="raw-data">
              <pre>{{ JSON.stringify(selectedConfigDeployment, null, 2) }}</pre>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="configDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getNDeployments } from '@/api/nebula'
import axios from 'axios'

const loading = ref(false)
const ndeployments = ref([])
const detailsDialogVisible = ref(false)
const selectedNDeployment = ref(null)
const showColumnSettings = ref(false)
const configDialogVisible = ref(false)
const selectedConfigDeployment = ref(null)
const showRawData = ref(false)
const route = useRoute()
const router = useRouter()

// 模型数据
const models = ref([])
const betaModels = ref([])
const deploymentToModels = ref(new Map())

// 列配置
const columnConfig = ref([
  { key: 'name', label: '名称', visible: true },
  { key: 'namespace', label: '命名空间', visible: true },
  { key: 'hot', label: '热副本', visible: true },
  { key: 'cold', label: '冷副本', visible: true },
  { key: 'status', label: '状态', visible: true },
  { key: 'customerRegion', label: '客户区域', visible: true },
  { key: 'gpuRequests', label: 'GPU Request', visible: true },
  { key: 'models', label: '关联模型', visible: true },
  { key: 'image', label: '镜像', visible: true },
  { key: 'age', label: '创建时间', visible: true },
  { key: 'actions', label: '操作', visible: true }
])

// 过滤条件
const filters = ref({
  name: '',
  status: '',
  customerRegion: '',
  gpuRequests: '',
  image: ''
})

// 恢复过滤条件
const loadFilters = () => {
  const saved = localStorage.getItem('nebula-deployment-filters')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      Object.assign(filters.value, parsed)
    } catch (e) {}
  }
}

// 监听filters变化，自动保存
watch(filters, (val) => {
  localStorage.setItem('nebula-deployment-filters', JSON.stringify(val))
}, { deep: true })

// 动态生成过滤选项
const namespaceOptions = computed(() => {
  const namespaces = [...new Set(ndeployments.value.map(ndeployment => ndeployment.metadata?.namespace).filter(Boolean))]
  return namespaces.sort()
})

const statusOptions = computed(() => {
  const statuses = [...new Set(ndeployments.value.map(ndeployment => ndeployment.status).filter(Boolean))]
  return statuses.sort()
})

const customerRegionOptions = computed(() => {
  const regions = [...new Set(ndeployments.value.map(ndeployment => ndeployment.customerRegion).filter(Boolean))]
  return regions.sort()
})

const gpuRequestsOptions = computed(() => {
  const gpus = [...new Set(ndeployments.value.map(ndeployment => ndeployment.resources?.gpuRequests?.number).filter(Boolean))]
  return gpus.sort((a, b) => a - b)
})

const imageOptions = computed(() => {
  const images = [...new Set(ndeployments.value.map(ndeployment => ndeployment.image).filter(Boolean))]
  return images.sort()
})

// 过滤后的数据
const filteredNDeployments = computed(() => {
  let filtered = ndeployments.value

  // 名称过滤
  if (filters.value.name) {
    filtered = filtered.filter(ndeployment => 
      ndeployment.metadata?.name?.toLowerCase().includes(filters.value.name.toLowerCase()) ||
      ndeployment.customerRegion?.toLowerCase().includes(filters.value.name.toLowerCase())
    )
  }

  // 状态过滤
  if (filters.value.status) {
    filtered = filtered.filter(ndeployment => ndeployment.status === filters.value.status)
  }

  // 客户区域过滤
  if (filters.value.customerRegion) {
    filtered = filtered.filter(ndeployment => ndeployment.customerRegion === filters.value.customerRegion)
  }

  // GPU请求数量过滤
  if (filters.value.gpuRequests) {
    filtered = filtered.filter(ndeployment => ndeployment.resources?.gpuRequests?.number === filters.value.gpuRequests)
  }

  // 镜像过滤
  if (filters.value.image) {
    filtered = filtered.filter(ndeployment => ndeployment.image && ndeployment.image.includes(filters.value.image))
  }

  return filtered
})

// Deployment汇总统计
const deploymentSummary = computed(() => {
  let hotReplicas = 0
  let coldReplicas = 0

  filteredNDeployments.value.forEach(ndeployment => {
    hotReplicas += ndeployment.hot || 0
    coldReplicas += ndeployment.cold || 0
  })

  return {
    hotReplicas,
    coldReplicas
  }
})

// 获取列是否可见
const getColumnVisible = (key) => {
  const column = columnConfig.value.find(col => col.key === key)
  return column ? column.visible : true
}

// 保存列配置到本地存储
const saveColumnConfig = () => {
  localStorage.setItem('nebula-deployment-column-config', JSON.stringify(columnConfig.value))
}

// 从本地存储加载列配置
const loadColumnConfig = () => {
  const saved = localStorage.getItem('nebula-deployment-column-config')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      // 确保所有列都存在，避免新增列时丢失
      const defaultConfig = columnConfig.value
      columnConfig.value = defaultConfig.map(defaultCol => {
        const savedCol = parsed.find(col => col.key === defaultCol.key)
        return savedCol || defaultCol
      })
    } catch (error) {
      console.error('Failed to parse saved column config:', error)
    }
  }
}

// 清除所有过滤条件
const clearFilters = () => {
  filters.value = {
    name: '',
    status: '',
    customerRegion: '',
    gpuRequests: '',
    image: ''
  }
}

// 获取部署状态类型
const getDeploymentStatusType = (status) => {
  switch (status) {
    case 'Ready':
      return 'success'
    case 'NotReady':
      return 'danger'
    default:
      return 'info'
  }
}

// 格式化时间
const formatAge = (age) => {
  if (!age) return '-'
  // 如果是 ISO 格式的时间戳，转换为相对时间
  if (age.includes('T')) {
    const date = new Date(age)
    const now = new Date()
    const diffMs = now - date
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    const diffHours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    if (diffDays > 0) {
      return `${diffDays}天${diffHours}小时`
    } else {
      return `${diffHours}小时`
    }
  }
  return age
}

// 刷新数据
const refreshData = async () => {
  await fetchNDeployments()
}

// 获取 NDeployment 列表
const fetchNDeployments = async () => {
  loading.value = true
  try {
    // 先获取模型数据
    await fetchModelsData()
    
    const response = await getNDeployments()
    console.log('NDeployment API response:', response)
    ndeployments.value = response.data?.data || []
    console.log('NDeployments data:', ndeployments.value)
  } catch (error) {
    console.error('Failed to fetch NDeployments:', error)
  } finally {
    loading.value = false
  }
}

// 获取模型数据
const fetchModelsData = async () => {
  try {
    // 获取普通模型列表
    const modelsResponse = await axios.get('/models')
    models.value = modelsResponse.data?.data?.data || []

    // 获取 beta 模型列表
    const betaModelsResponse = await axios.get('/models', { params: { platform: 'beta' } })
    betaModels.value = betaModelsResponse.data?.data?.data || []

    // 建立 deployment 到模型的映射
    buildDeploymentToModelsMapping()
  } catch (error) {
    console.error('Failed to fetch models data:', error)
  }
}

// 从 2.0 endpoint URL 中提取 deployment 名称
const extractDeploymentFromURL = (url) => {
  if (!url) return ''
  // 2.0 format: https://{deployment}.us-01.sls2.alpha.ai...
  const match = url.match(/https:\/\/([^.]+)\.us-01\.sls2\.alpha\.ai/)
  return match ? match[1] : ''
}

// 建立 deployment 到模型的映射关系
const buildDeploymentToModelsMapping = () => {
  const deploymentToModelsMap = new Map()

  // 处理普通模型
  models.value.forEach(model => {
    if (model.endpoints) {
      model.endpoints.forEach(endpoint => {
        if (endpoint.url && endpoint.url.includes('sls2.alpha.ai')) {
          const deploymentName = extractDeploymentFromURL(endpoint.url)
          if (deploymentName) {
            if (!deploymentToModelsMap.has(deploymentName)) {
              deploymentToModelsMap.set(deploymentName, new Map())
            }
            deploymentToModelsMap.get(deploymentName).set(model.model_name, {
              name: model.model_name,
              source: 'normal'
            })
          }
        }
      })
    }
  })

  // 处理 beta 模型
  betaModels.value.forEach(model => {
    if (model.endpoints) {
      model.endpoints.forEach(endpoint => {
        if (endpoint.url && endpoint.url.includes('sls2.alpha.ai')) {
          const deploymentName = extractDeploymentFromURL(endpoint.url)
          if (deploymentName) {
            if (!deploymentToModelsMap.has(deploymentName)) {
              deploymentToModelsMap.set(deploymentName, new Map())
            }
            deploymentToModelsMap.get(deploymentName).set(model.model_name + '_beta', {
              name: model.model_name,
              source: 'beta'
            })
          }
        }
      })
    }
  })

  deploymentToModels.value = deploymentToModelsMap
}

// 获取 deployment 关联的模型
const getDeploymentModels = (deploymentName) => {
  const modelsMap = deploymentToModels.value.get(deploymentName)
  return modelsMap ? Array.from(modelsMap.values()) : []
}

// 判断是否是 beta 模型
const isbetaModel = (model) => {
  return model.source === 'beta'
}

// 打开模型详情
const openModelDetail = (model) => {
  const platform = model.source === 'beta' ? 'beta' : ''
  router.push({
    name: 'EndpointManager',
    params: { modelName: model.name },
    query: platform ? { platform } : undefined
  })
}

// 查看详情
const viewDetails = (ndeployment) => {
  selectedNDeployment.value = ndeployment
  detailsDialogVisible.value = true
}

// 查看配置
const viewConfig = (ndeployment) => {
  selectedConfigDeployment.value = ndeployment
  configDialogVisible.value = true
}

// 组件挂载时获取数据
onMounted(() => {
  loadColumnConfig()
  loadFilters()
  fetchNDeployments()
  // query参数自动过滤（只支持name）
  if (route.query.name) {
    filters.value.name = route.query.name
  }
})

// 监听路由变化，支持后续跳转
watch(() => route.query.name, (val) => {
  if (val) {
    filters.value.name = val
  }
})

// 跳转到Worker页面并带过滤条件
const jumpToWorker = (ndeployment, type) => {
  const deploymentName = ndeployment.metadata?.name
  if (!deploymentName) return
  router.push({ name: 'NebulaWorker', query: { deployment: deploymentName, phase: type } })
}
</script>

<style scoped>
.ndeployment-container {
  padding: 20px;
  margin-top: 0 !important;
}

.data-summary {
  display: flex;
  gap: 30px;
  margin-bottom: 20px;
  margin-top: 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.summary-label {
  font-weight: 500;
  color: #495057;
}

.summary-value {
  font-weight: bold;
  font-size: 18px;
  color: #007bff;
}

.summary-unit {
  color: #6c757d;
  font-size: 14px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.table-filters {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.filter-item {
  width: 180px;
}

.clear-filters-btn {
  margin-left: 10px;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.image-text {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #666;
  word-break: break-all;
  line-height: 1.4;
}

.model-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.model-tag {
  cursor: pointer;
}

.model-tag:hover {
  opacity: 0.8;
  transform: scale(1.05);
  transition: all 0.2s ease;
}

.column-settings {
  max-height: 300px;
  overflow-y: auto;
}

.column-item {
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.column-item:last-child {
  border-bottom: none;
}

/* 配置对话框样式 */
.config-content {
  max-height: 500px;
  overflow-y: auto;
}

.config-section {
  margin-bottom: 25px;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  overflow: hidden;
}

.config-section:last-child {
  margin-bottom: 0;
}

.section-title {
  margin: 0;
  padding: 12px 16px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  font-size: 16px;
  font-weight: 600;
  color: #495057;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title i {
  color: #007bff;
}

.section-content {
  padding: 16px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.config-item:last-child {
  border-bottom: none;
}

.config-label {
  font-weight: 500;
  color: #495057;
  min-width: 150px;
}

.config-value {
  font-weight: 600;
  color: #007bff;
  text-align: right;
  flex: 1;
  margin-left: 16px;
}

.raw-data {
  margin-top: 16px;
  padding: 12px;
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  max-height: 300px;
  overflow-y: auto;
}

.raw-data pre {
  margin: 0;
  font-size: 12px;
  line-height: 1.4;
  color: #495057;
}

/* 可能存在的全局样式修正 */
:deep(.el-card__body),
:deep(.el-main),
:deep(.el-container) {
  margin-top: 0 !important;
  padding-top: 0 !important;
}
</style> 