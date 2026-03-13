<template>
  <div id="app" :class="{ 'login-page': !isLoggedIn }">
    <nav class="navbar navbar-expand-lg custom-navbar" v-if="isLoggedIn">
      <div class="container-fluid">
        <router-link class="navbar-brand" to="/">LLM OPS</router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav me-auto">
            <li class="nav-item">
              <router-link 
                class="nav-link" 
                to="/models"
                :class="{ active: isModelsActive }"
              >模型管理(NOVITA)</router-link>
            </li>
            <li class="nav-item">
              <router-link 
                class="nav-link" 
                to="/ppio-models"
                :class="{ active: isPPIOModelsActive }"
              >模型管理(PPIO)</router-link>
            </li>
            <li class="nav-item" v-if="isAdmin || isOps">
              <router-link 
                class="nav-link" 
                to="/resources"
                :class="{ active: isResourcesActive }"
              >资源管理</router-link>
            </li>
            <li class="nav-item" v-if="isAdmin || isOps">
              <router-link 
                class="nav-link" 
                to="/serverless"
                :class="{ active: isServerlessActive }"
              >Serverless管理</router-link>
            </li>
            <li class="nav-item" v-if="isAdmin || isOps">
              <router-link 
                class="nav-link" 
                to="/nebula"
                :class="{ active: isNebulaActive }"
              >Nebula管理</router-link>
            </li>
            <li class="nav-item" v-if="(isAdmin || isOps) && showModelApiTab">
              <router-link 
                class="nav-link" 
                to="/model-api"
                :class="{ active: isModelApiActive }"
              >模型API管理</router-link>
            </li>
            <li class="nav-item" v-if="isAdmin || isOps">
              <router-link 
                class="nav-link" 
                to="/fusion"
                :class="{ active: isFusionActive }"
              >Fusion</router-link>
            </li>
            <li class="nav-item" v-if="isAdmin">
              <router-link 
                class="nav-link" 
                to="/users"
                :class="{ active: isUsersActive }"
              >用户管理</router-link>
            </li>
          </ul>
          <div class="d-flex align-items-center">
            <el-button
              type="primary"
              plain
              @click="goToAuditLog"
              class="audit-log-btn me-3"
            >
              <i class="el-icon-document"></i>
              操作日志
            </el-button>
            <el-dropdown trigger="click" @command="handleCommand">
              <span class="user-dropdown me-3">
                {{ username }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <button class="btn-logout" @click="logout">退出</button>
          </div>
        </div>
      </div>
    </nav>
    
    <!-- 添加全局提醒 -->
    <!-- <div v-if="isLoggedIn && showAlert" class="global-alert">
      <div class="alert-content">
        <span class="alert-text">
          客户放量SOP：确认客户最大 RPM -> 调整并确认 RPM 频控配置 (@viktor / @李洛克 / @悟饭) -> 资源池扩充以满足最大 RPM -> 同步客户开始放量
        </span>
        <div class="alert-actions">
          <el-button type="warning" size="small" @click="hideAlert">我知道了</el-button>
          <button class="close-button" @click="hideAlert">
            <i class="el-icon-close"></i>
          </button>
        </div>
      </div>
    </div> -->

    <!-- keep-alive 只缓存模型管理页 -->
    <keep-alive>
      <router-view v-slot="{ Component }" v-if="$route.name === 'Models'">
        <component :is="Component" />
      </router-view>
    </keep-alive>
    <router-view v-if="$route.name !== 'Models'" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const route = useRoute()

const isLoggedIn = computed(() => {
  return !!localStorage.getItem('token')
})

const username = computed(() => {
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  return user.username || ''
})

const isAdmin = computed(() => {
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  return user.role === 'admin'
})

const isOps = computed(() => {
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  return user.role === 'operator'
})

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

// 判断是否在资源管理页面
const isResourcesActive = computed(() => {
  return route.path === '/resources'
})

// 判断是否在 Serverless 管理页面
const isServerlessActive = computed(() => {
  return route.path.startsWith('/serverless')
})

const isNebulaActive = computed(() => {
  return route.path.startsWith('/nebula')
})

// 控制模型API管理TAB的可见性（用于测试）
const showModelApiTab = ref(false)

// 判断是否在模型 API 管理页面
const isModelApiActive = computed(() => {
  return route.path === '/model-api'
})

// 判断是否在用户管理页面
const isUsersActive = computed(() => {
  return route.path === '/users'
})

// 判断是否在 Fusion 管理页面
const isFusionActive = computed(() => {
  return route.path === '/fusion'
})

const handleCommand = async (command) => {
  if (command === 'changePassword') {
    try {
      const { value: formData } = await ElMessageBox.prompt(
        '请输入新密码',
        '修改密码',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          inputType: 'password'
        }
      )
      
      if (formData) {
        await axios.post('/users/change-password', {
          new_password: formData
        })
        ElMessage.success('密码修改成功')
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error.response?.data?.message || '密码修改失败')
      }
    }
  } else if (command === 'logout') {
    logout()
  }
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}

const goToAuditLog = () => {
  router.push('/audit-log')
}

// 添加 alert 相关的状态
const showAlert = ref(true)

// 隐藏提醒
const hideAlert = () => {
  showAlert.value = false
}

// 每次组件挂载时重置显示状态
onMounted(() => {
  showAlert.value = true
})
</script>

<style>
:root {
  --primary-color: #1a73e8;
  --secondary-color: #5f6368;
  --background-color: #f8f9fa;
  --navbar-height: 64px;
  --sidebar-width: 250px;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  background-color: var(--background-color);
}

.custom-navbar {
  background: white;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  height: var(--navbar-height);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.custom-navbar .navbar-brand {
  color: var(--primary-color);
  font-weight: 600;
  font-size: 1.25rem;
}

.custom-navbar .nav-link {
  color: var(--secondary-color);
  font-weight: 500;
  padding: 0.5rem 1rem;
  transition: all 0.2s ease;
  position: relative;
  margin: 0 0.25rem;
  border-radius: 6px;
}

.custom-navbar .nav-link:hover {
  color: var(--primary-color);
  background-color: rgba(26, 115, 232, 0.04);
}

.custom-navbar .nav-link.active {
  color: var(--primary-color);
  font-weight: 600;
  background-color: rgba(26, 115, 232, 0.08);
  box-shadow: inset 0 0 0 1px rgba(26, 115, 232, 0.16);
}

.custom-navbar .nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 3px;
  background: var(--primary-color);
  border-radius: 3px;
}

.custom-navbar .nav-link.router-link-active:not(.active) {
  color: var(--secondary-color);
  background: none;
  box-shadow: none;
}

.custom-navbar .nav-link.router-link-active:not(.active)::after {
  display: none;
}

.user-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--secondary-color);
  font-weight: 500;
}

.user-dropdown:hover {
  color: var(--primary-color);
}

.btn-logout {
  padding: 0.5rem 1rem;
  border: 1px solid var(--primary-color);
  background: transparent;
  color: var(--primary-color);
  border-radius: 4px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-logout:hover {
  background: var(--primary-color);
  color: white;
}

.router-view-container {
  padding-top: calc(var(--navbar-height) + 52px) !important;
}

:deep(.el-dropdown-menu) {
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

:deep(.el-dropdown-menu__item) {
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
}

:deep(.el-dropdown-menu__item:hover) {
  background-color: #f3f4f6;
  color: var(--primary-color);
}

.audit-log-btn {
  margin-right: 12px;
}

/* 全局提醒样式 */
.global-alert {
  position: fixed;
  top: var(--navbar-height);
  left: 0;
  right: 0;
  z-index: 99;
  background-color: #fff3cd;
  border: 1px solid #ffeeba;
  padding: 12px 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.alert-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.alert-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.alert-text {
  color: #856404;
  font-size: 14px;
  line-height: 1.5;
  flex: 1;
}

.close-button {
  background: none;
  border: none;
  color: #856404;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s;
  margin-left: 8px;
}

.close-button:hover {
  opacity: 0.7;
}

/* 调整警告按钮样式 */
:deep(.el-button--warning) {
  background-color: #e6a23c;
  border-color: #e6a23c;
  color: white;
}

:deep(.el-button--warning:hover) {
  background-color: #ebb563;
  border-color: #ebb563;
  color: white;
}
</style>