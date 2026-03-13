package handlers

import (
	"context"
	"encoding/json"
	"llm-ops/config"
	"llm-ops/services"
	"net/http"

	"github.com/gorilla/mux"
)

type NodeResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func ListNodes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["cluster_name"]

	// 获取集群配置
	clusterConfig, ok := config.NexusClusters[clusterName]
	if !ok {
		json.NewEncoder(w).Encode(NodeResponse{
			Success: false,
			Message: "Cluster not found",
		})
		return
	}

	// 创建集群服务
	service, err := services.NewClusterService(clusterConfig.GrpcURL, clusterConfig.Token)
	if err != nil {
		json.NewEncoder(w).Encode(NodeResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	defer service.Close()

	// 获取节点列表
	nodes, err := service.ListNodes(context.Background())
	if err != nil {
		json.NewEncoder(w).Encode(NodeResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 返回结果
	json.NewEncoder(w).Encode(NodeResponse{
		Success: true,
		Data:    nodes,
	})
}
