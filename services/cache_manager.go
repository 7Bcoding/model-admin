package services

import (
	"sync"
	"time"
)

// CacheManager 缓存管理器
type CacheManager struct {
	caches map[string]*CacheService
	mutex  sync.RWMutex
}

// 全局缓存管理器实例
var globalCacheManager *CacheManager
var cacheManagerOnce sync.Once

// GetCacheManager 获取全局缓存管理器实例（单例模式）
func GetCacheManager() *CacheManager {
	cacheManagerOnce.Do(func() {
		globalCacheManager = &CacheManager{
			caches: make(map[string]*CacheService),
		}
	})
	return globalCacheManager
}

// GetCache 获取指定名称的缓存服务，如果不存在则创建
func (cm *CacheManager) GetCache(name string, ttl time.Duration) *CacheService {
	cm.mutex.RLock()
	cache, exists := cm.caches[name]
	cm.mutex.RUnlock()

	if exists {
		return cache
	}

	// 缓存不存在，创建新的
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 双重检查，防止并发创建
	if cache, exists := cm.caches[name]; exists {
		return cache
	}

	cache = NewCacheService(ttl)
	cm.caches[name] = cache
	return cache
}

// ClearCache 清空指定名称的缓存
func (cm *CacheManager) ClearCache(name string) {
	cm.mutex.RLock()
	cache, exists := cm.caches[name]
	cm.mutex.RUnlock()

	if exists {
		cache.Clear()
	}
}

// ClearAllCaches 清空所有缓存
func (cm *CacheManager) ClearAllCaches() {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for _, cache := range cm.caches {
		cache.Clear()
	}
}

// GetAllStats 获取所有缓存的统计信息
func (cm *CacheManager) GetAllStats() map[string]interface{} {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := make(map[string]interface{})
	for name, cache := range cm.caches {
		stats[name] = cache.GetStats()
	}

	return stats
}

const (
	// 缓存名称常量
	ModelListCache       = "model_list"
	ModelInspectionCache = "model_inspection"
	ModelDeploymentCache = "model_deployment"

	// 缓存TTL常量 - 默认1分钟
	DefaultTTL    = 5 * time.Minute
	ModelListTTL  = 5 * time.Minute
	InspectionTTL = 1 * time.Minute
	DeploymentTTL = 1 * time.Minute
)
