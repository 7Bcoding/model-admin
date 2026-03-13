<template>
  <div class="cnode-container">
    <!-- 数据总结 -->
    <div class="data-summary">
      <div class="summary-item">
        <span class="summary-label">当前显示：</span>
        <span class="summary-value">{{ filteredCNodes.length }}</span>
        <span class="summary-unit">条数据</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">GPU使用情况：</span>
        <span class="summary-value">{{ gpuUsageSummary.used }}</span>
        <span class="summary-unit">/ {{ gpuUsageSummary.total }}</span>
        <span class="summary-unit">({{ gpuUsageSummary.usageRate }}%)</span>
      </div>
    </div>

    <div class="table-header">
      <div class="table-filters">
        <el-input
          v-model="filters.name"
          placeholder="搜索 CNode 名称" 
          prefix-icon="el-icon-search"
          clearable
          class="filter-item search-input"
          style="width: 300px;"
        />
        <el-select
          v-model="filters.region"
          placeholder="选择区域"
          clearable
          class="filter-item filter-item-small"
        >
          <el-option
            v-for="region in regionOptions"
            :key="region"
            :label="region"
            :value="region"
          />
        </el-select>
        <el-select
          v-model="filters.gpuProduct"
          placeholder="选择GPU类型"
          clearable
          class="filter-item filter-item-small"
        >
          <el-option
            v-for="gpu in gpuProductOptions"
            :key="gpu"
            :label="gpu"
            :value="gpu"
          />
        </el-select>
        <el-select
          v-model="filters.policy"
          placeholder="选择策略"
          clearable
          class="filter-item filter-item-small"
        >
          <el-option
            v-for="policy in policyOptions"
            :key="policy"
            :label="policy"
            :value="policy"
          />
        </el-select>
        <el-select
          v-model="filters.status"
          placeholder="选择状态"
          clearable
          class="filter-item filter-item-small"
        >
          <el-option
            v-for="status in statusOptions"
            :key="status"
            :label="status"
            :value="status"
          />
        </el-select>
        <el-select
          v-model="filters.hostname"
          placeholder="选择主机名"
          clearable
          filterable
          allow-create
          class="filter-item filter-item-small"
        >
          <el-option
            v-for="hostname in hostnameOptions"
            :key="hostname"
            :label="hostname"
            :value="hostname"
          />
        </el-select>
        <el-select
          v-model="filters.image"
          placeholder="选择或输入Image缓存"
          clearable
          filterable
          allow-create
          default-first-option
          class="filter-item filter-item-large"
          @change="handleImageFilterChange"
        >
          <el-option
            v-for="image in imageOptions"
            :key="image"
            :label="image"
            :value="image"
          />
        </el-select>
        <el-select
          v-model="filters.model"
          placeholder="选择或输入Model缓存"
          clearable
          filterable
          allow-create
          default-first-option
          class="filter-item filter-item-large"
          @change="handleModelFilterChange"
        >
          <el-option
            v-for="model in modelOptions"
            :key="model"
            :label="model"
            :value="model"
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
      :data="filteredCNodes" 
      style="width: 100%"
      border
    >
      <el-table-column 
        label="名称" 
        min-width="200"
        prop="name"
        sortable
        v-if="getColumnVisible('name')"
      >
        <template #default="scope">
          <span>{{ scope.row.name || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="提供商" 
        width="100"
        prop="provider"
        sortable
        v-if="getColumnVisible('provider')"
      >
        <template #default="scope">
          <span>{{ scope.row.provider || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="区域" 
        width="120"
        prop="region"
        sortable
        v-if="getColumnVisible('region')"
      >
        <template #default="scope">
          <span>{{ scope.row.region || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="主机名" 
        width="150"
        prop="hostName"
        sortable
        v-if="getColumnVisible('hostname')"
      >
        <template #default="scope">
          <span>{{ scope.row.hostName || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="CPU 产品" 
        min-width="200"
        prop="cpuProduct"
        v-if="getColumnVisible('cpuProduct')"
      >
        <template #default="scope">
          <span class="product-text">{{ scope.row.cpuProduct || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="GPU 产品" 
        min-width="150"
        prop="gpuProduct"
        v-if="getColumnVisible('gpuProduct')"
      >
        <template #default="scope">
          <span class="product-text">{{ scope.row.gpuProduct || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="策略" 
        width="120"
        prop="policy"
        sortable
        v-if="getColumnVisible('policy')"
      >
        <template #default="scope">
          <el-tag 
            :type="getPolicyType(scope.row.policy)"
            size="small"
          >
            {{ scope.row.policy || '-' }}
          </el-tag>
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
            :type="getStatusType(scope.row.status)"
            size="small"
          >
            {{ scope.row.status || 'Unknown' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column 
        label="GPU 状态" 
        min-width="150"
        prop="gpuState"
        v-if="getColumnVisible('gpuState')"
      >
        <template #default="scope">
          <span class="gpu-state">{{ formatGpuState(scope.row.gpuState) }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="Image缓存" 
        min-width="200"
        prop="imageCache"
        v-if="getColumnVisible('imageCache')"
      >
        <template #default="scope">
          <div v-if="scope.row.imageCache && scope.row.imageCache.length > 0">
            <div 
              v-for="(item, index) in scope.row.imageCache" 
              :key="index"
              class="cache-item"
            >
              <el-tooltip 
                :content="item.name" 
                placement="top" 
                :show-after="500"
                popper-class="copyable-tooltip"
              >
                <span class="cache-name" @click="copyToClipboard(item.name)">
                  {{ item.name }}
                </span>
              </el-tooltip>
              <el-tag 
                :type="getCacheStatusType(item.status)"
                size="mini"
                class="cache-status"
              >
                {{ item.status }}
              </el-tag>
              <span class="cache-size">{{ formatSize(item.size) }}</span>
            </div>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="Model缓存" 
        min-width="200"
        prop="modelCache"
        v-if="getColumnVisible('modelCache')"
      >
        <template #default="scope">
          <div v-if="scope.row.modelCache && scope.row.modelCache.length > 0">
            <div 
              v-for="(item, index) in scope.row.modelCache" 
              :key="index"
              class="cache-item"
            >
              <el-tooltip 
                :content="item.name" 
                placement="top" 
                :show-after="500"
                popper-class="copyable-tooltip"
              >
                <span class="cache-name" @click="copyToClipboard(item.name)">
                  {{ item.name }}
                </span>
              </el-tooltip>
              <el-tag 
                :type="getCacheStatusType(item.status)"
                size="mini"
                class="cache-status"
              >
                {{ item.status }}
              </el-tag>
              <span class="cache-size">{{ formatSize(item.size) }}</span>
            </div>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="运行时间" 
        width="100"
        prop="age"
        sortable
        v-if="getColumnVisible('age')"
      >
        <template #default="scope">
          <span>{{ scope.row.age || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        label="操作" 
        width="120" 
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
      title="CNode 详情"
      width="60%"
      :close-on-click-modal="false"
    >
      <div v-if="selectedCNode">
        <pre>{{ JSON.stringify(selectedCNode, null, 2) }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getCNodes } from '@/api/nebula'

const loading = ref(false)
const cnodes = ref([])
const detailsDialogVisible = ref(false)
const selectedCNode = ref(null)
const showColumnSettings = ref(false)

// 列配置
const columnConfig = ref([
  { key: 'name', label: '名称', visible: true },
  { key: 'provider', label: '提供商', visible: true },
  { key: 'region', label: '区域', visible: true },
  { key: 'hostname', label: '主机名', visible: true },
  { key: 'cpuProduct', label: 'CPU 产品', visible: true },
  { key: 'gpuProduct', label: 'GPU 产品', visible: true },
  { key: 'policy', label: '策略', visible: true },
  { key: 'status', label: '状态', visible: true },
  { key: 'gpuState', label: 'GPU 状态', visible: true },
  { key: 'imageCache', label: 'Image缓存', visible: true },
  { key: 'modelCache', label: 'Model缓存', visible: true },
  { key: 'age', label: '运行时间', visible: true },
  { key: 'actions', label: '操作', visible: true }
])

// 过滤条件
const filters = ref({
  name: '',
  region: '',
  gpuProduct: '',
  policy: '',
  status: '',
  hostname: '',
  image: '',
  model: ''
})

// 恢复过滤条件
const loadFilters = () => {
  const saved = localStorage.getItem('nebula-cnode-filters')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      Object.assign(filters.value, parsed)
    } catch (e) {}
  }
}

// 监听filters变化，自动保存
watch(filters, (val) => {
  localStorage.setItem('nebula-cnode-filters', JSON.stringify(val))
}, { deep: true })

// 动态生成过滤选项
const regionOptions = computed(() => {
  const regions = [...new Set(cnodes.value.map(cnode => cnode.region).filter(Boolean))]
  return regions.sort()
})

const gpuProductOptions = computed(() => {
  const gpus = [...new Set(cnodes.value.map(cnode => cnode.gpuProduct).filter(Boolean))]
  return gpus.sort()
})

const policyOptions = computed(() => {
  const policies = [...new Set(cnodes.value.map(cnode => cnode.policy).filter(Boolean))]
  return policies.sort()
})

const statusOptions = computed(() => {
  const statuses = [...new Set(cnodes.value.map(cnode => cnode.status).filter(Boolean))]
  return statuses.sort()
})

const hostnameOptions = computed(() => {
  const hostnames = [...new Set(cnodes.value.map(cnode => cnode.hostName).filter(Boolean))]
  return hostnames.sort()
})

const imageOptions = computed(() => {
  const images = new Set()
  cnodes.value.forEach(cnode => {
    if (cnode.imageCache && Array.isArray(cnode.imageCache)) {
      cnode.imageCache.forEach(item => {
        if (item.status === 'Ready') {
          images.add(item.name)
        }
      })
    }
  })
  return Array.from(images).sort()
})

const modelOptions = computed(() => {
  const models = new Set()
  cnodes.value.forEach(cnode => {
    if (cnode.modelCache && Array.isArray(cnode.modelCache)) {
      cnode.modelCache.forEach(item => {
        if (item.status === 'Ready') {
          models.add(item.name)
        }
      })
    }
  })
  return Array.from(models).sort()
})

// 过滤后的数据
const filteredCNodes = computed(() => {
  let filtered = cnodes.value

  // 名称过滤
  if (filters.value.name) {
    filtered = filtered.filter(cnode => 
      cnode.name?.toLowerCase().includes(filters.value.name.toLowerCase()) ||
      cnode.provider?.toLowerCase().includes(filters.value.name.toLowerCase())
    )
  }

  // 区域过滤
  if (filters.value.region) {
    filtered = filtered.filter(cnode => cnode.region === filters.value.region)
  }

  // GPU类型过滤
  if (filters.value.gpuProduct) {
    filtered = filtered.filter(cnode => cnode.gpuProduct === filters.value.gpuProduct)
  }

  // 策略过滤
  if (filters.value.policy) {
    filtered = filtered.filter(cnode => cnode.policy === filters.value.policy)
  }

  // 状态过滤
  if (filters.value.status) {
    filtered = filtered.filter(cnode => cnode.status === filters.value.status)
  }

  // 主机名过滤
  if (filters.value.hostname) {
    filtered = filtered.filter(cnode => cnode.hostName === filters.value.hostname)
  }

  // Image缓存过滤
  if (filters.value.image) {
    filtered = filtered.filter(cnode => {
      if (!cnode.imageCache || !Array.isArray(cnode.imageCache)) return false
      return cnode.imageCache.some(item => 
        item.name.includes(filters.value.image) && item.status === 'Ready'
      )
    })
  }

  // Model缓存过滤
  if (filters.value.model) {
    filtered = filtered.filter(cnode => {
      if (!cnode.modelCache || !Array.isArray(cnode.modelCache)) return false
      return cnode.modelCache.some(item => 
        item.name.includes(filters.value.model) && item.status === 'Ready'
      )
    })
  }

  return filtered
})

// GPU使用情况汇总
const gpuUsageSummary = computed(() => {
  let total = 0
  let used = 0

  filteredCNodes.value.forEach(cnode => {
    if (cnode.gpuState && Array.isArray(cnode.gpuState)) {
      total += cnode.gpuState.length
      used += cnode.gpuState.filter(state => state === 1).length
    }
  })

  const usageRate = total > 0 ? Math.round((used / total) * 100) : 0

  return {
    total,
    used,
    usageRate
  }
})

// 获取列是否可见
const getColumnVisible = (key) => {
  const column = columnConfig.value.find(col => col.key === key)
  return column ? column.visible : true
}

// 保存列配置到本地存储
const saveColumnConfig = () => {
  localStorage.setItem('nebula-cnode-column-config', JSON.stringify(columnConfig.value))
}

// 从本地存储加载列配置
const loadColumnConfig = () => {
  const saved = localStorage.getItem('nebula-cnode-column-config')
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
    region: '',
    gpuProduct: '',
    policy: '',
    status: '',
    hostname: '',
    image: '',
    model: ''
  }
}

// 获取策略类型
const getPolicyType = (policy) => {
  switch (policy) {
    case 'normal':
      return 'success'
    case 'evictWorkers':
      return 'warning'
    default:
      return 'info'
  }
}

// 获取状态类型
const getStatusType = (status) => {
  switch (status) {
    case 'Ready':
      return 'success'
    case 'NotReady':
      return 'danger'
    default:
      return 'info'
  }
}

// 获取缓存状态类型
const getCacheStatusType = (status) => {
  switch (status) {
    case 'Ready':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Failed':
    case 'FailedToPull':
      return 'danger'
    default:
      return 'info'
  }
}

// 格式化大小
const formatSize = (size) => {
  if (!size) return '-'
  if (typeof size === 'number') {
    if (size >= 1024 * 1024 * 1024) {
      return `${(size / (1024 * 1024 * 1024)).toFixed(1)}GB`
    } else if (size >= 1024 * 1024) {
      return `${(size / (1024 * 1024)).toFixed(1)}MB`
    } else if (size >= 1024) {
      return `${(size / 1024).toFixed(1)}KB`
    } else {
      return `${size}B`
    }
  }
  return size
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    // 可以添加一个简单的提示，比如使用Element Plus的Message组件
    ElMessage.success('已复制到剪贴板')
  } catch (err) {
    console.error('复制失败:', err)
    // 降级方案：使用传统的复制方法
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    try {
      document.execCommand('copy')
      ElMessage.success('已复制到剪贴板')
    } catch (fallbackErr) {
      console.error('降级复制也失败:', fallbackErr)
      ElMessage.error('复制失败')
    }
    document.body.removeChild(textArea)
  }
}

// 处理Image筛选变化
const handleImageFilterChange = (value) => {
  if (value && !imageOptions.value.includes(value)) {
    // 如果是自定义输入的值，可以在这里进行额外处理
    console.log('自定义Image筛选值:', value)
  }
}

// 处理Model筛选变化
const handleModelFilterChange = (value) => {
  if (value && !modelOptions.value.includes(value)) {
    // 如果是自定义输入的值，可以在这里进行额外处理
    console.log('自定义Model筛选值:', value)
  }
}

// 格式化 GPU 状态
const formatGpuState = (gpuState) => {
  if (!gpuState) return '-'
  if (Array.isArray(gpuState)) {
    return gpuState.join(', ')
  }
  return gpuState
}

// 刷新数据
const refreshData = async () => {
  await fetchCNodes()
}

// 获取 CNode 列表
const fetchCNodes = async () => {
  loading.value = true
  try {
    const response = await getCNodes()
    console.log('CNode API response:', response)
    cnodes.value = response.data?.data || []
    console.log('CNodes data:', cnodes.value)
  } catch (error) {
    console.error('Failed to fetch CNodes:', error)
  } finally {
    loading.value = false
  }
}

// 查看详情
const viewDetails = (cnode) => {
  selectedCNode.value = cnode
  detailsDialogVisible.value = true
}

// 组件挂载时获取数据
onMounted(() => {
  loadColumnConfig()
  loadFilters()
  fetchCNodes()
})
</script>

<style scoped>
.cnode-container {
  padding: 20px;
}

.data-summary {
  display: flex;
  gap: 30px;
  margin-bottom: 20px;
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

.filter-item-small {
  width: 120px;
}

.filter-item-large {
  width: 280px;
}

/* 自定义筛选框样式 */
.filter-item :deep(.el-input__wrapper) {
  background-color: #fff;
}

.filter-item :deep(.el-select .el-input__inner) {
  font-size: 13px;
}

.filter-item :deep(.el-select-dropdown__item) {
  font-size: 13px;
  padding: 8px 12px;
}

.filter-item :deep(.el-select-dropdown__item.selected) {
  background-color: #f0f9ff;
  color: #007bff;
}

.clear-filters-btn {
  margin-left: 10px;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.product-text {
  font-family: monospace;
  font-size: 12px;
}

.gpu-state {
  font-family: monospace;
  font-size: 12px;
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

.cache-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
}

.cache-item:last-child {
  border-bottom: none;
}

.cache-name {
  flex: 1;
  font-family: monospace;
  font-size: 12px;
  color: #495057;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  transition: color 0.2s;
  padding: 2px 4px;
  border-radius: 3px;
}

.cache-name:hover {
  color: #007bff;
  background-color: #f8f9fa;
}

.cache-status {
  flex-shrink: 0;
}

.cache-size {
  flex-shrink: 0;
  font-size: 11px;
  color: #6c757d;
  font-family: monospace;
}

/* 全局样式 - 可复制的tooltip */
:global(.copyable-tooltip) {
  max-width: 400px !important;
  word-break: break-all;
  font-family: monospace;
  font-size: 12px;
  line-height: 1.4;
}

:global(.copyable-tooltip .el-tooltip__content) {
  text-align: left;
}
</style> 