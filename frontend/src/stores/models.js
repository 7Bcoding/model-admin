import { defineStore } from 'pinia'
import { getModels } from '../api/models'

export const useModelsStore = defineStore('models', {
  state: () => ({
    models: [],
    lastFetchTime: null,
    cacheExpirationTime: 5 * 60 * 1000 // 5分钟缓存过期时间
  }),

  getters: {
    // 获取缓存的模型列表
    getModels: (state) => state.models,
    
    // 检查缓存是否有效
    isCacheValid: (state) => {
      if (!state.lastFetchTime) return false
      const now = Date.now()
      return (now - state.lastFetchTime) < state.cacheExpirationTime
    }
  },

  actions: {
    // 获取模型列表，forceRefresh 参数用于强制刷新
    async fetchModels(forceRefresh = false) {
      // 如果缓存有效且不强制刷新，直接返回缓存数据
      if (this.isCacheValid && !forceRefresh) {
        return this.models
      }

      try {
        const response = await getModels()
        if (response?.data?.data?.data) {
          const modelsList = response.data.data.data
          
          if (Array.isArray(modelsList)) {
            // 处理模型数据
            this.models = modelsList.map(item => ({
              model_name: item.model_name || '',
              description: item.description || '',
              status: item.status || 'unknown',
              private: item.private || false,
              max_tokens: item.max_tokens || 0,
              tags: item.tags || [],
              starred: item.starred || false,
              endpoints: item.endpoints || [],
              input_token_price: item.input_token_price || 0,
              output_token_price: item.output_token_price || 0,
              model: item.model || {},
              note: item.note || '',
              open_chat_id: item.open_chat_id || '',
              inference_engine: item.inference_engine || ''
            }))
            
            // 更新最后获取时间
            this.lastFetchTime = Date.now()
            return this.models
          }
        }
        throw new Error('Invalid response format')
      } catch (error) {
        console.error('Error fetching models:', error)
        throw error
      }
    },

    // 清除缓存
    clearCache() {
      this.models = []
      this.lastFetchTime = null
    }
  }
}) 