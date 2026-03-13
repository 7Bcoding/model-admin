package services

import (
	"fmt"
	"llm-ops/models"
	"log"
	"strconv"
	"sync"
	"time"

	"k8s.io/client-go/tools/cache"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	v1beta1 "gitlab.paigod.work/ai-cloud/nebula/api/v1beta1"
	nebulaclientset "gitlab.paigod.work/ai-cloud/nebula/pkg/client/clientset/versioned"
	"gitlab.paigod.work/ai-cloud/nebula/pkg/client/informers/externalversions"
	nebulaLister "gitlab.paigod.work/ai-cloud/nebula/pkg/client/listers/nebula.novita.ai/v1beta1"
)

type KubernetesServiceV2 struct {
	clientset     *kubernetes.Clientset
	config        *rest.Config
	mu            sync.RWMutex
	ndeployLister nebulaLister.DeploymentLister
	nworkerLister nebulaLister.WorkerLister

	cnodeLock sync.RWMutex
	cnodes    map[string]*v1beta1.Node

	deploymentLock sync.RWMutex
	deployments    map[string]*v1beta1.Deployment

	workerLock sync.RWMutex
	workers    map[string]*v1beta1.Worker

	nspLock         sync.RWMutex
	nsps            map[string]*v1beta1.ScalingPolicy
	nspByDeployment map[string]*v1beta1.ScalingPolicy
}

func NewKubernetesServiceV2(kubeconfigPath string, stopCh <-chan struct{}) (*KubernetesServiceV2, error) {
	log.Printf("Initializing serverless v2 service with config: %s", kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client, err := nebulaclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := externalversions.NewSharedInformerFactory(client, 1*time.Hour)
	ndeployLister := factory.Nebula().V1beta1().Deployments().Lister()
	nworkerLister := factory.Nebula().V1beta1().Workers().Lister()
	// factory.Start(stopCh)

	k := &KubernetesServiceV2{
		clientset:       clientset,
		config:          config,
		ndeployLister:   ndeployLister,
		nworkerLister:   nworkerLister,
		cnodes:          make(map[string]*v1beta1.Node),
		deployments:     make(map[string]*v1beta1.Deployment),
		workers:         make(map[string]*v1beta1.Worker),
		nsps:            make(map[string]*v1beta1.ScalingPolicy),
		nspByDeployment: make(map[string]*v1beta1.ScalingPolicy),
	}

	cnodeInformer := factory.Nebula().V1beta1().Nodes().Informer()
	cnodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cnode := obj.(*v1beta1.Node)
			k.cnodeLock.Lock()
			k.cnodes[cnode.Name] = cnode
			k.cnodeLock.Unlock()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			cnode := newObj.(*v1beta1.Node)
			k.cnodeLock.Lock()
			k.cnodes[cnode.Name] = cnode
			k.cnodeLock.Unlock()
		},
		DeleteFunc: func(obj interface{}) {
			cnode, ok := obj.(*v1beta1.Node)
			if !ok {
				return
			}
			k.cnodeLock.Lock()
			delete(k.cnodes, cnode.Name)
			k.cnodeLock.Unlock()
		},
	})
	factory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, cnodeInformer.HasSynced) {
		return nil, fmt.Errorf("failed to sync cnode cache")
	}
	log.Printf("cnode cache synced")

	deploymentInformer := factory.Nebula().V1beta1().Deployments().Informer()
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment := obj.(*v1beta1.Deployment)
			k.deploymentLock.Lock()
			k.deployments[deployment.Name] = deployment
			k.deploymentLock.Unlock()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deployment := newObj.(*v1beta1.Deployment)
			k.deploymentLock.Lock()
			k.deployments[deployment.Name] = deployment
			k.deploymentLock.Unlock()
		},
		DeleteFunc: func(obj interface{}) {
			deployment, ok := obj.(*v1beta1.Deployment)
			if !ok {
				return
			}
			k.deploymentLock.Lock()
			delete(k.deployments, deployment.Name)
			k.deploymentLock.Unlock()
		},
	})
	factory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, deploymentInformer.HasSynced) {
		return nil, fmt.Errorf("failed to sync deployment cache")
	}
	log.Printf("deployment cache synced")

	workerInformer := factory.Nebula().V1beta1().Workers().Informer()
	workerInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			worker := obj.(*v1beta1.Worker)
			k.workerLock.Lock()
			k.workers[worker.Name] = worker
			k.workerLock.Unlock()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			worker := newObj.(*v1beta1.Worker)
			k.workerLock.Lock()
			k.workers[worker.Name] = worker
			k.workerLock.Unlock()
		},
		DeleteFunc: func(obj interface{}) {
			worker, ok := obj.(*v1beta1.Worker)
			if !ok {
				return
			}
			k.workerLock.Lock()
			delete(k.workers, worker.Name)
			k.workerLock.Unlock()
		},
	})
	if !cache.WaitForCacheSync(stopCh, workerInformer.HasSynced) {
		return nil, fmt.Errorf("failed to sync worker cache")
	}
	factory.Start(stopCh)
	log.Printf("worker cache synced")

	nspInformer := factory.Nebula().V1beta1().ScalingPolicies().Informer()
	nspInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nsp := obj.(*v1beta1.ScalingPolicy)
			k.nspLock.Lock()
			k.nsps[nsp.Name] = nsp
			if nsp.Spec.ScaleTarget.Name != "" {
				k.nspByDeployment[nsp.Spec.ScaleTarget.Name] = nsp
			}
			k.nspLock.Unlock()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			nsp := newObj.(*v1beta1.ScalingPolicy)
			k.nspLock.Lock()
			k.nsps[nsp.Name] = nsp
			if nsp.Spec.ScaleTarget.Name != "" {
				k.nspByDeployment[nsp.Spec.ScaleTarget.Name] = nsp
			}
			k.nspLock.Unlock()
		},
		DeleteFunc: func(obj interface{}) {
			nsp, ok := obj.(*v1beta1.ScalingPolicy)
			if !ok {
				return
			}
			k.nspLock.Lock()
			delete(k.nsps, nsp.Name)
			if nsp.Spec.ScaleTarget.Name != "" {
				delete(k.nspByDeployment, nsp.Spec.ScaleTarget.Name)
			}
			k.nspLock.Unlock()
		},
	})
	factory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, nspInformer.HasSynced) {
		return nil, fmt.Errorf("failed to sync scaling policy cache")
	}
	log.Printf("scaling policy cache synced")
	return k, nil
}

func (k *KubernetesServiceV2) ListCNodes() ([]*models.CNodeListItem, error) {
	result := make([]*models.CNodeListItem, 0, len(k.cnodes))
	k.cnodeLock.RLock()
	defer k.cnodeLock.RUnlock()
	hostName := ""
	for _, cnode := range k.cnodes {
		for _, address := range cnode.Status.Addresses {
			if address.Type == "Hostname" {
				hostName = address.Address
				break
			}
		}
		item := &models.CNodeListItem{
			TypeMeta: metav1.TypeMeta{
				APIVersion: cnode.APIVersion,
				Kind:       cnode.Kind,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cnode.Name,
				Namespace: cnode.Namespace,
			},
			Name:       cnode.Name,
			Provider:   cnode.Status.Provider,
			Region:     cnode.Status.Region,
			HostName:   hostName,
			CPUProduct: cnode.Status.CpuProduct,
			GPUProduct: cnode.Status.GpuProduct,
			Policy:     string(cnode.Spec.SchedulePolicy),
			Status:     cnode.Status.Status,
			Age:        cnode.CreationTimestamp.Time.Format(time.RFC3339),
		}

		imageCache := make([]*models.ImageCache, 0)
		modelCache := make([]*models.ModelCache, 0)
		for _, image := range cnode.Status.Images {
			imageCache = append(imageCache, &models.ImageCache{
				Name:   image.Name,
				Status: image.Status,
				Size:   int(image.Size.IntValue()),
			})
		}
		for _, model := range cnode.Status.Volumes {
			modelCache = append(modelCache, &models.ModelCache{
				Name:   model.Name,
				Status: model.Status,
				Size:   int(model.Size.IntValue()),
			})
		}
		item.ImageCache = imageCache
		item.ModelCache = modelCache

		// Calculate age from creation timestamp
		if !cnode.CreationTimestamp.IsZero() {
			age := time.Since(cnode.CreationTimestamp.Time)
			if age.Hours() >= 24 {
				item.Age = fmt.Sprintf("%dd", int(age.Hours()/24))
			} else if age.Hours() >= 1 {
				item.Age = fmt.Sprintf("%dh", int(age.Hours()))
			} else {
				item.Age = fmt.Sprintf("%dm", int(age.Minutes()))
			}
		}

		item.GPUState = make([]int32, len(cnode.Status.GpuState))
		for i, state := range cnode.Status.GpuState {
			item.GPUState[i] = int32(state)
		}

		result = append(result, item)
	}
	return result, nil
}

func (k *KubernetesServiceV2) ListNDeployments(namespace string) ([]*models.NDeploymentItem, error) {
	result := make([]*models.NDeploymentItem, 0, len(k.deployments))
	k.deploymentLock.RLock()
	defer k.deploymentLock.RUnlock()
	for _, deployment := range k.deployments {
		item := &models.NDeploymentItem{
			TypeMeta: metav1.TypeMeta{
				APIVersion: deployment.APIVersion,
				Kind:       deployment.Kind,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployment.Name,
				Namespace: deployment.Namespace,
			},
			Status:         string(deployment.Status.Status),
			Hot:            *deployment.Spec.HotReplicas,
			Cold:           *deployment.Spec.ColdReplicas,
			HotReady:       *deployment.Status.HotReadyReplicas,
			ColdReady:      *deployment.Status.ColdReadyReplicas,
			CustomerRegion: deployment.Spec.WorkerTemplate.Spec.CustomerRegion,
		}
		item.Age = deployment.CreationTimestamp.Time.Format(time.RFC3339)
		item.Resources = &models.Resources{
			CPURequests: models.ResourceValue{
				Number: strconv.Itoa(deployment.Spec.WorkerTemplate.Spec.Resources.CpuNum),
			},
			MemoryRequests: models.ResourceValue{
				Number: strconv.Itoa(deployment.Spec.WorkerTemplate.Spec.Resources.Memory.IntValue()),
			},
			GPURequests: models.GPURequest{
				Number: deployment.Spec.WorkerTemplate.Spec.Resources.GpuNum,
				Models: []string{deployment.Spec.WorkerTemplate.Spec.Resources.GpuProduct},
			},
			EphemeralStorageReqs: models.ResourceValue{
				Number: strconv.Itoa(deployment.Spec.WorkerTemplate.Spec.Resources.Storage.IntValue()),
			},
		}
		if deployment.Spec.MaxBatchSize != nil {
			item.MaxBatchSize = *deployment.Spec.MaxBatchSize
		}

		simpleNsp := &models.SimpleNSP{}
		nsp, ok := k.nspByDeployment[deployment.Name]
		if ok {
			simpleNsp.BatchPerWorker = *nsp.Spec.BatchPerWorker
			simpleNsp.MinReplicas = *nsp.Spec.MinReplicas
			simpleNsp.MaxReplicas = *nsp.Spec.MaxReplicas
		}

		item.SimpleNSP = simpleNsp

		for _, container := range deployment.Spec.WorkerTemplate.Spec.Containers {
			item.Image = container.Image
			break
		}

		result = append(result, item)
	}
	return result, nil
}

func (k *KubernetesServiceV2) GetNDeployment(name string, namespace string) (*models.NDeploymentItem, error) {
	deployment, err := k.ndeployLister.Deployments(namespace).Get(name)
	if err != nil {
		return nil, err
	}

	item := &models.NDeploymentItem{
		TypeMeta: metav1.TypeMeta{
			APIVersion: deployment.APIVersion,
			Kind:       deployment.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
		Status:         string(deployment.Status.Status),
		Hot:            *deployment.Spec.HotReplicas,
		HotReady:       *deployment.Status.HotReadyReplicas,
		Cold:           *deployment.Spec.ColdReplicas,
		ColdReady:      *deployment.Status.ColdReadyReplicas,
		CustomerRegion: deployment.Spec.WorkerTemplate.Spec.CustomerRegion,
		Resources: &models.Resources{
			CPURequests: models.ResourceValue{
				Number: strconv.Itoa(deployment.Spec.WorkerTemplate.Spec.Resources.CpuNum),
			},
			MemoryRequests: models.ResourceValue{
				Number: strconv.Itoa(deployment.Spec.WorkerTemplate.Spec.Resources.Memory.IntValue()),
			},
			GPURequests: models.GPURequest{
				Number: deployment.Spec.WorkerTemplate.Spec.Resources.GpuNum,
				Models: []string{deployment.Spec.WorkerTemplate.Spec.Resources.GpuProduct},
			},
			EphemeralStorageReqs: models.ResourceValue{
				Number: strconv.Itoa(deployment.Spec.WorkerTemplate.Spec.Resources.Storage.IntValue()),
			},
		},
		NodeSelector: deployment.Spec.WorkerTemplate.Spec.NodeSelector,
		Containers:   deployment.Spec.WorkerTemplate.Spec.Containers,
		Volumes:      deployment.Spec.WorkerTemplate.Spec.Volumes,
	}

	// Set Image field from first container
	for _, container := range deployment.Spec.WorkerTemplate.Spec.Containers {
		item.Image = container.Image
		break
	}

	return item, nil
}

func (k *KubernetesServiceV2) ListNebulaWorkers(deployment string) ([]*models.NebulaWorkerItem, error) {
	workers, err := k.nworkerLister.Workers(deployment).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var workerItems []*models.NebulaWorkerItem
	for _, worker := range workers {
		// Calculate age
		age := ""
		if !worker.CreationTimestamp.IsZero() {
			age = time.Since(worker.CreationTimestamp.Time).Round(time.Second).String()
		}

		if len(worker.OwnerReferences) > 0 {
			owner := worker.OwnerReferences[0]
			if owner.Kind == "Deployment" {
				deployment = owner.Name
			}
		}

		workerItem := &models.NebulaWorkerItem{
			TypeMeta: metav1.TypeMeta{
				APIVersion: worker.APIVersion,
				Kind:       worker.Kind,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      worker.Name,
				Namespace: worker.Namespace,
			},
			Phase:      string(worker.Spec.Phase),
			GpuProduct: worker.Spec.Resources.GpuProduct,
			Gpu:        worker.Spec.Resources.GpuIds,
			Node:       worker.Spec.NodeName,
			Status:     string(worker.Status.Status),
			Age:        age,
			Deployment: deployment,
		}
		workerItems = append(workerItems, workerItem)
	}

	return workerItems, nil
}
