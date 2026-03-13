package handlers

import (
	"llm-ops/services"
	"llm-ops/utils"
	"net/http"

	"github.com/go-chi/chi"
)

// GetCacheStats 获取缓存统计信息
func GetCacheStats(w http.ResponseWriter, r *http.Request) {
	// 获取用户信息进行权限检查
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// 只有管理员可以查看缓存统计
	if claims.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	cacheManager := services.GetCacheManager()
	stats := cacheManager.GetAllStats()

	utils.SuccessResponse(w, stats, "Cache statistics retrieved successfully")
}

// ClearCache 清空指定缓存
func ClearCache(w http.ResponseWriter, r *http.Request) {
	// 获取用户信息进行权限检查
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// 只有管理员可以清空缓存
	if claims.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	cacheName := chi.URLParam(r, "cacheName")
	if cacheName == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Cache name is required")
		return
	}

	cacheManager := services.GetCacheManager()
	cacheManager.ClearCache(cacheName)

	utils.SuccessResponse(w, nil, "Cache cleared successfully")
}

// ClearAllCaches 清空所有缓存
func ClearAllCaches(w http.ResponseWriter, r *http.Request) {
	// 获取用户信息进行权限检查
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// 只有管理员可以清空所有缓存
	if claims.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	cacheManager := services.GetCacheManager()
	cacheManager.ClearAllCaches()

	utils.SuccessResponse(w, nil, "All caches cleared successfully")
}