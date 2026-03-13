import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import axios from 'axios'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import { ElMessage } from 'element-plus'

// Bootstrap
import 'bootstrap/dist/css/bootstrap.min.css'
import * as bootstrap from 'bootstrap'

// 将 bootstrap 添加到全局 window 对象
window.bootstrap = bootstrap

// 设置 axios 默认值
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL
axios.defaults.headers.common['Content-Type'] = 'application/json'

// 添加请求拦截器
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  // 添加日志以便调试
  console.log('Request headers:', config.headers)
  return config
})

// 添加响应拦截器
axios.interceptors.response.use(
  response => {
    return response
  },
  error => {
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

// 创建 Vue 应用
const app = createApp(App)

// 创建 Pinia 实例
const pinia = createPinia()

app.use(pinia)
app.use(ElementPlus)
app.use(Antd)
app.use(router)

// 添加全局错误处理
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue Error:', err)
  console.error('Error Info:', info)
}

app.mount('#app') 