# Fusion 代理功能

## 概述

为了提高安全性和简化架构，我们在后端服务器上实现了统一的 Fusion 代理功能。通过 header 驱动的路由选择机制，使用单一接口来处理前端发来的 `fusion-beta/*` 和 `fusion-alpha/*` 路径的请求，并转发到相应的后端服务。

## 配置

### 1. 配置文件设置

在配置文件中添加以下配置项：

```yaml
# Fusion 代理配置
beta_fusion:
  url: https://api.ppinfra.com
  token: # 加密后的 beta Fusion Token

alpha_fusion:
  url: https://api.alpha.ai
  token: # 加密后的 alpha Fusion Token
```

### 2. Token 加密

Token 需要使用 AES 加密后存储在配置文件中。使用与其他敏感信息相同的 `aes_key` 进行加密。

## 统一代理架构

### Header 驱动的路由选择

统一代理服务根据 `X-Fusion-Provider` header 来决定路由目标：
- `X-Fusion-Provider: beta` → 路由到 beta Fusion 服务

### 路由映射

#### 统一代理接口
- 前端请求路径：`/api/v1/fusion/*`
- 根据 `X-Fusion-Provider` header 转发到：
  - beta: `{beta_fusion.url}/admin/v1/*`
  - alpha: `{alpha_fusion.url}/admin/v1/*`

## 示例

### 请求示例
```
前端请求：GET /fusion-beta/models
Nginx 转发：GET /api/v1/fusion/models (设置 X-Fusion-Provider: beta)
后端代理到：GET https://api.ppinfra.com/admin/v1/models
```

```
前端请求：POST /fusion-alpha/deployments
Nginx 转发：POST /api/v1/fusion/deployments (设置 X-Fusion-Provider: alpha)
后端代理到：POST https://api.alpha.ai/admin/v1/deployments
```

## 架构优势

相比之前的双接口实现，统一接口的优势包括：

1. **简化后端逻辑**：只需要维护一个代理处理器
2. **统一接口**：前端只需要调用一个后端接口 `/api/v1/fusion/`
3. **Header 驱动**：通过 header 灵活控制路由目标
4. **易于扩展**：添加新的 Fusion 服务只需要在 switch 语句中增加 case
5. **减少代码重复**：避免了多个几乎相同的代理处理器

## 安全特性

1. **认证要求**：所有代理请求都需要通过后端的认证中间件
2. **Token 管理**：敏感的 API Token 存储在后端配置中，前端无法访问
3. **请求日志**：所有代理请求都会记录日志，便于审计和调试
4. **错误处理**：统一的错误处理和响应格式
5. **Header 过滤**：`X-Fusion-Provider` header 不会泄露到外部服务
6. **审计日志**：
   - 自动记录所有写操作（PUT/POST）的审计日志
   - 支持 accounts、models、providers 资源的操作追踪
   - 记录操作用户、操作类型、请求内容、响应状态等详细信息
   - 确保操作的可追溯性和合规性

## 实现细节

### 统一代理处理器
- `FusionProxy`：处理所有 Fusion 请求的统一处理器
- 根据 `X-Fusion-Provider` header 选择目标服务

### URL 重写
代理服务会自动处理 URL 重写：
- `/api/v1/fusion/*` → 去掉 `/api/v1` 前缀 → `/fusion/*` → `/admin/v1/*`

### 头部处理
- 复制原始请求的大部分头部
- 跳过 Host、Connection、Upgrade、X-Fusion-Provider 等不应转发的头部
- 根据 provider 添加正确的 Authorization 头部
- 设置正确的 Host 头部

## 部署步骤

1. **配置后端**：
   - 在配置文件中添加 Fusion 相关配置
   - 使用 AES 加密敏感的 Token
   - 重启后端服务

2. **更新前端配置**：
   - **开发环境**：`vite.config.js` 和 `frontend/vite.config.js` 已更新
     - 代理配置指向统一的 `/api/v1/fusion` 接口
     - 自动设置 `X-Fusion-Provider` header 来区分服务
   - **生产环境**：`nginx.conf.template` 已更新，通过设置 `X-Fusion-Provider` header 来区分服务
   - 前端代码保持不变，继续使用 `/fusion-beta/*` 和 `/fusion-alpha/*` 路径

3. **验证功能**：
   - 测试代理转发是否正常
   - 检查认证是否正确添加
   - 验证错误处理

## 部署注意事项

1. 确保配置文件中的 URL 和 Token 正确配置
2. Token 必须经过 AES 加密
3. 重启服务以加载新配置
4. 验证代理功能是否正常工作
5. 确保 Nginx 配置正确设置 `X-Fusion-Provider` header

## 故障排除

### 常见错误
1. **缺少 Provider Header**：检查 Nginx 或 Vite 配置是否正确设置 `X-Fusion-Provider` header
2. **无效的 Provider 值**：确保 header 值为 "beta" 或 "alpha"
3. **配置未找到**：检查配置文件中是否正确添加了 `beta_fusion` 和 `alpha_fusion` 配置
4. **Token 解密失败**：检查 Token 是否正确加密，AES 密钥是否正确
5. **连接失败**：检查目标 URL 是否可访问，网络连接是否正常
6. **认证失败**：检查 Token 是否有效，权限是否足够

### 日志查看
统一代理请求会记录详细的日志，包括：
- Provider 信息
- 原始请求路径和方法
- 转发的目标 URL
- 请求和响应的状态
- 错误信息（如果有）

查看日志示例：
```
Proxying Fusion request to beta: GET /api/v1/fusion/models -> https://api.ppinfra.com/admin/v1/models
Proxying Fusion request to alpha: POST /api/v1/fusion/deployments -> https://api.alpha.ai/admin/v1/deployments
```

## 审计日志功能

### 功能说明
Fusion 代理服务集成了完整的审计日志功能，自动记录所有写操作的详细信息，确保系统操作的可追溯性。

### 审计范围
- **操作类型**：仅记录写操作（PUT、POST 方法）
- **资源类型**：
  - `accounts`：账户管理操作
  - `models`：模型管理操作
  - `providers`：模型提供商管理操作（models 的子资源）

### 审计信息
每条审计日志记录包含：
- 操作用户（从认证上下文获取）
- 操作类型（如"创建账户管理"、"更新模型管理"）
- 请求详情（URL、方法、请求体）
- 响应状态（状态码、成功/失败）
- 目标资源（格式：`{provider}:{resource}`）
- 操作描述

### 路径识别
系统使用正则表达式智能识别 RESTful API 路径：
- `/api/v1/fusion/accounts/*` → accounts 资源
- `/api/v1/fusion/models/*` → models 资源
- `/api/v1/fusion/models/{id}/providers/*` → providers 资源

### 审计日志查看
审计日志存储在系统数据库中，可通过管理界面查看：
- 支持按用户、操作类型、时间范围等条件筛选
- 提供详细的操作记录和请求内容
- 确保敏感信息的适当保护

### 性能考虑
- 审计日志记录不会影响代理请求的性能
- 仅在写操作时触发，读操作无额外开销
- 异步处理，不阻塞主要业务流程