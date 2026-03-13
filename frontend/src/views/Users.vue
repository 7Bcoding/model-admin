<template>
  <div class="users-page">
    <div class="header d-flex justify-content-between align-items-center">
      <h2 class="page-title">用户管理</h2>
      <button class="btn btn-success" @click="showAddUserModal">
        <i class="bi bi-person-plus"></i> 添加用户
      </button>
    </div>

    <div class="card">
      <div class="card-body">
        <div class="table-responsive">
          <table class="table table-hover align-middle">
            <thead>
              <tr>
                <th>用户名</th>
                <th>账号</th>
                <th>角色</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>{{ user.username }}</td>
                <td>{{ user.account_name }}</td>
                <td>
                  <span class="badge" :class="getRoleBadgeClass(user.role)">
                    {{ getRoleText(user.role) }}
                  </span>
                </td>
                <td>
                  <div class="btn-group">
                    <button 
                      class="btn btn-sm btn-outline-primary"
                      @click="showEditRole(user)"
                      :disabled="user.username === currentUser?.username"
                    >
                      修改角色
                    </button>
                    <button 
                      class="btn btn-sm btn-outline-danger"
                      @click="confirmDelete(user)"
                      :disabled="user.username === currentUser?.username"
                    >
                      删除
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- 添加用户模态框 -->
    <div class="modal fade" id="addUserModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">添加用户</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <form @submit.prevent="handleAddUser" autocomplete="off">
            <div class="modal-body">
              <div class="mb-3">
                <label class="form-label">用户名</label>
                <input 
                  type="text" 
                  class="form-control" 
                  v-model="newUser.username" 
                  required
                  autocomplete="off"
                >
              </div>
              <div class="mb-3">
                <label class="form-label">账号</label>
                <input 
                  type="text" 
                  class="form-control" 
                  v-model="newUser.account_name" 
                  required
                  autocomplete="off"
                >
              </div>
              <div class="mb-3">
                <label class="form-label">密码</label>
                <input 
                  type="password" 
                  class="form-control" 
                  v-model="newUser.password" 
                  required
                  autocomplete="new-password"
                >
              </div>
              <div class="mb-3">
                <label class="form-label">角色</label>
                <select class="form-select" v-model="newUser.role" required>
                  <option value="admin">管理员</option>
                  <option value="operator">运维</option>
                  <option value="user">普通用户</option>
                </select>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
              <button type="submit" class="btn btn-primary">添加</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 修改角色模态框 -->
    <div class="modal fade" id="editRoleModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">修改角色</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <form @submit.prevent="handleUpdateRole">
            <div class="modal-body">
              <div class="mb-3">
                <label class="form-label">用户名</label>
                <input type="text" class="form-control" :value="selectedUser?.username" disabled>
              </div>
              <div class="mb-3">
                <label class="form-label">角色</label>
                <select class="form-select" v-model="selectedRole" required>
                  <option value="admin">管理员</option>
                  <option value="operator">运维</option>
                  <option value="user">普通用户</option>
                </select>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
              <button type="submit" class="btn btn-primary">保存</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useToast } from '@/composables/toast'
import { getUsers, addUser, updateUserRole, deleteUser } from '@/api/users'

export default {
  name: 'Users',
  setup() {
    const toast = useToast()
    const users = ref([])
    const newUser = ref({
      username: '',
      account_name: '',
      password: '',
      role: 'user'
    })
    const selectedUser = ref(null)
    const selectedRole = ref('')
    const addUserModal = ref(null)
    const editRoleModal = ref(null)

    const currentUser = JSON.parse(localStorage.getItem('user') || '{}')

    const fetchUsers = async () => {
      try {
        const response = await getUsers()
        console.log('API Response:', response)
        if (response.code === 200) {
          users.value = response.data.data || []
          console.log('Processed users:', users.value)
        } else {
          throw new Error(response.message || '获取用户列表失败')
        }
      } catch (error) {
        console.error('Error fetching users:', error)
        toast.show(error.message || '获取用户列表失败', 'error')
      }
    }

    const showAddUserModal = () => {
      if (!addUserModal.value) {
        addUserModal.value = new window.bootstrap.Modal(document.getElementById('addUserModal'))
      }
      newUser.value = { username: '', account_name: '', password: '', role: 'user' }
      addUserModal.value.show()
    }

    const handleAddUser = async () => {
      try {
        const response = await addUser(newUser.value)
        if (response.code === 200) {
          newUser.value = {
            username: '',
            account_name: '',
            password: '',
            role: 'user'
          }
          if (addUserModal.value) {
            addUserModal.value.hide()
          }
          await fetchUsers()
          toast.show('添加用户成功', 'success')
        } else {
          throw new Error(response.message || '添加用户失败')
        }
      } catch (error) {
        console.error('Error adding user:', error)
        if (error.response?.status === 403) {
          toast.show('没有权限执行此操作', 'error')
        } else {
          toast.show(error.message || '添加用户失败', 'error')
        }
      }
    }

    const showEditRole = (user) => {
      if (!editRoleModal.value) {
        editRoleModal.value = new window.bootstrap.Modal(document.getElementById('editRoleModal'))
      }
      selectedUser.value = user
      selectedRole.value = user.role
      editRoleModal.value.show()
    }

    const handleUpdateRole = async () => {
      try {
        const response = await updateUserRole(selectedUser.value.username, selectedRole.value)
        if (response.code === 200) {
          const modal = window.bootstrap.Modal.getInstance(document.getElementById('editRoleModal'))
          if (modal) {
            modal.hide()
            Promise.resolve().then(async () => {
              await fetchUsers()
              toast.show('修改角色成功', 'success')
            })
          }
        } else {
          throw new Error(response.message || '修改角色失败')
        }
      } catch (error) {
        console.error('Error updating role:', error)
        toast.show(error.message || '修改角色失败', 'error')
      }
    }

    const confirmDelete = async (user) => {
      if (confirm(`确定要删除用户 ${user.username} 吗？`)) {
        try {
          const response = await deleteUser(user.username)
          if (response.code === 200) {
            toast.show('删除用户成功', 'success')
            await fetchUsers()
          } else {
            throw new Error(response.message || '删除用户失败')
          }
        } catch (error) {
          console.error('Error deleting user:', error)
          if (error.response?.status === 403) {
            toast.show('没有权限执行此操作', 'error')
          } else {
            toast.show(error.message || '删除用户失败', 'error')
          }
        }
      }
    }

    const getRoleBadgeClass = (role) => {
      const classes = {
        admin: 'bg-danger',
        operator: 'bg-primary',
        user: 'bg-info'
      }
      return classes[role] || 'bg-secondary'
    }

    const getRoleText = (role) => {
      const texts = {
        admin: '管理员',
        operator: '运维',
        user: '普通用户'
      }
      return texts[role] || role
    }

    onMounted(() => {
      fetchUsers()
    })

    return {
      users,
      newUser,
      selectedUser,
      selectedRole,
      currentUser,
      showAddUserModal,
      handleAddUser,
      showEditRole,
      handleUpdateRole,
      confirmDelete,
      getRoleBadgeClass,
      getRoleText
    }
  }
}
</script>

<style scoped>
.users-page {
  padding: 24px;
  background: #f5f5f5;
  min-height: calc(100vh - var(--navbar-height));
}

.header {
  margin-bottom: 32px;
  padding: 0 12px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0;
  padding-bottom: 8px;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.table th {
  background: #f8f9fa;
  font-weight: 600;
}

.btn-group {
  gap: 4px;
}

.badge {
  padding: 6px 10px;
  font-weight: normal;
}
</style>