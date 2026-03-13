import { ElMessage } from 'element-plus'

export function useToast() {
  const show = (message, type = 'info') => {
    ElMessage({
      message,
      type,
      duration: 3000
    })
  }

  return {
    show
  }
} 