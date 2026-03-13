package handlers

import (
	"llm-ops/utils"
	"log"
	"net/http"
)

func ListCNodes(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling UpdateEndpoint Request ===")
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)

	_, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	cnodes, err := k8sServicev2.ListCNodes()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, cnodes, "")
}

func GetCNode(w http.ResponseWriter, r *http.Request) {

}
