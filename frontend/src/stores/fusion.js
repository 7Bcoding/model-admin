import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useFusionStore = defineStore('fusion', () => {
  // Backend selection state
  const selectedBackend = ref('beta') // 'beta' or 'alpha'
  
  // Backend configurations
  const backends = {
    beta: {
      name: 'beta',
      baseUrl: '/fusion-beta',
      description: 'beta Fusion Backend'
    },
    alpha: {
      name: 'alpha',
      baseUrl: '/fusion-alpha',
      description: 'alpha Fusion Backend'
    }
  }
  
  // Computed properties
  const currentBackend = computed(() => backends[selectedBackend.value])
  const currentBaseUrl = computed(() => currentBackend.value.baseUrl)
  
  // Actions
  const setBackend = (backend) => {
    console.log('Setting backend to:', backend)
    if (backends[backend]) {
      selectedBackend.value = backend
      console.log('Backend set successfully. Current base URL:', currentBaseUrl.value)
    } else {
      console.error('Invalid backend:', backend)
    }
  }
  
  const getBackendOptions = () => {
    return Object.entries(backends).map(([key, config]) => ({
      value: key,
      label: config.name,
      description: config.description
    }))
  }
  
  return {
    selectedBackend,
    backends,
    currentBackend,
    currentBaseUrl,
    setBackend,
    getBackendOptions
  }
}) 