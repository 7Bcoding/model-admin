import createFusionInstance from './fusion'

// Model management APIs
export const getModels = (params = {}) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get('/models', { params })
}

export const getModel = (name) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get(`/models/${name}`)
}

export const createModel = (modelData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.post('/models', modelData)
}

export const updateModel = (name, modelData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.put(`/models/${name}`, modelData)
}

export const deleteModel = (name) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.delete(`/models/${name}`)
}

// Provider management APIs
export const addModelProvider = (modelName, providerData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.post(`/models/${modelName}/providers`, providerData)
}

export const updateModelProvider = (modelName, providerName, providerData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.put(`/models/${modelName}/providers/${providerName}`, providerData)
}

export const removeModelProvider = (modelName, providerName, isFallback = false) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.delete(`/models/${modelName}/providers/${providerName}`, {
    params: { is_fallback: isFallback }
  })
}

export const updateModelProviderWeight = (modelName, providerName, weight, isFallback = false) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.put(`/models/${modelName}/providers/${providerName}/weight`, {}, {
    params: { weight, is_fallback: isFallback }
  })
} 