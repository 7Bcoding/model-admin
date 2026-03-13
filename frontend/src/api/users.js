import axios from 'axios'

// 获取用户列表
export const getUsers = async () => {
  try {
    console.log('Calling getUsers API...')
    console.log('Current token:', localStorage.getItem('token'))
    console.log('Current user:', localStorage.getItem('user'))
    
    const response = await axios.get('/users')
    console.log('API Response:', response.data)
    return response.data
  } catch (error) {
    console.error('Error details:', {
      response: error.response,
      request: error.request,
      message: error.message
    })
    throw error
  }
}

// 添加新用户
export const addUser = async (userData) => {
  try {
    const response = await axios.post('/users', userData)
    return response.data
  } catch (error) {
    if (error.response?.status === 403) {
      throw new Error('没有权限执行此操作')
    }
    throw error
  }
}

// 更新用户角色
export const updateUserRole = async (username, role) => {
  try {
    const response = await axios.put(`/users/${username}/role`, { role })
    return response.data
  } catch (error) {
    if (error.response?.status === 403) {
      throw new Error('没有权限执行此操作')
    }
    throw error
  }
}

// 删除用户
export const deleteUser = async (username) => {
  try {
    const response = await axios.delete(`/users/${username}`)
    return response.data
  } catch (error) {
    if (error.response?.status === 403) {
      throw new Error('没有权限执行此操作')
    }
    throw error
  }
} 