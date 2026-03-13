import { createStore } from 'vuex'

export default createStore({
  state: {
    user: {
      isAuthenticated: false,
      username: '',
      role: '',
      token: ''
    },
    activeMenu: ''
  },
  mutations: {
    setUser(state, user) {
      state.user.isAuthenticated = true
      state.user.username = user.username
      state.user.role = user.role
      state.user.token = user.token
    },
    clearUser(state) {
      state.user.isAuthenticated = false
      state.user.username = ''
      state.user.role = ''
      state.user.token = ''
    },
    setActiveMenu(state, path) {
      state.activeMenu = path
    }
  },
  actions: {
    login({ commit }, user) {
      commit('setUser', user)
    },
    logout({ commit }) {
      commit('clearUser')
    }
  }
}) 