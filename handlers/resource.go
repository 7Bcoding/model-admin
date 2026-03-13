package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"llm-ops/api/grpc"
	nexus_api_server_v1 "llm-ops/api/nexus/api/nexus-api-server/v1"
	"llm-ops/config"
	"llm-ops/services"
	"llm-ops/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

var (
	tracker *services.TrackerService
)

func ListNexusClusters(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, config.NexusClusters, "")
}

type DiskUsage struct {
	MountPoint string `json:"mountPoint"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}

type WrappedNode struct {
	*nexus_api_server_v1.Node
	Models []services.ModelInfo `json:"models"`
	Disks  []services.DiskInfo  `json:"disks"`
	// Models []string `json:"models"`
}

func ListNexusNodes(w http.ResponseWriter, r *http.Request) {
	clusterID := r.URL.Query().Get("cluster_id")
	if clusterID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "cluster_id is required")
		return
	}

	// 获取集群配置
	var cluster *config.NexusCluster
	for _, c := range config.NexusClusters {
		if c.ID == clusterID {
			cluster = c
			break
		}
	}
	if cluster == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "cluster not found")
		return
	}

	var client nexus_api_server_v1.NexusApiClient
	var err error

	if cluster.Client != nil {
		client = cluster.Client
	} else {
		conn := grpc.NewApiServerClient(&grpc.RemoteServerOptions{
			Address:  cluster.GrpcURL,
			Token:    cluster.Token,
			Insecure: false,
		})
		client, err = conn.Client()
		if err != nil {
			log.Println(fmt.Sprintf("failed to create client: %v, url: %s, token: %s", err, cluster.GrpcURL, cluster.Token))
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cluster.Client = client
	}

	log.Println(fmt.Sprintf("connected to %s", cluster.GrpcURL))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// 调用 ListNodes API
	resp, err := client.ListNodes(ctx, &nexus_api_server_v1.ListNodesRequest{})
	if err != nil {
		log.Println(fmt.Sprintf("failed to list nodes: %v, url: %s, token: %s", err, cluster.GrpcURL, cluster.Token))
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	nodeNames := make([]string, 0, len(resp.Nodes))
	nodeMap := make(map[string]*WrappedNode)
	respNode := make([]*WrappedNode, 0, len(resp.Nodes))
	for _, node := range resp.Nodes {
		nodeNames = append(nodeNames, node.Id)
		nodeMap[node.Id] = &WrappedNode{
			Node: node,
		}
		respNode = append(respNode, nodeMap[node.Id])
	}

	// split nodeNames into batches of 10
	batchSize := 100
	batchedNodeNames := make([][]string, 0)
	for i := 0; i < len(nodeNames); i += batchSize {
		end := i + batchSize
		if end > len(nodeNames) {
			end = len(nodeNames)
		}
		batchedNodeNames = append(batchedNodeNames, nodeNames[i:end])
	}

	// process each batch
	for _, batch := range batchedNodeNames {
		klog.Infof("processing batch: %v", batch)
		batchServers, err := trackerService.SearchModels(
			[]services.RegionFilter{
				{
					RegionID: "",
					Servers:  batch,
				},
			},
		)
		if err != nil {
			log.Printf("failed to search models for batch: %v, error: %v", batch, err)
			continue
		}
		// klog.Infof("batch servers: %v", batchServers)

		for _, region := range batchServers.Data.Regions {
			for _, node := range region.Servers {
				if nodeMap[node.ServerID] != nil {
					nodeMap[node.ServerID].Models = node.Models
					filterdDisks := make([]services.DiskInfo, 0, len(node.Disks))
					foundNetProxy := false
					for _, disk := range node.Disks {
						if disk.MountPoint == "/etc/nexus/proxy/net_proxy" {
							disk.DiskID = "net_proxy"
							filterdDisks = append(filterdDisks, disk)
							foundNetProxy = true
							break
						}
						// } else if disk.MountPoint == "/var/lib/docker" {
						// 	disk.DiskID = "docker"
						// 	filterdDisks = append(filterdDisks, disk)
						// }
					}
					if !foundNetProxy {
						for _, disk := range node.Disks {
							if disk.MountPoint == "/" {
								disk.DiskID = "net_proxy"
								filterdDisks = append(filterdDisks, disk)
								break
							}
						}
					}
					nodeMap[node.ServerID].Disks = filterdDisks
				}
			}
		}
	}
	// klog.Infof("nodeNames: %v", nodeNames)
	// servers, err := trackerService.SearchModels(
	// 	[]services.RegionFilter{
	// 		{
	// 			RegionID: "",
	// 			Servers:  nodeNames,
	// 		},
	// 	},
	// )
	// klog.Infof("servers: %v", servers)

	if err != nil {
		log.Printf("failed to search models: %v", err)
	} else {

	}

	utils.SuccessResponse(w, respNode, "")
}

func SetNodeSchedulable(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClusterID   string `json:"cluster_id"`
		NodeID      string `json:"node_id"`
		Schedulable bool   `json:"schedulable"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// 获取集群配置
	var cluster *config.NexusCluster
	for _, c := range config.NexusClusters {
		if c.ID == req.ClusterID {
			cluster = c
			break
		}
	}
	if cluster == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "cluster not found")
		return
	}

	if cluster.Client == nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "client not found")
		return
	}

	// 从上下文中获取用户信息
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	// 调用 GRPC 接口
	_, err := cluster.Client.SetNodeSchedulable(context.Background(), &nexus_api_server_v1.SetNodeSchedulableRequest{
		Id:          req.NodeID,
		Schedulable: req.Schedulable,
	})
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 记录审计日志
	action := "启用节点调度"
	if !req.Schedulable {
		action = "禁用节点调度"
	}
	CreateAuditLog(claims.Username, action, r, http.StatusOK, "操作成功", "", req.NodeID, "")

	utils.SuccessResponse(w, nil, "")
}

func SetNodeLabel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClusterID string            `json:"cluster_id"`
		NodeIds   []string          `json:"node_ids"`
		Labels    map[string]string `json:"labels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// 获取集群配置
	var cluster *config.NexusCluster
	for _, c := range config.NexusClusters {
		if c.ID == req.ClusterID {
			cluster = c
			break
		}
	}
	if cluster == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "cluster not found")
		return
	}

	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
		return
	}

	for _, nodeId := range req.NodeIds {
		_, err := cluster.Client.SetNodeLabel(context.Background(), &nexus_api_server_v1.SetNodeLabelRequest{
			Id:     nodeId,
			Labels: req.Labels,
		})
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// 记录审计日志
	detail := fmt.Sprintf("节点: %s, 标签: %v", req.NodeIds, req.Labels)
	CreateAuditLog(claims.Username, "设置节点标签", r, http.StatusOK, "操作成功", "", strings.Join(req.NodeIds, ","), detail)

	utils.SuccessResponse(w, nil, "")
}

func GetNodeDetail(w http.ResponseWriter, r *http.Request) {
	clusterID := r.URL.Query().Get("cluster_id")
	nodeID := r.URL.Query().Get("node_id")

	if clusterID == "" || nodeID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "cluster_id and node_id are required")
		return
	}

	// 获取集群配置
	var cluster *config.NexusCluster
	for _, c := range config.NexusClusters {
		if c.ID == clusterID {
			cluster = c
			break
		}
	}
	if cluster == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "cluster not found")
		return
	}

	if cluster.Client == nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "client not found")
		return
	}

	// 调用 GRPC 接口获取节点详情
	resp, err := cluster.Client.GetNode(context.Background(), &nexus_api_server_v1.GetNodeRequest{
		Id: nodeID,
	})
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, resp.Node, "")
}

func GetInstanceDetail(w http.ResponseWriter, r *http.Request) {
	clusterID := r.URL.Query().Get("cluster_id")
	clusterIndex := r.URL.Query().Get("cluster_index")
	instanceID := r.URL.Query().Get("instance_id")

	if (clusterID == "" && clusterIndex == "") || instanceID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "cluster_id/cluster_index and instance_id are required")
		return
	}

	// 获取集配置
	var cluster *config.NexusCluster
	if clusterID != "" {
		for _, c := range config.NexusClusters {
			if c.ID == clusterID {
				cluster = c
				break
			}
		}
	} else {
		for _, c := range config.NexusClusters {
			if strconv.Itoa(c.Index) == clusterIndex || strconv.Itoa(c.IndexAlias) == clusterIndex {
				cluster = c
				break
			}
		}
	}
	if cluster == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "cluster not found")
		return
	}

	if cluster.Client == nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "client not found")
		return
	}

	// 调用 GRPC 接口获取实例详情
	resp, err := cluster.Client.GetInstance(context.Background(), &nexus_api_server_v1.GetInstanceRequest{
		Id: instanceID,
	})
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, resp.Instance, "")
}
