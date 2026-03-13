<template>
  <div class="page-container">
    <div class="header">
      <div class="page-header">
        <div class="title-section">
          <i class="el-icon-cpu title-icon"></i>
          <h1 class="page-title">资源管理</h1>
        </div>
        <div class="header-divider"></div>
      </div>
    </div>

    <div class="cluster-tabs">
      <el-tabs 
        v-model="activeTab" 
        type="card" 
        @tab-click="handleTabClick"
        class="custom-tabs"
      >
        <el-tab-pane 
          v-for="cluster in clusters" 
          :key="cluster.id"
          :label="cluster.name"
          :name="cluster.id"
        >
          <div class="content-card">
            <div class="table-header">
              <div class="table-filters">
                <el-input
                  v-model="searchQuery"
                  placeholder="搜索节点 ID、标签、Provider、或已预热模型"
                  prefix-icon="el-icon-search"
                  clearable
                  class="search-input"
                />
                <el-select
                  v-model="selectedGpuType"
                  placeholder="GPU型号"
                  clearable
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  class="gpu-filter"
                >
                  <el-option
                    v-for="type in gpuTypes"
                    :key="type"
                    :label="type"
                    :value="type"
                  />
                </el-select>
                <el-select
                  v-model="selectedProvider"
                  placeholder="Provider"
                  clearable
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  class="provider-filter"
                >
                  <el-option
                    v-for="provider in providers"
                    :key="provider"
                    :label="provider"
                    :value="provider"
                  />
                </el-select>
                <el-checkbox v-model="showSchedulableOnly">
                  仅显示可调度节点
                </el-checkbox>
                <el-checkbox v-model="showAvailableGpuOnly">
                  仅看空闲卡
                </el-checkbox>
                <div class="gpu-total-info">
                  空闲 GPU: <span class="gpu-total-value">{{ computeTotalFreeGPUs }}</span>
                </div>
                <div class="record-info">
                  共 {{ filteredNodes.length }} 条记录
                  <template v-if="selectedNodes.length > 0">
                    （已选中 {{ selectedNodes.length }} 条）
                  </template>
                </div>
                <el-button 
                  type="primary" 
                  plain
                  size="small"
                  @click="columnSettingsVisible = true"
                >
                  <i class="el-icon-setting"></i>
                  列设置
                </el-button>
              </div>
              <div class="table-actions">
                <el-button 
                  v-if="selectedNodes.length > 0"
                  type="primary"
                  size="small"
                  @click="handleShowBatchLabelDialog"
                >
                  批量设置标签
                </el-button>
                <el-button 
                  v-if="selectedNodes.length > 0"
                  type="primary"
                  size="small"
                  @click="showBatchModelWarmupDialog"
                >
                  批量预热模型
                </el-button>
                <el-button type="primary" size="small" @click="refreshCurrentTab">
                  <i class="el-icon-refresh"></i> 刷新
                </el-button>
              </div>
            </div>
            <el-table 
              v-loading="loading[cluster.id]" 
              :data="filteredNodes" 
              :default-sort="{ prop: 'instanceCount', order: 'descending' }"
              style="width: 100%"
              :row-style="{
                height: '60px'
              }"
              :header-cell-style="{
                background: '#f5f7fa',
                color: '#333',
                fontWeight: 600,
                fontSize: '15px',
                padding: '16px 12px',
                height: '60px',
                textAlign: 'center'
              }"
              :cell-style="{
                color: '#333',
                fontSize: '15px',
                padding: '16px 12px'
              }"
              border
              fit
              @selection-change="handleSelectionChange"
            >
              <el-table-column type="selection" width="55" fixed="left" />
              <el-table-column prop="id" label="节点 ID" :width="columnWidths.id" resizable fixed="left">
                <template #default="scope">
                  <div class="node-info">
                    <span class="node-id clickable" @click="showNodeDetail(scope.row)">
                      {{ scope.row.id }}
                    </span>
                    <el-tooltip
                      v-if="!scope.row.schedulable"
                      content="节点不可调度"
                      placement="top"
                    >
                      <span class="status-icon unschedulable">🚫</span>
                    </el-tooltip>
                  </div>
                </template>
              </el-table-column>
              <el-table-column v-if="columnSettings.gpuModel" label="GPU型号" :width="columnWidths.gpuModel" resizable>
                <template #default="scope">
                  <span>{{ scope.row.spec.gpuProductName || '-' }}</span>
                </template>
              </el-table-column>
              <el-table-column v-if="columnSettings.cpuModel" label="CPU型号" :width="columnWidths.cpuModel" resizable>
                <template #default="scope">
                  <span>{{ scope.row.spec.cpuProductName }}</span>
                </template>
              </el-table-column>
              <el-table-column 
                v-if="columnSettings.provider" 
                label="Provider" 
                :width="columnWidths.provider" 
                resizable
                align="center"
              >
                <template #default="scope">
                  <span>{{ scope.row.spec.provider }}</span>
                </template>
              </el-table-column>
              <el-table-column v-if="columnSettings.cuda" label="CUDA" :width="columnWidths.cuda" align="center" resizable>
                <template #default="scope">
                  <span>{{ scope.row.spec.cudaVersion }}</span>
                </template>
              </el-table-column>
              <el-table-column v-if="columnSettings.cpuUsage" label="CPU使用" :width="columnWidths.cpuUsage" align="center" resizable>
                <template #default="scope">
                  <span class="resource-value">
                    {{ computeCpuUsage(scope.row) }}/{{ scope.row.capacity.cpuNum }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column 
                v-if="columnSettings.gpuUsage" 
                label="GPU使用" 
                :width="columnWidths.gpuUsage" 
                align="center" 
                resizable
                sortable
                :sort-method="(a, b) => {
                  const aUsage = computeGpuUsage(a)
                  const bUsage = computeGpuUsage(b)
                  return aUsage - bUsage
                }"
              >
                <template #default="scope">
                  <span class="resource-value">
                    {{ scope.row.capacity.gpuNum ? `${computeGpuUsage(scope.row)}/${scope.row.capacity.gpuNum}` : '-' }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column v-if="columnSettings.memory" label="内存(GB)" :width="columnWidths.memory" align="center" resizable>
                <template #default="scope">
                  {{ formatGb(scope.row.capacity.memorySize - scope.row.allocable.memorySize) }}/{{ formatGb(scope.row.capacity.memorySize) }}
                </template>
              </el-table-column>
              <el-table-column 
                v-if="columnSettings.instanceCount" 
                label="实例数" 
                :width="columnWidths.instanceCount"
                align="center"
                resizable
                prop="instanceCount"
                sortable
              >
                <template #default="scope">
                  <span class="instance-count">
                    {{ scope.row.state?.instIds?.length || 0 }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column 
                v-if="columnSettings.labels" 
                label="标签" 
                :width="columnWidths.labels"
                resizable
              >
                <template #default="scope">
                  <div class="label-container">
                    <el-tag
                      v-for="(value, key) in scope.row.label"
                      :key="key"
                      :style="{
                        backgroundColor: generateColorFromHash(`${key}=${value}`),
                        color: generateTextColorFromHash(`${key}=${value}`),
                        border: 'none'
                      }"
                      size="small"
                      class="label-tag"
                    >
                      {{ key }}={{ value }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="磁盘用量" width="200" align="center" fixed="right" v-if="columnSettings.disks"
                sortable
                :sort-method="(a, b) => {
                  // 获取第一个磁盘的使用率
                  const getDiskUsage = (node) => {
                    if (!node.disks || node.disks.length === 0) return 0
                    const disk = node.disks[0]
                    return (1 - disk.freeSpace/disk.totalSpace) * 100
                  }
                  return getDiskUsage(a) - getDiskUsage(b)
                }"
              >
                <template #default="scope">
                  <div v-if="scope.row.disks" class="disk-usage-container">
                    <div v-for="disk in scope.row.disks" :key="disk.diskID" class="disk-item">
                      <div class="disk-info" style="display: flex; justify-content: space-between;">
                        <span class="disk-name" style="text-align: left;">{{ disk.diskID }}:</span>
                        <span class="disk-numbers" style="text-align: right;">
                          {{ formatBytes(disk.totalSpace - disk.freeSpace) }}/{{ formatBytes(disk.totalSpace) }}
                        </span>
                      </div>
                      <el-progress 
                        :percentage="Math.round((1 - disk.freeSpace/disk.totalSpace) * 100)"
                        :stroke-width="12"
                        :show-text="false"
                        :status="(1 - disk.freeSpace/disk.totalSpace) > 0.9 ? 'exception' : ''"
                      />
                    </div>
                  </div>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column 
                label="模型列表" 
                width="400" 
                align="left" 
                fixed="right"
                v-if="columnSettings.models"
              >
                <template #default="scope">
                  <div v-if="scope.row.models && scope.row.models.length > 0" class="model-list">
                    <div 
                      v-for="model in sortedModels(scope.row.models)"
                      :key="model.modelName" 
                      class="model-item"
                    >
                      <span class="model-name">{{ model.modelName }}</span>
                      <el-tag 
                        :type="model.modelStatus === 'ready' ? 'success' : 'warning'"
                        size="small"
                      >
                        {{ model.modelStatus }}
                      </el-tag>
                    </div>
                  </div>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="260" align="center" fixed="right">
                <template #default="scope">
                  <div class="action-buttons" style="display: flex; gap: 2px;">
                    <el-button
                      size="default"
                      :type="scope.row.schedulable ? 'warning' : 'success'"
                      @click="handleScheduleChange(scope.row)"
                    >
                      {{ scope.row.schedulable ? '停止调度' : '打开调度' }}
                    </el-button>
                    <el-button
                      size="default"
                      type="primary"
                      @click="handleEditLabels(scope.row)"
                    >
                      设置Label
                    </el-button>
                    <el-button
                      size="default"
                      type="info"
                      @click="showModelWarmupDialog(scope.row)"
                    >
                      模型预热
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- Label编辑对话框 -->
    <el-dialog
      v-model="labelDialogVisible"
      title="编辑标签"
      width="800px"
    >
      <div class="label-edit-container">
        <div v-for="(value, key) in editingLabels" :key="key" class="label-item">
          <el-input 
            v-model="editingLabels[key]" 
            class="label-input"
            :disabled="value === '-'"
          >
            <template #prepend>
              <div class="label-key">{{ key }}</div>
            </template>
          </el-input>
          <el-button 
            type="danger" 
            @click="removeLabel(key)"
            v-if="value !== '-'"
          >
            删除
          </el-button>
          <el-button 
            type="success" 
            @click="editingLabels[key] = ''"
            v-else
          >
            恢复
          </el-button>
        </div>
        <div class="add-label">
          <el-input v-model="newLabelKey" placeholder="键" class="label-key-input" />
          <el-input v-model="newLabelValue" placeholder="值" class="label-value-input" />
          <el-button type="primary" @click="addLabel">添加</el-button>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="labelDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveLabels">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 节点详情对话框 -->
    <el-dialog
      v-model="nodeDetailVisible"
      title="节点详情"
      width="80%"
      :fullscreen="false"
      :top="'5vh'"
      destroy-on-close
    >
      <template #title>
        <div class="dialog-title-container">
          <div class="dialog-title">节点详情</div>
          <el-button
            type="primary"
            size="small"
            :loading="detailLoading"
            @click="refreshNodeDetail"
          >
            <i class="el-icon-refresh"></i>
            刷新
          </el-button>
        </div>
      </template>
      <div v-loading="detailLoading" class="node-detail">
        <template v-if="nodeDetail">
          <div class="detail-section">
            <div class="section-header">
              <i class="el-icon-info"></i>
              <h3>基本信息</h3>
            </div>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">节点 ID:</span>
                <span class="value">{{ nodeDetail.id }}</span>
              </div>
              <div class="info-item">
                <span class="label">Provider:</span>
                <span class="value">{{ nodeDetail.spec.provider }}</span>
              </div>
              <div class="info-item">
                <span class="label">调度状态:</span>
                <el-tag :type="nodeDetail.schedulable ? 'success' : 'danger'" size="small">
                  {{ nodeDetail.schedulable ? '可调度' : '不可调度' }}
                </el-tag>
              </div>
            </div>
          </div>
          
          <div class="detail-section">
            <div class="section-header">
              <i class="el-icon-cpu"></i>
              <h3>硬件信息</h3>
            </div>
            <div class="info-grid hardware-info">
              <div class="info-item">
                <span class="label">CPU型号:</span>
                <span class="value">{{ nodeDetail.spec.cpuProductName }}</span>
              </div>
              <div class="info-item">
                <span class="label">GPU型号:</span>
                <span class="value">{{ nodeDetail.spec.gpuProductName || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">CUDA版本:</span>
                <span class="value">{{ nodeDetail.spec.cudaVersion || '-' }}</span>
              </div>
            </div>
          </div>
          
          <div class="detail-section">
            <div class="section-header">
              <i class="el-icon-data-line"></i>
              <h3>资源使用</h3>
            </div>
            <div class="resource-usage-grid">
              <div class="usage-item">
                <div class="progress-label">
                  <div>CPU使用率</div>
                  <div>{{ computeCpuUsage(nodeDetail) }}/{{ nodeDetail?.capacity?.cpuNum }}</div>
                </div>
                <el-progress 
                  type="circle" 
                  :percentage="computeCpuUsagePercentage(nodeDetail)"
                  :color="getProgressColor"
                  :width="80"
                />
              </div>
              <div class="usage-item">
                <div class="progress-label">
                  <div>内存使用率</div>
                  <div>{{ formatGb(nodeDetail?.capacity?.memorySize - nodeDetail?.allocable?.memorySize) }}/{{ formatGb(nodeDetail?.capacity?.memorySize) }}GB</div>
                </div>
                <el-progress 
                  type="circle" 
                  :percentage="computeMemoryUsagePercentage(nodeDetail)"
                  :color="getProgressColor"
                  :width="80"
                />
              </div>
              <div class="usage-item">
                <div class="progress-label">
                  <div>系统盘使用率</div>
                  <div>{{ formatGb(nodeDetail?.capacity?.rootfsSize - nodeDetail?.allocable?.rootfsSize) }}/{{ formatGb(nodeDetail?.capacity?.rootfsSize) }}GB</div>
                </div>
                <el-progress 
                  type="circle" 
                  :percentage="computeRootfsUsagePercentage(nodeDetail)"
                  :color="getProgressColor"
                  :width="80"
                />
              </div>
              <div class="usage-item">
                <div class="progress-label">
                  <div>数据盘使用率</div>
                  <div>{{ formatGb(nodeDetail?.capacity?.localStorageSize - nodeDetail?.allocable?.localStorageSize) }}/{{ formatGb(nodeDetail?.capacity?.localStorageSize) }}GB</div>
                </div>
                <el-progress 
                  type="circle" 
                  :percentage="computeLocalStorageUsagePercentage(nodeDetail)"
                  :color="getProgressColor"
                  :width="80"
                />
              </div>
            </div>
            <div class="divider"></div>
            <div class="gpu-usage">
              <div class="gpu-header">
                <span>GPU使用率</span>
                <span>{{ computeGpuUsage(nodeDetail) }}/{{ nodeDetail?.capacity?.gpuNum }}</span>
              </div>
              <div class="gpu-cards-container">
                <div class="gpu-cards">
                  <div 
                    v-for="(state, index) in nodeDetail?.allocable?.gpuState" 
                    :key="index"
                    class="gpu-card"
                    :class="{ 'gpu-used': state === 1 }"
                  >
                    GPU {{ index }}
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <div class="detail-section">
            <div class="section-header">
              <i class="el-icon-s-operation"></i>
              <h3>运行实例 ({{ nodeDetail.state?.instIds?.length || 0 }})</h3>
            </div>
            <el-table 
              :data="nodeDetail.state?.instIds || []"
              style="width: 100%"
              :max-height="'calc(100vh - 400px)'"
              size="small"
              :header-cell-style="{
                background: '#f5f7fa',
                color: '#333',
                fontWeight: 600,
                padding: '4px 2px'
              }"
              :cell-style="{
                padding: '2px'
              }"
            >
              <el-table-column prop="id" label="实例 ID" min-width="150">
                <template #default="scope">
                  <el-tooltip 
                    :content="scope.row"
                    placement="top"
                    :show-after="1000"
                  >
                    <span 
                      class="instance-id clickable"
                      @click="showInstanceDetail(scope.row)"
                    >
                      {{ scope.row }}
                    </span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="镜像" min-width="500">
                <template #default="scope">
                  <el-tooltip 
                    :content="getInstanceImage(scope.row)"
                    placement="top"
                    :show-after="1000"
                  >
                    <span class="image-name">{{ getInstanceImage(scope.row) }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="Model" min-width="260">
                <template #default="scope">
                  <el-tooltip 
                    :content="getModelId(scope.row)"
                    placement="top"
                    :show-after="1000"
                  >
                    <span class="model-id">{{ getModelId(scope.row) }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="Worker ID" min-width="500" :show-overflow-tooltip="true">
                <template #default="scope">
                  <span class="worker-id">{{ getWorkerID(scope.row) }}</span>
                </template>
              </el-table-column>
              <el-table-column 
                label="创建时间" 
                min-width="150"
                :formatter="(row) => getInstanceCreateTime(row)"
              />
              <el-table-column label="操作" width="70" fixed="right" align="center">
                <template #default="scope">
                  <el-button
                    v-if="getWorkerID(scope.row) !== '-'"
                    type="danger"
                    size="small"
                    @click="handleEvictInstance(scope.row)"
                  >
                    驱逐
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </template>
      </div>
    </el-dialog>

    <!-- 添加列设置弹窗 -->
    <el-dialog
      v-model="columnSettingsVisible"
      title="列显示设置"
      width="400px"
    >
      <div class="column-settings">
        <el-checkbox v-model="columnSettings.gpuModel">GPU型号</el-checkbox>
        <el-checkbox v-model="columnSettings.cpuModel">CPU型号</el-checkbox>
        <el-checkbox v-model="columnSettings.provider">Provider</el-checkbox>
        <el-checkbox v-model="columnSettings.cuda">CUDA</el-checkbox>
        <el-checkbox v-model="columnSettings.cpuUsage">CPU使用</el-checkbox>
        <el-checkbox v-model="columnSettings.gpuUsage">GPU使用</el-checkbox>
        <el-checkbox v-model="columnSettings.memory">内存</el-checkbox>
        <el-checkbox v-model="columnSettings.instanceCount">实例数</el-checkbox>
        <el-checkbox v-model="columnSettings.labels">标签</el-checkbox>
        <el-checkbox v-model="columnSettings.disks">磁盘</el-checkbox>
        <el-checkbox v-model="columnSettings.models">模型</el-checkbox>
      </div>
    </el-dialog>

    <!-- 实例详情弹窗 -->
    <el-dialog
      v-model="instanceDetailVisible"
      title="实例详情"
      width="1200px"
      destroy-on-close
      class="instance-details-dialog"
    >
      <div class="dialog-header-actions">
        <el-button 
          type="info" 
          size="small"
          @click="showRawInstanceDetail"
        >
          查看原始信息
        </el-button>
      </div>
      <div v-loading="instanceDetailLoading" class="instance-detail">
        <template v-if="selectedInstance">
          <div class="detail-section">
            <div class="section-header">
              <div class="header-left">
                <i class="el-icon-info"></i>
                <h3>基本信息</h3>
              </div>
            </div>
            <div class="info-grid">
              <div class="info-item double">
                <div class="info-pair">
                  <span class="label">实例 ID:</span>
                  <span class="value">{{ selectedInstance.id }}</span>
                </div>
                <div class="info-pair">
                  <span class="label">节点 ID:</span>
                  <span class="value">{{ selectedInstance.nodeId }}</span>
                </div>
              </div>
              <div class="info-item double">
                <div class="info-pair">
                  <span class="label">状态:</span>
                  <el-tag :type="selectedInstance.state?.state === 'running' ? 'success' : 'warning'">
                    {{ selectedInstance.state?.state || '-' }}
                  </el-tag>
                </div>
                <div class="info-pair">
                  <span class="label">实例类型:</span>
                  <span class="value">{{ selectedInstance.kind || '-' }}</span>
                </div>
              </div>
              <div class="info-item">
                <span class="label">镜像:</span>
                <span class="value code">{{ selectedInstance.containers?.image || '-' }}</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <div class="section-header">
              <div class="header-left">
                <i class="el-icon-monitor"></i>
                <h3>资源配置</h3>
              </div>
            </div>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">CPU 核心数:</span>
                <span class="value">{{ selectedInstance.containers?.cpuNum || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">内存:</span>
                <span class="value">{{ formatGb(selectedInstance.containers?.memorySize) }}GB</span>
              </div>
              <div class="info-item">
                <span class="label">GPU 数量:</span>
                <span class="value">{{ selectedInstance.containers?.gpuNum || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">GPU 型号:</span>
                <span class="value">{{ selectedInstance.containers?.gpuProductName || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">系统盘:</span>
                <span class="value">{{ formatGb(selectedInstance.containers?.rootfsSize) }}GB</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <div class="section-header">
              <div class="header-left">
                <i class="el-icon-connection"></i>
                <h3>网络与访问</h3>
              </div>
            </div>
            <div class="info-grid">
              <div v-if="selectedInstance.exposePortStates?.length" class="info-item port-mapping">
                <span class="label">端口映射:</span>
                <div class="value-group">
                  <div v-for="port in selectedInstance.exposePortStates" :key="port.port" class="port-item">
                    <span class="value">
                      {{ port.exposeAddress }}
                      <el-tag size="small" :type="port.online ? 'success' : 'danger'" class="ml-2">
                        {{ port.online ? '在线' : '离线' }}
                      </el-tag>
                    </span>
                  </div>
                </div>
              </div>
              <div v-if="selectedInstance.apps?.sshState" class="info-item ssh-info">
                <div class="ssh-info-row">
                  <div class="ssh-item">
                    <span class="label">SSH 命令:</span>
                    <span class="code-value">{{ selectedInstance.apps.sshState.sshCommand }}</span>
                    <el-button
                      type="primary"
                      size="small"
                      class="copy-btn"
                      @click="copyToClipboard(selectedInstance.apps.sshState.sshCommand)"
                    >
                      复制命令
                    </el-button>
                  </div>
                  <div class="ssh-item">
                    <span class="label">SSH 密码:</span>
                    <span class="code-value">{{ selectedInstance.apps.sshState.sshPassword }}</span>
                    <el-button
                      type="primary"
                      size="small"
                      class="copy-btn"
                      @click="copyToClipboard(selectedInstance.apps.sshState.sshPassword)"
                    >
                      复制密码
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 日志信息 -->
          <div class="detail-section" v-if="selectedInstance.apps?.logState">
            <div class="section-header">
              <div class="header-left">
                <i class="el-icon-document"></i>
                <h3>日志信息</h3>
              </div>
            </div>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">系统日志:</span>
                <div class="value-group">
                  <a 
                    :href="`${selectedInstance.apps.logState.systemLogAddress}?follow=1&tail=100`"
                    target="_blank"
                    class="log-link"
                  >
                    {{ selectedInstance.apps.logState.systemLogAddress }}
                  </a>
                </div>
              </div>
              <div class="info-item">
                <span class="label">实例日志:</span>
                <div class="value-group">
                  <a 
                    :href="`${selectedInstance.apps.logState.instanceLogAddress}?follow=1&tail=100`"
                    target="_blank"
                    class="log-link"
                  >
                    {{ selectedInstance.apps.logState.instanceLogAddress }}
                  </a>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </el-dialog>

    <!-- 实例原始信息对话框 -->
    <el-dialog
      v-model="rawInstanceDetailVisible"
      title="实例原始信息"
      width="60%"
      append-to-body
    >
      <pre class="json-viewer">{{ JSON.stringify(selectedInstance, null, 2) }}</pre>
    </el-dialog>

    <!-- 添加模型预热对话框 -->
    <el-dialog
      v-model="modelWarmupVisible"
      title="模型预热管理"
      width="800px"
    >
      <div class="warmup-dialog-content">
        <!-- 已有模型列表 -->
        <div class="existing-models-section">
          <h4>当前节点模型列表</h4>
          <div v-if="currentNode?.models?.length" class="model-list">
            <div v-for="model in currentNode.models" :key="model.modelName" class="model-item">
              <div class="model-info">
                <span class="model-name">{{ model.modelName }}</span>
                <el-tag 
                  size="small" 
                  :type="model.modelStatus === 'ready' ? 'success' : 'warning'"
                >
                  {{ model.modelStatus }}
                </el-tag>
              </div>
              <el-button 
                type="danger" 
                size="small"
                @click="handleDeleteModel(model.modelName)"
                :loading="deletingModel === model.modelName"
              >
                删除
              </el-button>
            </div>
          </div>
          <div v-else class="no-models">暂无模型</div>
        </div>
        
        <!-- 模型预热输入 -->
        <div class="warmup-input-section">
          <h4>预热新模型</h4>
          <div class="warmup-input-group">
            <el-input
              v-model="modelToWarmup"
              placeholder="请输入模型名称"
              :disabled="warmingUp"
            />
          </div>
        </div>

        <!-- HF Token配置 -->
        <div class="warmup-input-section">
          <h4>HuggingFace Token</h4>
          <div class="warmup-input-group">
            <el-input
              v-model="hfToken"
              placeholder="HuggingFace Token"
              :disabled="warmingUp"
            />
          </div>
        </div>

        <div class="dialog-footer">
          <el-button
            type="primary"
            @click="handleModelWarmup"
            :loading="warmingUp"
          >
            预热
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 批量模型预热对话框 -->
    <el-dialog
      v-model="batchModelWarmupVisible"
      title="批量模型预热"
      width="800px"
    >
      <div class="warmup-dialog-content">
        <!-- 选中的节点信息 -->
        <div class="selected-nodes-section">
          <h4>已选择节点 ({{ selectedNodes.length }})</h4>
          <div class="selected-nodes-list">
            <el-tag
              v-for="node in selectedNodes"
              :key="node.id"
              class="selected-node-tag"
              closable
              @close="deselectNode(node)"
            >
              {{ node.id }}
            </el-tag>
          </div>
        </div>

        <!-- 模型列表 -->
        <div class="models-section">
          <div class="models-header">
            <h4>预热模型列表</h4>
            <el-button type="primary" size="small" @click="addModelToList">
              添加模型
            </el-button>
          </div>
          
          <div class="model-list">
            <div 
              v-for="(model, index) in modelsToWarmup" 
              :key="index"
              class="model-input-group"
            >
              <div class="model-inputs">
                <el-input
                  v-model="model.name"
                  placeholder="请输入模型名称"
                  :disabled="warmingUp"
                />
                <el-input
                  v-model="model.token"
                  placeholder="HuggingFace Token"
                  :disabled="warmingUp"
                />
                <el-button 
                  type="danger" 
                  size="small"
                  :icon="Delete"
                  circle
                  @click="removeBatchLabel(index)"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <el-button
            type="primary"
            @click="handleBatchModelWarmup"
            :loading="warmingUp"
          >
            开始预热
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 批量设置标签对话框 -->
    <el-dialog
      v-model="batchLabelDialogVisible"
      title="批量设置标签"
      width="600px"
    >
      <div class="batch-label-dialog">
        <!-- 选中的节点信息 -->
        <div class="selected-nodes-section">
          <h4>已选择节点 ({{ selectedNodes.length }})</h4>
          <div class="selected-nodes-list">
            <el-tag
              v-for="node in selectedNodes"
              :key="node.id"
              class="selected-node-tag"
              closable
              @close="deselectNode(node)"
            >
              {{ node.id }}
            </el-tag>
          </div>
        </div>

        <!-- 标签输入区域 -->
        <div class="labels-section">
          <h4>节点标签</h4>
          <div class="labels-list">
            <div 
              v-for="(label, index) in batchLabels" 
              :key="index"
              class="label-item"
            >
              <el-input
                v-model="label.key"
                placeholder="键"
                class="label-key"
              />
              <el-input
                v-model="label.value"
                placeholder="值"
                class="label-value"
              />
              <el-button
                type="danger"
                :icon="Delete"
                circle
                @click="handleRemoveBatchLabel(index)"
              />
            </div>
          </div>
          <div class="add-label">
            <el-button
              type="primary"
              size="small"
              @click="handleAddBatchLabel"
            >
              添加标签
            </el-button>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="batchLabelDialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            @click="handleBatchSetLabels"
            :loading="settingLabels"
          >
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import axios from 'axios'
import dayjs from 'dayjs'

const clusters = ref([])
const activeTab = ref('')
const clusterNodes = ref({})
const loading = ref({})
const searchQuery = ref('')
const showSchedulableOnly = ref(false)
const showAvailableGpuOnly = ref(false)
const selectedGpuType = ref([])
const selectedProvider = ref([])
const labelDialogVisible = ref(false)
const editingLabels = ref({})
const newLabelKey = ref('')
const newLabelValue = ref('')
const currentNode = ref(null)
const nodeDetailVisible = ref(false)
const nodeDetail = ref(null)
const detailLoading = ref(false)
const instanceDetails = ref(new Map())
const columnSettings = ref(JSON.parse(localStorage.getItem('nodeListColumns')) || {
  gpuModel: true,      // GPU型号
  cpuModel: true,      // CPU型号
  provider: true,      // Provider
  cuda: true,          // CUDA
  cpuUsage: true,      // CPU使用
  gpuUsage: true,      // GPU使用
  memory: true,        // 内存
  instanceCount: true, // 实例数
  labels: true,        // 标签
  disks: true,         // 磁盘
  models: true,        // 模型
})
const columnSettingsVisible = ref(false)
const columnWidths = ref(JSON.parse(localStorage.getItem('nodeListColumnWidths')) || {
  id: 200,
  gpuModel: 160,
  cpuModel: 200,
  provider: 120,
  cuda: 100,
  cpuUsage: 120,
  gpuUsage: 120,
  memory: 120,
  instanceCount: 100,
  labels: 500,
  disks: 200,
  operations: 200,
  models: 200
})
const instanceDetailVisible = ref(false)
const instanceDetailLoading = ref(false)
const selectedInstance = ref(null)
const rawInstanceDetailVisible = ref(false)
const modelWarmupVisible = ref(false)
const modelToWarmup = ref('')
const warmingUp = ref(false)
const deletingModel = ref('')
const hfToken = ref('')
const selectedNodes = ref([])
const batchModelWarmupVisible = ref(false)
const modelsToWarmup = ref([])
const defaultHfToken = ''
const batchLabelDialogVisible = ref(false)
const batchLabels = ref([])
const settingLabels = ref(false)

// 生成哈希值的函数
const hashString = (str) => {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // Convert to 32-bit integer
  }
  return Math.abs(hash)
}

// 从哈希值生成颜色
const generateColorFromHash = (str) => {
  const hash = hashString(str)
  // 使用黄金分割比来生成分散的色相值
  const goldenRatio = 0.618033988749895
  const hue = (hash * goldenRatio * 360) % 360
  return `hsl(${hue}, 85%, 95%)`  // 使用高亮度确保文字可读
}

// 从哈希值生成文字颜色
const generateTextColorFromHash = (str) => {
  const hash = hashString(str)
  // 使用黄金分割比来生成分散的色相值
  const goldenRatio = 0.618033988749895
  const hue = (hash * goldenRatio * 360) % 360
  return `hsl(${hue}, 85%, 35%)`  // 使用较低亮度确保文字清晰
}

// 当节点列表中提取所有不同的 GPU 型号
const gpuTypes = computed(() => {
  const nodes = clusterNodes.value[activeTab.value] || []
  const types = new Set(nodes.map(node => node.spec.gpuProductName))
  return Array.from(types)
})

// 从当前集群节点列表中提取所有不同的 Provider
const providers = computed(() => {
  const nodes = clusterNodes.value[activeTab.value] || []
  const types = new Set(nodes.map(node => node.spec.provider))
  return Array.from(types)
})

// 过滤后的点列表
const filteredNodes = computed(() => {
  const nodes = clusterNodes.value[activeTab.value] || []
  return nodes.map(node => ({
    ...node,
    instanceCount: node.state.instIds?.length || 0
  })).filter(node => {
    // 检查是否满足搜索条件（ID、标签、Provider或模型列表）
    const searchTerm = searchQuery.value.toLowerCase()
    const matchesId = node.id.toLowerCase().includes(searchTerm)
    const matchesLabel = Object.entries(node.label || {}).some(([key, value]) => {
      const labelStr = `${key}=${value}`.toLowerCase()
      return labelStr.includes(searchTerm)
    })
    const matchesProvider = node.spec.provider.toLowerCase().includes(searchTerm)
    const matchesModels = node.models?.some(model => {
      return model.modelName.toLowerCase().includes(searchTerm) && model.modelStatus === 'ready'
    })
    const matchesSearch = matchesId || matchesLabel || matchesProvider || matchesModels

    // 检查是否满足调度条件
    const matchesSchedulable = !showSchedulableOnly.value || node.schedulable
    
    // 检查是否满足GPU可用条件
    const hasAvailableGpu = !showAvailableGpuOnly.value || (() => {
      // 如果没有 allocable 或 gpuState，认为可用 GPU
      if (!node.allocable?.gpuState || !Array.isArray(node.allocable.gpuState)) {
        return false
      }
      // 检查是否有空的 GPU
      return node.allocable.gpuState.some(flag => flag === 0)
    })()

    // 检查是否满足GPU型号条件
    const matchesGpuType = selectedGpuType.value.length === 0 || 
      selectedGpuType.value.includes(node.spec.gpuProductName)
    // 检查是否满足Provider条件
    const matchesProviderFilter = selectedProvider.value.length === 0 ||
      selectedProvider.value.includes(node.spec.provider)

    return matchesSearch && matchesSchedulable && hasAvailableGpu && 
      matchesGpuType && matchesProviderFilter
  })
})

// 获取集群列表
const fetchClusters = async () => {
  try {
    const response = await axios.get('/nexusclusters')
    clusters.value = Object.entries(response.data.data).map(([id, cluster]) => ({
      id: id,
      name: id.toUpperCase(),
      ...cluster
    }))
    // 默认选中第一个集群
    if (clusters.value.length > 0) {
      activeTab.value = clusters.value[0].id
      await fetchNodes(clusters.value[0].id)
    }
  } catch (error) {
    console.error('Failed to fetch clusters:', error)
    ElMessage.error('获取集群列表失败')
  }
}

// 取节点列表
const fetchNodes = async (clusterId) => {
  loading.value[clusterId] = true
  try {
    const response = await axios.get(`/nexusclusters/nodes`, {
      params: {
        cluster_id: clusterId
      }
    })
    clusterNodes.value[clusterId] = response.data.data
    ElMessage.success('数据载入成功')
  } catch (error) {
    console.error('Failed to fetch nodes:', error)
    ElMessage.error('获取节点列表失败')
  } finally {
    loading.value[clusterId] = false
  }
}

// 处理标签页切换
const handleTabClick = async (tab) => {
  // 只在该集没有据时才获取
  if (!clusterNodes.value[tab.props.name]) {
    await fetchNodes(tab.props.name)
  }
}

// 刷新当前标签
const refreshCurrentTab = async () => {
  if (activeTab.value) {
    await fetchNodes(activeTab.value)
    ElMessage.success('数据已刷新')
  }
}

// 格式化字节数
const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

// 格式化节数GB
const formatGb = (bytes) => {
  if (!bytes) return 0
  return Math.round(bytes / (1024 * 1024 * 1024))
}

// 计算CPU使用
const computeCpuUsage = (node) => {
  return node.allocable.cpuState.filter(flag => flag === 1).length
}

// 计算GPU使用
const computeGpuUsage = (node) => {
  if (!node.capacity.gpuNum) return 0
  return node.allocable.gpuState.filter(flag => flag === 1).length
}

// 处理调度状变更
const handleScheduleChange = async (node) => {
  try {
    await ElMessageBox.confirm(
      `确定要${node.schedulable ? '停止' : '打开'}该节点的调度吗？`,
      '确认作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await axios.post(`/nexusclusters/nodes/schedulable`, {
      cluster_id: activeTab.value,
      node_id: node.id,
      schedulable: !node.schedulable
    })
    
    ElMessage.success('操作成功')
    await fetchNodes(activeTab.value)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to change node schedulable:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 处理标签编辑
const handleEditLabels = (node) => {
  currentNode.value = node
  editingLabels.value = { ...node.label }
  if (node.label) {
    Object.keys(node.label).forEach(key => {
      if (node.label[key] === '-') {
        editingLabels.value[key] = '-'
      }
    })
  }
  newLabelKey.value = ''
  newLabelValue.value = ''
  labelDialogVisible.value = true
}

// 添加新标签
const addLabel = () => {
  if (!newLabelKey.value || !newLabelValue.value) {
    ElMessage.warning('请输入标签的键值')
    return
  }
  editingLabels.value[newLabelKey.value] = newLabelValue.value
  newLabelKey.value = ''
  newLabelValue.value = ''
}

// 移除标签
const removeLabel = (key) => {
  editingLabels.value[key] = '-'
}

// 保存标签
const saveLabels = async () => {
  try {
    await axios.post(`/nexusclusters/nodes/label`, {
      cluster_id: activeTab.value,
      node_ids: [currentNode.value.id],
      labels: editingLabels.value
    })
    
    ElMessage.success('标签保存成功')
    labelDialogVisible.value = false
    await fetchNodes(activeTab.value)
  } catch (error) {
    console.error('Failed to save labels:', error)
    ElMessage.error('标签保存失败')
  }
}

// 初始化
onMounted(() => {
  fetchClusters()
})

// 示节点详情
const showNodeDetail = async (node) => {
  nodeDetailVisible.value = true
  detailLoading.value = true
  // 重实例详情
  instanceDetails.value.clear()
  
  try {
    const response = await axios.get(`/nexusclusters/nodes/detail`, {
      params: {
        cluster_id: activeTab.value,
        node_id: node.id
      }
    })
    nodeDetail.value = response.data.data
    
    // 等待所有实例情加载完成
    if (nodeDetail.value.state?.instIds?.length) {
      await Promise.all((nodeDetail.value.state.instIds).map(async (instanceId) => {
        try {
          const instanceResponse = await axios.get('/nexusclusters/instances/detail', {
            params: {
              cluster_id: activeTab.value,
              instance_id: instanceId
            }
          })
          // 存储实例详情
          instanceDetails.value.set(instanceId, instanceResponse.data.data)
        } catch (error) {
          console.error(`Failed to fetch instance detail for ${instanceId}:`, error)
          instanceDetails.value.set(instanceId, { error: true })
        }
      }))
    }
  } catch (error) {
    console.error('Failed to fetch node detail:', error)
    ElMessage.error('获取节点详情失败')
  } finally {
    detailLoading.value = false
  }
}

// 计算CPU使用率百分比
const computeCpuUsagePercentage = (node) => {
  const used = computeCpuUsage(node)
  const total = node.capacity.cpuNum
  return Math.round((used / total) * 100)
}

// 计算GPU使用百分比
const computeGpuUsagePercentage = (node) => {
  const used = computeGpuUsage(node)
  const total = node.capacity.gpuNum
  return Math.round((used / total) * 100)
}

// 计算内存使用率百分比
const computeMemoryUsagePercentage = (node) => {
  const used = node.capacity.memorySize - node.allocable.memorySize
  const total = node.capacity.memorySize
  return Math.round((used / total) * 100)
}

// 计算系统盘使用率百分比
const computeRootfsUsagePercentage = (node) => {
  if (!node?.capacity?.rootfsSize) return 0
  const used = node.capacity.rootfsSize - node.allocable.rootfsSize
  const total = node.capacity.rootfsSize
  return Math.round((used / total) * 100)
}

// 计算数据盘使用率百分比
const computeLocalStorageUsagePercentage = (node) => {
  if (!node?.capacity?.localStorageSize) return 0
  const used = node.capacity.localStorageSize - node.allocable.localStorageSize
  const total = node.capacity.localStorageSize
  return Math.round((used / total) * 100)
}

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage < 70) return '#67C23A'
  if (percentage < 90) return '#E6A23C'
  return '#F56C6C'
}

// 格式化日期
const formatDate = (timestamp) => {
  if (!timestamp) return '-'
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

// 获取实例镜像
const getInstanceImage = (instanceId) => {
  const instance = instanceDetails.value.get(instanceId)
  return instance?.containers?.image || '-'
}

// 获取实例的 worker ID
const getWorkerID = (instanceId) => {
  const instance = instanceDetails.value.get(instanceId)
  if (!instance?.containers?.env) return '-'
  // 境是字符串数，需要包 serverless.name 提取值
  const serverlessNameEnv = instance.containers.env.find(env => env.startsWith('serverless.name='))
  if (!serverlessNameEnv) return '-'
  // 分字符串获取部分
  const value = serverlessNameEnv.split('=')[1]
  return value || '-'
}

// 处理实例驱逐
const handleEvictInstance = async (instanceId) => {
  try {
    const workerId = getWorkerID(instanceId)
    if (workerId === '-') {
      ElMessage.error('无法获取 Worker ID')
      return
    }
    
    await ElMessageBox.confirm(
      '确定要驱逐该实例吗？驱逐后该实例将被删除。',
      '确认驱逐',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 直接删除 Worker CR
    const response = await axios.delete(`/workers/${workerId}`, {
      params: {
        se: 'internal-model-api'  // 使用固定的 namespace
      }
    })
    
    if (response.data.code === 200) {
      ElMessage.success('驱逐成功')
      // 等待秒后刷新节点详情，给后端一些处理时间
      setTimeout(async () => {
        await refreshNodeDetail()
      }, 1000)
    } else {
      throw new Error(response.data.message || '驱逐失败')
    }
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to evict instance:', error)
      ElMessage.error(error.response?.data?.message || error.message || '驱逐失败')
    }
  }
}

// 获取实例的 Model ID
const getModelId = (instanceId) => {
  const instance = instanceDetails.value.get(instanceId)
  if (!instance?.containers?.env) return '-'
  const modelIdEnv = instance.containers.env.find(env => env.startsWith('MODEL_ID='))
  if (!modelIdEnv) return '-'
  const value = modelIdEnv.split('=')[1]
  return value || '-'
}

// 获取实例的 创建时间
const getInstanceCreateTime = (instanceId) => {
  const instance = instanceDetails.value.get(instanceId)
  if (!instance?.containers?.env) return '-'
  const createTimeEnv = instance.containers.env.find(env => env.startsWith('serverless.create_time='))
  if (!createTimeEnv) return '-'
  const value = createTimeEnv.split('=')[1]
  return value || '-'
}

// 刷新节点详情
const refreshNodeDetail = async () => {
  if (!nodeDetail.value) return
  detailLoading.value = true
  try {
    const response = await axios.get(`/nexusclusters/nodes/detail`, {
      params: {
        cluster_id: activeTab.value,
        node_id: nodeDetail.value.id
      }
    })
    nodeDetail.value = response.data.data
    
    // 空并重新获取所有实例情
    instanceDetails.value.clear()
    if (nodeDetail.value.state?.instIds?.length) {
      await Promise.all((nodeDetail.value.state.instIds).map(async (instanceId) => {
        try {
          const instanceResponse = await axios.get('/nexusclusters/instances/detail', {
            params: {
              cluster_id: activeTab.value,
              instance_id: instanceId
            }
          })
          instanceDetails.value.set(instanceId, instanceResponse.data.data)
        } catch (error) {
          console.error(`Failed to fetch instance detail for ${instanceId}:`, error)
          instanceDetails.value.set(instanceId, { error: true })
        }
      }))
    }
    ElMessage.success('节点详情已刷新')
  } catch (error) {
    console.error('Failed to refresh node detail:', error)
    ElMessage.error('刷新节点详情失败')
  } finally {
    detailLoading.value = false
  }
}

// 保存列设置
const saveColumnSettings = () => {
  localStorage.setItem('nodeListColumns', JSON.stringify(columnSettings.value))
}

// 监听设置变保存
watch(columnSettings, () => {
  saveColumnSettings()
}, { deep: true })

// 保存列宽设置
const saveColumnWidths = () => {
  localStorage.setItem('nodeListColumnWidths', JSON.stringify(columnWidths.value))
}

// 处理列宽变化
const handleColumnWidthChange = (newWidth, column) => {
  const prop = column.property || column.label.toLowerCase()
  columnWidths.value[prop] = newWidth
  saveColumnWidths()
}

// 显示实例详情
const showInstanceDetail = async (instanceId) => {
  instanceDetailVisible.value = true
  instanceDetailLoading.value = true
  
  try {
    // 如果已经有实例详情，直接使用
    const existingInstance = instanceDetails.value.get(instanceId)
    if (existingInstance) {
      selectedInstance.value = existingInstance
    } else {
      // 否则重新获取实例详情
      const response = await axios.get('/nexusclusters/instances/detail', {
        params: {
          cluster_id: activeTab.value,
          instance_id: instanceId
        }
      })
      selectedInstance.value = response.data.data
      instanceDetails.value.set(instanceId, response.data.data)
    }
  } catch (error) {
    console.error('Failed to fetch instance detail:', error)
    ElMessage.error('获取实例详情失败')
  } finally {
    instanceDetailLoading.value = false
  }
}

// 显示原始信息
const showRawInstanceDetail = () => {
  rawInstanceDetailVisible.value = true
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('复制成功')
  } catch (error) {
    console.error('Copy failed:', error)
    ElMessage.error('复制失败')
  }
}

// 显示模型预热对话框
const showModelWarmupDialog = (node) => {
  currentNode.value = node
  modelWarmupVisible.value = true
  modelToWarmup.value = ''
  hfToken.value = '' // Reset to default value
}

// 处理模型删除
const handleDeleteModel = async (modelName) => {
  try {
    deletingModel.value = modelName
    // 构造符合API要求的请求数据
    const requestData = {
      models: [
        {
          modelName: modelName,
          regions: [
            {
              regionID: activeTab.value,
              servers: [currentNode.value.id]
            }
          ]
        }
      ]
    }
    
    // 使用 axios 的 delete 方法，并传入 data
    await axios.delete('/nodes/models/delete', {
      data: requestData
    })

    ElMessage.success('模型删除成功')
    // 刷新节点数据
    await refreshCurrentTab()
  } catch (error) {
    console.error('删除模型失败:', error)
    ElMessage.error(error.response?.data?.message || '删除模型失败')
  } finally {
    deletingModel.value = ''
  }
}

// 处理模型预热
const handleModelWarmup = async () => {
  if (!modelToWarmup.value) {
    ElMessage.error('请输入模型名称')
    return
  }

  warmingUp.value = true
  try {
    const requestData = {
      node_ids: [currentNode.value.id],
      models: [{
        modelName: modelToWarmup.value,
        gpuCount: currentNode.value.gpu_count
      }],
      hf_token: hfToken.value // 添加 HF Token
    }

    await axios.post('/nodes/models/warmup', requestData)
    ElMessage.success('预热请求已发送')
    modelWarmupVisible.value = false
    modelToWarmup.value = ''
  } catch (error) {
    console.error('Error warming up model:', error)
    ElMessage.error('预热请求失败')
  } finally {
    warmingUp.value = false
  }
}

// 添加模型排序方法
const sortedModels = (models) => {
  return [...models].sort((a, b) => {
    // 首先按状态排序：ready 状态优先
    if (a.modelStatus === 'ready' && b.modelStatus !== 'ready') return -1
    if (a.modelStatus !== 'ready' && b.modelStatus === 'ready') return 1
    
    // 状态相同时按模型名称字母顺序排序
    return a.modelName.localeCompare(b.modelName)
  })
}

// 计算总空闲 GPU 数量
const computeTotalFreeGPUs = computed(() => {
  if (!filteredNodes.value) return 0
  
  return filteredNodes.value.reduce((total, node) => {
    if (!node.capacity.gpuNum) return total
    const usedGPUs = computeGpuUsage(node)
    const freeGPUs = node.capacity.gpuNum - usedGPUs
    return total + freeGPUs
  }, 0)
})

// 处理多选
const handleSelectionChange = (selection) => {
  selectedNodes.value = selection
}

// 取消选择节点
const deselectNode = (node) => {
  const tableRef = document.querySelector('.el-table')
  if (tableRef) {
    tableRef.__vue__?.toggleRowSelection(node, false)
  }
}

// 显示批量预热对话框
const showBatchModelWarmupDialog = () => {
  modelsToWarmup.value = [{
    name: '',
    token: defaultHfToken
  }]
  batchModelWarmupVisible.value = true
}

// 添加模型到列表
const addModelToList = () => {
  modelsToWarmup.value.push({
    name: '',
    token: defaultHfToken
  })
}

// 从列表中移除模型
const removeBatchLabel = (index) => {
  modelsToWarmup.value.splice(index, 1)
}

// 处理批量预热
const handleBatchModelWarmup = async () => {
  // 验证输入
  if (modelsToWarmup.value.some(model => !model.name)) {
    ElMessage.error('请填写所有模型名称')
    return
  }

  warmingUp.value = true
  try {
    const requestData = {
      node_ids: selectedNodes.value.map(node => node.id),
      models: modelsToWarmup.value.map(model => ({
        modelName: model.name,
        hf_token: model.token
      }))
    }

    await axios.post('/nodes/models/warmup', requestData)
    ElMessage.success('预热请求已发送')
    batchModelWarmupVisible.value = false
  } catch (error) {
    console.error('Error warming up models:', error)
    ElMessage.error('预热请求失败')
  } finally {
    warmingUp.value = false
  }
}

// 批量标签操作相关的方法
const handleAddBatchLabel = () => {
  batchLabels.value.push({ key: '', value: '' })
}

const handleRemoveBatchLabel = (index) => {
  batchLabels.value.splice(index, 1)
}

const handleShowBatchLabelDialog = () => {
  // 如果选中的节点都有相同的标签，则预填充这些标签
  const commonLabels = selectedNodes.value.reduce((acc, node, index) => {
    const nodeLabels = node.labels || {}
    if (index === 0) {
      return Object.entries(nodeLabels).map(([key, value]) => ({ key, value }))
    }
    return acc.filter(label => 
      nodeLabels[label.key] !== undefined && nodeLabels[label.key] === label.value
    )
  }, [])
  
  batchLabels.value = commonLabels.length ? commonLabels : [{ key: '', value: '' }]
  batchLabelDialogVisible.value = true
}

// 处理批量设置标签
const handleBatchSetLabels = async () => {
  // 验证标签输入
  if (!batchLabels.value.length || batchLabels.value.some(label => !label.key.trim())) {
    ElMessage.warning('请至少输入一个有效的标签键')
    return
  }

  settingLabels.value = true
  try {
    const labels = batchLabels.value.reduce((acc, { key, value }) => {
      if (key.trim()) {
        acc[key.trim()] = value.trim()
      }
      return acc
    }, {})

    // 对每个选中的节点调用设置标签接口
    const promises = selectedNodes.value.map(node => 
      axios.post(`/nexusclusters/nodes/label`, {
        node_ids: [node.id],
        labels: labels,
        cluster_id: activeTab.value
      })
    )
    
    await Promise.all(promises)
    ElMessage.success('标签设置成功')
    batchLabelDialogVisible.value = false
    refreshCurrentTab()
  } catch (error) {
    console.error('Error setting labels:', error)
    ElMessage.error('标签设置失败')
  } finally {
    settingLabels.value = false
  }
}
</script>

<style scoped>
.page-container {
  padding: 24px;
  margin: 0 auto;
  min-height: calc(100vh - var(--navbar-height));
}

.header {
  margin-bottom: 24px;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 24px;
  color: #3b82f6;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.cluster-tabs {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 100%;
  overflow-x: hidden;
}

.content-card {
  padding: 16px;
  overflow-x: auto;
}

.label-tag {
  white-space: normal;
  word-break: break-all;
  overflow: visible;
  font-size: 14px;
  text-align: left;
  margin: 2px;
  transition: all 0.3s;
  border-radius: 4px;
}

:deep(.el-tabs__header) {
  margin-bottom: 0;
}

:deep(.el-tabs__nav) {
  border: none;
}

:deep(.el-tabs__item) {
  font-size: 14px;
  font-weight: 500;
  padding: 0 24px;
  height: 48px;
  line-height: 48px;
}

:deep(.el-tabs__item.is-active) {
  color: #3b82f6;
  font-weight: 600;
}

.table-header {
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.node-id {
  font-family: 'Roboto Mono', monospace;
  font-weight: 500;
  color: #1a56db;
  font-size: 15px;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-icon {
  font-size: 16px;
  cursor: help;
}

.unschedulable {
  color: #dc2626;
  margin-left: 4px;
  font-size: 14px;
  opacity: 0.9;
}

.resource-value {
  font-weight: 500;
  font-size: 15px;
}

.label-container {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  justify-content: flex-start;
  padding: 4px 0;
}

.label-tag {
  white-space: normal;
  word-break: break-all;
  overflow: visible;
  font-size: 14px;
  text-align: left;
  margin: 2px;
  transition: all 0.3s;
  border-radius: 4px;
}

.highlight {
  color: #dc2626;
  font-weight: 600;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

:deep(.el-table__header) {
  border-bottom: 2px solid #e5e7eb;
}

:deep(.el-table__row) {
  border-bottom: 1px solid #f3f4f6;
}

.table-filters {
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-input {
  width: 350px;
}

.gpu-filter {
  width: 300px;
}

.provider-filter {
  width: 200px;
}

.instance-count {
  display: inline-block;
  min-width: 30px;
  text-align: center;
  font-size: 15px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  height: 100%;
  align-items: center;
}

.label-edit-container {
  max-height: 400px;
  overflow-y: auto;
}

.label-item {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.label-input {
  flex: 1;
}

.add-label {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}

.label-key-input {
  width: 500px;
}

.label-value-input {
  width: 200px;
}

.label-key {
  min-width: 500px;
  max-width: 500px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: left;
  padding: 0 8px;
}

.node-detail {
  padding: 20px;
  background-color: #f9fafb;
  border-radius: 8px;
}

.detail-section {
  margin-bottom: 24px;
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  
  .section-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 0 0 16px;
    
    i {
      font-size: 20px;
      color: #3b82f6;
    }
    
    h3 {
      margin: 0;
      font-size: 16px;
      color: #374151;
    }
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
  align-items: center;
}

.info-item {
  display: flex;
  gap: 8px;
  
  .label {
    color: #6B7280;
    font-weight: 500;
  }
  
  .value {
    color: #111827;
  }
}

.resource-usage {
  display: flex;
  justify-content: flex-start;
  flex-wrap: nowrap;
  gap: 16px;
  align-items: center;
  padding: 0;
  width: 100%;
}

.usage-metrics {
  display: flex;
  gap: 16px;
  flex: 2;
  justify-content: center;
  background-color: #f8fafc;
  border-radius: 8px;
  padding: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.instance-id {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;
  font-size: 12px;
}

.gpu-usage-container {
  background-color: #f8fafc;
  border-radius: 8px;
  padding: 16px;
  flex: 3;
  display: flex;
  align-items: center;
  gap: 12px;
}

.gpu-usage-header {
  min-width: 100px;
  flex-shrink: 0;
  padding-right: 16px;
  border-right: 1px solid #e5e7eb;
}

.gpu-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(60px, 1fr));
  grid-auto-flow: row;
  gap: 12px;
  flex: 1;
  padding-left: 16px;
}

.gpu-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 6px;
  background-color: #e5e7eb;
  border-radius: 6px;
  transition: all 0.3s;
  justify-content: center;
  min-width: 60px;
  
  i {
    font-size: 18px;
    color: #4b5563;
  }
  
  span {
    font-size: 11px;
    color: #4b5563;
    white-space: nowrap;
  }

  &.gpu-used {
    background-color: #fee2e2;
    
    i, span {
      color: #dc2626;
    }
  }
}

.instance-id {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;
  font-size: 12px;
}

.image-name {
  display: inline-block;
  max-width: 460px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;
  font-size: 12px;
}

.model-id {
  display: inline-block;
  max-width: 330px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;
  font-size: 12px;
  color: #666;
}

.worker-id {
  display: inline-block;
  max-width: 480px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: monospace;
  color: #666;
  font-size: 12px;
  letter-spacing: -0.2px;
}

.clickable {
  cursor: pointer;
  &:hover {
    color: #3b82f6;
    text-decoration: underline;
  }
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  padding: 4px 0;
}

/* 硬件信息部分的特殊样式 */
.hardware-info {
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)) !important;
}

.progress-label {
  text-align: center;
  font-size: 14px;
  
  > div:first-child {
    color: #6B7280;
    margin-bottom: 4px;
  }
  
  > div:last-child {
    font-weight: 500;
    color: #111827;
    font-size: 15px;
  }
}

/* 确保进度条容器有固定宽度 */
:deep(.el-progress-circle) {
  margin: 0 auto;
}

/* 优化表格样式 */
:deep(.el-table) {
  border-spacing: 0;
  border-collapse: collapse;

  .el-table__cell {
    padding: 2px !important;
  }

  .cell {
    padding-left: 4px !important;
    padding-right: 4px !important;
  }

  .el-table__body {
    width: 100% !important;
  }

  .el-table__header {
    width: 100% !important;
  }
}

.dialog-title-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 100%;
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  padding: 4px 0;
}

/* 调整表格内的 el-tag 样式 */
:deep(.el-table .el-tag) {
  font-size: 14px;
  padding: 4px 8px;
}

/* 调整操作按钮的样式 */
:deep(.action-buttons .el-button) {
  padding: 8px 16px;
  height: 36px;
  font-size: 14px;
}

/* 自定义 tab 样式 */
.custom-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 20px;
    border-bottom: none;
  }

  :deep(.el-tabs__nav) {
    border: none;
    gap: 8px;
  }

  :deep(.el-tabs__item) {
    height: 40px;
    line-height: 40px;
    font-size: 15px;
    font-weight: 500;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    margin-right: 8px;
    transition: all 0.3s;
    
    &:hover {
      color: #3b82f6;
      border-color: #3b82f6;
    }
    
    &.is-active {
      color: #3b82f6;
      border-color: #3b82f6;
      background-color: #eff6ff;
    }
  }

  :deep(.el-tabs__nav-wrap::after) {
    display: none;
  }
}

.page-header {
  padding: 16px 0 24px;
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
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #111827;
  line-height: 1.2;
}

.header-divider {
  margin-top: 16px;
  height: 1px;
  background: linear-gradient(to right, #e5e7eb 0%, transparent 100%);
}

.resource-usage-grid {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  gap: 20px;
  padding: 0 16px;
}

.usage-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;
}

.progress-label {
  text-align: center;
  font-size: 13px;
  margin-bottom: 8px;
}

.progress-label > div:first-child {
  color: #666;
  font-size: 12px;
}

.progress-label > div:last-child {
  color: #333;
  font-size: 13px;
  font-weight: 500;
}

.gpu-usage {
  margin-bottom: 20px;
  padding: 0 16px;
  width: 100%;
}

.gpu-cards-container {
  width: 100%;
  overflow-x: hidden;
}

.gpu-cards {
  display: flex;
  gap: 12px;
  flex-wrap: nowrap;
  width: 100%;
  justify-content: space-between;
}

.gpu-card {
  padding: 8px 0;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 13px;
  color: #666;
  border: 1px solid #e5e7eb;
  flex: 1;
  text-align: center;
  min-width: 80px;
}

.gpu-card.gpu-used {
  background: #fef2f2;
  color: #ef4444;
  border-color: #fecaca;
}

.gpu-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  color: #666;
  font-size: 14px;
  font-weight: 500;
}

.divider {
  height: 1px;
  background: #e5e7eb;
  margin: 16px 0;
}

.label-key {
  min-width: 300px;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: left;
  padding: 0 8px;
}

:deep(.el-input-group__prepend) {
  padding: 0;
  min-width: 500px;
  max-width: 500px;
  background-color: #f5f7fa;
}

.column-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 8px;
}

.column-settings .el-checkbox {
  margin-right: 0;
}

:deep(.el-table .el-table__cell.is-left) {
  text-align: left !important;
}

/* 确保表格内容不会被截断 */
:deep(.el-table__body-wrapper) {
  overflow-y: auto !important;
}

/* 允许单元格内容换行 */
:deep(.el-table .cell) {
  white-space: normal;
  height: auto;
  line-height: 1.5;
}

.instance-id.clickable {
  cursor: pointer;
  color: #3b82f6;
  transition: color 0.2s;
  
  &:hover {
    color: #2563eb;
    text-decoration: underline;
  }
}

.instance-detail {
  padding: 16px;
  background-color: #f9fafb;
  border-radius: 8px;
}

.dialog-header-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-bottom: 16px;
}

.instance-details-dialog {
  :deep(.el-dialog__header) {
    padding: 16px 24px;
    margin-right: 0;
    background-color: #f0f7ff;
    border-bottom: 1px solid #e5e7eb;
    border-radius: 8px 8px 0 0;
  }

  :deep(.el-dialog__title) {
    font-size: 20px;
    font-weight: 700;
    color: #1d4ed8;
    font-family: system-ui, -apple-system, sans-serif;
  }

  :deep(.el-dialog__headerbtn) {
    top: 18px;
    right: 20px;
    
    .el-dialog__close {
      font-size: 18px;
      color: #6b7280;
      
      &:hover {
        color: #1d4ed8;
      }
    }
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }
}

.detail-section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 12px;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    bottom: -2px;
    left: 0;
    width: 100px;
    height: 2px;
    background: linear-gradient(90deg, #3b82f6, transparent);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;

    i {
      font-size: 18px;
      color: #1d4ed8;
      background: #e0e7ff;
      padding: 8px;
      border-radius: 8px;
    }

    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #1f2937;
    }
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.info-item.double {
  display: flex;
  justify-content: space-between;
  
  .info-pair {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
  }
}

.label {
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
  min-width: 100px;
}

.value {
  color: #111827;
  font-weight: 500;
  
  &.code {
    font-family: monospace;
    background: #f8fafc;
    padding: 6px 12px;
    border-radius: 4px;
    font-size: 13px;
    border: 1px solid #e5e7eb;
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.ssh-info-row {
  display: flex;
  gap: 24px;
  width: 100%;
}

.ssh-item {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.code-value {
  font-family: monospace;
  background: #f8fafc;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  border: 1px solid #e5e7eb;
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.copy-btn {
  flex-shrink: 0;
}

.log-link {
  color: #3b82f6;
  text-decoration: none;
  font-family: monospace;
  font-size: 13px;
  
  &:hover {
    text-decoration: underline;
    color: #2563eb;
  }
}

.dialog-header-actions {
  position: absolute;
  top: 20px;
  right: 50px;
  z-index: 10;
}

.port-mapping {
  grid-column: 1 / -1;
  margin-bottom: 8px;
}

.ssh-info {
  grid-column: 1 / -1;
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.ssh-info-row {
  display: flex;
  gap: 24px;
  width: 100%;
}

.ssh-item {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.port-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.value-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.warmup-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.existing-models-section,
.warmup-input-section {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
}

.existing-models-section h4,
.warmup-input-section h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.warmup-input-group {
  display: flex;
  gap: 12px;
}

.model-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.model-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.model-name {
  font-weight: 500;
  color: #374151;
}

.model-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.no-models {
  color: #6b7280;
  text-align: center;
  padding: 16px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  margin-top: 12px;
}

.gpu-total-info {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #606266;
  padding: 0 12px;
  border-left: 1px solid #dcdfe6;
}

.gpu-total-value {
  font-weight: 600;
  color: #3b82f6;
  margin-left: 4px;
}

.table-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.selected-nodes-section {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
}

.selected-nodes-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.selected-node-tag {
  margin: 4px;
}

.models-section {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
}

.models-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.model-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.model-input-group {
  background: white;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.model-inputs {
  display: flex;
  gap: 12px;
  align-items: center;
}

.model-inputs .el-input {
  flex: 1;
}

.record-info {
  color: #606266;
  font-size: 14px;
  margin-left: 16px;
  display: flex;
  align-items: center;
}

.batch-label-dialog {
  .selected-nodes-section {
    margin-bottom: 24px;
  }

  .labels-section {
    h4 {
      margin-bottom: 12px;
    }
  }

  .labels-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .label-item {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .label-key,
  .label-value {
    flex: 1;
  }

  .add-label {
    margin-top: 12px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style> 