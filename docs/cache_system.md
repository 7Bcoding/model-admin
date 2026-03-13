# 缓存系统使用指南

## 概述

本项目实现了一个通用的缓存系统，用于缓存远程API调用的结果。该系统支持：

- 基于请求URL的缓存键生成
- TTL（生存时间）支持，默认1分钟
- 线程安全操作
- 全局共享缓存（所有用户共享相同的缓存数据）
- 缓存统计和管理
- 详细的缓存命中/未命中日志

## 架构

### 核心组件

1. **CacheService** (`services/cache_service.go`)
   - 单个缓存实例的核心实现
   - 支持Get、Set、Clear操作
   - 自动TTL过期清理
   - 缓存命中时输出日志

2. **CacheManager** (`services/cache_manager.go`)
   - 管理多个缓存实例
   - 提供缓存统计和清理功能
   - 单例模式确保全局一致性

3. **Cache API** (`handlers/cache.go`)
   - 提供HTTP API用于缓存管理
   - 需要管理员权限

## 使用方法

### 在Handler中添加缓存

```go
func YourAPIHandler(w http.ResponseWriter, r *http.Request) {
    // 获取用户信息（仅用于认证）
    _, ok := utils.GetUserFromContext(r.Context())
    if !ok {
        utils.ErrorResponse(w, http.StatusUnauthorized, "User not found")
        return
    }

    // 获取缓存实例
    cacheManager := services.GetCacheManager()
    cache := cacheManager.GetCache("your_cache_name", services.DefaultTTL)
    
    // 生成缓存键（基于URL）
    cacheKey := cache.GenerateKeyFromURL(r.URL.String())
    
    // 尝试从缓存获取
    if cachedData, exists := cache.Get(cacheKey); exists {
        // 缓存命中，会自动输出日志
        utils.SuccessResponse(w, cachedData, "")
        return
    }

    // 缓存未命中，获取新数据
    data, err := yourService.GetData()
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    // 缓存结果
    cache.Set(cacheKey, data)
    
    utils.SuccessResponse(w, data, "")
}
```

### 预定义的缓存配置

在 `services/cache_manager.go` 中定义了以下缓存：

```go
const (
    ModelListCache       = "model_list"        // 模型列表缓存
    ModelInspectionCache = "model_inspection"  // 模型检查缓存
    ModelDeploymentCache = "model_deployment"  // 模型部署缓存

    DefaultTTL         = 1 * time.Minute      // 默认1分钟
    ModelListTTL       = 1 * time.Minute      // 1分钟
    InspectionTTL      = 1 * time.Minute      // 1分钟
    DeploymentTTL      = 1 * time.Minute      // 1分钟
)
```

### 缓存管理API

#### 获取缓存统计
```
GET /api/cache/stats
```

返回所有缓存的统计信息，包括命中率、大小等。

#### 清除特定缓存
```
DELETE /api/cache/clear/{cacheName}
```

清除指定名称的缓存。

#### 清除所有缓存
```
DELETE /api/cache/clear-all
```

清除所有缓存。

**注意：** 所有缓存管理API都需要管理员权限。

## 缓存键生成策略

缓存键基于以下信息生成：
- 完整的请求URL（包括查询参数）
- 使用MD5哈希确保键的唯一性和一致性

**重要特性：**
- **全局共享**：所有用户看到相同的缓存数据
- **环境区分**：不同的URL参数（如platform=alpha vs platform=beta）会生成不同的缓存键
- **一致性保证**：相同的请求URL总是生成相同的缓存键

## 日志输出

系统会输出详细的缓存操作日志：

### 缓存命中
```
缓存命中 - 返回模型列表缓存数据，缓存键: url_abc123...
```

### 缓存未命中
```
缓存未命中 - 获取新的模型列表数据，缓存键: url_abc123...
```

### 缓存设置
```
Cache set for key: url_abc123..., TTL: 1m0s
缓存模型列表响应，缓存键: url_abc123...
```

## 已实现缓存的API

1. **ListModels** - 模型列表API
   - 缓存时间：1分钟
   - 键包含：完整URL（包括platform等参数）

2. **ListModelsInspection** - 模型检查API
   - 缓存时间：1分钟
   - 键包含：完整URL

3. **GetModelDeployment** - 模型部署信息API
   - 缓存时间：1分钟
   - 键包含：完整URL（包括模型名称）

## 扩展指南

要为新的API添加缓存：

1. 在 `cache_manager.go` 中添加缓存名称和TTL常量（可选，也可使用DefaultTTL）
2. 在对应的handler中按照上述模式添加缓存逻辑
3. 确保缓存键包含所有相关的区分参数

## 性能考虑

- 缓存使用内存存储，重启服务会清空所有缓存
- TTL过期的数据会在下次访问时自动清理
- 缓存操作是线程安全的，支持并发访问
- 默认TTL为1分钟，适合快速变化的数据
- 所有用户共享缓存，提高缓存命中率

## 监控和调试

- 所有缓存操作都有详细的中文日志记录
- 缓存命中时会明确输出"缓存命中"日志
- 可以通过 `/api/cache/stats` 监控缓存性能
- 测试文件 `cache_service_test.go` 提供了完整的功能验证

## 设计理念

这个缓存系统采用了"全局共享"的设计理念：
- **一致性**：所有用户看到相同的数据，避免数据不一致
- **效率**：提高缓存命中率，减少重复的远程API调用
- **简单性**：不需要考虑用户隔离，简化了缓存键的生成和管理