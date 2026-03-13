package handlers

import (
	"encoding/json"
	"fmt"
	"llm-ops/services"
	"llm-ops/utils"
	"log"
	"net/http"
)

var trackerService *services.TrackerService

func InitTrackerService(service *services.TrackerService) {
	trackerService = service
}

// SearchNodes 处理节点搜索请求
func SearchNodes(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Filters []services.RegionFilter `json:"filters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := trackerService.SearchModels(req.Filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// WarmupModels 处理模型预热请求
func WarmupModels(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Models  []services.ModelRequest `json:"models"`
		HFToken string                  `json:"hf_token"`
		NodeIds []string                `json:"node_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range req.Models {
		req.Models[i].Regions = []services.RegionFilter{
			{
				RegionID: "",
				Servers:  req.NodeIds,
			},
		}
	}

	log.Printf("warmup models HFToken: %s, models: %+v, nodeId: %s", req.HFToken, req.Models, req.NodeIds)

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	resp, err := trackerService.WarmupModels(req.Models, req.HFToken, req.NodeIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestBody, _ := json.Marshal(req)

	// 记录审计日志
	detail := fmt.Sprintf("模型: %v", req.Models)
	if err != nil {
		CreateAuditLog(claims.Username, "预热模型", r, http.StatusInternalServerError, "预热失败: "+err.Error(), string(requestBody), fmt.Sprintf("%+v", req.Models), detail)
	} else {
		CreateAuditLog(claims.Username, "预热模型", r, http.StatusOK, "预热成功", string(requestBody), fmt.Sprintf("%+v", req.Models), detail)
	}

	log.Printf("WarmupModels response: %+v", resp.ResponseID)

	// 清理模型相关缓存
	cacheManager := services.GetCacheManager()
	cacheManager.ClearCache(services.ModelListCache)
	cacheManager.ClearCache(services.ModelInspectionCache)
	cacheManager.ClearCache(services.ModelDeploymentCache)
	log.Printf("已清理模型相关缓存（预热模型后）")

	json.NewEncoder(w).Encode(resp)
}

// DeleteModels 处理模型删除请求
func DeleteModels(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Models []services.ModelRequest `json:"models"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	resp, err := trackerService.DeleteModels(req.Models)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteModels response: %+v", resp.ResponseID)

	// 清理模型相关缓存
	cacheManager := services.GetCacheManager()
	cacheManager.ClearCache(services.ModelListCache)
	cacheManager.ClearCache(services.ModelInspectionCache)
	cacheManager.ClearCache(services.ModelDeploymentCache)
	log.Printf("已清理模型相关缓存（删除模型文件后）")

	requestBody, _ := json.Marshal(req)

	// 记录审计日志
	detail := fmt.Sprintf("模型: %v", req.Models)
	if err != nil {
		CreateAuditLog(claims.Username, "删除模型文件", r, http.StatusInternalServerError, "删除失败: "+err.Error(), string(requestBody), fmt.Sprintf("%+v", req.Models), detail)
	} else {
		CreateAuditLog(claims.Username, "删除模型文件", r, http.StatusOK, "删除成功", string(requestBody), fmt.Sprintf("%+v", req.Models), detail)
	}

	json.NewEncoder(w).Encode(resp)
}
