Below is the Go source code to consume the provided APIs:

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiKey = "<replace_api_key>"
	baseURL = "https://api.on-demand.io/chat/v1"
)

type CreateSessionResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func main() {
	// Replace with your external user ID
	externalUserId := "<replace_external_user_id>"

	// Step 1: Create Chat Session
	sessionID, err := createChatSession(externalUserId)
	if err != nil {
		fmt.Printf("Error creating chat session: %v\n", err)
		return
	}

	fmt.Printf("Chat session created successfully. Session ID: %s\n", sessionID)

	// Step 2: Submit Query
	query := "Put your query here"
	err = submitQuery(sessionID, query)
	if err != nil {
		fmt.Printf("Error submitting query: %v\n", err)
		return
	}

	fmt.Println("Query submitted successfully.")
}

func createChatSession(externalUserId string) (string, error) {
	url := fmt.Sprintf("%s/sessions", baseURL)

	body := map[string]interface{}{
		"pluginIds":       []string{},
		"externalUserId":  externalUserId,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Data.ID, nil
}

func submitQuery(sessionID, query string) error {
	url := fmt.Sprintf("%s/sessions/%s/query", baseURL, sessionID)

	body := map[string]interface{}{
		"endpointId":    "predefined-openai-gpt4o",
		"query":         query,
		"pluginIds":     []string{"plugin-1741871229"},
		"responseMode":  "sync",
		"reasoningMode": "medium",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response: %s\n", string(responseBody))
	return nil
}
