package handlers

import (
	"encoding/json"
	"llm-ops/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func GetModelEndpoints(w http.ResponseWriter, r *http.Request) {
	modelName := mux.Vars(r)["modelName"]

	// 获取用户信息 - 仅做认证检查
	_, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// 获取模型的 endpoints
	endpoints, err := novitaModelService.GetModelEndpoints(modelName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, endpoints, "")
}

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EndpointURL string `json:"endpoint_url"`
		ModelName   string `json:"model_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := novitaModelService.TestEndpoint(req.EndpointURL, req.ModelName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, result, "")
}

// ... 其他端点处理函数
