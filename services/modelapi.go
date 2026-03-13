package services

import (
	"llm-ops/models"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type ModelAPIService struct {
	db *gorm.DB
}

func NewModelAPIService(db *gorm.DB) *ModelAPIService {
	return &ModelAPIService{db: db}
}

// ApplyYAML parses and applies the YAML configuration
func (s *ModelAPIService) ApplyYAML(yamlContent []byte) error {
	var modelAPI models.ModelAPI
	if err := yaml.Unmarshal(yamlContent, &modelAPI); err != nil {
		return err
	}

	// Validate the configuration
	if err := modelAPI.Validate(); err != nil {
		return err
	}

	// Check if model API already exists
	var existing models.ModelAPI
	result := s.db.Where("name = ?", modelAPI.Name).First(&existing)

	if result.Error == nil {
		// Update existing model API
		return s.db.Model(&existing).Updates(modelAPI).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// Create new model API
		return s.db.Create(&modelAPI).Error
	}

	return result.Error
}
