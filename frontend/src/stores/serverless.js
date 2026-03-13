import { defineStore } from 'pinia'
import axios from 'axios'

export const useServerlessStore = defineStore('serverless', {
  state: () => ({
    endpoints: [],
    lastFetchTime: null,
    cacheExpirationTime: 20 * 60 * 1000 // 20分钟缓存过期时间
  }),

  getters: {
    // 获取缓存的 endpoints 列表
    getEndpoints: (state) => state.endpoints,
    
    // 检查缓存是否有效
    isCacheValid: (state) => {
      if (!state.lastFetchTime) return false
      const now = Date.now()
      return (now - state.lastFetchTime) < state.cacheExpirationTime
    }
  },

  actions: {
    // 获取 endpoints 列表，forceRefresh 参数用于强制刷新
    async fetchEndpoints(forceRefresh = false) {
      // 如果缓存有效且不强制刷新，直接返回缓存数据
      if (this.isCacheValid && !forceRefresh) {
        return this.endpoints
      }

      try {
        // 获取 endpoints 列表
        const endpointsResponse = await axios.get('/se')
        const endpoints = endpointsResponse.data.data

        // 获取普通模型列表
        const modelsResponse = await axios.get('/models')
        const models = modelsResponse.data?.data?.data || []

        // 获取 beta 模型列表
        const betaModelsResponse = await axios.get('/models', { params: { platform: 'beta' } })
        const betaModels = betaModelsResponse.data?.data?.data || []

        // 处理模型关联
        const extractEndpointFromURL = (url) => {
          if (!url) return ''
          url = url.replace('https://', '').replace('http://', '')
          const matches = url.split('.runsync.alpha.dev')
          return matches[0] || ''
        }

        // 创建 endpoint 到模型的映射
        const endpointToModels = new Map()

        // 处理普通模型
        models.forEach(model => {
          if (model.endpoints) {
            model.endpoints.forEach(endpoint => {
              const endpointName = extractEndpointFromURL(endpoint.url)
              if (endpointName) {
                if (!endpointToModels.has(endpointName)) {
                  endpointToModels.set(endpointName, new Map())
                }
                endpointToModels.get(endpointName).set(model.model_name, { 
                  name: model.model_name,
                  source: 'normal'
                })
              }
            })
          }
        })

        // 处理 beta 模型
        betaModels.forEach(model => {
          if (model.endpoints) {
            model.endpoints.forEach(endpoint => {
              const endpointName = extractEndpointFromURL(endpoint.url)
              if (endpointName) {
                if (!endpointToModels.has(endpointName)) {
                  endpointToModels.set(endpointName, new Map())
                }
                endpointToModels.get(endpointName).set(model.model_name + '_beta', {
                  name: model.model_name,
                  source: 'beta'
                })
              }
            })
          }
        })

        // 更新 endpoints 的模型列表
        this.endpoints = endpoints.map(endpoint => ({
          ...endpoint,
          models: Array.from(endpointToModels.get(endpoint.metadata.name)?.values() || [])
        }))
            
        // 更新最后获取时间
        this.lastFetchTime = Date.now()
        return this.endpoints
      } catch (error) {
        console.error('Error fetching endpoints:', error)
        throw error
      }
    },

    // 清除缓存
    clearCache() {
      this.endpoints = []
      this.lastFetchTime = null
    }
  }
}) 