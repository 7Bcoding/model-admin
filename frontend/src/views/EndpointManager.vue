<template>
  <div class="endpoint-manager">
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <div class="d-flex align-items-center">
          <h3 class="card-title mb-0 me-3">
            <span class="text-muted">模型：</span>{{ modelName }}
          </h3>
          <div class="btn-group">
            <button 
              class="btn btn-sm btn-outline-primary"
              @click="openMonitor('overview')"
              title="在 Grafana 中查看监控"
            >
              <i class="bi bi-graph-up"></i>
              SLA 
            </button>
            <button 
              class="btn btn-sm btn-outline-secondary"
              @click="openMonitor('vllm')"
              title="查看模型监控"
            >
              <i class="bi bi-graph-up"></i>
              {{ modelData?.inference_engine || 'vllm' }}
            </button>
            <button 
              v-if="hasEndpoints"
              class="btn btn-sm btn-outline-info"
              @click="openMonitor('serverless')"
              title="查看 Serverless 监控"
            >
              <i class="bi bi-graph-up"></i>
              serverless
            </button>
          </div>
        </div>
        <button 
          v-if="canModify"
          class="btn btn-success" 
          @click="showAddEndpointModal"
        >
          <i class="bi bi-plus-circle"></i> 添加 Endpoint
        </button>
      </div>

      <div class="card-body">
        <div class="table-responsive">
          <table class="table">
            <thead>
              <tr>
                <th>Endpoint ID</th>
                <th>URL</th>
                <th class="text-center">权重</th>
                <th class="text-center">健康检查</th>
                <th class="text-center">测试模式</th>
                <th class="text-center">GPU 资源</th>
                <th class="text-center">实例数</th>
                <th class="text-center">支持API</th>
                <th class="text-center">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="endpoint in endpoints" :key="endpoint.endpoint_id">
                <td>{{ endpoint.endpoint_id }}</td>
                <td>
                  <span>{{ endpoint.url }}</span>
                  <template v-if="getEndpointTypeTag(endpoint.url)">
                    <el-tag :style="{marginLeft: '6px', background: getEndpointTypeTag(endpoint.url).color, color: '#fff', border: 'none'}" size="small" :title="getEndpointTypeTag(endpoint.url).tip">
                      {{ getEndpointTypeTag(endpoint.url).tag }}
                    </el-tag>
                  </template>
                </td>
                <td class="text-center">
                  <input 
                    type="number" 
                    class="form-control form-control-sm weight-input"
                    v-model="endpoint.weight"
                    :disabled="!canModify"
                  >
                </td>
                <td class="text-center">
                  <div class="form-check d-flex justify-content-center">
                    <input 
                      type="checkbox"
                      class="form-check-input"
                      v-model="endpoint.enable_check_health"
                      :disabled="!canModify"
                    >
                  </div>
                </td>
                <td class="text-center">
                  <div class="form-check d-flex justify-content-center">
                    <input 
                      type="checkbox"
                      class="form-check-input"
                      v-model="endpoint.enable_testflight"
                      :disabled="!canModify"
                    >
                  </div>
                </td>
                <td class="text-center">
                  <div class="gpu-info">
                    <template v-if="endpoint.gpuInfo">
                      <span v-if="endpoint.gpuInfo.loading" class="spinner-border spinner-border-sm text-info" role="status">
                        <span class="visually-hidden">Loading...</span>
                      </span>
                      <span v-else-if="endpoint.gpuInfo.error" class="text-danger">
                        {{ endpoint.gpuInfo.error }}
                      </span>
                      <span v-else-if="endpoint.gpuInfo.model && endpoint.gpuInfo.number" class="gpu-spec">
                        {{ endpoint.gpuInfo.model }} × {{ endpoint.gpuInfo.number }}
                      </span>
                      <span v-else>-</span>
                    </template>
                    <span v-else>-</span>
                  </div>
                </td>
                <td class="text-center">
                  <template v-if="endpoint.worker_count === null">
                    <span class="spinner-border spinner-border-sm text-info" role="status">
                      <span class="visually-hidden">Loading...</span>
                    </span>
                  </template>
                  <template v-else>
                    {{ endpoint.worker_count }}
                  </template>
                </td>
                <td class="text-center">
                  <template v-if="endpoint.supported_api_flag === null">
                    <span class="spinner-border spinner-border-sm text-info" role="status">
                      <span class="visually-hidden">Loading...</span>
                    </span>
                  </template>
                  <template v-else>
                    <el-select 
                      v-model="endpoint.supported_api_flag" 
                      size="small"
                      :disabled="!canModify"
                      class="api-flag-select"
                    >
                      <el-option 
                        label="仅支持对话" 
                        value="ENDPOINT_SUPPORTED_API_FLAG_CHAT_COMPLETIONS" 
                      />
                      <el-option 
                        label="仅支持Completion" 
                        value="ENDPOINT_SUPPORTED_API_FLAG_COMPLETIONS" 
                      />
                      <el-option 
                        label="全部支持" 
                        value="ENDPOINT_SUPPORTED_API_FLAG_ALL" 
                      />
                      <el-option 
                        label="claude协议" 
                        value="ENDPOINT_SUPPORTED_API_FLAG_ANTHROPIC_MESSAGES" 
                      />
                    </el-select>
                  </template>
                </td> 
                <td class="text-center">
                  <div class="action-buttons">
                    <el-button
                      size="small"
                      class="action-btn"
                      type="success"
                      @click="showEndpointDetails(endpoint)"
                    >
                      详情
                    </el-button>
                    <el-button
                      size="small"
                      class="action-btn"
                      type="info"
                      @click="showDeploymentInfo(endpoint)"
                    >
                      参数
                    </el-button>
                    <el-button
                      size="small"
                      class="action-btn"
                      type="warning"
                      @click="openGrafanaLink(endpoint)"
                    >
                      面板
                    </el-button>
                    <el-button 
                      v-if="canModify"
                      size="small"
                      class="action-btn"
                      type="primary"
                      @click="updateEndpoint(endpoint)"
                    >
                      更新
                    </el-button>
                    <el-button
                      v-if="canModify"
                      size="small"
                      class="action-btn"
                      type="danger"
                      @click="toggleEndpointStatus(endpoint)"
                    >
                      {{ endpoint.status === 'ENDPOINT_STATUS_SERVING' ? '停用' : '启用' }}
                    </el-button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 权重分布图表 -->
        <div class="weight-distribution mt-4">
          <h5>权重分布</h5>
          <div class="row">
            <div class="col-md-8">
              <div ref="weightChart" class="chart-container"></div>
            </div>
            <div class="col-md-4">
              <div class="legend-container">
                <div v-for="endpoint in endpoints" :key="endpoint.endpoint_id" class="legend-item">
                  <span class="legend-color" :style="{ backgroundColor: getEndpointColor(endpoint) }"></span>
                  <div class="legend-info">
                    <div class="legend-name">
                      {{ endpoint.endpoint_id }}
                      <!-- 为包含"融"和"2.0"的endpoint添加标记 -->
                      <span v-if="endpoint.url.includes('/fusion/v1/')" class="special-tag rong-tag">融</span>
                      <span v-if="endpoint.url.includes('sls2.alpha.ai')" class="special-tag version-tag">2.0</span>
                    </div>
                    <div class="legend-url text-muted clickable" @click="goToServerless(getEndpointName(endpoint.url), endpoint.url)">
                      {{ getEndpointName(endpoint.url) }}
                    </div>
                    <div class="legend-weight">权重: {{ endpoint.weight }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加 Endpoint 对话框 -->
    <el-dialog
      v-model="addEndpointDialogVisible"
      title="添加 Endpoint"
      width="40%"
    >
      <el-form :model="newEndpoint" label-width="140px">
        <el-form-item label="Endpoint ID">
          <el-input v-model="newEndpoint.endpoint_id" placeholder="请输入 Endpoint ID"></el-input>
        </el-form-item>

        <el-form-item label="URL">
          <el-input v-model="newEndpoint.url" placeholder="请输入 URL">
            <template #append>
              <el-button type="primary" @click="testNewEndpoint">测试</el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="权重">
          <el-input-number v-model="newEndpoint.weight" :min="1" :max="100" :step="1"></el-input-number>
        </el-form-item>

        <el-form-item label="Model ID Override">
          <el-input 
            v-model="newEndpoint.model_id_override" 
            placeholder="可选，用于覆盖默认模型ID"
          ></el-input>
          <div class="text-muted small">如果需要使用不同的模型ID，请在此处填写</div>
        </el-form-item>

        <div class="options-row">
          <el-form-item class="checkbox-item">
            <el-checkbox v-model="newEndpoint.enable_check_health">启用健康检查</el-checkbox>
          </el-form-item>
          <el-form-item class="checkbox-item">
            <el-checkbox v-model="newEndpoint.enable_testflight">启用测试模式</el-checkbox>
          </el-form-item>
          <el-form-item class="checkbox-item">
            <el-checkbox v-model="newEndpoint.skip_testflight_while_add">强制跳过testflight</el-checkbox>
          </el-form-item>
          <el-form-item class="select-item" label="支持API">
            <el-select 
              v-model="newEndpoint.supported_api_flag" 
              placeholder="请选择支持的API类型"
              class="api-flag-select"
            >
              <el-option 
                label="仅支持对话" 
                value="ENDPOINT_SUPPORTED_API_FLAG_CHAT_COMPLETIONS" 
              />
              <el-option 
                label="仅支持Completion" 
                value="ENDPOINT_SUPPORTED_API_FLAG_COMPLETIONS" 
              />
              <el-option 
                label="全部支持" 
                value="ENDPOINT_SUPPORTED_API_FLAG_ALL" 
              />
              <el-option 
                label="claude协议" 
                value="ENDPOINT_SUPPORTED_API_FLAG_ANTHROPIC_MESSAGES" 
              />
            </el-select>
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addEndpointDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addEndpoint">添加</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 添加详情对话框 -->
    <el-dialog
      v-model="detailsDialogVisible"
      :title="currentEndpoint?.endpoint_id + ' 详情'"
      width="55%"
      class="endpoint-details-dialog"
    >
      <div v-loading="detailsLoading">
        <!-- SAP 信息 -->
        <div class="details-section">
          <div class="section-header">
            <h4>自动扩缩容配置</h4>
            <div class="header-right">
              <template v-if="!isEditingSap">
                <el-button type="primary" size="small" @click="startEditSap">
                  编辑
                </el-button>
              </template>
              <template v-else>
                <el-button type="success" size="small" @click="saveSapChanges">
                  保存
                </el-button>
                <el-button type="default" size="small" @click="cancelEditSap">
                  取消
                </el-button>
              </template>
              <el-button type="text" @click="showFullSapDetails">查看完整配置</el-button>
            </div>
          </div>
          <div class="auto-scaling-info">
            <div class="scaling-item">
              <div class="label">最小副本数</div>
              <div class="value">
                <template v-if="isEditingSap">
                  <el-input-number 
                    v-model="editingSapValues.minReplicas" 
                    :min="0"
                    :precision="0"
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="custom-number-input"
                  />
                </template>
                <template v-else>
                  {{ sapDetails.spec.minReplicas }}
                </template>
              </div>
            </div>
            <div class="scaling-item">
              <div class="label">最大副本数</div>
              <div class="value">
                <template v-if="isEditingSap">
                  <el-input-number 
                    v-model="editingSapValues.maxReplicas" 
                    :min="0"
                    :precision="0"
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="custom-number-input"
                  />
                </template>
                <template v-else>
                  {{ sapDetails.spec.maxReplicas }}
                </template>
              </div>
            </div>
            <div class="scaling-item">
              <div class="label">触发扩容并发</div>
              <div class="value">
                <template v-if="isEditingSap">
                  <el-input-number 
                    v-model="editingSapValues.concurrencyPerWorker" 
                    :min="0"
                    :precision="0"
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="custom-number-input"
                  />
                </template>
                <template v-else>
                  {{ getConcurrencyTarget(sapDetails.spec) || '-' }}
                </template>
              </div>
            </div>
            <div class="scaling-item">
              <div class="label">扩容窗口时间</div>
              <div class="value">
                <template v-if="isEditingSap">
                  <el-input-number 
                    v-model="editingSapValues.scaleUpWindow" 
                    :min="0"
                    :precision="0"
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="custom-number-input"
                  />
                </template>
                <template v-else>
                  {{ sapDetails.spec.behavior?.scaleUp?.stabilizationWindowSeconds || '-' }}
                </template>
              </div>
            </div>
            <div class="scaling-item">
              <div class="label">缩容窗口时间</div>
              <div class="value">
                <template v-if="isEditingSap">
                  <el-input-number 
                    v-model="editingSapValues.scaleDownWindow" 
                    :min="0"
                    :precision="0"
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="custom-number-input"
                  />
                </template>
                <template v-else>
                  {{ sapDetails.spec.behavior?.scaleDown?.stabilizationWindowSeconds || '-' }}
                </template>
              </div>
            </div>
          </div>
        </div>

        <!-- Worker 信息 -->
        <div class="details-section">
          <div class="section-header">
            <h4>Worker 列表</h4>
            <div class="header-right">
              <span class="worker-stats">
                总数: {{ workerDetails.length }} / 运行中: {{ servingWorkerCount }}
              </span>
              <el-button
                type="primary"
                link
                @click="refreshWorkerList"
                :loading="workerListRefreshing"
              >
                刷新列表
              </el-button>
            </div>
          </div>
          <el-table :data="workerDetails" style="width: 100%" border>
            <el-table-column label="名称" min-width="280">
              <template #default="{ row }">
                <div class="worker-name">{{ row.metadata.name }}</div>
              </template>
            </el-table-column>
            <el-table-column label="并发数" width="80" align="center">
              <template #default="{ row }">
                <span 
                  :class="['concurrency-value', { 'concurrency-updated': row.concurrencyUpdated }]"
                >
                  {{ row.concurrency || 0 }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="spec.provider" label="Provider" width="90" />
            <el-table-column prop="spec.clusterIDs[0]" label="集群" width="90" />
            <el-table-column label="实例 ID" width="180">
              <template #default="{ row }">
                <span 
                  class="instance-id clickable"
                  @click="showInstanceDetail(row.status.realInstanceID, row.spec.clusterIDs[0])"
                >
                  {{ row.status.realInstanceID }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="90" align="center">
              <template #default="{ row }">
                <el-button
                  type="warning"
                  size="small"
                  @click="openWorkerGrafana(row)"
                >
                  面板
                </el-button>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120" align="center">
              <template #default="{ row }">
                <div class="status-container">
                  <el-tag :type="getStateType(row.status.state)" size="small">
                    {{ row.status.state }}
                  </el-tag>
                  <el-tag 
                    :type="row.status.healthy ? 'success' : 'danger'"
                    size="small"
                    class="health-tag"
                  >
                    {{ row.status.healthy ? '健康' : '异常' }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column 
              prop="metadata.creationTimestamp" 
              label="创建时间" 
              width="180"
              sortable
            >
              <template #default="scope">
                {{ formatDate(scope.row.metadata.creationTimestamp) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-dialog>

    <!-- SAP 完整信息对话框 -->
    <el-dialog
      v-model="sapFullDetailsVisible"
      title="完整 SAP 配置"
      width="60%"
      append-to-body
    >
      <pre class="json-viewer">{{ JSON.stringify(sapDetails, null, 2) }}</pre>
    </el-dialog>

    <!-- 实例详情对话框 -->
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
        <template v-if="instanceDetail">
          <div class="detail-section">
            <div class="section-header">
              <div class="header-left">
                <i class="el-icon-info"></i>
                <h3>基本信息息</h3>
              </div>
            </div>
            <div class="info-grid">
              <div class="info-item double">
                <div class="info-pair">
                  <span class="label">实例 ID:</span>
                  <span class="value">{{ instanceDetail.id }}</span>
                </div>
                <div class="info-pair">
                  <span class="label">节点 ID:</span>
                  <span class="value">{{ instanceDetail.nodeId }}</span>
                </div>
              </div>
              <div class="info-item double">
                <div class="info-pair">
                  <span class="label">状态:</span>
                  <el-tag :type="instanceDetail.state?.state === 'running' ? 'success' : 'warning'">
                    {{ instanceDetail.state?.state || '-' }}
                  </el-tag>
                </div>
                <div class="info-pair">
                  <span class="label">实例类型:</span>
                  <span class="value">{{ instanceDetail.kind || '-' }}</span>
                </div>
              </div>
              <div class="info-item">
                <span class="label">镜像:</span>
                <span class="value code">{{ instanceDetail.containers?.image || '-' }}</span>
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
                <span class="value">{{ instanceDetail.containers?.cpuNum || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">内存:</span>
                <span class="value">{{ formatGb(instanceDetail.containers?.memorySize) }}GB</span>
              </div>
              <div class="info-item">
                <span class="label">GPU 数量:</span>
                <span class="value">{{ instanceDetail.containers?.gpuNum || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">GPU 型号:</span>
                <span class="value">{{ instanceDetail.containers?.gpuProductName || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">系统盘:</span>
                <span class="value">{{ formatGb(instanceDetail.containers?.rootfsSize) }}GB</span>
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
              <div v-if="instanceDetail.exposePortStates?.length" class="info-item full-width">
                <span class="label">端口映射:</span>
                <div class="value-group">
                  <div v-for="port in instanceDetail.exposePortStates" :key="port.port" class="port-item">
                    <span class="value">
                      {{ port.exposeAddress }}
                      <el-tag size="small" :type="port.online ? 'success' : 'danger'" class="ml-2">
                        {{ port.online ? '在线' : '离线' }}
                      </el-tag>
                    </span>
                  </div>
                </div>
              </div>
              <div v-if="instanceDetail.apps?.sshState" class="info-item full-width">
                <div class="ssh-info-row">
                  <div class="ssh-item">
                    <span class="label">SSH 命令:</span>
                    <span class="code-value">{{ instanceDetail.apps.sshState.sshCommand }}</span>
                    <el-button
                      type="primary"
                      size="small"
                      class="copy-btn"
                      @click="copyToClipboard(instanceDetail.apps.sshState.sshCommand)"
                    >
                      复制命令
                    </el-button>
                  </div>
                  <div class="ssh-item">
                    <span class="label">SSH 密码:</span>
                    <span class="code-value">{{ instanceDetail.apps.sshState.sshPassword }}</span>
                    <el-button
                      type="primary"
                      size="small"
                      class="copy-btn"
                      @click="copyToClipboard(instanceDetail.apps.sshState.sshPassword)"
                    >
                      复制密码
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 日志信息 -->
          <div class="detail-section" v-if="instanceDetail.apps?.logState">
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
                    :href="`${instanceDetail.apps.logState.systemLogAddress}?follow=1&tail=100`"
                    target="_blank"
                    class="log-link"
                  >
                    {{ instanceDetail.apps.logState.systemLogAddress }}
                  </a>
                </div>
              </div>
              <div class="info-item">
                <span class="label">实例日志:</span>
                <div class="value-group">
                  <a 
                    :href="`${instanceDetail.apps.logState.instanceLogAddress}?follow=1&tail=100`"
                    target="_blank"
                    class="log-link"
                  >
                    {{ instanceDetail.apps.logState.instanceLogAddress }}
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
      <pre class="json-viewer">{{ JSON.stringify(instanceDetail, null, 2) }}</pre>
    </el-dialog>

    <!-- 部署参数弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="部署参数"
      width="80%"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div v-loading="loading" element-loading-text="加载中...">
        <div class="deployment-item">
          <div class="endpoint-section">
            <div class="endpoint-header">
              <div class="header-item">
                <span class="label">Endpoint ID:</span>
                <span class="value endpoint-id">{{ deploymentDetails?.endpoint_id }}</span>
              </div>
              <div class="header-item">
                <span class="label">URL:</span>
                <span class="value url-value">
                  {{ deploymentDetails?.url }}
                  <template v-if="getEndpointTypeTag(deploymentDetails?.url)">
                    <el-tag :style="{marginLeft: '6px', background: getEndpointTypeTag(deploymentDetails.url).color, color: '#fff', border: 'none'}" size="small" :title="getEndpointTypeTag(deploymentDetails.url).tip">
                      {{ getEndpointTypeTag(deploymentDetails.url).tag }}
                    </el-tag>
                  </template>
                </span>
              </div>
            </div>
            
            <div v-if="deploymentDetails?.serverlessInfo" class="endpoint-content">
              <div class="info-item">
                <span class="label">Model Path:</span>
                <span class="value model-path">{{ deploymentDetails.serverlessInfo.modelPath }}</span>
                <el-button 
                  type="text" 
                  @click="openHuggingFace(deploymentDetails.serverlessInfo.modelPath)"
                  class="view-button"
                >
                  查看
                </el-button>
              </div>
              <div class="info-item">
                <span class="label">Image:</span>
                <span class="value image-value">{{ deploymentDetails.serverlessInfo.image }}</span>
              </div>
              <div class="info-item">
                <span class="label">VARS:</span>
                <div class="vars-display">
                  <div v-for="envVar in deploymentDetails.serverlessInfo.vars"   
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
                  <pre>{{ deploymentDetails.serverlessInfo.args.join('\n') }}</pre>
                  <el-button type="text" class="copy-button" @click="copyText(deploymentDetails.serverlessInfo.args.join('\n'))">复制</el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/composables/toast'
import * as echarts from 'echarts/core'
import { PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { Refresh } from '@element-plus/icons-vue'

// 注册必需的组件
echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  PieChart,
  CanvasRenderer
])

export default {
  name: 'EndpointManager',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const toast = useToast()
    const modelName = ref(route.params.modelName)
    const endpoints = ref([])
    const modelData = ref(null)
    const canModify = ref(false)
    const dialogVisible = ref(false)
    const deploymentDetails = ref(null)
    const loading = ref(false)

    // 从 URL 中取 serverless endpoint 名称
    const extractServerlessEndpoint = (url) => {
      if (!url.includes('runsync.alpha.dev')) return null
      const match = url.match(/https:\/\/(.*?)\.runsync\.alpha\.dev/)
      return match ? match[1] : null
    }

    // 获取 serverless endpoint 详情
    const fetchServerlessEndpoint = async (name) => {
      try {
        const platform = route.query.platform
        const response = await axios.get(`se/${name}${platform ? `?platform=${platform}` : ''}`)
        return response.data.data
      } catch (error) {
        console.error('Error fetching serverless endpoint:', error)
        ElMessage.error(`获取 endpoint 详情失败: ${error.message}`)
        return null
      }
    }

    // 从 args 中提取 model path
    const extractModelPath = (args) => {
      if (!Array.isArray(args)) return ''
      const modelIndex = args.findIndex(arg => arg === '--model')
      if (modelIndex !== -1 && modelIndex + 1 < args.length) {
        return args[modelIndex + 1]
      }
      return ''
    }

    // 显示部署参数
    const showDeploymentInfo = async (endpoint) => {
      dialogVisible.value = true
      deploymentDetails.value = null
      loading.value = true

      try {
        const detail = {
          endpoint_id: endpoint.endpoint_id,
          url: endpoint.url
        }

        const seName = extractServerlessEndpoint(endpoint.url)
        if (seName) {
          const seInfo = await fetchServerlessEndpoint(seName)
          if (seInfo) {
            const container = seInfo.spec.template.spec.containers[0]
            const modelPath = extractModelPath(container.args)
            detail.serverlessInfo = {
              modelPath: modelPath,
              image: container.image,
              args: container.args,
              vars: container.env
            }
          }
        }

        deploymentDetails.value = detail
      } catch (error) {
        console.error('Error processing deployment info:', error)
        ElMessage.error('部署信息发送错误')
      } finally {
        loading.value = false
      }
    }

    const openHuggingFace = (modelPath) => {
      const url = `https://huggingface.co/${modelPath}`
      window.open(url, '_blank')
    }

    const copyText = async (text) => {
      try {
        await navigator.clipboard.writeText(text)
        ElMessage.success('复制成功')
      } catch (error) {
        console.error('Copy failed:', error)
        ElMessage.error('复制失败')
      }
    }

    // 获取 GPU 信息
    const fetchGPUInfo = async (endpoint) => {
      if (!endpoint.gpuInfo) {
        endpoint.gpuInfo = { loading: true }
      }

      try {
        const seName = extractServerlessEndpoint(endpoint.url)
        if (seName) {
          const seInfo = await fetchServerlessEndpoint(seName)
          const gpuRequests = seInfo?.spec?.template?.spec?.containers?.[0]?.resources?.gpuRequests
          
          if (gpuRequests) {
            endpoint.gpuInfo = {
              loading: false,
              model: gpuRequests.models?.[0] || '-',
              number: gpuRequests.number || 0
            }
          } else {
            endpoint.gpuInfo = {
              loading: false,
              model: '-',
              number: 0
            }
          }
        } else {
          endpoint.gpuInfo = {
            loading: false,
            model: '-',
            number: 0
          }
        }
      } catch (error) {
        endpoint.gpuInfo = {
          loading: false,
          error: '获取失败'
        }
      }
    }

    // 添加检用户限方法
    const checkUserPermission = () => {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        try {
          const user = JSON.parse(userStr)
          // admin 户或者具有 admin 角色的用户可以修改
          canModify.value = user.role === 'admin' || user.role === 'operator'
        } catch (error) {
          console.error('Error parsing user data:', error)
          canModify.value = false
        }
      }
    }

    // 改 fetchModelData 方法
    const fetchModelData = async () => {
      try {
        // 从路由中获取 platform 参数
        const platform = route.query.platform
        const response = await axios.get(`/models${platform ? `?platform=${platform}` : ''}`)
        if (response.data?.data?.data) {
          const modelsList = response.data.data.data
          modelData.value = modelsList.find(model => model.model_name === modelName.value)
          
          if (modelData.value) {
            // 先设置基础数据
            endpoints.value = modelData.value.endpoints || []
            
            // 异步加载每个 endpoint 的细信息
            endpoints.value.forEach(endpoint => {
              endpoint.gpuInfo = { loading: true }
              endpoint.worker_count = null
              
              // 异步加载 GPU 息
              fetchGPUInfo(endpoint)
              
              // 异步加载 worker 
              const seName = extractServerlessEndpoint(endpoint.url)
              if (seName) {
                axios.get(`/workers?se=${seName}${platform ? `&platform=${platform}` : ''}`).then(response => {
                  endpoint.worker_count = response.data.data?.filter(worker => 
                    worker.status.state === 'Running' && worker.status.healthy
                  ).length || 0
                }).catch(error => {
                  console.error(`Failed to fetch workers for endpoint ${endpoint.url}:`, error)
                  endpoint.worker_count = 0
                })
              } else {
                endpoint.worker_count = 0
              }
              // endpoint.supported_api_flag = "test" 
            })

            // 检查用户权限
            checkUserPermission()
          } else {
            console.warn('Model not found:', modelName.value)
            toast.show('未找到模型信息', 'error')
          }
        }
      } catch (error) {
        console.error('Error fetching model data:', error)
        toast.show('获取模型信息失败', 'error')
      }
    }

    // 计算属性
    const hasEndpoints = computed(() => endpoints.value.length > 0)
    
    // 判断是否 Serverless 端点
    const isServerless = (endpoint) => {
      return endpoint.url.includes('.runsync.alpha.dev')
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
    const openMonitor = (type) => {
      if (!modelData.value) return

      const platform = route.query.platform
      let overviewUrl = `https://grafana.aicloud.pplabs.tech/d/430e32c1-c956-4f4d-a4e0-8c463180cbe3/ai-cloud-platform-sla-metrics-new?orgId=1&from=now-12h&to=now&orgId=1&var-CustomerUUID=ALL&var-Model=${encodeURIComponent(modelName.value)}&refresh=5s&var-Customer=All`
      if (platform === 'beta') {
        overviewUrl = `https://grafana.aicloud.paigod.work/d/a9131a24-3266-4b18-8aae-b5fc2ec45adc/ai-cloud-platform-sla-metrics?orgId=1&from=now-12h&to=now&orgId=1&var-CustomerUUID=ALL&var-Model=${encodeURIComponent(modelName.value)}&refresh=5s&var-Customer=All`
      }

      const monitorUrls = {
        overview: overviewUrl,
        vllm: modelData.value.inference_engine === 'sglang' ? 
          `https://grafana.aicloud.pplabs.tech/d/ddyp55uq7brpcc/sglang-dashboard?orgId=1&refresh=5s&var-model_id=${encodeURIComponent(getSglangModelId(modelName.value))}` :
          `https://grafana.aicloud.pplabs.tech/d/b281712d-8bff-41ef-9f3f-71ad43c05e9bdd/vllm-per-model?orgId=1&var-model_name=${encodeURIComponent(modelName.value)}`,
        serverless: endpoints.value ? 
          `https://grafana.aicloud.pplabs.tech/d/KOpAVRCIz1/prod-overview?orgId=1${
            endpoints.value
              .filter(endpoint => endpoint.url.includes('.runsync.alpha.dev'))
              .map(endpoint => `&var-endpoint=${endpoint.url.replace('https://', '').replace('http://', '').split('.runsync.alpha.dev')[0]}`)
              .join('')
          }&from=now-24h&to=now&refresh=5s` : null
      }
      
      const url = monitorUrls[type]
      if (url) {
        window.open(url, '_blank')
      }
    }

    // Endpoint 关方法
    const showAddEndpointModal = () => {
      // 重置表单
      newEndpoint.value = {
        endpoint_id: '',
        url: '',
        weight: 1,
        model_id_override: '',
        enable_check_health: false,
        enable_testflight: false,
        supported_api_flag: 'ENDPOINT_SUPPORTED_API_FLAG_ALL',
        skip_testflight_while_add: false // 新增字段，默认不选中
      }
      addEndpointDialogVisible.value = true
    }

    const testEndpoint = async (endpoint) => {
      try {
        const platform = route.query.platform
        const response = await axios.post(`/models/endpoint/test?${platform ? `platform=${platform}` : ''}`, {
          endpoint_url: endpoint.url,
          model_name: modelName.value
        })
        
        if (response.data.success) {
          toast.show('测试成功', 'success')
        } else {
          throw new Error(response.data.message || '测试失败')
        }
      } catch (error) {
        console.error('Error testing endpoint:', error)
        toast.show(error.message || '测试失败', 'error')
      }
    }

    // 改 fetchWorkerMetrics 方法
    const fetchWorkerMetrics = async (seName) => {
      console.log('Fetching metrics for:', seName)
      try {
        const platform = route.query.platform
        const response = await axios.get(`/metrics?endpoint=${seName}${platform ? `&platform=${platform}` : ''}`)
        console.log('Metrics response:', response.data)

        // 更新 worker 的并发数，不影响其他数据
        if (workerDetails.value) {
          workerDetails.value = workerDetails.value.map(worker => {
            const workerName = worker.metadata.name
            const workerMetric = response.data.worker_concurrency?.find(
              m => workerName.includes(m.worker_id) || m.worker_id.includes(workerName)
            )
            
            return {
              ...worker,
              concurrency: workerMetric?.concurrency || 0,
              concurrencyUpdated: worker.concurrency !== workerMetric?.concurrency
            }
          })

          // 清除更新标记
          setTimeout(() => {
            workerDetails.value = workerDetails.value.map(worker => ({
              ...worker,
              concurrencyUpdated: false
            }))
          }, 1000)
        }
      } catch (error) {
        console.error('Error fetching metrics:', error)
        // 错误不影响界面显示
      }
    }

    // 修改 showEndpointDetails 方法
    const showEndpointDetails = async (endpoint) => {
      currentEndpoint.value = endpoint
      detailsDialogVisible.value = true
      detailsLoading.value = true

      try {
        const seName = extractServerlessEndpoint(endpoint.url)
        if (!seName) {
          throw new Error('法解析 endpoint 名称')
        }

        const platform = route.query.platform
        // 先获取基础数据
        const [sapResponse, workerResponse] = await Promise.all([
          axios.get(`/sap?se=${seName}${platform ? `&platform=${platform}` : ''}`),
          axios.get(`/workers?se=${seName}${platform ? `&platform=${platform}` : ''}`)
        ])

        sapDetails.value = sapResponse.data.data[0]
        workerDetails.value = workerResponse.data.data

        // 关闭加载状态，示基础数据
        detailsLoading.value = false

        // 异步获取并数数据
        fetchWorkerMetrics(seName).catch(error => {
          console.error('Error fetching initial metrics:', error)
          // 并数据获取失败不影响主界面
        })

        // 设置定时器定期刷新并发数
        const metricsTimer = setInterval(() => {
          if (detailsDialogVisible.value) {
            fetchWorkerMetrics(seName).catch(error => {
              console.error('Error refreshing metrics:', error)
              // 并发数据刷新失败不影响主界面
            })
          } else {
            clearInterval(metricsTimer)
          }
        }, 5000)

        // 监对话框关闭
        watch(detailsDialogVisible, (newVal) => {
          if (!newVal && metricsTimer) {
            clearInterval(metricsTimer)
          }
        })

      } catch (error) {
        console.error('Error fetching details:', error)
        ElMessage.error(error.message || '获取详情失败')
        detailsLoading.value = false
      }
    }

    // 在 setup 中添加 openGrafanaLink 方法
    const openGrafanaLink = (endpoint) => {
      if (!endpoint || !endpoint.url) {
        ElMessage.warning('无法获取 endpoint URL')
        return
      }

      const platform = route.query.platform
      const deploymentName = getEndpointName(endpoint.url)
      let grafanaUrl = ''

      // 1. 2.0 (sls2) endpoint logic
      if (endpoint.url.includes('sls2.alpha.ai')) {
        grafanaUrl = `https://grafana.aicloud.pplabs.tech/d/cegvuwvbnppfkb/gateway-overview?orgId=1&refresh=5s&var-env=zjxCFBjVz&var-deployment=${deploymentName}&var-instance=All`
      }
      // 2. Fusion endpoint logic
      else if (endpoint.url.includes('/fusion/v1/')) {
        if (platform === 'beta') {
          // beta Fusion URL
          grafanaUrl = `https://grafana.aicloud.paigod.work/d/bejkkcdnlbo5ca/fusion-gateway-overview?orgId=1&refresh=10s&var-env=aI7J0Qa4k&var-model=${deploymentName}&var-provider=All&from=now-6h&to=now`
        } else {
          // alpha Fusion URL
          grafanaUrl = `https://grafana.aicloud.pplabs.tech/d/bejkkcdnlbo5ca/fusion-gateway-overview?orgId=1&refresh=5s&var-env=zjxCFBjVz&var-model=${deploymentName}&var-provider=All&from=now-6h&to=now`
        }
      } 
      // 3. Original Serverless endpoint logic
      else if (endpoint.url.includes('.runsync.alpha.dev')) {
        const seName = endpoint.url.split('.runsync.alpha.dev')[0].replace('https://', '')
        grafanaUrl = `https://grafana.aicloud.pplabs.tech/d/KOpAVRCIz1/prod-overview?var-endpoint=${seName}`
      }

      if (grafanaUrl) {
        window.open(grafanaUrl, '_blank')
      } else {
        ElMessage.warning('此 endpoint 类型不支持自动跳转面板')
      }
    }

    const toggleEndpointStatus = async (endpoint) => {
      const action = endpoint.status === 'ENDPOINT_STATUS_SERVING' ? 'stop' : 'start'
      const actionText = action === 'stop' ? '停用' : '启用'
      
      try {
        // 显示确认对话框
        await ElMessageBox.confirm(
          `确定要${actionText}该 endpoint ？`,
          '确认操作',
          {
            confirmButtonText: actionText,
            cancelButtonText: '取消',
            type: action === 'stop' ? 'warning' : 'info',
            confirmButtonClass: action === 'stop' ? 'el-button--danger' : 'el-button--primary'
          }
        )

        // 用户确认，显示加载提示
        ElMessage({
          message: `正在${actionText}...`,
          type: 'info',
          duration: 1000
        })

        const platform = route.query.platform
        const response = await axios.post(`endpoints/${action}?${platform ? `platform=${platform}` : ''}`, {
          endpoint_id: endpoint.endpoint_id
        })

        if (response.data.code === 200) {
          ElMessage({
            message: `${actionText}成功`,
            type: 'success',
            duration: 2000
          })

          // 操作成功后刷新数据
          await fetchModelData()
        } else {
          throw new Error(response.data.message || `${actionText}失败`)
        }
      } catch (error) {
        if (error === 'cancel') {
          // 用户取消操作，不需要提示
          return
        }
        
        console.error(`Error ${action}ing endpoint:`, error)
        ElMessage({
          message: error.response?.data?.message || error.message || `${actionText}失败`,
          type: 'error',
          duration: 3000
        })
        
        // 操作失败时刷新数据确保显示确的状态
        await fetchModelData()
      }
    }

    const updateEndpoint = async (endpoint) => {
      try {
        // 显示载提示
        ElMessage({
          message: '正在更新...',
          type: 'info',
          duration: 1000
        })

        const platform = route.query.platform
        const response = await axios.put(`endpoints/${endpoint.endpoint_id}?${platform ? `platform=${platform}` : ''}`, {
          model_name: modelName.value,
          weight: endpoint.weight,
          enable_check_health: endpoint.enable_check_health,
          enable_testflight: endpoint.enable_testflight,
          supported_api_flag: endpoint.supported_api_flag
        })

        if (response.data.code === 200) {
          // 更新成功提示
          ElMessage({
            message: '更新成功',
            type: 'success',
            duration: 2000
          })
        } else {
          throw new Error(response.data.message || '更新失败')
        }
      } catch (error) {
        console.error('Error updating endpoint:', error)
        // 更新失败提示
        ElMessage({
          message: error.response?.data?.message || error.message || '更新失败',
          type: 'error',
          duration: 3000
        })
        
        // 更新失败时回滚数据
        await fetchModelData()
      }
    }

    const weightChart = ref(null)
    const chartInstance = ref(null)

    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEEAD']

    const getEndpointColor = (endpoint) => {
      const index = endpoints.value.indexOf(endpoint)
      return colors[index % colors.length]
    }

    const initWeightChart = () => {
      if (!chartInstance.value) {
        const chartDom = weightChart.value
        if (chartDom) {
          chartInstance.value = echarts.init(chartDom)
        }
      }

      if (chartInstance.value) {
        const option = {
          tooltip: {
            trigger: 'item',
            formatter: (params) => {
              const endpoint = endpoints.value.find(e => e.endpoint_id === params.name)
              const endpointName = getEndpointName(endpoint?.url)
              return `
                <div style="padding: 12px; max-width: 300px;">
                  <div style="font-weight: bold; font-size: 14px; margin-bottom: 8px;">${params.name}</div>
                  <div style="font-size: 12px; color: #666; margin-bottom: 8px;">Endpoint: ${endpointName}</div>
                  <div style="font-size: 13px; margin-bottom: 8px;">权重: ${params.value} (${params.percent}%)</div>
                </div>
              `
            }
          },
          legend: {
            show: false  // 隐藏默认图例，使用自定义图例
          },
          series: [
            {
              type: 'pie',
              radius: ['40%', '70%'],  // 环图
              center: ['50%', '50%'],
              avoidLabelOverlap: true,
              itemStyle: {
                borderRadius: 4,
                borderColor: '#fff',
                borderWidth: 2
              },
              label: {
                show: true,
                position: 'outside',
                formatter: (params) => {
                  const endpoint = endpoints.value.find(e => e.endpoint_id === params.name)
                  const endpointName = getEndpointName(endpoint?.url)
                  return `${endpointName}\n${params.percent}%`
                },
                fontSize: 14,
                lineHeight: 20
              },
              labelLine: {
                length: 20,
                length2: 30,
                smooth: true
              },
              data: endpoints.value.map((endpoint, index) => ({
                name: endpoint.endpoint_id,
                value: endpoint.weight,
                itemStyle: {
                  color: colors[index % colors.length]
                }
              }))
            }
          ]
        }
        
        chartInstance.value.setOption(option)
      }
    }

    // 在组件卸载时销毁图表实例
    onUnmounted(() => {
      if (chartInstance.value) {
        chartInstance.value.dispose()
      }
    })

    // 监听窗口小变化，调整图表大小
    window.addEventListener('resize', () => {
      if (chartInstance.value) {
        chartInstance.value.resize()
      }
    })

    onMounted(() => {
      // 在组件载时检权限
      checkUserPermission()
      fetchModelData().then(() => {
        initWeightChart()
      })
    })

    watch(modelData, () => {
      if (chartInstance.value) {
        initWeightChart()
      }
    })

    // 添加新的响应式量
    const addEndpointDialogVisible = ref(false)
    const newEndpoint = ref({
      endpoint_id: '',
      url: '',
      weight: 1,
      model_id_override: '',
      enable_check_health: false,
      enable_testflight: false,
      supported_api_flag: 'ENDPOINT_SUPPORTED_API_FLAG_ALL',
      skip_testflight_while_add: false // 新增字段，默认不选中
    })

    // 测试新的 endpoint
    const testNewEndpoint = async () => {
      if (!newEndpoint.value.url) {
        ElMessage.warning('请先输入 URL')
        return
      }

      try {
        ElMessage({
          message: '正在测试...',
          type: 'info',
          duration: 1000
        })

        const response = await axios.post('/models/endpoint/test', {
          endpoint_url: newEndpoint.value.url,
          model_name: modelName.value
        })

        if (response.data.success) {
          ElMessage({
            message: '测试成功',
            type: 'success'
          })
        } else {
          throw new Error(response.data.message || '测试失败')
        }
      } catch (error) {
        console.error('Error testing endpoint:', error)
        ElMessage({
          message: error.response?.data?.message || error.message || '测试失败',
          type: 'error'
        })
      }
    }

    // 添加新的 endpoint
    const addEndpoint = async () => {
      try {
        // 基础验证
        if (!newEndpoint.value.endpoint_id || !newEndpoint.value.url) {
          ElMessage.warning('请填写必要的信息')
          return
        }

        ElMessage({
          message: '正在添加...',
          type: 'info',
          duration: 1000
        })

        // 准备请求数据，platform 参数移到 URL 查询参数中
        const platform = route.query.platform
        const requestData = {
          ...newEndpoint.value,
          model_name: modelName.value,
          skip_testflight_while_add: newEndpoint.value.skip_testflight_while_add // 新增参数
        }

        const response = await axios.post(`/endpoints/create${platform ? `?platform=${platform}` : ''}`, requestData)

        if (response.data.code === 200) {
          ElMessage({
            message: '添加成功',
            type: 'success'
          })
          addEndpointDialogVisible.value = false
          await fetchModelData() // 刷新数据
        } else {
          throw new Error(response.data.message || '添加失败')
        }
      } catch (error) {
        console.error('Error adding endpoint:', error)
        ElMessage({
          message: error.response?.data?.message || error.message || '添加失败',
          type: 'error'
        })
      }
    }

    // setup 中加
    const detailsDialogVisible = ref(false)
    const sapFullDetailsVisible = ref(false)
    const detailsLoading = ref(false)
    const currentEndpoint = ref(null)
    const sapDetails = ref(null)
    const workerDetails = ref([])
    const rawInstanceDetailVisible = ref(false)

    // 计算正在服的 worker 数量
    const servingWorkerCount = computed(() => {
      return workerDetails.value.filter(w => 
        w.status.state === 'Running' && w.status.healthy
      ).length
    })

    // 取状态类型
    const getStateType = (state) => {
      switch (state) {
        case 'Running':
          return 'success'
        case 'Pending':
          return 'warning'
        default:
          return 'info'
      }
    }

    // 显示完整 SAP 配置
    const showFullSapDetails = () => {
      sapFullDetailsVisible.value = true
    }

    const fetchEndpoints = async () => {
      loading.value = true
      try {
        const response = await axios.get(`/endpoints?model_name=${modelName}`)
        if (response.data.data) {
          // 获取每个 endpoint 的 worker 数
          const endpointsWithWorkers = await Promise.all(
            response.data.data.map(async (endpoint) => {
              try {
                const workerResponse = await axios.get(`/workers?endpoint=${endpoint.name}`)
                return {
                  ...endpoint,
                  worker_count: workerResponse.data.data?.length || 0
                }
              } catch (error) {
                console.error(`Failed to fetch workers for endpoint ${endpoint.name}:`, error)
                return {
                  ...endpoint,
                  worker_count: 0
                }
              }
            })
          )
          endpoints.value = endpointsWithWorkers
        }
      } catch (error) {
        console.error('Failed to fetch endpoints:', error)
        ElMessage.error('获取 endpoints 列表失败')
      } finally {
        loading.value = false
      }
    }

    // 实例详情话框
    const instanceDetailVisible = ref(false)
    const instanceDetailLoading = ref(false)
    const instanceDetail = ref(null)

    // 显示例详情
    const showInstanceDetail = async (instanceId, clusterId) => {
      instanceDetailVisible.value = true
      instanceDetailLoading.value = true
      try {
        const response = await axios.get('/nexusclusters/instances/detail', {
          params: {
            cluster_index: clusterId,
            instance_id: instanceId
          }
        })
        instanceDetail.value = response.data.data
      } catch (error) {
        console.error('Failed to fetch instance detail:', error)
        ElMessage.error('获取例详情失败')
      } finally {
        instanceDetailLoading.value = false
      }
    }

    // 格式化日期
    const formatDate = (timestamp) => {
      if (!timestamp) return '-'
      return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
    }

    // 格式化字节为GB
    const formatGb = (bytes) => {
      if (!bytes) return 0
      return (bytes / (1024 * 1024 * 1024)).toFixed(1)
    }

    const copyToClipboard = async (text) => {
      try {
        // 创建临时输入框
        const textarea = document.createElement('textarea')
        textarea.value = text
        document.body.appendChild(textarea)
        
        // 选择并复制
        textarea.select()
        document.execCommand('copy')
        
        // 移除临时输入框
        document.body.removeChild(textarea)
        
        ElMessage.success('复制成功')
      } catch (error) {
        console.error('Failed to copy:', error)
        ElMessage.error('复制失败')
      }
    }

    // 例原始信息话框
    const showRawInstanceDetail = () => {
      rawInstanceDetailVisible.value = true
    }

    const getConcurrencyTarget = (spec) => {
      if (!spec?.metrics) return null
      
      const concurrencyMetric = spec.metrics.find(metric => 
        metric.type === 'Concurrency' && 
        metric.resource?.target?.type === 'AverageValue'
      )
      
      return concurrencyMetric?.resource?.target?.averageValue
    }

    const workerListRefreshing = ref(false)
    
    const refreshWorkerList = async () => {
      if (!currentEndpoint.value) return
      
      workerListRefreshing.value = true
      try {
        const seName = extractServerlessEndpoint(currentEndpoint.value.url)
        if (!seName) return
        
        // 重新获取 worker 列表
        const workersResponse = await axios.get(`/workers?se=${seName}`)
        workerDetails.value = workersResponse.data.data || []
        
        // 更新运行中的 worker 数量
        servingWorkerCount.value = workerDetails.value.filter(
          worker => worker.status.state === 'Running' && worker.status.healthy
        ).length
        
        ElMessage.success('Worker 列表已刷新')
      } catch (error) {
        console.error('Error refreshing worker list:', error)
        ElMessage.error('刷新 Worker 表失败')
      } finally {
        workerListRefreshing.value = false
      }
    }

    // 添加新的响应式变量
    const isEditingSap = ref(false)
    const editingSapValues = ref({
      minReplicas: 0,
      maxReplicas: 0,
      concurrencyPerWorker: 0,
      scaleUpWindow: 0,
      scaleDownWindow: 0
    })

    // 开始编辑 SAP 配置
    const startEditSap = () => {
      isEditingSap.value = true
      editingSapValues.value = {
        minReplicas: sapDetails.value.spec.minReplicas,
        maxReplicas: sapDetails.value.spec.maxReplicas,
        concurrencyPerWorker: getConcurrencyTarget(sapDetails.value.spec) || 0,
        scaleUpWindow: sapDetails.value.spec.behavior?.scaleUp?.stabilizationWindowSeconds || 0,
        scaleDownWindow: sapDetails.value.spec.behavior?.scaleDown?.stabilizationWindowSeconds || 0
      }
    }

    // 取消编辑 SAP 配置
    const cancelEditSap = () => {
      isEditingSap.value = false
    }

    // 保存 SAP 配置
    const saveSapChanges = async () => {
      try {
        const seName = extractServerlessEndpoint(currentEndpoint.value.url)
        if (!seName) {
          throw new Error('无法解析endpoint名称')
        }

        // 准备修改内容的展示
        const changes = []
        if (editingSapValues.value.minReplicas !== sapDetails.value.spec.minReplicas) {
          changes.push(`最小副本数: ${sapDetails.value.spec.minReplicas} -> ${editingSapValues.value.minReplicas}`)
        }
        if (editingSapValues.value.maxReplicas !== sapDetails.value.spec.maxReplicas) {
          changes.push(`最大副本数: ${sapDetails.value.spec.maxReplicas} -> ${editingSapValues.value.maxReplicas}`)
        }
        if (editingSapValues.value.concurrencyPerWorker !== getConcurrencyTarget(sapDetails.value.spec)) {
          changes.push(`触发扩容并发: ${getConcurrencyTarget(sapDetails.value.spec) || 0} -> ${editingSapValues.value.concurrencyPerWorker}`)
        }
        if (editingSapValues.value.scaleUpWindow !== sapDetails.value.spec.behavior?.scaleUp?.stabilizationWindowSeconds) {
          changes.push(`扩容窗口时间: ${sapDetails.value.spec.behavior?.scaleUp?.stabilizationWindowSeconds || 0} -> ${editingSapValues.value.scaleUpWindow}`)
        }
        if (editingSapValues.value.scaleDownWindow !== sapDetails.value.spec.behavior?.scaleDown?.stabilizationWindowSeconds) {
          changes.push(`缩容窗口时间: ${sapDetails.value.spec.behavior?.scaleDown?.stabilizationWindowSeconds || 0} -> ${editingSapValues.value.scaleDownWindow}`)
        }

        if (changes.length === 0) {
          ElMessage.info('没有值被修改')
          isEditingSap.value = false
          return
        }

        // 显示确认对话框
        const confirmContent = `
          <div class="sap-confirm-dialog">
            <div class="confirm-title">确认以下修改：</div>
            <div class="changes-list">
              ${changes.map(change => `
                <div class="change-item">
                  <i class="el-icon-arrow-right"></i>
                  <span>${change}</span>
                </div>
              `).join('')}
            </div>
          </div>
          <style>
            .sap-confirm-dialog {
              padding: 16px 0;
            }
            .confirm-title {
              font-size: 16px;
              font-weight: 600;
              color: #1f2937;
              margin-bottom: 16px;
            }
            .changes-list {
              background: #f3f4f6;
              border-radius: 8px;
              padding: 12px;
              border: 1px solid #e5e7eb;
            }
            .change-item {
              display: flex;
              align-items: center;
              padding: 10px;
              border-bottom: 1px solid #e5e7eb;
              font-size: 14px;
              background: white;
              margin-bottom: 4px;
              border-radius: 4px;
            }
            .change-item:last-child {
              margin-bottom: 0;
              border-bottom: none;
            }
            .change-item i {
              color: #3b82f6;
              margin-right: 10px;
              font-size: 14px;
            }
            .change-item span {
              color: #374151;
              font-weight: 500;
            }
          </style>
        `
        await ElMessageBox.confirm(
          confirmContent,
          '确认修改',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
            dangerouslyUseHTMLString: true,
            customClass: 'sap-confirm-box'
          }
        )

        // 显示加载提示
        ElMessage({
          message: '正在更新...',
          type: 'info',
          duration: 1000
        })

        // 准备请求数据，确保发送所有字段，并转换为正确的类型
        const requestData = {
          minReplicas: parseInt(editingSapValues.value.minReplicas) || 0,
          maxReplicas: parseInt(editingSapValues.value.maxReplicas) || 0,
          concurrencyPerWorker: parseInt(editingSapValues.value.concurrencyPerWorker) || 0,
          scaleUpWindow: parseInt(editingSapValues.value.scaleUpWindow) || 0,
          scaleDownWindow: parseInt(editingSapValues.value.scaleDownWindow) || 0
        }

        // 发送更新请求
        const response = await axios.post(`/sap/${seName}?endpoint=${seName}`, requestData)

        if (response.data.code === 200) {
          ElMessage({
            message: '更新成功',
            type: 'success',
            duration: 2000
          })
          isEditingSap.value = false
          
          // 重新获取SAP信息
          const sapResponse = await axios.get(`/sap?se=${seName}`)
          sapDetails.value = sapResponse.data.data[0]
        } else {
          throw new Error(response.data.message || '更新失败')
        }
      } catch (error) {
        if (error === 'cancel') return
        
        console.error('Error updating SAP configuration:', error)
        ElMessage({
          message: error.response?.data?.message || error.message || '更新失败',
          type: 'error',
          duration: 3000
        })
      }
    }

    // 在 setup 中添加 getEndpointName 方法
    const getEndpointName = (url) => {
      if (!url) return '-'

      try {
        // 2.0 (sls2) format: https://{deployment}.us-01.sls2.alpha.ai...
        let match = url.match(/https?:\/\/([^.]+)\.us-01\.sls2\.alpha\.ai/)
        if (match && match[1]) {
          return match[1]
        }

        // Fusion format: http://thirdparty-api-adapter/fusion/v1/{modelname}
        if (url.includes('/fusion/v1/')) {
          const parts = url.split('/fusion/v1/')
          if (parts.length > 1 && parts[1]) {
            return parts[1]
          }
        }

        // Original serverless format: https://{endpointName}.runsync.alpha.dev
        match = url.match(/https?:\/\/([^.]+)\.runsync\.alpha\.dev/)
        if (match && match[1]) {
          return match[1]
        }
      } catch (e) {
        console.error("Error parsing URL:", url, e);
        return url; // return original url on error
      }
      
      return url // Fallback to original URL
    }

    // 在 setup 中添加 goToServerless 方法
    const goToServerless = (endpointName, endpointUrl = '') => {
      // 判断是否是 2.0 版本的 endpoint
      if (endpointUrl.includes('sls2.alpha.ai')) {
        // 2.0 版本跳转到 nebula 管理页面的 deployment tab
        router.push({
          path: '/nebula/deployment',
          query: { name: endpointName }
        })
      } else {
        // 1.0 版本跳转到 serverless 管理页面
        router.push({
          path: '/serverless',
          query: { search: endpointName }
        })
      }
    }

    // 在 setup 函数中添加 openWorkerGrafana 方法
    const openWorkerGrafana = (worker) => {
      try {
        // 从 worker container 的环境变量中获取 model_name
        const container = worker.spec.containers?.[0]
        let modelNameParam = modelName.value

        if (container?.env) {
          // 优先从 METRIC_MODEL 获取，如果没有则从 MODEL_ID 获取
          const metricModel = container.env.find(env => env.name === 'METRIC_MODEL')?.value
          const modelId = container.env.find(env => env.name === 'MODEL_ID')?.value
          modelNameParam = metricModel || modelId || modelName.value
        }

        // 从 worker 名称中获取设备 ID
        const deviceId = worker.metadata.name.split('-').pop()

        // 构建 Grafana URL
        const url = `https://grafana.aicloud.pplabs.tech/d/b281712d-8bff-41ef-9f3f-71ad43c05e9c/vllm-per-instance?orgId=1&refresh=5s&var-model_name=${encodeURIComponent(modelNameParam)}&var-provider=All&var-device_id=${deviceId}`
        
        window.open(url, '_blank')
      } catch (error) {
        console.error('Error opening worker Grafana:', error)
        ElMessage.error('打开监控面板失败')
      }
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

    return {
      modelName,
      endpoints,
      modelData,
      canModify,
      hasEndpoints,
      isServerless,
      openMonitor,
      showAddEndpointModal,
      testEndpoint,
      showEndpointDetails,
      openGrafanaLink,
      toggleEndpointStatus,
      updateEndpoint,
      weightChart,
      getEndpointColor,
      addEndpointDialogVisible,
      newEndpoint,
      testNewEndpoint,
      addEndpoint,
      detailsDialogVisible,
      sapFullDetailsVisible,
      detailsLoading,
      currentEndpoint,
      sapDetails,
      workerDetails,
      servingWorkerCount,
      getStateType,
      showFullSapDetails,
      fetchEndpoints,
      instanceDetailVisible,
      copyText,
      isEditingSap,
      editingSapValues,
      startEditSap,
      saveSapChanges,
      cancelEditSap,
      getEndpointName,
      goToServerless,
      openWorkerGrafana,
      getEndpointTypeTag
    }
  }
}
</script>

<style scoped>
.endpoint-manager {
  padding: 20px;
}

/* 添加新的样式规则 */
:deep(.el-dialog) {
  --el-dialog-padding-primary: 16px;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}

:deep(.el-form-item__label) {
  font-size: 14px;
}

:deep(.el-input__inner) {
  font-size: 14px;
}

.options-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 10px;
  margin-bottom: 18px;
}

.options-row .checkbox-item {
  margin-bottom: 0;
  flex-shrink: 0;
}

.options-row .select-item {
  margin-bottom: 0;
  margin-left: auto;
}

.api-flag-select {
  width: 140px;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.card-header {
  background: #f8f9fa;
  border-bottom: 1px solid #eee;
  padding: 1rem;
}

.card-title {
  font-size: 1.25rem;
}

.btn-group {
  display: flex;
  gap: 4px;
}

.table {
  margin-bottom: 0;
}

.table th {
  background: #f8f9fa;
  font-weight: 600;
}

.table td {
  vertical-align: middle;
}

.form-check {
  margin: 0;
}

.gpu-info {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 24px;
}

.gpu-spec {
  font-family: monospace;
  background: #f0f9ff;
  color: #0969da;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 13px;
  border: 1px solid #e5e7eb;
}

.text-danger {
  color: #dc3545;
  font-size: 12px;
}

.weight-distribution {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-top: 30px;
}

.weight-distribution h5 {
  margin-bottom: 20px;
  padding-left: 20px;
}

.chart-container {
  height: 400px;
  padding: 20px;
}

.legend-container {
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  height: 100%;
}

.legend-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 15px;
  padding: 10px;
  background: white;
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.legend-color {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  margin-right: 10px;
  margin-top: 4px;
}

.legend-info {
  flex: 1;
}

.legend-name {
  font-weight: 600;
  margin-bottom: 4px;
}

.legend-url {
  font-size: 12px;
  word-break: break-all;
  margin-bottom: 4px;
}

.legend-weight {
  font-size: 13px;
  color: #666;
}

.action-btn {
  min-width: 68px;
  height: 32px;
  padding: 4px 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 1px;
}

.bi-graph-up {
  font-size: 14px;
}

/* 响应式调 */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .btn-group {
    width: 100%;
    justify-content: center;
  }
}

.table th,
.table td {
  vertical-align: middle;
}

.weight-input {
  max-width: 80px;
  min-width: 75px; /* 调整为75px，保证五位数显示 */
  margin: 0 auto;
  text-align: center;
}

.form-check {
  margin: 0;
  padding: 0;
}

.form-check-input {
  margin: 0;
  float: none;
}

.btn-group {
  display: flex;
  gap: 4px;
  flex-wrap: nowrap;
}

.action-buttons {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.action-btn {
  margin: 0 !important;
  padding: 4px 12px !important;
}

:deep(.el-button--small) {
  padding: 4px 12px;
  font-size: 13px;
  height: auto;
  line-height: 1.5;
}

:deep(.el-button + .el-button) {
  margin-left: 4px;
}

:deep(.el-button--primary) {
  --el-button-bg-color: #409eff;
  --el-button-border-color: #409eff;
  --el-button-hover-bg-color: #66b1ff;
  --el-button-hover-border-color: #66b1ff;
}

:deep(.el-button--info) {
  --el-button-bg-color: #909399;
  --el-button-border-color: #909399;
  --el-button-hover-bg-color: #a6a9ad;
  --el-button-hover-border-color: #a6a9ad;
}

:deep(.el-button--success) {
  --el-button-bg-color: #67c23a;
  --el-button-border-color: #67c23a;
  --el-button-hover-bg-color: #85ce61;
  --el-button-hover-border-color: #85ce61;
}

:deep(.el-button--danger) {
  --el-button-bg-color: #f56c6c;
  --el-button-border-color: #f56c6c;
  --el-button-hover-bg-color: #f78989;
  --el-button-hover-border-color: #f78989;
}

:deep(.el-button--warning) {
  --el-button-bg-color: #e6a23c;
  --el-button-border-color: #e6a23c;
  --el-button-hover-bg-color: #ebb563;
  --el-button-hover-border-color: #ebb563;
}

.auto-scaling-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #f8fafc;
  border-radius: 8px;
}

.scaling-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: white;
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.scaling-item .label {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
  text-align: center;
  width: 100%;
}

.scaling-item .value {
  font-size: 20px;
  font-weight: 600;
  color: #3b82f6;
  text-align: center;
  width: 100%;
}

.details-section {
  background: #fff;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 8px;
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
    gap: 6px;
  }

  i {
    font-size: 16px;
    color: #1d4ed8;
    background: #e0e7ff;
    padding: 6px;
    border-radius: 8px;
  }

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    position: relative;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(500px, 1fr));
  gap: 8px;
  padding: 8px;
  background: #fafafa;
  border-radius: 8px;
}

.full-width {
  grid-column: 1 / -1;
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

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.ml-2 {
  margin-left: 8px;
}

.dialog-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 700;
  color: #ffffff;
}

.ssh-password {
  margin-left: 12px;
  padding-left: 12px;
  border-left: 1px solid #4ade80;
  color: #ffffff;
  font-family: monospace;
  font-size: 14px;
}

.log-link {
  text-decoration: none;
  color: #3b82f6;
  font-family: monospace;
  font-size: 12px;
  word-break: break-all;
  
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

.info-item.double {
  display: flex;
  justify-content: space-between;
  padding: 6px 8px;
  
  .info-pair {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    
    &:first-child {
      margin-right: 16px;
    }
  }
}

.value {
  color: #111827;
  font-weight: 500;
  flex: 1;
  
  &.code {
    font-family: monospace;
    background: #f8fafc;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 13px;
    border: 1px solid #e5e7eb;
    max-width: none;  /* 移除最大宽度限制 */
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
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

.ssh-info-row {
  display: flex;
  gap: 24px;
  align-items: center;
  width: 100%;
}

.ssh-item {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.label {
  white-space: nowrap;
  color: #666;
  font-size: 14px;
}

.code-value {
  font-family: monospace;
  background: #f8fafc;
  padding: 6px 12px;
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

/* 部署参数弹窗样式 */
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

.action-buttons {
  display: flex;
  gap: 4px;
  align-items: center;
  justify-content: center;
}

.deployment-btn {
  min-width: 88px;  /* 调整为与其他按钮相同的宽度 */
  height: 32px;     /* 调整为与其他按钮相同的高度 */
}

.table th:last-child,
.table td:last-child {
  width: 420px;
  max-width: 420px;
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.action-btn {
  margin: 0 !important;
  padding: 4px 12px !important;
}

:deep(.el-button--small) {
  padding: 4px 12px;
  font-size: 13px;
  height: auto;
  line-height: 1.5;
}

:deep(.el-button + .el-button) {
  margin-left: 4px;
}

:deep(.el-button--primary) {
  --el-button-bg-color: #409eff;
  --el-button-border-color: #409eff;
  --el-button-hover-bg-color: #66b1ff;
  --el-button-hover-border-color: #66b1ff;
}

:deep(.el-button--info) {
  --el-button-bg-color: #909399;
  --el-button-border-color: #909399;
  --el-button-hover-bg-color: #a6a9ad;
  --el-button-hover-border-color: #a6a9ad;
}

:deep(.el-button--success) {
  --el-button-bg-color: #67c23a;
  --el-button-border-color: #67c23a;
  --el-button-hover-bg-color: #85ce61;
  --el-button-hover-border-color: #85ce61;
}

:deep(.el-button--danger) {
  --el-button-bg-color: #f56c6c;
  --el-button-border-color: #f56c6c;
  --el-button-hover-bg-color: #f78989;
  --el-button-hover-border-color: #f78989;
}

:deep(.el-button--warning) {
  --el-button-bg-color: #e6a23c;
  --el-button-border-color: #e6a23c;
  --el-button-hover-bg-color: #ebb563;
  --el-button-hover-border-color: #ebb563;
}

/* 自定义数字输入框样式 */
:deep(.custom-number-input) {
  width: auto;
}

:deep(.custom-number-input .el-input__inner) {
  font-size: 20px;
  font-weight: 600;
  color: #3b82f6;
  text-align: center;
  height: auto;
  padding: 0;
  border: none;
  background: transparent;
}

:deep(.custom-number-input .el-input-number__decrease),
:deep(.custom-number-input .el-input-number__increase) {
  background: transparent;
  border: none;
  color: #3b82f6;
}

:deep(.custom-number-input:hover .el-input__inner) {
  border: none;
}

:deep(.custom-number-input .el-input__inner:focus) {
  border: none;
  box-shadow: none;
}

.legend-url.clickable {
  cursor: pointer;
  color: #3b82f6;
  transition: all 0.2s;
  
  &:hover {
    color: #2563eb;
    text-decoration: underline;
  }
}

.status-container {
  display: flex;
  gap: 4px;
  justify-content: center;
  align-items: center;
}

.health-tag {
  margin-left: 4px;
}

:deep(.el-tag--small) {
  height: 20px;
  padding: 0 6px;
  font-size: 11px;
}

.api-flag-select {
  width: 150px;
}

/* el-dialog 小屏幕适配，内容区滚动，底部按钮固定可见 */
:deep(.el-dialog) {
  max-width: 95vw !important;
  width: 95vw !important;
  max-height: 80vh !important;
  display: flex;
  flex-direction: column;
}
:deep(.el-dialog__body) {
  overflow-y: auto;
  max-height: 50vh;
}
@media (max-width: 768px) {
  :deep(.el-dialog) {
    max-width: 98vw !important;
    width: 98vw !important;
    max-height: 90vh !important;
  }
  :deep(.el-dialog__body) {
    max-height: 60vh;
  }
}

/* 标记说明样式 */
.tags-legend {
  border-top: 1px solid #e5e7eb;
  padding-top: 16px;
}

.tags-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.tag-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.tag-symbol {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: #f3f4f6;
  color: #374151;
  font-weight: bold;
  font-size: 12px;
}

.tag-text {
  color: #6b7280;
}

.special-tag {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 4px;
}

.rong-tag {
  background-color: #f59e42;
  color: #fff;
}

.version-tag {
  background-color: #2563eb;
  color: #fff;
}
</style>