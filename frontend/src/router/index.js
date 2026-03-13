import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Models from '../views/Models.vue'
import Users from '../views/Users.vue'
import EndpointManager from '@/views/EndpointManager.vue'
import Resources from '../views/Resources.vue'
import Serverless from '../views/Serverless.vue'
import WorkerList from '../views/WorkerList.vue'
import ModelApi from '../views/ModelApi.vue'
import ExcludeWorker from '../views/ExcludeWorker.vue'
import AuditLog from '../views/AuditLog.vue'
import FeishuCallback from '../views/FeishuCallback.vue'
import Nebula from '../views/Nebula.vue'
import NebulaCnode from '../views/NebulaCnde.vue'
import NebulaDeployment from '../views/NebulaDeployment.vue'
import NebulaWorker from '../views/NebulaWorker.vue'
import Fusion from '../views/Fusion.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { public: true }
  },
  {
    path: '/',
    redirect: '/models'
  },
  {
    path: '/models',
    name: 'Models',
    component: () => import('@/views/Models.vue')
  },
  {
    path: '/ppio-models',
    name: 'PPIOModels',
    component: () => import('@/views/PPIOModels.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/models/:modelName/edit-note',
    name: 'EditModelNote',
    component: () => import('@/views/EditModelNote.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { requiresAdmin: true }
  },
  {
    path: '/models/:modelName/endpoints',
    name: 'EndpointManager',
    component: EndpointManager,
    meta: { requiresAuth: true }
  },
  {
    path: '/resources',
    name: 'Resources',
    component: Resources,
    beforeEnter: (to, from, next) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (user.role === 'admin' || user.role === 'operator') {
        next()
      } else {
        next('/models')
      }
    },
    meta: { requiresAuth: true }
  },
  {
    path: '/serverless',
    name: 'Serverless',
    component: Serverless,
    meta: { requiresAuth: true }
  },
  {
    path: '/nebula',
    name: 'Nebula',
    component: Nebula,
    redirect: '/nebula/cnode',
    beforeEnter: (to, from, next) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (user.role === 'admin' || user.role === 'operator') {
        next()
      } else {
        next('/models')
      }
    },
    meta: { requiresAuth: true },
    children: [
      {
        path: 'cnode',
        name: 'NebulaCnode',
        component: NebulaCnode,
        beforeEnter: (to, from, next) => {
          const user = JSON.parse(localStorage.getItem('user') || '{}')
          if (user.role === 'admin' || user.role === 'operator') {
            next()
          } else {
            next('/models')
          }
        },
        meta: { requiresAuth: true }
      },
      {
        path: 'deployment',
        name: 'NebulaDeployment',
        component: NebulaDeployment,
        beforeEnter: (to, from, next) => {
          const user = JSON.parse(localStorage.getItem('user') || '{}')
          if (user.role === 'admin' || user.role === 'operator') {
            next()
          } else {
            next('/models')
          }
        },
        meta: { requiresAuth: true }
      },
      {
        path: 'worker',
        name: 'NebulaWorker',
        component: NebulaWorker,
        beforeEnter: (to, from, next) => {
          const user = JSON.parse(localStorage.getItem('user') || '{}')
          if (user.role === 'admin' || user.role === 'operator') {
            next()
          } else {
            next('/models')
          }
        },
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: '/serverless/:endpointId',
    name: 'WorkerList',
    component: WorkerList,
    meta: { requiresAuth: true }
  },
  {
    path: '/model-api',
    name: 'ModelApi',
    component: ModelApi,
    meta: { requiresAuth: true }
  },
  {
    path: '/fusion',
    name: 'Fusion',
    component: Fusion,
    beforeEnter: (to, from, next) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (user.role === 'admin' || user.role === 'operator') {
        next()
      } else {
        next('/models')
      }
    },
    meta: { requiresAuth: true }
  },
  {
    path: '/exclude-worker',
    name: 'ExcludeWorker',
    component: ExcludeWorker,
    meta: { requiresAuth: true }
  },
  {
    path: '/audit-log',
    name: 'AuditLog',
    component: AuditLog,
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/feishu/callback',
    name: 'FeishuCallback',
    component: FeishuCallback,
    meta: { 
      public: true,
      requiresAuth: false 
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  // 如果路由不需要认证，直接放行
  if (to.meta.requiresAuth === false) {
    next()
    return
  }
  
  // 如果是公开页面（登录页），检查是否需要重定向
  if (to.meta.public) {
    if (token && to.path === '/login') {
      // 如果已经登录且访问登录页，重定向到首页
      next('/')
    } else {
      next()
    }
    return
  }
  
  // 其他页面检查是否登录
  if (!token || token === 'undefined') {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    next('/login')
    return
  }
  
  // 检查管理员权限
  if (to.meta.requiresAdmin) {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    if (!user || user.role !== 'admin') {
      next('/models')
      return
    }
  }
  
  next()
})

export default router