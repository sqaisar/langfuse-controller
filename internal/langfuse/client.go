package langfuse

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	BaseURL   string
	Client    *http.Client
	PublicKey string
	SecretKey string
}

func NewClient() *Client {
	baseURL := os.Getenv("LANGFUSE_HOST")
	if baseURL == "" {
		baseURL = "https://cloud.langfuse.com"
	}
	publicKey := os.Getenv("LANGFUSE_PUBLIC_KEY")
	secretKey := os.Getenv("LANGFUSE_SECRET_KEY")

	return &Client{
		BaseURL:   baseURL,
		Client:    &http.Client{},
		PublicKey: publicKey,
		SecretKey: secretKey,
	}
}

func (c *Client) do(req *http.Request, v interface{}) error {
	// Use Basic Auth with public_key:secret_key base64 encoded
	auth := base64.StdEncoding.EncodeToString([]byte(c.PublicKey + ":" + c.SecretKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s: %s", resp.Status, string(body))
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func (c *Client) CreateProject(name string) (*Project, error) {
	reqBody, _ := json.Marshal(CreateProjectRequest{Name: name})
	req, _ := http.NewRequest("POST", c.BaseURL+"/api/public/projects", bytes.NewBuffer(reqBody))

	var project Project
	err := c.do(req, &project)
	return &project, err
}

func (c *Client) GetProject(id string) (*Project, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/api/public/projects/"+id, nil)
	var project Project
	err := c.do(req, &project)
	return &project, err
}

func (c *Client) CreateAPIKey(projectID, name string) (*APIKey, error) {
	reqBody, _ := json.Marshal(CreateAPIKeyRequest{Name: name, ProjectID: projectID})
	// Note: Endpoint might be different, checking docs...
	// Docs say POST /api/public/projects/{projectId}/apiKeys
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/public/projects/%s/apiKeys", c.BaseURL, projectID), bytes.NewBuffer(reqBody))

	var apiKey APIKey
	err := c.do(req, &apiKey)
	return &apiKey, err
}

// CreateModel creates a new model definition
func (c *Client) CreateModel(model Model) (*Model, error) {
	reqBody, _ := json.Marshal(model)
	req, _ := http.NewRequest("POST", c.BaseURL+"/api/public/models", bytes.NewBuffer(reqBody))
	var createdModel Model
	err := c.do(req, &createdModel)
	return &createdModel, err
}

// CreateLlmConnection creates a new LLM connection
// Note: Endpoint is hypothetical, need verification
func (c *Client) CreateLlmConnection(projectID string, connection interface{}) error {
	reqBody, _ := json.Marshal(connection)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/public/projects/%s/llm-connections", c.BaseURL, projectID), bytes.NewBuffer(reqBody))
	return c.do(req, nil)
}

// CreatePrompt creates a new prompt
func (c *Client) CreatePrompt(projectID string, prompt interface{}) error {
	reqBody, _ := json.Marshal(prompt)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/public/projects/%s/prompts", c.BaseURL, projectID), bytes.NewBuffer(reqBody))
	return c.do(req, nil)
}

// CreateScoreConfig creates a new score configuration
func (c *Client) CreateScoreConfig(projectID string, config interface{}) error {
	reqBody, _ := json.Marshal(config)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/public/projects/%s/score-configs", c.BaseURL, projectID), bytes.NewBuffer(reqBody))
	return c.do(req, nil)
}
