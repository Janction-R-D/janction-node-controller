package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"janction/dao/postgres"
	"janction/model"
	"janction/setting"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	jobStore   = make(map[string][]string)
	jobStoreMu sync.Mutex
)

// GetJob retrieves a job for the given nodeID
func GetJob(nodeID string) (*string, error) {
	jobStoreMu.Lock()
	defer jobStoreMu.Unlock()

	// Check if there are jobs available for the given nodeID
	if jobs, exists := jobStore[nodeID]; exists && len(jobs) > 0 {
		// Pop the job from the head of the array
		job := jobs[0]
		jobStore[nodeID] = jobs[1:]
		return &job, nil
	}

	// If jobs are less than 100, fetch new jobs from the backend
	if jobs, exists := jobStore[nodeID]; !exists || len(jobs) < 100 {
		nodeRegistration, err := postgres.GetNodeRegistrationByNodeID(nodeID);
		if err != nil {
			return nil, err
		}
		params := model.FormGetJobType{
			Architecture: nodeRegistration.ArchitectureType,
			UseCPU: nodeRegistration.UseCPU,
			UseGPU: nodeRegistration.UseGPU,
		}
		newJob, err := fetchJobType(params)
		if err != nil {
			return nil, err
		}

		// Append the new jobs to the existing jobs
		jobStore[nodeID] = append(jobStore[nodeID], *newJob)
	}

	// Check again if jobs are available after fetching
	if jobs, exists := jobStore[nodeID]; exists && len(jobs) > 0 {
		job := jobs[0]
		jobStore[nodeID] = jobs[1:]
		return &job, nil
	}

	return nil, errors.New("no jobs available")
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
