<template>
  <div class="fusion-models">
    <div class="toolbar">
      <div class="search-wrapper">
        <el-input
          v-model="searchQuery"
          placeholder="搜索模型名称"
          prefix-icon="Search"
          clearable
          style="width: 300px;"
        />
        <el-select
          v-model="selectedTags"
          placeholder="按标签筛选"
          multiple
          collapse-tags
          collapse-tags-tooltip
          clearable
          style="width: 300px; margin-left: 12px;"
        >
          <el-option
            v-for="tag in allTags"
            :key="tag"
            :label="tag"
            :value="tag"
          />
        </el-select>
      </div>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          创建模型
        </el-button>
      </div>
    </div>

    <div class="table-container">
      <el-table
        :data="paginatedModels"
        v-loading="loading"
        style="width: 100%"
        border
      >
        <el-table-column prop="name" label="模型名称" class-name="col-1" sortable>
          <template #default="{ row }">
            <span class="model-name">{{ row.name }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="主要供应商" class-name="col-1" prop="primaryProviderCount" sortable>
          <template #default="{ row }">
            <div class="provider-overview">
              <div v-if="row.providers && row.providers.length > 0" class="provider-content">
                <div class="provider-progress-container">
                  <div class="provider-progress-bar">
                    <div 
                      v-for="(provider, index) in row.providers" 
                      :key="index" 
                      class="provider-progress-segment" 
                      :style="getProviderProgressStyle(provider, row.providers)"
                      :title="`${provider.name}: ${getProviderWeightPercentage(provider.weight, row.providers)}%`"
                      @mouseenter="highlightProvider(index, 'primary')"
                      @mouseleave="unhighlightProvider()"
                    ></div>
                  </div>
                  <div class="provider-labels">
                    <span class="provider-count-badge">{{ row.providers?.length || 0 }}</span>
                    <div v-for="(provider, index) in row.providers" :key="index" class="provider-label-item" :style="getProviderLabelStyle(provider, row.providers)" :class="{ 'highlighted': highlightedProviderIndex === index && highlightedProviderType === 'primary' }" @mouseenter="showProviderTooltip(provider, $event)" @mouseleave="hideProviderTooltip">
                      <div class="provider-color-indicator" :style="{ backgroundColor: getProviderColor(provider, row.providers) }"></div>
                      <span class="provider-name">{{ provider.name }}</span>
                      <span class="provider-weight">{{ getProviderWeightPercentage(provider.weight, row.providers) }}%</span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else>
                暂无主要供应商
              </div>
            </div>
          </template>
        </el-table-column>
        

        
        <el-table-column label="标签" class-name="col-1" prop="tags">
          <template #default="{ row }">
            <div class="tags-container">
              <el-tag
                v-for="tag in (row.tags || []).filter(tag => tag && tag.trim() !== '')"
                :key="tag"
                size="small"
                class="model-tag"
                :style="{ backgroundColor: getTagColor(tag), borderColor: getTagColor(tag), color: '#606266' }"
              >
                {{ tag }}
              </el-tag>
              <span v-if="!row.tags || row.tags.length === 0 || (row.tags || []).filter(tag => tag && tag.trim() !== '').length === 0" class="no-tags">
                -
              </span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="备用供应商" class-name="col-1" prop="fallbackProviderCount" sortable>
          <template #default="{ row }">
            <div class="provider-overview">
              <div v-show="(row.fallback_providers || row.fallbackProviders) && (row.fallback_providers || row.fallbackProviders).length > 0" class="provider-content">
                <div class="provider-progress-container">
                  <div class="provider-progress-bar">
                    <div 
                      v-for="(provider, index) in (row.fallback_providers || row.fallbackProviders)" 
                      :key="index" 
                      class="provider-progress-segment" 
                      :style="getProviderProgressStyle(provider, row.fallback_providers || row.fallbackProviders)"
                      :title="`${provider.name}: ${getProviderWeightPercentage(provider.weight, row.fallback_providers || row.fallbackProviders)}%`"
                      @mouseenter="highlightProvider(index, 'fallback')"
                      @mouseleave="unhighlightProvider()"
                    ></div>
                  </div>
                  <div class="provider-labels">
                    <span class="provider-count-badge">{{ (row.fallback_providers || row.fallbackProviders)?.length || 0 }}</span>
                    <div v-for="(provider, index) in (row.fallback_providers || row.fallbackProviders)" :key="index" class="provider-label-item" :style="getProviderLabelStyle(provider, row.fallback_providers || row.fallbackProviders)" :class="{ 'highlighted': highlightedProviderIndex === index && highlightedProviderType === 'fallback' }" @mouseenter="showProviderTooltip(provider, $event)" @mouseleave="hideProviderTooltip">
                      <div class="provider-color-indicator" :style="{ backgroundColor: getProviderColor(provider, row.fallback_providers || row.fallbackProviders) }"></div>
                      <span class="provider-name">{{ provider.name }}</span>
                      <span class="provider-weight">{{ getProviderWeightPercentage(provider.weight, row.fallback_providers || row.fallbackProviders) }}%</span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-show="!(row.fallback_providers || row.fallbackProviders) || (row.fallback_providers || row.fallbackProviders).length === 0">
                暂无备用供应商
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" class-name="col-1" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="editModel(row)">
              编辑
            </el-button>
            <el-button type="success" size="small" @click="manageProviders(row)">
              管理提供商
            </el-button>
            <el-button type="danger" size="small" @click="deleteModel(row)">
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

    <!-- Create/Edit Model Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑模型' : '创建模型'"
      width="650px"
      class="model-dialog"
    >
      <el-form
        ref="modelFormRef"
        :model="modelForm"
        :rules="modelRules"
        label-width="120px"
        class="model-form"
      >
        <el-form-item label="模型名称" prop="name">
          <el-input v-model="modelForm.name" :disabled="isEdit" placeholder="请输入模型名称" />
        </el-form-item>
        
        <el-form-item label="功能特性">
          <div class="feature-grid-compact">
            <div class="feature-row-compact">
              <div class="feature-item-compact">
                <label class="feature-label-compact">
                  函数调用
                </label>
                <el-switch v-model="modelForm.feature.functionCall" size="small" />
              </div>
              <div class="feature-item-compact">
                <label class="feature-label-compact">
                  Prompt Caching
                </label>
                <el-switch v-model="modelForm.feature.promptCaching" size="small" />
              </div>
            </div>
            
            <div class="feature-row-compact">
              <div class="feature-item-compact">
                <label class="feature-label-compact">
                  最大输入长度
                </label>
                <el-input-number
                  v-model="modelForm.feature.maxInputToken"
                  placeholder="输入长度"
                  :min="0"
                  style="width: 120px;"
                  controls
                  size="small"
                />
              </div>
              <div class="feature-item-compact">
                <label class="feature-label-compact">
                  最大输出长度
                </label>
                <el-input-number
                  v-model="modelForm.feature.maxOutput"
                  placeholder="输出长度"
                  :min="0"
                  style="width: 120px;"
                  controls
                  size="small"
                />
              </div>
            </div>
            
            <div class="feature-row-compact">
              <div class="feature-item-compact">
                <label class="feature-label-compact">
                  字符集编码
                  <el-tooltip content="非deepseek模型请勿填写" placement="top">
                    <el-icon style="margin-left: 4px; color: #909399; cursor: help;">
                      <QuestionFilled />
                    </el-icon>
                  </el-tooltip>
                </label>
                <el-select
                  v-model="modelForm.feature.inputTokenEncoding"
                  placeholder="选择字符集编码"
                  style="width: 120px;"
                  clearable
                  size="small"
                >
                  <el-option
                    v-for="encoding in inputTokenEncodingOptions"
                    :key="encoding"
                    :label="encoding"
                    :value="encoding"
                  />
                </el-select>
              </div>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="标签">
          <div class="tags-edit-container-compact">
            <div class="tags-input-row-compact">
              <el-select
                v-model="selectedExistingTag"
                placeholder="选择已有标签"
                style="width: 160px;"
                clearable
                size="small"
                @change="addExistingTag"
              >
                <el-option
                  v-for="tag in availableTagsForSelection"
                  :key="tag"
                  :label="tag"
                  :value="tag"
                />
              </el-select>
              <el-input
                v-model="newTag"
                placeholder="输入新标签名称"
                style="width: 160px;"
                size="small"
                @keyup.enter="addTag"
                @blur="validateTagInput"
              />
              <el-button type="primary" size="small" @click="addTag">
                添加
              </el-button>
            </div>
            <div class="tags-display-compact">
              <el-tag
                v-for="tag in modelForm.tags.filter(tag => tag && tag.trim() !== '')"
                :key="tag"
                closable
                @close="removeTag(tag)"
                class="editable-tag-compact"
                size="small"
                :style="{ backgroundColor: getTagColor(tag), borderColor: getTagColor(tag), color: '#606266' }"
              >
                {{ tag }}
              </el-tag>
              <span v-if="modelForm.tags.length === 0 || modelForm.tags.filter(tag => tag && tag.trim() !== '').length === 0" class="no-tags-hint">
                -
              </span>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item v-if="isEdit" label="时间信息">
          <div class="timestamp-grid-compact">
            <div class="timestamp-item-compact">
              <span class="timestamp-label-compact">创建时间:</span>
              <span class="timestamp-value-compact">{{ formatTimestamp(modelForm.createdAt) }}</span>
            </div>
            <div class="timestamp-item-compact">
              <span class="timestamp-label-compact">更新时间:</span>
              <span class="timestamp-value-compact">{{ formatTimestamp(modelForm.updatedAt) }}</span>
            </div>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveModel" :loading="saving">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- Provider Management Dialog -->
    <el-dialog
      v-model="providerDialogVisible"
      :title="`管理提供商 ${currentModel.name}`"
      width="800px"
    >
      <div class="provider-management">
        <div class="provider-tabs">
          <el-tabs v-model="activeProviderTab">
            <el-tab-pane label="主要提供商" name="primary">
              <div class="provider-section">
                <div class="section-header">
                  <h4>主要提供商</h4>
                  <el-button type="primary" size="small" @click="showAddProviderDialog(false)">
                    添加提供商
                  </el-button>
                </div>
                <el-table :data="currentModel.providers || []" border :key="`primary-${currentModel.name}-${(currentModel.providers || []).length}`">
                  <el-table-column prop="name" label="提供商名称" class-name="col-4" />
                  <el-table-column label="权重" class-name="col-3">
                    <template #default="{ row }">
                      <el-input-number
                        v-model="row.weight"
                        :min="0"
                        size="small"
                        @change="updateProviderWeight(row, false)"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" class-name="col-4">
                    <template #default="{ row }">
                      <el-button type="primary" size="small" @click="editProvider(row, false)">
                        编辑
                      </el-button>
                      <el-button type="danger" size="small" @click="removeProvider(row, false)">
                        删除
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="备用提供商" name="fallback">
              <div class="provider-section">
                <div class="section-header">
                  <h4>备用提供商</h4>
                  <el-button type="primary" size="small" @click="showAddProviderDialog(true)">
                    添加备用提供商
                  </el-button>
                </div>
                <el-table :data="(currentModel.fallback_providers || currentModel.fallbackProviders) || []" border :key="`fallback-${currentModel.name}-${((currentModel.fallback_providers || currentModel.fallbackProviders) || []).length}`">
                  <el-table-column prop="name" label="提供商名称" class-name="col-4" />
                  <el-table-column label="权重" class-name="col-3">
                    <template #default="{ row }">
                      <el-input-number
                        v-model="row.weight"
                        :min="0"
                        size="small"
                        @change="updateProviderWeight(row, true)"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" class-name="col-4">
                    <template #default="{ row }">
                      <el-button type="primary" size="small" @click="editProvider(row, true)">
                        编辑
                      </el-button>
                      <el-button type="danger" size="small" @click="removeProvider(row, true)">
                        删除
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </el-dialog>

    <!-- Add/Edit Provider Dialog -->
    <el-dialog
      v-model="providerFormVisible"
      :title="isEditProvider ? '编辑提供商' : '添加提供商'"
      width="700px"
      class="provider-dialog"
    >
      <el-form
        ref="providerFormRef"
        :model="providerForm"
        :rules="providerRules"
        label-width="100px"
        class="provider-form"
      >
        <!-- 基本信息 Section -->
        <div class="form-section">
          <div class="section-title">
            <h3>基本信息</h3>
            <div class="section-divider"></div>
          </div>
          <div class="section-content basic-info-compact">
            <!-- 第一行：提供商名称和权重 -->
            <div class="basic-info-row-compact">
              <el-form-item label="提供商名称" prop="name" class="provider-name-compact">
                <template #label>
                  <span style="white-space: nowrap !important; word-break: keep-all !important; word-wrap: normal !important; overflow-wrap: normal !important; hyphens: none !important; display: inline-block !important; overflow: hidden !important; text-overflow: ellipsis !important; max-width: 80px !important; min-width: 80px !important; width: 80px !important; flex-shrink: 0 !important; text-align: left !important; padding: 0 6px 0 0 !important; margin: 0 !important; line-height: 28px; font-size: 12px; font-weight: 500;">提供商名称</span>
                </template>
                <el-input v-model="providerForm.name" :disabled="isEditProvider" />
              </el-form-item>
              
              <el-form-item label="权重" prop="weight" class="weight-compact">
                <el-input-number v-model="providerForm.weight" :min="0" style="width: 100%;" :controls="true" :precision="0" />
              </el-form-item>
            </div>
            
            <!-- 第二行：提供商模型名 -->
            <div class="basic-info-row-compact model-name-row">
              <el-form-item label="提供商模型名" prop="modelName" class="model-name-compact">
                <template #label>
                  <span style="white-space: nowrap !important; word-break: keep-all !important; word-wrap: normal !important; overflow-wrap: normal !important; hyphens: none !important; display: inline-block !important; overflow: hidden !important; text-overflow: ellipsis !important; max-width: 80px !important; min-width: 80px !important; width: 80px !important; flex-shrink: 0 !important; text-align: left !important; padding: 0 6px 0 0 !important; margin: 0 !important; line-height: 28px; font-size: 12px; font-weight: 500;">提供商模型名</span>
                </template>
                <el-input 
                  v-model="providerForm.modelName" 
                  :placeholder="`${currentModel.name}`" 
                  class="model-name-input-compact"
                />
              </el-form-item>
            </div>
          </div>
          
          <!-- 账户权重分布长条 -->
          <div class="weight-distribution-bar">
            <div class="weight-bar-title">账户权重分布</div>
            <div class="weight-bar-container">
              <div v-if="providerForm.accounts && providerForm.accounts.length > 0" class="shared-weight-bar">
                <!-- 账户名称和百分比显示 -->
                <div class="account-labels">
                  <div v-for="(account, index) in providerForm.accounts" :key="index" class="account-label-item" :data-account-index="index">
                    <span class="account-name-bar">{{ getAccountDisplayName(account.name) }}</span>
                    <span class="account-weight-bar">{{ getAccountWeightPercentage(account.weight) }}%</span>
                  </div>
                </div>
                <!-- 共享进度条 -->
                <div class="shared-progress-bar">
                  <div 
                    v-for="(account, index) in providerForm.accounts" 
                    :key="index" 
                    class="progress-segment" 
                    :data-account-index="index"
                    :style="getProgressSegmentStyle(account, index)"
                    :title="`${getAccountDisplayName(account.name)}: ${getAccountWeightPercentage(account.weight)}%`"
                  ></div>
                </div>
              </div>
              <div v-else class="no-accounts-bar">
                暂无账户配置
              </div>
            </div>
          </div>
        </div>
        
        <!-- 功能特性覆盖 Section -->
        <div class="form-section">
          <div class="section-title">
            <h3>功能特性覆盖</h3>
            <div class="section-divider"></div>
          </div>
          <div class="section-content">
            <div class="feature-override-container">
              <!-- 函数调用配置行 -->
              <div class="feature-row">
                <div class="feature-row-content">
                  <div class="feature-label">函数调用</div>
                  <div class="feature-options">
                    <el-radio-group v-model="providerForm.featureOverride.functionCall">
                      <el-radio value="useModel">使用模型默认 ({{ currentModel.feature?.functionCall ? '支持' : '不支持' }})</el-radio>
                      <el-radio :value="true">支持</el-radio>
                      <el-radio :value="false">不支持</el-radio>
                    </el-radio-group>
                  </div>
                </div>
              </div>
              
              <!-- Prompt Caching配置行 -->
              <div class="feature-row">
                <div class="feature-row-content">
                  <div class="feature-label">Prompt Caching</div>
                  <div class="feature-options">
                    <el-radio-group v-model="providerForm.featureOverride.promptCaching">
                      <el-radio value="useModel">使用模型默认 ({{ currentModel.feature?.promptCaching ? '支持' : '不支持' }})</el-radio>
                      <el-radio :value="true">支持</el-radio>
                      <el-radio :value="false">不支持</el-radio>
                    </el-radio-group>
                  </div>
                </div>
              </div>
              
              <!-- 最大输入长度配置行 -->
              <div class="feature-row">
                <div class="feature-row-content">
                  <div class="feature-label">最大输入长度</div>
                  <div class="feature-options token-options">
                    <el-radio-group v-model="providerForm.featureOverride.maxInputToken" class="token-radio-group">
                      <el-radio value="useModel">使用模型默认 ({{ currentModel.feature?.maxInputToken || 0 }})</el-radio>
                      <el-radio value="custom">自定义值</el-radio>
                    </el-radio-group>
                    <el-input
                      v-if="providerForm.featureOverride.maxInputToken === 'custom'"
                      v-model="providerForm.featureOverride.maxInputTokenValue"
                      placeholder="最大输入长度"
                      class="token-input-inline"
                    />
                  </div>
                </div>
              </div>
              
              <!-- 最大输出长度配置行 -->
              <div class="feature-row">
                <div class="feature-row-content">
                  <div class="feature-label">最大输出长度</div>
                  <div class="feature-options token-options">
                    <el-radio-group v-model="providerForm.featureOverride.maxOutput" class="token-radio-group">
                      <el-radio value="useModel">使用模型默认 ({{ currentModel.feature?.maxOutput || 0 }})</el-radio>
                      <el-radio value="custom">自定义值</el-radio>
                    </el-radio-group>
                    <el-input
                      v-if="providerForm.featureOverride.maxOutput === 'custom'"
                      v-model="providerForm.featureOverride.maxOutputValue"
                      placeholder="最大输出长度"
                      class="token-input-inline"
                    />
                  </div>
                </div>
              </div>
              
              <!-- 高级配置 Section -->
              <div class="advanced-config-section">
                <div class="advanced-config-header">
                  <h4>高级配置</h4>
                  <div class="advanced-config-toggle">
                    <el-switch
                      v-model="showAdvancedConfig"
                      size="small"
                      active-text="显示"
                      inactive-text="隐藏"
                    />
                  </div>
                </div>
                
                <div v-if="showAdvancedConfig" class="advanced-config-content">
                  <!-- 支持的请求协议 -->
                  <div class="advanced-config-item">
                    <div class="advanced-config-row">
                      <div class="advanced-config-label">
                        <span>支持的请求协议</span>
                        <el-tooltip content="配置该提供商支持的API端点协议" placement="top">
                          <el-icon style="margin-left: 4px; color: #909399; cursor: help;">
                            <QuestionFilled />
                          </el-icon>
                        </el-tooltip>
                      </div>
                      <div class="advanced-config-controls">
                        <el-select
                          v-model="selectedEndpointToAdd"
                          placeholder="选择要添加的协议"
                          style="width: 200px;"
                          @change="addSupportedEndpoint"
                          clearable
                        >
                          <el-option
                            v-for="endpoint in availableEndpointsToAdd"
                            :key="endpoint"
                            :label="endpoint"
                            :value="endpoint"
                          />
                        </el-select>
                      </div>
                    </div>
                    <div class="advanced-config-list">
                      <div v-for="(endpoint, index) in providerForm.featureOverride.supportedEndpoints" :key="index" class="config-item">
                        <span class="config-item-text">{{ endpoint }}</span>
                        <el-button
                          type="danger"
                          size="small"
                          plain
                          circle
                          @click="removeSupportedEndpoint(index)"
                          class="config-item-close-btn"
                          title="删除此协议"
                        >
                          <el-icon><Close /></el-icon>
                        </el-button>
                      </div>
                      <div v-if="providerForm.featureOverride.supportedEndpoints.length === 0" class="empty-config">
                        <span>暂未配置支持的协议</span>
                      </div>
                    </div>
                  </div>
                  
                  <!-- 结构化输出支持 -->
                  <div class="advanced-config-item">
                    <div class="advanced-config-row">
                      <div class="advanced-config-label">
                        <span>结构化输出支持</span>
                        <el-tooltip content="配置该提供商支持的输出格式，不填表示全支持" placement="top">
                          <el-icon style="margin-left: 4px; color: #909399; cursor: help;">
                            <QuestionFilled />
                          </el-icon>
                        </el-tooltip>
                      </div>
                      <div class="advanced-config-controls">
                        <el-select
                          v-model="selectedResponseFormatToAdd"
                          placeholder="选择要添加的格式"
                          style="width: 200px;"
                          @change="addSupportedResponseFormat"
                          clearable
                        >
                          <el-option
                            v-for="format in availableResponseFormatsToAdd"
                            :key="format"
                            :label="format"
                            :value="format"
                          />
                        </el-select>
                      </div>
                    </div>
                    <div class="advanced-config-list">
                      <div v-for="(format, index) in providerForm.featureOverride.supportedResponseFormats" :key="index" class="config-item">
                        <span class="config-item-text">{{ format }}</span>
                        <el-button
                          type="danger"
                          size="small"
                          plain
                          circle
                          @click="removeSupportedResponseFormat(index)"
                          class="config-item-close-btn"
                          title="删除此格式"
                        >
                          <el-icon><Close /></el-icon>
                        </el-button>
                      </div>
                      <div v-if="providerForm.featureOverride.supportedResponseFormats.length === 0" class="empty-config">
                        <span>暂未配置支持的格式</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 特殊Feature配置 Section -->
        <div class="form-section">
          <div class="section-title">
            <h3>特殊Feature配置</h3>
            <div class="section-divider"></div>
          </div>
          <div class="section-content">
            <div class="special-features-section">
              
              <!-- 添加新Feature的选择器 -->
              <div class="add-feature-section">
                <el-select
                  v-model="selectedFeatureToAdd"
                  placeholder="选择要添加的Feature"
                  style="width: 100%;"
                  @change="addSpecialFeature"
                  clearable
                >
                  <el-option
                    v-for="feature in availableFeaturesToAdd"
                    :key="feature.key"
                    :label="feature.label"
                    :value="feature.key"
                  >
                    <div class="feature-option">
                      <span class="feature-name">{{ feature.label }}</span>
                      <span class="feature-desc">{{ feature.description }}</span>
                    </div>
                  </el-option>
                </el-select>
              </div>
              
              <!-- 已配置的Feature列表 -->
              <div v-for="featureKey in configuredFeatureKeys" :key="featureKey" class="special-feature-item">
                <div class="special-feature-content">
                  <div class="special-feature-row">
                                          <div class="special-feature-label" :style="{ minWidth: getLabelWidth(featureKey) + 'px' }">
                        <label>{{ getFeatureLabel(featureKey) }}</label>
                        <el-tooltip :content="getFeatureDescription(featureKey)" placement="top">
                          <el-icon style="margin-left: 4px; color: #909399; cursor: help;">
                            <QuestionFilled />
                          </el-icon>
                        </el-tooltip>
                      </div>
                      <div class="special-feature-input-wrapper" :style="{ flex: '1 1 ' + getInputWidth(featureKey) + 'px' }">
                        <!-- 根据Feature类型显示不同的输入组件 -->
                        <el-input
                          v-if="getFeatureType(featureKey) === 'string'"
                          v-model="providerForm.specialFeatures[featureKey]"
                          class="special-feature-input"
                          clearable
                          :placeholder="getFeaturePlaceholder(featureKey)"
                        />
                        <el-input-number
                          v-else-if="getFeatureType(featureKey) === 'integer'"
                          v-model="providerForm.specialFeatures[featureKey]"
                          class="special-feature-input"
                          :min="0"
                          :placeholder="getFeaturePlaceholder(featureKey)"
                        />
                        <el-switch
                          v-else-if="getFeatureType(featureKey) === 'boolean'"
                          v-model="providerForm.specialFeatures[featureKey]"
                          active-text="启用"
                          inactive-text="禁用"
                          class="special-feature-switch"
                        />
                        <el-input
                          v-else
                          v-model="providerForm.specialFeatures[featureKey]"
                          class="special-feature-input"
                          clearable
                          :placeholder="getFeaturePlaceholder(featureKey)"
                        />
                      </div>
                  </div>
                  <el-button
                    type="danger"
                    size="small"
                    plain
                    circle
                    @click="removeSpecialFeature(featureKey)"
                    class="special-feature-close-btn"
                    title="删除此配置"
                  >
                    <el-icon><Close /></el-icon>
                  </el-button>
                </div>
              </div>
              
              <!-- 空状态提示 -->
              <div v-if="configuredFeatureKeys.length === 0" class="empty-features">
                <span>暂无特殊Feature配置</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 账户配置 Section -->
        <div class="form-section">
          <div class="section-title">
            <h3>账户配置</h3>
            <div class="section-divider"></div>
          </div>
          <div class="section-content">
            <div class="account-list">
              <div v-for="(account, index) in providerForm.accounts" :key="index" class="account-item">
                <div class="account-grid">
                  <div class="account-select" style="margin-right: 32px;">
                    <label>账户名称:</label>
                    <el-select 
                      v-model="account.name" 
                      placeholder="选择账户" 
                      @change="onAccountChange(index)"
                      filterable
                      :filter-method="filterAccounts"
                      :reserve-keyword="false"
                      class="account-name-select"
                    >
                      <el-option
                        v-for="acc in filteredAccounts"
                        :key="acc.name"
                        :label="acc.name"
                        :value="acc.name"
                      >
                        <div class="account-option">
                          <span class="account-name">{{ acc.name }}</span>
                          <span v-if="acc.extras?.description" class="account-description"> - {{ acc.extras.description }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </div>
                  
                  <div class="account-weight" style="margin-right: 32px;">
                    <label>权重:</label>
                    <el-input-number
                      v-model="account.weight"
                      :min="0"
                      placeholder="权重"
                      class="account-weight-input"
                      style="text-align: left;"
                    />
                  </div>
                  
                  <div v-if="getAccountRegions(account.name).length > 0" class="account-region" style="margin-right: 32px;">
                    <label>区域:</label>
                    <el-select v-model="account.region" placeholder="选择区域" class="account-region-select">
                      <el-option
                        v-for="region in getAccountRegions(account.name)"
                        :key="region"
                        :label="region"
                        :value="region"
                      />
                    </el-select>
                  </div>
                </div>
                
                <!-- 删除按钮移到右上角 -->
                <el-button
                  type="danger"
                  size="small"
                  plain
                  circle
                  @click="removeAccount(index)"
                  class="account-close-btn"
                  title="删除此账户"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
                
                <!-- 账户描述信息显示在下方 -->
                <div v-if="getSelectedAccountDescription(account.name)" class="selected-account-description">
                  {{ getSelectedAccountDescription(account.name) }}
                </div>
              </div>
              
              <div class="add-account-section">
                <el-button type="primary" size="small" @click="addAccount">
                  添加账户
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </el-form>
      
      <template #footer>
        <el-button @click="providerFormVisible = false">取消</el-button>
        <el-button type="primary" @click="saveProvider" :loading="savingProvider">
          {{ isEditProvider ? '更新' : '添加' }}
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 供应商详情弹出框 -->
    <div v-if="tooltipVisible" class="provider-tooltip" :style="tooltipStyle">
      <div class="tooltip-header">
        <h4>{{ currentTooltipProvider?.name }}</h4>
        <div class="tooltip-close" @click="hideProviderTooltip">×</div>
      </div>
      <div class="tooltip-content">
        <div class="accounts-section">
          <h5>账户配置</h5>
          <div v-if="currentTooltipProvider?.accounts && currentTooltipProvider.accounts.length > 0" class="accounts-list">
            <div v-for="(account, index) in currentTooltipProvider.accounts" :key="index" class="account-item-tooltip">
              <div class="account-info">
                <span class="account-name">{{ getAccountDisplayName(account.name) }}</span>
                <span class="account-weight">{{ account.weight }}</span>
                <span v-if="account.region" class="account-region">{{ account.region }}</span>
              </div>
            </div>
          </div>
          <div v-else class="no-accounts">
            暂无账户配置
          </div>
        </div>
        <div class="chart-section">
          <h5>权重分布</h5>
          <div v-if="currentTooltipProvider?.accounts && currentTooltipProvider.accounts.length > 0" class="chart-container">
            <div ref="pieChartRef" class="pie-chart"></div>
          </div>
          <div v-else class="no-chart">
            暂无数据
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, QuestionFilled, Close } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { 
  getModels, 
  createModel, 
  updateModel, 
  deleteModel as deleteModelApi,
  addModelProvider,
  updateModelProvider as updateModelProviderApi,
  removeModelProvider,
  updateModelProviderWeight as updateModelProviderWeightApi
} from '@/api/fusion-models'
import { getAccounts } from '@/api/fusion-accounts'
import { useFusionStore } from '@/stores/fusion'
import { SpecialFeatures, validateFeatureValue, formatFeatureValue, FeatureType, FeatureKeys, SupportedEdnpointFeatures, SupportedResponseFormatFeatures, InputTokenEncodingOptions } from '@/config/specialFeatures'

const fusionStore = useFusionStore()

// Reactive data
const loading = ref(false)
const saving = ref(false)
const savingProvider = ref(false)
const models = ref([])
const availableAccounts = ref([])
const searchQuery = ref('')
const accountFilterQuery = ref('')
const dialogVisible = ref(false)
const providerDialogVisible = ref(false)
const providerFormVisible = ref(false)
const isEdit = ref(false)
const isEditProvider = ref(false)
const activeProviderTab = ref('primary')
const modelFormRef = ref()
const providerFormRef = ref()
const currentModel = ref({})
const weightChartRef = ref()
let weightChart = null

// 高亮状态管理
const highlightedProviderIndex = ref(-1)
const highlightedProviderType = ref('')

// 弹出框状态管理
const tooltipVisible = ref(false)
const currentTooltipProvider = ref(null)
const tooltipStyle = ref({
  left: '0px',
  top: '0px'
})
const pieChartRef = ref()
let pieChart = null

// Pagination
const currentPage = ref(1)
const pageSize = ref(parseInt(localStorage.getItem('fusionModelsPageSize')) || 15)
const total = ref(0)
const nextPageToken = ref('')
const pageTokens = ref({}) // Store page tokens for navigation (use object instead of array)
const totalCalculated = ref(false) // Track if we've calculated the total

const modelForm = ref({
  name: '',
  feature: {
    functionCall: false,
    maxInputToken: 0,
    maxOutput: 0,
    inputTokenEncoding: '',
    promptCaching: false
  },
  tags: []
})

const providerForm = ref({
  name: '',
  modelName: '',
  weight: 100,
  featureOverride: {
    functionCall: 'useModel',
    maxInputToken: 'useModel',
    maxInputTokenValue: '',
    maxOutput: 'useModel',
    maxOutputValue: '',
    promptCaching: 'useModel',
    // 高级配置
    supportedEndpoints: [],
    supportedResponseFormats: []
  },
  specialFeatures: {},
  accounts: []
})

const modelRules = {
  name: [
    { required: true, message: '请输入模型名称', trigger: 'blur' },
    { 
      pattern: /^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$/, 
      message: '模型名称只能包含小写字母、数字、连字符和点，且必须以字母或数字开头和结尾', 
      trigger: 'blur' 
    }
  ]
}

const providerRules = {
  name: [
    { required: true, message: '请输入提供商名称', trigger: 'blur' },
    { 
      pattern: /^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$/, 
      message: '提供商名称只能包含小写字母、数字、连字符和点，且必须以字母或数字开头和结尾', 
      trigger: 'blur' 
    }
  ],
  weight: [
    { required: true, message: '请输入权重', trigger: 'blur' }
  ],
  modelName: [
    { required: true, message: '请输入提供商模型名', trigger: 'blur' }
  ]
}

// Computed
const filteredModels = computed(() => {
  let filtered = models.value
  
  // 文本搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(model => 
      model.name.toLowerCase().includes(query)
    )
  }
  
  // tags过滤
  if (selectedTags.value.length > 0) {
    filtered = filtered.filter(model => {
      if (!model.tags || !Array.isArray(model.tags)) return false
      return selectedTags.value.some(selectedTag => 
        model.tags.includes(selectedTag)
      )
    })
  }
  
  return filtered
})

// Filtered accounts for provider selection
const filteredAccounts = computed(() => {
  if (!accountFilterQuery.value) {
    return availableAccounts.value
  }
  
  const lowercaseQuery = accountFilterQuery.value.toLowerCase()
  return availableAccounts.value.filter(account => 
    account.name.toLowerCase().includes(lowercaseQuery) ||
    (account.extras?.description && account.extras.description.toLowerCase().includes(lowercaseQuery))
  )
})

// Account filtering method for el-select
const filterAccounts = (query) => {
  accountFilterQuery.value = query
  return filteredAccounts.value
}

const paginatedModels = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return filteredModels.value.slice(startIndex, endIndex)
})

const totalFiltered = computed(() => filteredModels.value.length)

// 辅助函数：为模型添加排序所需的计算属性
const addSortingProperties = (model) => ({
  ...model,
  primaryProviderCount: (model.providers || []).length,
  fallbackProviderCount: ((model.fallback_providers || model.fallbackProviders) || []).length
})

// Methods
const loadAllModels = async () => {
  loading.value = true
  try {
    const allModels = []
    let pageToken = ''
    let hasMore = true
    
    // Load all models by iterating through pages
    while (hasMore) {
      const params = {
        pageSize: 20, // Use backend's default page size
        pageToken: pageToken
      }
      
      const response = await getModels(params)
      const models = response.data.models || []
      allModels.push(...models)
      
      // Check if there are more pages
      const nextToken = response.data.nextPageToken || response.data.next_page_token || ''
      if (nextToken && models.length > 0) {
        pageToken = nextToken
      } else {
        hasMore = false
      }
    }
    
    // 为每个模型添加排序所需的计算属性
    const processedModels = allModels.map(addSortingProperties)
    
    models.value = processedModels
    total.value = processedModels.length
    totalCalculated.value = true
    
    console.log(`Loaded ${allModels.length} models`)
  } catch (error) {
    console.error('Failed to load models:', error)
    ElMessage.error('加载模型列表失败')
  } finally {
    loading.value = false
  }
}

const loadModels = async () => {
  await loadAllModels()
}

const loadAccounts = async () => {
  try {
    const allAccounts = []
    let pageToken = ''
    let hasMore = true
    
    // Load all accounts by iterating through pages
    while (hasMore) {
      const params = {
        pageSize: 50,
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
    
    availableAccounts.value = allAccounts
    console.log(`Loaded ${allAccounts.length} accounts for provider selection`)
  } catch (error) {
    console.error('Failed to load accounts:', error)
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  resetModelForm()
  dialogVisible.value = true
}

const editModel = (model) => {
  isEdit.value = true
  // Store the original model for update operations
  currentModel.value = model
  modelForm.value = {
    name: model.name,
    feature: {
      functionCall: Boolean(model.feature?.functionCall),
      maxInputToken: model.feature?.maxInputToken || model.feature?.max_input_token || 0,
      maxOutput: model.feature?.maxOutput || 0,
      inputTokenEncoding: model.feature?.inputTokenEncoding || model.feature?.input_token_encoding || '',
      promptCaching: Boolean(model.feature?.promptCaching)
    },
    tags: (model.tags || []).filter(tag => tag && tag.trim() !== ''),
    createdAt: model.createdAt || model.created_at,
    updatedAt: model.updatedAt || model.updated_at
  }
  dialogVisible.value = true
}

const saveModel = async () => {
  if (!modelFormRef.value) return
  
  try {
    await modelFormRef.value.validate()
    saving.value = true
    
    // 过滤掉空标签
    const filteredTags = modelForm.value.tags.filter(tag => tag && tag.trim() !== '')
    
    if (isEdit.value) {
      // For updates, merge the form data with the original model object
      const updatedModelData = {
        ...currentModel.value, // Include all original model fields
        ...modelForm.value,    // Override with form data
        feature: {
          functionCall: modelForm.value.feature.functionCall ? 1 : 0,
          maxInputToken: modelForm.value.feature.maxInputToken,
          maxOutput: modelForm.value.feature.maxOutput,
          inputTokenEncoding: modelForm.value.feature.inputTokenEncoding || null,
          promptCaching: modelForm.value.feature.promptCaching ? 1 : 0
        },
        tags: filteredTags
      }
      await updateModel(updatedModelData.name, updatedModelData)
      ElMessage.success('模型更新成功')
    } else {
      // For creation, just use the form data
      const modelData = {
        ...modelForm.value,
        feature: {
          functionCall: modelForm.value.feature.functionCall ? 1 : 0,
          maxInputToken: modelForm.value.feature.maxInputToken,
          maxOutput: modelForm.value.feature.maxOutput,
          inputTokenEncoding: modelForm.value.feature.inputTokenEncoding || null,
          promptCaching: modelForm.value.feature.promptCaching ? 1 : 0
        },
        tags: filteredTags
      }
      await createModel(modelData)
      ElMessage.success('模型创建成功')
    }
    
    dialogVisible.value = false
    await loadModels()
  } catch (error) {
    console.error('Failed to save model:', error)
    ElMessage.error(isEdit.value ? '更新模型失败' : '创建模型失败')
  } finally {
    saving.value = false
  }
}

const deleteModel = async (model) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除模型 "${model.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteModelApi(model.name)
    ElMessage.success('模型删除成功')
    await loadModels()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete model:', error)
      
      // 提取详细错误信息
      let errorMessage = '删除模型失败'
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

const manageProviders = async (model) => {
  currentModel.value = model
  await loadAccounts()
  providerDialogVisible.value = true
}

const showAddProviderDialog = (isFallback) => {
  isEditProvider.value = false
  resetProviderForm()
  providerForm.value.isFallback = isFallback
  providerFormVisible.value = true
}

const editProvider = (provider, isFallback) => {
  isEditProvider.value = true
  
  // Handle featureOverride fields - set to useModel if not provided, otherwise use the value
  const featureOverride = {}
  const providerFeatureOverride = provider.feature_override || provider.featureOverride || {}
  
  // Set functionCall - if not provided, use 'useModel', otherwise convert to boolean
  featureOverride.functionCall = providerFeatureOverride.hasOwnProperty('functionCall') 
    ? Boolean(providerFeatureOverride.functionCall) 
    : 'useModel'
  
  // Set maxInputToken - if not provided, use 'useModel', otherwise use 'custom' and set the value
  if (providerFeatureOverride.hasOwnProperty('maxInputToken')) {
    featureOverride.maxInputToken = 'custom'
    featureOverride.maxInputTokenValue = providerFeatureOverride.maxInputToken
  } else {
    featureOverride.maxInputToken = 'useModel'
    featureOverride.maxInputTokenValue = ''
  }
  
  // Set maxOutput - if not provided, use 'useModel', otherwise use 'custom' and set the value
  if (providerFeatureOverride.hasOwnProperty('maxOutput')) {
    featureOverride.maxOutput = 'custom'
    featureOverride.maxOutputValue = providerFeatureOverride.maxOutput
  } else {
    featureOverride.maxOutput = 'useModel'
    featureOverride.maxOutputValue = ''
  }
  
  // Set promptCaching - if not provided, use 'useModel', otherwise convert to boolean
  featureOverride.promptCaching = providerFeatureOverride.hasOwnProperty('promptCaching') 
    ? Boolean(providerFeatureOverride.promptCaching) 
    : 'useModel'
  
  // Extract special features from customFeatures
  const specialFeatures = {}
  if (provider.customFeatures) {
    Object.keys(provider.customFeatures).forEach(key => {
      const value = provider.customFeatures[key]
      const featureType = getFeatureType(key)
      
      if (featureType === 'boolean') {
        // 对于boolean类型，转换为布尔值
        specialFeatures[key] = value === 'true' || value === true
      } else if (featureType === 'integer') {
        // 对于integer类型，转换为数字
        const numValue = parseInt(value)
        specialFeatures[key] = isNaN(numValue) ? '' : numValue
      } else {
        // 对于其他类型，直接使用字符串值
        specialFeatures[key] = value
      }
    })
  }
  
  console.log('Received customFeatures:', provider.customFeatures)
  console.log('Processed specialFeatures:', specialFeatures)
  
  // 处理高级配置
  const supportedEndpoints = provider.feature_override?.supportedEndpoints || provider.featureOverride?.supportedEndpoints || []
  const supportedResponseFormats = provider.feature_override?.supportedResponseFormats || provider.featureOverride?.supportedResponseFormats || []
  
  // 根据高级配置内容决定是否默认展开
  showAdvancedConfig.value = supportedEndpoints.length > 0 || supportedResponseFormats.length > 0
  
  // 将高级配置添加到featureOverride中
  featureOverride.supportedEndpoints = supportedEndpoints
  featureOverride.supportedResponseFormats = supportedResponseFormats
  
  providerForm.value = {
    name: provider.name,
    modelName: provider.model_name || provider.modelName,
    weight: provider.weight,
    featureOverride,
    specialFeatures,
    accounts: provider.accounts || [],
    isFallback
  }
  providerFormVisible.value = true
}

const saveProvider = async () => {
  if (!providerFormRef.value) return
  
  try {
    await providerFormRef.value.validate()
    savingProvider.value = true
    
    // Process accounts to remove empty region values
    const processedAccounts = providerForm.value.accounts.map(account => {
      const processedAccount = {
        name: account.name,
        weight: account.weight
      }
      // Only include region if it has a value
      if (account.region && account.region.trim()) {
        processedAccount.region = account.region.trim()
      }
      return processedAccount
    })
    
    // Build featureOverride object - only include fields that have been explicitly set
    const featureOverride = {}
    
    // Check if functionCall has been explicitly set (not useModel)
    if (providerForm.value.featureOverride.functionCall !== 'useModel') {
      featureOverride.functionCall = providerForm.value.featureOverride.functionCall ? 1 : 0
    }
    
    // Check if maxInputToken has been explicitly set (custom value)
    if (providerForm.value.featureOverride.maxInputToken === 'custom' && 
        providerForm.value.featureOverride.maxInputTokenValue) {
      featureOverride.maxInputToken = providerForm.value.featureOverride.maxInputTokenValue
    }
    
    // Check if maxOutput has been explicitly set (custom value)
    if (providerForm.value.featureOverride.maxOutput === 'custom' && 
        providerForm.value.featureOverride.maxOutputValue) {
      featureOverride.maxOutput = providerForm.value.featureOverride.maxOutputValue
    }
    
    // Check if promptCaching has been explicitly set (not useModel)
    if (providerForm.value.featureOverride.promptCaching !== 'useModel') {
      featureOverride.promptCaching = providerForm.value.featureOverride.promptCaching ? 1 : 0
    }
    
    // Process customFeatures (special features) - 根据类型处理值
    const customFeatures = {}
    Object.keys(providerForm.value.specialFeatures || {}).forEach(key => {
      const value = providerForm.value.specialFeatures[key]
      const featureType = getFeatureType(key)
      
      if (featureType === 'boolean') {
        // 对于boolean类型，转换为字符串
        customFeatures[key] = value ? 'true' : 'false'
      } else if (featureType === 'integer') {
        // 对于integer类型，确保是数字
        if (value !== null && value !== undefined && value !== '') {
          customFeatures[key] = value.toString()
        }
      } else if (value && value.trim() !== '') {
        // 对于其他类型，直接透传用户输入的值
        customFeatures[key] = value.trim()
      }
    })
    
    console.log('Sending customFeatures:', customFeatures)
    
    const providerData = {
      provider: {
        name: providerForm.value.name,
        modelName: providerForm.value.modelName,
        weight: providerForm.value.weight,
        accounts: processedAccounts
      },
      is_fallback: providerForm.value.isFallback
    }
    
    // 处理高级配置 - 直接添加到featureOverride中
    if (providerForm.value.featureOverride.supportedEndpoints.length > 0) {
      featureOverride.supportedEndpoints = providerForm.value.featureOverride.supportedEndpoints
      console.log('Added supportedEndpoints to featureOverride:', providerForm.value.featureOverride.supportedEndpoints)
    }
    if (providerForm.value.featureOverride.supportedResponseFormats.length > 0) {
      featureOverride.supportedResponseFormats = providerForm.value.featureOverride.supportedResponseFormats
      console.log('Added supportedResponseFormats to featureOverride:', providerForm.value.featureOverride.supportedResponseFormats)
    }
    
    console.log('Final featureOverride object:', featureOverride)
    
    // Only include featureOverride if it has any fields (包括高级配置)
    if (Object.keys(featureOverride).length > 0) {
      providerData.provider.featureOverride = featureOverride
      console.log('Including featureOverride in request:', featureOverride)
    } else {
      console.log('No featureOverride fields to include')
    }
    
    // Only include customFeatures if it has any fields
    if (Object.keys(customFeatures).length > 0) {
      providerData.provider.customFeatures = customFeatures
    }
    
    let response
    if (isEditProvider.value) {
      response = await updateModelProviderApi(
        currentModel.value.name,
        providerForm.value.name,
        providerData
      )
      ElMessage.success('提供商更新成功')
    } else {
      response = await addModelProvider(
        currentModel.value.name,
        providerData
      )
      ElMessage.success('提供商添加成功')
    }
    
    providerFormVisible.value = false
    
    // Use the returned model data to update the current model and models list
    if (response && response.data ) {
      const updatedModel = response.data
      
      
      // Force reactivity by creating a new object
      currentModel.value = addSortingProperties(updatedModel)
      
      // Update the model in the models list
      const modelIndex = models.value.findIndex(m => m.name === updatedModel.name)
      if (modelIndex !== -1) {
        models.value[modelIndex] = addSortingProperties(updatedModel)
      }
    }
  } catch (error) {
    console.error('Failed to save provider:', error)
    ElMessage.error(isEditProvider.value ? '更新提供商失败' : '添加提供商失败')
  } finally {
    savingProvider.value = false
  }
}

const removeProvider = async (provider, isFallback) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除提供商 "${provider.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await removeModelProvider(currentModel.value.name, provider.name, isFallback)
    ElMessage.success('提供商删除成功')
    
    // Use the returned model data to update the current model and models list
    if (response && response.data) {
      const updatedModel = response.data
      
      // Force reactivity by creating a new object
      currentModel.value = addSortingProperties(updatedModel)
      
      // Update the model in the models list
      const modelIndex = models.value.findIndex(m => m.name === updatedModel.name)
      if (modelIndex !== -1) {
        models.value[modelIndex] = addSortingProperties(updatedModel)
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove provider:', error)
      const errorMessage = error.response?.data?.message || error.response?.data?.reason || error.message || '删除提供商失败'
      ElMessage.error(errorMessage)
    }
  }
}

const updateProviderWeight = async (provider, isFallback) => {
  try {
    const response = await updateModelProviderWeightApi(
      currentModel.value.name,
      provider.name,
      provider.weight,
      isFallback
    )
    ElMessage.success('权重更新成功')
    
    // Use the returned model data to update the current model and models list
    if (response && response.data) {
      const updatedModel = response.data
      
      // Force reactivity by creating a new object
      currentModel.value = addSortingProperties(updatedModel)
      
      // Update the model in the models list
      const modelIndex = models.value.findIndex(m => m.name === updatedModel.name)
      if (modelIndex !== -1) {
        models.value[modelIndex] = addSortingProperties(updatedModel)
      }
    }
  } catch (error) {
    console.error('Failed to update provider weight:', error)
    ElMessage.error('权重更新失败')
  }
}

const addAccount = () => {
  providerForm.value.accounts.push({
    name: '',
    weight: 100,
    region: ''
  })
}

const removeAccount = (index) => {
  providerForm.value.accounts.splice(index, 1)
}

const removeSpecialFeature = (featureKey) => {
  console.log('Removing feature:', featureKey)
  if (providerForm.value.specialFeatures.hasOwnProperty(featureKey)) {
    delete providerForm.value.specialFeatures[featureKey]
    console.log('Feature removed, current specialFeatures:', providerForm.value.specialFeatures)
  }
}

// 新增：选择要添加的Feature
const selectedFeatureToAdd = ref('')
const showAdvancedConfig = ref(false)
const selectedEndpointToAdd = ref('')
const selectedResponseFormatToAdd = ref('')
const inputTokenEncodingOptions = ref(InputTokenEncodingOptions)
const selectedTags = ref([])
const tagFilterVisible = ref(false)
const newTag = ref('')
const selectedExistingTag = ref('')

// 计算属性：已配置的Feature键列表（显示所有已添加的Feature）
const configuredFeatureKeys = computed(() => {
  const keys = Object.keys(providerForm.value.specialFeatures)
  console.log('Configured feature keys:', keys)
  return keys
})

// 计算属性：可添加的Feature列表（去重）
const availableFeaturesToAdd = computed(() => {
  const configuredKeys = Object.keys(providerForm.value.specialFeatures)
  return SpecialFeatures.filter(feature => !configuredKeys.includes(feature.key))
})

// 计算属性：可添加的端点协议列表（去重）
const availableEndpointsToAdd = computed(() => {
  const configuredEndpoints = providerForm.value.featureOverride.supportedEndpoints
  return SupportedEdnpointFeatures.filter(endpoint => !configuredEndpoints.includes(endpoint))
})

// 计算属性：可添加的响应格式列表（去重）
const availableResponseFormatsToAdd = computed(() => {
  const configuredFormats = providerForm.value.featureOverride.supportedResponseFormats
  return SupportedResponseFormatFeatures.filter(format => !configuredFormats.includes(format))
})

// 计算属性：获取所有已使用的tags
const allTags = computed(() => {
  const tagSet = new Set()
  models.value.forEach(model => {
    if (model.tags && Array.isArray(model.tags)) {
      model.tags.forEach(tag => {
        if (tag && tag.trim() !== '') {
          tagSet.add(tag.trim())
        }
      })
    }
  })
  return Array.from(tagSet).sort()
})

// 计算属性：获取可用于选择的标签（排除已选择的）
const availableTagsForSelection = computed(() => {
  return allTags.value.filter(tag => !modelForm.value.tags.includes(tag))
})

// 添加新的Feature
const addSpecialFeature = () => {
  if (selectedFeatureToAdd.value) {
    console.log('Adding feature:', selectedFeatureToAdd.value)
    // 如果该feature还没有配置，则添加
    if (!providerForm.value.specialFeatures.hasOwnProperty(selectedFeatureToAdd.value)) {
      const featureType = getFeatureType(selectedFeatureToAdd.value)
      // 根据feature类型设置默认值
      if (featureType === 'boolean') {
        providerForm.value.specialFeatures[selectedFeatureToAdd.value] = true
        console.log(`Added boolean feature "${selectedFeatureToAdd.value}" with default value: true`)
      } else if (featureType === 'integer') {
        providerForm.value.specialFeatures[selectedFeatureToAdd.value] = null
        console.log(`Added integer feature "${selectedFeatureToAdd.value}" with default value: null`)
      } else {
        providerForm.value.specialFeatures[selectedFeatureToAdd.value] = ''
        console.log(`Added string feature "${selectedFeatureToAdd.value}" with default value: ""`)
      }
      console.log('Feature added, current specialFeatures:', providerForm.value.specialFeatures)
    }
    // 清空选择器
    selectedFeatureToAdd.value = ''
  }
}

// 添加支持的端点协议
const addSupportedEndpoint = () => {
  if (selectedEndpointToAdd.value) {
    console.log('Adding endpoint:', selectedEndpointToAdd.value)
    if (!providerForm.value.featureOverride.supportedEndpoints.includes(selectedEndpointToAdd.value)) {
      providerForm.value.featureOverride.supportedEndpoints.push(selectedEndpointToAdd.value)
      console.log('Endpoint added, current supportedEndpoints:', providerForm.value.featureOverride.supportedEndpoints)
    }
    selectedEndpointToAdd.value = ''
  }
}

// 删除支持的端点协议
const removeSupportedEndpoint = (index) => {
  providerForm.value.featureOverride.supportedEndpoints.splice(index, 1)
  console.log('Endpoint removed, current supportedEndpoints:', providerForm.value.featureOverride.supportedEndpoints)
}

// 添加支持的响应格式
const addSupportedResponseFormat = () => {
  if (selectedResponseFormatToAdd.value) {
    console.log('Adding response format:', selectedResponseFormatToAdd.value)
    if (!providerForm.value.featureOverride.supportedResponseFormats.includes(selectedResponseFormatToAdd.value)) {
      providerForm.value.featureOverride.supportedResponseFormats.push(selectedResponseFormatToAdd.value)
      console.log('Response format added, current supportedResponseFormats:', providerForm.value.featureOverride.supportedResponseFormats)
    }
    selectedResponseFormatToAdd.value = ''
  }
}

// 删除支持的响应格式
const removeSupportedResponseFormat = (index) => {
  providerForm.value.featureOverride.supportedResponseFormats.splice(index, 1)
  console.log('Response format removed, current supportedResponseFormats:', providerForm.value.featureOverride.supportedResponseFormats)
}

// 添加标签
const addTag = () => {
  const trimmedTag = newTag.value?.trim()
  // 检查标签是否为空、null、undefined或只包含空白字符
  if (trimmedTag && trimmedTag !== '' && !modelForm.value.tags.includes(trimmedTag)) {
    modelForm.value.tags.push(trimmedTag)
    newTag.value = ''
  } else if (!trimmedTag || trimmedTag === '') {
    ElMessage.warning('标签不能为空')
  } else if (modelForm.value.tags.includes(trimmedTag)) {
    ElMessage.warning('标签已存在')
  }
}

// 删除标签
const removeTag = (tag) => {
  const index = modelForm.value.tags.indexOf(tag)
  if (index > -1) {
    modelForm.value.tags.splice(index, 1)
  }
}

// 添加现有标签
const addExistingTag = () => {
  const trimmedTag = selectedExistingTag.value?.trim()
  // 检查标签是否为空、null、undefined或只包含空白字符
  if (trimmedTag && trimmedTag !== '' && !modelForm.value.tags.includes(trimmedTag)) {
    modelForm.value.tags.push(trimmedTag)
    selectedExistingTag.value = ''
  } else if (!trimmedTag || trimmedTag === '') {
    ElMessage.warning('标签不能为空')
  } else if (modelForm.value.tags.includes(trimmedTag)) {
    ElMessage.warning('标签已存在')
  }
}

// 验证标签输入
const validateTagInput = () => {
  if (newTag.value && newTag.value.trim() === '') {
    newTag.value = ''
  }
}

// 获取Feature的标签
const getFeatureLabel = (featureKey) => {
  const feature = SpecialFeatures.find(f => f.key === featureKey)
  return feature ? feature.label : featureKey
}

// 获取Feature的描述
const getFeatureDescription = (featureKey) => {
  const feature = SpecialFeatures.find(f => f.key === featureKey)
  return feature ? feature.description : ''
}

// 获取Feature的类型
const getFeatureType = (featureKey) => {
  return FeatureType[featureKey] || 'string'
}

// 获取Feature的占位符
const getFeaturePlaceholder = (featureKey) => {
  if (featureKey === FeatureKeys.TTFT_TIMEOUT || featureKey === FeatureKeys.TPOT_TIMEOUT) {
    return '例如：1000ms 或 10s'
  }
  if (featureKey === FeatureKeys.LIMIT_MAX_TOKEN) {
    return '请输入最大token数'
  }
  return '请输入值'
}

// 计算label宽度
const getLabelWidth = (featureKey) => {
  const label = getFeatureLabel(featureKey)
  // 估算每个字符的宽度（中文字符约16px，英文字符约8px）
  let width = 0
  for (let char of label) {
    width += /[\u4e00-\u9fa5]/.test(char) ? 16 : 8
  }
  // 加上图标和间距的宽度
  width += 24 // 图标宽度 + 间距
  const finalWidth = Math.max(width, 80) // 最小宽度80px
  console.log(`Label width for ${featureKey}: ${finalWidth}px (label: "${label}")`)
  return finalWidth
}

// 计算输入框宽度
const getInputWidth = (featureKey) => {
  const labelWidth = getLabelWidth(featureKey)
  // 根据label宽度动态调整输入框宽度
  // label越短，输入框越宽
  const baseWidth = 200 // 基础宽度
  const maxWidth = 400 // 最大宽度
  const minWidth = 150 // 最小宽度
  
  // 根据label长度调整输入框宽度
  let inputWidth
  if (labelWidth <= 100) {
    inputWidth = maxWidth
  } else if (labelWidth <= 150) {
    inputWidth = baseWidth
  } else {
    inputWidth = minWidth
  }
  
  console.log(`Input width for ${featureKey}: ${inputWidth}px (label width: ${labelWidth}px)`)
  return inputWidth
}

const resetModelForm = () => {
  modelForm.value = {
    name: '',
    feature: {
      functionCall: false,
      maxInputToken: 0,
      maxOutput: 0,
      inputTokenEncoding: '',
      promptCaching: false
    },
    tags: []
  }
}

const resetProviderForm = () => {
  providerForm.value = {
    name: '',
    modelName: '',
    weight: 100,
    featureOverride: {
      functionCall: 'useModel',
      maxInputToken: 'useModel',
      maxInputTokenValue: '',
      maxOutput: 'useModel',
      maxOutputValue: '',
      promptCaching: 'useModel',
      // 高级配置
      supportedEndpoints: [],
      supportedResponseFormats: []
    },
    specialFeatures: {},
    accounts: []
  }
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

// 格式化字符集编码显示
const formatInputTokenEncoding = (encoding) => {
  if (!encoding || encoding.trim() === '') {
    return '未配置'
  }
  return encoding
}

// 生成标签颜色
const getTagColor = (tag) => {
  if (!tag) return '#909399'
  
  // 使用标签字符串的hash值来生成一致的颜色
  let hash = 0
  for (let i = 0; i < tag.length; i++) {
    const char = tag.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // Convert to 32bit integer
  }
  
  // 预定义的颜色数组 - 更柔和的颜色
  const colors = [
    '#e1f3ff', '#e8f5e8', '#fff3e0', '#ffeaea',
    '#f0f0f0', '#f3e5f5', '#fff8e1', '#efebe9',
    '#e3f2fd', '#e8eaf6', '#e0f2f1', '#f1f8e9'
  ]
  
  return colors[Math.abs(hash) % colors.length]
}

const getAccountRegions = (accountName) => {
  const account = availableAccounts.value.find(acc => acc.name === accountName)
  return account?.regions || []
}

const getSelectedAccountDescription = (accountName) => {
  const account = availableAccounts.value.find(acc => acc.name === accountName)
  return account?.extras?.description || ''
}

const onAccountChange = (index) => {
  const account = providerForm.value.accounts[index]
  const selectedAccount = availableAccounts.value.find(acc => acc.name === account.name)
  
  if (selectedAccount) {
    // Set the preferred region to account's defaultRegion
    account.region = selectedAccount.defaultRegion || selectedAccount.default_region || ''
  }
}

// 获取账户显示名称
const getAccountDisplayName = (accountId) => {
  const account = availableAccounts.value.find(acc => acc.name === accountId)
  const displayName = account?.name || accountId || '未知账户'
  
  // 如果账户名称太长，进行截断
  if (displayName.length > 20) {
    return displayName.substring(0, 17) + '...'
  }
  return displayName
}

// 获取账户颜色
const getAccountColor = (index) => {
  const colors = [
    '#409eff', '#67c23a', '#e6a23c', '#f56c6c', 
    '#909399', '#c71585', '#ff6347', '#32cd32',
    '#1e90ff', '#ff69b4', '#ffd700', '#8a2be2'
  ]
  return colors[index % colors.length]
}

// 计算账户权重百分比
const getAccountWeightPercentage = (weight) => {
  const accounts = providerForm.value.accounts || []
  const totalWeight = accounts.reduce((sum, acc) => sum + (acc.weight || 0), 0)
  
  if (totalWeight === 0) return 0
  
  return Math.round((weight / totalWeight) * 100)
}

// 计算进度条段的起始位置
const getProgressSegmentStart = (index) => {
  const accounts = providerForm.value.accounts || []
  let startPosition = 0
  
  for (let i = 0; i < index; i++) {
    startPosition += getAccountWeightPercentage(accounts[i].weight)
  }
  
  return startPosition
}

// 计算供应商权重百分比
const getProviderWeightPercentage = (weight, providers) => {
  if (!providers || providers.length === 0) return 0
  
  const totalWeight = providers.reduce((sum, provider) => sum + (provider.weight || 0), 0)
  
  if (totalWeight === 0) return 0
  
  return Math.round((weight / totalWeight) * 100)
}

// 获取供应商颜色
const getProviderColor = (provider, providers) => {
  const colors = [
    '#409eff', '#67c23a', '#e6a23c', '#f56c6c', 
    '#909399', '#c71585', '#ff6347', '#32cd32',
    '#1e90ff', '#ff69b4', '#ffd700', '#8a2be2'
  ]
  
  return colors[providers.indexOf(provider) % colors.length]
}

// 计算供应商进度条样式
const getProviderProgressStyle = (provider, providers) => {
  const percentage = getProviderWeightPercentage(provider.weight, providers)
  
  return {
    flex: percentage,
    backgroundColor: getProviderColor(provider, providers)
  }
}

// 计算供应商标签样式
const getProviderLabelStyle = (provider, providers) => {
  // 根据供应商名称长度动态调整宽度，确保完整显示
  const nameLength = provider.name.length
  const charWidth = 8 // 减少每个字符的宽度估算
  const padding = 16 // 为颜色指示器和权重留出空间
  const baseWidth = Math.max(nameLength * charWidth + padding, 60) // 最小60px
  
  return {
    minWidth: baseWidth + 'px',
    maxWidth: 'none' // 移除最大宽度限制，允许完整显示
  }
}

// 高亮供应商
const highlightProvider = (index, type) => {
  highlightedProviderIndex.value = index
  highlightedProviderType.value = type
}

// 取消高亮
const unhighlightProvider = () => {
  highlightedProviderIndex.value = -1
  highlightedProviderType.value = ''
}

// 显示供应商详情弹出框
const showProviderTooltip = (provider, event) => {
  currentTooltipProvider.value = provider
  tooltipVisible.value = true
  
  // 计算弹出框位置
  const rect = event.target.getBoundingClientRect()
  const tooltipWidth = 320
  const tooltipHeight = 450
  
  let left = rect.left + rect.width + 10
  let top = rect.top - tooltipHeight / 2
  
  // 确保弹出框不超出视窗
  if (left + tooltipWidth > window.innerWidth) {
    left = rect.left - tooltipWidth - 10
  }
  if (left < 0) {
    left = 10
  }
  if (top < 0) {
    top = 10
  }
  if (top + tooltipHeight > window.innerHeight) {
    top = window.innerHeight - tooltipHeight - 10
  }
  
  tooltipStyle.value = {
    left: left + 'px',
    top: top + 'px'
  }
  
  // 延迟初始化饼图，确保DOM已渲染
  nextTick(() => {
    initPieChart()
  })
}

// 隐藏供应商详情弹出框
const hideProviderTooltip = () => {
  tooltipVisible.value = false
  currentTooltipProvider.value = null
  if (pieChart) {
    pieChart.dispose()
    pieChart = null
  }
}

// 初始化饼图
const initPieChart = () => {
  if (!pieChartRef.value || !currentTooltipProvider.value?.accounts) return
  
  if (pieChart) {
    pieChart.dispose()
  }
  
  pieChart = echarts.init(pieChartRef.value)
  
  const accounts = currentTooltipProvider.value.accounts
  const totalWeight = accounts.reduce((sum, acc) => sum + (acc.weight || 0), 0)
  
  const data = accounts.map((account, index) => ({
    name: getAccountDisplayName(account.name),
    value: account.weight,
    itemStyle: {
      color: getAccountColor(index)
    }
  }))
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
      confine: true
    },
    legend: {
      orient: 'vertical',
      left: '5%',
      top: 'middle',
      textStyle: {
        fontSize: 9
      },
      itemWidth: 8,
      itemHeight: 8,
      itemGap: 4
    },
    series: [
      {
        name: '账户权重',
        type: 'pie',
        radius: ['35%', '65%'],
        center: ['65%', '50%'],
        data: data,
        emphasis: {
          itemStyle: {
            shadowBlur: 8,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.3)'
          }
        },
        label: {
          show: false
        },
        labelLine: {
          show: false
        }
      }
    ]
  }
  
  pieChart.setOption(option)
}

// 计算连接线的位置
const getConnectionLinePosition = (index) => {
  const accounts = providerForm.value.accounts || []
  let startPosition = 0
  
  for (let i = 0; i < index; i++) {
    startPosition += getAccountWeightPercentage(accounts[i].weight)
  }
  
  // 连接线位置是进度条段的中心位置
  const segmentWidth = getAccountWeightPercentage(accounts[index].weight)
  return startPosition + (segmentWidth / 2)
}

// 确保进度条段正确显示
const getProgressSegmentStyle = (account, index) => {
  return {
    flex: getAccountWeightPercentage(account.weight),
    backgroundColor: getAccountColor(index)
  }
}

// Helper method to explicitly set a feature override field
const setFeatureOverride = (field, value) => {
  providerForm.value.featureOverride[field] = value
}

// 饼图相关方法
const initWeightChart = () => {
  if (!weightChartRef.value) return
  
  weightChart = echarts.init(weightChartRef.value)
  updateWeightChart()
}

const updateWeightChart = () => {
  if (!weightChart) return
  
  const accounts = providerForm.value.accounts || []
  
  if (accounts.length === 0) {
    const option = {
      title: {
        text: '暂无账户配置',
        left: 'center',
        top: 'middle',
        textStyle: {
          color: '#909399',
          fontSize: 12
        }
      },
      graphic: {
        type: 'text',
        left: 'center',
        top: 'center',
        style: {
          text: '请先添加账户配置',
          fontSize: 12,
          fill: '#C0C4CC'
        }
      }
    }
    weightChart.setOption(option)
    return
  }
  
  const data = accounts.map(account => ({
    name: account.name || '未命名账户',
    value: account.weight || 0
  }))
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    series: [{
      name: '账户权重',
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['50%', '50%'],
      data: data,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      label: {
        fontSize: 10,
        formatter: '{b}\n{d}%'
      },
      labelLine: {
        length: 10,
        length2: 5
      }
    }]
  }
  
  weightChart.setOption(option)
}

const resizeChart = () => {
  if (weightChart) {
    weightChart.resize()
  }
}

// 监听账户配置变化，更新饼图
watch(
  () => providerForm.value.accounts,
  () => {
    nextTick(() => {
      updateWeightChart()
    })
  },
  { deep: true }
)

// 监听对话框打开状态，初始化饼图
watch(
  () => providerFormVisible.value,
  (newVal) => {
    if (newVal) {
      nextTick(() => {
        initWeightChart()
      })
    } else {
      // 销毁图表实例
      if (weightChart) {
        weightChart.dispose()
        weightChart = null
      }
    }
  }
)

// Pagination handlers
const handleSizeChange = (newSize) => {
  pageSize.value = newSize
  localStorage.setItem('fusionModelsPageSize', newSize.toString())
  currentPage.value = 1
}

const handleCurrentChange = (newPage) => {
  currentPage.value = newPage
}

// Watch for backend changes
watch(() => fusionStore.selectedBackend, () => {
  loadModels()
})

// Watch for search query changes to reset pagination
watch(searchQuery, () => {
  currentPage.value = 1
})

// Watch for localStorage changes (in case user changes pageSize in another tab)
const checkLocalStorage = () => {
  const storedValue = localStorage.getItem('fusionModelsPageSize')
  if (storedValue && parseInt(storedValue) !== pageSize.value) {
    pageSize.value = parseInt(storedValue)
  }
}

// Check localStorage periodically and on window focus
setInterval(checkLocalStorage, 1000)
window.addEventListener('focus', checkLocalStorage)

// Lifecycle
onMounted(() => {
  loadModels()
  // 添加窗口大小变化监听器
  window.addEventListener('resize', resizeChart)
})

// Cleanup on unmount
onUnmounted(() => {
  clearInterval(checkLocalStorage)
  window.removeEventListener('focus', checkLocalStorage)
  window.removeEventListener('resize', resizeChart)
  // 清理图表实例
  if (weightChart) {
    weightChart.dispose()
    weightChart = null
  }
})
</script>

<style scoped>
.fusion-models {
  padding: 20px;
}

/* 表格列宽优化 */
.fusion-models .el-table .col-3 {
  width: 25%;
}

.fusion-models .el-table .col-4 {
  width: 33.33%;
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

.model-name {
  font-weight: 500;
  color: #303133;
}

.features {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.provider-count {
  display: flex;
  align-items: center;
  gap: 5px;
}

.primary-count {
  font-weight: 500;
  color: #303133;
}

.fallback-count {
  font-weight: 500;
  color: #303133;
}

.provider-management {
  padding: 20px;
}

.provider-section {
  margin-bottom: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h4 {
  margin: 0;
  color: #303133;
}

/* Provider Dialog Styles */
.provider-dialog {
  --el-dialog-padding-primary: 0;
  font-size: 13px;
}

.provider-form {
  padding: 0;
  font-size: 13px;
}

.provider-form .el-form-item {
  margin-bottom: 8px;
}

.provider-form .el-form-item .el-form-item__label {
  padding: 0 !important;
  margin: 0 !important;
  text-align: left !important;
}

.provider-form .el-form-item .el-form-item__content {
  margin: 0 !important;
}

.form-section {
  margin-bottom: 6px;
  padding: 6px;
  background: #fafbfc;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-title {
  margin-bottom: 4px;
}

.section-title h3 {
  margin: 0 0 6px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.section-divider {
  height: 1px;
  background: linear-gradient(90deg, #409eff 0%, #e4e7ed 100%);
  border: none;
}

.section-content {
  background: white;
  padding: 4px;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

/* Basic Info Compact Layout */
.basic-info-compact {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.basic-info-row-compact {
  display: flex;
  gap: 4px;
  align-items: flex-start;
}

.model-name-row {
  margin-top: 2px;
}

.provider-name-compact {
  flex: 1.5;
}

.weight-compact {
  flex: 1.2;
  min-width: 180px;
}

/* 确保提供商名称和提供商模型名左对齐 */
.provider-name-compact .el-form-item__label,
.model-name-compact .el-form-item__label {
  text-align: left !important;
  justify-content: flex-start !important;
  width: 80px !important;
  min-width: 80px !important;
  max-width: 80px !important;
  flex-shrink: 0 !important;
  padding-left: 0 !important;
  margin-left: 0 !important;
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow-wrap: normal !important;
  hyphens: none !important;
  display: inline-block !important;
}

/* 强制所有标签左对齐 */
.basic-info-row-compact .el-form-item {
  display: flex !important;
  align-items: flex-start !important;
}

/* 特别针对提供商模型名标签，强制不换行 */
.provider-dialog .model-name-compact .el-form-item__label,
.provider-dialog .basic-info-row-compact .model-name-compact .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow-wrap: normal !important;
  hyphens: none !important;
  display: inline-block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  max-width: 80px !important;
  min-width: 80px !important;
  width: 80px !important;
  flex-shrink: 0 !important;
  text-align: left !important;
  padding: 0 6px 0 0 !important;
  margin: 0 !important;
  line-height: 28px;
  font-size: 12px;
  font-weight: 500;
}

.basic-info-row-compact .el-form-item .el-form-item__label {
  text-align: left !important;
  padding: 0 6px 0 0 !important;
  margin: 0 !important;
  line-height: 28px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap !important;
  flex-shrink: 0 !important;
  width: 80px !important;
  min-width: 80px !important;
  max-width: 80px !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow-wrap: normal !important;
  hyphens: none !important;
  display: inline-block !important;
}

.weight-compact .el-input-number {
  width: 100%;
}

.weight-compact .el-input-number .el-input__inner {
  text-align: right;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.model-name-compact {
  flex: 1;
}

.model-name-compact .el-form-item__label {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}

.model-name-row .el-form-item__label {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}

/* 强制所有表单标签不换行 */
.provider-form .el-form-item__label {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  min-width: 0 !important;
  flex-shrink: 0 !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
}

/* 特别针对提供商模型名标签 */
.basic-info-compact .el-form-item__label {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  min-width: 0 !important;
  flex-shrink: 0 !important;
  display: inline-block !important;
  max-width: 120px !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
}

/* 更具体的提供商模型名标签样式，确保不换行 */
.provider-form .basic-info-compact .model-name-compact .el-form-item__label {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  display: inline-block !important;
  max-width: 120px !important;
  min-width: 0 !important;
  flex-shrink: 0 !important;
}

/* 全局覆盖Element Plus的表单标签样式，确保不换行 */
.provider-dialog .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}

/* 特别针对提供商模型名标签的额外保护 */
.provider-dialog .el-form-item[prop="modelName"] .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  display: inline-block !important;
  max-width: 120px !important;
}

/* 基于label文本内容的选择器，确保提供商模型名标签不换行 */
.provider-dialog .el-form-item__label:contains("提供商模型名"),
.provider-dialog .el-form-item__label[for*="modelName"] {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  display: inline-block !important;
  max-width: 120px !important;
  min-width: 0 !important;
  flex-shrink: 0 !important;
}

/* 使用更具体的选择器来确保提供商模型名标签不换行 */
.provider-dialog .model-name-compact .el-form-item__label,
.provider-dialog .basic-info-row-compact .model-name-compact .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  display: inline-block !important;
  max-width: 120px !important;
  min-width: 0 !important;
  flex-shrink: 0 !important;
  line-height: 28px !important;
  font-size: 12px !important;
  font-weight: 500 !important;
}

/* 最终的保护措施：使用!important覆盖所有可能的换行样式 */
.provider-dialog .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow-wrap: normal !important;
  hyphens: none !important;
}

/* 特别针对所有表单标签，强制不换行 */
.provider-dialog .el-form-item .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow-wrap: normal !important;
  hyphens: none !important;
  display: inline-block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}

/* 针对Element Plus的特定样式覆盖 */
.provider-dialog .el-form-item {
  display: flex !important;
  flex-wrap: nowrap !important;
}

.provider-dialog .el-form-item .el-form-item__label {
  flex-shrink: 0 !important;
  min-width: 0 !important;
  max-width: 120px !important;
}

/* 特别针对提供商模型名标签的最终保护 */
.provider-dialog .el-form-item[prop="modelName"] .el-form-item__label,
.provider-dialog .model-name-compact .el-form-item__label {
  white-space: nowrap !important;
  word-break: keep-all !important;
  word-wrap: normal !important;
  overflow-wrap: normal !important;
  hyphens: none !important;
  display: inline-block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  max-width: 80px !important;
  min-width: 80px !important;
  width: 80px !important;
  flex-shrink: 0 !important;
  text-align: left !important;
  padding: 0 6px 0 0 !important;
  margin: 0 !important;
  line-height: 28px;
  font-size: 12px;
  font-weight: 500;
  position: relative !important;
}

.basic-info-row-compact .el-form-item {
  margin-bottom: 0;
}

.basic-info-row-compact .el-form-item .el-form-item__label {
  text-align: left !important;
  padding: 0 6px 0 0 !important;
  margin: 0 !important;
  line-height: 28px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap !important;
}

.basic-info-row-compact .el-form-item .el-form-item__content {
  margin: 0 !important;
}

.model-name-input-compact {
  width: 100%;
}

/* Weight Distribution Bar */
.weight-distribution-bar {
  margin-top: 4px;
  padding: 8px;
  background: #f8f9fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.weight-bar-title {
  font-size: 12px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 4px;
}

.weight-bar-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.shared-weight-bar {
  display: flex;
  flex-direction: column;
  gap: 6px;
  position: relative;
}

/* 连线样式 */
.shared-weight-bar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 1;
}

/* 移除连接线相关样式 */

/* 移除旧的连线样式 */

.account-labels {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 2px;
}

.account-label-item {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 11px;
  padding: 2px 4px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 3px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
}

.account-label-item:hover {
  background: rgba(255, 255, 255, 1);
  border-color: rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.account-name-bar {
  font-weight: 600;
  color: #303133;
  font-size: 10px;
}

.account-weight-bar {
  color: #409eff;
  font-weight: 600;
  font-size: 10px;
}

.shared-progress-bar {
  position: relative;
  height: 16px;
  background: #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e4e7ed;
  margin-top: 6px;
  display: flex;
}

.progress-segment {
  flex: 1;
  height: 100%;
  transition: all 0.3s ease;
  border-right: 1px solid rgba(255, 255, 255, 0.3);
  cursor: pointer;
  z-index: 1;
}

.progress-segment:last-child {
  border-right: none;
}

.progress-segment:hover {
  filter: brightness(1.2);
  transform: scaleY(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.progress-segment:last-child {
  border-right: none;
}

.no-accounts-bar {
  text-align: center;
  color: #909399;
  font-size: 11px;
  padding: 8px;
  font-style: italic;
}

/* Feature Override Container */
.feature-override-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.feature-row {
  background: #fafbfc;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 8px 12px;
  transition: all 0.3s ease;
  margin-bottom: 6px;
}

.feature-row:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.feature-row-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.feature-label {
  font-weight: 500;
  color: #606266;
  font-size: 13px;
  min-width: 100px;
  flex-shrink: 0;
  line-height: 28px;
}

.feature-options {
  flex: 1;
  display: flex;
  flex-direction: row;
  gap: 2px;
  align-items: center;
}

.feature-options .el-radio-group {
  display: flex;
  flex-wrap: nowrap;
  gap: 12px;
  align-items: center;
}

.token-options {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
}

.token-radio-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
  white-space: nowrap;
}

.token-radio-group .el-radio {
  margin-right: 8px;
  white-space: nowrap;
}

.token-input-inline {
  min-width: 100px;
  max-width: 120px;
  flex-shrink: 0;
  margin-left: 8px;
}

.token-input {
  max-width: 200px;
}

/* Special Features Section */
.special-features-section {
  padding: 0;
  background: transparent;
  border: none;
}

.special-features-description {
  color: #606266;
  font-size: 12px;
  margin-bottom: 12px;
  line-height: 1.4;
  padding: 8px;
  background: #f0f9ff;
  border-left: 3px solid #409eff;
  border-radius: 4px;
}

.add-feature-section {
  margin-bottom: 8px;
}

.feature-option {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.feature-name {
  font-weight: 500;
  color: #303133;
}

.feature-desc {
  font-size: 12px;
  color: #909399;
}

.special-feature-item {
  margin-bottom: 6px;
  padding: 8px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background-color: #ffffff;
  transition: all 0.3s ease;
  position: relative;
}

.special-feature-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.special-feature-item:last-child {
  margin-bottom: 0;
}

.special-feature-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.special-feature-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.special-feature-label {
  font-weight: 500;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  flex-shrink: 0;
}

.special-feature-input-wrapper {
  min-width: 0;
  margin-right: 32px;
  transition: all 0.3s ease;
}

.special-feature-close-btn {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 20px;
  height: 20px;
  padding: 0;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-color: #f56c6c;
  color: #f56c6c;
  background-color: #ffffff;
  z-index: 2;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.special-feature-close-btn:hover {
  background-color: #f56c6c;
  color: white;
  border-color: #f56c6c;
  transform: scale(1.1);
}

.special-feature-input {
  width: 100%;
}

.special-feature-switch {
  width: 100%;
}

.special-feature-input .el-input-number {
  width: 100%;
}

.empty-features {
  text-align: center;
  padding: 12px;
  color: #909399;
  font-size: 12px;
  background-color: #fafbfc;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
  margin-top: 8px;
}

.empty-features span {
  margin: 0;
}

/* Account Configuration */
.account-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.account-item {
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 8px 12px 8px 8px;
  background-color: #ffffff;
  transition: all 0.3s ease;
  margin-bottom: 8px;
  position: relative;
}

/* 供应商总览样式 */
.provider-overview {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.provider-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.provider-count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 6px;
  background-color: #409eff;
  color: white;
  border-radius: 3px;
  font-size: 9px;
  font-weight: 600;
  flex-shrink: 0;
  margin-right: 4px;
}

/* 供应商详情弹出框样式 */
.provider-tooltip {
  position: fixed;
  z-index: 9999;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 320px;
  max-height: 450px;
  overflow: hidden;
}

.tooltip-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.tooltip-header h4 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.tooltip-close {
  cursor: pointer;
  font-size: 18px;
  color: #909399;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.tooltip-close:hover {
  background-color: #f56c6c;
  color: white;
}

.tooltip-content {
  padding: 16px;
  max-height: 370px;
  overflow-y: auto;
}

.accounts-section,
.chart-section {
  margin-bottom: 12px;
}

.accounts-section h5,
.chart-section h5 {
  margin: 0 0 6px 0;
  font-size: 11px;
  font-weight: 600;
  color: #606266;
}

.accounts-list {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.account-item-tooltip {
  padding: 4px 6px;
  background: #fafbfc;
  border-radius: 3px;
  border: 1px solid #e4e7ed;
}

.account-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  min-height: 20px;
}

.account-name {
  font-weight: 500;
  color: #303133;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.account-weight {
  color: #409eff;
  font-weight: 600;
  min-width: 35px;
  text-align: right;
  flex-shrink: 0;
}

.account-region {
  color: #909399;
  font-size: 10px;
  background: #f0f2f5;
  padding: 2px 4px;
  border-radius: 2px;
  flex-shrink: 0;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-accounts,
.no-chart {
  text-align: center;
  color: #c0c4cc;
  font-size: 10px;
  padding: 12px;
  background: #fafbfc;
  border-radius: 3px;
  border: 1px dashed #d9d9d9;
}

.chart-container {
  height: 180px;
  position: relative;
  margin-top: 8px;
}

.pie-chart {
  width: 100%;
  height: 100%;
}

.provider-label {
  color: #909399;
  font-size: 11px;
}

.provider-progress-container {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.provider-progress-bar {
  display: flex;
  height: 6px;
  border-radius: 3px;
  overflow: hidden;
  background-color: #f5f7fa;
}

.provider-progress-segment {
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
}

.provider-progress-segment:hover {
  filter: brightness(1.2);
  transform: scaleY(1.1);
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.3);
}

.provider-progress-segment::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, transparent 30%, rgba(255, 255, 255, 0.1) 50%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.provider-progress-segment:hover::after {
  opacity: 1;
}

.provider-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 3px;
  margin-top: 1px;
  align-items: center;
}

.provider-label-item {
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 9px;
  padding: 2px 4px;
  background-color: #f5f7fa;
  border-radius: 3px;
  border: 1px solid #e4e7ed;
  transition: all 0.2s ease;
  cursor: pointer;
  flex-shrink: 0;
}

.provider-label-item:hover,
.provider-label-item.highlighted {
  background-color: #ecf5ff;
  border-color: #409eff;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.provider-label-item.highlighted {
  background-color: #e6f7ff;
  border-color: #1890ff;
  box-shadow: 0 2px 6px rgba(24, 144, 255, 0.2);
}

.provider-color-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* 高级配置样式 */
.advanced-config-section {
  margin-top: 20px;
  padding: 16px;
  background: #fafbfc;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.advanced-config-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.advanced-config-header h4 {
  margin: 0;
  color: #303133;
  font-size: 14px;
  font-weight: 600;
}

.advanced-config-toggle {
  display: flex;
  align-items: center;
}

.advanced-config-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.advanced-config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.advanced-config-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.advanced-config-label {
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  flex-shrink: 0;
}

.advanced-config-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.advanced-config-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  min-height: 32px;
  padding: 8px;
  background: white;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.config-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 4px;
  font-size: 12px;
  color: #0369a1;
}

.config-item-text {
  font-weight: 500;
}

.config-item-close-btn {
  padding: 2px;
  color: #f56c6c;
  transition: all 0.2s ease;
}

.config-item-close-btn:hover {
  color: #f56c6c;
  background-color: #fef0f0;
}

.empty-config {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
  font-size: 12px;
  font-style: italic;
  width: 100%;
  height: 32px;
}

/* 字符集编码显示样式 */
.encoding-info {
  font-size: 12px;
  color: #606266;
  padding: 2px 6px;
  background-color: #f5f7fa;
  border-radius: 3px;
  border: 1px solid #e4e7ed;
}

/* 模型对话框样式优化 */
.model-dialog .el-dialog__body {
  padding: 16px 20px;
}

.model-form .el-form-item {
  margin-bottom: 16px;
}

/* Compact Feature Grid */
.feature-grid-compact {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
  background-color: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.feature-row-compact {
  display: flex;
  gap: 24px;
  align-items: center;
}

.feature-item-compact {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.feature-label-compact {
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
  min-width: 80px;
}

.feature-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.feature-row {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.feature-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.feature-label {
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 4px;
}

.feature-label {
  font-weight: 500;
  color: #606266;
}

.timestamp-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.timestamp-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.timestamp-label {
  font-size: 13px;
  color: #909399;
  min-width: 70px;
}

.timestamp-value {
  font-size: 13px;
  color: #606266;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background-color: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  border: 1px solid #e4e7ed;
}

/* Compact Timestamp */
.timestamp-grid-compact {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px 12px;
  background-color: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.timestamp-item-compact {
  display: flex;
  align-items: center;
  gap: 8px;
}

.timestamp-label-compact {
  font-size: 12px;
  color: #909399;
  min-width: 60px;
}

.timestamp-value-compact {
  font-size: 12px;
  color: #606266;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Tags样式 */
.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  min-height: 24px;
}

.model-tag {
  margin: 0;
}

.no-tags {
  color: #c0c4cc;
  font-size: 12px;
}

.tags-edit-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tags-input-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

/* Compact Tags */
.tags-edit-container-compact {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tags-input-row-compact {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.tags-display-compact {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  min-height: 24px;
  align-items: center;
}

.editable-tag-compact {
  margin: 0;
}

.tags-display {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  min-height: 32px;
  padding: 8px;
  background-color: #fafbfc;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.editable-tag {
  margin: 0;
}

.no-tags-hint {
  color: #c0c4cc;
  font-size: 12px;
}

.provider-name {
  color: #606266;
  font-weight: 500;
  white-space: nowrap;
  flex: 0 1 auto;
  min-width: 0;
}

.provider-weight {
  color: #909399;
  font-size: 8px;
  font-weight: 600;
  flex-shrink: 0;
  margin-left: 2px;
}



/* 账户删除按钮样式 */
.account-close-btn {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 20px;
  height: 20px;
  padding: 0;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-color: #f56c6c;
  color: #f56c6c;
  background-color: #ffffff;
  z-index: 2;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.account-close-btn:hover {
  background-color: #f56c6c;
  color: white;
  border-color: #f56c6c;
  transform: scale(1.1);
}

.account-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.account-grid {
  display: flex;
  gap: 8px;
  align-items: flex-end;
  flex-wrap: nowrap;
}

.account-select {
  flex: 1.3;
  min-width: 0;
}

.account-weight {
  flex: 0.7;
  min-width: 110px;
}

.account-region {
  flex: 0.9;
  min-width: 120px;
}

.account-select,
.account-weight,
.account-region {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-select label,
.account-weight label,
.account-region label {
  font-size: 12px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 2px;
  white-space: nowrap;
}

/* 特别优化权重标签 */
.account-weight label {
  font-size: 12px;
  margin-bottom: 2px;
}

.account-actions {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-top: 20px;
}

.selected-account-description {
  font-size: 11px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
  padding: 4px 8px;
  background: #f8f9fa;
  border-radius: 4px;
  border-left: 3px solid #409eff;
}

/* 账户名称选择框样式 */
.account-name-select {
  width: 100%;
}

.account-name-select .el-input__wrapper {
  padding: 0 8px;
}

/* 权重输入框样式 */
.account-weight-input {
  width: 100%;
}

.account-weight-input .el-input-number__decrease,
.account-weight-input .el-input-number__increase {
  width: 19px;
  height: 19px;
}

.account-weight-input .el-input__inner {
  text-align: left;
  font-size: 12px;
}

/* 区域选择框样式 */
.account-region-select {
  width: 100%;
}

.account-option {
  display: flex;
  align-items: center;
  gap: 4px;
}

.account-name {
  font-weight: 500;
  color: #303133;
  font-size: 12px;
}

.account-description {
  color: #909399;
  font-size: 11px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 150px;
}

.add-account-section {
  display: flex;
  justify-content: center;
  padding: 20px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  background-color: #fafbfc;
  transition: all 0.3s ease;
}

.add-account-section:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
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

/* Responsive Design */
@media (max-width: 768px) {
  .basic-info-row-compact {
    flex-direction: column;
    gap: 4px;
  }
  
  .provider-name-compact,
  .weight-compact,
  .model-name-compact {
    flex: 1;
    min-width: auto;
  }
  
  .weight-distribution-bar {
    margin-top: 4px;
  }
  
  .weight-bar-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .weight-bar-info {
    min-width: auto;
    width: 100%;
  }
  
  .account-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .account-actions {
    justify-content: flex-start;
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>