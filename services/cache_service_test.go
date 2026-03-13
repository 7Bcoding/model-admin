package services

import (
	"testing"
	"time"
)

func TestCacheService(t *testing.T) {
	// 创建缓存服务实例
	cache := NewCacheService(1 * time.Second) // 1秒TTL用于测试

	// 测试数据
	testKey := "test_key"
	testValue := map[string]interface{}{
		"message": "hello world",
		"number":  42,
	}

	// 测试Set和Get
	cache.Set(testKey, testValue)
	
	if value, exists := cache.Get(testKey); exists {
		if data, ok := value.(map[string]interface{}); ok {
			if data["message"] != "hello world" {
				t.Errorf("Expected 'hello world', got %v", data["message"])
			}
		} else {
			t.Error("Failed to cast cached value to expected type")
		}
	} else {
		t.Error("Expected cached value to exist")
	}

	// 测试TTL过期
	time.Sleep(2 * time.Second)
	if _, exists := cache.Get(testKey); exists {
		t.Error("Expected cached value to be expired")
	}

	// 测试Clear
	cache.Set(testKey, testValue)
	cache.Clear()
	if _, exists := cache.Get(testKey); exists {
		t.Error("Expected cache to be cleared")
	}
}

func TestGenerateKeyFromURL(t *testing.T) {
	cache := NewCacheService(5 * time.Minute)
	
	// 测试URL key生成
	url1 := "/api/models?platform=novita"
	url2 := "/api/models?platform=ppio"
	
	key1 := cache.GenerateKeyFromURL(url1)
	key2 := cache.GenerateKeyFromURL(url2)
	
	if key1 == key2 {
		t.Error("Different URLs should generate different cache keys")
	}
	
	// 相同URL应该生成相同的key
	key3 := cache.GenerateKeyFromURL(url1)
	if key1 != key3 {
		t.Error("Same URL should generate same cache key")
	}
}

func TestCacheManager(t *testing.T) {
	manager := GetCacheManager()
	
	// 获取缓存实例
	cache1 := manager.GetCache("test_cache", 1*time.Minute)
	cache2 := manager.GetCache("test_cache", 1*time.Minute)
	
	// 应该返回相同的实例
	if cache1 != cache2 {
		t.Error("Expected same cache instance for same cache name")
	}
	
	// 测试缓存操作
	testKey := "manager_test"
	testValue := "test_value"
	
	cache1.Set(testKey, testValue)
	
	if value, exists := cache2.Get(testKey); !exists || value != testValue {
		t.Error("Cache instances should share data")
	}
	
	// 测试清除缓存
	manager.ClearCache("test_cache")
	if _, exists := cache1.Get(testKey); exists {
		t.Error("Expected cache to be cleared")
	}
}