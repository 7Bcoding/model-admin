import createFusionInstance from './fusion'

// Model prices management APIs
export const getModelPrices = (params = {}) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get('/model-prices', { params })
}

export const getModelPrice = (id) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get(`/model-prices/${id}`)
}

export const createModelPrice = (priceData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.post('/model-prices', priceData)
}

export const updateModelPrice = (id, priceData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.put(`/model-prices/${id}`, priceData)
}

export const deleteModelPrice = (id) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.delete(`/model-prices/${id}`)
}