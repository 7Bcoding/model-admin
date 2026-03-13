// 区域定义
export const regions = {
  '1': {
    name: 'US-01',
    volumes: {
      mount: '/data',
      size: 300,
      network_id: '4c546248-5531-4146-86be-674a8ff971ab'
    }
  },
  '2': {
    name: 'US-CA-01',
    volumes: {
      mount: '/data',
      size: 300,
      network_id: '425f26e8-3753-4349-bf93-c93e2692a338'
    }
  },
  '3': {
    name: 'US-CA-02',
  },
  '4': {
    name: 'US-NYC-01',
  }
}

// GPU型号定义
export const gpuModels = [
  { value: 'nvidia/3090', label: 'NVIDIA 3090', memory: '24GB', provider: '3090' },
  { value: 'nvidia/4090', label: 'NVIDIA 4090', memory: '24GB', provider: '4090' },
  { value: 'nvidia/4090D', label: 'NVIDIA 4090D', memory: '24GB', provider: '4090d' },
  { value: 'nvidia/L20', label: 'NVIDIA L20', memory: '24GB', provider: 'l20' },
  { value: 'nvidia/L40', label: 'NVIDIA L40', memory: '48GB', provider: 'l40' },
  { value: 'nvidia/A6000', label: 'NVIDIA A6000', memory: '48GB', provider: 'a6000' },
  { value: 'nvidia/RTX6000Ada', label: 'NVIDIA RTX 6000 Ada', memory: '48GB', provider: 'rtx6000ada' },
  { value: 'nvidia/A100_SXM4_80GB', label: 'NVIDIA A100 SXM4', memory: '80GB', provider: 'a100' },
  { value: 'nvidia/A100_PCIE_80GB', label: 'NVIDIA A100 PCIe', memory: '80GB', provider: 'a100' },
  { value: 'nvidia/A100_SXM4_40GB', label: 'NVIDIA A100 SXM4', memory: '40GB', provider: 'a100' },
  { value: 'nvidia/H800', label: 'NVIDIA H800', memory: '80GB', provider: 'h800' },
  { value: 'nvidia/H20', label: 'NVIDIA H20', memory: '80GB', provider: 'h20' },
  { value: 'nvidia/H100', label: 'NVIDIA H100', memory: '80GB', provider: 'h100' }
]

// 预设的环境变量
export const defaultEnv = {
  VLLM_ENGINE_ITERATION_TIMEOUT_S: '240',
  LOG_UPLOAD_ENDPOINT: 'aws',
  LOG_UPLOAD_LOCATION: '/var/log/supervisor/',
  HF_TOKEN: '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'
}

// 获取GPU Provider名称
export const getGpuProvider = (region, gpuType) => {
  const regionMap = {
    '1': 'us-01',
    '2': 'ca-01',
    '3': 'ca-02',
    '4': 'nyc-01'
  }
  const gpu = gpuModels.find(g => g.value === gpuType)
  return gpu ? `novita-serverless-${regionMap[region]}-${gpu.provider}` : undefined
}

// 获取HF_HOME值
export const getHfHome = (region) => {
  return regions[region]?.volumes ? '/data' : '/hf_data'
}

// 预设模板定义
export const templates = {
  'us-01': {
    name: 'us-01-template',
    args: [
      '--disable-log-requests',
      '--port',
      '8080',
      '--model',
      'meta-llama/Llama-3.1-70B-Instruct',
      '--served-model-name',
      'meta-llama/llama-3.1-70b-instruct',
      '--swap-space',
      '16',
      '--max-num-seqs',
      '32',
      '--gpu-memory-utilization',
      '0.9',
      '--max-model-len',
      '32768'
    ],
    env: {
      ...defaultEnv,
      HF_HOME: '/data',
      MODEL_ID: 'meta-llama/llama-3.1-70b-instruct',
      GPU_PROVIDER: 'novita-serverless-us-01-h20'
    },
    cluster: '1'
  },
  'us-ca-01': {
    name: 'ca-01-template',
    args: [
      '--disable-log-requests',
      '--port',
      '8080',
      '--model',
      'meta-llama/Llama-3.1-70B-Instruct',
      '--served-model-name',
      'meta-llama/llama-3.1-70b-instruct',
      '--swap-space',
      '16',
      '--max-num-seqs',
      '32',
      '--gpu-memory-utilization',
      '0.9',
      '--max-model-len',
      '32768'
    ],
    env: {
      ...defaultEnv,
      HF_HOME: '/data',
      MODEL_ID: 'meta-llama/llama-3.1-70b-instruct',
      GPU_PROVIDER: 'novita-serverless-ca-01-h20'
    },
    cluster: '2'
  },
  'us-ca-02': {
    name: 'ca-02-template',
    args: [
      '--disable-log-requests',
      '--port',
      '8080',
      '--model',
      'meta-llama/Llama-3.1-70B-Instruct',
      '--served-model-name',
      'meta-llama/llama-3.1-70b-instruct',
      '--swap-space',
      '16',
      '--max-num-seqs',
      '32',
      '--gpu-memory-utilization',
      '0.9',
      '--max-model-len',
      '32768'
    ],
    env: {
      ...defaultEnv,
      HF_HOME: '/hf_data',
      MODEL_ID: 'meta-llama/llama-3.1-70b-instruct',
      GPU_PROVIDER: 'novita-serverless-ca-02-h20'
    },
    cluster: '3'
  }
}

// 固定的镜像密钥配置
export const imageSecrets = {
  registry: 'ghcr.io',
  username: 'AllenShen',
  password: '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'
} 