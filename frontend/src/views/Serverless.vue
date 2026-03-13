<template>
  <div class="page-container">
    <div class="header">
      <div class="title-section">
        <i class="el-icon-cloudy title-icon"></i>
        <h2 class="page-title">Serverless 管理</h2>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="goToAuditLog">
          <i class="el-icon-document"></i> 操作日志
        </el-button>
      </div>
    </div>

    <div class="content-card">
      <div class="table-header">
        <div class="table-filters">
          <el-input
            v-model="searchQuery"
            placeholder="搜索 Endpoint 名称或关联模型" 
            prefix-icon="el-icon-search"
            clearable
            class="filter-item search-input"
          />
          <el-select
            v-model="imageSearchQuery"
            placeholder="搜索镜像" 
            clearable
            filterable
            class="filter-item image-search-input"
          >
            <el-option
              v-for="image in uniqueImages"
              :key="image"
              :label="image"
              :value="image"
            />
          </el-select>
          <el-select
            v-model="selectedCluster"
            placeholder="集群"
            clearable
            class="cluster-filter"
          >
            <el-option
              v-for="cluster in clusters"
              :key="cluster"
              :label="cluster"
              :value="cluster"
            />
          </el-select>
          <el-select
            v-model="selectedGpuType"
            placeholder="GPU 型号"
            clearable
            class="gpu-filter"
          >
            <el-option
              v-for="type in gpuTypes"
              :key="type"
              :label="type"
              :value="type"
            />
          </el-select>
          <el-checkbox v-model="showInstancesOnly" class="instance-filter">
            仅看有实例
          </el-checkbox>
        </div>
        <div class="action-buttons">
          <el-button type="primary" size="small" @click="goToExcludeWorker">
            断流Worker管理
          </el-button>
          <el-button type="primary" size="small" @click="refreshData">
            <i class="el-icon-refresh"></i> 刷新
          </el-button>
        </div>
      </div>

      <el-table 
        v-loading="loading" 
        :data="filteredEndpoints" 
        style="width: 100%"
        border
      >
        <el-table-column 
          label="Endpoint" 
          min-width="200"
          prop="metadata.name"
          sortable
        >
          <template #default="scope">
            <router-link 
              :to="`/serverless/${scope.row.metadata.name}`"
              class="endpoint-link"
            >
              {{ scope.row.metadata.name }}
            </router-link>
          </template>
        </el-table-column>

        <el-table-column 
          label="实例数" 
          width="80" 
          align="center"
          prop="workerCount"
          sortable
        >
          <template #default="scope">
            <span>{{ scope.row.workerCount || 0 }}</span>
          </template>
        </el-table-column>

        <el-table-column 
          label="GPU" 
          min-width="100" 
          align="center"
        >
          <template #default="scope">
            <span v-if="scope.row.spec.template.spec.containers[0]?.resources?.gpuRequests">
              {{ scope.row.spec.template.spec.containers[0].resources.gpuRequests.models[0] }}
              × {{ scope.row.spec.template.spec.containers[0].resources.gpuRequests.number }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="Cluster" width="80" align="center">
          <template #default="scope">
            <div class="cluster-tags">
              <template v-if="scope.row.spec.clusterIDs && scope.row.spec.clusterIDs.length">
                <el-tag 
                  v-for="clusterId in scope.row.spec.clusterIDs" 
                  :key="clusterId"
                  size="small"
                  class="cluster-tag"
                >
                  {{ clusterId }}
                </el-tag>
              </template>
              <span v-else>-</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="关联模型" 
          min-width="250"
        >
          <template #default="scope">
            <div class="model-tags">
              <template v-if="scope.row.models?.length">
                <el-tag
                  v-for="model in scope.row.models"
                  :key="model.name + '_' + model.source"
                  size="small"
                  :type="isPPIOModel(model) ? 'warning' : 'success'"
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
          prop="spec.template.spec.containers.0.image"
          sortable
        >
          <template #default="scope">
            <span class="image-text">
              {{ scope.row.spec.template.spec.containers[0]?.image || '-' }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="配置" width="160" align="center">
          <template #default="scope">
            <div class="yaml-buttons">
              <el-button
                type="primary"
                size="small"
                plain
                @click="showEndpointYaml(scope.row.metadata.name)"
              >
                YAML
              </el-button>
              <el-button
                type="info"
                size="small"
                plain
                @click="showSapYaml(scope.row.metadata.name)"
              >
                SAP
              </el-button>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" align="center">
          <template #default="scope">
            <div class="action-buttons">
              <el-button
                type="primary"
                size="small"
                @click="showScaleConfig(scope.row)"
              >
                扩缩容
              </el-button>
              <el-button
                type="primary"
                size="small"
                @click="showImageUpdate(scope.row)"
              >
                镜像
              </el-button>
              <el-button
                type="primary"
                size="small"
                @click="showMaxConcurrencyUpdate(scope.row)"
              >
                负载
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="yamlDialogVisible"
      :title="yamlDialogTitle"
      width="60%"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div class="yaml-content" v-loading="yamlLoading">
        <pre>{{ selectedYaml }}</pre>
        <el-button
          type="primary"
          size="small"
          class="copy-button"
          @click="copyYaml"
        >
          复制
        </el-button>
      </div>
    </el-dialog>

    <el-dialog
      v-model="scaleConfigVisible"
      title="扩缩容配置"
      width="500px"
      :close-on-click-modal="false"
    >
      <div class="scale-config-form">
        <el-form 
          ref="scaleConfigForm"
          :model="scaleConfig"
          :rules="scaleConfigRules"
          label-width="140px"
        >
          <el-form-item label="最小副本数" prop="minReplicas">
            <el-input-number 
              v-model="scaleConfig.minReplicas"
              :min="0"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="最大副本数" prop="maxReplicas">
            <el-input-number 
              v-model="scaleConfig.maxReplicas"
              :min="0"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="并发触发扩容数" prop="concurrencyPerWorker">
            <el-input-number 
              v-model="scaleConfig.concurrencyPerWorker"
              :min="1"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="扩容窗口时间(秒)" prop="scaleUpWindow">
            <el-input-number 
              v-model="scaleConfig.scaleUpWindow"
              :min="0"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="缩容窗口时间(秒)" prop="scaleDownWindow">
            <el-input-number 
              v-model="scaleConfig.scaleDownWindow"
              :min="0"
              controls-position="right"
            />
          </el-form-item>
        </el-form>
        <div class="dialog-footer">
          <el-button @click="scaleConfigVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmScaleConfig" :loading="updating">
            确认修改
          </el-button>
        </div>
      </div>
    </el-dialog>

    <el-dialog
      v-model="confirmDialogVisible"
      title="确认修改"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="confirm-content">
        <p>请确认以下修改内容：</p>
        <div class="confirm-item">
          <span class="label">最小副本数：</span>
          <span class="value">{{ scaleConfig.minReplicas }}</span>
        </div>
        <div class="confirm-item">
          <span class="label">最大副本数：</span>
          <span class="value">{{ scaleConfig.maxReplicas }}</span>
        </div>
        <div class="confirm-item">
          <span class="label">并发触发扩容数：</span>
          <span class="value">{{ scaleConfig.concurrencyPerWorker }}</span>
        </div>
        <div class="confirm-item">
          <span class="label">扩容窗口时间：</span>
          <span class="value">{{ scaleConfig.scaleUpWindow }}秒</span>
        </div>
        <div class="confirm-item">
          <span class="label">缩容窗口时间：</span>
          <span class="value">{{ scaleConfig.scaleDownWindow }}秒</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="confirmDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="updateScaleConfig" :loading="updating">
          确认
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="imageDialogVisible" title="更新镜像" width="800px">
      <el-form :model="imageForm" label-width="120px">
        <el-form-item label="新镜像地址" required>
          <el-input v-model="imageForm.image" placeholder="请输入新的镜像地址" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="imageDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUpdateImage">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="maxConcurrencyDialogVisible" title="更新最大负载" width="500px">
      <el-form :model="maxConcurrencyForm" label-width="120px">
        <el-form-item label="最大负载值" required>
          <el-input-number v-model="maxConcurrencyForm.maxConcurrency" :min="1" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="maxConcurrencyDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUpdateMaxConcurrency">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import { useServerlessStore } from '../stores/serverless'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const endpoints = ref([])
const searchQuery = ref('')
const imageSearchQuery = ref('')
const selectedCluster = ref('')
const selectedGpuType = ref('')
const showInstancesOnly = ref(false)
const clusters = ref([])
const gpuTypes = ref([])
const serverlessStore = useServerlessStore()

// 从 localStorage 恢复筛选条件
const restoreFilters = () => {
  const savedFilters = localStorage.getItem('serverlessFilters')
  if (savedFilters) {
    const filters = JSON.parse(savedFilters)
    searchQuery.value = filters.searchQuery || ''
    selectedCluster.value = filters.selectedCluster || ''
    selectedGpuType.value = filters.selectedGpuType || ''
    imageSearchQuery.value = filters.imageSearchQuery || ''
    showInstancesOnly.value = filters.showInstancesOnly || false
  }
}

// 保存筛选条件到 localStorage
const saveFilters = () => {
  const filters = {
    searchQuery: searchQuery.value,
    selectedCluster: selectedCluster.value,
    selectedGpuType: selectedGpuType.value,
    imageSearchQuery: imageSearchQuery.value,
    showInstancesOnly: showInstancesOnly.value
  }
  localStorage.setItem('serverlessFilters', JSON.stringify(filters))
}

// 监听筛选条件变化
watch([searchQuery, selectedCluster, selectedGpuType, imageSearchQuery, showInstancesOnly], () => {
  saveFilters()
}, { deep: true })

// YAML 相关的响应式变量
const yamlDialogVisible = ref(false)
const selectedYaml = ref('')
const yamlDialogTitle = ref('')
const yamlLoading = ref(false)

const scaleConfigVisible = ref(false)
const confirmDialogVisible = ref(false)
const updating = ref(false)
const currentEndpoint = ref(null)
const scaleConfigForm = ref(null)

// 扩缩容配置表单数据
const scaleConfig = ref({
  minReplicas: 0,
  maxReplicas: 1,
  concurrencyPerWorker: 1,
  scaleUpWindow: 0,
  scaleDownWindow: 0
})

// 表单验证规则
const scaleConfigRules = {
  minReplicas: [
    { required: true, message: '请输入最小副本数', trigger: 'blur' },
    { validator: validateMinReplicas, trigger: 'blur' }
  ],
  maxReplicas: [
    { required: true, message: '请输入最大副本数', trigger: 'blur' }
  ],
  concurrencyPerWorker: [
    { required: true, message: '请输入并发触发扩容数', trigger: 'blur' }
  ],
  scaleUpWindow: [
    { required: true, message: '请输入扩容窗口时间', trigger: 'blur' }
  ],
  scaleDownWindow: [
    { required: true, message: '请输入缩容窗口时间', trigger: 'blur' }
  ]
}

// 验证最小副本数不大于最大副本数
function validateMinReplicas(rule, value, callback) {
  if (value > scaleConfig.value.maxReplicas) {
    callback(new Error('最小副本数不能大于最大副本数'))
  } else {
    callback()
  }
}

// 过滤后的 endpoints
const filteredEndpoints = computed(() => {
  return endpoints.value.filter(endpoint => {
    // 实例数过滤
    const hasInstances = !showInstancesOnly.value || (endpoint.workerCount && endpoint.workerCount > 0)
    if (!hasInstances) return false

    // 搜索匹配（SE名称或关联模型）
    const searchTerm = searchQuery.value.toLowerCase()
    const matchesSearch = endpoint.metadata.name.toLowerCase().includes(searchTerm) ||
      endpoint.models?.some(model => model.name.toLowerCase().includes(searchTerm))

    // 镜像搜索匹配
    const imageSearchTerm = imageSearchQuery.value.toLowerCase()
    const matchesImage = !imageSearchTerm || 
      (endpoint.spec.template.spec.containers[0]?.image || '').toLowerCase().includes(imageSearchTerm)

    // 集群匹配
    const matchesCluster = !selectedCluster.value || 
      (endpoint.spec.clusterIDs && endpoint.spec.clusterIDs.includes(selectedCluster.value))

    // GPU型号匹配
    const matchesGpu = !selectedGpuType.value ||
      (endpoint.spec.template.spec.containers[0]?.resources?.gpuRequests?.models?.includes(selectedGpuType.value))
    
    return matchesSearch && matchesCluster && matchesGpu && matchesImage
  })
})

// 获取数据
const fetchData = async (forceRefresh = false) => {
  loading.value = true
  try {
    // 使用 store 获取数据
    endpoints.value = await serverlessStore.fetchEndpoints(forceRefresh)

    // 提取所有唯一的 cluster_id
    clusters.value = [...new Set(endpoints.value
      .filter(e => e.spec.clusterIDs)
      .flatMap(e => e.spec.clusterIDs)
    )].sort()

    // 提取所有唯一的 GPU 型号
    gpuTypes.value = [...new Set(endpoints.value
      .filter(e => e.spec.template.spec.containers[0]?.resources?.gpuRequests?.models)
      .flatMap(e => e.spec.template.spec.containers[0].resources.gpuRequests.models)
    )].sort()

  } catch (error) {
    console.error('Failed to fetch data:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  fetchData(true)
}

// 添加页面刷新时的处理
const handleBeforeUnload = () => {
  localStorage.removeItem('serverlessFilters')
}

// 从路由参数更新搜索框
const updateSearchFromRoute = () => {
  if (route.query.search) {
    searchQuery.value = route.query.search
  }
}

onMounted(() => {
  restoreFilters()
  updateSearchFromRoute()
  fetchData()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

// 监听路由变化
watch(() => route.query.search, (newSearch) => {
  if (newSearch) {
    searchQuery.value = newSearch
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

// 显示 Endpoint YAML
const showEndpointYaml = async (name) => {
  try {
    yamlLoading.value = true
    yamlDialogTitle.value = 'Endpoint YAML'
    yamlDialogVisible.value = true
    const response = await axios.get(`/se/${name}`)
    selectedYaml.value = response.data.data.yaml
  } catch (error) {
    console.error('获取 Endpoint YAML 失败:', error)
    ElMessage.error('获取 YAML 失败')
  } finally {
    yamlLoading.value = false
  }
}

// 显示 SAP YAML
const showSapYaml = async (name) => {
  try {
    yamlLoading.value = true
    yamlDialogTitle.value = 'SAP YAML'
    yamlDialogVisible.value = true
    const response = await axios.get(`/se/${name}`)
    selectedYaml.value = response.data.data.sapYaml
  } catch (error) {
    console.error('获取 SAP YAML 失败:', error)
    ElMessage.error('获取 YAML 失败')
  } finally {
    yamlLoading.value = false
  }
}

// 复制 YAML 内容
const copyYaml = async () => {
  try {
    await navigator.clipboard.writeText(selectedYaml.value)
    ElMessage.success('YAML 已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

// 显示扩缩容配置对话框
const showScaleConfig = async (endpoint) => {
  currentEndpoint.value = endpoint
  try {
    const response = await axios.get(`/se/${endpoint.metadata.name}`)
    const sapParam = response.data.data.sapParam
    if (sapParam) {
      scaleConfig.value = {
        minReplicas: sapParam.minReplicas,
        maxReplicas: sapParam.maxReplicas,
        concurrencyPerWorker: sapParam.concurrencyPerWorker,
        scaleUpWindow: sapParam.scaleUpWindow,
        scaleDownWindow: sapParam.scaleDownWindow
      }
    } else {
      // 如果没有 SAP 配置，使用默认值
      scaleConfig.value = {
        minReplicas: 0,
        maxReplicas: 1,
        concurrencyPerWorker: 1,
        scaleUpWindow: 0,
        scaleDownWindow: 0
      }
    }
    scaleConfigVisible.value = true
  } catch (error) {
    console.error('获取扩缩容配置失败:', error)
    ElMessage.error('获取配置失败')
  }
}

// 确认修改配置
const confirmScaleConfig = () => {
  scaleConfigForm.value.validate((valid) => {
    if (valid) {
      confirmDialogVisible.value = true
    }
  })
}

// 更新扩缩容配置
const updateScaleConfig = async () => {
  try {
    updating.value = true
    await axios.post('/sap/' + currentEndpoint.value.metadata.name, scaleConfig.value, {
      params: {
        endpoint: currentEndpoint.value.metadata.name
      }
    })
    ElMessage.success('更新成功')
    confirmDialogVisible.value = false
    scaleConfigVisible.value = false
    refreshData()
  } catch (error) {
    console.error('更新扩缩容配置失败:', error)
    ElMessage.error(error.response?.data?.message || '更新失败')
  } finally {
    updating.value = false
  }
}

// 判断是否是 PPIO 模型
const isPPIOModel = (model) => {
  return model.source === 'ppio'
}

// 打开模型详情
const openModelDetail = (model) => {
  const platform = model.source === 'ppio' ? 'ppio' : ''
  const baseUrl = platform ? '/ppio-models' : '/models'
  router.push({
    name: 'EndpointManager',
    params: { modelName: model.name },
    query: platform ? { platform } : undefined
  })
}

// 跳转到断流worker管理页面
const goToExcludeWorker = () => {
  router.push('/exclude-worker')
}

// 镜像更新相关
const imageDialogVisible = ref(false)
const imageForm = ref({
  image: '',
  endpoint: null
})

const showImageUpdate = (endpoint) => {
  imageForm.value.endpoint = endpoint
  imageForm.value.image = endpoint.image // 预填充当前镜像
  imageDialogVisible.value = true
}

const handleUpdateImage = async () => {
  try {
    const response = await axios.put(`/se/${imageForm.value.endpoint.metadata.name}/image`, {
      image: imageForm.value.image
    })
    if (response.status === 200) {
      ElMessage.success('更新镜像成功')
      imageDialogVisible.value = false
      fetchData()
    }
  } catch (error) {
    console.error('更新镜像失败:', error)
    const errorMsg = error.response?.data?.message || error.message || '未知错误'
    ElMessage.error(`更新镜像失败: ${errorMsg}`)
  }
}

// 最大负载更新相关
const maxConcurrencyDialogVisible = ref(false)
const maxConcurrencyForm = ref({
  maxConcurrency: 1,
  endpoint: null
})

const showMaxConcurrencyUpdate = (endpoint) => {
  maxConcurrencyForm.value.endpoint = endpoint
  maxConcurrencyForm.value.maxConcurrency = endpoint.maxConcurrency // 预填充当前最大负载
  maxConcurrencyDialogVisible.value = true
}

const handleUpdateMaxConcurrency = async () => {
  try {
    const response = await axios.put(`/se/${maxConcurrencyForm.value.endpoint.metadata.name}/maxConcurrency`, {
      maxConcurrency: maxConcurrencyForm.value.maxConcurrency
    })
    if (response.status === 200) {
      ElMessage.success('更新最大负载成功')
      maxConcurrencyDialogVisible.value = false
      fetchData()
    }
  } catch (error) {
    console.error('更新最大负载失败:', error)
    const errorMsg = error.response?.data?.message || error.message || '未知错误'
    ElMessage.error(`更新最大负载失败: ${errorMsg}`)
  }
}

// 跳转到操作日志页面
const goToAuditLog = () => {
  // 从 localStorage 获取用户信息
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
  if (!userInfo.isAdmin) {
    ElMessage.error('只有管理员可以查看操作日志')
    return
  }
  router.push('/audit-log')
}

// 在 script setup 部分添加计算属性
const uniqueImages = computed(() => {
  return [...new Set(endpoints.value
    .map(e => e.spec.template.spec.containers[0]?.image)
    .filter(Boolean)
  )].sort()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
}

.header {
  margin-bottom: 32px;
  padding: 0 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 28px;
  color: #3b82f6;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0;
  padding-bottom: 8px;
}

.content-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.table-filters {
  display: flex;
  gap: 16px;
  align-items: center;
}

.filter-item {
  min-width: 350px;
}

.search-input {
  width: 600px;
}

.image-search-input {
  width: 300px;
}

.cluster-filter {
  width: 80px;
}

.gpu-filter {
  width: 250px;
}

.instance-filter {
  width: 100px;
}

.endpoint-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 500;
  
  &:hover {
    text-decoration: underline;
  }
}

.gpu-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.gpu-spec {
  font-family: monospace;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 13px;
}

.loading-spinner {
  color: #3b82f6;
}

.cluster-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  justify-content: center;
}

.cluster-tag {
  margin: 2px;
}

.yaml-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.yaml-content {
  position: relative;
  background: #1e1e1e;
  border-radius: 8px;
  padding: 16px;
  margin: 0;
  min-height: 200px;
}

.yaml-content pre {
  margin: 0;
  padding: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: #e6e6e6;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.copy-button {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0.8;
}

.copy-button:hover {
  opacity: 1;
}

.scale-config-form {
  padding: 20px;
}

.dialog-footer {
  margin-top: 24px;
  text-align: right;
}

.confirm-content {
  padding: 16px;
}

.confirm-item {
  margin: 12px 0;
  display: flex;
  align-items: center;
}

.confirm-item .label {
  width: 140px;
  color: #606266;
}

.confirm-item .value {
  color: #333;
  font-weight: 500;
}

.model-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.model-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.model-tag:hover {
  opacity: 0.8;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.action-buttons {
  display: flex;
  gap: 4px;
  justify-content: center;
  flex-wrap: nowrap;
}

.action-buttons .el-button {
  padding: 8px 12px;
  flex-shrink: 0;
}

.header-right {
  display: flex;
  gap: 12px;
}

.image-text {
  font-family: monospace;
  word-break: break-all;
  display: block;
  line-height: 1.4;
}
</style> 