import { ref, markRaw } from 'vue'
import Toast from '@/components/Toast.vue'

const toast = ref(null)

export function useToast() {
  const show = (message, type = 'info', duration = 3000) => {
    if (toast.value) {
      toast.value.show(message, type, duration)
    }
  }

  const setToastRef = (el) => {
    toast.value = el
  }

  return {
    show,
    setToastRef
  }
} 