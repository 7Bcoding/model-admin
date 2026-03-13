<template>
  <div class="model-api-container">
    <div class="layout-container">
      <!-- 左侧编辑区域 -->
      <div class="edit-section">
        <a-card title="模型API配置">
          <!-- 基本信息 -->
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="名称" name="name" 
                :help="'建议使用小写字母、数字和连字符'">
                <a-input v-model:value="formState.name" placeholder="请输入名称" />
              </a-form-item>
            </a-col>
            <a-col :span="4">
              <a-form-item label="区域" name="region">
                <a-select v-model:value="selectedRegion" placeholder="选择区域" @change="handleRegionChange">
                  <a-select-option v-for="(region, key) in regions" :key="key" :value="key">
                    {{ region.name }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="镜像" name="image">
                <a-input v-model:value="formState.image" placeholder="请输入镜像" />
              </a-form-item>
            </a-col>
          </a-row>

          <!-- Args Config -->
          <a-card title="启动参数配置" size="small" class="mb-4">
            <a-textarea
              v-model:value="argsText"
              :rows="8"
              placeholder="请输入启动参数，每行一个参数"
              @change="handleArgsChange"
            />
            <a-form-item label="单Worker最大并发数" name="max_concurrent_per_worker"
              :help="'确保该值小于启动参数中的--max-num-seqs'"
              :rules="[{ required: true, message: '请输入最大并发数' }]"
              style="margin-top: 16px; margin-bottom: 0;">
              <a-input-number v-model:value="formState.max_concurrent_per_worker" :min="1" :max="32" style="width: 200px" />
            </a-form-item>
          </a-card>

          <!-- Resources -->
          <a-card title="资源配置" size="small" class="mb-4">
            <a-row :gutter="4" type="flex" align="middle">
              <a-col :span="14">
                <a-form-item label="GPU型号" :name="['resources', 'gpu', 'model']" :rules="[{ required: true, message: '请选择GPU型号' }]" :validate-status="formState.resources.gpu.model ? 'success' : undefined" class="no-margin">
                  <a-select v-model:value="formState.resources.gpu.model" placeholder="请选择GPU型号" style="width: 100%">
                    <a-select-option v-for="gpu in gpuModels" :key="gpu.value" :value="gpu.value">
                      {{ gpu.label }} ({{ gpu.memory }})
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="5">
                <a-form-item label="GPU数量" :name="['resources', 'gpu', 'count']" 
                  :rules="[{ required: true, message: '请输入GPU数量' }, { type: 'number', min: 1, max: 8, message: 'GPU数量必须在1-8之间' }]"
                  class="no-margin"
                >
                  <a-input-number 
                    v-model:value="formState.resources.gpu.count" 
                    :min="1" 
                    :max="8" 
                    style="width: 80px" 
                  />
                </a-form-item>
              </a-col>
              <a-col :span="5">
                <a-form-item label="存储大小(GB)" :name="['resources', 'storage']"
                  :rules="[{ required: true, message: '请输入存储大小' }, { type: 'number', min: 10, message: '存储大小必须大于10GB' }]"
                  class="no-margin"
                >
                  <a-input-number 
                    v-model:value="formState.resources.storage" 
                    :min="10" 
                    style="width: 100px" 
                  />
                </a-form-item>
              </a-col>
            </a-row>
          </a-card>

          <!-- 扩缩容配置 -->
          <a-card title="扩缩容配置" size="small" class="mb-4">
            <a-row :gutter="16">
              <a-col :span="4">
                <a-form-item 
                  label="最小副本数" 
                  :name="['scaling_policy', 'min_replicas']"
                  :rules="[{ required: true, message: '请输入最小副本数' }]"
                  class="scale-form-item"
                >
                  <a-input-number 
                    v-model:value="formState.scaling_policy.min_replicas" 
                    :min="1" 
                    :max="16" 
                    style="width: 100%" 
                  />
                </a-form-item>
              </a-col>

              <a-col :span="4">
                <a-form-item 
                  label="最大副本数" 
                  :name="['scaling_policy', 'max_replicas']"
                  :rules="[{ required: true, message: '请输入最大副本数' }]"
                  class="scale-form-item"
                >
                  <a-input-number 
                    v-model:value="formState.scaling_policy.max_replicas" 
                    :min="1" 
                    :max="16" 
                    style="width: 100%" 
                  />
                </a-form-item>
              </a-col>

              <a-col :span="4">
                <a-form-item 
                  label="目标并发数" 
                  :name="['scaling_policy', 'concurrency', 'target']"
                  :rules="[{ required: true, message: '请输入目标并发数' }]"
                  class="scale-form-item"
                >
                  <a-input-number 
                    v-model:value="formState.scaling_policy.concurrency.target" 
                    :min="1" 
                    style="width: 100%" 
                  />
                </a-form-item>
              </a-col>

              <a-col :span="6">
                <a-form-item 
                  label="扩容窗口(秒)" 
                  :name="['scaling_policy', 'scale_up_window_seconds']"
                  :rules="[{ required: true, message: '请输入扩容窗口' }]"
                  class="scale-form-item"
                >
                  <a-input-number 
                    v-model:value="formState.scaling_policy.scale_up_window_seconds" 
                    :min="1" 
                    style="width: 100%" 
                  />
                </a-form-item>
              </a-col>

              <a-col :span="6">
                <a-form-item 
                  label="缩容窗口(秒)" 
                  :name="['scaling_policy', 'scale_down_window_seconds']"
                  :rules="[{ required: true, message: '请输入缩容窗口' }]"
                  class="scale-form-item"
                >
                  <a-input-number 
                    v-model:value="formState.scaling_policy.scale_down_window_seconds" 
                    :min="1" 
                    style="width: 100%" 
                  />
                </a-form-item>
              </a-col>
            </a-row>
          </a-card>

          <!-- Environment Variables -->
          <a-card title="环境变量配置" size="small" class="mb-4">
            <div v-for="key in Object.keys(formState.env)" :key="key" style="margin-bottom: 8px">
              <a-space style="width: 100%" align="start">
                <a-input :value="key" placeholder="Key" style="width: 200px" disabled />
                <a-textarea
                  v-model:value="formState.env[key]"
                  placeholder="Value"
                  :auto-size="{ minRows: 1, maxRows: 10 }"
                  style="width: 400px"
                />
                <a-button type="link" @click="() => handleDeleteEnv(key)"
                  :disabled="Object.keys(defaultEnv).includes(key)">删除</a-button>
              </a-space>
            </div>
            <div v-if="showNewEnvInput" style="margin-bottom: 8px">
              <a-space style="width: 100%" align="start">
                <a-input v-model:value="newEnvKey" placeholder="Key" style="width: 200px" 
                  @pressEnter="handleAddEnv" />
                <a-textarea
                  v-model:value="newEnvValue"
                  placeholder="Value"
                  :auto-size="{ minRows: 1, maxRows: 10 }"
                  style="width: 400px"
                  @pressEnter="handleAddEnv"
                />
                <a-button type="link" @click="handleAddEnv">确定</a-button>
                <a-button type="link" @click="cancelAddEnv">取消</a-button>
              </a-space>
            </div>
            <a-button type="dashed" @click="showAddEnv" block v-if="!showNewEnvInput">
              添加环境变量
            </a-button>
          </a-card>

          <!-- Actions -->
          <a-form-item>
            <a-space>
              <a-button type="primary" @click="handleSubmit">生成YAML</a-button>
              <a-button @click="resetForm">重置</a-button>
            </a-space>
          </a-form-item>
        </a-card>
      </div>

      <!-- 右侧YAML预览 -->
      <div class="preview-section">
        <a-card title="YAML预览" class="preview-card">
          <template #extra>
            <a-space>
              <a-button type="primary" @click="applyYaml" :disabled="!yamlContent">应用YAML</a-button>
              <a-button type="primary" @click="downloadYaml">下载YAML</a-button>
              <a-button @click="copyYaml">复制内容</a-button>
            </a-space>
          </template>
          <pre v-html="highlightedYaml || '暂无内容'"></pre>
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import yaml from 'js-yaml'
import { message } from 'ant-design-vue'
import hljs from 'highlight.js/lib/core'
import yamlLang from 'highlight.js/lib/languages/yaml'
import 'highlight.js/styles/github.css'
import { regions, gpuModels, templates, imageSecrets, defaultEnv, getGpuProvider, getHfHome } from '@/config/regions'

hljs.registerLanguage('yaml', yamlLang)

const formRef = ref()
const yamlContent = ref('')
const selectedRegion = ref('')

const formState = reactive({
  name: '',
  image: '',
  args: [],
  resources: {
    gpu: {
      model: '',
      count: 1,
    },
    storage: 70,
  },
  env: { ...defaultEnv },
  max_concurrent_per_worker: 16,
  scaling_policy: {
    min_replicas: 1,
    max_replicas: 16,
    concurrency: {
      target: 8
    },
    scale_up_window_seconds: 10,
    scale_down_window_seconds: 120
  }
})

const rules = {
  name: [
    { required: true, message: '请输入模型名称' },
    { pattern: /^[a-z0-9-]+$/, message: '模型名称只能包含小写字母、数字和连字符' }
  ],
  image: [{ required: true, message: '请输入镜像地址' }],
}

// 高亮显示的YAML
const highlightedYaml = computed(() => {
  if (!yamlContent.value) return ''
  return hljs.highlight(yamlContent.value, { language: 'yaml' }).value
})

const argsText = ref('')

// 将数组转换为文本
const arrayToText = (arr) => arr.join('\n')

// 将文本转换为数组
const textToArray = (text) => text.split('\n').filter(line => line.trim() !== '')

// 处理参数变化
const handleArgsChange = () => {
  formState.args = textToArray(argsText.value)
}

// 监听GPU型号和区域变化，更新GPU_PROVIDER
watch(
  [() => formState.resources.gpu.model, selectedRegion],
  ([gpuModel, region]) => {
    if (gpuModel && region) {
      formState.env.GPU_PROVIDER = getGpuProvider(region, gpuModel)
    }
  }
)

const handleRegionChange = (value) => {
  const region = regions[value]
  const template = templates[`us-${region.name.toLowerCase()}`]
  
  if (template) {
    formState.name = template.name
    formState.args = [...template.args]
    argsText.value = arrayToText(template.args)
    formState.env = { ...defaultEnv, ...template.env }
  } else {
    // 如果没有找到模板，设置基本的环境变量和根据区域设置HF_HOME
    formState.env = { 
      ...defaultEnv,
      HF_HOME: getHfHome(value)
    }
    if (formState.resources.gpu.model) {
      formState.env.GPU_PROVIDER = getGpuProvider(value, formState.resources.gpu.model)
    }
  }
}

// 修改 handleSubmit，移除 try-catch，实时生成YAML
const generateYaml = () => {
  const region = regions[selectedRegion.value]
  const yamlObj = {
    name: formState.name,
    image: formState.image,
    args: formState.args,
    image_secrets: [imageSecrets],
    resources: formState.resources,
    cluster: selectedRegion.value,
    volumes: region ? [region.volumes] : undefined,
    env: formState.env,
    port: 8080,
    health_check: {
      path: '/health'
    },
    max_concurrent_per_worker: 16,
    scaling_policy: {
      min_replicas: 1,
      max_replicas: 16,
      concurrency: {
        target: 8
      },
      scale_up_window_seconds: 10,
      scale_down_window_seconds: 120
    }
  }
  yamlContent.value = yaml.dump(yamlObj, { indent: 2, lineWidth: -1 })
}

// 监听表单变化
watch(
  [() => formState, selectedRegion],
  () => {
    generateYaml()
  },
  { deep: true }
)

// 修改 handleSubmit
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    generateYaml()
  } catch (error) {
    console.error('Validation failed:', error)
  }
}

const resetForm = () => {
  formRef.value.resetFields()
  yamlContent.value = ''
  selectedRegion.value = ''
  argsText.value = ''
}

const downloadYaml = () => {
  const blob = new Blob([yamlContent.value], { type: 'text/yaml' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${formState.name}.yaml`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}

const copyYaml = async () => {
  try {
    await navigator.clipboard.writeText(yamlContent.value)
    message.success('已复制到剪贴板')
  } catch (err) {
    message.error('复制失败')
  }
}

const showNewEnvInput = ref(false)
const newEnvKey = ref('')
const newEnvValue = ref('')

const showAddEnv = () => {
  showNewEnvInput.value = true
  newEnvKey.value = ''
  newEnvValue.value = ''
}

const cancelAddEnv = () => {
  showNewEnvInput.value = false
  newEnvKey.value = ''
  newEnvValue.value = ''
}

const handleAddEnv = () => {
  if (!newEnvKey.value) {
    message.error('请输入环境变量名称')
    return
  }
  if (formState.env[newEnvKey.value] !== undefined) {
    message.error('环境变量名称已存在')
    return
  }
  formState.env[newEnvKey.value] = newEnvValue.value
  showNewEnvInput.value = false
  newEnvKey.value = ''
  newEnvValue.value = ''
}

const handleDeleteEnv = (key) => {
  if (Object.keys(defaultEnv).includes(key)) {
    message.warning('不能删除预设的环境变量')
    return
  }
  delete formState.env[key]
}

// Add watch for image changes
watch(
  () => formState.image,
  (newImage) => {
    if (newImage) {
      const tagMatch = newImage.match(/:([^:\/]+)$/)
      if (tagMatch) {
        formState.env.VLLM_VERSION = tagMatch[1]
      } else {
        formState.env.VLLM_VERSION = 'latest'
      }
    }
  }
)

const applyYaml = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/modelapi`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/yaml',
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
      body: yamlContent.value
    })
    
    if (!response.ok) {
      const errorText = await response.text()
      let errorMessage = `应用YAML失败 (${response.status})`
      try {
        // 尝试解析错误响应为 JSON
        const errorJson = JSON.parse(errorText)
        errorMessage += `: ${errorJson.message || errorJson.error || errorText}`
      } catch {
        // 如果不是 JSON 格式，直接使用错误文本
        errorMessage += `: ${errorText}`
      }
      throw new Error(errorMessage)
    }
    
    message.success('YAML应用成功')
  } catch (error) {
    message.error(error.message)
  }
}
</script>

<style scoped>
.model-api-container {
  padding: 24px;
}
.layout-container {
  display: flex;
  gap: 24px;
  min-height: calc(100vh - 140px);
}
.edit-section {
  flex: 1;
  min-width: 0;
}
.preview-section {
  width: 45%;
  min-width: 500px;
}
.preview-card {
  position: sticky;
  top: 24px;
  height: calc(100vh - 140px);
  overflow-y: auto;
}
.mb-4 {
  margin-bottom: 16px;
}
pre {
  background-color: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  min-height: 200px;
}
:deep(.hljs) {
  background-color: #f5f5f5;
  padding: 0;
}
:deep(.ant-card-body) {
  padding: 16px;
}
:deep(.ant-form-item) {
  margin-bottom: 16px;
}
:deep(.ant-input-number) {
  width: 100%;
}
.scale-config-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0 8px;
}

.scale-form-item {
  margin-bottom: 0 !important;
}

.scale-form-item :deep(.ant-form-item-label) {
  padding: 0 0 8px;
}

.scale-form-item :deep(.ant-form-item-label > label) {
  height: unset;
}

.no-margin {
  margin-bottom: 0 !important;
}

:deep(.ant-form-item-label) {
  padding-bottom: 4px !important;
}

:deep(.ant-row) {
  flex-wrap: nowrap !important;
}
</style> 