<template>
  <div class="page-container">
    <div class="header">
      <div class="header-left">
        <div class="title-section">
          <i class="el-icon-cpu title-icon"></i>
          <h2 class="page-title">模型管理</h2>
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
                    <!-- 如果是私有模型且有用户邮箱，显示邮箱 -->
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

  <!-- 修改部署信息弹窗内容 -->
  <el-dialog
    v-model="dialogVisible"
    title="部署信息"
    width="80%"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    destroy-on-close
  >
    <div v-loading="loading" element-loading-text="加载中...">
      <div v-for="(deployment, index) in deploymentDetails" :key="index" class="deployment-item">
        <div class="endpoint-section">
          <div class="endpoint-header">
            <div class="header-item">
              <span class="label">Endpoint ID:</span>
              <span class="value endpoint-id">{{ deployment.endpoint_id }}</span>
            </div>
            <div class="header-item">
              <span class="label">URL:</span>
              <span class="value url-value">
                {{ deployment.url }}
                <template v-if="getEndpointTypeTag(deployment.url)">
                  <el-tag :style="{marginLeft: '6px', background: getEndpointTypeTag(deployment.url).color, color: '#fff', border: 'none'}" size="small" :title="getEndpointTypeTag(deployment.url).tip">
                    {{ getEndpointTypeTag(deployment.url).tag }}
                  </el-tag>
                </template>
              </span>
            </div>
          </div>
          
          <div v-if="deployment.serverlessInfo" class="endpoint-content">
            <div class="info-item">
              <span class="label">Model Path:</span>
              <span class="value model-path">{{ deployment.serverlessInfo.modelPath }}</span>
              <el-button 
                type="text" 
                @click="openHuggingFace(deployment.serverlessInfo.modelPath)"
                class="view-button"
              >
                查看
              </el-button>
            </div>
            <div class="info-item">
              <span class="label">Image:</span>
              <span class="value image-value">{{ deployment.serverlessInfo.image }}</span>
            </div>
            <div class="info-item">
              <span class="label">VARS:</span>
              <div class="vars-display">
                <div v-for="envVar in deployment.serverlessInfo.vars"   
                     :key="envVar.name" 
                     class="var-item"
                >
                  <span class="var-name">{{ envVar.name }}</span>
                  <span class="var-value">{{ envVar.value }}</span>
                </div>
              </div>
            </div>
            <div class="info-item">
              <span class="label">Args:</span>
              <div class="args-display">
                <pre>{{ deployment.serverlessInfo.args.join('\n') }}</pre>
                <el-button type="text" class="copy-button" @click="copyText(deployment.serverlessInfo.args.join('\n'))">复制</el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>

  <!-- 添加备注对话框 -->
  <el-dialog
    v-model="noteDialogVisible"
    :title="`${currentModel?.note ? '编辑' : '添加'}备注`"
    width="30%"
  >
    <el-form>
      <el-form-item label="备注">
        <el-input
          v-model="noteInput"
          type="textarea"
          :rows="3"
          placeholder="请输入备注内容"
        />
      </el-form-item>
      <el-form-item label="群组ID">
        <el-input
          v-model="openChatIdInput"
          placeholder="请输入飞书群组ID"
        >
          <template #append>
            <el-tooltip
              content="飞书群组的唯一标识，用于模型相关通知"
              placement="top"
            >
              <i class="el-icon-question"></i>
            </el-tooltip>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="推理引擎">
        <el-input
          v-model="inferenceEngineInput"
          placeholder="请输入推理引擎，默认为vllm"
        >
          <template #append>
            <el-tooltip
              content="模型的推理引擎类型"
              placement="top"
            >
              <i class="el-icon-question"></i>
            </el-tooltip>
          </template>
        </el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="noteDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveNote">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { getModels, starModel, unstarModel, updateModelNote } from '../api/models'
import { useToast } from '../composables/toast'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'
import { Edit, ChatSquare } from '@element-plus/icons-vue'
import { useModelsStore } from '@/stores/models'

export default {
  name: 'Models',
  setup() {
    const models = ref([])
    const search = ref('')
    const showStarred = ref(false)
    const toast = useToast()
    const selectedModel = ref(null)
    let modelDetailsModal = null
    const showRawData = ref(false)
    const router = useRouter()
    const dialogVisible = ref(false)
    const deploymentDetails = ref([])
    const loading = ref(false)
    const noteDialogVisible = ref(false)
    const noteInput = ref('')
    const currentModel = ref(null)
    const openChatIdInput = ref('')
    const inferenceEngineInput = ref('')
    const isAdmin = computed(() => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      return user.role === 'admin'
    })
    const showActive = ref(true)
    const modelsStore = useModelsStore()

    const filteredModels = computed(() => {
      if (!Array.isArray(models.value)) {
        console.warn('models.value is not an array:', models.value)
        return []
      }

      const filtered = models.value.filter(model => {
        // 基本验证
        if (!model || typeof model !== 'object') {
          console.warn('Invalid model object:', model)
          return false
        }

        // 收藏过滤
        if (showStarred.value && !model.starred) {
          return false
        }

        // 搜索过滤
        if (search.value) {
          const searchTerm = search.value.toLowerCase()
          const modelName = String(model.model_name || '').toLowerCase()
          const note = String(model.note || '').toLowerCase()
          const searchTerms = searchTerm.split(/\s+/).filter(term => term)
          
          if (!searchTerms.every(term => 
            modelName.includes(term) || note.includes(term)
          )) {
            return false
          }
        }

        // 状态过滤 - 修复这里
        if (showActive.value) {
          // 只显示状态为MODEL_STATUS_SERVING的模型
          if (model.status !== 'MODEL_STATUS_SERVING') {
            return false
          }
        }

        return true
      })

      // 对过滤后的结果进行排序，收藏的排在前面
      return filtered.sort((a, b) => {
        if (a.starred !== b.starred) {
          return a.starred ? -1 : 1
        }
        return a.model_name.localeCompare(b.model_name)
      })
    })

    const getStatusClass = (status) => {
      switch (status) {
        case 'MODEL_STATUS_SERVING':
          return 'bg-success'
        case 'MODEL_STATUS_PENDING':
          return 'bg-warning'
        case 'MODEL_STATUS_DELETED':
          return 'bg-danger'
        default:
          return 'bg-secondary'
      }
    }

    const getStatusText = (status) => {
      switch (status) {
        case 'MODEL_STATUS_SERVING':
          return '运行中'
        case 'MODEL_STATUS_PENDING':
          return '等待中'
        case 'MODEL_STATUS_DELETED':
          return '已删除'
        default:
          return '未知'
      }
    }

    const showDetails = (model) => {
      selectedModel.value = model
      if (!modelDetailsModal) {
        const modalEl = document.getElementById('modelDetailsModal')
        if (window.bootstrap) {
          modelDetailsModal = new window.bootstrap.Modal(modalEl)
          modalEl.addEventListener('hidden.bs.modal', () => {
            selectedModel.value = null
          })
        } else {
          console.error('Bootstrap is not loaded')
          return
        }
      }
      modelDetailsModal.show()
    }

    const fetchData = async () => {
      loading.value = true
      try {
        const fetchedModels = await modelsStore.fetchModels()
        models.value = fetchedModels
      } catch (error) {
        console.error('Error fetching models:', error)
        toast.show('获取模型列表失败', 'error')
      } finally {
        loading.value = false
      }
    }

    const refreshData = async () => {
      loading.value = true
      try {
        const fetchedModels = await modelsStore.fetchModels(true) // 强制刷新
        models.value = fetchedModels
        toast.show('刷新成功', 'success')
      } catch (error) {
        console.error('Error refreshing models:', error)
        toast.show('刷新失败', 'error')
      } finally {
        loading.value = false
      }
    }

    const toggleStar = async (model) => {
      try {
        if (model.starred) {
          await unstarModel(model.model_name)
          toast.show('已取消收藏', 'success')
        } else {
          await starModel(model.model_name)
          toast.show('已添加收藏', 'success')
        }
        model.starred = !model.starred
      } catch (error) {
        console.error('Failed to toggle star:', error)
        toast.show('操作失败', 'error')
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

    // 监控相关方法
    const openMonitor = (model, type) => {
      if (!model || !model.model_name) {
        ElMessage.warning('模型信息不完整')
        return
      }

      const monitorUrls = {
        overview: `https://grafana.aicloud.pplabs.tech/d/430e32c1-c956-4f4d-a4e0-8c463180cbe3/ai-cloud-platform-sla-metrics-new?orgId=1&from=now-12h&to=now&orgId=1&var-CustomerUUID=ALL&var-Model=${encodeURIComponent(model.model_name)}&refresh=5s&var-Customer=All`,
        vllm: model.inference_engine === 'sglang' ? 
          `https://grafana.aicloud.pplabs.tech/d/ddyp55uq7brpcc/sglang-dashboard?orgId=1&refresh=5s&var-model_id=${encodeURIComponent(getSglangModelId(model.model_name))}` :
          `https://grafana.aicloud.pplabs.tech/d/b281712d-8bff-41ef-9f3f-71ad43c05e9bdd/vllm-per-model?orgId=1&var-model_name=${encodeURIComponent(model.model_name)}`,
        serverless: model.endpoints ? 
          `https://grafana.aicloud.pplabs.tech/d/KOpAVRCIz1/prod-overview?orgId=1${
            model.endpoints
              .filter(endpoint => endpoint.url.includes('.runsync.alpha.dev'))
              .map(endpoint => `&var-endpoint=${endpoint.url.replace('https://', '').replace('http://', '').split('.runsync.alpha.dev')[0]}`)
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

    // 模型管理相关方法
    const manageModel = (model) => {
      window.location.href = `/models/${encodeURIComponent(model.model_name)}/manage`
    }

    const showEndpoints = (model) => {
      // 使用 router 导航到 endpoints 管理页面
      router.push({
        name: 'EndpointManager',
        params: {
          modelName: model.model_name
        }
      })
    }

    // 从 URL 中提取 serverless endpoint 名称
    const extractServerlessEndpoint = (url) => {
      if (!url.includes('runsync.alpha.dev')) return null
      const match = url.match(/https:\/\/(.*?)\.runsync\.alpha\.dev/)
      return match ? match[1] : null
    }

    // 获取 serverless endpoint 
    const fetchServerlessEndpoint = async (name) => {
      try {
        const response = await axios.get(`se/${name}`)
        return response.data.data
      } catch (error) {
        console.error('Error fetching serverless endpoint:', error)
        ElMessage.error(`获取 endpoint 详情失败: ${error.message}`)
        return null
      }
    }

    // 修改部署信息显示方法
    const showDeploymentInfo = async (model) => {
      dialogVisible.value = true
      deploymentDetails.value = []
      loading.value = true

      try {
        const deployments = model.endpoints || []
        for (const endpoint of deployments) {
          const detail = {
            endpoint_id: endpoint.endpoint_id,
            url: endpoint.url,
            collapsed: false  // 添加折叠状态
          }

          const seName = extractServerlessEndpoint(endpoint.url)
          if (seName) {
            const seInfo = await fetchServerlessEndpoint(seName)
            if (seInfo) {
              const container = seInfo.spec.template.spec.containers[0]
              // 格式化参数
              const formattedArgs = []
              const args = container.args
              for (let i = 0; i < args.length; i++) {
                const arg = args[i]
                if (arg.startsWith('--')) {
                  // 如果下一个参数不是以--开头，则认为是当前参数的值
                  const value = i + 1 < args.length && !args[i + 1].startsWith('--') ? args[++i] : ''
                  formattedArgs.push({
                    key: arg,
                    value: value
                  })
                } else {
                  formattedArgs.push(arg)
                }
              }

              // 提取 Model ID
              const modelArg = formattedArgs.find(arg => arg.key === '--model')
              detail.serverlessInfo = {
                modelPath: modelArg ? modelArg.value : '',
                image: container.image,
                args: container.args,
                formattedArgs: formattedArgs,
                vars: container.env
              }
            }
          }

          deploymentDetails.value.push(detail)
        }
      } catch (error) {
        console.error('Error processing deployment info:', error)
        ElMessage.error('处理部署信息发生错误')
      } finally {
        loading.value = false
      }
    }

    // 复制文本功能
    const copyText = async (text) => {
      try {
        await navigator.clipboard.writeText(text)
        ElMessage.success('复制成功')
      } catch (error) {
        console.error('Copy failed:', error)
        ElMessage.error('复制失败')
      }
    }

    const toggleRawData = () => {
      showRawData.value = !showRawData.value
    }

    // 添加提取 endpoint 名称的方法
    const extractEndpointName = (url) => {
      const match = url.match(/https:\/\/(.*?)\.runsync\.alpha\.dev/)
      return match ? match[1] : url
    }

    // 添加打开 Hugging Face 的方法
    const openHuggingFace = (modelPath) => {
      const url = `https://huggingface.co/${modelPath}`
      window.open(url, '_blank')
    }

    const toggleEndpoint = (index) => {
      if (!deploymentDetails.value[index].collapsed) {
        deploymentDetails.value[index].collapsed = true
      } else {
        deploymentDetails.value[index].collapsed = false
      }
    }

    // 显示备注对话框
    const showNoteDialog = (model) => {
      currentModel.value = model
      noteInput.value = model.note || ''
      openChatIdInput.value = model.open_chat_id || ''
      inferenceEngineInput.value = model.inference_engine || 'vllm'
      noteDialogVisible.value = true
    }

    // 保存备注
    const saveNote = async () => {
      try {
        await updateModelNote(currentModel.value.model_name, noteInput.value, openChatIdInput.value, inferenceEngineInput.value)

        // 更新本地数据
        currentModel.value.note = noteInput.value
        currentModel.value.open_chat_id = openChatIdInput.value
        currentModel.value.inference_engine = inferenceEngineInput.value
        
        ElMessage({
          message: '保存成功',
          type: 'success'
        })
        
        noteDialogVisible.value = false
      } catch (error) {
        console.error('Error saving note:', error)
        ElMessage({
          message: error.response?.data?.message || '保存失败',
          type: 'error'
        })
      }
    }

    const openFeishuChat = (openChatId) => {
      const url = `https://applink.feishu.cn/client/chat/open?openChatId=${openChatId}`
      // open feishu chat in new window and auto close after 3 seconds
      const newWindow = window.open(url, '_blank')
      setTimeout(() => {
        if (newWindow) {
          newWindow.close()
        }
      }, 2000)
      return
    }

    // 新增 endpoint 类型解析方法
    const getEndpointTypeTag = (url) => {
      if (!url) return ''
      // serverless2.0
      if (/^https:\/\/[\w-]+\.us-01\.sls2\.alpha\.ai/.test(url)) {
        return { tag: '2.0', color: '#2563eb', tip: 'Serverless 2.0' }
      }
      // 融合服
      if (/^http:\/\/thirdparty-api-adapter\/fusion\/v1\//.test(url)) {
        return { tag: '融', color: '#f59e42', tip: '融合服' }
      }
      return null
    }

    // 页面挂载时从 localStorage 恢复搜索条件
    onMounted(async () => {
      console.log('Models component mounted')
      
      // 恢复搜索框内容
      const savedSearch = localStorage.getItem('models_search_query')
      if (savedSearch !== null) {
        search.value = savedSearch
      }
      
      // 恢复过滤选项
      const savedShowStarred = localStorage.getItem('models_show_starred')
      if (savedShowStarred !== null) {
        showStarred.value = savedShowStarred === 'true'
      }
      
      const savedShowActive = localStorage.getItem('models_show_active')
      if (savedShowActive !== null) {
        showActive.value = savedShowActive === 'true'
      }
      
      // 加载数据
      await fetchData()
    })

    // 监听搜索条件变化，保存到 localStorage
    watch(search, (val) => {
      localStorage.setItem('models_search_query', val)
    })

    watch(showStarred, (val) => {
      localStorage.setItem('models_show_starred', val)
    })

    watch(showActive, (val) => {
      localStorage.setItem('models_show_active', val)
    })

    onUnmounted(() => {
      // modelsStore.clearCache() // 取消注释如果需要在离开页面时清除缓存
    })

    return {
      models,
      search,
      showStarred,
      filteredModels,
      toggleStar,
      getStatusClass,
      getStatusText,
      showDetails,
      openMonitor,
      manageModel,
      showEndpoints,
      showDeploymentInfo,
      selectedModel,
      showRawData,
      toggleRawData,
      dialogVisible,
      deploymentDetails,
      loading,
      copyText,
      extractEndpointName,
      openHuggingFace,
      toggleEndpoint,
      noteDialogVisible,
      noteInput,
      currentModel,
      showNoteDialog,
      saveNote,
      isAdmin,
      showActive,
      openChatIdInput,
      inferenceEngineInput,
      openFeishuChat,
      ChatSquare,
      getEndpointTypeTag
    }
  }
}
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

.action-buttons {
  display: flex;
  gap: 8px;
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

h6 {
  margin-bottom: 1rem;
  color: #495057;
}

label.fw-bold {
  font-size: 0.9rem;
  color: #6c757d;
  margin-bottom: 0.25rem;
}

.deployment-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.deployment-item {
  margin-bottom: 16px;
}

.endpoint-section {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.endpoint-header {
  background: #f8fafc;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.header-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.header-item:last-child {
  margin-bottom: 0;
}

.endpoint-content {
  padding: 16px;
  background: #ffffff;
}

.info-item {
  margin-bottom: 12px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.label {
  font-weight: 500;
  color: #64748b;
  width: 100px;
  flex-shrink: 0;
}

.value {
  word-break: break-all;
  flex: 1;
}

.args-display {
  background: #1e293b;
  padding: 16px;
  border-radius: 6px;
  position: relative;
  width: 100%;
  margin-top: 8px;
}

.args-display pre {
  margin: 0;
  font-size: 13px;
  color: #e2e8f0;
  font-family: 'Fira Code', monospace;
}

.copy-button {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(255, 255, 255, 0.1);
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  color: #e2e8f0;
  
  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
}

.endpoint-id {
  color: #0f172a;
  font-family: monospace;
  font-weight: 600;
  background: #e2e8f0;
  padding: 2px 8px;
  border-radius: 4px;
}

.url-value {
  color: #2563eb;
  font-family: monospace;
}

.model-path {
  color: #0f766e;
  font-family: monospace;
}

.image-value {
  color: #7c3aed;
  font-family: monospace;
}

.view-button {
  color: #2563eb;
  font-weight: 500;
  
  &:hover {
    color: #1d4ed8;
    text-decoration: underline;
  }
}

.vars-display {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  width: 100%;
  margin-top: 4px;
}

.var-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 4px 0;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
}

.var-item:not(:last-child) {
  border-bottom: 1px dashed #e2e8f0;
}

.var-name {
  color: #0f766e;
  font-weight: 500;
  width: 280px;
  flex-shrink: 0;
}

.var-value {
  color: #7c3aed;
  word-break: break-all;
}

.model-metadata {
  display: flex;
  flex-direction: column;
  margin-left: 24px;
  gap: 4px;
}

.model-note {
  color: #6b7280;
  font-weight: 400;
  font-size: 0.85em;
  font-style: italic;
  opacity: 0.9;
}

.model-chat-id {
  color: #6b7280;
  font-size: 0.85em;
  display: flex;
  align-items: center;
  gap: 4px;
}

.model-chat-id i {
  font-size: 14px;
  color: #8b5cf6;
}

.chat-icon {
  font-size: 16px;
  color: #8b5cf6;
  cursor: pointer;
  margin-left: 6px;
  transition: all 0.2s ease;
}

.chat-icon:hover {
  color: #6d28d9;
  transform: scale(1.1);
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
</style>