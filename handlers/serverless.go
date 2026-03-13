package handlers

import (
	"encoding/json"
	"fmt"
	"llm-ops/models"
	"llm-ops/serverless/v1beta1"
	"llm-ops/services"
	"llm-ops/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var k8sService *services.KubernetesService

func InitK8sService(ks *services.KubernetesService) {
	k8sService = ks
}

// ListEndpoints 获取所有 Endpoint 资源
func ListEndpoints(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Listing endpoints in namespace: %s", namespace)
	endpoints, err := k8sService.ListEndpoints(namespace)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	for i := range endpoints {
		workers, err := k8sService.ListWorkers(namespace, endpoints[i].Metadata.Name)
		if err == nil {
			endpoints[i].WorkerCount = len(workers)
		}
	}

	modelService := getModelService(r)

	// get model lists
	// 获取模型列表
	modelsResponse, err := modelService.ListModels(r.URL.String())
	if err != nil {
		log.Printf("Error fetching models: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	extractEndpointFromURL := func(url string) string {
		if url == "" {
			return ""
		}
		// 移除 https:// 或 http:// 前缀
		url = strings.TrimPrefix(url, "https://")
		url = strings.TrimPrefix(url, "http://")
		if matches := strings.Split(url, ".runsync.novita.dev"); len(matches) > 0 {
			name := matches[0]
			return name
		}
		return ""
	}

	modelsByEndpoint := make(map[string][]string)
	for _, model := range modelsResponse {
		for _, endpoint := range model.Endpoints {
			endpointName := extractEndpointFromURL(endpoint.URL)
			// log.Printf("Endpoint name: %s, model name: %s", endpointName, model.Model.ModelName)
			if endpointName == "" {
				continue
			}
			modelsByEndpoint[endpointName] = append(modelsByEndpoint[endpointName], model.Model.ModelName)
		}
	}

	for i, endpoint := range endpoints {
		endpoints[i].Models = modelsByEndpoint[endpoint.Metadata.Name]
	}

	utils.SuccessResponse(w, endpoints, "")
}

// GetEndpoint 获取指定的 Endpoint 资源
func GetEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Getting endpoint %s in namespace: %s", name, namespace)
	endpoint, err := k8sService.GetEndpoint(namespace, name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	saps, err := k8sService.ListServerlessPolicies(namespace, "")
	if err != nil {
		log.Printf("Error listing SAPs: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, sap := range saps {
		if sap.Spec.ScaleTargetRef.Name != name {
			continue
		}
		sapYaml, err := yaml.Marshal(sap)
		if err != nil {
			log.Printf("Error marshaling SAP to YAML: %v", err)
			continue
		}
		endpoint.SapYaml = string(sapYaml)
		endpoint.SapParam = models.SapParamFromSap(sap)
		break
	}

	utils.SuccessResponse(w, endpoint, "")
}

// UpdateEndpointImage 更新 Endpoint 镜像
func UpdateEndpointImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	var req struct {
		Image string `json:"image"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		klog.Errorf("更新Endpoint镜像失败: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 验证镜像格式
	if !strings.Contains(req.Image, ":") {
		klog.Errorf("镜像格式错误,必须包含tag: %s", req.Image)
		utils.ErrorResponse(w, http.StatusBadRequest, "镜像格式错误,必须包含tag")
		return
	}

	se, err := k8sService.GetEndpoint(namespace, name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if se.Spec.Template.Spec.Containers[0].Image == req.Image {
		utils.ErrorResponse(w, http.StatusBadRequest, "镜像已是最新")
		return
	}

	// 记录审计日志
	requestBody, _ := json.Marshal(req)
	detail := fmt.Sprintf("更新镜像为: %s", req.Image)
	CreateAuditLog(claims.Username, "更新Endpoint镜像", r, http.StatusOK, "更新成功", string(requestBody), name, detail)

	err = k8sService.UpdateEndpointImage(namespace, se.Metadata.Name, req.Image)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, se, "")
}

// UpdateEndpointMaxConcurrency 更新 Endpoint 最大并发数
func UpdateEndpointMaxConcurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	var req struct {
		MaxConcurrency int `json:"maxConcurrency"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.MaxConcurrency <= 0 || req.MaxConcurrency > 1000 {
		utils.ErrorResponse(w, http.StatusBadRequest, "maxConcurrency must be between 1 and 1000")
		return
	}

	err = k8sService.UpdateEndpointMaxConcurrency(namespace, name, req.MaxConcurrency)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, nil, "")
}

// ListClusters 获取所有 Cluster 资源
func ListClusters(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Listing clusters in namespace: %s", namespace)
	clusters, err := k8sService.ListClusters(namespace)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, clusters, "")
}

// GetCluster 获取指定的 Cluster 资源
func GetCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Getting cluster %s in namespace: %s", name, namespace)
	cluster, err := k8sService.GetCluster(namespace, name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"data": cluster,
	}, "")
}

// 添加辅助函数来从 URL 提取 SE 名称
func extractSEFromURL(url string) string {
	if url == "" {
		return ""
	}
	// 匹配 https://{name}.runsync.novita.dev 格式
	if matches := strings.Split(url, ".runsync.novita.dev"); len(matches) > 0 {
		name := matches[0]
		// 移除 https:// 或 http:// 前缀
		name = strings.TrimPrefix(name, "https://")
		name = strings.TrimPrefix(name, "http://")
		return name
	}
	return ""
}

// ListSAPs 获取 SAP 列表，支持按 SE 名称过滤
func ListSAPs(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// 从查询参数获取 endpoint URL
	endpointURL := r.URL.Query().Get("se")
	// 从 URL 提取 SE 名称
	seName := extractSEFromURL(endpointURL)

	log.Printf("Listing SAPs in namespace: %s with SE filter: %s (from URL: %s)",
		namespace, seName, endpointURL)

	saps, err := k8sService.ListServerlessPolicies(namespace, seName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, saps, "")
}

// GetSAP 获取指定的 ServerlessAutoscalingPolicy 资源
func GetSAP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Getting SAP %s in namespace: %s", name, namespace)
	policy, err := k8sService.GetServerlessPolicy(namespace, name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	policy.ObjectMeta.ManagedFields = nil

	utils.SuccessResponse(w, policy, "")
}

// SetSAPValue 设置 SAP 的值
func SetSAPValue(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	endpointName := r.URL.Query().Get("endpoint")
	// get sap by endpoint name
	saps, err := k8sService.ListServerlessPolicies(namespace, endpointName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var targetSap *v1beta1.ServerlessAutoScalingPolicy
	for _, sap := range saps {
		if sap.Spec.ScaleTargetRef.Name == endpointName {
			targetSap = sap
			break
		}
	}

	if targetSap == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "SAP not found for endpoint: "+endpointName)
		return
	}

	var sapParam models.SapParam
	err = json.NewDecoder(r.Body).Decode(&sapParam)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate sapParam
	if sapParam.MinReplicas > sapParam.MaxReplicas {
		utils.ErrorResponse(w, http.StatusBadRequest, "minReplicas must be less than maxReplicas")
		return
	}

	if sapParam.MinReplicas < 0 || sapParam.MaxReplicas < 0 || sapParam.ConcurrencyPerWorker < 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "minReplicas, maxReplicas and concurrencyPerWorker must larger than 0")
		return
	}
	if sapParam.ScaleUpWindow < 0 || sapParam.ScaleDownWindow < 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "scaleUpWindow and scaleDownWindow must larger than 0")
		return
	}

	if targetSap.Spec.Metrics == nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "metrics is required")
		return
	}
	metric := targetSap.Spec.Metrics[0]
	if metric.Resource.Target.Type != "AverageValue" {
		utils.ErrorResponse(w, http.StatusBadRequest, "metrics type must be AverageValue")
		return
	}

	targetSap.Spec.MinReplicas = &sapParam.MinReplicas
	targetSap.Spec.MaxReplicas = sapParam.MaxReplicas
	targetSap.Spec.Metrics[0].Resource.Target.AverageValueAsInt = &sapParam.ConcurrencyPerWorker
	targetSap.Spec.Behavior.ScaleUp.StabilizationWindowSeconds = &sapParam.ScaleUpWindow
	targetSap.Spec.Behavior.ScaleDown.StabilizationWindowSeconds = &sapParam.ScaleDownWindow

	log.Printf("Updating SAP: %s", targetSap.Name)
	// 记录审计日志
	detail := fmt.Sprintf("MinReplicas: %d, MaxReplicas: %d, ConcurrencyPerWorker: %d, ScaleUpWindow: %d, ScaleDownWindow: %d",
		sapParam.MinReplicas, sapParam.MaxReplicas, sapParam.ConcurrencyPerWorker, sapParam.ScaleUpWindow, sapParam.ScaleDownWindow)
	requestBody, _ := json.Marshal(sapParam)

	err = k8sService.UpdateServerlessPolicy(namespace, targetSap)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		CreateAuditLog(claims.Username, "更新SAP配置", r, http.StatusInternalServerError, "更新失败: "+err.Error(), string(requestBody), endpointName, detail)
		return
	}

	CreateAuditLog(claims.Username, "更新SAP配置", r, http.StatusOK, "更新成功", string(requestBody), endpointName, detail)

	utils.SuccessResponse(w, targetSap, "")
}

// ListWorkers 获取 Worker 列表，支持按 SE 名称过滤
func ListWorkers(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// 从查询参数获取 endpoint URL
	endpointURL := r.URL.Query().Get("se")
	// 从 URL 提取 SE 名称
	seName := extractSEFromURL(endpointURL)
	if len(seName) == 0 {
		seName = r.URL.Query().Get("endpoint")
	}

	// log.Printf("Listing Workers in namespace: %s with SE filter: %s (from URL: %s)",
	// 	namespace, seName, endpointURL)

	workers, err := k8sService.ListWorkers(namespace, seName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, worker := range workers {
		worker.ObjectMeta.ManagedFields = nil
	}

	utils.SuccessResponse(w, workers, "")
}

// GetWorker 获取指定的 Worker 资源
func GetWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Getting worker %s in namespace: %s", name, namespace)
	worker, err := k8sService.GetWorker(namespace, name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	worker.ObjectMeta.ManagedFields = nil
	utils.SuccessResponse(w, worker, "")
}

// DeleteWorker 删除指定的 Worker 资源
func DeleteWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	log.Printf("Deleting worker %s in namespace: %s", name, namespace)
	err := k8sService.DeleteWorker(namespace, name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 记录审计日志
	CreateAuditLog(claims.Username, "删除Worker", r, http.StatusOK, "删除成功", "", name, "")
	utils.SuccessResponse(w, nil, "")
}
