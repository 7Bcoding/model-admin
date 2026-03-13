package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type Cluster struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Cluster.
	// +optional
	Spec ClusterSpec `json:"spec"`

	// Most recently observed status of the Cluster.
	// +optional
	Status ClusterStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

type ClusterSpec struct {
	// +optional
	Provider string `json:"provider"`
	// +optional
	ClusterID string `json:"clusterID"`
	// Unschedulable controls cluster schedulability of new workers. By default, worker is schedulable.
	// +optional
	Unschedulable bool `json:"unschedulable,omitempty"`
}

type ClusterStatus struct {
	// Allocatable represents the resources of a cluster that are available for scheduling.
	// Defaults to Capacity.
	// +optional
	Allocatable []ClusterNodeResource `json:"allocatable,omitempty"`

	// Conditions is an array of current observed cluster conditions.
	// Refer to: https://kubernetes.io/docs/concepts/nodes/node/#condition
	// +optional
	Conditions []ClusterCondition `json:"conditions,omitempty"`
	// List of addresses reachable to the cluster.
	// +optional
	Addresses []ClusterAddress `json:"addresses,omitempty"`
	// ServerlessEndpoints running on the Cluster.
	// +optional
	ServerlessEndpoints []NamespacedName `json:"serverlessEndpoints,omitempty"`
	// Set of infos to uniquely identify the cluster.
	// +optional
	ClusterInfo ClusterSystemInfo `json:"clusterInfo,omitempty"`
}

type ClusterNodeResource struct {
	NodeName         string            `json:"nodeName"`
	GPUModel         string            `json:"gpuModel"`
	GPUNumber        uint32            `json:"gpuNumber"`
	CPU              resource.Quantity `json:"cpu"`
	Memory           resource.Quantity `json:"memory"`
	EphemeralStorage resource.Quantity `json:"ephemeralStorage"`
	// cpu per gpu
	CpuPer int64 `json:"cpuPer,omitempty"`

	Schedulable bool `json:"schedulable"`
}

// ClusterSystemInfo is a set of infos to uniquely identify the cluster.
type ClusterSystemInfo struct {
	// +optional
	Version string `json:"version"`
	// +optional
	Region string `json:"region"`
}

// ClusterAddress contains information for the node's address.
type ClusterAddress struct {
	// Cluster address type, one of Hostname or IP.
	Type ClusterAddressType `json:"type"`
	// The cluster address.
	Address string `json:"address"`
}

type ClusterAddressType string

const (
	ClusterHostName ClusterAddressType = "Hostname"
	ClusterIP       ClusterAddressType = "IP"
)

type ClusterConditionType string

// These are valid but not exhaustive conditions of cluster.
// Relevant events contain "ClusterReady", "ClusterNotReady", "ClusterSchedulable", and "ClusterNotSchedulable".
const (
	// ClusterReady means cluster is healthy and ready to accept workers.
	ClusterReady ClusterConditionType = "Ready"
	// ClusterMemoryPressure means the cluster is under pressure due to insufficient available memory.
	ClusterMemoryPressure ClusterConditionType = "MemoryPressure"
	// ClusterDiskPressure means the cluster is under pressure due to insufficient available disk.
	ClusterDiskPressure ClusterConditionType = "DiskPressure"
	// ClusterVolumePressure means the cluster is under pressure due to insufficient network volume.
	ClusterVolumePressure ClusterConditionType = "VolumePressure"
	// ClusterNetworkUnavailable means that network for the cluster is not correctly configured.
	ClusterNetworkUnavailable ClusterConditionType = "NetworkUnavailable"
)

// ClusterCondition contains condition information for a cluster.
type ClusterCondition struct {
	// Type of cluster condition.
	Type ClusterConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status"`
	// Last time we got an update on a given condition.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty"`
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

type NamespacedName struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}
