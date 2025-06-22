package ai_client

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
	url        string // URL of the AI service
	httpClient *http.Client
}

func NewClient(url string, timeout int) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second, // Set a reasonable timeout

		},
	}
}

func (c *Client) ProcessQuery(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	requestURL := fmt.Sprintf("%s/api/v1/azure-openai/process-query", c.url)

	if request.UserID == "" || request.MessageID == "" {
		return nil, fmt.Errorf("internal error no userID or messageID provided") // Return an empty response if UserID or MessageID is missing
	}

	request.Query = append(request.Query, defaultPrompts...)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response ChatResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
