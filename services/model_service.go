package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"llm-ops/config"
	"llm-ops/db"
	"log"
	"net/http"
)

type ModelService struct {
	apiHost      string
	apiKey       string
	cacheManager *CacheManager
}

func NewModelService() *ModelService {
	return &ModelService{
		apiHost:      config.Config.Alpha.ApiURL,
		apiKey:       config.Config.Alpha.ApiKey,
		cacheManager: GetCacheManager(),
	}
}

func NewBetaModelService() *ModelService {
	return &ModelService{
		apiHost:      config.Config.Beta.ApiURL,
		apiKey:       config.Config.Beta.ApiKey,
		cacheManager: GetCacheManager(),
	}
}

func (s *ModelService) GetUserStars(userID int) (map[string]bool, error) {
	stars := make(map[string]bool)
	var modelStars []struct{ ModelName string }
	err := db.DB.Table("model_stars").Select("model_name").Where("user_id = ?", userID).Find(&modelStars).Error
	if err != nil {
		return nil, err
	}
	for _, star := range modelStars {
		stars[star.ModelName] = true
	}
	return stars, nil
}

func (s *ModelService) GetUserNotes(userID int) (map[string]string, map[string]string, map[string]string, error) {
	notes := make(map[string]string)
	openChatIds := make(map[string]string)
	engines := make(map[string]string)
	rows, err := db.DB.Table("model_notes").Select("model_name, note, open_chat_id, inference_engine").Rows()
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var modelName, note, openChatId, inferenceEngine string
		if err := rows.Scan(&modelName, &note, &openChatId, &inferenceEngine); err != nil {
			return nil, nil, nil, err
		}
		notes[modelName] = note
		openChatIds[modelName] = openChatId
		engines[modelName] = inferenceEngine
	}
	return notes, openChatIds, engines, nil
}

func (s *ModelService) StarModel(userID int, modelName string) error {
	return db.DB.Exec("INSERT INTO model_stars (user_id, model_name) VALUES (?, ?)", userID, modelName).Error
}

func (s *ModelService) UnstarModel(userID int, modelName string) error {
	return db.DB.Exec("DELETE FROM model_stars WHERE user_id = ? AND model_name = ?", userID, modelName).Error
}

func (s *ModelService) UpdateNote(modelName, note, openChatId, inferenceEngine string) error {
	query := `
		INSERT INTO model_notes (model_name, note, open_chat_id, inference_engine) 
		VALUES (?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE note = ?, open_chat_id = ?, inference_engine = ?
	`

	err := db.DB.Exec(query, modelName, note, openChatId, inferenceEngine, note, openChatId, inferenceEngine).Error
	if err != nil {
		log.Printf("Error updating note: %v", err)
		return fmt.Errorf("failed to update note: %v", err)
	}

	return nil
}

func (s *ModelService) listModelsFromAPI() ([]LLMModelInfo, error) {
	// 生成缓存键 - 基于 API URL
	cacheKey := fmt.Sprintf("models_api_%s", s.apiHost)
	
	// 尝试从缓存获取数据
	cache := s.cacheManager.GetCache(ModelListCache, ModelListTTL)
	if cachedData, exists := cache.Get(cacheKey); exists {
		log.Printf("缓存命中 - 返回模型 API 缓存数据，缓存键: %s", cacheKey)
		if models, ok := cachedData.([]LLMModelInfo); ok {
			return models, nil
		}
		// 如果类型转换失败，清除这个缓存项并继续从 API 获取
		cache.Delete(cacheKey)
		log.Printf("缓存数据类型错误，已清除缓存项: %s", cacheKey)
	}

	log.Printf("缓存未命中 - 从 API 获取新的模型数据，缓存键: %s", cacheKey)

	url := fmt.Sprintf("%s/v3/admin/llm/model", s.apiHost)
	log.Printf("Requesting models from URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// log.Printf("Request headers: %+v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 检查返回状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// 解析返回的JSON到结构化数据
	var modelResp ModelResponse
	if err := json.Unmarshal(body, &modelResp); err != nil {
		return nil, fmt.Errorf("failed to parse model response: %v", err)
	}

	// 缓存结果
	cache.Set(cacheKey, modelResp.Models)
	log.Printf("缓存模型 API 响应，缓存键: %s，模型数量: %d", cacheKey, len(modelResp.Models))

	return modelResp.Models, nil
}

func (s *ModelService) GetModelDeployment(modelName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v3/admin/llm/model/%s/deployment", s.apiHost, modelName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result, nil
}

func (s *ModelService) GetModelEndpoints(modelName string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v3/admin/llm/model/%s/endpoints", s.apiHost, modelName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Endpoints []map[string]interface{} `json:"endpoints"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Endpoints, nil
}

func (s *ModelService) TestEndpoint(endpointURL, modelName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v3/admin/llm/endpoint/test", s.apiHost)

	reqBody := map[string]string{
		"endpoint_url": endpointURL,
		"model_name":   modelName,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result, nil
}

func (s *ModelService) UpdateEndpoint(modelName string, endpointID string, weight int, enableCheckHealth bool, enableTestflight bool, supportedApiFlag string) error {
	log.Printf("Updating endpoint - Model: %s, EndpointID: %s", modelName, endpointID)

	url := fmt.Sprintf("%s/v3/admin/llm/endpoint", s.apiHost)
	log.Printf("Request URL: %s", url)

	data := map[string]interface{}{
		"endpoint_id": endpointID,
	}
	data["new_data"] = map[string]interface{}{
		"weight":              weight,
		"enable_check_health": enableCheckHealth,
		"enable_testflight":   enableTestflight,
		"supported_api_flag":  supportedApiFlag,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling request data: %v", err)
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	log.Printf("Request payload: %s", string(jsonData))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	// log.Printf("Request headers: %+v", req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully updated endpoint")
	
	// 清理模型相关缓存
	s.ClearModelCaches()
	
	return nil
}

func (s *ModelService) StopEndpoint(endpointID string) error {
	log.Printf("Stopping endpoint: %s", endpointID)

	url := fmt.Sprintf("%s/v3/admin/llm/endpoint", s.apiHost)
	log.Printf("Request URL: %s", url)

	data := map[string]interface{}{
		"endpoint_id": endpointID,
	}
	data["new_data"] = map[string]interface{}{
		"status": "ENDPOINT_STATUS_DELETED",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling request data: %v", err)
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	log.Printf("Request payload: %s", string(jsonData))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	// log.Printf("Request headers: %+v", req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully stopped endpoint")
	
	// 清理模型相关缓存
	s.ClearModelCaches()
	
	return nil
}

func (s *ModelService) AddEndpoint(modelName string, endpointID string, epurl string, weight int, modelIDOverride string, enableCheckHealth bool,
	enableTestflight bool, supportedApiFlag string, skipTestflightWhileAdd bool) error {
	log.Printf("Adding endpoint - Model: %s, EndpointID: %s", modelName, endpointID)

	// 构建请求 URL
	url := fmt.Sprintf("%s/v3/admin/llm/endpoint", s.apiHost)
	log.Printf("Request URL: %s", url)

	// 构建请求数据
	endpoint := map[string]interface{}{
		"model_name":                modelName,
		"endpoint_id":               endpointID,
		"url":                       epurl,
		"weight":                    weight,
		"check_health_url":          "/health",
		"provider":                  "alpha-serverless",
		"status":                    "ENDPOINT_STATUS_SERVING",
		"supported_api_flag":        supportedApiFlag,
		"skip_testflight_while_add": skipTestflightWhileAdd,
	}

	if modelIDOverride != "" {
		endpoint["model_name_override"] = modelIDOverride
	}

	endpoint["enable_check_health"] = enableCheckHealth
	endpoint["enable_testflight"] = enableTestflight

	data := map[string]interface{}{
		"endpoint": endpoint,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling request data: %v", err)
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	log.Printf("Request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	// log.Printf("Request headers: %+v", req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully added endpoint")
	
	// 清理模型相关缓存
	s.ClearModelCaches()
	
	return nil
}

// ListModels 获取模型列表（底层已有缓存）
func (s *ModelService) ListModels(requestURL string) ([]LLMModelInfo, error) {
	return s.listModelsFromAPI()
}

// ListModelsWithUserData 获取包含用户数据的模型列表
func (s *ModelService) ListModelsWithUserData(url string, userID int) (map[string]interface{}, error) {
	// 获取基础模型列表（已在 listModelsFromAPI 中缓存）
	modelsResponse, err := s.listModelsFromAPI()
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %v", err)
	}

	// 获取用户的收藏信息
	stars, err := s.GetUserStars(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user stars: %v", err)
	}

	// 获取用户的模型备注信息
	notes, openChatIds, engines, err := s.GetUserNotes(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user notes: %v", err)
	}

	// 处理模型数据
	var processedModels []map[string]interface{}
	for _, item := range modelsResponse {
		modelName := item.Model.ModelName
		if engines[modelName] == "" {
			engines[modelName] = "vllm"
		}

		processedModel := map[string]interface{}{
			"model":              item.Model,
			"starred":            stars[modelName],
			"note":               notes[modelName],
			"endpoints":          item.Endpoints,
			"model_name":         modelName,
			"private":            item.Model.IsPrivate,
			"max_tokens":         item.Model.ContextSize,
			"description":        item.Model.Description,
			"status":             item.Model.Status,
			"tags":               item.Model.Tags,
			"input_token_price":  item.Model.InputTokenPricePerMillion,
			"output_token_price": item.Model.OutputTokenPricePerMillion,
			"rank":               item.Model.Rank,
			"open_chat_id":       openChatIds[modelName],
			"inference_engine":   engines[modelName],
		}
		processedModels = append(processedModels, processedModel)
	}

	response := map[string]interface{}{
		"data": processedModels,
	}

	return response, nil
}

// ListModelsInspection 获取模型检查数据（带缓存）
func (s *ModelService) ListModelsInspection(url string) ([]map[string]interface{}, error) {
	// 对检查数据进行缓存
	cache := s.cacheManager.GetCache(ModelInspectionCache, ModelListTTL)
	cacheKey := cache.GenerateKeyFromURL(url)
	
	// 尝试从缓存获取数据
	if cachedData, exists := cache.Get(cacheKey); exists {
		log.Printf("缓存命中 - 返回模型检查缓存数据，缓存键: %s", cacheKey)
		return cachedData.([]map[string]interface{}), nil
	}

	log.Printf("缓存未命中 - 获取新的模型检查数据，缓存键: %s", cacheKey)

	// 获取模型列表（已缓存）
	modelsResponse, err := s.listModelsFromAPI()
	if err != nil {
		return nil, err
	}

	// 过滤可用模型
	availableModels := []LLMModelInfo{}
	for _, model := range modelsResponse {
		if model.Model.Status == "MODEL_STATUS_SERVING" {
			availableModels = append(availableModels, model)
		}
	}

	// 处理模型数据
	processedModels := []map[string]interface{}{}
	for _, item := range availableModels {
		modelName := item.Model.ModelName
		endpoints := make([]map[string]interface{}, 0)
		for _, endpoint := range item.Endpoints {
			if endpoint.Weight == 0 {
				continue
			}
			override_name := endpoint.ModelNameOverride
			if override_name == "" {
				override_name = item.Model.ModelName
			}
			endpoints = append(endpoints, map[string]interface{}{
				"url":        endpoint.URL,
				"model_name": override_name,
				"weight":     endpoint.Weight,
			})
		}
		processedModel := map[string]interface{}{
			"endpoints":  endpoints,
			"model_name": modelName,
			"max_tokens": item.Model.ContextSize,
		}
		processedModels = append(processedModels, processedModel)
	}

	// 缓存结果
	cache.Set(cacheKey, processedModels)
	log.Printf("缓存模型检查响应，缓存键: %s", cacheKey)

	return processedModels, nil
}

// GetModelDeploymentWithCache 带缓存的模型部署信息获取方法
func (s *ModelService) GetModelDeploymentWithCache(modelName, url string) (map[string]interface{}, error) {
	// 获取缓存实例
	cache := s.cacheManager.GetCache(ModelDeploymentCache, DeploymentTTL)
	cacheKey := cache.GenerateKeyFromURL(url)
	
	// 尝试从缓存获取数据
	if cachedData, exists := cache.Get(cacheKey); exists {
		log.Printf("缓存命中 - 返回模型部署缓存数据，模型: %s，缓存键: %s", modelName, cacheKey)
		return cachedData.(map[string]interface{}), nil
	}

	log.Printf("缓存未命中 - 获取新的模型部署数据，模型: %s，缓存键: %s", modelName, cacheKey)

	// 从API获取数据
	deployment, err := s.GetModelDeployment(modelName)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	cache.Set(cacheKey, deployment)
	log.Printf("缓存模型部署响应，模型: %s，缓存键: %s", modelName, cacheKey)

	return deployment, nil
}

// ClearModelCaches 清理所有模型相关的缓存
func (s *ModelService) ClearModelCaches() {
	log.Printf("清理所有模型相关缓存")
	
	// 清理模型列表缓存
	s.cacheManager.ClearCache(ModelListCache)
	log.Printf("已清理模型列表缓存")
	
	// 清理模型检查缓存
	s.cacheManager.ClearCache(ModelInspectionCache)
	log.Printf("已清理模型检查缓存")
	
	// 清理模型部署缓存
	s.cacheManager.ClearCache(ModelDeploymentCache)
	log.Printf("已清理模型部署缓存")
}
