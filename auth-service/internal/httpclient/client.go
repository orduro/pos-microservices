package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient              *http.Client
	communicationServiceURL string
}

type Config struct {
	CommunicationServiceURL string
	Timeout                 time.Duration
}

func NewHttpClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		httpClient:              client,
		communicationServiceURL: config.CommunicationServiceURL,
	}
}

// represents a request to any service
type ServiceRequest struct {
	Method  string
	Path    string
	Body    any
	Headers map[string]string
}

// represents the response from any service
type ServiceResponse struct {
	StatusCode int
	Body       string
	Headers    http.Header
}

func (c *Client) CallService(ctx context.Context, req ServiceRequest) (*ServiceResponse, error) {
	url := fmt.Sprintf("%s%s", c.communicationServiceURL, req.Path)

	var reqBody io.Reader
	if req.Body != nil {
		jsonData, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal json: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	// create request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	// set default headers
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// perform request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	return &ServiceResponse{
		StatusCode: resp.StatusCode,
		Body:       string(responseBody),
		Headers:    resp.Header,
	}, nil
}
