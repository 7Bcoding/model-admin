package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=se
type ServerlessEndpoint struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Endpoint.
	// +optional
	Spec EndpointSpec `json:"spec"`

	// Most recently observed status of the Endpoint.
	// +optional
	Status EndpointStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServerlessEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerlessEndpoint `json:"items"`
}

// +k8s:deepcopy-gen=true
type EndpointSpec struct {
	// +kubebuilder:default:=0
	// +optional
	Replicas *int32 `json:"replicas"`
	// +kubebuilder:default:=false
	// +optional
	FastBoot bool `json:"fastBoot"`
	// +optional
	Provider string `json:"provider"`
	// +optional
	ClusterIDs []string `json:"clusterIDs"`

	// +kubebuilder:validation:Minimum:=1
	ConcurrencyPerWorker int `json:"concurrencyPerWorker"`

	// Label selector for workers. Existing workers are selected by this
	// will be the ones affected by this Endpoint.
	// It must match the pod template's labels.
	Selector *metav1.LabelSelector `json:"selector"`

	// Template describes the workers that will be created.
	// The only allowed template.spec.restartPolicy value is "Always".
	Template WorkerTemplateSpec `json:"template"`

	// The endpoint strategy to use to replace existing workers with new ones.
	// +optional
	Strategy EndpointStrategy `json:"strategy,omitempty"`

	// The number of old history to retain to allow rollback.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Defaults to 10.
	// +kubebuilder:default:=10
	// +optional
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`

	// +optional
	Initialization *InitializationStrategy `json:"initialization,omitempty"`
}

// PodTemplateSpec describes the data a pod should have when created from a template
type WorkerTemplateSpec struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the pod.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec WorkerSpec `json:"spec,omitempty"`
}

// EndpointStrategy describes how to replace existing workers with new ones.
type EndpointStrategy struct {
	// Type of endpoint. Can be "OnDelete" or "RollingUpdate". Default is RollingUpdate.
	// +kubebuilder:default:="RollingUpdate"
	// +optional
	Type EndpointStrategyType `json:"type,omitempty"`

	// Rolling update config params. Present only if EndpointStrategyType =
	// RollingUpdate.
	//---
	// TODO: Update this to follow our convention for oneOf, whatever we decide it
	// to be.
	// +kubebuilder:default:={}
	// +optional
	RollingUpdate *RollingUpdateEndpoint `json:"rollingUpdate,omitempty"`
}

// +enum
type EndpointStrategyType string

const (
	// Replace the old workers by new one using rolling update i.e gradually scale down the old workers and scale up the new one.
	RollingUpdateEndpointStrategyType EndpointStrategyType = "RollingUpdate"

	// OnDeleteEndpointStrategyType triggers the legacy behavior. Version
	// tracking and ordered rolling restarts are disabled. Workers are recreated
	// from the EndpointSpec when they are manually deleted. When a scale
	// operation is performed with this strategy, specification version indicated
	// by the Endpoint's currentRevision.
	OnDeleteEndpointStrategyType EndpointStrategyType = "OnDelete"
)

// Spec to control the desired behavior of rolling update.
type RollingUpdateEndpoint struct {
	// The maximum number of workers that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired workers (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	// This can not be 0 if MaxSurge is 0.
	// Defaults to 0.
	// +kubebuilder:default:="25%"
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,1,opt,name=maxUnavailable"`

	// The maximum number of workers that can be scheduled above the desired number of
	// workers.
	// Value can be an absolute number (ex: 5) or a percentage of desired workers (ex: 10%).
	// This can not be 0 if MaxUnavailable is 0.
	// Absolute number is calculated from percentage by rounding up.
	// Defaults to 25%.
	// +kubebuilder:default:="25%"
	// +optional
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty" protobuf:"bytes,2,opt,name=maxSurge"`
}

type InitializationStrategy struct {
	Type InitializationStrategyType `json:"type,omitempty"`
	// +optional
	Creation *CreationInitializationStrategy `json:"creation,omitempty"`
}
type InitializationStrategyType string

const (
	PullImageInitializationStrategyType InitializationStrategyType = "pullImage"
)

type CreationInitializationStrategy struct {
	Count          int  `json:"count,omitempty"`
	AcrossClusters bool `json:"acrossClusters,omitempty"`
	AcrossNodes    bool `json:"acrossNodes,omitempty"`
}

type EndpointStatus struct {
	State EndpointState `json:"state"`

	// Total number of non-terminated workers targeted by this endpoint (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	MQTopicName     string `json:"mqTopicName,omitempty"`
	MQConsumerGroup string `json:"mqConsumerGroup,omitempty"`

	// Count of hash collisions for the Endpoint. The Endpoint controller
	// uses this field as a collision avoidance mechanism when it needs to
	// create the name for the newest ControllerRevision.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty"`
}

type EndpointState string

const (
	EndpointStateInitializing EndpointState = "Initializing"
	EndpointStateStopped      EndpointState = "Stopped"
	EndpointStateServing      EndpointState = "Serving"
	EndpointStateFailed       EndpointState = "Failed"
	EndpointStateTerminating  EndpointState = "Terminating"
)

type EndpointEnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EndpointVolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	Size      uint64 `json:"size"`
}

type EndpointService struct {
	Port             uint16                          `json:"port"`
	MaxConcurreny    uint16                          `json:"maxConcurreny"`
	HealthCheckProbe EndpointServiceHealthCheckProbe `json:"healthCheckProbe"`
}

type EndpointServiceHealthCheckProbe struct {
	Method              string `json:"method"`
	Path                string `json:"path"`
	SuccessStatusCode   int    `json:"successStatusCode"`
	InitialDelaySeconds uint16 `json:"initialDelaySeconds"`
	PeriodSeconds       uint16 `json:"periodSeconds"`
	TimeoutSeconds      uint16 `json:"timeoutSeconds"`
	SuccessThreshold    uint16 `json:"successThreshold"`
}

// +k8s:deepcopy-gen=true
type EndpointScalePolicy struct {
	MinReplicas uint16                `json:"minReplicas"`
	MaxReplicas uint16                `json:"maxReplicas"`
	Strategy    EndpointScaleStrategy `json:"strategy"`
}

// +k8s:deepcopy-gen=true
type EndpointScaleStrategy struct {
	QueueDepth   *EndpointScaleByQueueDepth   `json:"queueDepth,omitempty"`
	RequestDelay *EndpointScaleByRequestDelay `json:"requestDelay,omitempty"`
}

type EndpointScaleByQueueDepth struct {
	QueueDepthThreshold uint16 `json:"queueDepthThreshold"`
	// ScaleUpRatio        uint16 `json:"scaleUpRatio"`
	// ScaleDownRatio      uint16 `json:"scaleDownRatio"`
}

type EndpointScaleByRequestDelay struct {
	RequestDelayThreshold uint16 `json:"requestDelayThreshold"`
	// ScaleUpRatio          uint16 `json:"scaleUpRatio"`
	// ScaleDownRatio        uint16 `json:"scaleDownRatio"`
}

type EndpointResources struct {
	CPU    EndpointResourceCPU    `json:"cpu"`
	Memory EndpointResourceMemory `json:"memory"`
	GPU    EndpointResourceGPU    `json:"gpu"`
}

type EndpointResourceCPU struct {
	Number uint16 `json:"number"`
}

type EndpointResourceMemory struct {
	Size uint64 `json:"size"`
}

type EndpointResourceGPU struct {
	Number      uint16 `json:"number"`
	Model       string `json:"model"`
	CUDAVersion string `json:"cudaVersion"`
}
