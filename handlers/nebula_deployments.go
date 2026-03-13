package handlers

import (
	"llm-ops/utils"
	"log"
	"net/http"
)

func ListNDeployments(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling UpdateEndpoint Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	_, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// Get namespace from query parameter, default to "default"
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "nebula"
	}

	ndeployments, err := k8sServicev2.ListNDeployments(namespace)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, ndeployments, "")
}
