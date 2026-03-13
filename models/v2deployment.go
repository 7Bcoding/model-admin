package models

import (
	"gitlab.paigod.work/ai-cloud/nebula/api/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageCache struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Size   int    `json:"size,omitempty"`
}

type ModelCache struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Size   int    `json:"size,omitempty"`
}

// CNodeListItem represents a compute node list item with kubectl output format
type CNodeListItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Fields matching kubectl output
	Name       string  `json:"name,omitempty"`       // NAME
	Provider   string  `json:"provider,omitempty"`   // PROVIDER
	Region     string  `json:"region,omitempty"`     // REGION
	HostName   string  `json:"hostName,omitempty"`   // HOSTNAME
	CPUProduct string  `json:"cpuProduct,omitempty"` // CPUPRODUCT
	GPUProduct string  `json:"gpuProduct,omitempty"` // GPUPRODUCT
	Policy     string  `json:"policy,omitempty"`     // POLICY
	Status     string  `json:"status,omitempty"`     // STATUS
	GPUState   []int32 `json:"gpuState,omitempty"`   // GPUSTATE
	Age        string  `json:"age,omitempty"`        // AGE

	ImageCache []*ImageCache `json:"imageCache,omitempty"`
	ModelCache []*ModelCache `json:"modelCache,omitempty"`
}

type NDeploymentItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Hot               int32  `json:"hot,omitempty"`
	HotReady          int32  `json:"hotReady,omitempty"`
	Cold              int32  `json:"cold,omitempty"`
	ColdReady         int32  `json:"coldReady,omitempty"`
	Status            string `json:"status,omitempty"`

	CustomerRegion string            `json:"customerRegion,omitempty"`
	Resources      *Resources        `json:"resources"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	// Affinity       *Affinity            `json:"affinity,omitempty"`
	Containers []*v1beta1.Container `json:"containers,omitempty"`

	Volumes []*v1beta1.Volume `json:"volumes,omitempty"` // volumes

	Image string `json:"image,omitempty"`

	Age          string     `json:"age,omitempty"`
	MaxBatchSize int32      `json:"maxBatchSize,omitempty"`
	SimpleNSP    *SimpleNSP `json:"simpleNSP,omitempty"`
}

type SimpleNSP struct {
	MinReplicas int32 `json:"minReplicas,omitempty"`
	MaxReplicas int32 `json:"maxReplicas,omitempty"`

	BatchPerWorker   int32 `json:"batchPerWorker,omitempty"`
	MaxScaleDownRate int32 `json:"maxScaleDownRate,omitempty"`
}

// NebulaWorkerItem represents a nebula worker list item with kubectl output format
type NebulaWorkerItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Fields matching kubectl output
	Name       string `json:"name,omitempty"`       // NAME
	Phase      string `json:"phase,omitempty"`      // PHASE
	GpuProduct string `json:"gpuProduct,omitempty"` // GPUPRODUCT
	Gpu        []int  `json:"gpu,omitempty"`        // GPU
	Node       string `json:"node,omitempty"`       // NODE
	Status     string `json:"status,omitempty"`     // STATUS
	Age        string `json:"age,omitempty"`        // AGE

	Deployment string `json:"deployment,omitempty"`
}
