package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"llm-ops/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type EndpointRequest struct {
	ModelName         string `json:"model_name"`
	Weight            int    `json:"weight"`
	EndpointID        string `json:"endpoint_id"`
	EnableCheckHealth bool   `json:"enable_check_health"`
	EnableTestflight  bool   `json:"enable_testflight"`
	SupportedApiFlag  string `json:"supported_api_flag"`

	URL                    string `json:"url"`
	ModelIDOverride        string `json:"model_id_override"`
	SkipTestflightWhileAdd bool   `json:"skip_testflight_while_add"`
}

func (req *EndpointRequest) BuildDetail() string {
	return fmt.Sprintf("ModelName: %s, EndpointID: %s, URL: %s, Weight: %d, ModelIDOverride: %s",
		req.ModelName, req.EndpointID, req.URL, req.Weight, req.ModelIDOverride)
}

// UpdateEndpoint 处理更新 endpoint 的请求
func UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling UpdateEndpoint Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 获取请求体内容用于审计日志
	var requestBody []byte
	var err error
	requestBody, err = utils.ReadRequestBody(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "读取请求体失败")
		return
	}

	vars := mux.Vars(r)
	endpointID := vars["id"]
	log.Printf("Endpoint ID from URL: %s", endpointID)

	var req EndpointRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		CreateAuditLog(claims.Username, "更新Endpoint", r, http.StatusBadRequest, "请求体解析失败", string(requestBody), endpointID, req.BuildDetail())
		return
	}
	req.EndpointID = endpointID

	log.Printf("Request data - Model: %s, Weight: %d, CheckHealth: %v, Testflight: %v",
		req.ModelName, req.Weight, req.EnableCheckHealth, req.EnableTestflight)

	modelService := getModelService(r)
	err = modelService.UpdateEndpoint(
		req.ModelName,
		endpointID,
		req.Weight,
		req.EnableCheckHealth,
		req.EnableTestflight,
		req.SupportedApiFlag,
	)

	if err != nil {
		log.Printf("Error updating endpoint: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		CreateAuditLog(claims.Username, "更新Endpoint", r, http.StatusInternalServerError, "更新失败: "+err.Error(), string(requestBody), req.ModelName, req.BuildDetail())
		return
	}

	log.Printf("Successfully updated endpoint %s for model %s, supported_api_flag: %s", endpointID, req.ModelName, req.SupportedApiFlag)
	CreateAuditLog(claims.Username, "更新Endpoint", r, http.StatusOK, "更新成功", string(requestBody), req.ModelName, req.BuildDetail())
	utils.SuccessResponse(w, nil, "Successfully updated endpoint")
}

// StopEndpoint 处理停用 endpoint 的请求
func StopEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling StopEndpoint Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 获取请求体内容用于审计日志
	var (
		requestBody []byte
		err         error
	)
	requestBody, err = utils.ReadRequestBody(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "读取请求体失败")
		return
	}

	var req EndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		CreateAuditLog(claims.Username, "停用Endpoint", r, http.StatusBadRequest, "请求体解析失败", string(requestBody), req.EndpointID, req.BuildDetail())
		return
	}

	log.Printf("Request to stop endpoint: %s", req.EndpointID)

	modelService := getModelService(r)
	err = modelService.StopEndpoint(req.EndpointID)
	if err != nil {
		log.Printf("Error stopping endpoint: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		CreateAuditLog(claims.Username, "停用Endpoint", r, http.StatusInternalServerError, "停用失败: "+err.Error(), string(requestBody), req.EndpointID, req.BuildDetail())
		return
	}

	log.Printf("Successfully stopped endpoint: %s", req.EndpointID)
	CreateAuditLog(claims.Username, "停用Endpoint", r, http.StatusOK, "停用成功", string(requestBody), req.EndpointID, req.BuildDetail())
	utils.SuccessResponse(w, nil, "Successfully stopped endpoint")
}

// AddEndpoint 处理添加 endpoint 的请求
func AddEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling AddEndpoint Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 获取请求体内容用于审计日志
	var (
		requestBody []byte
		err         error
	)
	requestBody, err = utils.ReadRequestBody(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "读取请求体失败")
		return
	}

	var req EndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		CreateAuditLog(claims.Username, "添加Endpoint", r, http.StatusBadRequest, "请求体解析失败", string(requestBody), req.EndpointID, req.BuildDetail())
		return
	}

	log.Printf("Request data: %+v", req)

	modelService := getModelService(r)
	err = modelService.AddEndpoint(
		req.ModelName,
		req.EndpointID,
		req.URL,
		req.Weight,
		req.ModelIDOverride,
		req.EnableCheckHealth,
		req.EnableTestflight,
		req.SupportedApiFlag,
		req.SkipTestflightWhileAdd,
	)

	if err != nil {
		log.Printf("Error adding endpoint: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		CreateAuditLog(claims.Username, "添加Endpoint", r, http.StatusInternalServerError, "添加失败: "+err.Error(), string(requestBody), req.EndpointID, req.BuildDetail())
		return
	}

	log.Printf("Successfully added endpoint %s for model %s", req.EndpointID, req.ModelName)
	CreateAuditLog(claims.Username, "添加Endpoint", r, http.StatusOK, "添加成功", string(requestBody), req.EndpointID, req.BuildDetail())
	utils.SuccessResponse(w, nil, "Successfully added endpoint")
}

// GetEndpointMetrics 代理获取 endpoint metrics 的请求
func GetEndpointMetrics(w http.ResponseWriter, r *http.Request) {
	seName := r.URL.Query().Get("endpoint")
	if seName == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "endpoint parameter is required")
		return
	}

	log.Printf("=== Handling GetEndpointMetrics Request for endpoint: %s ===", seName)
	// 构建请求 URL
	url := fmt.Sprintf("https://metric.novita.dev/metrics?endpoint=%s", seName)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
