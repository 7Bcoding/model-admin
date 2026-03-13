import createFusionInstance from './fusion'

// Account management APIs
export const getAccounts = (params = {}) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get('/accounts', { params })
}

export const getAccount = (name) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.get(`/accounts/${name}`)
}

export const createAccount = (accountData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.post('/accounts', accountData)
}

export const updateAccount = (name, accountData) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.put(`/accounts/${name}`, accountData)
}

export const deleteAccount = (name) => {
  const fusionInstance = createFusionInstance()
  return fusionInstance.delete(`/accounts/${name}`)
} 