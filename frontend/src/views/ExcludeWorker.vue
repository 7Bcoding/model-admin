<template>
  <div class="exclude-worker-container">
    <div class="page-header">
      <el-button type="text" size="small" @click="$router.push('/serverless')" class="back-button">
        <i class="el-icon-arrow-left"></i> 返回
      </el-button>
      <h2>断流Worker管理</h2>
      <el-button type="primary" @click="showAddDialog" class="add-button">
        添加
      </el-button>
    </div>
    <el-card class="box-card">
      <!-- Worker列表 -->
      <el-table :data="workerList" style="width: 100%">
        <el-table-column prop="worker" label="Worker ID" width="500" />
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(scope.row.worker)"
            >
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加Worker对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="添加断流Worker"
      width="500px"
    >
      <el-form :model="form" label-width="120px">
        <el-form-item label="Worker ID">
          <el-input v-model="form.worker" placeholder="请输入Worker ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAdd">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const workerList = ref([])
const dialogVisible = ref(false)
const form = ref({
  worker: ''
})

// 获取断流worker列表
const fetchWorkerList = async () => {
  try {
    const response = await axios.get('/exclude-workers')
    if (response.status !== 200) {
      throw new Error('获取数据失败')
    }
    workerList.value = response.data.map(worker => ({ worker }))
  } catch (error) {
    ElMessage.error('获取断流Worker列表失败')
    console.error('Error:', error)
  }
}

// 显示添加对话框
const showAddDialog = () => {
  form.value.worker = ''
  dialogVisible.value = true
}

// 添加断流worker
const handleAdd = async () => {
  if (!form.value.worker) {
    ElMessage.warning('请输入Worker ID')
    return
  }

  try {
    const response = await axios.post('/exclude-workers', null, {
      params: {
        worker: form.value.worker
      }
    })
    
    if (response.status !== 200) {
      throw new Error('添加失败')
    }

    ElMessage.success('添加成功')
    dialogVisible.value = false
    fetchWorkerList()
  } catch (error) {
    ElMessage.error('添加断流Worker失败')
    console.error('Error:', error)
  }
}

// 删除断流worker
const handleDelete = async (worker) => {
  try {
    const response = await axios.delete('/exclude-workers', {
      params: {
        worker: worker
      }
    })
    
    if (response.status !== 200) {
      throw new Error('删除失败')
    }

    ElMessage.success('删除成功')
    fetchWorkerList()
  } catch (error) {
    const errorMessage = error.response?.data?.message || error.response?.data?.reason || error.message || '删除断流Worker失败'
    ElMessage.error(errorMessage)
    console.error('Error:', error)
  }
}

onMounted(() => {
  fetchWorkerList()
})
</script>

<style scoped>
.exclude-worker-container {
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
}

.back-button {
  margin-right: 10px;
  font-size: 14px;
}

.add-button {
  margin-left: auto;
}

h2 {
  margin: 0;
  font-size: 20px;
}

.action-buttons {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-bottom: 20px;
}
</style>