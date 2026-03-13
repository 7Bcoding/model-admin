import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useFusionStore = defineStore('fusion', () => {
  // Backend selection state
  const selectedBackend = ref('ppio') // 'ppio' or 'novita'
  
  // Backend configurations
  const backends = {
    ppio: {
      name: 'PPIO',
      baseUrl: '/fusion-ppio',
      description: 'PPIO Fusion Backend'
    },
    novita: {
      name: 'Novita',
      baseUrl: '/fusion-novita',
      description: 'Novita Fusion Backend'
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