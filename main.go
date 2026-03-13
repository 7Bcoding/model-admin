package main

import (
	"fmt"
	"llm-ops/config"
	"llm-ops/db"
	"llm-ops/handlers"
	"llm-ops/middleware"
	"llm-ops/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	userService    *services.UserService
	trackerService *services.TrackerService
	auditService   *services.AuditService
	store          = sessions.NewCookieStore([]byte("your-secret-key"))
)

func getConfigPath() string {
	// 优先使用环境变量指定的配置文件
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		return configFile
	}
	// 默认使用 /app/config/config.yaml
	return "/app/config/config.yaml"
}

func init() {
	// 加载配置文件
	configPath := getConfigPath()
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 移除 users.json 文件的依赖，因为我们使用数据库存储用户信息
	userService = services.NewUserService(db.GetDB())
	handlers.InitUserService(userService)

	// 初始化 tracker service
	trackerService = services.NewTrackerService(
		config.Config.Tracker.BaseURL,
	)
	handlers.InitTrackerService(trackerService)

	// 初始化审计日志服务
	auditService = services.NewAuditService(db.GetDB())
	handlers.InitAuditService(auditService)

	// 自动迁移审计日志表
	// if err := db.GetDB().AutoMigrate(&models.AuditLog{}); err != nil {
	// 	log.Fatalf("Failed to migrate audit_logs table: %v", err)
	// }
}

func main() {
	// 创建路由器
	r := mux.NewRouter()

	// 应用 CORS 中间件
	r.Use(middleware.CORS)

	// add stop channel
	stopCh := make(chan struct{})
	defer close(stopCh)

	// 初始化 K8s 服务
	k8sService, err := services.NewKubernetesService(config.Config.Kubernetes.ConfigPath, stopCh)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes service: %v", err)
	}
	handlers.InitK8sService(k8sService)

	k8sServiceV2, err := services.NewKubernetesServiceV2(config.Config.Kubernetes.ConfigPathV2, stopCh)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes service v2: %v", err)
	}
	handlers.InitK8sServiceV2(k8sServiceV2)

	log.Printf("Kubernetes service v2 initialized")

	// 初始化模型服务
	novitaModelService := services.NewModelService()
	ppioModelService := services.NewPPIOModelService()
	handlers.InitModelService(novitaModelService, ppioModelService)

	// API 路由
	api := r.PathPrefix("/api/v1").Subrouter()

	// 公开接口 - 不需要认证
	api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/logout", handlers.Logout).Methods("POST", "OPTIONS")
	api.HandleFunc("/feishu/callback", handlers.HandleFeishuLogin).Methods("POST", "OPTIONS")
	// 需要认证的接口
	// 用户管理相关 API
	api.HandleFunc("/users", middleware.AuthRequired(middleware.AdminRequired(handlers.ListUsers))).Methods("GET", "OPTIONS")
	api.HandleFunc("/users", middleware.AuthRequired(middleware.AdminRequired(handlers.AddUser))).Methods("POST", "OPTIONS")
	api.HandleFunc("/users/{username}/role", middleware.AuthRequired(middleware.AdminRequired(handlers.UpdateUserRole))).Methods("PUT", "OPTIONS")
	api.HandleFunc("/users/{username}", middleware.AuthRequired(middleware.AdminRequired(handlers.DeleteUser))).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/users/change-password", middleware.AuthRequired(handlers.ChangePassword)).Methods("POST", "OPTIONS")

	// 审计日志相关 API
	api.HandleFunc("/audit-logs", middleware.AuthRequired(middleware.AdminRequired(handlers.ListAuditLogs))).Methods("GET", "OPTIONS")

	// 缓存管理相关 API
	api.HandleFunc("/cache/stats", middleware.AuthRequired(middleware.AdminRequired(handlers.GetCacheStats))).Methods("GET", "OPTIONS")
	api.HandleFunc("/cache/clear/{cacheName}", middleware.AuthRequired(middleware.AdminRequired(handlers.ClearCache))).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/cache/clear-all", middleware.AuthRequired(middleware.AdminRequired(handlers.ClearAllCaches))).Methods("DELETE", "OPTIONS")

	// 模型相关 API
	api.HandleFunc("/models", middleware.AuthRequired(handlers.ListModels)).Methods("GET", "OPTIONS")
	api.HandleFunc("/models-inspection", middleware.AuthRequired(handlers.ListModelsInspection)).Methods("GET", "OPTIONS")
	api.HandleFunc("/models/star", middleware.AuthRequired(handlers.StarModel)).Methods("POST", "OPTIONS")
	api.HandleFunc("/models/unstar", middleware.AuthRequired(handlers.UnstarModel)).Methods("POST", "OPTIONS")
	api.HandleFunc("/models/note", middleware.AuthRequired(handlers.UpdateModelNote)).Methods("POST", "OPTIONS")

	// 集群相关 API
	api.HandleFunc("/nexusclusters", middleware.AuthRequired(handlers.ListNexusClusters)).Methods("GET", "OPTIONS")
	api.HandleFunc("/nexusclusters/nodes", middleware.AuthRequired(handlers.ListNexusNodes)).Methods("GET", "OPTIONS")
	api.HandleFunc("/nexusclusters/nodes/detail", middleware.AuthRequired(handlers.GetNodeDetail)).Methods("GET", "OPTIONS")
	api.HandleFunc("/nexusclusters/instances/detail", middleware.AuthRequired(handlers.GetInstanceDetail)).Methods("GET", "OPTIONS")
	api.HandleFunc("/nexusclusters/nodes/schedulable", middleware.AuthRequired(handlers.SetNodeSchedulable)).Methods("POST", "OPTIONS")
	api.HandleFunc("/nexusclusters/nodes/label", middleware.AuthRequired(handlers.SetNodeLabel)).Methods("POST", "OPTIONS")

	// 模型相关路由
	api.HandleFunc("/models/{modelName}/deployment", middleware.AuthRequired(handlers.GetModelDeployment)).Methods("GET", "OPTIONS")

	// Kubernetes 资源相关路由
	api.HandleFunc("/se", middleware.AuthRequired(handlers.ListEndpoints)).Methods("GET", "OPTIONS")
	api.HandleFunc("/se/{name}", middleware.AuthRequired(handlers.GetEndpoint)).Methods("GET", "OPTIONS")
	api.HandleFunc("/se/{name}/image", middleware.AuthRequired(handlers.UpdateEndpointImage)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/se/{name}/maxConcurrency", middleware.AuthRequired(handlers.UpdateEndpointMaxConcurrency)).Methods("PUT", "OPTIONS")

	api.HandleFunc("/clusters", middleware.AuthRequired(handlers.ListClusters)).Methods("GET", "OPTIONS")
	api.HandleFunc("/clusters/{name}", middleware.AuthRequired(handlers.GetCluster)).Methods("GET", "OPTIONS")

	api.HandleFunc("/sap", middleware.AuthRequired(handlers.ListSAPs)).Methods("GET", "OPTIONS")
	api.HandleFunc("/sap/{name}", middleware.AuthRequired(handlers.GetSAP)).Methods("GET", "OPTIONS")
	// set sap value
	api.HandleFunc("/sap/{name}", middleware.AuthRequired(handlers.SetSAPValue)).Methods("POST", "OPTIONS")

	api.HandleFunc("/workers", middleware.AuthRequired(handlers.ListWorkers)).Methods("GET", "OPTIONS")
	api.HandleFunc("/workers/{name}", middleware.AuthRequired(handlers.GetWorker)).Methods("GET", "OPTIONS")
	api.HandleFunc("/workers/{name}", middleware.AuthRequired(handlers.DeleteWorker)).Methods("DELETE", "OPTIONS")

	// 添加 endpoint 更路由
	api.HandleFunc("/endpoints/{id}", middleware.AuthRequired(handlers.UpdateEndpoint)).Methods("PUT", "OPTIONS")
	// 添加 endpoint 停用路由
	api.HandleFunc("/endpoints/stop", middleware.AuthRequired(handlers.StopEndpoint)).Methods("POST", "OPTIONS")
	// 添加 endpoint 创建路由
	api.HandleFunc("/endpoints/create", middleware.AuthRequired(handlers.AddEndpoint)).Methods("POST", "OPTIONS")

	api.HandleFunc("/nebula/cnodes", middleware.AuthRequired(middleware.NebulaRequired(handlers.ListCNodes))).Methods("GET", "OPTIONS")
	api.HandleFunc("/nebula/cnodes/{name}", middleware.AuthRequired(middleware.NebulaRequired(handlers.GetCNode))).Methods("GET", "OPTIONS")
	api.HandleFunc("/nebula/ndeployments", middleware.AuthRequired(middleware.NebulaRequired(handlers.ListNDeployments))).Methods("GET", "OPTIONS")
	// api.HandleFunc("/nebula/ndeployments/{name}", middleware.AuthRequired(middleware.NebulaRequired(handlers.GetNDeployment))).Methods("GET", "OPTIONS")
	api.HandleFunc("/nebula/workers", middleware.AuthRequired(middleware.NebulaRequired(handlers.ListNebulaWorkers))).Methods("GET", "OPTIONS")
	// api.HandleFunc("/nebula/workers/{name}", middleware.AuthRequired(middleware.NebulaRequired(handlers.GetWorker))).Methods("GET", "OPTIONS")
	// api.HandleFunc("/nebula/workers/{name}", middleware.AuthRequired(middleware.NebulaRequired(handlers.DeleteWorker))).Methods("DELETE", "OPTIONS")
	// api.HandleFunc("/nebula/scalingpolicies", middleware.AuthRequired(middleware.NebulaRequired(handlers.ListScalingPolicies))).Methods("GET", "OPTIONS")
	// api.HandleFunc("/nebula/scalingpolicies/{name}", middleware.AuthRequired(middleware.NebulaRequired(handlers.GetScalingPolicy))).Methods("GET", "OPTIONS")
	// api.HandleFunc("/nebula/scalingpolicies/{name}", middleware.AuthRequired(middleware.NebulaRequired(handlers.DeleteScalingPolicy))).Methods("DELETE", "OPTIONS")

	// 添加 metrics 代理路由
	api.HandleFunc("/metrics", middleware.AuthRequired(handlers.GetEndpointMetrics)).Methods("GET", "OPTIONS")

	handlers.InitTrackerService(trackerService)

	// 添加 tracker 相关路由
	api.HandleFunc("/nodes/search", middleware.AuthRequired(handlers.SearchNodes)).Methods("POST", "OPTIONS")
	api.HandleFunc("/nodes/models/warmup", middleware.AuthRequired(handlers.WarmupModels)).Methods("POST", "OPTIONS")
	api.HandleFunc("/nodes/models/delete", middleware.AuthRequired(handlers.DeleteModels)).Methods("DELETE", "OPTIONS")

	// 添加断流worker相关路由
	api.HandleFunc("/exclude-workers", middleware.AuthRequired(handlers.ListExcludeWorkers)).Methods("GET", "OPTIONS")
	api.HandleFunc("/exclude-workers", middleware.AuthRequired(handlers.AddExcludeWorker)).Methods("POST", "OPTIONS")
	api.HandleFunc("/exclude-workers", middleware.AuthRequired(handlers.DeleteExcludeWorker)).Methods("DELETE", "OPTIONS")

	// Model API routes
	api.HandleFunc("/modelapi", middleware.AuthRequired(handlers.ApplyModelAPI)).Methods("POST", "OPTIONS")

	// Fusion 统一代理路由 - 需要认证
	api.PathPrefix("/fusion/").HandlerFunc(middleware.AuthRequired(handlers.FusionProxy))

	// 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
