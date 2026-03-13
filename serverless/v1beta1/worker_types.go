package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=".spec.provider",description="Provider of the worker"
// +kubebuilder:printcolumn:name="ClusterID",type=string,JSONPath=".spec.clusterIDs",description="ClusterIDs of the worker"
// +kubebuilder:printcolumn:name="Stage",type=string,JSONPath=".status.state",description="status of the worker"
// +kubebuilder:printcolumn:name="Healthy",type=string,JSONPath=".status.healthy",description="healthy of the worker"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
type Worker struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Worker.
	// +optional
	Spec WorkerSpec `json:"spec"`

	// Most recently observed status of the Worker.
	// +optional
	Status WorkerStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WorkerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Worker `json:"items"`
}

// +k8s:deepcopy-gen=true
type WorkerSpec struct {
	// +optional
	Provider string `json:"provider"`
	// +optional
	ClusterIDs []string `json:"clusterIDs"`
	// +optional
	Status WorkerState `json:"status"`

	// List of containers belonging to the worker.
	// Containers cannot currently be added or removed.
	// There must be at least one container in a worker.
	// Cannot be updated.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Containers []Container `json:"containers"`

	// Restart policy for all containers within the worker.
	// One of Always, OnFailure, Never. In some contexts, only a subset of those values may be permitted.
	// Default to Always.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
	// +optional
	//RestartPolicy RestartPolicy `json:"restartPolicy,omitempty"`

	// If specified, the worker's scheduling constraints
	// +optional
	Affinity *Affinity `json:"affinity,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets"`

	// List of volumes that can be mounted by containers belonging to the worker.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	// +listType=map
	// +listMapKey=name
	Volumes []Volume `json:"volumes,omitempty"`
}

// Volume represents a named volume in a worker that may be accessed by any container in the worker.
type Volume struct {
	// name of the volume.
	// Must be a DNS_LABEL and unique within the worker.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// volumeSource represents the location and type of the mounted volume.
	// If not specified, the Volume is implied to be an EmptyDir.
	// This implied behavior is deprecated and will be removed in a future version.
	VolumeSource `json:",inline" protobuf:"bytes,2,opt,name=volumeSource"`
}

// Represents the source of a volume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	NitorVolume *NitorVolume `json:"nitorVolume"`
}

type NitorVolume struct {
	Size            uint64 `json:"size"`
	IsNetworkVolume *bool  `json:"isNetworkVolume"`
	// +optional
	NetworkVolumeID string `json:"networkVolumeID"`
}

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
// +structType=atomic
type LocalObjectReference struct {
	// Name of the referent.
	// This field is effectively required, but due to backwards compatibility is
	// allowed to be empty. Instances of this type with an empty value here are
	// almost certainly wrong.
	// TODO: Add other useful fields. apiVersion, kind, uid?
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	// +default=""
	// +kubebuilder:default=""
	// TODO: Drop `kubebuilder:default` when controller-gen doesn't need it https://github.com/kubernetes-sigs/kubebuilder/issues/3896.
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
}

// Affinity is a group of affinity scheduling rules.
type Affinity struct {
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *NodeAffinity `json:"nodeAffinity,omitempty"`
}

// Node affinity is a group of node affinity scheduling rules.
type NodeAffinity struct {
	// NOT YET IMPLEMENTED. TODO: Uncomment field once it is implemented.
	// If the affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to an update), the system
	// will try to eventually evict the pod from its node.
	// +optional
	// RequiredDuringSchedulingRequiredDuringExecution *NodeSelector `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to an update), the system
	// may or may not try to eventually evict the pod from its node.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution *NodeSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,opt,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule pods to nodes that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a node that violates one or more of the expressions. The node that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the node matches the corresponding matchExpressions; the
	// node(s) with the highest sum are the most preferred.
	// +optional
	// +listType=atomic
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// An empty preferred scheduling term matches all objects with implicit weight 0
// (i.e. it's a no-op). A null preferred scheduling term matches no objects (i.e. is also a no-op).
type PreferredSchedulingTerm struct {
	// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
	// A node selector term, associated with the corresponding weight.
	Preference NodeSelectorTerm `json:"preference" protobuf:"bytes,2,opt,name=preference"`
}

// A node selector represents the union of the results of one or more label queries
// over a set of nodes; that is, it represents the OR of the selectors represented
// by the node selector terms.
// +structType=atomic
type NodeSelector struct {
	// Required. A list of node selector terms. The terms are ORed.
	// +listType=atomic
	NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms" protobuf:"bytes,1,rep,name=nodeSelectorTerms"`
}

// A null or empty node selector term matches no objects. The requirements of
// them are ANDed.
// The TopologySelectorTerm type implements a subset of the NodeSelectorTerm.
// +structType=atomic
type NodeSelectorTerm struct {
	// A list of node selector requirements by node's labels.
	// +optional
	// +listType=atomic
	MatchExpressions []NodeSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,1,rep,name=matchExpressions"`
	// A list of node selector requirements by node's fields.
	// +optional
	// +listType=atomic
	MatchFields []NodeSelectorRequirement `json:"matchFields,omitempty" protobuf:"bytes,2,rep,name=matchFields"`
}

// A node selector requirement is a selector that contains values, a key, and an operator
// that relates the key and values.
type NodeSelectorRequirement struct {
	// The label key that the selector applies to.
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// Represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	Operator NodeSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=NodeSelectorOperator"`
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. If the operator is Gt or Lt, the values
	// array must have a single element, which will be interpreted as an integer.
	// This array is replaced during a strategic merge patch.
	// +optional
	// +listType=atomic
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// A node selector operator is the set of operators that can be used in
// a node selector requirement.
// +enum
type NodeSelectorOperator string

const (
	NodeSelectorOpIn           NodeSelectorOperator = "In"
	NodeSelectorOpNotIn        NodeSelectorOperator = "NotIn"
	NodeSelectorOpExists       NodeSelectorOperator = "Exists"
	NodeSelectorOpDoesNotExist NodeSelectorOperator = "DoesNotExist"
	NodeSelectorOpGt           NodeSelectorOperator = "Gt"
	NodeSelectorOpLt           NodeSelectorOperator = "Lt"
)

// RestartPolicy describes how the container should be restarted.
// Only one of the following restart policies may be specified.
// If none of the following policies is specified, the default one
// is RestartPolicyAlways.
// +enum
type RestartPolicy string

const (
	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"

	RestartPolicyRecreateAfterFailedRestart RestartPolicy = "RecreateAfterFailedRestart"
)

// A single application container that you want to run within a worker.
type Container struct {
	// Name of the container specified as a DNS_LABEL.
	// Each container in a worker must have a unique name (DNS_LABEL).
	// Cannot be updated.
	Name string `json:"name"`
	// Container image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +required
	Image string `json:"image,omitempty"`
	// Entrypoint array. Not executed within a shell.
	// The container image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// +listType=atomic
	Command []string `json:"command,omitempty"`
	// Arguments to the entrypoint.
	// The container image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// +listType=atomic
	Args []string `json:"args,omitempty"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Env []EnvVar `json:"env,omitempty"`
	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`

	// Pod volumes to mount into the container's filesystem.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=mountPath
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=mountPath
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`

	// Periodic probe of container service readiness.
	// Container will be removed from service endpoints if the probe fails.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`
	// Periodic probe of container liveness.
	// Container will be restarted if the probe fails.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	LivenessProbe *Probe `json:"livenessProbe,omitempty"`

	// List of ports to expose from the container. Not specifying a port here
	// DOES NOT prevent that port from being exposed. Any port which is
	// listening on the default "0.0.0.0" address inside a container will be
	// accessible from the network.
	Ports []string `json:"ports,omitempty"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Optional: no more than one of the following may be specified.

	// Variable references $(VAR_NAME) are expanded
	// using the previously defined environment variables in the container and
	// any service environment variables. If a variable cannot be resolved,
	// the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e.
	// "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)".
	// Escaped references will never be expanded, regardless of whether the variable
	// exists or not.
	// Defaults to "".
	// +optional
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	//ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	//ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	//ResourceStorage ResourceName = "storage"
	// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	//ResourceEphemeralStorage ResourceName = "ephemeral-storage"

	ResourceNumber ResourceName = "number"
)

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]resource.Quantity

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value. Requests cannot exceed Limits.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	CPURequests ResourceList `json:"cpuRequests,omitempty"`
	// +optional
	MemoryRequests ResourceList `json:"memoryRequests,omitempty"`
	// +optional
	GPURequests GPUResource `json:"gpuRequests,omitempty"`
	// +optional
	EphemeralStorageRequests ResourceList `json:"ephemeralStorageRequests,omitempty"`
}

type GPUResource struct {
	Models []string `json:"models,omitempty"`
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=8
	Number uint32 `json:"number"`
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	// This must match the Name of a Volume.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`

	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty" protobuf:"bytes,4,opt,name=subPath"`
}

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
type Probe struct {
	// The action taken to determine the health of a container
	ProbeHandler `json:",inline" protobuf:"bytes,1,opt,name=handler"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty" protobuf:"varint,2,opt,name=initialDelaySeconds"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum:=1
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty" protobuf:"varint,3,opt,name=timeoutSeconds"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +kubebuilder:default:=10
	// +kubebuilder:validation:Minimum:=1
	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty" protobuf:"varint,4,opt,name=periodSeconds"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum:=1
	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty" protobuf:"varint,5,opt,name=successThreshold"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +kubebuilder:default:=3
	// +kubebuilder:validation:Minimum:=1
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty" protobuf:"varint,6,opt,name=failureThreshold"`
}

// ProbeHandler defines a specific action that should be taken in a probe.
// One and only one of the fields must be specified.
type ProbeHandler struct {
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty" protobuf:"bytes,1,opt,name=httpGet"`
	// TCPSocket specifies an action involving a TCP port.
	// +optional
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty" protobuf:"bytes,2,opt,name=tcpSocket"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
type HTTPHeader struct {
	// The header field name.
	// This will be canonicalized upon output, so case-variant names will be understood as the same header.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// The header field value
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

// HTTPGetAction describes an action based on HTTP Get requests.
type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
	// Number of the port to access on the container.
	// Valid number must be in the range 1 to 65535.
	// When it's 0, use WorkerStatusAddress Port instead.
	// Defaults to WorkerStatusAddress Port.
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=65535
	// +optional
	Port int `json:"port" protobuf:"bytes,2,opt,name=port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional
	Host string `json:"host,omitempty" protobuf:"bytes,3,opt,name=host"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +kubebuilder:default:="HTTP"
	// +optional
	Scheme URIScheme `json:"scheme,omitempty" protobuf:"bytes,4,opt,name=scheme,casttype=URIScheme"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	// +listType=atomic
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty" protobuf:"bytes,5,rep,name=httpHeaders"`
}

// URIScheme identifies the scheme used for connection to a host for Get actions
// +enum
type URIScheme string

const (
	// URISchemeHTTP means that the scheme used will be http://
	URISchemeHTTP URIScheme = "HTTP"
	// URISchemeHTTPS means that the scheme used will be https://
	URISchemeHTTPS URIScheme = "HTTPS"
)

// TCPSocketAction describes an action based on opening a socket
type TCPSocketAction struct {
	// Number of the port to access on the container.
	// Valid number must be in the range 1 to 65535.
	// When it's 0, use WorkerStatusAddress Port instead.
	// Defaults to WorkerStatusAddress Port.
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=65535
	// +optional
	Port int `json:"port" protobuf:"bytes,1,opt,name=port"`
	// Optional: Host name to connect to, defaults to the WorkerStatusAddress ipOrDNSName.
	// +optional
	Host string `json:"host,omitempty" protobuf:"bytes,2,opt,name=host"`
}

type WorkerStatus struct {
	// +kubebuilder:default:="Pending"
	// +optional
	State WorkerState `json:"state"`

	// +optional
	RealCreateTime metav1.Time `json:"realCreateTime"`
	// +optional
	RealStartTime metav1.Time `json:"realStartTime"`

	// +optional
	ErrorMessage string `json:"errorMessage"`
	// +optional
	RealInstanceID string `json:"realInstanceID"`

	// +kubebuilder:default:=false
	// +optional
	Healthy bool `json:"healthy"`
	// +optional
	Addresses []WorkerStatusAddress `json:"addresses"`

	// Current service state of worker.
	// +optional
	Conditions []WorkerCondition `json:"conditions,omitempty"`

	// The list has one entry per container in the manifest.
	// +optional
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
}

type WorkerService struct {
	Port             uint16                        `json:"port"`
	MaxConcurrency   uint16                        `json:"maxConcurrency"`
	HealthCheckProbe WorkerServiceHealthCheckProbe `json:"healthCheckProbe"`
}

type WorkerServiceHealthCheckProbe struct {
	Method              string `json:"method"`
	Path                string `json:"path"`
	SuccessStatusCode   int    `json:"successStatusCode"`
	InitialDelaySeconds uint16 `json:"initialDelaySeconds"`
	PeriodSeconds       uint16 `json:"periodSeconds"`
	TimeoutSeconds      uint16 `json:"timeoutSeconds"`
	SuccessThreshold    uint16 `json:"successThreshold"`
}

type WorkerEnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WorkerVolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type WorkerResources struct {
	CPU    WorkerResourceCPU    `json:"cpu"`
	Memory WorkerResourceMemory `json:"memory"`
	GPU    WorkerResourceGPU    `json:"gpu"`
}

type WorkerResourceCPU struct {
	Number uint16 `json:"number"`
}

type WorkerResourceMemory struct {
	Size uint64 `json:"size"`
}

const (
	WorkerResourceGPUModelNvidia3090          = "nvidia/3090"
	WorkerResourceGPUModelNvidia4090          = "nvidia/4090"
	WorkerResourceGPUModelNvidia4090D         = "nvidia/4090D"
	WorkerResourceGPUModelNvidiaL20           = "nvidia/L20"
	WorkerResourceGPUModelNvidiaL40           = "nvidia/L40"
	WorkerResourceGPUModelNvidiaA6000         = "nvidia/A6000"
	WorkerResourceGPUModelNvidiaRTX6000Ada    = "nvidia/RTX6000Ada"
	WorkerResourceGPUModelNvidiaA100SXM4_80GB = "nvidia/A100_SXM4_80GB"
	WorkerResourceGPUModelNvidiaA100PCIE_80GB = "nvidia/A100_PCIE_80GB"
	WorkerResourceGPUModelNvidiaA100SXM4_40GB = "nvidia/A100_SXM4_40GB"
	WorkerResourceGPUModelNvidiaH800          = "nvidia/H800"
	WorkerResourceGPUModelNvidiaH20           = "nvidia/H20"
	WorkerResourceGPUModelNvidiaH100          = "nvidia/H100"
)

type WorkerResourceGPU struct {
	Number      uint16 `json:"number"`
	Model       string `json:"model"`
	CUDAVersion string `json:"cudaVersion"`
}

type WorkerState string

const (
	WorkerStatePending  WorkerState = "Pending"
	WorkerStateCreating WorkerState = "Creating"
	WorkerStateError    WorkerState = "Error"

	WorkerStateRunning  WorkerState = "Running"
	WorkerStatePaused   WorkerState = "Paused"
	WorkerStateUnpaused WorkerState = "Unpaused"
	WorkerStateStopped  WorkerState = "Stopped"

	WorkerStateFailed WorkerState = "Failed"

	WorkerStateUnScheduling WorkerState = "UnScheduling"
	WorkerStateTerminating  WorkerState = "Terminating"

	WorkerStateUnknown WorkerState = "Unknown"
)

type WorkerStatusAddress struct {
	IPOrDNSName string `json:"ipOrDNSName"`
	// +optional
	Port uint16 `json:"port"`
	// +optional
	HTTPHeaders map[string]string `json:"httpHeaders"`
}

// WorkerConditionType is a valid value for WorkerCondition.Type
type WorkerConditionType string

// These are built-in conditions of worker. An application may use a custom condition not listed here.
const (
	// ContainersReady indicates whether all containers in the worker are ready.
	ContainersReady WorkerConditionType = "ContainersReady"
	// WorkerInitialized means that all init containers in the worker have started successfully.
	WorkerInitialized WorkerConditionType = "Initialized"
	// WorkerReady means the worker is able to service requests and should be added to the
	// load balancing pools of all matching services.
	WorkerReady WorkerConditionType = "Ready"
	// WorkerScheduled represents status of the scheduling process for this worker.
	WorkerScheduled WorkerConditionType = "WorkerScheduled"
	// LivenessProbeSucceeded indicates whether all containers' LivenessProbe in the worker succeed
	LivenessProbeSucceeded WorkerConditionType = "LivenessProbeSucceeded"

	WorkerRealStatus WorkerConditionType = "WorkerRealStatus"
)

const (
	WorkerRealAvailable   WorkerConditionType = "WorkerRealAvailable"
	WorkerRealUnavailable WorkerConditionType = "WorkerRealUnavailable"
)

// These are reasons for a worker's transition to a condition.
const (
	// WorkerReasonUnschedulable reason in WorkerScheduled WorkerCondition means that the scheduler
	// can't schedule the worker right now, for example due to insufficient resources in the cluster.
	WorkerReasonUnschedulable = "Unschedulable"

	// WorkerReasonSchedulingFailure reason in WorkerScheduled WorkerCondition means that the scheduler
	// skips scheduling the worker because something is wrong with this worker, for example due to the wrong config.
	WorkerReasonSchedulingFailure = "SchedulingFailure"

	// WorkerReasonSchedulerError reason in WorkerScheduled WorkerCondition means that some internal error happens
	// during scheduling, for example due to nodeAffinity parsing errors.
	WorkerReasonSchedulerError = "SchedulerError"
)

// WorkerCondition contains details for the current condition of this worker.
type WorkerCondition struct {
	// Type is the type of the condition.
	Type WorkerConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus `json:"status"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// ContainerStateWaiting is a waiting state of a container.
type ContainerStateWaiting struct {
	// (brief) reason the container is not yet running.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,1,opt,name=reason"`
	// Message regarding why the container is not yet running.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
}

// ContainerStateRunning is a running state of a container.
type ContainerStateRunning struct {
	// Time at which the container was last (re-)started
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,1,opt,name=startedAt"`
}

// ContainerStateTerminated is a terminated state of a container.
type ContainerStateTerminated struct {
	// Exit status from the last termination of the container
	ExitCode int32 `json:"exitCode" protobuf:"varint,1,opt,name=exitCode"`
	// Signal from the last termination of the container
	// +optional
	Signal int32 `json:"signal,omitempty" protobuf:"varint,2,opt,name=signal"`
	// (brief) reason from the last termination of the container
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// Message regarding the last termination of the container
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// Time at which previous execution of the container started
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,5,opt,name=startedAt"`
	// Time at which the container last terminated
	// +optional
	FinishedAt metav1.Time `json:"finishedAt,omitempty" protobuf:"bytes,6,opt,name=finishedAt"`
	// Container's ID in the format '<type>://<container_id>'
	// +optional
	ContainerID string `json:"containerID,omitempty" protobuf:"bytes,7,opt,name=containerID"`
}

// ContainerState holds a possible state of container.
// Only one of its members may be specified.
// If none of them is specified, the default one is ContainerStateWaiting.
type ContainerState struct {
	// Details about a waiting container
	// +optional
	Waiting *ContainerStateWaiting `json:"waiting,omitempty"`
	// Details about a running container
	// +optional
	Running *ContainerStateRunning `json:"running,omitempty"`
	// Details about a terminated container
	// +optional
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
	// Details about a paused container
	// +optional
	Paused *ContainerStateWaiting `json:"paused,omitempty"`
	// Details about a stopped container
	// +optional
	Stopped *ContainerStateWaiting `json:"stopped,omitempty"`
}

// ContainerStatus contains details for the current status of this container.
type ContainerStatus struct {
	// Name is a DNS_LABEL representing the unique name of the container.
	// Each container in a pod must have a unique name across all container types.
	// Cannot be updated.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// State holds details about the container's current condition.
	// +optional
	State ContainerState `json:"state,omitempty" protobuf:"bytes,2,opt,name=state"`
	// LastTerminationState holds the last termination state of the container to
	// help debug container crashes and restarts. This field is not
	// populated if the container is still running and RestartCount is 0.
	// +optional
	LastTerminationState ContainerState `json:"lastState,omitempty" protobuf:"bytes,3,opt,name=lastState"`
	// Ready specifies whether the container is currently passing its readiness check.
	// The value will change as readiness probes keep executing. If no readiness
	// probes are specified, this field defaults to true once the container is
	// fully started (see Started field).
	//
	// The value is typically used to determine whether a container is ready to
	// accept traffic.
	Ready bool `json:"ready" protobuf:"varint,4,opt,name=ready"`
	// RestartCount holds the number of times the container has been restarted.
	// Kubelet makes an effort to always increment the value, but there
	// are cases when the state may be lost due to node restarts and then the value
	// may be reset to 0. The value is never negative.
	RestartCount int32 `json:"restartCount" protobuf:"varint,5,opt,name=restartCount"`
	// Image is the name of container image that the container is running.
	// The container image may not match the image used in the PodSpec,
	// as it may have been resolved by the runtime.
	// More info: https://kubernetes.io/docs/concepts/containers/images.
	Image string `json:"image" protobuf:"bytes,6,opt,name=image"`
	// ImageID is the image ID of the container's image. The image ID may not
	// match the image ID of the image used in the PodSpec, as it may have been
	// resolved by the runtime.
	ImageID string `json:"imageID" protobuf:"bytes,7,opt,name=imageID"`
	// ContainerID is the ID of the container in the format '<type>://<container_id>'.
	// Where type is a container runtime identifier, returned from Version call of CRI API
	// (for example "containerd").
	// +optional
	ContainerID string `json:"containerID,omitempty" protobuf:"bytes,8,opt,name=containerID"`
}
