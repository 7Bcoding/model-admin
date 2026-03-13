import axios from 'axios'
import { useToast } from '@/composables/toast'
import { useFusionStore } from '@/stores/fusion'
import router from '@/router'
import { ElMessage } from 'element-plus'

// Create a separate axios instance for Fusion backend
const createFusionInstance = () => {
  const fusionStore = useFusionStore()
  
  console.log('Creating Fusion API instance with base URL:', fusionStore.currentBaseUrl)
  
  const fusionInstance = axios.create({
    baseURL: fusionStore.currentBaseUrl,
    // Add timeout for development
    timeout: 10000
  })
  
  // Request interceptor, add token
  fusionInstance.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (token && token !== 'undefined') {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  }, error => {
    return Promise.reject(error)
  })
  
  // Response interceptor
  fusionInstance.interceptors.response.use(
    response => {
      return response
    },
    error => {
      const toast = useToast()
      
      // 处理401未授权错误
      if (error.response && error.response.status === 401) {
        // 清除本地存储的token和用户信息
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        
        // 显示提示消息
        ElMessage.error('登录已过期，请重新登录')
        
        // 跳转到登录页面
        router.push('/login')
        
        return Promise.reject(error)
      }
      
      if (error.response) {
        // Show specific error message
        const message = error.response.data?.message || error.message
        toast.show(message, 'error')
        console.error('Fusion API Error:', {
          status: error.response.status,
          message: message,
          url: error.config.url,
          method: error.config.method
        })
      } else {
        toast.show('Network error', 'error')
        console.error('Fusion Network Error:', error)
      }
      return Promise.reject(error)
    }
  )
  
  return fusionInstance
}

export default createFusionInstance