export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token'),
    user: JSON.parse(localStorage.getItem('user') || '{}'),
  }),

  actions: {
    hasRole(role) {
      if (Array.isArray(role)) {
        return role.includes(this.user.role)
      }
      if (role === 'operator') {
        return this.user.role === 'admin' || this.user.role === 'operator'
      }
      return this.user.role === role
    },
    
    setUserInfo(userInfo) {
      this.user = userInfo
      this.token = userInfo.token
      localStorage.setItem('user', JSON.stringify(userInfo))
      localStorage.setItem('token', userInfo.token)
    },

    logout() {
      this.user = {}
      this.token = null
      localStorage.removeItem('user')
      localStorage.removeItem('token')
    }
  }
}) 