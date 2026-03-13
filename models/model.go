package models

type Model struct {
	ID          int    `json:"id"`
	ModelName   string `json:"model_name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type ModelStar struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ModelName string `json:"model_name"`
	Note      string `json:"note"`
}
