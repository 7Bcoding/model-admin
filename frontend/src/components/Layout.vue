<template>
  <div class="app-layout">
    <NavBar />
    <main class="main-content">
      <slot></slot>
    </main>
  </div>
</template>

<script>
import NavBar from './NavBar.vue'

export default {
  name: 'Layout',
  components: {
    NavBar
  }
}
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 20px;
  background: #f5f5f5;
}
</style> 