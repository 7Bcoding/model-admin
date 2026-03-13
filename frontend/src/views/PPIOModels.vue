<template>
  <div class="page-container">
    <div class="header">
      <div class="header-left">
        <div class="title-section">
          <i class="el-icon-cpu title-icon"></i>
          <h2 class="page-title">模型管理(PPIO)</h2>
        </div>
      </div>
      <div class="filters">
        <div class="search-wrapper">
          <i class="el-icon-search search-icon"></i>
          <input 
            v-model="search" 
            type="text" 
            class="search-input" 
            placeholder="搜索模型名称或备注"
          >
        </div>
        <div class="filter-options">
          <label class="filter-checkbox">
            <i class="el-icon-star-on checkbox-icon"></i>
            <input type="checkbox" v-model="showStarred">
            <span class="checkbox-text">只看收藏</span>
          </label>
          <label class="filter-checkbox">
            <i class="el-icon-check checkbox-icon"></i>
            <input type="checkbox" v-model="showActive" checked>
            <span class="checkbox-text">只看启用中</span>
          </label>
        </div>
      </div>
    </div>

    <div class="content-card">
      <table class="data-table">
        <thead>
          <tr>
            <th><i class="el-icon-document"></i> 模型名称</th>
            <th><i class="el-icon-loading"></i> 状态</th>
            <th><i class="el-icon-lock"></i> 私有</th>
            <th><i class="el-icon-data-line"></i> 最大长度</th>
            <th><i class="el-icon-connection"></i> EP数</th>
            <th><i class="el-icon-monitor"></i> 监控</th>
            <th><i class="el-icon-setting"></i> 操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="model in filteredModels" :key="model.model_name">
            <td class="model-cell">
              <div class="model-name">
                <div class="model-info">
                  <span class="star" @click.stop="toggleStar(model)">
                    {{ model.starred ? '★' : '☆' }}
                  </span>
                  <div class="name-container">
                    <span class="name clickable" @click="showEndpoints(model)">{{ model.model_name }}</span>
                    <span 
                      v-if="model.open_chat_id"
                      class="msg-link" 
                      @click.stop="openFeishuChat(model.open_chat_id)"
                    >msg</span>
                  </div>
                  <span v-if="model.note" class="model-note">{{ model.note }}</span>
                </div>
                <div class="note-button-wrapper">
                  <el-button
                    type="text"
                    size="small"
                    class="note-button"
                    @click="showNoteDialog(model)"
                  >
                    {{ model.note ? '编辑' : '备注' }}
                  </el-button>
                </div>
              </div>
            </td>
            <td>
              <span class="badge" :class="getStatusClass(model.status)">
                {{ getStatusText(model.status) }}
              </span>
            </td>
            <td>
              <span class="badge" :class="model.private ? 'bg-warning' : 'bg-info'">
                {{ model.private ? '是' : '否' }}
              </span>
            </td>
            <td class="text-center">{{ model.max_tokens || '-' }}</td>
            <td class="text-center">{{ model.endpoints?.length || 0 }}</td>
            <td>
              <div class="monitor-buttons">
                <button 
                  class="monitor-btn-overview" 
                  @click="openMonitor(model, 'overview')"
                >SLA</button>
                <button 
                  class="monitor-btn-vllm"
                  @click="openMonitor(model, 'vllm')"
                >{{ model.inference_engine || 'vllm' }}</button>
                <button 
                  class="monitor-btn-serverless"
                  @click="openMonitor(model, 'serverless')"
                >serverless</button>
              </div>
            </td>
            <td>
              <div class="action-buttons">
                <button 
                  class="btn btn-sm btn-primary" 
                  @click="showDetails(model)"
                >详情</button>
                <button 
                  class="btn btn-sm btn-info"
                  @click="showDeploymentInfo(model)"
                >部署信息</button>
                <button 
                  class="btn btn-sm btn-warning"
                  @click="showEndpoints(model)"
                >endpoints</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 添加模型详情模态框 -->
    <div class="modal fade" id="modelDetailsModal" tabindex="-1">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">模型详情</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <div class="basic-info mb-4">
              <h6>基本信息</h6>
              <div class="row">
                <div class="col-md-6">
                  <div class="mb-3">
                    <label class="fw-bold">模型名称</label>
                    <div>{{ selectedModel?.model_name }}</div>
                  </div>
                  <div class="mb-3">
                    <label class="fw-bold">状态</label>
                    <div>
                      <span class="badge" :class="getStatusClass(selectedModel?.status)">
                        {{ getStatusText(selectedModel?.status) }}
                      </span>
                    </div>
                  </div>
                  <div class="mb-3">
                    <label class="fw-bold">最大长度</label>
                    <div>{{ selectedModel?.max_tokens }}</div>
                  </div>
                  <div class="mb-3">
                    <label class="fw-bold">私有模型</label>
                    <div>
                      <span class="badge" :class="selectedModel?.private ? 'bg-warning' : 'bg-info'">
                        {{ selectedModel?.private ? '是' : '否' }}
                      </span>
                      <span v-if="selectedModel?.private && selectedModel?.model?.user_email" class="ms-2 text-muted">
                        ({{ selectedModel.model.user_email }})
                      </span>
                    </div>
                  </div>
                </div>
                <div class="col-md-6">
                  <div class="mb-3">
                    <label class="fw-bold">输入价格</label>
                    <div>{{ selectedModel?.input_token_price || 0 }}/tokens</div>
                  </div>
                  <div class="mb-3">
                    <label class="fw-bold">输出价格</label>
                    <div>{{ selectedModel?.output_token_price || 0 }}/百万 tokens</div>
                  </div>
                  <div class="mb-3" v-if="selectedModel?.tags?.length">
                    <label class="fw-bold">标签</label>
                    <div class="tags">
                      <span v-for="tag in selectedModel.tags" :key="tag" class="tag">{{ tag }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="description mb-4" v-if="selectedModel?.description">
              <h6>模型描述</h6>
              <p>{{ selectedModel.description }}</p>
            </div>

            <div class="raw-data">
              <h6 class="d-flex justify-content-between align-items-center">
                原始数据
                <button class="btn btn-sm btn-outline-secondary" @click="toggleRawData">
                  {{ showRawData ? '隐藏' : '显示' }}
                </button>
              </h6>
              <pre v-if="showRawData" class="bg-light p-3 rounded"><code>{{ JSON.stringify(selectedModel?.model, null, 2) }}</code></pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加备注对话框 -->
    <div class="modal fade" id="noteDialog" tabindex="-1" aria-labelledby="noteDialogLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="noteDialogLabel">编辑备注</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label for="noteInput" class="form-label">备注</label>
              <textarea 
                id="noteInput"
                v-model="noteInput"
                class="form-control"
                rows="3"
                placeholder="请输入备注内容"
              ></textarea>
            </div>
            <div class="mb-3">
              <label for="openChatIdInput" class="form-label">飞书群组</label>
              <input 
                id="openChatIdInput"
                v-model="openChatIdInput"
                class="form-control"
                placeholder="请输入飞书群组ID"
              >
            </div>
            <div class="mb-3">
              <label for="inferenceEngineInput" class="form-label">推理引擎</label>
              <input 
                id="inferenceEngineInput"
                v-model="inferenceEngineInput"
                class="form-control"
                placeholder="请输入推理引擎"
              >
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
            <button type="button" class="btn btn-primary" @click="saveNote">保存</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { usePPIOModelsStore } from '../stores/ppioModels'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { starModel, unstarModel, updateModelNote } from '../api/models'

const router = useRouter()
const ppioModelsStore = usePPIOModelsStore()
const loading = ref(false)
const search = ref('')
const showStarred = ref(false)
const showActive = ref(true)
const selectedModel = ref(null)
const showRawData = ref(false)
const noteInput = ref('')
const openChatIdInput = ref('')
const inferenceEngineInput = ref('')
const noteDialogVisible = ref(false)
let modelDetailsModal = null

// 获取模型列表
const fetchModels = async (forceRefresh = false) => {
  loading.value = true
  try {
    await ppioModelsStore.fetchModels(forceRefresh)
  } catch (error) {
    ElMessage.error('获取模型列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 过滤后的模型列表
const filteredModels = computed(() => {
  let models = ppioModelsStore.getModels
  
  if (!models || !Array.isArray(models)) {
    console.warn('Models is not an array:', models)
    return []
  }
  
  return models.filter(model => {
    if (!model || typeof model !== 'object') {
      console.warn('Invalid model object:', model)
      return false
    }
    
    // 收藏过滤
    if (showStarred.value && !model.starred) return false
    
    // 状态过滤
    if (showActive.value && model.status !== 'MODEL_STATUS_SERVING') {
      return false
    }
    
    // 搜索过滤
    if (search.value) {
      const searchTerm = search.value.toLowerCase()
      const modelName = String(model.model_name || '').toLowerCase()
      const note = String(model.note || '').toLowerCase()
      const searchTerms = searchTerm.split(/\s+/).filter(term => term)
      
      return searchTerms.every(term => 
        modelName.includes(term) || note.includes(term)
      )
    }
    
    return true
  }).sort((a, b) => {
    // 首先按照rank排序
    if (a.rank !== b.rank) {
      return a.rank - b.rank
    }
    // 然后按照starred排序
    if (a.starred !== b.starred) {
      return a.starred ? -1 : 1
    }
    // 最后按照名称排序
    return a.model_name.localeCompare(b.model_name)
  })
})

// 获取状态样式
const getStatusClass = (status) => {
  switch (status) {
    case 'MODEL_STATUS_SERVING':
      return 'bg-success'
    case 'MODEL_STATUS_DELETED':
      return 'bg-danger'
    case 'MODEL_STATUS_MAINTAINING':
      return 'bg-warning'
    default:
      return 'bg-secondary'
  }
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'MODEL_STATUS_SERVING': '在线',
    'MODEL_STATUS_DELETED': '已删除',
    'MODEL_STATUS_MAINTAINING': '维护中'
  }
  return statusMap[status] || status
}

// 切换收藏状态
const toggleStar = async (model) => {
  try {
    if (model.starred) {
      await unstarModel(model.model_name, 'ppio')
      ElMessage.success('取消收藏成功')
    } else {
      await starModel(model.model_name, 'ppio')
      ElMessage.success('收藏成功')
    }
    await fetchModels(true)
  } catch (error) {
    console.error('Star toggle error:', error)
    ElMessage.error(model.starred ? '取消收藏失败' : '收藏失败')
  }
}

// 显示端点信息
const showEndpoints = (model) => {
  router.push({
    name: 'EndpointManager',
    params: { modelName: model.model_name },
    query: { platform: 'ppio' }
  })
}

// 显示部署信息
const showDeploymentInfo = (model) => {
  router.push({
    name: 'ModelApi',
    query: { 
      modelName: model.model_name,
      platform: 'ppio'
    }
  })
}

// 打开飞书群组
const openFeishuChat = (openChatId) => {
  if (!openChatId) {
    ElMessage.warning('未配置飞书群组')
    return
  }
  const url = `https://applink.feishu.cn/client/chat/open?openChatId=${openChatId}`
  window.open(url, '_blank')
}

// 显示备注对话框
const showNoteDialog = (model) => {
  selectedModel.value = model
  noteInput.value = model.note || ''
  openChatIdInput.value = model.open_chat_id || ''
  inferenceEngineInput.value = model.inference_engine || ''
  
  // 初始化并显示 Modal
  const modalEl = document.getElementById('noteDialog')
  if (modalEl && window.bootstrap) {
    const modal = new window.bootstrap.Modal(modalEl)
    modal.show()
  } else {
    console.error('Bootstrap Modal initialization failed')
  }
}

// 保存备注
const saveNote = async () => {
  try {
    await updateModelNote(
      selectedModel.value.model_name, 
      noteInput.value, 
      openChatIdInput.value,
      inferenceEngineInput.value,
      'ppio'
    )

    // 更新本地数据
    selectedModel.value.note = noteInput.value
    selectedModel.value.open_chat_id = openChatIdInput.value
    selectedModel.value.inference_engine = inferenceEngineInput.value
    
    ElMessage({
      message: '保存成功',
      type: 'success'
    })
    
    // 关闭 Modal
    const modalEl = document.getElementById('noteDialog')
    if (modalEl && window.bootstrap) {
      const modal = window.bootstrap.Modal.getInstance(modalEl)
      if (modal) {
        modal.hide()
      }
    }
  } catch (error) {
    console.error('Error saving note:', error)
    ElMessage({
      message: error.response?.data?.message || '保存失败',
      type: 'error'
    })
  }
}

// 显示模型详情
const showDetails = (model) => {
  selectedModel.value = model
  if (!modelDetailsModal) {
    const modalEl = document.getElementById('modelDetailsModal')
    if (modalEl && window.bootstrap) {
      modelDetailsModal = new window.bootstrap.Modal(modalEl)
      modalEl.addEventListener('hidden.bs.modal', () => {
        selectedModel.value = null
        showRawData.value = false
      })
    } else {
      console.error('Bootstrap Modal initialization failed')
      return
    }
  }
  modelDetailsModal.show()
}

// 切换原始数据显示
const toggleRawData = () => {
  showRawData.value = !showRawData.value
}

// 打开监控
const openMonitor = (model, type) => {
  if (!model || !model.model_name) {
    ElMessage.warning('模型信息不完整')
    return
  }

  const monitorUrls = {
    overview: `https://grafana.aicloud.paigod.work/d/a9131a24-3266-4b18-8aae-b5fc2ec45adc/ai-cloud-platform-sla-metrics?orgId=1&from=now-12h&to=now&orgId=1&var-CustomerUUID=ALL&var-Model=${encodeURIComponent(model.model_name)}&refresh=5s&var-Customer=All`,
    vllm: model.inference_engine === 'sglang' ? 
      `https://grafana.aicloud.pplabs.tech/d/ddyp55uq7brpcc/sglang-dashboard?orgId=1&refresh=5s&var-model_id=${encodeURIComponent(getSglangModelId(model.model_name))}` :
      `https://grafana.aicloud.pplabs.tech/d/b281712d-8bff-41ef-9f3f-71ad43c05e9bdd/vllm-per-model?orgId=1&var-model_name=${encodeURIComponent(model.model_name)}`,
    serverless: model.endpoints ? 
      `https://grafana.aicloud.pplabs.tech/d/KOpAVRCIz1/prod-overview?orgId=1${
        model.endpoints
          .filter(endpoint => endpoint.url.includes('.runsync.novita.dev'))
          .map(endpoint => `&var-endpoint=${endpoint.url.replace('https://', '').replace('http://', '').split('.runsync.novita.dev')[0]}`)
          .join('')
      }&from=now-24h&to=now&refresh=5s` : null
  }
  
  const url = monitorUrls[type]
  if (url) {
    window.open(url, '_blank')
  } else {
    ElMessage.warning('无法获取监控地址')
  }
}

const getSglangModelId = (modelName) => {
  const modelMapping = {
    'deepseek/deepseek-r1': 'deepseek-r1',
    'deepseek/deepseek-r1/community': 'deepseek-r1-community',
    'deepseek/deepseek-v3': 'deepseek-v3',
    'deepseek/deepseek-v3/community': 'deepseek-v3-community'
  }
  return modelMapping[modelName] || modelName
}

// 页面挂载时从 localStorage 恢复搜索条件
onMounted(async () => {
  // 恢复搜索框内容
  const savedSearch = localStorage.getItem('ppio_models_search_query')
  if (savedSearch !== null) {
    search.value = savedSearch
  }
  // 恢复过滤选项
  const savedShowStarred = localStorage.getItem('ppio_models_show_starred')
  if (savedShowStarred !== null) {
    showStarred.value = savedShowStarred === 'true'
  }
  const savedShowActive = localStorage.getItem('ppio_models_show_active')
  if (savedShowActive !== null) {
    showActive.value = savedShowActive === 'true'
  }
  await fetchModels()
})

// 监听搜索条件变化，保存到 localStorage
watch(search, (val) => {
  localStorage.setItem('ppio_models_search_query', val)
})
watch(showStarred, (val) => {
  localStorage.setItem('ppio_models_show_starred', val)
})
watch(showActive, (val) => {
  localStorage.setItem('ppio_models_show_active', val)
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  margin: 0 auto;
  background: #f5f5f5;
  min-height: calc(100vh - var(--navbar-height));
}

.header {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-bottom: 32px;
  padding: 0 12px;
  gap: 48px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  position: relative;
  margin: 0;
  padding-bottom: 8px;
}

.page-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 40px;
  height: 3px;
  background: #3b82f6;
  border-radius: 2px;
}

.filters {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-wrapper {
  position: relative;
  width: 480px;
  transition: all 0.3s ease;
  background: white;
  border-radius: 8px;
  box-shadow: 
    0 1px 2px rgba(0, 0, 0, 0.1),
    inset 0 -1px 0 rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.search-wrapper:focus-within {
  transform: translateY(-1px);
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.05);
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 16px;
}

.search-input {
  width: 100%;
  padding: 10px 16px 10px 40px;
  border: 1px solid transparent;
  border-radius: 8px;
  font-size: 14px;
  background-color: transparent;
  transition: all 0.3s ease;
  height: 42px;
}

.search-input:focus {
  border-color: transparent;
  box-shadow: 
    0 0 0 2px #3b82f6,
    0 4px 8px -2px rgba(59, 130, 246, 0.2);
  outline: none;
}

.filter-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.checkbox-text {
  color: #4b5563;
  font-size: 14px;
  font-weight: 500;
}

.content-card {
  background: white;
  border-radius: 12px;
  box-shadow: 
    0 1px 3px rgba(0, 0, 0, 0.1),
    0 4px 6px -1px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  position: relative;
  border: 1px solid #e5e7eb;
}

.data-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
}

.data-table th {
  padding: 8px 16px;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  font-size: 13px;
  font-weight: 600;
  color: #4b5563;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: inset 0 -1px 0 #e5e7eb;
  text-align: center;
}

.data-table th i {
  margin-right: 6px;
  font-size: 16px;
  vertical-align: middle;
  color: #6b7280;
}

.data-table th:nth-child(5) {
  font-weight: 600;
  width: 80px;
  text-align: center;
}

.data-table td {
  padding: 8px 16px;
  border-bottom: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.data-table td:nth-child(4),
.data-table td:nth-child(5) {
  font-family: 'SF Mono', SFMono-Regular, ui-monospace, Menlo, Monaco, Consolas, monospace;
  font-weight: 600;
  font-size: 14px;
  color: #1f2937;
  text-align: center;
  padding-right: 16px;
  padding-left: 16px;
}

.data-table tr:hover td {
  background-color: rgba(59, 130, 246, 0.02);
  box-shadow: 
    inset 1px 0 0 #f3f4f6,
    inset -1px 0 0 #f3f4f6;
}

.badge {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-width: unset;
  padding: 6px 8px;
}

.bg-success {
  background: #10b981;
  color: white;
  border: none;
}

.bg-warning {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid rgba(251, 191, 36, 0.4);
}

.bg-info {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid rgba(96, 165, 250, 0.4);
}

.bg-danger {
  background: #fee2e2;
  color: #b91c1c;
  border: 1px solid rgba(239, 68, 68, 0.4);
}

.star {
  cursor: pointer;
  color: #ffd700;
  font-size: 1.2em;
  transition: transform 0.2s;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.star:hover {
  transform: scale(1.1);
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
}

.model-note {
  color: #6b7280;
  font-weight: 400;
  font-size: 0.85em;
  margin-bottom: 1px;
  margin-left: 24px;
  font-style: italic;
  opacity: 0.9;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  position: relative;
  height: 32px;
  overflow: hidden;
}

.btn::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0));
  opacity: 0;
  transition: opacity 0.3s;
}

.btn:hover::after {
  opacity: 1;
}

.btn-primary {
  background: #3b82f6;
  color: white;
  border: none;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-info {
  background: #ec4899;
  color: white;
  border: none;
}

.btn-info:hover {
  background: #db2777;
}

.btn-warning {
  background: #8b5cf6;
  color: white;
  border: none;
}

.btn-warning:hover {
  background: #7c3aed;
}

.model-name {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 8px;
}

.model-info {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  flex: 1;
}

.name {
  font-weight: 700;
  color: #303133;
  font-size: 16px;
  letter-spacing: -0.2px;
}

.name.clickable {
  cursor: pointer;
}

.name.clickable:hover {
  opacity: 0.8;
}

.note-button-wrapper {
  margin-left: auto;
  opacity: 0;
  transition: opacity 0.2s;
}

.model-cell:hover .note-button-wrapper {
  opacity: 1;
}

.note-button {
  padding: 2px 4px;
  height: auto;
  font-size: 12px;
  color: #909399;
}

.note-button:hover {
  color: #409EFF;
}

.monitor-buttons,
.action-buttons {
  display: flex;
  gap: 8px;
}

.monitor-btn-overview,
.monitor-btn-vllm,
.monitor-btn-serverless {
  padding: 4px 12px;
  border-radius: 2px;
  font-size: 12px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.monitor-btn-overview {
  background-color: #e6f4ff;
  border-color: #91caff;
  color: #0958d9;
}

.monitor-btn-overview:hover {
  background-color: #bae0ff;
  border-color: #0958d9;
}

.monitor-btn-vllm {
  background-color: #fff1f0;
  border-color: #ffccc7;
  color: #cf1322;
}

.monitor-btn-vllm:hover {
  background-color: #ffa39e;
  border-color: #cf1322;
}

.monitor-btn-serverless {
  background-color: #f6ffed;
  border-color: #b7eb8f;
  color: #389e0d;
}

.monitor-btn-serverless:hover {
  background-color: #d9f7be;
  border-color: #389e0d;
}

.name-container {
  display: flex;
  align-items: center;
}

.msg-link {
  font-size: 12px;
  color: #6366f1;
  cursor: pointer;
  margin-left: 6px;
  transition: all 0.2s ease;
  padding: 2px 4px;
  border-radius: 4px;
}

.msg-link:hover {
  color: #4f46e5;
  background-color: #eef2ff;
}

.modal-lg {
  max-width: 800px;
}

.modal-body {
  max-height: 80vh;
  overflow-y: auto;
}

.basic-info, .description, .raw-data {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 1rem;
}

pre {
  margin: 0;
  max-height: 400px;
  overflow-y: auto;
}

code {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style> 