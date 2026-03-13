import axios from 'axios'

// 创建 axios 实例
const instance = axios.create({
  baseURL: '/api/v1'  // 设置基础 URL
})

export default instance 