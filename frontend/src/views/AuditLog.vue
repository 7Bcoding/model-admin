<template>
  <div class="page-container">
    <div class="header">
      <div class="title-section">
        <i class="el-icon-document title-icon"></i>
        <h2 class="page-title">操作日志</h2>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="refreshData">
          <i class="el-icon-refresh"></i> 刷新
        </el-button>
      </div>
    </div>

    <div class="content-card">
      <div class="filter-section">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="操作人">
            <el-input
              v-model="filterForm.operator"
              placeholder="请输入操作人"
              clearable
              @clear="handleFilter"
            />
          </el-form-item>
          <el-form-item label="操作类型">
            <el-select
              v-model="filterForm.operationType"
              placeholder="请选择操作类型"
              clearable
              @clear="handleFilter"
              style="width: 220px"
            >
              <el-option
                v-for="type in operationTypes"
                :key="type"
                :label="type"
                :value="type"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="资源名称">
            <el-input
              v-model="filterForm.resourceName"
              placeholder="请输入资源名称"
              clearable
              @clear="handleFilter"
              style="width: 220px"
            />
          </el-form-item>
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="filterForm.timeRange"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              value-format="YYYY-MM-DD HH:mm:ss"
              :default-time="['00:00:00', '23:59:59']"
              @clear="handleFilter"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleFilter">查询</el-button>
            <el-button @click="resetFilter">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        v-loading="loading"
        :data="auditLogs"
        style="width: 100%"
        border
      >
        <el-table-column
          prop="operator"
          label="操作人"
          min-width="120"
        />
        <el-table-column
          prop="action"
          label="操作类型"
          min-width="150"
        />
        <el-table-column
          prop="target"
          label="资源名称"
          min-width="150"
        />
        <el-table-column
          prop="request_url"
          label="请求URL"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column
          prop="detail"
          label="操作详情"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column
          prop="request_body"
          label="请求内容"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column
          prop="created_at"
          label="操作时间"
          min-width="180"
          sortable
        >
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column
          label="状态"
          width="100"
          align="center"
        >
          <template #default="scope">
            <el-tag :type="scope.row.status === 200 ? 'success' : 'danger'">
              {{ scope.row.result }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const loading = ref(false)
const auditLogs = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const operationTypes = ref([])

const filterForm = ref({
  operator: '',
  operationType: '',
  resourceName: '',
  timeRange: []
})

// 格式化日期时间
const formatDateTime = (datetime) => {
  if (!datetime) return ''
  const date = new Date(datetime)
  return date.toLocaleString()
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      operator: filterForm.value.operator,
      action: filterForm.value.operationType,
      target: filterForm.value.resourceName,
      startTime: filterForm.value.timeRange?.[0],
      endTime: filterForm.value.timeRange?.[1]
    }

    const response = await axios.get('/audit-logs', { params })
    auditLogs.value = response.data.data.items
    total.value = response.data.data.total
    
    // 获取所有唯一的操作类型
    const types = new Set(response.data.data.items.map(log => log.action))
    operationTypes.value = Array.from(types)
  } catch (error) {
    console.error('获取审计日志失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

// 处理筛选
const handleFilter = () => {
  currentPage.value = 1
  fetchData()
}

// 重置筛选
const resetFilter = () => {
  filterForm.value = {
    operator: '',
    operationType: '',
    resourceName: '',
    timeRange: []
  }
  handleFilter()
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchData()
}

// 处理每页条数变化
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
}

.header {
  margin-bottom: 32px;
  padding: 0 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 28px;
  color: #3b82f6;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

.content-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.filter-section {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 