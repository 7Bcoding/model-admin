<template>
  <div class="fusion-vendors">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索Vendor名称"
          style="width: 300px;"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      <div class="toolbar-right">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          创建Vendor
        </el-button>
      </div>
    </div>

    <!-- Vendor列表 -->
    <div class="vendors-table">
      <el-table
        :data="filteredVendors"
        v-loading="loading"
        style="width: 100%"
        :header-cell-style="{ background: '#f5f7fa' }"
      >
        <el-table-column label="名称" prop="name" min-width="150" />
        <el-table-column label="Vendor" prop="vendor" min-width="150">
          <template #default="{ row }">
            <span>{{ row.vendor || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="账户名" prop="accountName" min-width="150">
          <template #default="{ row }">
            <span>{{ row.accountName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述" prop="desc" min-width="200">
          <template #default="{ row }">
            <span>{{ row.desc || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="支付方式" prop="paymentMethod" width="120">
          <template #default="{ row }">
            <el-tag :type="getPaymentMethodType(row.paymentMethod)">
              {{ getPaymentMethodLabel(row.paymentMethod) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="余额" prop="balance" width="120">
          <template #default="{ row }">
            <span>{{ formatBalance(row.balance) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="关联API账号与模型" min-width="300">
          <template #header>
            <div class="accounts-header">
              <span>关联API账号与模型</span>
              <el-switch
                v-model="showAccountDetails"
                size="small"
                active-text="详情"
                inactive-text="概览"
                style="margin-left: 8px;"
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="related-accounts-container">
              <div v-if="accountsLoading || modelsLoading" class="loading-accounts">
                <span class="loading-text">加载中...</span>
              </div>
              <div v-else-if="getRelatedAccounts(row.name).length > 0" class="related-accounts-content">
                <!-- 概览模式 -->
                <div v-if="!showAccountDetails && !getRowDetailMode(row.name)" class="accounts-overview" @click="toggleRowDetailMode(row.name)">
                  <div class="overview-info">
                    <span class="overview-icon">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="#67c23a"/>
                      </svg>
                    </span>
                    <span class="overview-text">{{ getRelatedAccounts(row.name).length }} 个API账号</span>
                    <span class="models-overview-text">{{ getVendorRelatedModels(row.name).length }} 个模型</span>
                    <span class="expand-hint">点击展开</span>
                  </div>
                </div>
                <!-- 详情模式（全局开关或单行开关） -->
                <div v-else class="accounts-detail">
                  <div class="accounts-header-row">
                    <span class="accounts-title">关联详情</span>
                    <el-button 
                      v-if="!showAccountDetails"
                      type="text" 
                      size="small" 
                      @click="toggleRowDetailMode(row.name)"
                      class="collapse-btn"
                    >
                      <el-icon><ArrowUp /></el-icon>
                    </el-button>
                  </div>
                  <div class="accounts-list">
                    <div v-for="(account, index) in getRelatedAccounts(row.name)" :key="index" class="account-item">
                      <div class="account-header">
                        <span class="account-name">{{ account.name }}</span>
                        <el-tag size="small" :type="getApiTypeTagType(account.api_type || account.apiType)">
                          {{ getApiTypeLabel(account.api_type || account.apiType) }}
                        </el-tag>
                      </div>
                      <div class="account-models">
                        <div v-if="getAccountRelatedModels(account.name).length > 0" class="models-list">
                          <div v-for="(modelInfo, modelIndex) in getAccountRelatedModels(account.name)" :key="modelIndex" class="model-item">
                            <span class="model-name">{{ modelInfo.modelName }}</span>
                            <el-tag size="small" type="info" class="provider-tag">{{ modelInfo.providerName }}</el-tag>
                            <el-tag size="small" :type="modelInfo.type === 'primary' ? 'success' : 'warning'" class="type-tag">
                              {{ modelInfo.type === 'primary' ? '主要' : '备用' }}
                            </el-tag>
                          </div>
                        </div>
                        <div v-else class="no-models">
                          <span class="no-models-text">未关联模型</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="no-related-accounts">
                <span class="no-accounts-text">未关联</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" prop="updated_at" width="160">
          <template #default="{ row }">
            <span>{{ formatTimestamp(row.updated_at || row.updatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="editVendor(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="deleteVendor(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>



    <!-- 创建/编辑Vendor对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑Vendor' : '创建Vendor'"
      width="600px"
      class="vendor-dialog"
    >
      <el-form
        ref="vendorFormRef"
        :model="vendorForm"
        :rules="vendorRules"
        label-width="100px"
        class="vendor-form"
      >
        <el-form-item label="名称" prop="name">
          <el-input 
            v-model="vendorForm.name" 
            :disabled="isEdit"
            placeholder="请输入名称" 
          />
        </el-form-item>
        
        <el-form-item label="Vendor" prop="vendor">
          <el-input 
            v-model="vendorForm.vendor" 
            placeholder="请输入Vendor" 
          />
        </el-form-item>
        
        <el-form-item label="账户名" prop="accountName">
          <el-input 
            v-model="vendorForm.accountName" 
            placeholder="请输入账户名" 
          />
        </el-form-item>
        
        <el-form-item label="描述" prop="desc">
          <el-input
            v-model="vendorForm.desc"
            type="textarea"
            :rows="3"
            placeholder="请输入Vendor描述"
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="vendorForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="支付方式" prop="payment_method">
          <el-select v-model="vendorForm.payment_method" placeholder="请选择支付方式">
            <el-option 
              v-for="(label, value) in PaymentMethodLabels" 
              :key="value" 
              :label="label" 
              :value="parseInt(value)" 
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="余额" prop="balance">
          <el-input-number
            v-model="vendorForm.balance"
            :precision="2"
            :step="0.01"
            :min="0"
            placeholder="请输入余额"
            style="width: 100%"
          />
        </el-form-item>
        
        <!-- 时间信息（仅在编辑时显示） -->
        <div v-if="isEdit" class="time-info-section">
          <el-form-item label="创建时间">
            <span class="time-display">{{ formatTimestamp(vendorForm.createdAt) }}</span>
          </el-form-item>
          <el-form-item label="更新时间">
            <span class="time-display">{{ formatTimestamp(vendorForm.updatedAt) }}</span>
          </el-form-item>
        </div>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveVendor" :loading="saving">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, ArrowUp } from '@element-plus/icons-vue'
import { 
  getVendors, 
  createVendor, 
  updateVendor, 
  deleteVendor as deleteVendorApi 
} from '@/api/fusion-vendors'
import { getAccounts } from '@/api/fusion-accounts'
import { getModels } from '@/api/fusion-models'
import { PaymentMethod, PaymentMethodLabels, PaymentMethodTypes } from '@/config/specialFeatures'

// Reactive data
const loading = ref(false)
const saving = ref(false)
const vendors = ref([])
const accounts = ref([])
const accountsLoading = ref(false)
const models = ref([])
const modelsLoading = ref(false)
const searchQuery = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const vendorFormRef = ref()

// 控制关联账号详情显示
const showAccountDetails = ref(false)
const rowDetailModes = ref(new Set()) // 控制单行详情显示

// Vendor form
const vendorForm = ref({
  name: '',
  vendor: '',
  accountName: '',
  desc: '',
  status: 1,
  payment_method: 0,
  balance: null
})

// Form validation rules
const vendorRules = {
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在2到50个字符', trigger: 'blur' }
  ],
  vendor: [
    { required: false, message: '请输入Vendor', trigger: 'blur' },
    { max: 100, message: 'Vendor长度不能超过100个字符', trigger: 'blur' }
  ],
  accountName: [
    { required: false, message: '请输入账户名', trigger: 'blur' },
    { max: 100, message: '账户名长度不能超过100个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  payment_method: [
    { required: true, message: '请选择支付方式', trigger: 'change' }
  ]
}

// Computed properties
const filteredVendors = computed(() => {
  if (!searchQuery.value) {
    return vendors.value
  }
  const query = searchQuery.value.toLowerCase()
  return vendors.value.filter(vendor => 
    vendor.name.toLowerCase().includes(query) ||
    (vendor.desc && vendor.desc.toLowerCase().includes(query))
  )
})

// 获取供应商关联的API账号
const getRelatedAccounts = (vendorName) => {
  if (!vendorName || !accounts.value.length) {
    return []
  }
  
  // 根据vendor字段匹配供应商name
  return accounts.value.filter(account => 
    account.vendor === vendorName
  )
}

// 获取API账号关联的模型信息
const getAccountRelatedModels = (accountName) => {
  if (!accountName || !models.value.length) {
    return []
  }
  
  // 查找使用该账号的模型
  const relatedModels = []
  models.value.forEach(model => {
    // 检查主要提供商
    if (model.providers && model.providers.length > 0) {
      model.providers.forEach(provider => {
        if (provider.accounts && provider.accounts.some(acc => acc.name === accountName)) {
          relatedModels.push({
            modelName: model.name,
            providerName: provider.name,
            type: 'primary'
          })
        }
      })
    }
    
    // 检查备用提供商
    const fallbackProviders = model.fallback_providers || model.fallbackProviders || []
    if (fallbackProviders.length > 0) {
      fallbackProviders.forEach(provider => {
        if (provider.accounts && provider.accounts.some(acc => acc.name === accountName)) {
          relatedModels.push({
            modelName: model.name,
            providerName: provider.name,
            type: 'fallback'
          })
        }
      })
    }
  })
  
  return relatedModels
}

// 获取供应商关联的所有模型信息（通过API账号）
const getVendorRelatedModels = (vendorName) => {
  const relatedAccounts = getRelatedAccounts(vendorName)
  const allModels = []
  
  relatedAccounts.forEach(account => {
    const accountModels = getAccountRelatedModels(account.name)
    accountModels.forEach(modelInfo => {
      // 添加账号信息到模型信息中
      allModels.push({
        ...modelInfo,
        accountName: account.name,
        accountType: account.api_type || account.apiType
      })
    })
  })
  
  // 去重（同一个模型可能通过多个账号关联）
  const uniqueModels = []
  const seen = new Set()
  
  allModels.forEach(modelInfo => {
    const key = `${modelInfo.modelName}-${modelInfo.providerName}-${modelInfo.accountName}`
    if (!seen.has(key)) {
      seen.add(key)
      uniqueModels.push(modelInfo)
    }
  })
  
  return uniqueModels
}

// 切换单行详情显示模式
const toggleRowDetailMode = (vendorName) => {
  if (rowDetailModes.value.has(vendorName)) {
    rowDetailModes.value.delete(vendorName)
  } else {
    rowDetailModes.value.add(vendorName)
  }
}

// 检查单行是否显示详情
const getRowDetailMode = (vendorName) => {
  return rowDetailModes.value.has(vendorName)
}

// Methods
const loadVendors = async () => {
  loading.value = true
  try {
    const response = await getVendors()
    vendors.value = response.data.vendors || []
  } catch (error) {
    console.error('Failed to load vendors:', error)
    ElMessage.error('加载Vendor列表失败')
  } finally {
    loading.value = false
  }
}

const loadAccounts = async () => {
  accountsLoading.value = true
  try {
    const response = await getAccounts()
    accounts.value = response.data.accounts || []
    console.log('Loaded accounts for vendor relations:', accounts.value.length)
  } catch (error) {
    console.error('Failed to load accounts:', error)
    ElMessage.error('加载API账号失败')
  } finally {
    accountsLoading.value = false
  }
}

const loadModels = async () => {
  modelsLoading.value = true
  try {
    const allModels = []
    let pageToken = ''
    let hasMore = true
    
    // 分页加载所有模型
    while (hasMore) {
      const params = {
        pageSize: 20,
        pageToken: pageToken
      }
      
      const response = await getModels(params)
      const modelList = response.data.models || []
      allModels.push(...modelList)
      
      // 检查是否还有更多页面
      const nextToken = response.data.nextPageToken || response.data.next_page_token || ''
      if (nextToken && modelList.length > 0) {
        pageToken = nextToken
      } else {
        hasMore = false
      }
    }
    
    models.value = allModels
    console.log('Loaded models for vendor relations:', models.value.length)
  } catch (error) {
    console.error('Failed to load models:', error)
    ElMessage.error('加载模型数据失败')
  } finally {
    modelsLoading.value = false
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  resetVendorForm()
  dialogVisible.value = true
}

const editVendor = (vendor) => {
  isEdit.value = true
  vendorForm.value = {
    id: vendor.id,
    name: vendor.name,
    vendor: vendor.vendor || '',
    accountName: vendor.accountName || '',
    desc: vendor.desc || '',
    status: vendor.status,
    payment_method: vendor.paymentMethod,
    balance: vendor.balance,
    createdAt: vendor.createdAt || vendor.created_at,
    updatedAt: vendor.updatedAt || vendor.updated_at
  }
  dialogVisible.value = true
}

const saveVendor = async () => {
  try {
    await vendorFormRef.value.validate()
    saving.value = true
    
    if (isEdit.value) {
      // 更新时传递Vendor信息（不包含id）
      const vendorData = {
        name: vendorForm.value.name,
        vendor: vendorForm.value.vendor || '',
        accountName: vendorForm.value.accountName || '',
        desc: vendorForm.value.desc || '',
        status: vendorForm.value.status,
        payment_method: vendorForm.value.payment_method,
        balance: vendorForm.value.balance
      }
      await updateVendor(vendorForm.value.name, vendorData)
      ElMessage.success('Vendor更新成功')
    } else {
      // 创建时只需要必要字段
      const vendorData = {
        name: vendorForm.value.name,
        vendor: vendorForm.value.vendor || '',
        accountName: vendorForm.value.accountName || '',
        desc: vendorForm.value.desc || '',
        status: vendorForm.value.status,
        payment_method: vendorForm.value.payment_method,
        balance: vendorForm.value.balance
      }
      await createVendor(vendorData)
      ElMessage.success('Vendor创建成功')
    }
    
    dialogVisible.value = false
    loadVendors()
  } catch (error) {
    console.error('Failed to save vendor:', error)
    ElMessage.error(isEdit.value ? 'Vendor更新失败' : 'Vendor创建失败')
  } finally {
    saving.value = false
  }
}

const deleteVendor = async (vendor) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除Vendor "${vendor.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteVendorApi(vendor.name)
    ElMessage.success('Vendor删除成功')
    loadVendors()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete vendor:', error)
      
      // 提取详细错误信息
      let errorMessage = 'Vendor删除失败'
      if (error.response?.data) {
        const errorData = error.response.data
        if (errorData.message) {
          errorMessage = errorData.message
        } else if (errorData.reason) {
          errorMessage = errorData.reason
        }
      } else if (error.message) {
        errorMessage = error.message
      }
      
      ElMessage.error(errorMessage)
    }
  }
}

const resetVendorForm = () => {
  vendorForm.value = {
    name: '',
    vendor: '',
    accountName: '',
    desc: '',
    status: 1,
    payment_method: 0,
    balance: null,
    createdAt: null,
    updatedAt: null
  }
}

const handleSearch = () => {
  // 搜索功能保持不变，但不需要重置页码
}

// Utility methods
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'
  
  let date
  if (typeof timestamp === 'string') {
    date = new Date(timestamp)
  } else if (timestamp.seconds) {
    // 处理protobuf Timestamp格式
    date = new Date(timestamp.seconds * 1000)
  } else {
    date = new Date(timestamp)
  }
  
  return date.toLocaleString('zh-CN')
}

// 获取支付方式标签
const getPaymentMethodLabel = (method) => {
  // 确保method是数字类型
  const methodValue = parseInt(method)
  return PaymentMethodLabels[methodValue] || `未知(${method})`
}

// 获取支付方式标签类型
const getPaymentMethodType = (method) => {
  // 确保method是数字类型
  const methodValue = parseInt(method)
  return PaymentMethodTypes[methodValue] || 'info'
}

const formatBalance = (balance) => {
  if (balance === null || balance === undefined) return '-'
  return `¥${balance.toFixed(2)}`
}

// API类型标签处理
const getApiTypeLabel = (apiType) => {
  const typeMap = {
    'API_OPENAI': 'OpenAI',
    'API_ANTHROPIC': 'Anthropic',
    'API_VOLC': 'Volc',
    'API_QWEN': 'Qwen',
    'API_OPENROUTER': 'OpenRouter',
    'API_FIREWORKS': 'Fireworks',
    'API_BAIDU': 'Baidu',
    'API_TRT': 'TRT',
    'API_AWS': 'AWS',
    'API_AZURE': 'Azure',
    'API_VERTEX': 'Vertex',
    'API_VERTEX_GEMINI': 'Vertex Gemini'
  }
  return typeMap[apiType] || apiType || '未知'
}

const getApiTypeTagType = (apiType) => {
  const typeMap = {
    'API_OPENAI': 'primary',
    'API_ANTHROPIC': 'success',
    'API_VOLC': 'warning',
    'API_QWEN': 'info',
    'API_OPENROUTER': 'primary',
    'API_FIREWORKS': 'danger',
    'API_BAIDU': 'warning',
    'API_TRT': 'info',
    'API_AWS': 'success',
    'API_AZURE': 'primary',
    'API_VERTEX': 'success',
    'API_VERTEX_GEMINI': 'success'
  }
  return typeMap[apiType] || 'info'
}



onMounted(() => {
  loadVendors()
  loadAccounts()
  loadModels()
})
</script>

<style scoped>
.fusion-vendors {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.vendors-table {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}



.vendor-dialog .el-dialog__body {
  padding: 20px;
}

.vendor-form .el-form-item {
  margin-bottom: 20px;
}

.time-info-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

.time-display {
  color: #606266;
  font-size: 14px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

/* 关联API账号样式 */
.related-accounts-container {
  min-height: 24px;
}

.loading-accounts {
  display: flex;
  align-items: center;
  color: #909399;
  font-size: 12px;
}

.loading-text {
  margin-left: 4px;
}

.related-accounts-content {
  width: 100%;
}

/* 表头控制开关 */
.accounts-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

/* 概览模式 */
.accounts-overview {
  width: 100%;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.accounts-overview:hover {
  background-color: #f5f7fa;
}

.overview-info {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.overview-icon {
  display: flex;
  align-items: center;
}

.overview-text {
  color: #409eff;
  font-weight: 500;
  font-size: 12px;
}

.models-overview-text {
  color: #67c23a;
  font-weight: 500;
  font-size: 12px;
}

.expand-hint {
  color: #909399;
  font-size: 11px;
  font-style: italic;
}

/* 详情模式 */
.accounts-detail {
  width: 100%;
}

.accounts-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.accounts-title {
  font-size: 12px;
  color: #67c23a;
  font-weight: 500;
}

.collapse-btn {
  padding: 0 !important;
  min-height: auto !important;
  height: 16px !important;
  color: #909399 !important;
}

.collapse-btn:hover {
  color: #409eff !important;
}

.accounts-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-item {
  padding: 8px;
  background-color: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.account-name {
  color: #303133;
  font-weight: 500;
  flex: 1;
  margin-right: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-related-accounts {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
  font-size: 12px;
  height: 24px;
}

.no-accounts-text {
  font-style: italic;
}

.account-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.account-models {
  margin-top: 6px;
}

.models-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.model-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 6px;
  background-color: #ffffff;
  border-radius: 4px;
  border: 1px solid #e1e6f0;
  flex-wrap: wrap;
}

.model-name {
  font-weight: 500;
  color: #2c3e50;
  font-size: 13px;
}

.provider-tag {
  font-size: 11px;
}

.type-tag {
  font-size: 11px;
}

.no-models {
  padding: 4px 6px;
  color: #909399;
  font-style: italic;
  font-size: 12px;
}

.no-models-text {
  color: #c0c4cc;
}
</style>
