package handlers

import (
	"encoding/json"
	"llm-ops/services"
	"llm-ops/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// 声明为包级变量
var AlphaModelService *services.ModelService
var BetaModelService *services.ModelService

// InitModelService 初始化模型服务
func InitModelService(alphaService *services.ModelService, betaService *services.ModelService) {
	AlphaModelService = alphaService
	BetaModelService = betaService
}

type ModelResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// getModelService 根据请求的platform参数返回对应的模型服务
// platform=beta 返回 BetaModelService，否则返回 AlphaModelService
func getModelService(r *http.Request) *services.ModelService {
	// get platform from query params
	platform := r.URL.Query().Get("platform")

	// if platform is beta, return beta service
	if platform == "beta" {
		return BetaModelService
	}

	// default return alpha service
	return AlphaModelService
}

func ListModelsInspection(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling ListModelsInspection Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	// 如果是 OPTIONS 请求，直接返回
	if r.Method == "OPTIONS" {
		return
	}

	modelService := getModelService(r)

	// 使用 ModelService 的缓存功能
	response, err := modelService.ListModelsInspection(r.URL.String())
	if err != nil {
		log.Printf("Error fetching models inspection: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully retrieved models inspection data")

	// 返回数据
	utils.SuccessResponse(w, response, "")
}

func ListModels(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling ListModels Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	// 如果是 OPTIONS 请求，直接返回
	if r.Method == "OPTIONS" {
		return
	}

	// 获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		log.Printf("No user found in context")
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	log.Printf("User from context - ID: %d, Username: %s, Role: %s",
		claims.UserID, claims.Username, claims.Role)

	modelService := getModelService(r)

	// 使用 ModelService 的缓存功能
	response, err := modelService.ListModelsWithUserData(r.URL.String(), claims.UserID)
	if err != nil {
		log.Printf("Error fetching models with user data: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully retrieved models with user data")

	// 返回数据
	utils.SuccessResponse(w, response, "")
}

func StarModel(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling StarModel Request ===")

	var req struct {
		ModelName string `json:"model_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("Request to star model: %s", req.ModelName)

	// 从 context 中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		log.Printf("No user found in context")
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	log.Printf("User from context - ID: %d, Username: %s", claims.UserID, claims.Username)

	// 检查 modelService 是否已初始化
	if AlphaModelService == nil {
		log.Printf("modelService is nil")
		utils.ErrorResponse(w, http.StatusInternalServerError, "Service not initialized")
		return
	}

	// 调用 modelService 收藏模型
	err := AlphaModelService.StarModel(claims.UserID, req.ModelName)
	if err != nil {
		log.Printf("Error starring model: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully starred model %s for user %d", req.ModelName, claims.UserID)
	
	// 清理模型相关缓存（用户数据变更）
	cacheManager := services.GetCacheManager()
	cacheManager.ClearCache(services.ModelListCache)
	log.Printf("已清理模型列表缓存（收藏模型后）")
	
	utils.SuccessResponse(w, nil, "Successfully starred model")
}

func UnstarModel(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling UnstarModel Request ===")

	var req struct {
		ModelName string `json:"model_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("Request to unstar model: %s", req.ModelName)

	// 从 context 中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		log.Printf("No user found in context")
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	log.Printf("User from context - ID: %d, Username: %s", claims.UserID, claims.Username)

	// 检查 modelService 是否已初始化
	if AlphaModelService == nil {
		log.Printf("modelService is nil")
		utils.ErrorResponse(w, http.StatusInternalServerError, "Service not initialized")
		return
	}

	// 调用 modelService 取消收藏模型
	err := AlphaModelService.UnstarModel(claims.UserID, req.ModelName)
	if err != nil {
		log.Printf("Error unstarring model: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully unstarred model %s for user %d", req.ModelName, claims.UserID)
	
	// 清理模型相关缓存（用户数据变更）
	cacheManager := services.GetCacheManager()
	cacheManager.ClearCache(services.ModelListCache)
	log.Printf("已清理模型列表缓存（取消收藏模型后）")
	
	utils.SuccessResponse(w, nil, "Successfully unstarred model")
}

func UpdateModelNote(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling UpdateModelNote Request ===")

	var req struct {
		ModelName       string `json:"model_name"`
		Note            string `json:"note"`
		OpenChatId      string `json:"open_chat_id"`
		InferenceEngine string `json:"inference_engine"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("Request to update note for model: %s", req.ModelName)

	// 从 context 中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		log.Printf("No user found in context")
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	log.Printf("User from context - ID: %d, Username: %s", claims.UserID, claims.Username)

	// 检查 modelService 是否已初始化
	if AlphaModelService == nil {
		log.Printf("modelService is nil")
		utils.ErrorResponse(w, http.StatusInternalServerError, "Service not initialized")
		return
	}

	// 调用 modelService 更新笔记和飞书群组ID
	err := AlphaModelService.UpdateNote(req.ModelName, req.Note, req.OpenChatId, req.InferenceEngine)
	if err != nil {
		log.Printf("Error updating note: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully updated note and metadata for model %s by user %d", req.ModelName, claims.UserID)
	
	// 清理模型相关缓存（用户数据变更）
	cacheManager := services.GetCacheManager()
	cacheManager.ClearCache(services.ModelListCache)
	log.Printf("已清理模型列表缓存（更新模型笔记后）")
	
	utils.SuccessResponse(w, nil, "Successfully updated note and metadata")
}

func GetModelDeployment(w http.ResponseWriter, r *http.Request) {
	modelName := chi.URLParam(r, "modelName")

	// 获取用户信息 - 仅做认证检查
	_, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// 获取模型部署信息
	modelService := getModelService(r)
	deploymentInfo, err := modelService.GetModelDeploymentWithCache(modelName, r.URL.String())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, deploymentInfo, "")
}
