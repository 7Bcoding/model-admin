package models

// ModelAPI represents the configuration for a model API deployment
type ModelAPI struct {
	Name      string            `yaml:"name" json:"name"`
	Region    string            `yaml:"region" json:"region"`
	Image     string            `yaml:"image" json:"image"`
	Args      []string          `yaml:"args" json:"args"`
	Env       map[string]string `yaml:"env" json:"env"`
	Resources struct {
		GPU struct {
			Model string `yaml:"model" json:"model"`
			Count int    `yaml:"count" json:"count"`
		} `yaml:"gpu" json:"gpu"`
		Storage int `yaml:"storage" json:"storage"`
	} `yaml:"resources" json:"resources"`
	ScalingPolicy struct {
		MinReplicas int `yaml:"min_replicas" json:"min_replicas"`
		MaxReplicas int `yaml:"max_replicas" json:"max_replicas"`
		Concurrency struct {
			Target int `yaml:"target" json:"target"`
		} `yaml:"concurrency" json:"concurrency"`
		ScaleUpWindowSeconds   int `yaml:"scale_up_window_seconds" json:"scale_up_window_seconds"`
		ScaleDownWindowSeconds int `yaml:"scale_down_window_seconds" json:"scale_down_window_seconds"`
	} `yaml:"scaling_policy" json:"scaling_policy"`
	MaxConcurrentPerWorker int `yaml:"max_concurrent_per_worker" json:"max_concurrent_per_worker"`
}

// Validate performs validation on the ModelAPI configuration
func (m *ModelAPI) Validate() error {
	// TODO: Add validation logic
	return nil
}
