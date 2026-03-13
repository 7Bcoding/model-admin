package handlers

import (
	"fmt"
	"io"
	"llm-ops/db"
	"llm-ops/models"
	"llm-ops/services"
	"log"
	"net/http"
)

// ApplyModelAPI handles the application of YAML configuration for model API
func ApplyModelAPI(w http.ResponseWriter, r *http.Request) {
	// Read YAML content
	yamlContent, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Apply YAML configuration
	modelAPIService := services.NewModelAPIService(db.GetDB())
	if err := modelAPIService.ApplyYAML(yamlContent); err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply YAML: %v", err), http.StatusBadRequest)
		return
	}

	// Create audit log
	operator := r.Context().Value("user").(string)
	auditService := services.NewAuditService(db.GetDB())
	auditLog := &models.AuditLogRequest{
		Operator:    operator,
		RequestURL:  r.URL.String(),
		Method:      r.Method,
		Action:      "应用ModelAPI配置",
		Target:      "ModelAPI",
		Detail:      "应用YAML配置文件",
		RequestBody: string(yamlContent),
		Status:      http.StatusOK,
		Result:      "成功",
	}
	if err := auditService.CreateAuditLog(auditLog); err != nil {
		log.Printf("Failed to create audit log: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("YAML configuration applied successfully"))
}
