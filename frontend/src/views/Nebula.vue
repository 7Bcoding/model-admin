<template>
  <div class="page-container">
    <div class="header">
      <div class="title-section">
        <h2 class="page-title">Nebula 管理</h2>
      </div>
    </div>
    <div class="content-card">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="CNode" name="cnode"></el-tab-pane>
        <el-tab-pane label="Deployment" name="deployment"></el-tab-pane>
        <el-tab-pane label="Worker" name="worker"></el-tab-pane>
      </el-tabs>
      <div class="tab-content">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const activeTab = ref('cnode')

// 处理标签页变化
const handleTabChange = (tabName) => {
  console.log('Tab changed to:', tabName)
  activeTab.value = tabName
  // 更新 URL，但不重新加载页面
  router.push(`/nebula/${tabName}`, { replace: true })
}

// 根据路由设置活动标签页
const setActiveTabFromRoute = () => {
  const path = route.path
  if (path.includes('/cnode')) {
    activeTab.value = 'cnode'
  } else if (path.includes('/deployment')) {
    activeTab.value = 'deployment'
  } else if (path.includes('/worker')) {
    activeTab.value = 'worker'
  } else {
    activeTab.value = 'cnode'
  }
}

// 监听路由变化
watch(() => route.path, () => {
  setActiveTabFromRoute()
})

onMounted(() => {
  setActiveTabFromRoute()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.title-section {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  margin-left: 10px;
}

.content-card {
  background-color: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.tab-content {
  margin-top: 20px;
}
</style> 