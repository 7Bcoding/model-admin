package models

// EndpointInfo 表示 ServerlessEndpoint 的关键信息
type EndpointInfo struct {
	APIVersion     string           `json:"apiVersion"`
	Kind           string           `json:"kind"`
	Metadata       EndpointMetadata `json:"metadata"`
	Spec           EndpointSpec     `json:"spec"`
	WorkerCount    int              `json:"workerCount"`
	Yaml           string           `json:"yaml"`
	SapYaml        string           `json:"sapYaml"`
	SapParam       *SapParam        `json:"sapParam"`
	Models         []string         `json:"models"`
	Image          string           `json:"image"`
	MaxConcurrency int              `json:"maxConcurrency"`
}

type EndpointMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type EndpointSpec struct {
	ClusterIDs           []string         `json:"clusterIDs"`
	ConcurrencyPerWorker float64          `json:"concurrencyPerWorker"`
	FastBoot             bool             `json:"fastBoot"`
	Provider             string           `json:"provider"`
	Replicas             float64          `json:"replicas"`
	Template             EndpointTemplate `json:"template"`
}

type EndpointTemplate struct {
	Spec EndpointTemplateSpec `json:"spec"`
}

type EndpointTemplateSpec struct {
	Containers []Container `json:"containers"`
	Affinity   Affinity    `json:"affinity,omitempty"`
}

// Affinity 是一个简化的亲和性定义
type Affinity struct {
	NodeAffinity NodeAffinity `json:"nodeAffinity"`
}

type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution NodeSelector `json:"requiredDuringSchedulingIgnoredDuringExecution"`
}

type NodeSelector struct {
	NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms"`
}

type NodeSelectorTerm struct {
	MatchExpressions []MatchExpression `json:"matchExpressions"`
}

type MatchExpression struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// EnvVar 环境变量
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
