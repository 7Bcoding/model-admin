package services

import (
	"context"
	"fmt"
	"llm-ops/generated/clientset/versioned"
	"llm-ops/generated/listers/serverless/v1beta1"
	"llm-ops/models"
	serverlessv1beta1 "llm-ops/serverless/v1beta1"
	"log"
	"sync"
	"time"

	nitorclientset "llm-ops/generated/clientset/versioned"
	informers "llm-ops/generated/informers/externalversions"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesService struct {
	clientset      *kubernetes.Clientset
	config         *rest.Config
	mu             sync.RWMutex
	workerLister   v1beta1.WorkerLister
	endpointLister v1beta1.ServerlessEndpointLister
	sapLister      v1beta1.ServerlessAutoScalingPolicyLister

	endpoints map[string]*serverlessv1beta1.ServerlessEndpoint
	workers   map[string]*serverlessv1beta1.Worker
	saps      map[string]*serverlessv1beta1.ServerlessAutoScalingPolicy
}

func NewKubernetesService(kubeconfigPath string, stopCh <-chan struct{}) (*KubernetesService, error) {
	log.Printf("Initializing Kubernetes service with config: %s", kubeconfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating clientset: %v", err)
		return nil, err
	}

	client, err := nitorclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(client, time.Minute)
	workerLister := factory.Nitor().V1beta1().Workers().Lister()
	endpointLister := factory.Nitor().V1beta1().ServerlessEndpoints().Lister()
	sapLister := factory.Nitor().V1beta1().ServerlessAutoScalingPolicies().Lister()

	// 启动 informers
	factory.Start(stopCh)

	ret := &KubernetesService{
		clientset:      clientset,
		config:         config,
		workerLister:   workerLister,
		endpointLister: endpointLister,
		sapLister:      sapLister,
		endpoints:      make(map[string]*serverlessv1beta1.ServerlessEndpoint),
		workers:        make(map[string]*serverlessv1beta1.Worker),
		saps:           make(map[string]*serverlessv1beta1.ServerlessAutoScalingPolicy),
	}

	serverlessEndpointInformer := factory.Nitor().V1beta1().ServerlessEndpoints().Informer()

	serverlessEndpointInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			endpoint := obj.(*serverlessv1beta1.ServerlessEndpoint)
			ret.mu.Lock()
			defer ret.mu.Unlock()
			ret.endpoints[endpoint.Name] = endpoint
			// log.Printf("Added endpoint: %s", endpoint.Name)
		},
		DeleteFunc: func(obj interface{}) {
			endpoint, ok := obj.(*serverlessv1beta1.ServerlessEndpoint)
			if !ok {
				log.Printf("Error deleting endpoint: %v", obj)
				return
			}
			ret.mu.Lock()
			defer ret.mu.Unlock()
			delete(ret.endpoints, endpoint.Name)
			// log.Printf("Deleted endpoint: %s", endpoint.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			endpoint := newObj.(*serverlessv1beta1.ServerlessEndpoint)
			ret.mu.Lock()
			defer ret.mu.Unlock()
			ret.endpoints[endpoint.Name] = endpoint
			// log.Printf("Updated endpoint: %s", endpoint.Name)
		},
	})

	workerInformers := factory.Nitor().V1beta1().Workers().Informer()
	workerInformers.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			worker := obj.(*serverlessv1beta1.Worker)
			ret.mu.Lock()
			defer ret.mu.Unlock()
			ret.workers[worker.Name] = worker
			// log.Printf("Added worker: %s", worker.Name)
		},
		DeleteFunc: func(obj interface{}) {
			worker, ok := obj.(*serverlessv1beta1.Worker)
			if !ok {
				log.Printf("Error deleting worker: %v", obj)
				return
			}
			ret.mu.Lock()
			defer ret.mu.Unlock()
			delete(ret.workers, worker.Name)
			// log.Printf("Deleted worker: %s", worker.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			worker := newObj.(*serverlessv1beta1.Worker)
			ret.mu.Lock()
			defer ret.mu.Unlock()
			ret.workers[worker.Name] = worker
			// log.Printf("Updated worker: %s", worker.Name)
		},
	})

	sspInformer := factory.Nitor().V1beta1().ServerlessAutoScalingPolicies().Informer()
	sspInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nsp := obj.(*serverlessv1beta1.ServerlessAutoScalingPolicy)
			ret.mu.Lock()
			defer ret.mu.Unlock()
			ret.saps[nsp.Name] = nsp
			// log.Printf("Added nsp: %s", nsp.Name)
		},
		DeleteFunc: func(obj interface{}) {
			nsp, ok := obj.(*serverlessv1beta1.ServerlessAutoScalingPolicy)
			if !ok {
				log.Printf("Error deleting nsp: %v", obj)
				return
			}
			ret.mu.Lock()
			defer ret.mu.Unlock()
			delete(ret.saps, nsp.Name)
			// log.Printf("Deleted nsp: %s", nsp.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			nsp := newObj.(*serverlessv1beta1.ServerlessAutoScalingPolicy)
			ret.mu.Lock()
			defer ret.mu.Unlock()
			ret.saps[nsp.Name] = nsp
			// log.Printf("Updated nsp: %s", nsp.Name)
		},
	})

	// 等待缓存同步
	if !cache.WaitForCacheSync(stopCh, serverlessEndpointInformer.HasSynced, workerInformers.HasSynced, sspInformer.HasSynced) {
		log.Fatal("Failed to sync cache")
	}
	log.Printf("================ Endpoint and Worker and SAP cache synced")

	return ret, nil
}

// ListEndpoints 使用生成的客户端获取 Endpoints
func (k *KubernetesService) ListEndpoints(namespace string) ([]models.EndpointInfo, error) {
	log.Printf("Listing endpoints in namespace: %s", namespace)

	k.mu.RLock()
	endpoints := k.endpoints
	k.mu.RUnlock()

	// use informer to list all endpoints
	// endpoints, err := k.endpointLister.List(labels.Everything())
	// if err != nil {
	// 	log.Printf("Error listing endpoints: %v", err)
	// 	return nil, err
	// }

	result := make([]models.EndpointInfo, 0, len(endpoints))
	for _, endpoint := range endpoints {
		endpoint.ObjectMeta.ManagedFields = nil
		info := models.EndpointInfo{
			APIVersion: endpoint.APIVersion,
			Kind:       endpoint.Kind,
			Metadata: models.EndpointMetadata{
				Name:        endpoint.Name,
				Namespace:   endpoint.Namespace,
				Labels:      endpoint.Labels,
				Annotations: endpoint.Annotations,
			},
			Spec: models.EndpointSpec{
				ClusterIDs:           endpoint.Spec.ClusterIDs,
				ConcurrencyPerWorker: float64(endpoint.Spec.ConcurrencyPerWorker),
				FastBoot:             endpoint.Spec.FastBoot,
				Provider:             endpoint.Spec.Provider,
				Replicas:             float64(*endpoint.Spec.Replicas),
				Template: models.EndpointTemplate{
					Spec: convertWorkerTemplateToEndpointTemplate(endpoint.Spec.Template),
				},
			},
			MaxConcurrency: int(endpoint.Spec.ConcurrencyPerWorker),
			Image:          endpoint.Spec.Template.Spec.Containers[0].Image,
		}
		// yaml, err := yaml.Marshal(endpoint)
		// if err != nil {
		// 	log.Printf("Error marshaling endpoint to YAML: %v", err)
		// 	continue
		// }
		// // info.Yaml = string(yaml)
		// info.SapYaml = sapMap[endpoint.Name]
		result = append(result, info)
	}

	log.Printf("Found %d endpoints", len(result))
	return result, nil
}

// convertWorkerTemplateToEndpointTemplate 将 WorkerTemplateSpec 转换为 EndpointTemplateSpec
func convertWorkerTemplateToEndpointTemplate(template serverlessv1beta1.WorkerTemplateSpec) models.EndpointTemplateSpec {
	spec := models.EndpointTemplateSpec{
		Containers: make([]models.Container, len(template.Spec.Containers)),
	}

	// 转换容器信息
	for i, c := range template.Spec.Containers {
		container := models.Container{
			Name:  c.Name,
			Image: c.Image,
			Args:  c.Args,
			Env:   make([]models.EnvVar, len(c.Env)),
			Resources: models.Resources{
				CPURequests: models.ResourceValue{
					Number: fmt.Sprintf("%v", c.Resources.CPURequests),
				},
				MemoryRequests: models.ResourceValue{
					Number: fmt.Sprintf("%v", c.Resources.MemoryRequests),
				},
				GPURequests: models.GPURequest{
					Number: int(c.Resources.GPURequests.Number),
					Models: c.Resources.GPURequests.Models,
				},
				EphemeralStorageReqs: models.ResourceValue{
					Number: fmt.Sprintf("%v", c.Resources.EphemeralStorageRequests),
				},
			},
		}

		// 转换环境变量
		for j, env := range c.Env {
			container.Env[j] = models.EnvVar{
				Name:  env.Name,
				Value: env.Value,
			}
		}

		// 转换探针
		if c.ReadinessProbe != nil {
			container.ReadinessProbe = convertProbeToModels(c.ReadinessProbe)
		}

		spec.Containers[i] = container
	}

	// 转换亲和性设置
	if template.Spec.Affinity != nil {
		spec.Affinity = convertAffinityToModels(template.Spec.Affinity)
	}

	return spec
}

// 辅助转换函数
func convertEnvVarsToModels(envs []serverlessv1beta1.EnvVar) []models.EnvVar {
	result := make([]models.EnvVar, len(envs))
	for i, env := range envs {
		result[i] = models.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		}
	}
	return result
}

func convertProbeToModels(probe *serverlessv1beta1.Probe) *models.Probe {
	if probe == nil {
		return nil
	}
	return &models.Probe{
		HTTPGet: &models.HTTPGetAction{
			Path:   probe.HTTPGet.Path,
			Port:   float64(probe.HTTPGet.Port),
			Scheme: string(probe.HTTPGet.Scheme),
		},
		InitialDelaySeconds: float64(probe.InitialDelaySeconds),
		TimeoutSeconds:      float64(probe.TimeoutSeconds),
		PeriodSeconds:       float64(probe.PeriodSeconds),
		SuccessThreshold:    float64(probe.SuccessThreshold),
		FailureThreshold:    float64(probe.FailureThreshold),
	}
}

func convertAffinityToModels(affinity *serverlessv1beta1.Affinity) models.Affinity {
	if affinity == nil || affinity.NodeAffinity == nil {
		return models.Affinity{}
	}

	return models.Affinity{
		NodeAffinity: convertNodeAffinityToModels(affinity.NodeAffinity),
	}
}

func convertNodeAffinityToModels(nodeAffinity *serverlessv1beta1.NodeAffinity) models.NodeAffinity {
	if nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		return models.NodeAffinity{}
	}

	return models.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: convertNodeSelectorToModels(nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution),
	}
}

func convertNodeSelectorToModels(selector *serverlessv1beta1.NodeSelector) models.NodeSelector {
	terms := make([]models.NodeSelectorTerm, len(selector.NodeSelectorTerms))
	for i, term := range selector.NodeSelectorTerms {
		terms[i] = convertNodeSelectorTermToModels(term)
	}
	return models.NodeSelector{
		NodeSelectorTerms: terms,
	}
}

func convertNodeSelectorTermToModels(term serverlessv1beta1.NodeSelectorTerm) models.NodeSelectorTerm {
	exprs := make([]models.MatchExpression, len(term.MatchExpressions))
	for i, expr := range term.MatchExpressions {
		exprs[i] = models.MatchExpression{
			Key:      expr.Key,
			Operator: string(expr.Operator),
			Values:   expr.Values,
		}
	}
	return models.NodeSelectorTerm{
		MatchExpressions: exprs,
	}
}

// GetEndpoint 获取单个 Endpoint
func (k *KubernetesService) GetEndpoint(namespace, name string) (*models.EndpointInfo, error) {
	log.Printf("Getting endpoint %s in namespace %s", name, namespace)

	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint client: %v", err)
	}

	k.mu.Lock()
	endpoint, ok := k.endpoints[name]
	if !ok {
		k.mu.Unlock()
		log.Printf("Endpoint %s not found", name)
		return nil, fmt.Errorf("endpoint not found")
	}
	k.mu.Unlock()

	endpoint, err = client.NitorV1beta1().ServerlessEndpoints(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting endpoint: %v", err)
		return nil, err
	}
	endpoint.ObjectMeta.ManagedFields = nil
	yaml, err := yaml.Marshal(endpoint)
	if err != nil {
		log.Printf("Error marshaling endpoint to YAML: %v", err)
		return nil, err
	}

	info := &models.EndpointInfo{
		APIVersion: endpoint.APIVersion,
		Kind:       endpoint.Kind,
		Metadata: models.EndpointMetadata{
			Name:        endpoint.Name,
			Namespace:   endpoint.Namespace,
			Labels:      endpoint.Labels,
			Annotations: endpoint.Annotations,
		},
		Spec: models.EndpointSpec{
			ClusterIDs:           endpoint.Spec.ClusterIDs,
			ConcurrencyPerWorker: float64(endpoint.Spec.ConcurrencyPerWorker),
			FastBoot:             endpoint.Spec.FastBoot,
			Provider:             endpoint.Spec.Provider,
			Replicas:             float64(*endpoint.Spec.Replicas),
			Template: models.EndpointTemplate{
				Spec: convertWorkerTemplateToEndpointTemplate(endpoint.Spec.Template),
			},
		},
	}
	info.Yaml = string(yaml)

	return info, nil
}

// ListClusters 获取集群列表
func (k *KubernetesService) ListClusters(namespace string) ([]*serverlessv1beta1.Cluster, error) {
	log.Printf("Listing clusters in namespace: %s", namespace)

	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster client: %v", err)
	}

	clusters, err := client.NitorV1beta1().Clusters().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing clusters: %v", err)
		return nil, err
	}

	// 转换为指针切片
	result := make([]*serverlessv1beta1.Cluster, len(clusters.Items))
	for i := range clusters.Items {
		result[i] = &clusters.Items[i]
	}

	return result, nil
}

// GetCluster 获取单个集群
func (k *KubernetesService) GetCluster(namespace, name string) (*serverlessv1beta1.Cluster, error) {
	log.Printf("Getting cluster %s in namespace %s", name, namespace)

	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster client: %v", err)
	}

	return client.NitorV1beta1().Clusters().Get(context.Background(), name, metav1.GetOptions{})
}

// ListServerlessPolicies 获取 SAP 列表，支持按 SE 名称过滤
func (k *KubernetesService) ListServerlessPolicies(namespace string, seName string) ([]*serverlessv1beta1.ServerlessAutoScalingPolicy, error) {
	log.Printf("Fetching SAPs in namespace: %s with SE filter: %s", namespace, seName)

	// client, err := versioned.NewForConfig(k.config)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create SAP client: %v", err)
	// }

	k.mu.Lock()
	defer k.mu.Unlock()

	saps := k.saps

	// saps, err := client.NitorV1beta1().ServerlessAutoScalingPolicies(namespace).List(context.Background(), metav1.ListOptions{})
	// if err != nil {
	// 	log.Printf("Error listing SAPs: %v", err)
	// 	return nil, err
	// }

	var err error
	var filteredSAPs []*serverlessv1beta1.ServerlessAutoScalingPolicy
	for _, sap := range saps {
		sap.ObjectMeta.ManagedFields = nil
		// 如果指定了 SE 名称，进行过滤
		if seName == "" || sap.Spec.ScaleTargetRef.Name == seName {
			filteredSAPs = append(filteredSAPs, sap)
			if err != nil {
				log.Printf("Error converting SAP to map: %v", err)
				continue
			}
			filteredSAPs = append(filteredSAPs, sap)
		}
	}

	return filteredSAPs, nil
}

// GetServerlessPolicy 获取单个 ServerlessAutoscalingPolicy
func (k *KubernetesService) GetServerlessPolicy(namespace, name string) (*serverlessv1beta1.ServerlessAutoScalingPolicy, error) {
	log.Printf("Getting SAP %s in namespace %s", name, namespace)

	k.mu.Lock()

	sap, ok := k.saps[name]
	if !ok {
		k.mu.Unlock()
		client, err := versioned.NewForConfig(k.config)
		if err != nil {
			return nil, fmt.Errorf("failed to create SAP client: %v", err)
		}

		return client.NitorV1beta1().ServerlessAutoScalingPolicies(namespace).Get(context.Background(), name, metav1.GetOptions{})
	}
	k.mu.Unlock()
	return sap, nil
}

func (k *KubernetesService) UpdateEndpointImage(namespace string, seName string, image string) error {
	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return fmt.Errorf("failed to create SAP client: %v", err)
	}

	endpoint, err := client.NitorV1beta1().ServerlessEndpoints(namespace).Get(context.Background(), seName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting endpoint: %v", err)
		return err
	}
	endpoint.Spec.Template.Spec.Containers[0].Image = image

	return k.updateEndpoint(namespace, endpoint)
}

// UpdateEndpointMaxConcurrency 更新 Endpoint 最大并发数
func (k *KubernetesService) UpdateEndpointMaxConcurrency(namespace string, seName string, maxConcurrency int) error {
	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return fmt.Errorf("failed to create SAP client: %v", err)
	}

	endpoint, err := client.NitorV1beta1().ServerlessEndpoints(namespace).Get(context.Background(), seName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting endpoint: %v", err)
		return err
	}
	endpoint.Spec.ConcurrencyPerWorker = int(maxConcurrency)

	return k.updateEndpoint(namespace, endpoint)
}

// UpdateEndpoint 更新 Endpoint
func (k *KubernetesService) updateEndpoint(namespace string, se *serverlessv1beta1.ServerlessEndpoint) error {
	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return fmt.Errorf("failed to create SAP client: %v", err)
	}

	_, err = client.NitorV1beta1().ServerlessEndpoints(namespace).Update(context.Background(), se, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error updating SAP: %v", err)
		return err
	}

	return nil
}

// UpdateServerlessPolicy 更新 ServerlessAutoScalingPolicy
func (k *KubernetesService) UpdateServerlessPolicy(namespace string, sap *serverlessv1beta1.ServerlessAutoScalingPolicy) error {
	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return fmt.Errorf("failed to create SAP client: %v", err)
	}

	_, err = client.NitorV1beta1().ServerlessAutoScalingPolicies(namespace).Update(context.Background(), sap, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error updating SAP: %v", err)
		return err
	}

	return nil
}

// ListWorkers 获取 Worker 列表，支持按 SE 名称过滤
func (k *KubernetesService) ListWorkers(namespace string, seName string) ([]*serverlessv1beta1.Worker, error) {
	// log.Printf("Fetching Workers in namespace: %s with SE filter: %s", namespace, seName)

	var err error
	k.mu.Lock()
	defer k.mu.Unlock()
	workers := k.workers
	// if len(workers) == 0 {
	// 	// use informer to list all workers
	// 	workers, err = k.workerLister.List(labels.Everything())
	// 	// workers, err := client.NitorV1beta1().Workers(namespace).List(context.Background(), metav1.ListOptions{})
	// 	if err != nil {
	// 		log.Printf("Error listing workers: %v", err)
	// 		return nil, err
	// 	}
	// }

	var filteredWorkers []*serverlessv1beta1.Worker
	for _, worker := range workers {
		// 如果指定了 SE 名称，检查 ownerReferences
		if seName == "" {
			filteredWorkers = append(filteredWorkers, worker)
			if err != nil {
				log.Printf("Error converting Worker to map: %v", err)
				continue
			}
			filteredWorkers = append(filteredWorkers, worker)
		} else {
			for _, ownerRef := range worker.ObjectMeta.OwnerReferences {
				if ownerRef.Name == seName {
					filteredWorkers = append(filteredWorkers, worker)
					break
				}
			}
		}
	}

	return filteredWorkers, nil
}

// GetWorker 获取单个 Worker
func (k *KubernetesService) GetWorker(namespace, name string) (*serverlessv1beta1.Worker, error) {
	log.Printf("Getting worker %s in namespace %s", name, namespace)

	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create worker client: %v", err)
	}

	return client.NitorV1beta1().Workers(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

// ... 其他方法保持不变 ...
// DeleteWorker 删除指定的 Worker 资源
func (k *KubernetesService) DeleteWorker(namespace, name string) error {
	client, err := versioned.NewForConfig(k.config)
	if err != nil {
		return fmt.Errorf("failed to create worker client: %v", err)
	}
	return client.NitorV1beta1().Workers(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
