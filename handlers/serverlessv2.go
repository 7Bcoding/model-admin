package handlers

import (
	"llm-ops/services"
	"llm-ops/utils"
	"net/http"
)

var k8sServicev2 *services.KubernetesServiceV2

func InitK8sServiceV2(ks *services.KubernetesServiceV2) {
	k8sServicev2 = ks
}

func ListNDeploys(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	ndeploys, err := k8sServicev2.ListNDeployments(namespace)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, ndeploys, "success")
}

func GetNDeploy(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	ndeploy, err := k8sServicev2.GetNDeployment(name, namespace)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, ndeploy, "success")
}
