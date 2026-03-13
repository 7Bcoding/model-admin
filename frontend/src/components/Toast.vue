<template>
  <transition name="toast">
    <div v-if="visible" class="toast" :class="type">
      {{ message }}
    </div>
  </transition>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'Toast',
  setup() {
    const visible = ref(false)
    const message = ref('')
    const type = ref('info')

    const show = (msg, t = 'info', duration = 3000) => {
      message.value = msg
      type.value = t
      visible.value = true
      setTimeout(() => {
        visible.value = false
      }, duration)
    }

    return {
      visible,
      message,
      type,
      show
    }
  }
}
</script>

<style scoped>
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 10px 20px;
  border-radius: 4px;
  color: white;
  z-index: 9999;
}

.error {
  background-color: #dc3545;
}

.info {
  background-color: #17a2b8;
}

.success {
  background-color: #28a745;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style> 