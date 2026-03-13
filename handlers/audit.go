package handlers

import (
	"llm-ops/models"
	"llm-ops/services"
	"llm-ops/utils"
	"net/http"
	"strconv"
	"time"
)

var auditService *services.AuditService

func InitAuditService(service *services.AuditService) {
	auditService = service
}

// ListAuditLogs 获取审计日志列表
func ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	// 获取分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 获取筛选参数
	operator := r.URL.Query().Get("operator")
	action := r.URL.Query().Get("action")
	target := r.URL.Query().Get("target")

	// 处理时间范围
	var startTime, endTime time.Time
	if startTimeStr := r.URL.Query().Get("startTime"); startTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
			startTime = t
		}
	}
	if endTimeStr := r.URL.Query().Get("endTime"); endTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
			endTime = t
		}
	}

	// 查询审计日志
	logs, total, err := auditService.ListAuditLogs(page, pageSize, operator, action, target, startTime, endTime)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回结果
	utils.SuccessResponse(w, map[string]interface{}{
		"total": total,
		"items": logs,
		"page": map[string]interface{}{
			"current": page,
			"size":    pageSize,
		},
	}, "")
}

// CreateAuditLog 创建审计日志的辅助函数
func CreateAuditLog(operator, action string, r *http.Request, status int, result string, requestBody string, target string, detail string) {
	if auditService == nil {
		return
	}

	// 创建审计日志
	req := &models.AuditLogRequest{
		Operator:    operator,
		RequestURL:  r.URL.String(),
		Method:      r.Method,
		Action:      action,
		Target:      target,
		Detail:      detail,
		Result:      result,
		Status:      status,
		RequestBody: requestBody,
	}

	if err := auditService.CreateAuditLog(req); err != nil {
		// 记录错误但不影响主流程
		utils.LogError("Failed to create audit log", err)
	}
}
