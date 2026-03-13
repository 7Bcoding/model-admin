package services

import (
	"crypto/md5"
	"fmt"
	"log"
	"sync"
	"time"
)

// CacheItem 缓存项结构
type CacheItem struct {
	Data      interface{} `json:"data"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// CacheService 通用缓存服务
type CacheService struct {
	cache map[string]*CacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewCacheService 创建新的缓存服务实例
func NewCacheService(ttl time.Duration) *CacheService {
	service := &CacheService{
		cache: make(map[string]*CacheItem),
		ttl:   ttl,
	}
	
	// 启动清理过期缓存的goroutine
	go service.cleanupExpiredItems()
	
	return service
}

// GenerateKeyFromURL 根据URL生成缓存键（不包含用户ID，所有用户共享缓存）
func (c *CacheService) GenerateKeyFromURL(requestURL string) string {
	// 使用URL的哈希作为缓存键，不包含用户ID
	hash := md5.Sum([]byte(requestURL))
	return fmt.Sprintf("url_%x", hash)
}

// Get 从缓存中获取数据
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.cache[key]
	if !exists {
		return nil, false
	}
	
	// 检查是否过期
	if time.Now().After(item.ExpiresAt) {
		// 过期了，删除并返回false
		delete(c.cache, key)
		return nil, false
	}
	
	log.Printf("Cache hit for key: %s", key)
	return item.Data, true
}

// Set 设置缓存数据
func (c *CacheService) Set(key string, data interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.cache[key] = &CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	
	log.Printf("Cache set for key: %s, TTL: %v", key, c.ttl)
}

// Delete 删除缓存数据
func (c *CacheService) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.cache, key)
	log.Printf("Cache deleted for key: %s", key)
}

// Clear 清空所有缓存
func (c *CacheService) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.cache = make(map[string]*CacheItem)
	log.Printf("Cache cleared")
}

// GetStats 获取缓存统计信息
func (c *CacheService) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	totalItems := len(c.cache)
	expiredItems := 0
	
	now := time.Now()
	for _, item := range c.cache {
		if now.After(item.ExpiresAt) {
			expiredItems++
		}
	}
	
	return map[string]interface{}{
		"total_items":   totalItems,
		"expired_items": expiredItems,
		"active_items":  totalItems - expiredItems,
		"ttl_seconds":   c.ttl.Seconds(),
	}
}

// cleanupExpiredItems 定期清理过期的缓存项
func (c *CacheService) cleanupExpiredItems() {
	ticker := time.NewTicker(time.Minute * 5) // 每5分钟清理一次
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			now := time.Now()
			deletedCount := 0
			
			for key, item := range c.cache {
				if now.After(item.ExpiresAt) {
					delete(c.cache, key)
					deletedCount++
				}
			}
			
			if deletedCount > 0 {
				log.Printf("Cleaned up %d expired cache items", deletedCount)
			}
			c.mutex.Unlock()
		}
	}
}

// CacheWrapper 缓存包装器函数类型
type CacheWrapper func() (interface{}, error)

// GetOrSet 获取缓存数据，如果不存在则执行函数并缓存结果
func (c *CacheService) GetOrSet(key string, fn CacheWrapper) (interface{}, error) {
	// 先尝试从缓存获取
	if data, exists := c.Get(key); exists {
		return data, nil
	}
	
	// 缓存不存在，执行函数获取数据
	data, err := fn()
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	c.Set(key, data)
	
	return data, nil
}