<template>
  <div class="page-container">
    <div class="header">
      <div class="header-left">
        <el-button @click="$router.back()" class="back-button">
          <i class="el-icon-arrow-left"></i> 返回
        </el-button>
        <div class="title-section">
          <h2 class="page-title">
            Worker 列表
            <span class="endpoint-id">{{ endpointId }}</span>
          </h2>
        </div>
      </div>
      <el-button type="primary" @click="refreshData">
        <i class="el-icon-refresh"></i> 刷新
      </el-button>
    </div>

    <div class="content-card">
      <el-table 
        v-loading="loading" 
        :data="workers"
        style="width: 100%"
        :default-sort="{ prop: 'metadata.creationTimestamp', order: 'descending' }"
      >
        <el-table-column prop="metadata.name" label="Worker ID" min-width="200" />
        
        <el-table-column 
          label="状态" 
          width="120" 
          align="center"
          sortable
          :sort-method="sortByState"
          prop="status.state"
        >
          <template #default="scope">
            <el-tag 
              :type="getStateType(scope.row.status.state)"
              size="small"
            >
              {{ scope.row.status.state }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="健康状态" width="120" align="center">
          <template #default="scope">
            <el-tag 
              :type="scope.row.status.healthy ? 'success' : 'danger'"
              size="small"
            >
              {{ scope.row.status.healthy ? '健康' : '不健康' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="Instance" min-width="200">
          <template #default="scope">
            <template v-if="scope.row.status.realInstanceID">
              <el-button 
                type="text" 
                @click="viewInstance(scope.row, scope.row.status.realInstanceID)"
                class="instance-link"
              >
                {{ scope.row.status.realInstanceID }}
              </el-button>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column 
          label="创建时间" 
          min-width="180"
          sortable
          prop="metadata.creationTimestamp"
        >
          <template #default="scope">
            {{ formatTime(scope.row.metadata.creationTimestamp) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" align="center">
          <template #default="scope">
            <div class="action-buttons">
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(scope.row)"
                :loading="deletingWorker === scope.row.metadata.name"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="instanceDialogVisible"
      title="实例详情"
      width="800px"
      destroy-on-close
    >
      <div v-loading="instanceLoading">
        <div v-if="instanceInfo" class="instance-detail">
          <div class="detail-card">
            <div class="section">
              <div class="section-title">基本信息</div>
              <div class="info-section">
                <div class="info-row">
                  <span class="info-label">实例 ID:</span>
                  <span class="info-value">{{ instanceInfo.id }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">节点 ID:</span>
                  <span class="info-value">{{ instanceInfo.nodeId }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">状态:</span>
                  <span class="info-value">
                    <el-tag size="small" :type="instanceInfo.state?.state === 'running' ? 'success' : 'warning'">
                      {{ instanceInfo.state?.state }}
                    </el-tag>
                  </span>
                </div>
                <div class="info-row">
                  <span class="info-label">实例类型:</span>
                  <span class="info-value">{{ instanceInfo.kind }}</span>
                </div>
              </div>
            </div>

            <div class="section">
              <div class="section-title">镜像信息</div>
              <div class="info-section">
                <div class="info-row">
                  <span class="info-value monospace">{{ instanceInfo.containers?.image }}</span>
                </div>
              </div>
            </div>

            <div class="section">
              <div class="section-title">资源配置</div>
              <div class="info-section">
                <div class="info-row">
                  <span class="info-label">CPU 核心数:</span>
                  <span class="info-value">{{ instanceInfo.containers?.cpuNum }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">内存:</span>
                  <span class="info-value">{{ (instanceInfo.containers?.memorySize / 1024 / 1024 / 1024).toFixed(1) }}GB</span>
                </div>
                <div class="info-row">
                  <span class="info-label">系统盘:</span>
                  <span class="info-value">{{ (instanceInfo.containers?.rootfsSize / 1024 / 1024 / 1024).toFixed(1) }}GB</span>
                </div>
                <div class="info-row" v-if="instanceInfo.containers?.gpuNum">
                  <span class="info-label">GPU:</span>
                  <span class="info-value">{{ instanceInfo.containers.gpuProductName }} × {{ instanceInfo.containers.gpuNum }}</span>
                </div>
              </div>
            </div>

            <div class="section">
              <div class="section-title">网络访问</div>
              <div class="info-section">
                <div class="info-row">
                  <span class="info-label">端口映射:</span>
                  <span class="info-value">
                    <template v-if="instanceInfo.exposePortStates?.length">
                      <div v-for="port in instanceInfo.exposePortStates" :key="port.port">
                        {{ port.exposeAddress }}
                        <el-tag size="small" :type="port.online ? 'success' : 'warning'" class="port-status">
                          {{ port.online ? '在线' : '离线' }}
                        </el-tag>
                      </div>
                    </template>
                    <span v-else>-</span>
                  </span>
                </div>
                <div class="info-row">
                  <span class="info-label">SSH 命令:</span>
                  <span class="info-value monospace">{{ instanceInfo.apps?.sshState?.sshCommand }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">SSH 密码:</span>
                  <span class="info-value monospace">{{ instanceInfo.apps?.sshState?.sshPassword }}</span>
                </div>
              </div>
            </div>

            <div class="section">
              <div class="section-title">日志信息</div>
              <div class="info-section">
                <div class="info-row">
                  <span class="info-label">系统日志:</span>
                  <a :href="instanceInfo.apps?.logState?.systemLogAddress" target="_blank" class="info-value link">
                    {{ instanceInfo.apps?.logState?.systemLogAddress }}
                  </a>
                </div>
                <div class="info-row">
                  <span class="info-label">实例日志:</span>
                  <a :href="`${instanceInfo.apps?.logState?.instanceLogAddress}?follow=1&tail=100`" target="_blank" class="info-value link">
                    {{ instanceInfo.apps?.logState?.instanceLogAddress }}
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const endpointId = ref(route.params.endpointId)
const loading = ref(false)
const workers = ref([])

const instanceDialogVisible = ref(false)
const instanceLoading = ref(false)
const instanceInfo = ref(null)

const deletingWorker = ref('')

const getStateType = (state) => {
  const types = {
    'Running': 'success',
    'Pending': 'warning',
    'Failed': 'danger'
  }
  return types[state] || 'info'
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

const canDelete = (worker) => {
  return worker.status.state === 'Running'
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/workers?endpoint=${endpointId.value}`)
    workers.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch workers:', error)
    ElMessage.error('获取 worker 列表失败')
  } finally {
    loading.value = false
  }
}

const handleDelete = async (worker) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Worker ${worker.metadata.name} 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    deletingWorker.value = worker.metadata.name
    await axios.delete(`/workers/${worker.metadata.name}`)
    ElMessage.success('删除成功')
    await refreshData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 Worker 失败:', error)
      const errorMessage = error.response?.data?.message || error.response?.data?.reason || error.message || '删除失败'
      ElMessage.error(errorMessage)
    }
  } finally {
    deletingWorker.value = ''
  }
}

const refreshData = () => {
  fetchData()
}

const viewInstance = async (worker, instanceID) => {
  try {
    instanceLoading.value = true
    instanceDialogVisible.value = true
    const clusterIndex = worker.spec?.clusterIDs?.[0] || '1'
    const response = await axios.get('/nexusclusters/instances/detail', {
      params: {
        cluster_index: clusterIndex,
        instance_id: instanceID,
        namespace: 'default'
      }
    })
    instanceInfo.value = response.data.data
  } catch (error) {
    console.error('Failed to view instance:', error)
    ElMessage.error('查看实例失败')
  } finally {
    instanceLoading.value = false
  }
}

// 状态排序的优先级映射
const stateOrder = {
  'Running': 1,
  'Pending': 2,
  'Failed': 3
}

// 自定义状态排序方法
const sortByState = (a, b) => {
  const stateA = stateOrder[a.status.state] || 999
  const stateB = stateOrder[b.status.state] || 999
  return stateA - stateB
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-button {
  padding: 8px 16px;
}

.title-section {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.endpoint-id {
  color: #3b82f6;
  font-family: monospace;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 16px;
}

.content-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.instance-link {
  color: #3b82f6;
  text-decoration: none;
  
  &:hover {
    text-decoration: underline;
  }
}

.detail-section {
  margin-bottom: 24px;

  h3 {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 16px;
    color: #374151;
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label {
  font-size: 13px;
  color: #6b7280;
}

.value {
  font-size: 14px;
  color: #111827;
}

.image-info {
  font-family: monospace;
  background: #f3f4f6;
  padding: 12px;
  border-radius: 4px;
  font-size: 14px;
  word-break: break-all;
}

.detail-card {
  background-color: #f9fafb;
  border-radius: 8px;
}

.section {
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.section:last-child {
  border-bottom: none;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 16px;
}

.info-section {
  margin-bottom: 0;
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  line-height: 1.5;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  width: 100px;
  font-size: 13px;
  color: #6b7280;
}

.info-value {
  font-size: 14px;
  color: #111827;
  flex: 1;
}

.monospace {
  font-family: monospace;
  font-size: 14px;
  word-break: break-all;
  color: #374151;
}

.port-status {
  margin-left: 8px;
}

.link {
  color: #3b82f6;
  text-decoration: none;
  
  &:hover {
    text-decoration: underline;
  }
}
</style>