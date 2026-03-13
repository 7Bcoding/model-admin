require('dotenv').config()
const express = require('express')
const path = require('path')
const cors = require('cors')
const { createProxyMiddleware } = require('http-proxy-middleware')

const app = express()
const port = process.env.PORT || 3000
const backendUrl = process.env.BACKEND_URL || 'http://127.0.0.1:8080'

// 启用 CORS
app.use(cors())

// API 代理
app.use('/api', createProxyMiddleware({
  target: backendUrl,
  changeOrigin: true,
  pathRewrite: {
    '^/api': '/api/v1' // 重写路径
  }
}))

// 静态文件服务
app.use(express.static(path.join(__dirname, '../dist')))

// 所有其他请求返回 index.html
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, '../dist/index.html'))
})

app.listen(port, () => {
  console.log(`Frontend server running at http://localhost:${port}`)
  console.log(`Proxying API requests to ${backendUrl}`)
}) 