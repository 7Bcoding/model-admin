package handlers

// ListSAPs 获取 SAP 列表，支持按 SE 名称过滤
// func ListSAPs(w http.ResponseWriter, r *http.Request) {
// 	seName := r.URL.Query().Get("se")
// 	log.Printf("Listing SAPs with SE filter: %s", seName)

// 	saps, err := k8sService.ListSAPs(seName)
// 	if err != nil {
// 		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	utils.SuccessResponse(w, saps, "")
// }

// ListWorkers 获取 Worker 列表，支持按 SE 名称过滤
// func ListWorkers(w http.ResponseWriter, r *http.Request) {
// 	seName := r.URL.Query().Get("se")
// 	log.Printf("Listing Workers with SE filter: %s", seName)

// 	workers, err := k8sService.ListWorkers(seName)
// 	if err != nil {
// 		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	utils.SuccessResponse(w, workers, "")
// }
