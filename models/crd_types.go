package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerlessEndpoint CRD 结构体定义
type ServerlessEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServerlessEndpointSpec   `json:"spec"`
	Status            ServerlessEndpointStatus `json:"status,omitempty"`
}

type ServerlessEndpointSpec struct {
	ClusterIDs           []string        `json:"clusterIDs"`
	ConcurrencyPerWorker float64         `json:"concurrencyPerWorker"`
	FastBoot             bool            `json:"fastBoot"`
	Provider             string          `json:"provider"`
	Replicas             float64         `json:"replicas"`
	RevisionHistoryLimit *float64        `json:"revisionHistoryLimit,omitempty"`
	Template             PodTemplateSpec `json:"template"`
}

type ServerlessEndpointStatus struct {
	AvailableReplicas float64 `json:"availableReplicas"`
	ReadyReplicas     float64 `json:"readyReplicas"`
	Replicas          float64 `json:"replicas"`
	UpdatedReplicas   float64 `json:"updatedReplicas"`
}

// PodTemplateSpec 描述 Pod 的模板
type PodTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PodSpec `json:"spec"`
}

// PodSpec 描述 Pod 的规格
type PodSpec struct {
	Containers    []Container       `json:"containers"`
	Affinity      *corev1.Affinity  `json:"affinity,omitempty"`
	RestartPolicy string            `json:"restartPolicy,omitempty"`
	NodeSelector  map[string]string `json:"nodeSelector,omitempty"`
	ClusterIDs    []string          `json:"clusterIDs,omitempty"`
	Provider      string            `json:"provider,omitempty"`
	Status        string            `json:"status,omitempty"`
}

// Container 描述容器的配置
type Container struct {
	Name           string    `json:"name"`
	Image          string    `json:"image"`
	Args           []string  `json:"args,omitempty"`
	Env            []EnvVar  `json:"env,omitempty"`
	Resources      Resources `json:"resources"`
	Ports          []float64 `json:"ports,omitempty"`
	ReadinessProbe *Probe    `json:"readinessProbe,omitempty"`
}

// Resources 描述资源需求
type Resources struct {
	CPURequests          ResourceValue `json:"cpuRequests"`
	MemoryRequests       ResourceValue `json:"memoryRequests"`
	GPURequests          GPURequest    `json:"gpuRequests"`
	EphemeralStorageReqs ResourceValue `json:"ephemeralStorageRequests"`
}

// ResourceValue 描述资源值
type ResourceValue struct {
	Number string `json:"number"`
}

// GPURequest 描述 GPU 需求
type GPURequest struct {
	Number int      `json:"number"`
	Models []string `json:"models"`
}

// Probe 描述健康检查
type Probe struct {
	HTTPGet             *HTTPGetAction `json:"httpGet,omitempty"`
	InitialDelaySeconds float64        `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      float64        `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       float64        `json:"periodSeconds,omitempty"`
	SuccessThreshold    float64        `json:"successThreshold,omitempty"`
	FailureThreshold    float64        `json:"failureThreshold,omitempty"`
}

// HTTPGetAction 描述 HTTP GET 请求
type HTTPGetAction struct {
	Path   string  `json:"path"`
	Port   float64 `json:"port"`
	Scheme string  `json:"scheme,omitempty"`
}

// 其他相关类型定义...
