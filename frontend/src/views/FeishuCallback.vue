<template>
  <div class="callback-container">
    <div class="callback-box">
      <h2>飞书登录处理中...</h2>
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

export default {
  name: 'FeishuCallback',
  setup() {
    const router = useRouter()

    onMounted(async () => {
      try {
        // 获取 URL 中的 code 参数
        const urlParams = new URLSearchParams(window.location.search)
        const code = urlParams.get('code')
        const state = urlParams.get('state')

        if (!code) {
          throw new Error('未获取到授权码')
        }

        // 发送 code 到后端处理
        const response = await axios.post('/feishu/callback', {
          code,
          state
        })

        if (response.data && response.data.data) {
          const { token, user } = response.data.data
          
          // 保存 token 和用户信息
          localStorage.setItem('token', token)
          localStorage.setItem('user', JSON.stringify(user))
          
          // 设置 axios 默认 header
          axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
          
          // 跳转到首页
          router.push('/')
        }
      } catch (error) {
        console.error('飞书登录失败:', error)
        router.push('/login')
      }
    })

    return {}
  }
}
</script>

<style scoped>
.callback-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
}

.callback-box {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.callback-box h2 {
  color: #2c3e50;
  margin: 0;
}
</style>