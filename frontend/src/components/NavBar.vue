<template>
  <nav class="navbar">
    <div class="logo">
      LLM OPS
    </div>
    <div class="nav-links">
      <router-link 
        to="/models" 
        :class="{ active: isModelsActive }"
      >模型管理(NOVITA)</router-link>
      <router-link 
        to="/ppio-models" 
        :class="{ active: isPPIOModelsActive }"
      >模型管理(PPIO)</router-link>
      <router-link to="/resources">资源管理</router-link>
      <router-link to="/serverless">Serverless管理</router-link>
      <router-link to="/model-api">模型API管理</router-link>
      <router-link v-if="isAdmin" to="/users">用户管理</router-link>
    </div>
    <div class="user-actions">
      <button class="logout-btn" @click="handleLogout">退出</button>
    </div>
  </nav>
</template>

<script>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { logout } from '@/api/auth'

export default {
  name: 'NavBar',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    
    const isAdmin = computed(() => user.role === 'admin')

    // 判断是否在模型管理页面
    const isModelsActive = computed(() => {
      // 如果是直接访问模型管理页面
      if (route.path === '/models') return true
      // 如果是在 endpoint 详情页面，且没有 platform 参数
      if (route.name === 'EndpointManager' && !route.query.platform) return true
      return false
    })

    // 判断是否在 PPIO 模型管理页面
    const isPPIOModelsActive = computed(() => {
      // 如果是直接访问 PPIO 模型管理页面
      if (route.path === '/ppio-models') return true
      // 如果是在 endpoint 详情页面，且 platform 参数为 ppio
      if (route.name === 'EndpointManager' && route.query.platform === 'ppio') return true
      return false
    })

    const handleLogout = async () => {
      console.log('Logout button clicked')
      try {
        const response = await logout()
        console.log('Logout API response:', response)
        
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        console.log('Local storage cleared')
        
        window.location.href = '/login'
      } catch (error) {
        console.error('Logout error:', error)
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }
    }

    return {
      isAdmin,
      handleLogout,
      isModelsActive,
      isPPIOModelsActive
    }
  }
}
</script>

<style scoped>
.navbar {
  background: #fff;
  padding: 0;
  display: flex;
  align-items: center;
  height: 60px;
  border-bottom: 1px solid #e4e7ed;
}

.logo {
  font-size: 1.2rem;
  font-weight: bold;
  padding: 0 20px;
  color: #333;
}

.nav-links {
  display: flex;
  height: 100%;
}

.nav-links a {
  color: #333;
  text-decoration: none;
  padding: 0 20px;
  height: 100%;
  display: flex;
  align-items: center;
  font-size: 14px;
  position: relative;
  transition: color 0.3s;
}

.nav-links a:hover {
  color: #409EFF;
}

.nav-links a.active {
  color: #409EFF !important;
}

.nav-links a.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: #409EFF;
}

.user-actions {
  margin-left: auto;
  padding-right: 20px;
}

.logout-btn {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 14px;
  padding: 8px 15px;
}

.logout-btn:hover {
  color: #409EFF;
}
</style> 