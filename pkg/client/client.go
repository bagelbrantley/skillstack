package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bagelbrantley/skillstack/pkg/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type VerifyCodeResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	User    struct {
		UID   string `json:"uid"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
	Error string `json:"error,omitempty"`
}

func (c *Client) VerifyCode(code string) (*VerifyCodeResponse, error) {
	reqBody, _ := json.Marshal(map[string]string{"code": code})
	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/cli/verify-code", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("connection to SkillStack server failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &errResp)
		if errResp.Error != "" {
			return nil, fmt.Errorf("%s", errResp.Error)
		}
		return nil, fmt.Errorf("server returned HTTP %d", resp.StatusCode)
	}

	var result VerifyCodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid server response: %w", err)
	}
	return &result, nil
}

type Lesson struct {
	ID              string `json:"id"`
	Number          string `json:"number"`
	Title           string `json:"title"`
	CheckpointTitle string `json:"checkpointTitle"`
	CLICommand      string `json:"cliCommand"`
	Instructions    string `json:"instructions"`
	ExpectedOutput  string `json:"expectedOutput"`
}

type Chapter struct {
	ID      string   `json:"id"`
	Number  int      `json:"number"`
	Title   string   `json:"title"`
	Lessons []Lesson `json:"lessons"`
}

type GuidedProject struct {
	ID               string    `json:"id"`
	Slug             string    `json:"slug"`
	Title            string    `json:"title"`
	Level            string    `json:"level"`
	Duration         string    `json:"duration"`
	Language         string    `json:"language"`
	Summary          string    `json:"summary"`
	CheckpointsCount int       `json:"checkpointsCount"`
	Chapters         []Chapter `json:"chapters"`
}

type ProjectsResponse struct {
	Projects []GuidedProject `json:"projects"`
}

func (c *Client) GetProjects() ([]GuidedProject, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/cli/projects")
	if err != nil {
		return nil, fmt.Errorf("failed to reach SkillStack server: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res ProjectsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}
	return res.Projects, nil
}

func (c *Client) GetProject(id string) (*GuidedProject, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/cli/projects/" + id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project '%s' not found", id)
	}

	var proj GuidedProject
	if err := json.Unmarshal(body, &proj); err != nil {
		return nil, fmt.Errorf("failed to parse project details: %w", err)
	}
	return &proj, nil
}

type SubmissionPayload struct {
	ProjectID        string `json:"projectId"`
	LessonID         string `json:"lessonId"`
	CheckpointNumber string `json:"checkpointNumber"`
	ExitCode         int    `json:"exitCode"`
	Output           string `json:"output"`
}

type SubmissionResponse struct {
	Success          bool   `json:"success"`
	Verified         bool   `json:"verified"`
	User             string `json:"user"`
	ProjectID        string `json:"projectId"`
	LessonID         string `json:"lessonId"`
	CheckpointNumber string `json:"checkpointNumber"`
	XPAwarded        int    `json:"xpAwarded"`
	Message          string `json:"message"`
	Timestamp        string `json:"timestamp"`
	Error            string `json:"error,omitempty"`
}

func (c *Client) SubmitCheckpoint(payload SubmissionPayload) (*SubmissionResponse, error) {
	creds, err := config.LoadCredentials()
	if err != nil || creds == nil || creds.Token == "" {
		return nil, fmt.Errorf("unauthenticated. Please run `skillstack login` first")
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", c.BaseURL+"/api/cli/submit", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+creds.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send submission to server: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("session expired or invalid token. Please log in again using `skillstack login`")
	}

	var subResp SubmissionResponse
	if err := json.Unmarshal(body, &subResp); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %w", err)
	}

	return &subResp, nil
}
