<template>
  <div class="fusion-model-prices">
    <div class="toolbar">
      <div class="search-wrapper">
        <el-input
          v-model="searchQuery"
          placeholder="搜索模型名称、API账号或备注..."
          prefix-icon="Search"
          clearable
          style="width: 300px; margin-right: 10px"
          @input="handleSearch"
        />
        <el-select
          v-model="selectedApiAccount"
          placeholder="筛选API账号"
          clearable
          style="width: 200px; margin-right: 10px"
          @change="handleSearch"
        >
          <el-option
            v-for="account in apiAccounts"
            :key="account"
            :label="account"
            :value="account"
          />
        </el-select>
        <el-select
          v-model="selectedCurrency"
          placeholder="筛选货币类型"
          style="width: 150px"
          clearable
          @change="handleSearch"
        >
          <el-option label="USD" :value="0" />
          <el-option label="CNY" :value="1" />
        </el-select>
      </div>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          添加价格配置
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div class="table-container">
      <el-table
        :data="paginatedData"
        v-loading="loading"
        style="width: 100%"
        border
        :header-cell-style="{
          background: '#f5f7fa',
          color: '#333',
          fontWeight: 600,
          fontSize: '14px',
          padding: '12px 8px',
          textAlign: 'center'
        }"
        :cell-style="{
          color: '#333',
          fontSize: '14px',
          padding: '12px 8px'
        }"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="modelName" label="模型名称" min-width="140" />
        <el-table-column prop="apiAccount" label="API账号" min-width="120" />
        <el-table-column prop="promptPrice" label="prompt价格" min-width="110" align="right" sortable>
          <template #default="scope">
            {{ formatPrice(scope.row.promptPrice, scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="completionPrice" label="completions价格" min-width="110" align="right" sortable>
          <template #default="scope">
            {{ formatPrice(scope.row.completionPrice, scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="promptUncachePrice" label="未缓存prompt价格" min-width="130" align="right" sortable>
          <template #default="scope">
            {{ formatPrice(scope.row.promptUncachePrice, scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="promptCacheWritePrice" label="缓存创建价格" min-width="120" align="right" sortable>
          <template #default="scope">
            {{ formatPrice(scope.row.promptCacheWritePrice, scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="promptCacheReadPrice" label="缓存读取价格" min-width="120" align="right" sortable>
          <template #default="scope">
            {{ formatPrice(scope.row.promptCacheReadPrice, scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="reasoningContentPrice" label="reasoning内容价格" min-width="120" align="right" sortable>
          <template #default="scope">
            {{ formatPrice(scope.row.reasoningContentPrice, scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="currency" label="货币" width="80">
          <template #default="scope">
            {{ getCurrencyText(scope.row.currency) }}
          </template>
        </el-table-column>
        <el-table-column prop="desc" label="备注" min-width="150">
          <template #default="scope">
            <span :title="scope.row.desc" class="desc-text">
              {{ scope.row.desc || '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="140">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="showEditDialog(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- Pagination -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[15, 30, 50, 100]"
          :total="filteredData.length"
          :layout="filteredData.length > pageSize ? 'total, sizes, prev, pager, next, jumper' : 'total, sizes'"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑价格配置' : '添加价格配置'"
      width="650px"
      class="price-dialog"
    >
      <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          label-position="top"
          class="price-form"
        >
        <!-- 基本信息 -->
        <div class="form-section">
          <div class="section-title">
            <span class="section-icon">📋</span>
            基本信息
          </div>
          <div class="form-row">
            <el-form-item label="模型名称" prop="modelName" class="form-item-third">
              <!-- 编辑时显示只读输入框，创建时显示下拉选择 -->
              <el-input 
                v-if="isEdit"
                v-model="form.modelName" 
                readonly
                style="width: 100%"
                placeholder="模型名称"
              />
              <el-select 
                v-else
                v-model="form.modelName" 
                placeholder="请选择模型名称" 
                style="width: 100%"
                filterable
                :loading="modelsLoading"
              >
                <el-option
                  v-for="model in modelsList"
                  :key="model.name"
                  :label="model.name"
                  :value="model.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="API账号" prop="apiAccount" class="form-item-third">
              <!-- 编辑时显示只读输入框，创建时显示下拉选择 -->
              <el-input 
                v-if="isEdit"
                v-model="form.apiAccount" 
                readonly
                style="width: 100%"
                placeholder="API账号"
              />
              <el-select 
                v-else
                v-model="form.apiAccount" 
                placeholder="请选择API账号" 
                style="width: 100%"
                filterable
                :loading="accountsLoading"
              >
                <el-option
                  v-for="account in accountsList"
                  :key="account.name"
                  :label="account.name"
                  :value="account.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="货币" prop="currency" class="form-item-third">
              <el-select v-model="form.currency" placeholder="请选择货币" style="width: 100%">
                <el-option label="USD" :value="0" />
                <el-option label="CNY" :value="1" />
              </el-select>
            </el-form-item>
          </div>
        </div>

        <!-- 基础价格配置 -->
         <div class="form-section">
           <div class="section-title">
             <span class="section-icon">💰</span>
             基础价格配置
             <span class="price-unit">（每百万 Token）</span>
           </div>
           <div class="form-row">
             <el-form-item label="prompt Token 价格" prop="promptPrice" class="form-item-half">
               <el-input-number
                 v-model="form.promptPrice"
                 :precision="4"
                 :step="0.0001"
                 :min="0"
                 placeholder="prompt价格"
                 style="width: 100%"
               />
             </el-form-item>
             <el-form-item label="completions Token 价格" prop="completionPrice" class="form-item-half">
               <el-input-number
                 v-model="form.completionPrice"
                 :precision="4"
                 :step="0.0001"
                 :min="0"
                 placeholder="completions价格"
                 style="width: 100%"
               />
             </el-form-item>
           </div>
         </div>

         <!-- 缓存相关价格 -->
         <div class="form-section">
           <div class="section-title">
             <span class="section-icon">🗄️</span>
             缓存相关价格
             <span class="price-unit">（每百万 Token）</span>
           </div>
           <div class="form-row">
             <el-form-item label="未缓存prompt价格" prop="promptUncachePrice" class="form-item-half">
               <el-input-number
                 v-model="form.promptUncachePrice"
                 :precision="4"
                 :step="0.0001"
                 :min="0"
                 placeholder="未缓存prompt价格"
                 style="width: 100%"
               />
             </el-form-item>
             <el-form-item label="缓存创建价格" prop="promptCacheWritePrice" class="form-item-half">
               <el-input-number
                 v-model="form.promptCacheWritePrice"
                 :precision="4"
                 :step="0.0001"
                 :min="0"
                 placeholder="缓存创建价格"
                 style="width: 100%"
               />
             </el-form-item>
           </div>
           <div class="form-row">
             <el-form-item label="缓存读取价格" prop="promptCacheReadPrice" class="form-item-half">
               <el-input-number
                 v-model="form.promptCacheReadPrice"
                 :precision="4"
                 :step="0.0001"
                 :min="0"
                 placeholder="缓存读取价格"
                 style="width: 100%"
               />
             </el-form-item>
             <el-form-item label="reasoning内容价格" prop="reasoningContentPrice" class="form-item-half">
               <el-input-number
                 v-model="form.reasoningContentPrice"
                 :precision="4"
                 :step="0.0001"
                 :min="0"
                 placeholder="reasoning内容价格"
                 style="width: 100%"
               />
             </el-form-item>
           </div>
         </div>

        <!-- 其他信息 -->
        <div class="form-section">
          <div class="section-title">
            <span class="section-icon">📝</span>
            其他信息
          </div>
          <el-form-item label="描述" prop="desc">
            <el-input
              v-model="form.desc"
              type="textarea"
              :rows="3"
              placeholder="请输入描述信息"
            />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import {
  getModelPrices,
  createModelPrice,
  updateModelPrice,
  deleteModelPrice
} from '@/api/fusion-model-prices'
import { getModels } from '@/api/fusion-models'
import { getAccounts } from '@/api/fusion-accounts'
import { useFusionStore } from '@/stores/fusion'

// 响应式数据
const fusionStore = useFusionStore()
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const selectedApiAccount = ref('')
const selectedCurrency = ref(null)
const currentPage = ref(1)
const pageSize = ref(15)
const modelPrices = ref([])
const formRef = ref()
const modelsList = ref([])
const accountsList = ref([])
const modelsLoading = ref(false)
const accountsLoading = ref(false)

// 表单数据
const form = reactive({
  id: null,
  modelName: '',
  apiAccount: '',
  promptPrice: null,
  completionPrice: null,
  promptUncachePrice: null,
  promptCacheWritePrice: null,
  promptCacheReadPrice: null,
  reasoningContentPrice: null,
  currency: 0,
  desc: ''
})

// 表单验证规则
const rules = {
  modelName: [
    { required: true, message: '请输入模型名称', trigger: 'blur' }
  ],
  apiAccount: [
    { required: true, message: '请输入API账号', trigger: 'blur' }
  ],
  currency: [
    { required: true, message: '请选择货币', trigger: 'change' }
  ]
}

// 计算属性
const apiAccounts = computed(() => {
  const accountSet = new Set(modelPrices.value.map(item => item.apiAccount))
  return Array.from(accountSet).sort()
})

const filteredData = computed(() => {
  let data = modelPrices.value
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item => 
      item.modelName?.toLowerCase().includes(query) ||
      item.apiAccount?.toLowerCase().includes(query) ||
      item.desc?.toLowerCase().includes(query)
    )
  }
  
  if (selectedApiAccount.value) {
    data = data.filter(item => item.apiAccount === selectedApiAccount.value)
  }
  
  if (selectedCurrency.value !== null && selectedCurrency.value !== undefined) {
    data = data.filter(item => item.currency === selectedCurrency.value)
  }
  
  return data
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

// 方法
const fetchModelPrices = async () => {
  loading.value = true
  try {
    const response = await getModelPrices({
      page: 1,
      limit: 1000 // 获取所有数据，前端分页
    })
    console.log('API Response:', response.data)
    modelPrices.value = response.data?.modelPrices || []
    console.log('Model prices loaded:', modelPrices.value.length, 'items')
  } catch (error) {
    console.error('获取模型价格失败:', error)
    const errorMessage = error.response?.data?.message || error.response?.data?.reason || error.message || '获取模型价格失败'
    ElMessage.error(errorMessage)
  } finally {
    loading.value = false
  }
}

const fetchModels = async () => {
  try {
    modelsLoading.value = true
    const response = await getModels()
    modelsList.value = response.data?.models || []
    console.log('Loaded models:', modelsList.value.length, 'items')
  } catch (error) {
    console.error('获取模型列表失败:', error)
    ElMessage.error('获取模型列表失败')
    modelsList.value = []
  } finally {
    modelsLoading.value = false
  }
}

const fetchAccounts = async () => {
  try {
    accountsLoading.value = true
    const response = await getAccounts()
    accountsList.value = response.data?.accounts || []
    console.log('Loaded accounts:', accountsList.value.length, 'items')
  } catch (error) {
    console.error('获取账号列表失败:', error)
    ElMessage.error('获取账号列表失败')
    accountsList.value = []
  } finally {
    accountsLoading.value = false
  }
}

const refreshData = () => {
  fetchModelPrices()
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

const showCreateDialog = async () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
  // 加载模型和账号列表
  await Promise.all([fetchModels(), fetchAccounts()])
}

const showEditDialog = async (row) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, { ...row })
  // 加载模型和账号列表
  await Promise.all([fetchModels(), fetchAccounts()])
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(form, {
    id: null,
    modelName: '',
    apiAccount: '',
    promptPrice: null,
    completionPrice: null,
    promptUncachePrice: null,
    promptCacheWritePrice: null,
    promptCacheReadPrice: null,
    reasoningContentPrice: null,
    currency: 0,
    desc: ''
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    const data = { ...form }
    delete data.id
    
    if (isEdit.value) {
      await updateModelPrice(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await createModelPrice(data)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    await fetchModelPrices()
  } catch (error) {
    if (error.response) {
      console.error('提交失败:', error)
      const errorMessage = error.response?.data?.message || error.response?.data?.reason || error.message || '操作失败'
      ElMessage.error(errorMessage)
    }
  } finally {
    submitting.value = false
  }
}

const formatPrice = (price, currency = 0) => {
  if (price === null || price === undefined) return '-'
  const symbol = currency === 1 ? '¥' : '$'
  return `${symbol}${price.toFixed(4)}`
}

const getCurrencyText = (currency) => {
  switch (currency) {
    case 0:
      return 'USD'
    case 1:
      return 'CNY'
    default:
      return 'USD'
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除模型 "${row.modelName}" 的价格配置吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteModelPrice(row.id)
    ElMessage.success('删除成功')
    await fetchModelPrices()
  } catch (error) {
    if (error === 'cancel') return
    console.error('删除失败:', error)
    const errorMessage = error.response?.data?.message || error.response?.data?.reason || error.message || '删除失败'
    ElMessage.error(errorMessage)
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  console.log('FusionModelPrices component mounted')
  console.log('Current fusion store state:', fusionStore.currentBaseUrl)
  fetchModelPrices()
  fetchModels()
  fetchAccounts()
})
</script>

<style scoped>
.fusion-model-prices {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.search-wrapper {
  display: flex;
  align-items: center;
}

.actions {
  display: flex;
  gap: 10px;
}

.table-container {
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.pagination-container {
  padding: 20px;
  display: flex;
  justify-content: center;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
}

.price-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.price-form {
  max-height: 60vh;
  overflow-y: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 表格样式优化 */
.table-container :deep(.el-table) {
  border: none;
}

.table-container :deep(.el-table__header-wrapper) {
  border-radius: 8px 8px 0 0;
}

.table-container :deep(.el-table__body-wrapper) {
  border-radius: 0 0 8px 8px;
}

.table-container :deep(.el-table td),
.table-container :deep(.el-table th) {
  border-bottom: 1px solid #f0f0f0;
}

.table-container :deep(.el-table--border) {
  border: 1px solid #e9ecef;
}

.table-container :deep(.el-table--border::after) {
  display: none;
}

/* 表单分组样式 */
.form-section {
  margin-bottom: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  background: #fafbfc;
  overflow: hidden;
}

.section-title {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #f5f7fa;
  color: #606266;
  font-weight: 500;
  font-size: 12px;
  margin: 0;
  border-bottom: 1px solid #e4e7ed;
}

.section-icon {
  margin-right: 4px;
  font-size: 12px;
}

.price-unit {
  margin-left: auto;
  font-size: 10px;
  color: #909399;
  font-weight: normal;
}

/* 两列布局 */
.form-row {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
}

.form-item-half {
  flex: 1;
  margin: 0 !important;
}

.form-item-third {
  margin: 0 !important;
}

/* 三列布局中的特殊宽度分配 */
.form-item-third:nth-child(1) {
  flex: 2.3; /* 模型名称需要更多空间 */
}

.form-item-third:nth-child(2) {
  flex: 2.3; /* API账号需要更多空间 */
}

.form-item-third:nth-child(3) {
  flex: 1.1; /* 货币选择适中宽度 */
}

.form-section .el-form-item {
  margin: 8px 12px 0 12px;
  background: white;
  padding: 8px;
  border-radius: 3px;
  border: 1px solid #f0f0f0;
}

.form-section .el-form-item:last-child {
  margin-bottom: 8px;
}

.form-section .el-form-item:deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
  font-size: 12px;
  line-height: 1.2;
  margin-bottom: 6px;
  text-align: center;
  width: 100% !important;
  padding: 0 !important;
}

.form-section .el-form-item:deep(.el-input-number) {
  width: 100%;
}

.form-section .el-form-item:deep(.el-input-number .el-input__inner) {
  text-align: center;
  height: 32px;
  font-size: 13px;
}

.form-section .el-form-item:deep(.el-input__inner) {
  height: 32px;
  font-size: 13px;
  text-align: center;
}

.form-section .el-form-item:deep(.el-select) {
  width: 100%;
}

.form-section .el-form-item:deep(.el-select .el-input__inner) {
  text-align: center;
}

.form-section .el-form-item:deep(.el-textarea__inner) {
  font-size: 13px;
  text-align: center;
}

/* 只读输入框样式 */
.form-section .el-form-item:deep(.el-input.is-disabled .el-input__inner),
.form-section .el-form-item:deep(.el-input[readonly] .el-input__inner) {
  background-color: #f5f7fa !important;
  border-color: #e4e7ed !important;
  color: #606266 !important;
  cursor: not-allowed;
  text-align: center;
}

/* 不同section的柔和主题色 */
.form-section:nth-child(1) .section-title {
  background: #f0f2f5;
  color: #595959;
}

.form-section:nth-child(2) .section-title {
  background: #f6f8fc;
  color: #595959;
}

.form-section:nth-child(3) .section-title {
  background: #f0f9ff;
  color: #595959;
}

.form-section:nth-child(4) .section-title {
  background: #f0fdf4;
  color: #595959;
}

.form-section:nth-child(5) .section-title {
  background: #fefce8;
  color: #595959;
}

/* 对话框样式优化 */
.price-dialog :deep(.el-dialog) {
  border-radius: 8px;
  overflow: hidden;
}

.price-dialog :deep(.el-dialog__header) {
  background: #f5f7fa;
  color: #303133;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.price-dialog :deep(.el-dialog__title) {
  color: #303133;
  font-weight: 600;
  font-size: 16px;
}

.price-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: #909399;
}

.price-dialog :deep(.el-dialog__body) {
  padding: 24px;
  background: #f8f9fa;
}

.price-form {
  max-height: 65vh;
  overflow-y: auto;
  padding-right: 8px;
}

/* 滚动条样式 */
.price-form::-webkit-scrollbar {
  width: 6px;
}

.price-form::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.price-form::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.price-form::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 备注文本样式 */
.desc-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #606266;
  font-size: 13px;
  line-height: 1.4;
}
</style>