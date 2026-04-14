package handlers

import (
	"bytes"
	"io"
	"llm-ops/config"
	"llm-ops/utils"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// extractResourceFromPath 从URL路径中提取资源名称
func extractResourceFromPath(path string) (string, string) {
	// 移除 /api/v1/fusion 前缀
	cleanPath := strings.TrimPrefix(path, "/api/v1/fusion")
	if cleanPath == "" || cleanPath == "/" {
		return "", ""
	}

	// 使用正则表达式匹配 RESTful API 路径
	// 匹配 /accounts, /models, /models/{id}/providers 等
	accountsPattern := regexp.MustCompile(`^/accounts(/.*)?$`)
	modelsPattern := regexp.MustCompile(`^/models(/.*)?$`)
	providersPattern := regexp.MustCompile(`^/models/[^/]+/providers(/.*)?$`)

	if accountsPattern.MatchString(cleanPath) {
		return "accounts", "账户管理"
	} else if providersPattern.MatchString(cleanPath) {
		return "providers", "模型提供商管理"
	} else if modelsPattern.MatchString(cleanPath) {
		return "models", "模型管理"
	}

	return "", ""
}

// getActionFromMethod 根据HTTP方法获取操作描述
func getActionFromMethod(method string, resource string) string {
	switch strings.ToUpper(method) {
	case "POST":
		return "Fusion-创建" + resource
	case "PUT":
		return "Fusion-更新" + resource
	default:
		return "Fusion-操作" + resource
	}
}

// FusionProxy 统一的代理处理器，根据 X-Fusion-Provider header 来选择路由目标
func FusionProxy(w http.ResponseWriter, r *http.Request) {
	// 检查是否需要记录审计日志（仅针对写操作）
	shouldAudit := r.Method == "POST" || r.Method == "PUT"
	var auditData struct {
		username     string
		resource     string
		resourceDesc string
		action       string
		requestBody  string
	}

	// 如果需要审计，提取相关信息
	if shouldAudit {
		// 获取用户信息
		if claims, ok := utils.GetUserFromContext(r.Context()); ok {
			auditData.username = claims.Username
		}

		// 提取资源信息
		auditData.resource, auditData.resourceDesc = extractResourceFromPath(r.URL.Path)
		auditData.action = getActionFromMethod(r.Method, auditData.resourceDesc)

		// 读取请求体用于审计日志
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				auditData.requestBody = string(bodyBytes)
				// 重新设置请求体，以便后续使用
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}
	}

	// 从 header 中获取服务提供商
	provider := r.Header.Get("X-Fusion-Provider")
	if provider == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Missing X-Fusion-Provider header")
		return
	}

	var targetURL *url.URL
	var token string
	var err error

	// 根据 provider 选择配置
	switch strings.ToLower(provider) {
	case "beta":
		if config.Config.BetaFusion.URL == "" {
			utils.ErrorResponse(w, http.StatusInternalServerError, "beta Fusion URL not configured")
			return
		}
		targetURL, err = url.Parse(config.Config.BetaFusion.URL)
		token = config.Config.BetaFusion.Token
	case "alpha":
		if config.Config.AlphaFusion.URL == "" {
			utils.ErrorResponse(w, http.StatusInternalServerError, "alpha Fusion URL not configured")
			return
		}
		targetURL, err = url.Parse(config.Config.AlphaFusion.URL)
		token = config.Config.AlphaFusion.Token
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid X-Fusion-Provider: "+provider)
		return
	}

	if err != nil {
		log.Printf("Error parsing Fusion URL for provider %s: %v", provider, err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Invalid Fusion URL")
		return
	}

	// URL 重写：/api/v1/fusion/* -> /admin/v1/*
	originalPath := r.URL.Path
	pathWithoutAPIV1 := strings.TrimPrefix(originalPath, "/api/v1")
	newPath := strings.Replace(pathWithoutAPIV1, "/fusion", "/admin/v1", 1)

	targetURL.Path = newPath
	targetURL.RawQuery = r.URL.RawQuery

	log.Printf("Proxying Fusion request to %s: %s %s -> %s", provider, r.Method, originalPath, targetURL.String())

	// 创建新的请求
	var body io.Reader
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			utils.ErrorResponse(w, http.StatusBadRequest, "Error reading request body")
			return
		}
		body = bytes.NewReader(bodyBytes)
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), body)
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error creating proxy request")
		return
	}

	// 复制原始请求的头部
	for name, values := range r.Header {
		// 跳过一些不应该转发的头部和我们的自定义头部
		if name == "Host" || name == "Connection" || name == "Upgrade" || name == "X-Fusion-Provider" {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// 添加外部服务的API认证头部
	// 注意：这里会覆盖用户的JWT token，因为外部服务需要自己的API token
	if token != "" {
		proxyReq.Header.Set("Authorization", "Bearer "+token)
	}

	// 设置正确的 Host 头部
	proxyReq.Host = targetURL.Host

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error sending proxy request to %s: %v", provider, err)
		utils.ErrorResponse(w, http.StatusBadGateway, "Error sending proxy request")
		return
	}
	defer resp.Body.Close()

	// 复制响应头部
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}

	// 记录审计日志（仅针对写操作且有资源信息的请求）
	if shouldAudit && auditData.resource != "" && auditData.username != "" {
		var result string
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result = "操作成功"
		} else {
			result = "操作失败"
		}

		// 构建目标描述，包含提供商和资源信息
		target := provider + ":" + auditData.resource
		detail := auditData.resourceDesc

		CreateAuditLog(
			auditData.username,
			auditData.action,
			r,
			resp.StatusCode,
			result,
			auditData.requestBody,
			target,
			detail,
		)
	}
}
