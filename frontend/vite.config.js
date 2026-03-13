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
      '/fusion-beta': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/fusion-beta/, '/api/v1/fusion'),
        configure: (proxy, options) => {
          proxy.on('error', (err, req, res) => {
            console.log('fusion-beta proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // 添加 X-Fusion-Provider header 来标识 beta 服务
            proxyReq.setHeader('X-Fusion-Provider', 'beta');
            console.log('Sending fusion-beta Request to Backend:', req.method, req.url, '-> /api/v1/fusion');
          });
          proxy.on('proxyRes', (proxyRes, req, res) => {
            console.log('Received fusion-beta Response from Backend:', proxyRes.statusCode, req.url);
          });
        },
      },
      '/fusion-alpha': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/fusion-alpha/, '/api/v1/fusion'),
        configure: (proxy, options) => {
          proxy.on('error', (err, req, res) => {
            console.log('fusion-alpha proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // 添加 X-Fusion-Provider header 来标识 alpha 服务
            proxyReq.setHeader('X-Fusion-Provider', 'alpha');
            console.log('Sending fusion-alpha Request to Backend:', req.method, req.url, '-> /api/v1/fusion');
          });
          proxy.on('proxyRes', (proxyRes, req, res) => {
            console.log('Received fusion-alpha Response from Backend:', proxyRes.statusCode, req.url);
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