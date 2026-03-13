<template>
  <div class="worker-container">
    <!-- 数据总结 -->
    <div class="data-summary">
      <div class="summary-item">
        <span class="summary-label">当前显示：</span>
        <span class="summary-value">{{ filteredWorkers.length }}</span>
        <span class="summary-unit">条数据</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">热阶段：</span>
        <span class="summary-value">{{ workerSummary.hotCount }}</span>
        <span class="summary-unit">个</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">冷阶段：</span>
        <span class="summary-value">{{ workerSummary.coldCount }}</span>
        <span class="summary-unit">个</span>
      </div>
    </div>

    <div class="table-header">
      <div class="table-filters">
        <el-input
          v-model="filters.name"
          placeholder="搜索 Worker 名称" 
          prefix-icon="el-icon-search"
          clearable
          class="filter-item search-input"
          style="width: 300px;"
        />
        <el-select
          v-model="filters.phase"
          placeholder="选择阶段"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="phase in phaseOptions"
            :key="phase"
            :label="phase"
            :value="phase"
          />
        </el-select>
        <el-select
          v-model="filters.gpuProduct"
          placeholder="选择GPU产品"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="gpu in gpuProductOptions"
            :key="gpu"
            :label="gpu"
            :value="gpu"
          />
        </el-select>
        <el-select
          v-model="filters.node"
          placeholder="选择Node"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="node in nodeOptions"
            :key="node"
            :label="node"
            :value="node"
          />
        </el-select>
        <el-select
          v-model="filters.deployment"
          placeholder="选择Deployment"
          clearable
          class="filter-item"
        >
          <el-option
            v-for="dep in deploymentOptions"
            :key="dep"
            :label="dep"
            :value="dep"
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
      :data="filteredWorkers" 
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
        label="阶段" 
        width="100" 
        align="center"
        prop="phase"
        sortable
        v-if="getColumnVisible('phase')"
      >
        <template #default="scope">
          <el-tag 
            :type="getPhaseType(scope.row.phase)"
            size="small"
          >
            {{ scope.row.phase || 'Unknown' }}
          </el-tag>
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
        label="Deployment" 
        min-width="180"
        prop="deployment"
        v-if="getColumnVisible('deployment')"
      >
        <template #default="scope">
          <el-link type="primary" @click="jumpToDeployment(scope.row.deployment)">
            {{ scope.row.deployment || '-' }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column 
        label="Node" 
        min-width="180"
        prop="node"
        v-if="getColumnVisible('node')"
      >
        <template #default="scope">
          <el-link type="primary" @click="jumpToCNode(scope.row.node)">
            {{ scope.row.node || '-' }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column 
        label="运行时间" 
        width="120"
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
      title="Nebula Worker 详情"
      width="60%"
      :close-on-click-modal="false"
    >
      <div v-if="selectedWorker">
        <pre>{{ JSON.stringify(selectedWorker, null, 2) }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getNebulaWorkers } from '@/api/nebula'

const loading = ref(false)
const workers = ref([])
const detailsDialogVisible = ref(false)
const selectedWorker = ref(null)
const showColumnSettings = ref(false)
const router = useRouter()
const route = useRoute()

// 列配置
const columnConfig = ref([
  { key: 'name', label: '名称', visible: true },
  { key: 'phase', label: '阶段', visible: true },
  { key: 'gpuProduct', label: 'GPU 产品', visible: true },
  { key: 'deployment', label: 'Deployment', visible: true },
  { key: 'node', label: 'Node', visible: true },
  { key: 'age', label: '运行时间', visible: true },
  { key: 'actions', label: '操作', visible: true }
])

// 过滤条件
const filters = ref({
  name: '',
  phase: '',
  gpuProduct: '',
  deployment: '',
  node: ''
})

// 恢复过滤条件
const loadFilters = () => {
  const saved = localStorage.getItem('nebula-worker-filters')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      Object.assign(filters.value, parsed)
    } catch (e) {}
  }
}

// 监听filters变化，自动保存
watch(filters, (val) => {
  localStorage.setItem('nebula-worker-filters', JSON.stringify(val))
}, { deep: true })

// 动态生成过滤选项
const phaseOptions = computed(() => {
  const phases = [...new Set(workers.value.map(worker => worker.phase).filter(Boolean))]
  return phases.sort()
})

const gpuProductOptions = computed(() => {
  const gpus = [...new Set(workers.value.map(worker => worker.gpuProduct).filter(Boolean))]
  return gpus.sort()
})

const deploymentOptions = computed(() => {
  const deps = [...new Set(workers.value.map(worker => worker.deployment).filter(Boolean))]
  return deps.sort()
})

const nodeOptions = computed(() => {
  const nodes = [...new Set(workers.value.map(worker => worker.node).filter(Boolean))]
  return nodes.sort()
})

// 过滤后的数据
const filteredWorkers = computed(() => {
  let filtered = workers.value

  // 名称过滤
  if (filters.value.name) {
    filtered = filtered.filter(worker => 
      worker.metadata?.name?.toLowerCase().includes(filters.value.name.toLowerCase()) ||
      worker.gpuProduct?.toLowerCase().includes(filters.value.name.toLowerCase())
    )
  }

  // 阶段过滤
  if (filters.value.phase) {
    filtered = filtered.filter(worker => worker.phase === filters.value.phase)
  }

  // GPU产品过滤
  if (filters.value.gpuProduct) {
    filtered = filtered.filter(worker => worker.gpuProduct === filters.value.gpuProduct)
  }

  // Node过滤
  if (filters.value.node) {
    filtered = filtered.filter(worker => worker.node === filters.value.node)
  }

  // Deployment过滤
  if (filters.value.deployment) {
    filtered = filtered.filter(worker => worker.deployment === filters.value.deployment)
  }

  return filtered
})

// Worker汇总统计
const workerSummary = computed(() => {
  let hotCount = 0
  let coldCount = 0

  filteredWorkers.value.forEach(worker => {
    if (worker.phase === 'hot') {
      hotCount++
    } else if (worker.phase === 'cold') {
      coldCount++
    }
  })

  return {
    hotCount,
    coldCount
  }
})

// 获取列是否可见
const getColumnVisible = (key) => {
  const column = columnConfig.value.find(col => col.key === key)
  return column ? column.visible : true
}

// 保存列配置到本地存储
const saveColumnConfig = () => {
  localStorage.setItem('nebula-worker-column-config', JSON.stringify(columnConfig.value))
}

// 从本地存储加载列配置
const loadColumnConfig = () => {
  const saved = localStorage.getItem('nebula-worker-column-config')
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
    phase: '',
    gpuProduct: '',
    deployment: '',
    node: ''
  }
}

// 获取阶段类型
const getPhaseType = (phase) => {
  switch (phase) {
    case 'hot':
      return 'danger'
    case 'cold':
      return 'info'
    default:
      return 'warning'
  }
}

// 跳转到Deployment页面并带过滤条件
const jumpToDeployment = (deploymentName) => {
  if (!deploymentName) return
  router.push({ name: 'NebulaDeployment', query: { name: deploymentName } })
}

// 跳转到CNode页面并带过滤条件
const jumpToCNode = (nodeName) => {
  if (!nodeName) return
  router.push({ name: 'NebulaCnode', query: { name: nodeName } })
}

// 刷新数据
const refreshData = async () => {
  await fetchWorkers()
}

// 获取 Nebula Worker 列表
const fetchWorkers = async () => {
  loading.value = true
  try {
    const response = await getNebulaWorkers()
    console.log('Nebula Worker API response:', response)
    workers.value = response.data?.data || []
    console.log('Nebula Workers data:', workers.value)
  } catch (error) {
    console.error('Failed to fetch Nebula Workers:', error)
  } finally {
    loading.value = false
  }
}

// 查看详情
const viewDetails = (worker) => {
  selectedWorker.value = worker
  detailsDialogVisible.value = true
}

// 组件挂载时获取数据
onMounted(() => {
  loadColumnConfig()
  loadFilters()
  fetchWorkers()
  // query参数自动过滤
  if (route.query.deployment) {
    filters.value.deployment = route.query.deployment
  }
  if (route.query.phase) {
    filters.value.phase = route.query.phase
  }
})

watch(() => route.query.deployment, (val) => {
  if (val) {
    filters.value.deployment = val
  }
})

watch(() => route.query.phase, (val) => {
  if (val) {
    filters.value.phase = val
  }
})
</script>

<style scoped>
.worker-container {
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
</style> 