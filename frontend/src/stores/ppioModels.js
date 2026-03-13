import { defineStore } from 'pinia'
import { getModels } from '../api/models'

export const usebetaModelsStore = defineStore('betaModels', {
  state: () => ({
    models: [],
    lastFetchTime: null,
    cacheExpirationTime: 5 * 60 * 1000 // 5 minutes cache expiration time
  }),

  getters: {
    getModels: (state) => state.models,
    
    isCacheValid: (state) => {
      if (!state.lastFetchTime) return false
      const now = Date.now()
      return (now - state.lastFetchTime) < state.cacheExpirationTime
    }
  },

  actions: {
    async fetchModels(forceRefresh = false) {
      if (this.isCacheValid && !forceRefresh) {
        return this.models
      }

      try {
        // Add platform=beta parameter to the request
        const response = await getModels({ platform: 'beta' })
        if (response?.data?.data?.data) {
          const modelsList = response.data.data.data
          
          if (Array.isArray(modelsList)) {
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
            
            this.lastFetchTime = Date.now()
            return this.models
          }
        }
        throw new Error('Invalid response format')
      } catch (error) {
        console.error('Error fetching beta models:', error)
        throw error
      }
    },

    clearCache() {
      this.models = []
      this.lastFetchTime = null
    }
  }
}) 