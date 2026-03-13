package handlers

import (
	"io/ioutil"
	"llm-ops/utils"
	"net/http"
)

const excludeWorkerBaseURL = "https://metric.alpha.dev/exclude-worker"

// ListExcludeWorkers 获取断流worker列表
func ListExcludeWorkers(w http.ResponseWriter, r *http.Request) {
	// 创建请求
	req, err := http.NewRequest("GET", excludeWorkerBaseURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get exclude workers", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// AddExcludeWorker 添加断流worker
func AddExcludeWorker(w http.ResponseWriter, r *http.Request) {
	// 获取worker参数
	worker := r.URL.Query().Get("worker")
	if worker == "" {
		http.Error(w, "Worker parameter is required", http.StatusBadRequest)
		return
	}

	// 构造请求URL
	url := excludeWorkerBaseURL + "?worker=" + worker

	// 创建POST请求
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to add exclude worker", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 记录审计日志
	if resp.StatusCode == http.StatusOK {
		CreateAuditLog(claims.Username, "添加断流Worker", r, http.StatusOK, "添加成功", "", worker, "")
	} else {
		CreateAuditLog(claims.Username, "添加断流Worker", r, resp.StatusCode, "添加失败", "", worker, "")
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)
}

// DeleteExcludeWorker 删除断流worker
func DeleteExcludeWorker(w http.ResponseWriter, r *http.Request) {
	// 获取worker参数
	worker := r.URL.Query().Get("worker")
	if worker == "" {
		http.Error(w, "Worker parameter is required", http.StatusBadRequest)
		return
	}

	// 构造请求URL
	url := excludeWorkerBaseURL + "?worker=" + worker

	// 创建DELETE请求
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to delete exclude worker", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 记录审计日志
	if resp.StatusCode == http.StatusOK {
		CreateAuditLog(claims.Username, "删除断流Worker", r, http.StatusOK, "删除成功", "", worker, "")
	} else {
		CreateAuditLog(claims.Username, "删除断流Worker", r, resp.StatusCode, "删除失败", "", worker, "")
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)
}
