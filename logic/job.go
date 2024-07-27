package logic

import (
	"encoding/json"
	"fmt"
	"janction/model"
	"janction/setting"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	jobStore   = make([]string, 0, 100)
	jobStoreMu sync.Mutex
)

// GetJob retrieves a job
func GetJob(nodeID string) (*string, error) {
	jobStoreMu.Lock()
	defer jobStoreMu.Unlock()

	// Check if there are jobs available
	if len(jobStore) > 0 {
		// Pop the job from the head of the array
		job := jobStore[0]
		jobStore = jobStore[1:]
		return &job, nil
	}

	return nil, nil
}

// FetchJob fetches jobs from the backend and fills the jobStore
func FetchJob() error {
	if len(jobStore) == 100 {
		return nil
	}
	
	params := model.FormGetJobType{
		Architecture: "amd64",
		UseCPU:       1,
		UseGPU:       0,
	}
	newJob, err := fetchJobType(params)
	if err != nil {
		return err
	}

	jobStoreMu.Lock()
	defer jobStoreMu.Unlock()

	if len(jobStore) < 100 {
		jobStore = append(jobStore, *newJob)
	} 

	return nil
}

// fetchJobType sends a request to the backend to fetch job types
func fetchJobType(params model.FormGetJobType) (*string, error) {
	baseURL := setting.Config.UrlConfig.JanctionBackend

	// Construct the request URL with query parameters
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		zap.L().Error("Failed to parse base URL", zap.Error(err))
		return nil, err
	}

	query := reqURL.Query()
	query.Set("architecture", params.Architecture)
	query.Set("use_gpu", fmt.Sprintf("%d", params.UseGPU))
	query.Set("use_cpu", fmt.Sprintf("%d", params.UseCPU))
	reqURL.RawQuery = query.Encode()

	// Create a new HTTP client and set a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the GET request
	resp, err := client.Get(reqURL.String())
	if err != nil {
		zap.L().Error("Failed to fetch job type", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Code != 1000 {
		return nil, fmt.Errorf("request failed: %s", response.Msg)
	}

	jobType := response.Data.(string)

	return &jobType, nil
}
