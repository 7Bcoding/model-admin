import axios from 'axios'
import router from '@/router'
import { ElMessage } from 'element-plus'

const API_URL = import.meta.env.VITE_API_BASE_URL || ''

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
  return config
})

api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    console.error('Nebula API Error:', error)
    
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

// CNode APIs
export const getCNodes = (params = {}) => {
  return api.get('/nebula/cnodes', { params })
}

export const getCNode = (name) => {
  return api.get(`/nebula/cnodes/${name}`)
}

// NDeployment APIs
export const getNDeployments = (params = {}) => {
  return api.get('/nebula/ndeployments', { params })
}

// Nebula Worker APIs
export const getNebulaWorkers = (params = {}) => {
  return api.get('/nebula/workers', { params })
} 