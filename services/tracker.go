package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"llm-ops/utils"
	"log"
	"net/http"
	"time"
)

type TrackerService struct {
	baseURL string
	client  *http.Client
}

type RegionFilter struct {
	RegionID string   `json:"regionID"`
	Servers  []string `json:"servers,omitempty"`
}

type ModelInfo struct {
	ModelName         string `json:"modelName"`
	ModelStatus       string `json:"modelStatus"`
	TargetModelStatus string `json:"targetModelStatus,omitempty"`
}

type DiskInfo struct {
	DiskID     string `json:"diskID"`
	MountPoint string `json:"mountPoint"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}

type ServerInfo struct {
	ServerID string      `json:"serverID"`
	Disks    []DiskInfo  `json:"disks"`
	IPs      []string    `json:"ip"`
	Models   []ModelInfo `json:"models"`
}

type RegionInfo struct {
	RegionID string       `json:"regionID"`
	Servers  []ServerInfo `json:"servers"`
}

type SearchResponse struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	ResponseID string `json:"responseId"`
	Data       struct {
		Regions []RegionInfo `json:"regions"`
	} `json:"data"`
}

type ModelRequest struct {
	ModelName string         `json:"modelName"`
	HFToken   string         `json:"hfToken,omitempty"`
	Regions   []RegionFilter `json:"regions,omitempty"`
}

type BaseResponse struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	ResponseID string `json:"responseId"`
}

func NewTrackerService(baseURL string) *TrackerService {
	return &TrackerService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

//add auth header

func (s *TrackerService) SearchModels(filters []RegionFilter) (*SearchResponse, error) {
	reqBody := struct {
		Filters []RegionFilter `json:"filters"`
	}{
		Filters: filters,
	}

	resp := &SearchResponse{}
	err := s.doRequest("POST", "/v2/models/search", reqBody, resp)
	if err != nil {
		return nil, fmt.Errorf("search models failed: %w", err)
	}

	return resp, nil
}

func (s *TrackerService) WarmupModels(models []ModelRequest, hfToken string, _ []string) (*BaseResponse, error) {
	reqBody := struct {
		Models []ModelRequest `json:"models"`
	}{
		Models: models,
	}

	if hfToken == "" {
		hfToken = "xxxx"
	}
	// 如果hfToken不为空，则设置HFToken
	for i := range reqBody.Models {
		reqBody.Models[i].HFToken = hfToken
	}

	resp := &BaseResponse{}
	err := s.doRequest("POST", "/v2/models", reqBody, resp)
	if err != nil {
		return nil, fmt.Errorf("warmup models failed: %w", err)
	}

	return resp, nil
}

func (s *TrackerService) DeleteModels(models []ModelRequest) (*BaseResponse, error) {
	reqBody := struct {
		Models []ModelRequest `json:"models"`
	}{
		Models: models,
	}

	resp := &BaseResponse{}
	err := s.doRequest("DELETE", "/v2/models", reqBody, resp)
	if err != nil {
		return nil, fmt.Errorf("delete models failed: %w", err)
	}

	return resp, nil
}

func (s *TrackerService) doRequest(method, path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest(method, s.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 添加认证头
	authHeaders := utils.GenerateAuthHeaders()
	req.Header.Set("sign", authHeaders.Sign)
	req.Header.Set("nonce", authHeaders.Nonce)
	req.Header.Set("timestamp", authHeaders.Timestamp)

	// print log of request, header and body
	// log.Printf("Request: %+v", req)
	// log.Printf("Headers: %+v", req.Header)
	log.Printf("Body: %+v", string(jsonBody))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	log.Printf("Tracker Response: %+v", resp)

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}

	return nil
}
