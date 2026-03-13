package services

import (
	"llm-ops/models"
	"time"

	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// CreateAuditLog 创建审计日志
func (s *AuditService) CreateAuditLog(req *models.AuditLogRequest) error {
	log := &models.AuditLog{
		Operator:    req.Operator,
		RequestURL:  req.RequestURL,
		Method:      req.Method,
		Action:      req.Action,
		Target:      req.Target,
		Detail:      req.Detail,
		Result:      req.Result,
		Status:      req.Status,
		RequestBody: req.RequestBody,
		CreatedAt:   time.Now(),
	}
	return s.db.Create(log).Error
}

// ListAuditLogs 查询审计日志列表
func (s *AuditService) ListAuditLogs(page, pageSize int, operator, action, target string, startTime, endTime time.Time) ([]models.AuditLog, int64, error) {
	query := s.db.Model(&models.AuditLog{})

	// 添加筛选条件，使用 LIKE 进行模糊搜索
	if operator != "" {
		query = query.Where("operator LIKE ?", "%"+operator+"%")
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if target != "" {
		query = query.Where("target LIKE ?", "%"+target+"%")
	}
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.AuditLog
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}
