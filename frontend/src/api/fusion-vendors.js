import createFusionInstance from './fusion'

// Vendor management APIs
export const getVendors = (params = {}) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get('/vendors', { params })
}

export const getVendor = (name) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get(`/vendors/${name}`)
}

export const createVendor = (vendorData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.post('/vendors', vendorData)
}

export const updateVendor = (name, vendorData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.put(`/vendors/${name}`, vendorData)
}

export const deleteVendor = (name) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.delete(`/vendors/${name}`)
}
