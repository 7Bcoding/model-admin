import axios from 'axios'
import { useToast } from '@/composables/toast'
import router from '@/router'
import { ElMessage } from 'element-plus'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL
})

// 请求拦截器，添加 token
instance.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token && token !== 'undefined') {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}, error => {
  return Promise.reject(error)
})

// 响应拦截器
instance.interceptors.response.use(
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
      // 显示具体的错误信息
      const message = error.response.data?.message || error.message
      toast.show(message, 'error')
      console.error('API Error:', {
        status: error.response.status,
        message: message,
        url: error.config.url,
        method: error.config.method
      })
    } else {
      toast.show('Network error', 'error')
      console.error('Network Error:', error)
    }
    return Promise.reject(error)
  }
)

export default instance 