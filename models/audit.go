package models

import "time"

// AuditLog 审计日志模型
type AuditLog struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Operator    string    `json:"operator" gorm:"type:varchar(255);index:idx_operator,length:255;not null;comment:'操作者'"`
	RequestURL  string    `json:"request_url" gorm:"type:text;not null;comment:'请求URL'"`
	Method      string    `json:"method" gorm:"type:varchar(10);not null;comment:'HTTP方法'"`
	Action      string    `json:"action" gorm:"type:varchar(255);index:idx_action,length:255;not null;comment:'操作类型'"`
	Target      string    `json:"target" gorm:"type:varchar(255);index:idx_target,length:255;comment:'操作对象'"`
	Detail      string    `json:"detail" gorm:"type:text;comment:'操作详情'"`
	Result      string    `json:"result" gorm:"type:text;not null;comment:'操作结果'"`
	Status      int       `json:"status" gorm:"not null;comment:'HTTP状态码'"`
	RequestBody string    `json:"request_body,omitempty" gorm:"type:text;comment:'请求体'"`
	CreatedAt   time.Time `json:"created_at" gorm:"index;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
}

// AuditLogRequest 创建审计日志的请求结构
type AuditLogRequest struct {
	Operator    string `json:"operator"`
	RequestURL  string `json:"request_url"`
	Method      string `json:"method"`
	Action      string `json:"action"`
	Target      string `json:"target"`
	Detail      string `json:"detail"`
	Result      string `json:"result"`
	Status      int    `json:"status"`
	RequestBody string `json:"request_body,omitempty"`
}
