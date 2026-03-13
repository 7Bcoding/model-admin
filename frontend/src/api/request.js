import axios from 'axios'
import router from '@/router'
import { ElMessage } from 'element-plus'

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 5000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    error => {
        console.error('Request error:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => response,
    error => {
        console.error('Response error:', error)
        
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

export default request 