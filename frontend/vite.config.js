import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver, AntDesignVueResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver(), AntDesignVueResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver(), AntDesignVueResolver({
        importStyle: false,
      })],
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      'bootstrap': 'bootstrap'
    }
  },
  server: {
    port: process.env.VITE_PORT || 3000,
    cors: true,
    proxy: {
      '/fusion-ppio': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/fusion-ppio/, '/api/v1/fusion'),
        configure: (proxy, options) => {
          proxy.on('error', (err, req, res) => {
            console.log('fusion-ppio proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // 添加 X-Fusion-Provider header 来标识 PPIO 服务
            proxyReq.setHeader('X-Fusion-Provider', 'ppio');
            console.log('Sending fusion-ppio Request to Backend:', req.method, req.url, '-> /api/v1/fusion');
          });
          proxy.on('proxyRes', (proxyRes, req, res) => {
            console.log('Received fusion-ppio Response from Backend:', proxyRes.statusCode, req.url);
          });
        },
      },
      '/fusion-novita': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/fusion-novita/, '/api/v1/fusion'),
        configure: (proxy, options) => {
          proxy.on('error', (err, req, res) => {
            console.log('fusion-novita proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // 添加 X-Fusion-Provider header 来标识 Novita 服务
            proxyReq.setHeader('X-Fusion-Provider', 'novita');
            console.log('Sending fusion-novita Request to Backend:', req.method, req.url, '-> /api/v1/fusion');
          });
          proxy.on('proxyRes', (proxyRes, req, res) => {
            console.log('Received fusion-novita Response from Backend:', proxyRes.statusCode, req.url);
          });
        },
      }
    }
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
        modifyVars: {
          'primary-color': '#1a73e8',
        },
      },
    },
  },
})