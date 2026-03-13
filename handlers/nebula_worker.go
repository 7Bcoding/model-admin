package handlers

import (
	"llm-ops/utils"
	"log"
	"net/http"
)

func ListNebulaWorkers(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling UpdateEndpoint Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	_, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	deployment := r.URL.Query().Get("deployment")

	workers, err := k8sServicev2.ListNebulaWorkers(deployment)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, workers, "获取Nebula Worker列表成功")

}
