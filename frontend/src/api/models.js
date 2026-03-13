import axios from 'axios'
import router from '@/router'
import { ElMessage } from 'element-plus'

const API_URL = import.meta.env.VITE_API_BASE_URL || ''
console.log('API_URL:', API_URL)

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  console.log('Request config:', config)
  return config
})

api.interceptors.response.use(
  response => {
    console.log('Response:', response)
    return response
  },
  error => {
    console.error('API Error:', error)
    
    // 处理401未授权错误
    if (error.response && error.response.status === 401) {
      // 清除本地存储的token和用户信息
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      
      // 显示提示消息
      ElMessage.error('登录已过期，请重新登录')
      
      // 跳转到登录页面
      router.push('/login')
    }
    
    return Promise.reject(error)
  }
)

export const getModels = (params = {}) => {
  return api.get('/models', { params })
}

export const starModel = async (modelName, platform) => {
  return api.post('/models/star', { 
    model_name: modelName,
    platform
  })
}

export const unstarModel = async (modelName, platform) => {
  return api.post('/models/unstar', { 
    model_name: modelName,
    platform
  })
}

export const updateModelNote = (modelName, note, openChatId, inferenceEngine, platform) => {
  return api.post('/models/note', { 
    model_name: modelName,
    note: note,
    open_chat_id: openChatId,
    inference_engine: inferenceEngine,
    platform
  })
} 