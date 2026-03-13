<template>
  <div class="login-container">
    <div class="login-box">
      <h1>LLM OPS</h1>
      <div class="login-form">
        <button type="button" @click="handleFeishuLogin" class="login-button">
          <img src="@/assets/feishu-logo.png" alt="飞书" class="feishu-logo">
          <span>飞书登录</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/composables/toast'
import axios from 'axios'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const toast = useToast()

    const handleFeishuLogin = async () => {
      const clientId = 'f3c9bbc3ec62ed982b4a'; // 替换为你的飞书App ID  
      const redirectUri = encodeURIComponent(import.meta.env.VITE_FEISHU_CALLBACK_URL);
      const state = Math.random().toString(36).substring(2); // 防CSRF
  
      window.location.href = `https://iam.beta.cloud/login/oauth/authorize?` +
      `client_id=${clientId}&response_type=code&redirect_uri=${redirectUri}&scope=profile&state=xaspedfi1x`;
    };

    return {
      handleFeishuLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1001;
}

.login-box {
  background: white;
  padding: 3rem 2rem;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 420px;
  text-align: center;
}

.login-box h1 {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 2.5rem;
  font-size: 2.5rem;
  font-weight: 700;
}

.login-form {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-button {
  background: #00d6b9;
  color: white;
  padding: 18px 32px;
  border: none;
  border-radius: 10px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  box-shadow: 0 4px 16px rgba(0, 214, 185, 0.3);
  min-height: 56px;
  width: 100%;
  max-width: 280px;
  line-height: 1.2;
}

.login-button:hover {
  background: #00b8a0;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 214, 185, 0.4);
}

.login-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(0, 214, 185, 0.3);
}

.feishu-logo {
  width: 28px;
  height: 28px;
  object-fit: contain;
  flex-shrink: 0;
}

.login-button span {
  white-space: nowrap;
  overflow: visible;
}

@media (max-width: 480px) {
  .login-box {
    margin: 1rem;
    padding: 2rem 1.5rem;
  }
  
  .login-box h1 {
    font-size: 2rem;
    margin-bottom: 2rem;
  }
  
  .login-button {
    padding: 16px 24px;
    font-size: 16px;
    min-height: 52px;
  }
  
  .feishu-logo {
    width: 24px;
    height: 24px;
  }
}
</style> 