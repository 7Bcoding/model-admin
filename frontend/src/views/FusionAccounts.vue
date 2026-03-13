<template>
  <div class="fusion-accounts">
    <div class="toolbar">
      <div class="search-wrapper">
        <el-input
          v-model="searchQuery"
          placeholder="搜索账户名称或供应商"
          prefix-icon="Search"
          clearable
          style="width: 300px;"
        />
      </div>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          创建账户
        </el-button>
      </div>
    </div>

    <div class="table-container">
      <el-table
        :data="paginatedAccounts"
        v-loading="loading"
        style="width: 100%"
        border
      >
        <el-table-column prop="name" label="账户名称" class-name="col-3">
          <template #default="{ row }">
            <span class="account-name">{{ row.name }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="vendor" label="供应商" class-name="col-2">
          <template #header>
            <div class="header-select-container">
              <el-select 
                v-model="selectedVendorFilter" 
                placeholder="供应商" 
                clearable 
                style="width: 100%"
                @change="handleVendorFilterChange"
              >
                <el-option label="全部" value="" />
                <el-option
                  v-for="vendor in availableVendors"
                  :key="vendor"
                  :label="vendor"
                  :value="vendor"
                />
              </el-select>
            </div>
          </template>
          <template #default="{ row }">
            <el-tag :type="getVendorTagType(row.vendor)">
              {{ row.vendor || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="api_type" label="API类型" class-name="col-2">
          <template #header>
            <div class="header-select-container">
              <el-select 
                v-model="selectedApiTypeFilter" 
                placeholder="API类型" 
                clearable 
                style="width: 100%"
                @change="handleApiTypeFilterChange"
              >
                <el-option label="全部" value="" />
                <el-option
                  v-for="apiType in availableApiTypes"
                  :key="apiType.value"
                  :label="apiType.label"
                  :value="apiType.value"
                />
              </el-select>
            </div>
          </template>
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ getApiTypeText(row.api_type || row.apiType) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="description" label="描述" class-name="col-3">
          <template #default="{ row }">
            <span class="description-text">
              {{ row.extras?.description || '-' }}
            </span>
          </template>
        </el-table-column>
        
        <el-table-column label="关联模型" class-name="col-3">
          <template #header>
            <div class="column-header-with-toggle">
              <span>关联模型</span>
              <el-switch
                v-model="showDetailedModels"
                size="small"
                active-text="详情"
                inactive-text="概览"
                style="margin-left: 8px;"
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="related-models-container">
              <div v-if="!modelsLoaded" class="loading-models">
                <span class="loading-text">加载中...</span>
              </div>
              <div v-else-if="getRelatedModels(row.name).length > 0" class="related-models-content">
                <!-- 概览模式 -->
                <div v-if="!showDetailedModels && !getRowDetailedMode(row.name)" class="models-overview" @click="toggleRowDetailedMode(row.name)">
                  <div class="overview-info">
                    <span class="overview-icon">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M12 2C13.1 2 14 2.9 14 4C14 5.1 13.1 6 12 6C10.9 6 10 5.1 10 4C10 2.9 10.9 2 12 2ZM21 9V7L15 1H5C3.9 1 3 1.9 3 3V21C3 22.1 3.9 23 5 23H19C20.1 23 21 22.1 21 21V9ZM19 21H5V3H13V9H19V21Z" fill="#409eff"/>
                      </svg>
                    </span>
                    <span class="overview-text">共在 {{ getRelatedModels(row.name).length }} 个模型中使用</span>
                    <span v-if="hasMultipleUsage(row.name)" class="multiple-usage-indicator" title="存在重复使用情况">
                      <svg width="10" height="10" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="#f56c6c"/>
                      </svg>
                    </span>
                    <span class="expand-hint">点击展开</span>
                  </div>
                </div>
                <!-- 详情模式（全局开关或单行开关） -->
                <div v-else class="related-models-list">
                  <div class="models-header">
                    <span class="models-title">关联模型详情</span>
                    <el-button 
                      type="text" 
                      size="small" 
                      @click="toggleRowDetailedMode(row.name)"
                      class="collapse-btn"
                    >
                      <el-icon><Close /></el-icon>
                    </el-button>
                  </div>
                  <div v-for="(modelInfo, index) in getRelatedModels(row.name)" :key="index" class="model-item">
                    <div class="model-info">
                      <span class="model-name">{{ modelInfo.modelName }}</span>
                      <span v-if="modelInfo.usageCount > 1" class="usage-count-badge" :title="`在${modelInfo.modelName}中被使用${modelInfo.usageCount}次`">
                        {{ modelInfo.usageCount }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="no-related-models">
                <span class="no-models-text">未使用</span>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" class-name="col-2" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="editAccount(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="deleteAccount(row)">
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
          :total="totalFiltered"
          :layout="totalFiltered > pageSize ? 'total, sizes, prev, pager, next, jumper' : 'total, sizes'"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- Create/Edit Account Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑账户' : '创建账户'"
      width="600px"
    >
      <el-form
        ref="accountFormRef"
        :model="accountForm"
        :rules="accountRules"
        label-width="120px"
        autocomplete="off"
      >
        <el-form-item label="账户名称" prop="name">
          <el-input v-model="accountForm.name" :disabled="isEdit" autocomplete="off" />
        </el-form-item>
        
        <el-form-item label="供应商">
          <div class="vendor-select-container">
            <el-select 
              v-model="accountForm.vendor" 
              placeholder="选择或输入供应商" 
              filterable 
              allow-create 
              default-first-option
              style="width: 100%"
              autocomplete="off"
            >
            <!-- 供应商账号管理中的vendor -->
                <el-option-group v-if="vendorOptionsFromManagement.length > 0" label="供应商账号管理">
              <el-option
                v-for="vendor in vendorOptionsFromManagement"
                :key="`managed-${vendor.value}`"
                :label="vendor.label"
                :value="vendor.value"
              >
                <span>{{ vendor.value }}</span>
                <span style="float: right; color: #8cc5ff; font-size: 12px;">
                  <el-icon><Star /></el-icon>
                </span>
              </el-option>
            </el-option-group>
            
            <!-- 现有账户中的vendor -->
            <el-option-group v-if="vendorOptionsFromAccounts.length > 0" label="现有账户">
              <el-option
                v-for="vendor in vendorOptionsFromAccounts"
                :key="`existing-${vendor.value}`"
                :label="vendor.label"
                :value="vendor.value"
              >
                <span>{{ vendor.value }}</span>
                <span style="float: right; color: #a8abb2; font-size: 12px;">已使用</span>
              </el-option>
             </el-option-group>
           </el-select>
           <div class="vendor-help-text">
                <el-icon><Star /></el-icon>
                优先选择供应商账号管理中的名称，可获得更好的管理体验
              </div>
         </div>
       </el-form-item>
        
        <el-form-item label="API类型" prop="api_type">
          <el-select v-model="accountForm.api_type" placeholder="选择API类型" style="width: 100%" autocomplete="off">
            <el-option label="OpenAI" value="API_OPENAI" />
            <el-option label="Anthropic" value="API_ANTHROPIC" />
            <el-option label="Volc" value="API_VOLC" />
            <el-option label="Qwen" value="API_QWEN" />
            <el-option label="OpenRouter" value="API_OPENROUTER" />
            <el-option label="Fireworks" value="API_FIREWORKS" />
            <el-option label="Baidu" value="API_BAIDU" />
            <el-option label="TRT" value="API_TRT" />
            <el-option label="AWS" value="API_AWS" />
            <el-option label="Azure" value="API_AZURE" />
            <el-option label="Vertex" value="API_VERTEX" />
            <el-option label="Vertex Gemini" value="API_VERTEX_GEMINI" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="Base URL" prop="base_url">
          <el-input v-model="accountForm.base_url" placeholder="https://api.provider.com" autocomplete="off" />
        </el-form-item>
        
        <el-form-item label="API Key" prop="api_key">
          <el-input 
            v-if="accountForm.api_type !== 'API_VERTEX' && accountForm.api_type !== 'API_VERTEX_GEMINI'"
            v-model="accountForm.api_key" 
            type="password" 
            show-password 
            :placeholder="getApiKeyPlaceholder()"
            autocomplete="off"
          />
          <el-input 
            v-else
            v-model="accountForm.api_key" 
            type="textarea" 
            :rows="8"
            :placeholder="getApiKeyPlaceholder()"
            autocomplete="off"
          />
        </el-form-item>
        
        <el-form-item v-if="accountForm.api_type === 'API_AZURE'" label="API Version">
          <el-input 
            v-model="accountForm.extras['azure/api-version']" 
            placeholder="2024-02-15-preview"
            autocomplete="off"
          />
        </el-form-item>
        
        <el-form-item label="描述">
          <el-input v-model="accountForm.extras.description" type="textarea" :rows="3" autocomplete="off" />
        </el-form-item>
        
        <el-form-item label="区域">
          <el-select 
            v-model="accountForm.regions" 
            multiple 
            placeholder="选择或输入区域" 
            filterable 
            allow-create 
            default-first-option
            style="width: 100%"
            autocomplete="off"
          >
            <el-option
              v-for="region in availableRegions"
              :key="region"
              :label="region"
              :value="region"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="默认区域">
          <el-input v-model="accountForm.default_region" autocomplete="off" />
        </el-form-item>
        
        <el-form-item label="状态 (待开发)">
          <el-switch v-model="accountForm.disable" active-text="禁用" inactive-text="启用" />
        </el-form-item>
        
        <el-form-item v-if="isEdit" label="创建:">
          <div class="timestamp-info">
            <span class="timestamp">{{ formatTimestamp(accountForm.createdAt) }}</span>
            <span class="timestamp-label">更新:</span>
            <span class="timestamp">{{ formatTimestamp(accountForm.updatedAt) }}</span>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAccount" :loading="saving">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Close, Star } from '@element-plus/icons-vue'
import { getAccounts, createAccount, updateAccount, deleteAccount as deleteAccountApi } from '@/api/fusion-accounts'
import { getModels } from '@/api/fusion-models'
import { getVendors } from '@/api/fusion-vendors'
import { useFusionStore } from '@/stores/fusion'

const fusionStore = useFusionStore()

// Reactive data
const loading = ref(false)
const saving = ref(false)
const accounts = ref([])
const models = ref([]) // 存储模型数据
const modelsLoaded = ref(false) // 标记模型是否已加载
const vendors = ref([]) // 存储供应商数据
const vendorsLoaded = ref(false) // 标记供应商是否已加载
const showDetailedModels = ref(false) // 控制是否显示详细的模型信息
const rowDetailedModes = ref(new Set()) // 控制单行是否显示详细模型信息
const searchQuery = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const accountFormRef = ref()

// Filter states
const selectedVendorFilter = ref('')
const selectedApiTypeFilter = ref('')

// Pagination
const currentPage = ref(1)
const pageSize = ref(parseInt(localStorage.getItem('fusionAccountsPageSize')) || 15)
const total = ref(0)
const nextPageToken = ref('')
const pageTokens = ref({}) // Store page tokens for navigation (use object instead of array)
const totalCalculated = ref(false) // Track if we've calculated the total

const accountForm = ref({
  name: '',
  vendor: '',
  api_type: '',
  base_url: '',
  api_key: '',
  regions: [],
  default_region: '',
  disable: false,
  extras: {
    description: ''
  }
})

const accountRules = {
  name: [
    { required: true, message: '请输入账户名称', trigger: 'submit' },
    { 
      pattern: /^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$/, 
      message: '账户名称只能包含小写字母、数字、连字符和点，且必须以字母或数字开头和结尾', 
      trigger: 'submit' 
    }
  ],
  api_type: [
    { required: true, message: '请选择API类型', trigger: 'submit' }
  ],
  base_url: [
    { required: true, message: '请输入Base URL', trigger: 'submit' }
  ],
  api_key: [
    { required: true, message: '请输入API Key', trigger: 'submit' }
  ]
}

// Computed
const filteredAccounts = computed(() => {
  let filtered = accounts.value
  
  // Apply vendor filter
  if (selectedVendorFilter.value) {
    filtered = filtered.filter(account => account.vendor === selectedVendorFilter.value)
  }
  
  // Apply API type filter
  if (selectedApiTypeFilter.value) {
    filtered = filtered.filter(account => {
      const accountApiType = account.api_type || account.apiType
      return accountApiType === selectedApiTypeFilter.value
    })
  }
  
  // Apply search query
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(account => 
      account.name.toLowerCase().includes(query) ||
      (account.vendor && account.vendor.toLowerCase().includes(query)) ||
      (account.extras?.description && account.extras.description.toLowerCase().includes(query))
    )
  }
  
  return filtered
})

const paginatedAccounts = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return filteredAccounts.value.slice(startIndex, endIndex)
})

const totalFiltered = computed(() => filteredAccounts.value.length)

// 获取账户关联的模型信息
const getRelatedModels = (accountName) => {
  const relatedModels = []
  
  if (!models.value || models.value.length === 0) {
    console.log('No models available')
    return relatedModels
  }
  
  console.log(`=== Checking account: ${accountName} ===`)
  console.log(`Total models to check: ${models.value.length}`)
  console.log('All model names:', models.value.map(m => m.name))
  
  let checkedModels = 0
  let modelsWithProviders = 0
  let modelsWithFallbackProviders = 0
  let totalProviders = 0
  let totalAccounts = 0
  
  models.value.forEach((model, modelIndex) => {
    checkedModels++
    let usageCount = 0
    
    console.log(`\n--- Checking Model ${modelIndex + 1}/${models.value.length}: ${model.name} ---`)
    
    // 检查主要供应商
    if (model.providers && Array.isArray(model.providers)) {
      modelsWithProviders++
      totalProviders += model.providers.length
      console.log(`Primary providers (${model.providers.length}):`, model.providers.map(p => p.name))
      
      model.providers.forEach((provider, providerIndex) => {
        console.log(`  Provider ${providerIndex + 1}: ${provider.name}`)
        
        if (provider.accounts && Array.isArray(provider.accounts)) {
          totalAccounts += provider.accounts.length
          console.log(`    Accounts (${provider.accounts.length}):`, provider.accounts.map(a => a.name))
          
          provider.accounts.forEach((account, accountIndex) => {
            console.log(`      Account ${accountIndex + 1}: "${account.name}" vs "${accountName}"`)
            if (account.name === accountName) {
              usageCount++
              console.log(`      *** FOUND MATCH in provider ${provider.name}! ***`)
            }
          })
        } else {
          console.log(`    No accounts found in provider ${provider.name}`)
        }
      })
    } else {
      console.log('No primary providers found')
    }
    
    // 检查备用供应商
    const fallbackProviders = model.fallback_providers || model.fallbackProviders
    if (fallbackProviders && Array.isArray(fallbackProviders)) {
      modelsWithFallbackProviders++
      totalProviders += fallbackProviders.length
      console.log(`Fallback providers (${fallbackProviders.length}):`, fallbackProviders.map(p => p.name))
      
      fallbackProviders.forEach((provider, providerIndex) => {
        console.log(`  Fallback Provider ${providerIndex + 1}: ${provider.name}`)
        
        if (provider.accounts && Array.isArray(provider.accounts)) {
          totalAccounts += provider.accounts.length
          console.log(`    Accounts (${provider.accounts.length}):`, provider.accounts.map(a => a.name))
          
          provider.accounts.forEach((account, accountIndex) => {
            console.log(`      Fallback Account ${accountIndex + 1}: "${account.name}" vs "${accountName}"`)
            if (account.name === accountName) {
              usageCount++
              console.log(`      *** FOUND MATCH in fallback provider ${provider.name}! ***`)
            }
          })
        } else {
          console.log(`    No accounts found in fallback provider ${provider.name}`)
        }
      })
    } else {
      console.log('No fallback providers found')
    }
    
    if (usageCount > 0) {
      relatedModels.push({
        modelName: model.name,
        usageCount: usageCount
      })
      console.log(`*** Added model ${model.name} with usage count ${usageCount} ***`)
    }
  })
  
  console.log(`\n=== Summary ===`)
  console.log(`Checked ${checkedModels} models`)
  console.log(`Models with primary providers: ${modelsWithProviders}`)
  console.log(`Models with fallback providers: ${modelsWithFallbackProviders}`)
  console.log(`Total providers: ${totalProviders}`)
  console.log(`Total accounts: ${totalAccounts}`)
  console.log(`Final related models:`, relatedModels)
  
  return relatedModels
}

// 检查账户是否有重复使用情况
const hasMultipleUsage = (accountName) => {
  const relatedModels = getRelatedModels(accountName)
  return relatedModels.some(model => model.usageCount > 1)
}

// 获取单行的详细模式状态
const getRowDetailedMode = (accountName) => {
  return rowDetailedModes.value.has(accountName)
}

// 切换单行的详细模式
const toggleRowDetailedMode = (accountName) => {
  if (rowDetailedModes.value.has(accountName)) {
    rowDetailedModes.value.delete(accountName)
  } else {
    rowDetailedModes.value.add(accountName)
  }
}

// Filter data
const vendorFilters = computed(() => {
  const vendors = new Set()
  accounts.value.forEach(account => {
    const vendor = account.vendor
    if (vendor && vendor.trim()) {
      vendors.add(vendor.trim())
    }
  })
  return Array.from(vendors).map(vendor => ({
    text: vendor,
    value: vendor
  }))
})

const apiTypeFilters = computed(() => {
  const apiTypes = new Set()
  accounts.value.forEach(account => {
    const apiType = account.api_type || account.apiType
    if (apiType) {
      apiTypes.add(apiType)
    }
  })
  return Array.from(apiTypes).map(apiType => ({
    text: getApiTypeText(apiType),
    value: apiType
  }))
})

// 供应商管理中的vendor选项
const vendorOptionsFromManagement = computed(() => {
  const vendorSet = new Set()
  vendors.value.forEach(vendor => {
    if (vendor.name && vendor.name.trim()) {
      vendorSet.add(vendor.name.trim())
    }
  })
  return Array.from(vendorSet).sort().map(vendor => ({
    label: vendor,
    value: vendor
  }))
})

// 现有账户中的vendor选项（排除已在供应商管理中的）
const vendorOptionsFromAccounts = computed(() => {
  const managedVendors = new Set(vendorOptionsFromManagement.value.map(v => v.value))
  const vendorSet = new Set()
  
  accounts.value.forEach(account => {
    const vendor = account.vendor
    if (vendor && vendor.trim() && !managedVendors.has(vendor.trim())) {
      vendorSet.add(vendor.trim())
    }
  })
  
  return Array.from(vendorSet).sort().map(vendor => ({
    label: vendor,
    value: vendor
  }))
})

// Computed for available vendors - 从供应商管理中获取
const availableVendors = computed(() => {
  const vendorSet = new Set()
  
  // 首先添加供应商管理中的name字段
  vendors.value.forEach(vendor => {
    if (vendor.name && vendor.name.trim()) {
      vendorSet.add(vendor.name.trim())
    }
  })
  
  // 然后添加现有账户中的vendor（向后兼容）
  accounts.value.forEach(account => {
    const vendor = account.vendor
    if (vendor && vendor.trim()) {
      vendorSet.add(vendor.trim())
    }
  })
  
  return Array.from(vendorSet).sort()
})

// Computed for available API types
const availableApiTypes = computed(() => {
  const apiTypes = new Set()
  accounts.value.forEach(account => {
    const apiType = account.api_type || account.apiType
    if (apiType) {
      apiTypes.add(apiType)
    }
  })
  return Array.from(apiTypes).map(apiType => ({
    value: apiType,
    label: getApiTypeText(apiType)
  })).sort((a, b) => a.label.localeCompare(b.label))
})

// Computed for available regions (filtered by same API type)
const availableRegions = computed(() => {
  const regions = new Set()
  const currentApiType = accountForm.value.api_type
  
  if (currentApiType) {
    accounts.value.forEach(account => {
      const accountApiType = account.api_type || account.apiType
      if (accountApiType === currentApiType) {
        const accountRegions = account.regions || []
        accountRegions.forEach(region => {
          if (region && region.trim()) {
            regions.add(region.trim())
          }
        })
      }
    })
  }
  
  return Array.from(regions).sort()
})

// Methods
const loadAllAccounts = async () => {
  loading.value = true
  try {
    const allAccounts = []
    let pageToken = ''
    let hasMore = true
    
    // Load all accounts by iterating through pages
    while (hasMore) {
      const params = {
        pageSize: 20, // Use backend's default page size
        pageToken: pageToken
      }
      
      const response = await getAccounts(params)
      const accounts = response.data.accounts || []
      allAccounts.push(...accounts)
      
      // Check if there are more pages
      const nextToken = response.data.nextPageToken || response.data.next_page_token || ''
      if (nextToken && accounts.length > 0) {
        pageToken = nextToken
      } else {
        hasMore = false
      }
    }
    
    accounts.value = allAccounts
    total.value = allAccounts.length
    totalCalculated.value = true
    
    console.log(`Loaded ${allAccounts.length} accounts`)
  } catch (error) {
    console.error('Failed to load accounts:', error)
    ElMessage.error('加载账户列表失败')
  } finally {
    loading.value = false
  }
}

const loadAllModels = async () => {
  // 如果已经加载过，直接返回
  if (modelsLoaded.value) {
    return
  }
  
  try {
    console.log('Starting to load all models...')
    
    // 获取所有模型（分页加载）
    const allModels = []
    let pageToken = ''
    let hasMore = true
    let pageCount = 0
    
    while (hasMore) {
      pageCount++
      console.log(`Loading models page ${pageCount}...`)
      
      const params = {
        pageSize: 50, // 使用较大的页面大小
        pageToken: pageToken
      }
      
      const response = await getModels(params)
      const models = response.data.models || []
      allModels.push(...models)
      
      console.log(`Page ${pageCount}: loaded ${models.length} models`)
      
      // 检查是否还有更多页面
      const nextToken = response.data.nextPageToken || response.data.next_page_token || ''
      if (nextToken && models.length > 0) {
        pageToken = nextToken
      } else {
        hasMore = false
      }
    }
    
    models.value = allModels
    modelsLoaded.value = true
    console.log(`Total loaded ${allModels.length} models from ${pageCount} pages`)
    
    // 打印所有模型的名称
    console.log('All model names:', allModels.map(m => m.name))
    
    // 调试：打印第一个模型的数据结构
    if (allModels.length > 0) {
      console.log('Sample model data:', allModels[0])
      console.log('Sample model providers:', allModels[0].providers)
      if (allModels[0].providers && allModels[0].providers.length > 0) {
        console.log('Sample provider accounts:', allModels[0].providers[0].accounts)
      }
    }
  } catch (error) {
    console.error('Failed to load models:', error)
    // 不显示错误消息，因为模型数据不是必需的
  }
}

const loadAllVendors = async () => {
  // 如果已经加载过，直接返回
  if (vendorsLoaded.value) {
    return
  }
  
  try {
    console.log('Starting to load all vendors...')
    const response = await getVendors()
    vendors.value = response.data.vendors || []
    vendorsLoaded.value = true
    console.log(`Loaded ${vendors.value.length} vendors`)
    
    // 打印供应商的name字段
    console.log('Available vendor names:', vendors.value.map(v => v.name).filter(Boolean))
  } catch (error) {
    console.error('Failed to load vendors:', error)
    // 不显示错误消息，因为供应商数据不是必需的
  }
}

const loadAccounts = async () => {
  await Promise.all([
    loadAllAccounts(),
    loadAllModels(),
    loadAllVendors()
  ])
}

const showCreateDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const editAccount = (account) => {
  isEdit.value = true
  accountForm.value = {
    name: account.name,
    vendor: account.vendor || '',
    api_type: account.api_type || account.apiType || '',
    base_url: account.base_url || account.baseUrl || '',
    api_key: account.api_key || account.apiKey || '',
    regions: account.regions || [],
    default_region: account.default_region || account.defaultRegion || '',
    disable: account.disable || false,
    extras: {
      description: account.extras?.description || '',
      'azure/api-version': account.extras?.['azure/api-version'] || ''
    },
    createdAt: account.createdAt || account.created_at,
    updatedAt: account.updatedAt || account.updated_at
  }
  dialogVisible.value = true
}

const saveAccount = async () => {
  if (!accountFormRef.value) return
  
  try {
    await accountFormRef.value.validate()
    saving.value = true
    
    const accountData = {
      ...accountForm.value,
      extras: {
        description: accountForm.value.extras.description,
        'azure/api-version': accountForm.value.extras['azure/api-version']
      }
    }
    
    if (isEdit.value) {
      await updateAccount(accountForm.value.name, accountData)
      ElMessage.success('账户更新成功')
    } else {
      await createAccount(accountData)
      ElMessage.success('账户创建成功')
    }
    
    dialogVisible.value = false
    await loadAccounts()
  } catch (error) {
    console.error('Failed to save account:', error)
    ElMessage.error(isEdit.value ? '更新账户失败' : '创建账户失败')
  } finally {
    saving.value = false
  }
}

const deleteAccount = async (account) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除账户 "${account.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteAccountApi(account.name)
    ElMessage.success('账户删除成功')
    await loadAccounts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete account:', error)
      
      // 提取详细错误信息
      let errorMessage = '删除账户失败'
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

const resetForm = () => {
  accountForm.value = {
    name: '',
    vendor: '',
    api_type: '',
    base_url: '',
    api_key: '',
    regions: [],
    default_region: '',
    disable: false,
    extras: {
      description: '',
      'azure/api-version': ''
    }
  }
}

const getVendorTagType = (vendor) => {
  const types = {
    'openai': 'success',
    'anthropic': 'warning',
    'aws': 'danger',
    'google': 'info'
  }
  return types[vendor?.toLowerCase()] || 'info'
}

const getApiTypeText = (apiType) => {
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
  return typeMap[apiType] || apiType
}

const getApiKeyPlaceholder = () => {
  const apiType = accountForm.value.api_type
  if (apiType === 'API_AWS') {
    return 'ak-xxx|sk-xxx'
  }
  if (apiType === 'API_ANTHROPIC') {
    return `sk-xxx`
  }
  if (apiType === 'API_VERTEX' || apiType === 'API_VERTEX_GEMINI') {
    return `{
  "type": "service_account",
  "project_id": "myproject",
  "private_key_id": "24136268...d0",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "xxx@yyy.iam.gserviceaccount.com",
  "client_id": "123456789",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/xxx.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`
  }
  return '请输入API Key'
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timestamp
  }
}

const filterVendor = (value, row) => {
  return row.vendor === value
}

const filterApiType = (value, row) => {
  const rowApiType = row.api_type || row.apiType
  return rowApiType === value
}

const handleVendorFilterChange = () => {
  currentPage.value = 1
}

const handleApiTypeFilterChange = () => {
  currentPage.value = 1
}

// Pagination handlers
const handleSizeChange = (newSize) => {
  pageSize.value = newSize
  localStorage.setItem('fusionAccountsPageSize', newSize.toString())
  currentPage.value = 1
}

const handleCurrentChange = (newPage) => {
  currentPage.value = newPage
}

// Watch for backend changes
watch(() => fusionStore.selectedBackend, () => {
  loadAccounts()
})

// Watch for search query changes to reset pagination
watch(searchQuery, () => {
  currentPage.value = 1
})

// Watch for localStorage changes (in case user changes pageSize in another tab)
const checkLocalStorage = () => {
  const storedValue = localStorage.getItem('fusionAccountsPageSize')
  if (storedValue && parseInt(storedValue) !== pageSize.value) {
    pageSize.value = parseInt(storedValue)
  }
}

// Check localStorage periodically and on window focus
setInterval(checkLocalStorage, 1000)
window.addEventListener('focus', checkLocalStorage)

// Lifecycle
onMounted(() => {
  loadAccounts()
})

// Cleanup on unmount
onUnmounted(() => {
  clearInterval(checkLocalStorage)
  window.removeEventListener('focus', checkLocalStorage)
})
</script>

<style scoped>
.fusion-accounts {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.table-container {
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.account-name {
  font-weight: 500;
  color: #303133;
}

.description-text {
  color: #606266;
  font-size: 14px;
}

/* 表格列宽度样式 */
.fusion-accounts .el-table .col-2 {
  width: 16.66%;
}

.fusion-accounts .el-table .col-3 {
  width: 25%;
}

.fusion-accounts .el-table .col-4 {
  width: 33.33%;
}

/* 关联模型列样式 */
.related-models-container {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.no-related-models {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
}

.no-models-text {
  color: #c0c4cc;
  font-size: 11px;
  font-style: italic;
}

.loading-models {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
}

.loading-text {
  color: #909399;
  font-size: 11px;
  font-style: italic;
}

/* 列标题样式 */
.column-header-with-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.column-header-with-toggle span {
  font-weight: 600;
  color: #303133;
}

/* 概览模式样式 */
.models-overview {
  padding: 4px 0;
}

.overview-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border: 1px solid #bae6fd;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.overview-info:hover {
  background: linear-gradient(135deg, #e0f2fe 0%, #bae6fd 100%);
  border-color: #7dd3fc;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
}

.expand-hint {
  color: #409eff;
  font-size: 10px;
  font-style: italic;
  margin-left: auto;
}

/* 详情模式样式 */
.related-models-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.models-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 4px;
}

.models-title {
  font-size: 11px;
  font-weight: 500;
  color: #606266;
}

.collapse-btn {
  padding: 2px;
  color: #909399;
  transition: color 0.2s ease;
}

.collapse-btn:hover {
  color: #f56c6c;
}

.overview-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}

.overview-text {
  color: #0369a1;
  font-size: 11px;
  font-weight: 500;
  flex: 1;
}

.multiple-usage-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 12px;
  height: 12px;
  flex-shrink: 0;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
  100% {
    opacity: 1;
  }
}

/* 详情模式样式优化 */
.related-models-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.related-models-list {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.model-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 6px;
  background-color: #f8fafc;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
}

.model-item:hover {
  background-color: #f1f5f9;
  border-color: #cbd5e1;
  transform: translateY(-1px);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.model-info {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
}

.model-name {
  font-size: 11px;
  color: #334155;
  font-weight: 500;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.usage-count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: linear-gradient(135deg, #f56c6c 0%, #e53e3e 100%);
  color: white;
  border-radius: 8px;
  font-size: 9px;
  font-weight: 600;
  flex-shrink: 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.timestamp-info {
  display: flex;
  align-items: center;
  gap: 15px;
  justify-content: flex-start;
}

.timestamp-label {
  color: #606266;
  font-size: 14px;
  font-weight: 500;
}

.timestamp {
  color: #909399;
  font-size: 14px;
  font-family: monospace;
}

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 20px;
  background: white;
  border-top: 1px solid #ebeef5;
}

/* Header select styles */
.header-select-container {
  width: 100%;
}

.header-select-container .el-select {
  width: 100%;
}

.header-select-container .el-select .el-input__wrapper {
  border-radius: 6px;
  border: 1px solid #dcdfe6;
  transition: all 0.3s ease;
}

.header-select-container .el-select .el-input__wrapper:hover {
  border-color: #409eff;
}

.header-select-container .el-select .el-input__wrapper.is-focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

/* Beautiful responsive select styles - only for form selects */
.fusion-accounts .el-dialog .el-select {
  width: 100%;
}

.fusion-accounts .el-dialog .el-select .el-input__wrapper {
  border-radius: 8px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.fusion-accounts .el-dialog .el-select .el-input__wrapper:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: translateY(-1px);
}

.fusion-accounts .el-dialog .el-select .el-input__wrapper.is-focus {
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

@media (max-width: 768px) {
  .fusion-accounts .el-dialog .el-select .el-input__wrapper {
    font-size: 14px;
    padding: 8px 12px;
  }
  
  .toolbar {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .search-wrapper .el-input {
    width: 100% !important;
  }
}

/* Vendor select styles */
.vendor-select-container {
  width: 100%;
}

.vendor-help-text {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  font-size: 12px;
  color: #8cc5ff;
  font-weight: 500;
}

.vendor-help-text .el-icon {
  font-size: 12px;
  color: #8cc5ff;
}

/* Option group styles */
.el-select-dropdown .el-select-group__title {
  font-weight: 600;
  color: #409eff;
  font-size: 12px;
  padding: 8px 12px 4px;
}

.el-select-dropdown .el-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.el-select-dropdown .el-option:hover {
  background-color: #f5f7fa;
}

@media (max-width: 576px) {
  .fusion-accounts .el-dialog .el-select .el-input__wrapper {
    font-size: 13px;
    padding: 6px 10px;
  }
  
  .fusion-accounts {
    padding: 10px;
  }
  
  .el-dialog {
    width: 95% !important;
    margin: 5vh auto;
  }
  
  .vendor-help-text {
    font-size: 11px;
  }
}
</style>